package timestamp

import (
	"encoding/json"
	"strconv"
	"time"
)

// Timestamp is a structure used for timestamp responses.
type Timestamp struct {
	Timestamp string `json:"timestamp"`
}

// GetTimestamp returns the current timestamp in JSON format.
func GetTimestamp() (rsp []byte, err error) {
	time := Timestamp{Timestamp: strconv.FormatInt(time.Now().Unix(), 10)}
	rsp, err = json.Marshal(time)
	return
}
