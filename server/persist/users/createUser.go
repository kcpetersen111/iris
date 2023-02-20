package users

import (
	"github.com/kcpetersen111/iris/server/adapter/http/user"
	"github.com/kcpetersen111/iris/server/persist"
)

func CreateUser(username, password string) {

}

func (u user.User) GetUserByEmail(db *persist.DBInterface) (user.User, error) {
	tx, err := db.Database.Begin()
	if err != nil {
		return nil, err
	}
	row, err := tx.Query(`
		SELECT UserId, Name, Email, Role
		FROM users
		WHERE ? = Email;
	`, u.Email)
	defer row.Close()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	row.Next()
	var UserId, name, email, role string
	row.Scan(&UserId, &name, &email, &role)

	return user.User{
		UserId: UserId,
		Name:   name,
		Email:  email,
		Role:   role,
	}, nil
}
