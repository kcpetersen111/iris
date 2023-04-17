package persist

import (
	"fmt"
	"log"
)

type Platform struct {
	PlatformID   string `json:"platformID"`
	PlatformName string `json:"platformName"`
}

type Message struct {
	Name      string `json:"name"`
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Platform  string `json:"platform"`
	TimeStamp string `json:"timestamp"`
}

func (db *DBInterface) GetPlatform(userID string) ([]Platform, error) {
	tx, err := db.Database.Begin()
	if err != nil {
		log.Println(fmt.Sprintf("Error starting transaction: %v", err))
		tx.Rollback()
		return nil, err
	}
	row, err := tx.Query(`
		SELECT 
			platformID,
			platformName
		FROM platforms
		WHERE
			? = userId
		ORDER BY platformName;
	`, userID)
	if err != nil {
		log.Println(fmt.Sprintf("Error in sql request: %v", err))
		tx.Rollback()
		return nil, err
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

func (db *DBInterface) GetMessages(platform string) ([]Message, error) {

	tx, err := db.Database.Begin()
	if err != nil {
		log.Println(fmt.Sprintf("Error starting transaction: %v", err))
		tx.Rollback()
		return nil, err
	}

	row, err := tx.Query(`
		SELECT
			users.email,
			messages.sender,
			messages.message,
			messages.platform,
			messages.timeStamp
		FROM messages
		INNER JOIN users ON users.userID = messages.sender
		WHERE
			messages.platform = ? AND
			messages.isCall = 0
		ORDER BY timeStamp
		LIMIT 200;
	`, platform)
	if err != nil {
		log.Println(fmt.Sprintf("Error in sql request: %v", err))
		tx.Rollback()
		return nil, err
	}
	defer row.Close()

	var messageList []Message
	var name, sender, message, plat, timestamp string
	for row.Next() {
		err = row.Scan(&name, &sender, &message, &plat, &timestamp)
		if err != nil {
			log.Println(fmt.Sprintf("Error in reading sql response: %v", err))
			tx.Rollback()
			return nil, err
		}
		messageList = append(messageList, Message{
			Name:      name,
			Sender:    sender,
			Message:   message,
			Platform:  plat,
			TimeStamp: timestamp,
		})
	}
	tx.Commit()

	return messageList, nil
}
