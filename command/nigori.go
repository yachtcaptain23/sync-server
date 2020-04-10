package command

import (
	"database/sql"
	"time"

	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/brave-experiments/sync-server/utils"
	"github.com/satori/go.uuid"
)

const (
	name             string = "Nigori"
	serverDefinedTag string = "google_chrome_nigori"
)

func createNewClientEntity(clientID string) *sync_pb.SyncEntity {
	now := time.Now().Unix()
	deleted := false
	folder := true
	version := int64(1)
	parentIDString := "0"
	idString := uuid.NewV4().String()

	nigoriSpecific := &sync_pb.NigoriSpecifics{}
	specific := &sync_pb.EntitySpecifics_Nigori{Nigori: nigoriSpecific}
	specifics := &sync_pb.EntitySpecifics{SpecificsVariant: specific}

	pbEntity := &sync_pb.SyncEntity{
		Ctime: &now, Mtime: &now, Deleted: &deleted, Folder: &folder,
		Name: utils.String(name), ServerDefinedUniqueTag: utils.String(serverDefinedTag),
		Version: &version, ParentIdString: &parentIDString,
		IdString: &idString, Specifics: specifics}

	return pbEntity
}

// GetNewClientEntity gets the nigori top level folder entity for new clients.
// If it is not existed in the DB yet, a new one will be created.
func GetNewClientEntity(pg *datastore.Postgres, clientID string) (pbEntity *sync_pb.SyncEntity, dbEntity *datastore.SyncEntity, err error) {
	dbEntity, err = pg.GetServerDefinedUniqueEntity(serverDefinedTag, clientID)
	if err == sql.ErrNoRows { // create a new entity
		pbEntity = createNewClientEntity(clientID)
		dbEntity, err = datastore.CreateDBSyncEntity(pbEntity, "", clientID)
		if err != nil {
			return nil, nil, err
		}
		err = pg.InsertSyncEntity(dbEntity)
	} else if err != nil {
		return nil, nil, err
	} else {
		pbEntity, err = datastore.CreatePBSyncEntity(dbEntity)
	}

	return
}
