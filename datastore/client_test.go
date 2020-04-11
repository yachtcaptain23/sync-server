package datastore

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
}

func (suite *ClientTestSuite) SetupSuite() {
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

func (suite *ClientTestSuite) SetupTest() {
	tables := []string{"clients", "sync_entities"}

	pg, err := NewPostgres(false)
	suite.Require().NoError(err, "Failed to get postgres conn")

	for _, table := range tables {
		_, err = pg.DB.Exec("delete from " + table)
		suite.Require().NoError(err, "Failed to get clean table")
	}
}

func (suite *ClientTestSuite) TestInsertClient() {
	pg, err := NewPostgres(false)
	suite.Require().NoError(err)

	// Insert a non-exist client
	client := &Client{ID: "id", Token: uuid.NewV4().String(), ExpireAt: time.Now().Unix()}
	var savedClient *Client
	savedClient, err = pg.InsertClient(client.ID, client.Token, client.ExpireAt)
	suite.Require().NoError(err, "Insert client should succeed")
	suite.Assert().Equal(client, savedClient)

	// Insert a client with the same ID should update the entry
	client2 := &Client{ID: "id", Token: uuid.NewV4().String(), ExpireAt: time.Now().Unix()}
	savedClient, err = pg.InsertClient(client2.ID, client2.Token, client2.ExpireAt)
	suite.Require().NoError(err, "Insert client should succeed")
	suite.Assert().Equal(client2, savedClient)
}

func (suite *ClientTestSuite) TestGetClient() {
	pg, err := NewPostgres(false)
	suite.Require().NoError(err)

	// Empty string and nil error should be returned when it's not a valid token.
	var id string
	client := &Client{ID: "id", Token: uuid.NewV4().String(), ExpireAt: time.Now().Unix()}
	id, err = pg.GetClient(client.Token)
	suite.Require().NoError(err, "Get client should succeed")
	suite.Assert().Equal("", id)

	// Outdated token should return empty string and nil error.
	_, err = pg.InsertClient(client.ID, client.Token, client.ExpireAt)
	suite.Require().NoError(err, "Insert client should succeed")
	id, err = pg.GetClient(client.Token)
	suite.Require().NoError(err, "Get client should succeed")
	suite.Assert().Equal("", id)

	// Target client ID and nil error should be returned for a valid token.
	client2 := &Client{ID: "id2", Token: uuid.NewV4().String(),
		ExpireAt: time.Now().Add(time.Duration(5 * time.Minute)).Unix()}
	_, err = pg.InsertClient(client2.ID, client2.Token, client2.ExpireAt)
	suite.Require().NoError(err, "Insert client should succeed")
	id, err = pg.GetClient(client2.Token)
	suite.Require().NoError(err, "Get client should succeed")
	suite.Assert().Equal(client2.ID, id)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
