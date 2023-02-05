package persist

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func DbSetupConnection(setUpDB bool) (*sql.DB, error) {
	log.Printf("Starting database setup")
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/irisDB")
	if err != nil {
		// log.Printf("Error starting up the database: %v\n", err)
		return nil, err
	}
	log.Printf("Connected to database successfully")
	if setUpDB {
		err = DbSetup(db)
	}
	if err != nil {
		return nil, err
	}

	log.Printf("Database successfully setup")
	return db, nil
}
