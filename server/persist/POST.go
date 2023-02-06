package persist

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *DBInterface) InsertTest(message string) error {

	tx, err := db.Database.Begin()
	if err != nil {
		return err
	}
	tx.Exec(
		fmt.Sprintf(`INSERT INTO test
		Values ("%v", 1);`, message))
	return tx.Commit()

}
