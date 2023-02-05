package persist

import (
	"database/sql"
)

//This will do all of the first time set up for the database
//doing things like creating tables
func DbSetup(db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		`CREATE TABLE test(
			testfield TEXT,
			numberfield INT
			);`,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()

	return err
}
