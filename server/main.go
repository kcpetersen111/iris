package main

import (
	"flag"
	"log"

	httpServer "github.com/kcpetersen111/iris/server/adapter/http"
	"github.com/kcpetersen111/iris/server/persist"
)

func main() {
	var setUpDB bool
	flag.BoolVar(&setUpDB, "s", false, "Will run script to set up the database, Only needs to be done the first time starting up a server")
	flag.Parse()

	db, err := persist.DbSetupConnection(setUpDB)
	if err != nil {
		log.Fatalf("Error in starting up the database: %v", err)
	}
	server := httpServer.NewIrisServer("localhost:4444", db)
	server.Serve()
}
