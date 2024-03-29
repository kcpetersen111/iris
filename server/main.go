package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	httpServer "github.com/kcpetersen111/iris/server/adapter/http"
	"github.com/kcpetersen111/iris/server/persist"
)

func main() {
	var setUpDB bool
	var address string
	var port int
	var sqlPort int
	flag.BoolVar(&setUpDB, "s", false, "Will run script to set up the database, Only needs to be done the first time starting up a server (default False)")
	flag.StringVar(&address, "a", "localhost", "The address of the server")
	flag.IntVar(&port, "p", 4444, "The port of the server")
	flag.IntVar(&sqlPort, "P", 3306, "The port of the database")
	flag.Parse()

	db, err := persist.DbSetupConnection(setUpDB, sqlPort)
	db.Database.SetConnMaxLifetime(time.Second * 10)
	if err != nil {
		log.Fatalf("Error in starting up the database: %v", err)
	}

	server := httpServer.NewIrisServer(fmt.Sprintf("%v:%v", address, port), db)
	server.Serve()
}
