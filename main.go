package main

import (
	"github.com/brave-experiments/sync-server/server"
	_ "github.com/lib/pq"
)

func main() {
	server.StartServer()
}
