package command

import (
	"encoding/binary"
	"fmt"

	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/brave-experiments/sync-server/utils"
)

const (
	// Always return hard-coded value of store birthday for now
	STORE_BIRTHDAY                string = "1"
	MAX_COMMIT_BATCH_SIZE         int32  = 90
	MAX_GU_BATCH_SIZE             int32  = 500
	SESSIONS_COMMIT_DELAY_SECONDS int32  = 11
	SET_SYNC_POLL_INTERVAL        int32  = 30
	NIGORI_TYPE_ID                int32  = 47745
)

func HandleGetUpdatesRequest(guMsg *sync_pb.GetUpdatesMessage, guRsp *sync_pb.GetUpdatesResponse, pg *datastore.Postgres) (*sync_pb.SyncEnums_ErrorType, error) {
	fmt.Println("GET UPDATE RECEIVED")
	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later

	changesRemaining := int64(0)
	guRsp.ChangesRemaining = &changesRemaining

	if guMsg.FromProgressMarker == nil { // nothing to process
		return &errCode, nil
	}

	fetchFolders := true
	if guMsg.FetchFolders != nil {
		fetchFolders = *guMsg.FetchFolders
	}

	maxSize := MAX_GU_BATCH_SIZE
	if guMsg.BatchSize != nil && *guMsg.BatchSize < MAX_GU_BATCH_SIZE {
		maxSize = *guMsg.BatchSize
	}

	// Process from_progress_marker
	guRsp.NewProgressMarker = make([]*sync_pb.DataTypeProgressMarker, len(guMsg.FromProgressMarker))
	guRsp.Entries = make([]*sync_pb.SyncEntity, 0, maxSize)
	for i, fromProgressMarker := range guMsg.FromProgressMarker {
		guRsp.NewProgressMarker[i] = &sync_pb.DataTypeProgressMarker{}
		guRsp.NewProgressMarker[i].DataTypeId = fromProgressMarker.DataTypeId

		// Default token value is client's token, otherwise 0.
		// This token will be updated when we return the updated entities.
		if len(fromProgressMarker.Token) > 0 {
			guRsp.NewProgressMarker[i].Token = fromProgressMarker.Token
		} else {
			guRsp.NewProgressMarker[i].Token = make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(guRsp.NewProgressMarker[i].Token, int64(0))
		}

		if *fromProgressMarker.DataTypeId == NIGORI_TYPE_ID &&
			*guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_NEW_CLIENT {
			fmt.Println("NEW CLIENT")

			// Bypassing chromium's restriction here, our server won't provide the
			// initial encryption keys like chromium does, this will be overwritten
			// by our client.
			guRsp.EncryptionKeys = make([][]byte, 1)
			guRsp.EncryptionKeys[0] = []byte("1234")

			// Create a nigori top folder entry and save into DB
			syncEntity := GetNewClientEntity()
			entityToCommit, err := datastore.CreateDBSyncEntity(syncEntity, "")
			if err != nil {
				errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
				return &errCode, nil
			}
			err = pg.InsertSyncEntity(entityToCommit)
			if err != nil {
				errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
				return &errCode, nil
			}

			guRsp.Entries = append(guRsp.Entries, syncEntity)
			guRsp.NewProgressMarker[i].Token = make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(guRsp.NewProgressMarker[i].Token, entityToCommit.Mtime)
		} else { // Query DB to get update entries for a given data type
			token, n := binary.Varint(guRsp.NewProgressMarker[i].Token)
			if n <= 0 {
				// TODO: return bad message instead
				return &errCode, fmt.Errorf("Failed at decoding token value %v", token)
			}

			entities, err := pg.GetUpdatesForType(*fromProgressMarker.DataTypeId, token, fetchFolders)
			if err != nil {
				fmt.Println("pg.GetUpdatesForType error:", err.Error())
				errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
				return &errCode, nil
			}

			// Fill the PB entry from above DB entries until maxSize is reached.
			j := 0
			for ; j < len(entities) && len(guRsp.Entries) < cap(guRsp.Entries); j++ {
				entity, err := datastore.CreatePBSyncEntity(&entities[j])
				if err != nil {
					errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
					return &errCode, nil
				}
				guRsp.Entries = append(guRsp.Entries, entity)
			}
			changesRemaining += int64(len(entities) - j)

			// If entities are appended, use the lastest mtime as returned token.
			if j != 0 {
				guRsp.NewProgressMarker[i].Token = make([]byte, binary.MaxVarintLen64)
				binary.PutVarint(guRsp.NewProgressMarker[i].Token, *guRsp.Entries[j-1].Mtime)
			}
		}
	}

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
			entityToCommit, err := datastore.CreateDBSyncEntity(v, *commitMsg.CacheGuid)
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
				match, err := pg.CheckVersion(entityToCommit.Id, entityToCommit.Version)
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
			entryRsp.IdString = utils.String(entityToCommit.Id)
			if entityToCommit.ParentId.Valid {
				entryRsp.ParentIdString = utils.String(entityToCommit.ParentId.String)
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
		pbRsp.ErrorCode, err = HandleGetUpdatesRequest(pb.GetUpdates, guRsp, pg)
	} else if *pb.MessageContents == sync_pb.ClientToServerMessage_COMMIT {
		commitRsp := &sync_pb.CommitResponse{}
		pbRsp.Commit = commitRsp
		pbRsp.ErrorCode, err = HandleCommitRequest(pb.Commit, commitRsp, pg)
	} else {
		return fmt.Errorf("Unsupported message type of ClientToServerMessage.")
	}

	return err
}
