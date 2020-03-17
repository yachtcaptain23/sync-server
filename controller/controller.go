package controller

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/go-chi/chi"
)

func SyncRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/sync", Sync)
	return r
}

func Sync(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
}
