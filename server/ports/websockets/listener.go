package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
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
	maxMessageSize = 4096
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
	Caller string `json:"Caller"`
	Callee string `json:"Callee"`
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
	if err != nil {
		log.Printf(fmt.Sprintf("Incoming call error in reading: %v", err))
		conn.Close()
		return
	}
	if msgType != 1 {
		log.Printf("Incoming call send incorect opening message\n")
		conn.Close()
		return
	}
	var init InitInfo
	err = json.Unmarshal(msg, &init)

	if err != nil {
		log.Println(fmt.Sprintf("Error in reading the startup metadata: %v", err))
		conn.Close()
		return
	}
	requestRetry := 0

persistCallSetup:
	err = hub.db.RequestCall(init.Callee, init.Caller)
	if err != nil {
		if requestRetry < 5 {
			requestRetry++
			goto persistCallSetup
		} else {
			log.Printf("Error in requesting call: %v", err)
			return
		}
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	client := &Client{hub: hub, conn: conn, send: nil, metaData: init, receive: make(chan []byte, 256), wg: &wg}
	client.hub.register <- client

	wg.Wait()
	go client.handleInput()
	go client.handleOutput()
}

func (c *Client) handleInput() {
	defer func() {
		log.Println("Closing input from client")
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
			log.Println(fmt.Sprintf("Error reading message, closing Input, Error: %v", err))
			break
		}
		if c.send != nil {
			c.send <- msg
		}
	}
}

func (c *Client) handleOutput() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		// c.receive = nil
		// c.send = nil
		log.Println("logging call end")
		c.ctxCancel()
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
			fmt.Println("pinging client")
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error pinging client")
				return
			}
		case <-c.ctx.Done():
			log.Println("Context done")
			return

		}
	}
}
