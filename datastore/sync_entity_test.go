package datastore

import (
	"database/sql"
	"testing"
	"time"

	"github.com/brave-experiments/sync-server/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

var (
	nullStr = utils.StringOrNull(nil)
)

type SyncEntityTestSuite struct {
	suite.Suite
}

func (suite *SyncEntityTestSuite) SetupSuite() {
	pg, err := NewPostgres(false)
	suite.Require().NoError(err, "Failed to get postgres conn")

	m, err := pg.NewMigrate()
	suite.Require().NoError(err, "Failed to create migrate instance")

	ver, dirty, _ := m.Version()
	if dirty {
		suite.Require().NoError(m.Force(int(ver)))
	}
	if ver > 0 {
		suite.Require().NoError(m.Down(), "Failed to migrate down cleanly")
	}

	suite.Require().NoError(pg.Migrate(), "Failed to fully migrate")
}

func (suite *SyncEntityTestSuite) SetupTest() {
	tables := []string{"sync_entities", "clients"}

	pg, err := NewPostgres(false)
	suite.Require().NoError(err, "Failed to get postgres conn")

	for _, table := range tables {
		_, err = pg.DB.Exec("delete from " + table)
		suite.Require().NoError(err, "Failed to get clean table")
	}

	// Insert a dummy client
	client := &Client{ID: "id", Token: uuid.NewV4().String(),
		ExpireAt: time.Now().Add(time.Duration(5 * time.Minute)).Unix()}
	_, err = pg.InsertClient(client.ID, client.Token, client.ExpireAt)
	suite.Require().NoError(err, "Insert dummy client should succeed")
}

func createSyncEntity() *SyncEntity {
	t := time.Now().Unix()
	bytes := []byte{}
	return &SyncEntity{
		ID:                     uuid.NewV4().String(),
		ParentID:               nullStr,
		OldParentID:            nullStr,
		Version:                1,
		Mtime:                  t,
		Ctime:                  t,
		Name:                   utils.StringOrNull(utils.String("test")),
		NonUniqueName:          nullStr,
		ServerDefinedUniqueTag: nullStr,
		DeletedAt:              utils.Int64OrNull(nil),
		OriginatorCacheGUID:    nullStr,
		OriginatorClientItemID: nullStr,
		Specifics:              bytes,
		DataTypeID:             123,
		Folder:                 false,
		ClientDefinedUniqueTag: nullStr,
		UniquePosition:         bytes,
		ClientID:               "id",
	}
}

func (suite *SyncEntityTestSuite) TestInsertSyncEntity() {
	pg, err := NewPostgres(false)
	suite.Require().NoError(err)

	entity := createSyncEntity()
	err = pg.InsertSyncEntity(entity)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	entity.ID = uuid.NewV4().String()
	entity.ServerDefinedUniqueTag = utils.StringOrNull(utils.String("server_tag"))
	err = pg.InsertSyncEntity(entity)
	suite.Require().NoError(err, "Insert sync entity should succeed")
	insertedID := entity.ID
	entity.ID = uuid.NewV4().String()
	err = pg.InsertSyncEntity(entity)
	suite.Require().Error(err, "Insert duplicate server_defined_unique_tag should fail")

	entity.ID = insertedID
	entity.DeletedAt = utils.Int64OrNull(utils.Int64(time.Now().Unix()))
	err = pg.UpdateSyncEntity(entity)
	suite.Require().NoError(err, "Update sync entity should succeed")
	entity.ID = uuid.NewV4().String()
	entity.DeletedAt = utils.Int64OrNull(nil)
	err = pg.InsertSyncEntity(entity)
	suite.Require().NoError(err,
		"Insert a sync entity where the server_defined_unique_tag is the same as a deleted sync entity should succeed")

	entity.ID = uuid.NewV4().String()
	entity.ServerDefinedUniqueTag = nullStr
	entity.ClientDefinedUniqueTag = utils.StringOrNull(utils.String("client_tag"))
	err = pg.InsertSyncEntity(entity)
	suite.Require().NoError(err, "Insert sync entity should succeed")
	insertedID = entity.ID
	entity.ID = uuid.NewV4().String()
	err = pg.InsertSyncEntity(entity)
	suite.Require().Error(err, "Insert duplicate client_defined_unique_tag should fail")

	entity.ID = insertedID
	entity.DeletedAt = utils.Int64OrNull(utils.Int64(time.Now().Unix()))
	err = pg.UpdateSyncEntity(entity)
	suite.Require().NoError(err, "Update sync entity should succeed")
	entity.ID = uuid.NewV4().String()
	entity.DeletedAt = utils.Int64OrNull(nil)
	err = pg.InsertSyncEntity(entity)
	suite.Require().NoError(err,
		"Insert a sync entity where the client_defined_unique_tag is the same as a deleted sync entity should succeed")
}

func (suite *SyncEntityTestSuite) TestUpdateSyncEntity() {
	pg, err := NewPostgres(false)
	suite.Require().NoError(err)

	entity := createSyncEntity()
	err = pg.InsertSyncEntity(entity)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	// Test updating name
	var name sql.NullString
	entity.Name = utils.StringOrNull(utils.String("test2"))
	err = pg.UpdateSyncEntity(entity)
	suite.Require().NoError(err, "Update sync entity should succeed")
	err = pg.Get(&name, "SELECT name FROM sync_entities WHERE id = $1", entity.ID)
	suite.Require().NoError(err, "Get sync entity should succeed")
	suite.Assert().Equal(sql.NullString{String: "test2", Valid: true}, name)

	// Test soft delete
	t := time.Now().Unix()
	entity.DeletedAt = utils.Int64OrNull(utils.Int64(t))
	err = pg.UpdateSyncEntity(entity)
	suite.Require().NoError(err, "Update sync entity should succeed")
	var deletedAt sql.NullInt64
	err = pg.Get(&deletedAt, "SELECT deleted_at FROM sync_entities WHERE id = $1", entity.ID)
	suite.Require().NoError(err, "Get sync entity should succeed")
	suite.Assert().Equal(sql.NullInt64{Int64: t, Valid: true}, deletedAt)

	// Delete a deleted entry should error out
	t = time.Now().Unix()
	entity.DeletedAt = utils.Int64OrNull(utils.Int64(t))
	err = pg.UpdateSyncEntity(entity)
	suite.Require().EqualError(err, "No rows updated")
}

func (suite *SyncEntityTestSuite) TestCheckVersion() {
	// Insert a record with version 1
	pg, err := NewPostgres(false)
	suite.Require().NoError(err)

	entity := createSyncEntity()
	err = pg.InsertSyncEntity(entity)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	match, err := pg.CheckVersion(entity.ID, 1)
	suite.Require().NoError(err, "Check version should succeed")
	suite.Require().True(match, "Should return true when server and client versions are matched")

	match, err = pg.CheckVersion(entity.ID, 2)
	suite.Require().NoError(err, "Check version should succeed")
	suite.Require().False(match, "Should return false when server and client versions are not matched")
}

func (suite *SyncEntityTestSuite) TestGetUpdatesForType() {
	pg, err := NewPostgres(false)
	suite.Require().NoError(err)

	entity := createSyncEntity()
	entity.Mtime -= 2
	err = pg.InsertSyncEntity(entity)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	clientToken := entity.Mtime + 1
	entity2 := createSyncEntity()
	// entity2.Mtime = entity.Mtime + 1
	err = pg.InsertSyncEntity(entity2)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	entity3 := createSyncEntity()
	// entity3.Mtime = entity2.Mtime + 1
	entity3.Folder = true
	err = pg.InsertSyncEntity(entity3)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	entity4 := createSyncEntity()
	entity4.Mtime = entity3.Mtime + 1
	entity4.DataTypeID = 456
	err = pg.InsertSyncEntity(entity4)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	// Get updates for type 123 including folders should return entity2 and 3.
	dataType := int32(123)
	expectedEntities := []SyncEntity{*entity2, *entity3}
	savedEntities, err := pg.GetUpdatesForType(dataType, clientToken, true, entity.ClientID)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(expectedEntities, savedEntities)

	// Get updates for type 123 without folders should return entity2.
	expectedEntities = []SyncEntity{*entity2}
	savedEntities, err = pg.GetUpdatesForType(dataType, clientToken, false, entity.ClientID)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(expectedEntities, savedEntities)

	// Get updates for type 456 should return entity4.
	expectedEntities = []SyncEntity{*entity4}
	dataType = 456
	savedEntities, err = pg.GetUpdatesForType(dataType, clientToken, false, entity.ClientID)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(expectedEntities, savedEntities)
}

func (suite *SyncEntityTestSuite) TestGetServerDefinedUniqueEntity() {
	pg, err := NewPostgres(false)
	suite.Require().NoError(err)

	// Insert a entry without server tag.
	entity := createSyncEntity()
	err = pg.InsertSyncEntity(entity)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	// Insert a deleted entry with server tag.
	entity2 := createSyncEntity()
	entity2.DeletedAt = utils.Int64OrNull(utils.Int64(time.Now().Unix()))
	entity2.ServerDefinedUniqueTag = utils.StringOrNull(utils.String("server"))
	err = pg.InsertSyncEntity(entity2)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	// Insert a entry with server tag.
	entity3 := createSyncEntity()
	entity3.ServerDefinedUniqueTag = utils.StringOrNull(utils.String("server"))
	err = pg.InsertSyncEntity(entity3)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	// Insert another dummy client
	client := &Client{ID: "id2", Token: uuid.NewV4().String(),
		ExpireAt: time.Now().Add(time.Duration(5 * time.Minute)).Unix()}
	_, err = pg.InsertClient(client.ID, client.Token, client.ExpireAt)
	suite.Require().NoError(err, "Insert dummy client should succeed")

	// Insert a entry with server tag and another client ID.
	entity4 := createSyncEntity()
	entity4.ServerDefinedUniqueTag = utils.StringOrNull(utils.String("server"))
	entity4.ClientID = "id2"
	err = pg.InsertSyncEntity(entity4)
	suite.Require().NoError(err, "Insert sync entity should succeed")

	// Should return just one entry with the target server tag and client ID.
	savedEntity, err := pg.GetServerDefinedUniqueEntity("server", entity.ClientID)
	suite.Require().NoError(err, "GetServerDefinedUniqueEntity should succeed")
	suite.Assert().Equal(entity3, savedEntity)
}

func (suite *SyncEntityTestSuite) TestCreateDBSyncEntity() {
	// TODO:
	//	- Test case with cacheGUID passed
	//	- Test case without cacheGUID passed
	//	- Test case with marshalling specifics
	//	- Test case with marshalling unique postitions
	//	- Test deleted and folder values
}

func (suite *SyncEntityTestSuite) TestCreatePBSyncEntity() {
	// TODO:
	//	- Test specifics and unique position could be unmarshalled correctly
}

func TestSyncEntityTestSuite(t *testing.T) {
	suite.Run(t, new(SyncEntityTestSuite))
}
