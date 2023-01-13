package handler

import (
	"fmt"
	"github.com/intrntsrfr/vue-ws-test"
	"github.com/intrntsrfr/vue-ws-test/database"
	"github.com/intrntsrfr/vue-ws-test/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	User  *api.User
	Conn  *websocket.Conn
	Ready bool
}

// Hub maintains a list of connected clients and broadcasts messages to them
type Hub struct {
	Clients    []*Client
	Messages   []*api.Message
	Message    chan *api.Message
	Register   chan *Client
	Unregister chan *Client
	db         database.DB
}

// NewHub returns a default Hub
func NewHub(db database.DB) *Hub {
	hub := &Hub{
		Clients:    []*Client{},
		Messages:   []*api.Message{},
		Message:    make(chan *api.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		db:         db,
	}
	return hub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handler returns the websocket handler for the Hub
func (h *Hub) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		fmt.Println("new connection")

		username := c.Query("username")
		if username == "" {
			return
		}

		// ideally grab a user from some db?
		client := &Client{User: &api.User{ID: uuid.New(), Username: username, Created: time.Now().Format(time.RFC3339)}, Conn: conn}
		h.Register <- client // this will move later

		for {
			var message api.Message
			err := conn.ReadJSON(&message)
			if err != nil {
				fmt.Println(err)
				break
			}
			message.Author = client.User
			message.Timestamp = time.Now().Format(time.RFC3339)
			h.Message <- &message
		}

		h.Unregister <- client
	}
}

// Listen starts a loop for reading events that come through
func (h *Hub) Listen() {
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
	Messages []*api.Message `json:"messages"`
	Users    []*api.User    `json:"users"`
}

func (h *Hub) userConnected(c *Client) {
	msgs := h.Messages[len(h.Messages)-util.Min(len(h.Messages), 50):]
	users := []*api.User{}
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
	User *api.User `json:"user"`
}

func (h *Hub) userJoin(user *api.User) {
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
	User *api.User `json:"user"`
}

func (h *Hub) userLeave(user *api.User) {
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
	Message *api.Message `json:"message"`
}

func (h *Hub) userMessage(msg *api.Message) {
	data := Packet{
		Op: UserMessage,
		Data: UserMessageData{
			Message: msg,
		},
	}
	h.broadcast(data)
}
