package auth

import (
	"encoding/json"
	"net/http"
)

type AuthRsp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// GetAuthReply processes client's requests and provides the reply in JSON.
func GetAuthRsp(r *http.Request) (rsp []byte, err error) {
	authRsp := AuthRsp{"brave5566", int64(86400)} // a stub reply for now
	rsp, err = json.Marshal(authRsp)
	return
}
