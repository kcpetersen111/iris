package users

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type MessageJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateUser(username, password string) {

}

func CreateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	var input MessageJson
	rawRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request: %v", err)
	}
	err = json.Unmarshal(rawRequest, &input)
	if err != nil {
		log.Printf("Error unmarshaling request: %v", err)
	}
	CreateUser(input.Username, input.Password)
	w.WriteHeader(http.StatusOK)
}
