package main

import (
	"log"

	"github.com/kcpetersen111/iris/server/persist"
)

func main() {
	_, err := persist.DbSetupConnection(true)
	if err != nil {
		log.Fatalf("Error in starting up the database: %v", err)
	}
	http.startServer()
}
