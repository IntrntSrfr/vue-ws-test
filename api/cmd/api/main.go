package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type OpCode int

const (
	UserJoin OpCode = iota
	UserLeave
	UserMessage
	UserPing
	UserReady
)

type Packet struct {
	Op   int         `json:"op"`
	Data interface{} `json:"data"`
}

// Message is a struct that represents a message sent over the websocket
type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

// Client is a struct that represents a connected websocket client
type Client struct {
	Username string
	Conn     *websocket.Conn
}

// Hub is a struct that maintains a list of connected clients and broadcasts messages to them
type Hub struct {
	Clients    []*Client
	Messages   []*Message
	Message    chan *Message
	Register   chan *Client
	Unregister chan *Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients = append(h.Clients, client)
			h.userJoin(client.Username)
		case client := <-h.Unregister:
			for i, c := range h.Clients {
				if c == client {
					h.Clients = append(h.Clients[:i], h.Clients[i+1:]...)
					break
				}
			}
			h.userLeave(client.Username)
		case message := <-h.Message:
			h.userMessage(message)
		}
	}
}

func (h *Hub) broadcast(msg interface{}) {
	for _, client := range h.Clients {
		client.Conn.WriteJSON(msg)
	}
}

func (h *Hub) userConnected(c *Client) {
	// send all clients and last 50 messages
	// some opcode
	// after this let everyone else know the user joined
	//c.Conn.WriteJSON()
}

func (h *Hub) userJoin(username string) {
	data := struct {
		Op       OpCode
		Username string
	}{Op: UserJoin, Username: username}
	h.broadcast(data)
}

func (h *Hub) userLeave(username string) {
	data := struct {
		Op       OpCode
		Username string
	}{Op: UserLeave, Username: username}
	h.broadcast(data)
}

func (h *Hub) userMessage(msg *Message) {
	data := struct {
		Op      OpCode
		Message *Message
	}{Op: UserMessage, Message: msg}
	h.broadcast(data)
}

func (h *Hub) websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("new connection")

	username := r.URL.Query().Get("username")
	if username == "" {
		return
	}

	client := &Client{Username: username, Conn: conn}
	h.Register <- client

	for {
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			fmt.Println(err)
			break
		}
		h.Message <- &message
	}

	h.Unregister <- client
}

func main() {
	hub := Hub{
		Message:    make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	go hub.run()

	http.HandleFunc("/ws", hub.websocketHandler)
	http.ListenAndServe(":8080", nil)
}
