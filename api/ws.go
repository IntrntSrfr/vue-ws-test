package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
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

// Packet represents data send over the websocket
type Packet struct {
	Op   OpCode      `json:"op"`
	Data interface{} `json:"data"`
}

// Client represents a connected websocket client
type Client struct {
	User *User
	Conn *websocket.Conn
}

// Hub maintains a list of connected clients and broadcasts messages to them
type Hub struct {
	Clients    []*Client
	Messages   []*Message
	Message    chan *Message
	Register   chan *Client
	Unregister chan *Client
}

// NewHub returns a default Hub
func NewHub() *Hub {
	hub := &Hub{
		Clients:    []*Client{},
		Messages:   []*Message{},
		Message:    make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	return hub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handler returns the websocket handler for the Hub
func (h *Hub) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// ideally grab a user from some db?
		client := &Client{User: &User{ID: uuid.New(), Username: username, Created: time.Now().Format(time.RFC3339)}, Conn: conn}
		h.Register <- client // this will move later

		for {
			var message Message
			err := conn.ReadJSON(&message)
			if err != nil {
				fmt.Println(err)
				break
			}
			message.Username = username
			message.Timestamp = time.Now().Format(time.RFC3339)
			h.Message <- &message
		}

		h.Unregister <- client
	}
}

// Run starts a loop for reading events that come through
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.userConnected(client)
			h.Clients = append(h.Clients, client)
			h.userJoin(client.User)
		case client := <-h.Unregister:
			for i, c := range h.Clients {
				if c == client {
					h.Clients = append(h.Clients[:i], h.Clients[i+1:]...)
					break
				}
			}
			h.userLeave(client.User)
		case message := <-h.Message:
			fmt.Println(message)
			h.Messages = append(h.Messages, message)
			h.userMessage(message)
		}
	}
}

func (h *Hub) broadcast(msg interface{}) {
	for _, client := range h.Clients {
		client.Conn.WriteJSON(msg)
	}
}

// UserReadyData is data that is sent to the user when the server has loaded their data
type UserReadyData struct {
	Messages []*Message `json:"messages"`
	Users    []*User    `json:"users"`
}

func (h *Hub) userConnected(c *Client) {
	msgs := h.Messages[len(h.Messages)-min(len(h.Messages), 50):]
	users := []*User{}
	for _, client := range h.Clients {
		users = append(users, client.User)
	}

	// send all clients and last 50 messages
	// some opcode
	// after this let everyone else know the user joined
	//c.Conn.WriteJSON()

	data := Packet{
		Op: UserReady,
		Data: UserReadyData{
			Messages: msgs,
			Users:    users,
		},
	}
	c.Conn.WriteJSON(data)
}

// UserJoinData is the data to be sent when a user joins
type UserJoinData struct {
	User *User `json:"user"`
}

func (h *Hub) userJoin(user *User) {
	data := Packet{
		Op: UserJoin,
		Data: UserJoinData{
			User: user,
		},
	}
	h.broadcast(data)
}

// UserLeaveData is the data to be sent when a user leaves
type UserLeaveData struct {
	User *User `json:"user"`
}

func (h *Hub) userLeave(user *User) {
	data := Packet{
		Op: UserLeave,
		Data: UserLeaveData{
			User: user,
		},
	}
	h.broadcast(data)
}

// UserMessageData is the data to be sent when a user sends a message
type UserMessageData struct {
	Message *Message `json:"message"`
}

func (h *Hub) userMessage(msg *Message) {
	data := Packet{
		Op: UserMessage,
		Data: UserMessageData{
			Message: msg,
		},
	}
	h.broadcast(data)
}
