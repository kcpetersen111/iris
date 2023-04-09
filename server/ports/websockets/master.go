package websockets

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/kcpetersen111/iris/server/persist"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Data about who is calling who
	metaData InitInfo

	//inbound data
	receive chan []byte
}

func NewHub(db *persist.DBInterface) *Hub {
	h := &Hub{
		db:         db,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			//figure out who they are suppose to be talking to
			for c := range h.clients {
				//if someone is here waiting for you and you try to join them
				// if c.metaData.Callee == client.metaData.Caller && c.metaData.Caller == client.metaData.Callee {
				//connect these two

				//for testing
				if c.metaData.Caller != client.metaData.Caller {

					c.send = client.receive
					client.send = c.receive
					fmt.Println("call started")
				}
				//if client is the first one here
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.receive)
			}
			// case message := <-h.broadcast:
			// 	for client := range h.clients {
			// 		select {
			// 		case client.send <- message:
			// 		default:
			// 			close(client.send)
			// 			delete(h.clients, client)
			// 		}
			// 	}
		}
	}
}
