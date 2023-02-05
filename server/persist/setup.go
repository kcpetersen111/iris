package persist

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func DbSetupConnection(setUpDB bool) (*sql.DB, error) {
	log.Printf("Starting database setup")
	db, err := sql.Open("mysql", ":password@/irisDB")
	if err != nil {
		log.Printf("Error starting up the database: %v\n", err)
	}
	if setUpDB {
		DbSetup(db)
	}
	log.Printf("Database successfully setup")
	return db, nil
}
