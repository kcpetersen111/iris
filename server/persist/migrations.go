package persist

//This will do all of the first time set up for the database
//doing things like creating tables
func (db *DBInterface) DbSetup() error {

	tx, err := db.Database.Begin()
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

	_, err = tx.Exec(
		`CREATE TABLE users(
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
	_, err = tx.Exec(
		`CREATE TABLE messages(
			sender TEXT,
			receiver TEXT,
			message TEXT,
			platform TEXT,
			call INT
			);`,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}
