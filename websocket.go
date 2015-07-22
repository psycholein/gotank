package main

import (
	"encoding/json"
	"fmt"
	"gotank/event"
	"gotank/libs/websocket"
	"net/http"
	"time"
)

const TIMEOUT = 5 * time.Second

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

type message struct {
	connection *connection
	message    []byte
}

type connectionPool struct {
	connections map[*connection]bool
	receive     chan *message
	register    chan *connection
	unregister  chan *connection
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	pool = connectionPool{
		connections: make(map[*connection]bool),
		receive:     make(chan *message),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
)

func startConnectionHandler() {
	go pool.run()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "No Websocket Connection", http.StatusBadRequest)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	pool.register <- c
	go c.writePump()
	c.readPump()
}

func (c *connection) readPump() {
	defer func() {
		pool.unregister <- c
		c.ws.Close()
	}()

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		r := &message{connection: c, message: msg}
		pool.receive <- r
	}
}

func (c *connection) writePump() {
	defer c.ws.Close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(TIMEOUT))
	return c.ws.WriteMessage(mt, payload)
}

func (pool *connectionPool) run() {
	for {
		select {
		case c := <-pool.register:
			pool.connections[c] = true
			fmt.Println("connected", c)
			sendModulesToWeb(c)
		case c := <-pool.unregister:
			if _, ok := pool.connections[c]; ok {
				delete(pool.connections, c)
				close(c.send)
			}
		case r := <-pool.receive:
			handleReceive(r)
		}
	}
}

func sendDataToAll(e event.Event) {
	data, err := json.Marshal(e)
	if err != nil {
		return
	}
	for c := range pool.connections {
		c.send <- data
	}
}

func sendData(c *connection, e event.Event) {
	data, err := json.Marshal(e)
	if err != nil {
		return
	}
	c.send <- data
}

func handleReceive(m *message) {
	var e event.Event
	if err := json.Unmarshal(m.message, &e); err != nil {
		return
	}
	go e.SendEventToLocal()
}
