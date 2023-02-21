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

func (db *DBInterface) InsertMessage(message, fromUser, platform string, toUser ...string) error {
	if toUser == nil {
		return fmt.Errorf("No users selected to send to")
	}

	tx, err := db.Database.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	for toU := range toUser {

		_, err = tx.Exec(`
		INSERT INTO messages VALUE
		(?, ?, ?, ?)
		`,
			fromUser,
			toU,
			message,
			platform,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
