package persist

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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

func (db *DBInterface) InsertMessage(message, fromUser, platform string) error {

	tx, err := db.Database.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	// for toU := range toUser {

	_, err = tx.Exec(`
		INSERT INTO messages VALUE
		(?, ?, ?, ?, 0,?);
		`,
		fromUser,
		"",
		message,
		platform,
		time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	// }

	tx.Commit()

	return nil
}

func (db *DBInterface) CreatePlatform(name string, emails []string) error {

	tx, err := db.Database.Begin()
	if err != nil {
		log.Printf("Error while starting the transaction: %v", err)
		return err
	}

	id := uuid.New().String()
	for _, email := range emails {

		row, err := db.Database.Query(`
		SELECT
			userID
		FROM users
		WHERE email = ?;
		`, email)
		if err != nil {
			tx.Rollback()
			log.Printf("Error while finding the userId that goes with an email: %v", err)
			return err
		}

		defer row.Close()
		var userID string
		b := row.Next()
		if b {
			err = row.Scan(&userID)
		} else {
			return fmt.Errorf("User does not exist")
		}
		if err != nil {
			log.Println(fmt.Sprintf("Error in reading sql response: %v", err))
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(`
		INSERT INTO platforms VALUE
		(?, ?, ?);
	`, id, name, userID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error while inserting the transaction: %v", err)
			return err
		}

	}
	tx.Commit()

	return nil
}
