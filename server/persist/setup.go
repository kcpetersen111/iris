package persist

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBInterface struct {
	Database *sql.DB
}

func DbSetupConnection(setUpDB bool) (*DBInterface, error) {
	log.Printf("Starting database setup")
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/irisDB")
	if err != nil {
		// log.Printf("Error starting up the database: %v\n", err)
		return &DBInterface{}, err
	}
	log.Printf("Connected to database successfully")

	DBInter := DBInterface{
		Database: db,
	}

	if setUpDB {
		err = DBInter.DbSetup()
	}
	if err != nil {
		return &DBInterface{}, err
	}

	log.Printf("Database successfully setup")
	return &DBInter, nil
}
