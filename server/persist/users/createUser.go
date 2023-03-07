package users

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kcpetersen111/iris/server/persist"
)

type User struct {
	UserID   string `json:"userID"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u User) CreateUser(db *persist.DBInterface) (uuid.UUID, error) {
	fmt.Println(u)
	if u.Email == "" || u.Name == "" || u.Password == "" {
		return uuid.Nil, fmt.Errorf("Email, Name, or password is empty")
	}
	uu := uuid.New()
	tx, err := db.Database.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	_, err = tx.Query(`
	INSERT INTO users VALUES
	(?, ?, ?, ?, ? );
	`,
		uu,
		u.Name,
		u.Email,
		u.Role,
		u.Password,
	)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	tx.Commit()

	return uu, nil

}

func (u User) GetUserByEmail(db *persist.DBInterface) (User, error) {
	tx, err := db.Database.Begin()
	if err != nil {
		return User{}, err
	}
	row, err := tx.Query(`
		SELECT UserId, Name, Email, Role
		FROM users
		WHERE ? = Email;
	`, u.Email)

	if err != nil {
		tx.Rollback()
		return User{}, err
	}
	defer row.Close()
	var userId, name, email, role string
	//should only return one but this way of doing it looks a bit cleaner
	for row.Next() {
		err = row.Scan(&userId, &name, &email, &role)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return User{}, err
		}

	}

	tx.Commit()
	return User{
		UserID: userId,
		Name:   name,
		Email:  email,
		Role:   role,
	}, nil
}

func (u Authentication) GetUserByEmail(db *persist.DBInterface) (User, error) {
	tx, err := db.Database.Begin()
	if err != nil {
		return User{}, err
	}
	row, err := tx.Query(`
		SELECT UserId, Name, Email, Role, Password
		FROM users
		WHERE ? = Email;
	`, u.Email)

	if err != nil {
		tx.Rollback()
		return User{}, err
	}
	defer row.Close()
	var userId, name, email, role, password string
	//should only return one but this way of doing it looks a bit cleaner
	for row.Next() {
		err = row.Scan(&userId, &name, &email, &role, &password)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return User{}, err
		}

	}

	tx.Commit()
	return User{
		UserID:   userId,
		Name:     name,
		Email:    email,
		Role:     role,
		Password: password,
	}, nil
}
