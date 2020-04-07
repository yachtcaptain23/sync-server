package timestamp

import (
	"encoding/json"
	"time"
)

// Timestamp is a structure used for timestamp responses.
type Timestamp struct {
	Timestamp int64 `json:"timestamp"`
}

// GetTimestamp returns the current timestamp in JSON format.
func GetTimestamp() (rsp []byte, err error) {
	time := Timestamp{Timestamp: time.Now().Unix()}
	rsp, err = json.Marshal(time)
	return
}
