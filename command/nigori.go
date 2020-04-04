package command

import (
	"time"

	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/satori/go.uuid"
)

const (
	NAME               string = "Nigori"
	SERVER_DEFINED_TAG string = "google_chrome_nigori"
)

func GetNewClientEntity() *sync_pb.SyncEntity {
	now := time.Now().Unix()
	deleted := false
	folder := true
	name := NAME
	serverDefinedTag := SERVER_DEFINED_TAG
	version := int64(1)
	parentIdString := "0"
	idString := uuid.NewV4().String()

	nigoriSpecific := &sync_pb.NigoriSpecifics{}
	specific := &sync_pb.EntitySpecifics_Nigori{Nigori: nigoriSpecific}
	specifics := &sync_pb.EntitySpecifics{SpecificsVariant: specific}

	syncEntity := &sync_pb.SyncEntity{
		Ctime: &now, Mtime: &now, Deleted: &deleted, Folder: &folder,
		Name: &name, ServerDefinedUniqueTag: &serverDefinedTag,
		Version: &version, ParentIdString: &parentIdString,
		IdString: &idString, Specifics: specifics}

	return syncEntity
}
