package auth

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	fmt.Println("post form:", r.PostForm)
	req := &Request{
		PublicKey:       r.PostFormValue("client_id"),
		Timestamp:       r.PostFormValue("timestamp"),
		SignedTimestamp: r.PostFormValue("client_secret"),
	}

	fmt.Println("sig")
	// Verify the signature.
	publicKey, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		return nil, err
	}
	timestampBytes, err := hex.DecodeString(req.Timestamp)
	if err != nil {
		return nil, err
	}
	signedTimestamp, err := hex.DecodeString(req.SignedTimestamp)
	if err != nil {
		return nil, err
	}
	if !ed25519.Verify(publicKey, timestampBytes, signedTimestamp) {
		fmt.Println("signature verification failed")
		return nil, fmt.Errorf("signature verification failed")
	}

	var timestamp int64
	timestamp, err = strconv.ParseInt(string(timestampBytes), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse timestamp error")
	}
	fmt.Println("Verify timestamp:", timestamp)

	// Verify the timestamp is not outdated
	if time.Now().Unix()-timestamp > timestampMaxDuration {
		fmt.Println("timestamp is outdated")
		return nil, fmt.Errorf("timestamp is outdated")
	}

	fmt.Println("insert")
	// Create a new token, save it in DB, and return it.
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

// Authorize extracts the authorize token from the HTTP request and query the
// database to return the clientID associated with that token if the token is
// valid, otherwise, an empty string will be returned.
func Authorize(pg *datastore.Postgres, r *http.Request) (clientID string) {
	var token string
	// Extract token from the header.
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if token == "" {
		return
	}

	// Query clients table for the token to return the clientID.
	return pg.GetClient(token)
}
