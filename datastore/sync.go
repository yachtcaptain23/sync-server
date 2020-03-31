package datastore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
)

type SyncEntity struct {
	ID                     string         `json:"id_string" db: "id"`
	ParentID               sql.NullString `json:"parent_id_string" db:"parent_id"`
	OldParentID            sql.NullString `json:"old_parent_id" db:"old_parent_id"`
	Version                int64          `json:"version" db:"version"`
	Mtime                  int64          `json:"mtime" db:"mtime"`
	Ctime                  int64          `json:"ctime" db:"ctime"`
	Name                   sql.NullString `json:"name" db:"name"`
	NonUniqueName          sql.NullString `json:"non_unique_name" db:"non_unique_name"`
	ServerDefinedUniqueTag sql.NullString `json:"server_defined_unique_tag" db:"server_defined_unique_tag"`
	Deleted                sql.NullBool   `json:"deleted" db:"deleted"`
	OriginatorCacheGUID    string         `json:"originator_cache_guid" db:"originator_cache_guid"`
	OriginatorClientItemID string         `json:"originator_client_item_id" db:"originator_client_item_id"`
	Specifics              []byte         `json:"specifics" db:"specifics"`
	Folder                 sql.NullBool   `json:"folder" db:"folder"`
	ClientDefinedUniqueTag sql.NullString `json:"client_defined_unique_tag" db:"client_defined_unique_tag"`
	UniquePosition         []byte         `json:"unique_position" db:"unique_position"`
}

func InsertSyncEntity(entity *SyncEntity) error {
	return nil
}

func CreateSyncEntity(entity *sync_pb.SyncEntity, cacheGuid string) (*SyncEntity, error) {
	var err error
	parentId := sql.NullString{"", false}
	if entity.ParentIdString != nil {
		parentId = sql.NullString{*entity.ParentIdString, true}
	}
	oldParentId := sql.NullString{"", false}
	if entity.OldParentId != nil {
		oldParentId = sql.NullString{*entity.OldParentId, true}
	}
	name := sql.NullString{"", false}
	if entity.Name != nil {
		name = sql.NullString{*entity.Name, true}
	}
	nonUniqueName := sql.NullString{"", false}
	if entity.NonUniqueName != nil {
		nonUniqueName = sql.NullString{*entity.NonUniqueName, true}
	}
	serverDefinedUniqueTag := sql.NullString{"", false}
	if entity.ServerDefinedUniqueTag != nil {
		serverDefinedUniqueTag = sql.NullString{*entity.ServerDefinedUniqueTag, true}
	}
	clientDefinedUniqueTag := sql.NullString{"", false}
	if entity.ClientDefinedUniqueTag != nil {
		clientDefinedUniqueTag = sql.NullString{*entity.ClientDefinedUniqueTag, true}
	}
	var specifics []byte
	if entity.Specifics != nil {
		specifics, err = proto.Marshal(entity.Specifics)
		if err != nil {
			fmt.Println("Marshal Error", err.Error())
			return nil, err
		}
	}
	var uniquePosition []byte
	if entity.UniquePosition != nil {
		uniquePosition, err = proto.Marshal(entity.UniquePosition)
		if err != nil {
			fmt.Println("Marshal Error", err.Error())
			return nil, err
		}
	}

	deleted := sql.NullBool{false, false}
	if entity.Deleted != nil {
		deleted = sql.NullBool{*entity.Deleted, true}
	}
	folder := sql.NullBool{false, false}
	if entity.Folder != nil {
		folder = sql.NullBool{*entity.Folder, true}
	}

	return &SyncEntity{
		ID:                     uuid.NewV4().String(),
		ParentID:               parentId,
		OldParentID:            oldParentId,
		Version:                *entity.Version,
		Ctime:                  *entity.Mtime,
		Mtime:                  time.Now().Unix(),
		Name:                   name,
		NonUniqueName:          nonUniqueName,
		ServerDefinedUniqueTag: serverDefinedUniqueTag,
		Deleted:                deleted,
		OriginatorCacheGUID:    cacheGuid,
		OriginatorClientItemID: *entity.IdString,
		ClientDefinedUniqueTag: clientDefinedUniqueTag,
		Specifics:              specifics,
		Folder:                 folder,
		UniquePosition:         uniquePosition,
	}, nil
}
