package controller

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/brave-experiments/sync-server/auth"
	"github.com/brave-experiments/sync-server/command"
	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/brave-experiments/sync-server/timestamp"
	"github.com/go-chi/chi"
	"github.com/golang/protobuf/proto"
)

// SyncRouter add routers for command and auth endpoint requests.
func SyncRouter(datastore *datastore.Postgres) chi.Router {
	r := chi.NewRouter()
	r.Post("/command/", Command(datastore))
	r.Post("/auth", Auth(datastore))
	r.Get("/timestamp", Timestamp)
	return r
}

func sendJSONRsp(body []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(body)
	if err != nil {
		fmt.Println("Error writing response body:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Timestamp returns a current timestamp back to sync clients.
func Timestamp(w http.ResponseWriter, r *http.Request) {
	body, err := timestamp.GetTimestamp()
	if err != nil {
		fmt.Println("error when getting timestamp:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONRsp(body, w)
}

// Auth handles authentication requests from sync clients.
func Auth(pg *datastore.Postgres) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AUTH")
		body, err := auth.Authenticate(r, pg)
		if err != nil {
			fmt.Println("authenticate error:", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJSONRsp(body, w)
	})
}

// Command handles GetUpdates and Commit requests from sync clients.
func Command(pg *datastore.Postgres) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorize
		clientID, err := auth.Authorize(pg, r)
		if clientID == "" {
			if err != nil {
				fmt.Println("error while authorizing:", err.Error())
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Decompress
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

		pbRsp := &sync_pb.ClientToServerResponse{}
		err = command.HandleClientToServerMessage(pb, pbRsp, pg, clientID)
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
