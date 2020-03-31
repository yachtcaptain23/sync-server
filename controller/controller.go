package controller

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/brave-experiments/sync-server/utils"
	"github.com/go-chi/chi"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
)

const (
	// Always return hard-coded value of store birthday for now
	STORE_BIRTHDAY                string = "1"
	MAX_COMMIT_BATCH_SIZE         int32  = 90
	SESSIONS_COMMIT_DELAY_SECONDS int32  = 11
	SET_SYNC_POLL_INTERVAL        int32  = 30
)

func SyncRouter(datastore *datastore.Postgres) chi.Router {
	r := chi.NewRouter()
	r.Post("/command/", Command(datastore))
	return r
}

func Dump(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(requestDump))
}

func Command(datastore *datastore.Postgres) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Dump(w, r)

		// Decompress
		var err error
		var message []byte
		var gr *gzip.Reader
		if r.Header.Get("Content-Encoding") == "gzip" {
			gr, err = gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Decompress error", http.StatusInternalServerError)
				return
			}
			message, err = ioutil.ReadAll(gr)
		} else {
			message, err = ioutil.ReadAll(r.Body)
		}
		if err != nil {
			fmt.Println("Error while reading data from client: ", err.Error())
			http.Error(w, "Decompress error", http.StatusInternalServerError)
			return
		}

		// Unmarshal into ClientToServerMessage
		pb := &sync_pb.ClientToServerMessage{}
		err = proto.Unmarshal(message, pb)
		if err != nil {
			fmt.Println("Error while unmarshalling protocol buf: ", err.Error())
			http.Error(w, "Unmarshal error", http.StatusInternalServerError)
			return
		}
		fmt.Println("Received ClientToServerMessage:", pb)

		// Create ClientToServerResponse and fill general fields for both GU and
		// Commit.
		pbRsp := &sync_pb.ClientToServerResponse{}

		pbRsp.StoreBirthday = utils.String(STORE_BIRTHDAY)
		errCode := sync_pb.SyncEnums_SUCCESS
		pbRsp.ErrorCode = &errCode
		pbRsp.ClientCommand = &sync_pb.ClientCommand{
			SetSyncPollInterval:        utils.Int32(SET_SYNC_POLL_INTERVAL),
			MaxCommitBatchSize:         utils.Int32(MAX_COMMIT_BATCH_SIZE),
			SessionsCommitDelaySeconds: utils.Int32(SESSIONS_COMMIT_DELAY_SECONDS)}

		// Create GU response and fill it into the response
		if *pb.MessageContents == sync_pb.ClientToServerMessage_GET_UPDATES {
			fmt.Println("GET UPDATE RECEIVED")
			guMsg := *pb.GetUpdates
			guRsp := &sync_pb.GetUpdatesResponse{}
			pbRsp.GetUpdates = guRsp

			// TODO: process FetchFolders

			// Process from_progress_marker
			// TODO: query DB to get update entries
			if guMsg.FromProgressMarker != nil {
				fmt.Println("Len of from_progress_marker: ", len(guMsg.FromProgressMarker))
				guRsp.NewProgressMarker = make([]*sync_pb.DataTypeProgressMarker, len(guMsg.FromProgressMarker))

				for i := 0; i < len(guMsg.FromProgressMarker); i++ {
					guRsp.NewProgressMarker[i] = &sync_pb.DataTypeProgressMarker{}
					guRsp.NewProgressMarker[i].DataTypeId = guMsg.FromProgressMarker[i].DataTypeId
					// TODO: latest timestamp of recoords?
					guRsp.NewProgressMarker[i].Token, _ = time.Now().MarshalJSON()
				}
			}

			if *pb.GetUpdates.GetUpdatesOrigin == sync_pb.SyncEnums_NEW_CLIENT {
				fmt.Println("New client")
				guRsp.Entries = make([]*sync_pb.SyncEntity, 1)
				ctime := int64(0)
				mtime := time.Now().Unix()
				deleted := false
				folder := true
				name := "Nigori"
				serverDefinedTag := "google_chrome_nigori"
				version := time.Now().Unix()
				parentIdString := "0"
				idString := uuid.NewV4().String()

				nigoriSpecific := &sync_pb.NigoriSpecifics{}
				specific := &sync_pb.EntitySpecifics_Nigori{Nigori: nigoriSpecific}
				specifics := &sync_pb.EntitySpecifics{SpecificsVariant: specific}
				guRsp.Entries[0] = &sync_pb.SyncEntity{
					Ctime: &ctime, Mtime: &mtime, Deleted: &deleted, Folder: &folder,
					Name: &name, ServerDefinedUniqueTag: &serverDefinedTag,
					Version: &version, ParentIdString: &parentIdString,
					IdString: &idString, Specifics: specifics}
			}

			// TODO: Implement batch reply and update the value accordingly
			changesRemaining := int64(0)
			guRsp.ChangesRemaining = &changesRemaining
		} else if *pb.MessageContents == sync_pb.ClientToServerMessage_COMMIT {

		}

		out, err := proto.Marshal(pbRsp)
		if err != nil {
			fmt.Println("Error while marshalling protocol buf: ", err.Error())
			http.Error(w, "Marshal Error", http.StatusInternalServerError)
			return
		}

		/*
			pbRspOut := &sync_pb.ClientToServerResponse{}
			err = proto.Unmarshal(out, pbRspOut)
			if err != nil {
				fmt.Println("Error while unmarshalling protocal buf: ", err.Error())
				http.Error(w, "Unmarshal Error", http.StatusInternalServerError)
				return
			}
			fmt.Println("pbRspOut: ", pbRspOut)
		*/
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(out)
		if err != nil {
			fmt.Println("Error writing response body: ", err.Error())
			return
		}
	})
}
