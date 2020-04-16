package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/brave-experiments/sync-server/utils"
	"github.com/golang/protobuf/proto"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

// SyncEntity represents the underline DB schema of sync entities.
type SyncEntity struct {
	ID                     string         `db:"id"`
	ParentID               sql.NullString `db:"parent_id"`
	OldParentID            sql.NullString `db:"old_parent_id"`
	Version                int64          `db:"version"`
	Mtime                  int64          `db:"mtime"`
	Ctime                  int64          `db:"ctime"`
	Name                   sql.NullString `db:"name"`
	NonUniqueName          sql.NullString `db:"non_unique_name"`
	ServerDefinedUniqueTag sql.NullString `db:"server_defined_unique_tag"`
	DeletedAt              sql.NullInt64  `db:"deleted_at"`
	OriginatorCacheGUID    sql.NullString `db:"originator_cache_guid"`
	OriginatorClientItemID sql.NullString `db:"originator_client_item_id"`
	Specifics              []byte         `db:"specifics"`
	DataTypeID             int            `db:"data_type_id"`
	Folder                 bool           `db:"folder"`
	ClientDefinedUniqueTag sql.NullString `db:"client_defined_unique_tag"`
	UniquePosition         []byte         `db:"unique_position"`
	ClientID               string         `db:"client_id"`
}

// InsertSyncEntity inserts a new sync entity into postgres database.
func (pg *Postgres) InsertSyncEntity(entity *SyncEntity) error {
	stmt := `INSERT INTO sync_entities(id, parent_id, old_parent_id, version, mtime, ctime, name, non_unique_name, server_defined_unique_tag, deleted_at, originator_cache_guid, originator_client_item_id, specifics, data_type_id, folder, client_defined_unique_tag, unique_position, client_id) VALUES(:id, :parent_id, :old_parent_id, :version, :mtime, :ctime, :name, :non_unique_name, :server_defined_unique_tag, :deleted_at, :originator_cache_guid, :originator_client_item_id, :specifics, :data_type_id, :folder, :client_defined_unique_tag, :unique_position, :client_id)`
	_, err := pg.NamedExec(stmt, *entity)
	if err != nil {
		fmt.Println("Insert error: ", err.Error())
	}
	return err
}

// InsertSyncEntities inserts an array of entities in a single transaction.
func (pg *Postgres) InsertSyncEntities(entities []*SyncEntity) error {
	stmt := `INSERT INTO sync_entities(id, parent_id, old_parent_id, version, mtime, ctime, name, non_unique_name, server_defined_unique_tag, deleted_at, originator_cache_guid, originator_client_item_id, specifics, data_type_id, folder, client_defined_unique_tag, unique_position, client_id) VALUES(:id, :parent_id, :old_parent_id, :version, :mtime, :ctime, :name, :non_unique_name, :server_defined_unique_tag, :deleted_at, :originator_cache_guid, :originator_client_item_id, :specifics, :data_type_id, :folder, :client_defined_unique_tag, :unique_position, :client_id)`
	tx, err := pg.DB.Beginx()
	if err != nil {
		return err
	}
	defer pg.RollbackTx(tx)

	for _, entity := range entities {
		_, err := tx.NamedExec(stmt, *entity)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

// UpdateSyncEntity updates a sync entity in postgres database.
func (pg *Postgres) UpdateSyncEntity(entity *SyncEntity) error {
	stmt := `UPDATE sync_entities SET deleted_at = :deleted_at, parent_id = :parent_id, old_parent_id = :old_parent_id, version = :version, mtime = :mtime, name = :name, non_unique_name = :non_unique_name, specifics = :specifics, folder = :folder, unique_position = :unique_position WHERE id = :id`

	result, err := pg.NamedExec(stmt, *entity)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 || err != nil {
		return errors.New("No rows updated")
	}

	return nil
}

// CheckVersion get the sync entry with id saved in the database and checks
// the saved server version against the client version.
func (pg *Postgres) CheckVersion(id string, clientVersion int64) (bool, error) {
	var serverVersion int64
	err := pg.Get(&serverVersion, "SELECT version FROM sync_entities WHERE id = $1", id)
	if err != nil {
		fmt.Println("Get version error: ", err.Error(), "id: ", id)
		return false, err
	}

	return clientVersion == serverVersion, nil
}

// GetUpdatesForType retrieves sync entities of a data type where it's mtime
// is later than the client token.
func (pg *Postgres) GetUpdatesForType(dataType int32, clientToken int64, fetchFolders bool, clientID string) (entities []SyncEntity, err error) {
	stmt := "SELECT * FROM sync_entities WHERE data_type_id = $1 AND mtime > $2 AND client_id = $3"
	if !fetchFolders {
		stmt += " AND folder = false"
	}
	stmt += " ORDER BY mtime"
	err = pg.Select(&entities, stmt, dataType, clientToken, clientID)
	return
}

// GetServerDefinedUniqueEntity returns the entity where client_id and
// server_defined_unique_tag is equal to the parameters.
func (pg *Postgres) GetServerDefinedUniqueEntity(tag string, clientID string) (*SyncEntity, error) {
	var entity SyncEntity
	err := pg.Get(&entity,
		"SELECT * FROM sync_entities WHERE server_defined_unique_tag = $1 AND client_id = $2 AND deleted_at IS NULL",
		tag, clientID)
	return &entity, err
}

// IsServerDefinedUniqueEntitiesReady returns wehter sync entities with server
// defined unique tags for a specific client are ready in DB.
func (pg *Postgres) IsServerDefinedUniqueEntitiesReady(tags []string, clientID string) (bool, error) {
	var count int
	query, args, err := sqlx.In(
		`SELECT COUNT(*) FROM sync_entities WHERE client_id = ? AND server_defined_unique_tag IN (?);`,
		clientID, tags)
	if err != nil {
		return false, err
	}
	query = pg.DB.Rebind(query)
	err = pg.Get(&count, query, args...)
	return count == len(tags), err
}

// CreateDBSyncEntity converts a protobuf sync entity into a DB sync entity.
func CreateDBSyncEntity(entity *sync_pb.SyncEntity, cacheGUID string, clientID string) (*SyncEntity, error) {
	var err error
	specifics := []byte{}
	if entity.Specifics != nil { // TODO: make sure this is present in the validator
		specifics, err = proto.Marshal(entity.Specifics)
		if err != nil {
			fmt.Println("Marshal Error", err.Error())
			return nil, err
		}
	}

	// TODO: wrap getting type ID into an util function
	structField := reflect.ValueOf(entity.Specifics.SpecificsVariant).Elem().Type().Field(0)
	tag := structField.Tag.Get("protobuf")
	s := strings.Split(tag, ",")
	dataTypeID, _ := strconv.Atoi(s[1])

	uniquePosition := []byte{}
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
	var deletedAt sql.NullInt64
	if deleted {
		deletedAt = sql.NullInt64{Int64: time.Now().Unix(), Valid: true}
	} else {
		deletedAt = sql.NullInt64{Int64: 0, Valid: false}
	}

	folder := false
	if entity.Folder != nil {
		folder = *entity.Folder
	}

	id := *entity.IdString
	originatorCacheGUID := utils.StringOrNull(nil)
	originatorClientItemID := utils.StringOrNull(nil)
	if len(cacheGUID) > 0 {
		if *entity.Version == 0 {
			id = uuid.NewV4().String()
		}
		originatorCacheGUID = utils.StringOrNull(&cacheGUID)
		originatorClientItemID = utils.StringOrNull(entity.IdString)
	}

	now := time.Now().Unix()
	// ctime is only used when inserting a new entity, here we use client passed
	// ctime if it is passed, otherwise, use current server time as the creation
	// time. When updating, ctime will be ignored later in the query statement.
	cTime := now
	if entity.Ctime != nil {
		cTime = *entity.Ctime
	}

	return &SyncEntity{
		ID:                     id,
		ParentID:               utils.StringOrNull(entity.ParentIdString),
		OldParentID:            utils.StringOrNull(entity.OldParentId),
		Version:                *entity.Version,
		Ctime:                  cTime,
		Mtime:                  now,
		Name:                   utils.StringOrNull(entity.Name),
		NonUniqueName:          utils.StringOrNull(entity.NonUniqueName),
		ServerDefinedUniqueTag: utils.StringOrNull(entity.ServerDefinedUniqueTag),
		DeletedAt:              deletedAt,
		OriginatorCacheGUID:    originatorCacheGUID,
		OriginatorClientItemID: originatorClientItemID,
		ClientDefinedUniqueTag: utils.StringOrNull(entity.ClientDefinedUniqueTag),
		Specifics:              specifics,
		Folder:                 folder,
		UniquePosition:         uniquePosition,
		DataTypeID:             dataTypeID,
		ClientID:               clientID,
	}, nil
}

// CreatePBSyncEntity converts a DB sync entity to a protobuf sync entity.
func CreatePBSyncEntity(entity *SyncEntity) (*sync_pb.SyncEntity, error) {
	pbEntity := &sync_pb.SyncEntity{
		IdString:               &entity.ID,
		ParentIdString:         utils.StringPtr(&entity.ParentID),
		Version:                &entity.Version,
		Mtime:                  &entity.Mtime,
		Ctime:                  &entity.Ctime,
		Name:                   utils.StringPtr(&entity.Name),
		NonUniqueName:          utils.StringPtr(&entity.NonUniqueName),
		ServerDefinedUniqueTag: utils.StringPtr(&entity.ServerDefinedUniqueTag),
		ClientDefinedUniqueTag: utils.StringPtr(&entity.ClientDefinedUniqueTag),
		OriginatorCacheGuid:    utils.StringPtr(&entity.OriginatorCacheGUID),
		OriginatorClientItemId: utils.StringPtr(&entity.OriginatorClientItemID),
		Deleted:                &entity.DeletedAt.Valid,
		Folder:                 &entity.Folder}

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

	return pbEntity, nil
}
