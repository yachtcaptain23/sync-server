package command

import (
	"time"

	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/brave-experiments/sync-server/utils"
	"github.com/satori/go.uuid"
)

const (
	name             string = "Nigori"
	serverDefinedTag string = "google_chrome_nigori"
)

// GetNewClientEntity generates the initial nigori folder entity for new
// clients.
func GetNewClientEntity() *sync_pb.SyncEntity {
	now := time.Now().Unix()
	deleted := false
	folder := true
	version := int64(1)
	parentIDString := "0"
	idString := uuid.NewV4().String()

	nigoriSpecific := &sync_pb.NigoriSpecifics{}
	specific := &sync_pb.EntitySpecifics_Nigori{Nigori: nigoriSpecific}
	specifics := &sync_pb.EntitySpecifics{SpecificsVariant: specific}

	syncEntity := &sync_pb.SyncEntity{
		Ctime: &now, Mtime: &now, Deleted: &deleted, Folder: &folder,
		Name: utils.String(name), ServerDefinedUniqueTag: utils.String(serverDefinedTag),
		Version: &version, ParentIdString: &parentIDString,
		IdString: &idString, Specifics: specifics}

	return syncEntity
}
