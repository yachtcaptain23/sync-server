package command

import (
	"fmt"
	"time"

	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/brave-experiments/sync-server/utils"
	"github.com/satori/go.uuid"
)

const (
	// Always return hard-coded value of store birthday for now
	STORE_BIRTHDAY                string = "1"
	MAX_COMMIT_BATCH_SIZE         int32  = 90
	SESSIONS_COMMIT_DELAY_SECONDS int32  = 11
	SET_SYNC_POLL_INTERVAL        int32  = 30
)

func HandleGetUpdateRequest(guMsg *sync_pb.GetUpdatesMessage, guRsp *sync_pb.GetUpdatesResponse, pg *datastore.Postgres) (*sync_pb.SyncEnums_ErrorType, error) {
	fmt.Println("GET UPDATE RECEIVED")
	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later
	// TODO: process FetchFolders

	// Process from_progress_marker
	// TODO: query DB to get update entries
	if guMsg.FromProgressMarker != nil {
		guRsp.NewProgressMarker = make([]*sync_pb.DataTypeProgressMarker, len(guMsg.FromProgressMarker))

		for i := 0; i < len(guMsg.FromProgressMarker); i++ {
			guRsp.NewProgressMarker[i] = &sync_pb.DataTypeProgressMarker{}
			guRsp.NewProgressMarker[i].DataTypeId = guMsg.FromProgressMarker[i].DataTypeId
			// TODO: latest timestamp of records instead?
			guRsp.NewProgressMarker[i].Token, _ = time.Now().MarshalJSON()
		}
	}

	if *guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_NEW_CLIENT {
		fmt.Println("New client")

		// Create nigori top folder
		guRsp.Entries = make([]*sync_pb.SyncEntity, 1)
		now := time.Now().Unix()
		deleted := false
		folder := true
		name := "Nigori"
		serverDefinedTag := "google_chrome_nigori"
		version := int64(1)
		parentIdString := "0"
		idString := uuid.NewV4().String()

		nigoriSpecific := &sync_pb.NigoriSpecifics{}
		specific := &sync_pb.EntitySpecifics_Nigori{Nigori: nigoriSpecific}
		specifics := &sync_pb.EntitySpecifics{SpecificsVariant: specific}
		syncEntry := &sync_pb.SyncEntity{
			Ctime: &now, Mtime: &now, Deleted: &deleted, Folder: &folder,
			Name: &name, ServerDefinedUniqueTag: &serverDefinedTag,
			Version: &version, ParentIdString: &parentIdString,
			IdString: &idString, Specifics: specifics}
		entityToCommit, err := datastore.CreateSyncEntity(syncEntry, "")
		if err != nil {
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		} else {
			err = pg.InsertSyncEntity(entityToCommit)
			if err != nil {
				errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			} else {
				guRsp.Entries = make([]*sync_pb.SyncEntity, 1)
				guRsp.Entries[0] = syncEntry
			}
		}
		// Bypassing chromium's restriction here, our server won't provide the
		// initial encryption keys like chromium does, this will be overwritten
		// by our client.
		guRsp.EncryptionKeys = make([][]byte, 1)
		guRsp.EncryptionKeys[0] = []byte("1234")
	}

	// TODO: Implement batch reply and update the value accordingly
	changesRemaining := int64(0)
	guRsp.ChangesRemaining = &changesRemaining
	return &errCode, nil
}

func HandleCommitRequest(commitMsg *sync_pb.CommitMessage, commitRsp *sync_pb.CommitResponse, pg *datastore.Postgres) (*sync_pb.SyncEnums_ErrorType, error) {
	fmt.Println("COMMIT RECEIVED")
	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later

	commitRsp.Entryresponse = make([]*sync_pb.CommitResponse_EntryResponse, len(commitMsg.Entries))
	if commitMsg.Entries != nil {
		for i, v := range commitMsg.Entries {
			// TODO: Verified fields here before processing the values, early
			// return bad message when any required fields are invalid.
			entryRsp := &sync_pb.CommitResponse_EntryResponse{}
			commitRsp.Entryresponse[i] = entryRsp
			entityToCommit, err := datastore.CreateSyncEntity(v, *commitMsg.CacheGuid)
			if err != nil {
				rspType := sync_pb.CommitResponse_INVALID_MESSAGE
				entryRsp.ResponseType = &rspType
				continue
			}

			if *v.Version == 0 { // Create
				entityToCommit.Version++
				err = pg.InsertSyncEntity(entityToCommit)
				if err != nil {
					rspType := sync_pb.CommitResponse_INVALID_MESSAGE
					entryRsp.ResponseType = &rspType
					continue
				}
			} else { // Update
				match, err := pg.CheckVersion(entityToCommit.ID, entityToCommit.Version)
				if err != nil {
					rspType := sync_pb.CommitResponse_INVALID_MESSAGE
					entryRsp.ResponseType = &rspType
					continue
				}
				if !match {
					rspType := sync_pb.CommitResponse_CONFLICT
					entryRsp.ResponseType = &rspType
					continue
				}
				entityToCommit.Version++
				err = pg.UpdateSyncEntity(entityToCommit)
				if err != nil {
					rspType := sync_pb.CommitResponse_INVALID_MESSAGE
					entryRsp.ResponseType = &rspType
					continue
				}
			}

			// Prepare success response
			rspType := sync_pb.CommitResponse_SUCCESS
			entryRsp.ResponseType = &rspType
			entryRsp.IdString = utils.String(entityToCommit.ID)
			if entityToCommit.ParentID.Valid {
				entryRsp.ParentIdString = utils.String(entityToCommit.ParentID.String)
			}
			entryRsp.Version = &entityToCommit.Version
			if entityToCommit.Name.Valid {
				entryRsp.Name = utils.String(entityToCommit.Name.String)
			}
			if entityToCommit.NonUniqueName.Valid {
				entryRsp.NonUniqueName = utils.String(entityToCommit.NonUniqueName.String)
			}
			entryRsp.Mtime = &entityToCommit.Mtime
		}
	}
	return &errCode, nil
}

func HandleClientToServerMessage(pb *sync_pb.ClientToServerMessage, pbRsp *sync_pb.ClientToServerResponse, pg *datastore.Postgres) error {
	// Create ClientToServerResponse and fill general fields for both GU and
	// Commit.
	pbRsp.StoreBirthday = utils.String(STORE_BIRTHDAY)
	pbRsp.ClientCommand = &sync_pb.ClientCommand{
		SetSyncPollInterval:        utils.Int32(SET_SYNC_POLL_INTERVAL),
		MaxCommitBatchSize:         utils.Int32(MAX_COMMIT_BATCH_SIZE),
		SessionsCommitDelaySeconds: utils.Int32(SESSIONS_COMMIT_DELAY_SECONDS)}

	var err error
	if *pb.MessageContents == sync_pb.ClientToServerMessage_GET_UPDATES {
		guRsp := &sync_pb.GetUpdatesResponse{}
		pbRsp.GetUpdates = guRsp
		pbRsp.ErrorCode, err = HandleGetUpdateRequest(pb.GetUpdates, guRsp, pg)
	} else if *pb.MessageContents == sync_pb.ClientToServerMessage_COMMIT {
		commitRsp := &sync_pb.CommitResponse{}
		pbRsp.Commit = commitRsp
		pbRsp.ErrorCode, err = HandleCommitRequest(pb.Commit, commitRsp, pg)
	} else {
		return fmt.Errorf("Unsupported message type of ClientToServerMessage.")
	}

	return err
}
