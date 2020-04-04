package datastore

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
)

type SyncEntity struct {
	Id                     string         `json:"id_string" db:"id"`
	ParentId               sql.NullString `json:"parent_id_string" db:"parent_id"`
	OldParentId            sql.NullString `json:"old_parent_id" db:"old_parent_id"`
	Version                int64          `json:"version" db:"version"`
	Mtime                  int64          `json:"mtime" db:"mtime"`
	Ctime                  int64          `json:"ctime" db:"ctime"`
	Name                   sql.NullString `json:"name" db:"name"`
	NonUniqueName          sql.NullString `json:"non_unique_name" db:"non_unique_name"`
	ServerDefinedUniqueTag sql.NullString `json:"server_defined_unique_tag" db:"server_defined_unique_tag"`
	Deleted                bool           `json:"deleted" db:"deleted"`
	OriginatorCacheGuid    sql.NullString `json:"originator_cache_guid" db:"originator_cache_guid"`
	OriginatorClientItemId sql.NullString `json:"originator_client_item_id" db:"originator_client_item_id"`
	Specifics              []byte         `json:"specifics" db:"specifics"`
	DataTypeID             int            `json:"data_type_id" db:"data_type_id"`
	Folder                 bool           `json:"folder" db:"folder"`
	ClientDefinedUniqueTag sql.NullString `json:"client_defined_unique_tag" db:"client_defined_unique_tag"`
	UniquePosition         []byte         `json:"unique_position" db:"unique_position"`
}

func (pg *Postgres) InsertSyncEntity(entity *SyncEntity) error {
	stmt := `INSERT INTO sync_entities(id, parent_id, old_parent_id, version, mtime, ctime, name, non_unique_name, server_defined_unique_tag, deleted, originator_cache_guid, originator_client_item_id, specifics, data_type_id, folder, client_defined_unique_tag, unique_position) VALUES(:id, :parent_id, :old_parent_id, :version, :mtime, :ctime, :name, :non_unique_name, :server_defined_unique_tag, :deleted, :originator_cache_guid, :originator_client_item_id, :specifics, :data_type_id, :folder, :client_defined_unique_tag, :unique_position)`
	_, err := pg.NamedExec(stmt, *entity)
	if err != nil {
		fmt.Println("Insert error: ", err.Error())
	}
	return err
}

func (pg *Postgres) UpdateSyncEntity(entity *SyncEntity) error {
	if entity.Deleted {
		return pg.DeleteSyncEntity(entity.Id)
	}

	stmt := `UPDATE sync_entities SET parent_id = :parent_id, old_parent_id = :old_parent_id, version = :version, mtime = :mtime, name = :name, non_unique_name = :non_unique_name, specifics = :specifics, folder = :folder, unique_position = :unique_position WHERE id = :id`
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

func (pg *Postgres) GetUpdatesForType(dataType int32, clientToken int64, fetchFolders bool) (entities []SyncEntity, err error) {
	stmt := "SELECT * FROM sync_entities WHERE data_type_id = $1 AND mtime > $2"
	if !fetchFolders {
		stmt += "AND folder = false"
	}
	stmt += " ORDER BY mtime"
	err = pg.Select(&entities, stmt, dataType, clientToken)
	return
}

func CreateDBSyncEntity(entity *sync_pb.SyncEntity, cacheGuid string) (*SyncEntity, error) {
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
	// if entity.Specifics != nil {  // TODO: make sure this is present in the
	// validator
	specifics, err = proto.Marshal(entity.Specifics)
	if err != nil {
		fmt.Println("Marshal Error", err.Error())
		return nil, err
	}
	// }

	// TODO: wrap getting type ID into an util function
	structField := reflect.ValueOf(entity.Specifics.SpecificsVariant).Elem().Type().Field(0)
	tag := structField.Tag.Get("protobuf")
	s := strings.Split(tag, ",")
	dataTypeId, _ := strconv.Atoi(s[1])

	var uniquePosition []byte
	if entity.UniquePosition != nil {
		uniquePosition, err = proto.Marshal(entity.UniquePosition)
		if err != nil {
			fmt.Println("Marshal Error", err.Error())
			return nil, err
		}
	}

	deleted := false
	if entity.Deleted != nil {
		deleted = *entity.Deleted
	}
	folder := false
	if entity.Folder != nil {
		folder = *entity.Folder
	}

	id := *entity.IdString
	originatorCacheGuid := sql.NullString{"", false}
	originatorClientItemId := sql.NullString{"", false}
	if len(cacheGuid) > 0 {
		if *entity.Version == 0 {
			id = uuid.NewV4().String()
		}
		originatorCacheGuid = sql.NullString{cacheGuid, true}
		originatorClientItemId = sql.NullString{*entity.IdString, true}
	}

	return &SyncEntity{
		Id:                     id,
		ParentId:               parentId,
		OldParentId:            oldParentId,
		Version:                *entity.Version,
		Ctime:                  *entity.Mtime,
		Mtime:                  time.Now().Unix(),
		Name:                   name,
		NonUniqueName:          nonUniqueName,
		ServerDefinedUniqueTag: serverDefinedUniqueTag,
		Deleted:                deleted,
		OriginatorCacheGuid:    originatorCacheGuid,
		OriginatorClientItemId: originatorClientItemId,
		ClientDefinedUniqueTag: clientDefinedUniqueTag,
		Specifics:              specifics,
		Folder:                 folder,
		UniquePosition:         uniquePosition,
		DataTypeID:             dataTypeId,
	}, nil
}

func CreatePBSyncEntity(entity *SyncEntity) (*sync_pb.SyncEntity, error) {
	pbEntity := &sync_pb.SyncEntity{
		IdString: &entity.Id, Version: &entity.Version, Mtime: &entity.Mtime,
		Ctime: &entity.Ctime, Deleted: &entity.Deleted, Folder: &entity.Folder}
	specifics := &sync_pb.EntitySpecifics{}
	err := proto.Unmarshal(entity.Specifics, specifics)
	if err != nil {
		fmt.Println("[CreatePBSyncEntity] Error when unmarshalling specifics:", err.Error())
		return nil, err
	}
	pbEntity.Specifics = specifics

	if entity.UniquePosition != nil {
		uniquePosition := &sync_pb.UniquePosition{}
		err := proto.Unmarshal(entity.UniquePosition, uniquePosition)
		if err != nil {
			fmt.Println("[CreatePBSyncEntity] Error when unmarshalling specifics:", err.Error())
			return nil, err
		}
		pbEntity.UniquePosition = uniquePosition
	}

	if entity.ParentId.Valid {
		pbEntity.ParentIdString = &entity.ParentId.String
	}

	if entity.Name.Valid {
		pbEntity.Name = &entity.Name.String
	}

	if entity.NonUniqueName.Valid {
		pbEntity.NonUniqueName = &entity.NonUniqueName.String
	}

	if entity.ServerDefinedUniqueTag.Valid {
		pbEntity.ServerDefinedUniqueTag = &entity.ServerDefinedUniqueTag.String
	}

	if entity.ClientDefinedUniqueTag.Valid {
		pbEntity.ClientDefinedUniqueTag = &entity.ClientDefinedUniqueTag.String
	}

	if entity.OriginatorCacheGuid.Valid {
		pbEntity.OriginatorCacheGuid = &entity.OriginatorCacheGuid.String
	}

	if entity.OriginatorClientItemId.Valid {
		pbEntity.OriginatorClientItemId = &entity.OriginatorClientItemId.String
	}

	return pbEntity, nil
}
