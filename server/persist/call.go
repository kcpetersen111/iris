package persist

import (
	"fmt"
	"time"
)

func (db *DBInterface) RequestCall(callerId, calleeId string) error {
	tx, err := db.Database.Begin()
	if err != nil {
		return err
	}
	tx.Exec(`
		INSERT INTO messages VALUE
		(?,?,?,"",1);
	`,
		calleeId,
		calleeId,
		fmt.Sprintf("%v", time.Now()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
