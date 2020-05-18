package controller_test

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave-experiments/sync-server/controller"
	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/datastore/datastoretest"
	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/brave-experiments/sync-server/utils"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
}

func (suite *ControllerTestSuite) SetupSuite() {
	datastore.Table = "client-entity-token-test-controllor"
	var err error
	suite.dynamo, err = datastore.NewDynamo()
	suite.Require().NoError(err, "Failed to get dynamoDB session")
}

func (suite *ControllerTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetTable(suite.dynamo), "Failed to reset table")
}

func (suite *ControllerTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamo), "Failed to delete table")
}

func (suite *ControllerTestSuite) TestCommand() {
	// Generate request body.
	commitMsg := &sync_pb.CommitMessage{
		Entries: []*sync_pb.SyncEntity{
			{
				IdString: aws.String("id"),
				Version:  aws.Int64(1),
				Deleted:  aws.Bool(false),
				Folder:   aws.Bool(false),
				Specifics: &sync_pb.EntitySpecifics{
					SpecificsVariant: &sync_pb.EntitySpecifics_Nigori{
						Nigori: &sync_pb.NigoriSpecifics{},
					},
				},
			},
		},
		CacheGuid: aws.String("cache_guid"),
	}
	commit := sync_pb.ClientToServerMessage_COMMIT
	msg := &sync_pb.ClientToServerMessage{
		MessageContents: &commit,
		Commit:          commitMsg,
		Share:           aws.String(""),
	}

	body, err := proto.Marshal(msg)
	suite.Require().NoError(err, "proto.Marshal should succeed")

	req, err := http.NewRequest("POST", "v2/command/", bytes.NewBuffer(body))
	suite.Require().NoError(err, "NewRequest should succeed")
	req.Header.Set("Authorization", "Bearer token")

	handler := controller.Command(suite.dynamo)

	// Test unauthorized response.
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusUnauthorized, rr.Code)

	// Add token into DB to simulate authenicate.
	ts := utils.UnixMilli(time.Now().Add(time.Minute * 30))
	suite.Require().NoError(suite.dynamo.InsertClientToken("key", "token", ts))

	// Test message without gzip.
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusOK, rr.Code)

	// Test message with gzip.
	buf := new(bytes.Buffer)
	zw := gzip.NewWriter(buf)
	_, err = zw.Write(body)
	suite.Require().NoError(err, "gzip write should succeed")
	err = zw.Close()
	suite.Require().NoError(err, "gzip close should succeed")

	req, err = http.NewRequest("POST", "v2/command/", buf)
	suite.Require().NoError(err, "NewRequest should succeed")
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Encoding", "gzip")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusOK, rr.Code)
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}
