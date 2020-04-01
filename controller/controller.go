package controller

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/brave-experiments/sync-server/command"
	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/sync_pb"
	//	"github.com/brave-experiments/sync-server/utils"
	"github.com/go-chi/chi"
	"github.com/golang/protobuf/proto"
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

func Command(pg *datastore.Postgres) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dump(w, r)

		// Decompress
		var err error
		var msg []byte
		var gr *gzip.Reader
		if r.Header.Get("Content-Encoding") == "gzip" {
			gr, err = gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Decompress error", http.StatusInternalServerError)
				return
			}
			msg, err = ioutil.ReadAll(gr)
		} else {
			msg, err = ioutil.ReadAll(r.Body)
		}
		if err != nil {
			fmt.Println("Error while reading data from client: ", err.Error())
			http.Error(w, "Decompress error", http.StatusInternalServerError)
			return
		}

		// Unmarshal into ClientToServerMessage
		pb := &sync_pb.ClientToServerMessage{}
		err = proto.Unmarshal(msg, pb)
		if err != nil {
			fmt.Println("Error while unmarshalling protocol buf: ", err.Error())
			http.Error(w, "Unmarshal error", http.StatusInternalServerError)
			return
		}
		// fmt.Println("Received ClientToServerMessage:", pb)

		pbRsp := &sync_pb.ClientToServerResponse{}
		err = command.HandleClientToServerMessage(pb, pbRsp, pg)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		out, err := proto.Marshal(pbRsp)
		if err != nil {
			fmt.Println("Error while marshalling protocol buf: ", err.Error())
			http.Error(w, "Marshal Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(out)
		if err != nil {
			fmt.Println("Error writing response body: ", err.Error())
			return
		}
	})
}
