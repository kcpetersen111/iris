package persist

import (
	"fmt"
	"log"
)

type Platform struct {
	PlatformID   string `json:"platformID"`
	PlatformName string `json:"platformName"`
}

func (db *DBInterface) GetPlatform(userID string) ([]Platform, error) {
	tx, err := db.Database.Begin()
	if err != nil {
		log.Println(fmt.Sprintf("Error starting transaction: %v", err))
	}
	row, err := tx.Query(`
		SELECT 
			platformID,
			platfromName
		FROM platforms
		WHERE
			? = platformID
		ORDER BY platformName;
	`)
	if err != nil {
		log.Println(fmt.Sprintf("Error in sql request: %v", err))
		tx.Rollback()
		return nil, err
		// return User{}, err
	}
	defer row.Close()
	var pId, pN string
	var platforms []Platform
	for row.Next() {
		err = row.Scan(&pId, &pN)
		if err != nil {
			log.Println(fmt.Sprintf("Error in reading sql response: %v", err))
			tx.Rollback()
			return nil, err
		}
		platforms = append(platforms, Platform{
			PlatformID:   pId,
			PlatformName: pN,
		})
	}
	tx.Commit()
	return platforms, nil
}
