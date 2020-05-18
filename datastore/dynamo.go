package datastore

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	// Strings for the primary key
	pk     string = "ClientID"
	projPk string = "ClientID, ID"

	// Strings for (ClientID, DataTypeMtime) GSI
	clientIDDataTypeMtimeIdx   string = "ClientIDDataTypeMtimeIndex"
	clientIDDataTypeMtimeIdxPk string = "ClientID"
	clientIDDataTypeMtimeIdxSk string = "DataTypeMtime"

	// Strings for (ID, ExpireAt) GSI
	idExpireAtIdx   string = "IDExpireAtIndex"
	idExpireAtIdxPk string = "ID"
	idExpireAtIdxSk string = "ExpireAt"
)

var (
	// Table is the name of the table in dynamoDB, could be modified in tests.
	Table string = os.Getenv("TABLE_NAME")
)

// PrimaryKey struct is used to represent the primary key of our table.
type PrimaryKey struct {
	ClientID string // Hash key
	ID       string // Range key
}

// Dynamo is a Datastore wrapper around a dynamoDB.
type Dynamo struct {
	*dynamodb.DynamoDB
}

// NewDynamo returns a dynamoDB client to be used.
func NewDynamo() (*Dynamo, error) {
	config := &aws.Config{
		Region:   aws.String(os.Getenv("AWS_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_ENDPOINT")),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	db := dynamodb.New(sess)
	return &Dynamo{db}, nil
}
