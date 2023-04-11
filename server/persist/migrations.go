package persist

import "log"

//This will do all of the first time set up for the database
//doing things like creating tables
func (db *DBInterface) DbSetup() error {
	log.Println("Setting up database")

	tx, err := db.Database.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS test(
			testfield TEXT,
			numberfield INT
			);`,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS users(
			userID VARCHAR(45),
			name TEXT,
			email TEXT,
			role TEXT,
			password TEXT,
			PRIMARY KEY(userID)
			);`,
	)
	if err != nil {
		return err
	}

	//I am going to say platform is where they sent it ie: discord server, dm
	//will probably need to add message id at some point if I want to add editing or deleting messages
	//should take out receiver and just have that be the platform
	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS messages(
			sender VARCHAR(45),
			receiver VARCHAR(45),
			message TEXT,
			platform VARCHAR(45),
			isCall INT,
			timeStamp DATETIME
			);`,
	)
	if err != nil {
		return err
	}

	//I am going to say platform is where they sent it ie: discord server, dm

	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS platforms(
			platformID VARCHAR(45),
			platformName TEXT,
			userId VARCHAR(45)
			);`,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}
