package auth

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/brave-experiments/sync-server/datastore"
	"github.com/satori/go.uuid"
)

const (
	timestampMaxDuration int64 = 120
	tokenMaxDuration     int64 = 86400
)

// Request is a struct used for authenication requests.
type Request struct {
	PublicKey       string `json:"public_key"`
	Timestamp       string `json:"timestamp"`
	SignedTimestamp string `json:"signed_timestamp"`
}

// Response is a struct used for authenication responses.
type Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Authenticate validates client's auth requests and provides the reply.
func Authenticate(r *http.Request, pg *datastore.Postgres) ([]byte, error) {
	var rsp []byte
	// Unmarshal request
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	// Verify signature
	publicKey, err := base64.StdEncoding.DecodeString(req.PublicKey)
	if err != nil {
		return nil, err
	}

	timestamp, err := base64.StdEncoding.DecodeString(req.Timestamp)
	if err != nil {
		return nil, err
	}

	signedTimestamp, err := base64.StdEncoding.DecodeString(req.SignedTimestamp)
	if err != nil {
		return nil, err
	}
	if !ed25519.Verify(publicKey, timestamp, signedTimestamp) {
		return nil, fmt.Errorf("signature verification failed")
	}

	// TODO: Verify the timestamp is not outdated

	// Create a new token, save it in DB, and return it
	expireAt := time.Now().Add(time.Duration(tokenMaxDuration) * time.Second).Unix()
	token := uuid.NewV4().String()
	err = pg.InsertClient(req.PublicKey, token, expireAt)
	if err != nil {
		return nil, err
	}

	authRsp := Response{AccessToken: token, ExpiresIn: tokenMaxDuration}
	rsp, err = json.Marshal(authRsp)
	return rsp, err
}
