package persist

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBInterface struct {
	Database *sql.DB
}

func DbSetupConnection(setUpDB bool, mysqlPort int) (*DBInterface, error) {
	log.Printf("Starting database")
	db, err := sql.Open("mysql", fmt.Sprintf("root:password@tcp(localhost:%v)/irisDB", mysqlPort))
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

	log.Printf("Database successfully started")
	return &DBInter, nil
}
