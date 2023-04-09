package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kcpetersen111/iris/server/persist"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	db         *persist.DBInterface
}

type InitInfo struct {
	Caller string
	Callee string
}

func StartCall(hub *Hub, w http.ResponseWriter, r *http.Request) {
	fmt.Println("New call")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	msgType, msg, err := conn.ReadMessage()
	if msgType != 1 {
		log.Printf("Incoming call send incorect opening message\n")
		conn.Close()
		return
	}
	var init InitInfo
	json.Unmarshal(msg, &init)
	persist.RequestCall(init.Callee, init.Caller)

	client := &Client{hub: hub, conn: conn, send: nil, metaData: init, receive: make(chan []byte, 256)}
	client.hub.register <- client

	//have incoming call

	//write to db to say there is a call waiting for client 2

	//wait 30 secs for the them to pick up if not close the connection

	// go client.handleInput()
	go client.handleInput()
}

func (c *Client) handleInput() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if c.send != nil {
			c.send <- msg
		}
		// fmt.Println(msg)
		// err = c.conn.WriteMessage(2, msg)
		// if err != nil {
		// 	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 		log.Printf("error: %v", err)
		// 	}
		// 	break
		// }

		// c.conn.UnderlyingConn().Write(message)

		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// c.hub.broadcast <- message
	}
}

func (c *Client) handleOutput() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.receive:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Printf("Error in creating a new writer: %v", err)
				return
			}
			w.Write(message)

			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
