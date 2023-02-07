package users

import (
	"github.com/gorilla/mux"
	"github.com/kcpetersen111/iris/server/persist"
)

type UserRoutes struct {
	DB *persist.DBInterface
}

func CreateUserRouter(db *persist.DBInterface) {
	userStruct := UserRoutes{
		DB: db,
	}
	router := mux.NewRouter()
}
