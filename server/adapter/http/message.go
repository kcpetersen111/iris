package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kcpetersen111/iris/server/persist"
)

type MessageRoutes struct {
	DB *persist.DBInterface
}

type MessagePlatform struct {
	Platform string `json:"platformID"`
}

type Message struct {
	Message  string `json:"message"`
	Sender   string `json:"sender"`
	Platform string `json:"platform"`
}

func (m *MessageRoutes) GetMessage(w http.ResponseWriter, r *http.Request) {
	var platId MessagePlatform
	err := json.NewDecoder(r.Body).Decode(&platId)
	if err != nil {
		err := fmt.Errorf("Error in reading body: %v", err)
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	msgs, err := m.DB.GetMessages(platId.Platform)
	if err != nil {
		err := fmt.Errorf("Error in getting messages: %v", err)
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(msgs); err != nil {
		log.Printf("%v\n", err)
	}

}

func (m *MessageRoutes) PostMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		err := fmt.Errorf("Error in reading body: %v", err)
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	err = m.DB.InsertMessage(msg.Message, msg.Sender, msg.Platform)
	if err != nil {
		err := fmt.Errorf("Error in creating message: %v", err)
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
