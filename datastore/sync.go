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
	ID                     string         `json:"id_string" db:"id"`
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

func (pg *Postgres) InsertSyncEntity(entity *SyncEntity) error {
	stmt := `INSERT INTO sync_entities(id, parent_id, old_parent_id, version, mtime, ctime, name, non_unique_name, server_defined_unique_tag, deleted, originator_cache_guid, originator_client_item_id, specifics, folder, client_defined_unique_tag, unique_position) VALUES(:id, :parent_id, :old_parent_id, :version, :mtime, :ctime, :name, :non_unique_name, :server_defined_unique_tag, :deleted, :originator_cache_guid, :originator_client_item_id, :specifics, :folder, :client_defined_unique_tag, :unique_position)`
	_, err := pg.NamedExec(stmt, *entity)
	if err != nil {
		fmt.Println("Insert error: ", err.Error())
	}
	return err
}

func (pg *Postgres) UpdateSyncEntity(entity *SyncEntity) error {
	if entity.Deleted.Valid && entity.Deleted.Bool {
		return pg.DeleteSyncEntity(entity.ID)
	}

	stmt := `UPDATE sync_entites SET parent_id = :parent_id, old_parent_id = :old_parent_id, version = :version, mtime = :mtime, name = :name, non_unique_name = :non_unique_name, specifics = :specifics, folder = :folder, unique_position = :unique_position WHERE id = :id`
	_, err := pg.NamedExec(stmt, *entity)
	if err != nil {
		fmt.Println("Update error: ", err.Error())
	}
	return err
}

func (pg *Postgres) DeleteSyncEntity(id string) error {
	_, err := pg.Exec(`UPDATE sync_entities SET deleted = true WHERE id = $1`, id)
	if err != nil {
		fmt.Println("Delete error: ", err.Error())
	}
	return err
}

func (pg *Postgres) GetSyncEntity(entity *SyncEntity) error {
	return nil
}

func (pg *Postgres) CheckVersion(id string, clientVersion int64) (bool, error) {
	var serverVersion int64
	err := pg.Get(&serverVersion, "SELECT version FROM sync_entities WHERE id = $1", id)
	if err != nil {
		fmt.Println("Get version error: ", err.Error(), "id: ", id)
		return false, nil
	}

	return clientVersion == serverVersion, nil
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
