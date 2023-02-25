package handler

import (
	"errors"
	"fmt"
	"net/http"

	api "github.com/intrntsrfr/vue-ws-test"
	"github.com/intrntsrfr/vue-ws-test/database"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type OpCode int

const (
	Identify OpCode = iota
	Ping
	PingACK
	Action
	Error
)

type ActionCode int

const (
	UserReady ActionCode = iota
	UserJoin
	UserLeave
	UserMessage
)

type WSEvent struct {
	Client *Client
	Event  *Event
}

// Event represents data send over the websocket
type Event struct {
	Op   OpCode      `json:"op"`
	Data interface{} `json:"data"`
}

func NewEvent(op OpCode, data interface{}) *Event {
	return &Event{
		Op:   op,
		Data: data,
	}
}

type IdentifyData struct {
	Token string `json:"token"`
}

type PingData struct {
	Sequence int `json:"sequence"`
}

type ErrorData struct {
	Code int `json:"code"`
}

// Client represents a connected websocket client
type Client struct {
	User       *api.User
	Conn       *websocket.Conn
	Identified bool
}

// Hub maintains a list of connected clients and broadcasts messages to them
type Hub struct {
	Clients    []*Client
	Messages   []*api.Message
	EventCh    chan *WSEvent
	Register   chan *Client
	Unregister chan *Client
	db         database.DB
}

// NewHub returns a default Hub
func NewHub(db database.DB) *Hub {
	hub := &Hub{
		Clients:    []*Client{},
		Messages:   []*api.Message{},
		EventCh:    make(chan *WSEvent),
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

		// &api.User{ID: uuid.New(), Username: username, Created: time.Now().Format(time.RFC3339)}
		client := &Client{User: nil, Conn: conn, Identified: false}
		h.Register <- client

		for {
			var evt WSEvent
			err := conn.ReadJSON(&evt)
			if err != nil {
				fmt.Println(err)
				break
			}
			h.EventCh <- &evt
		}

		h.Unregister <- client
	}
}

// Run starts a loop for reading events that come through
func (h *Hub) Run() {
	//h.heartbeats()
	h.listenEvents()
}

func (h *Hub) listenEvents() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.removeClient(client)
			h.userLeave(nil, client.User)
		case evt := <-h.EventCh:
			if evt.Event.Op == Identify {
				if e, ok := evt.Event.Data.(*IdentifyData); ok {
					h.identifyClient(evt.Client, e)
				}
			} else if evt.Event.Op == PingACK {
				if e, ok := evt.Event.Data.(*PingData); ok {
					h.handlePingACK(evt.Client, e)
				}
			}
		}
	}
}

func (h *Hub) identifyClient(client *Client, evt *IdentifyData) {
	// check JWT and check againt database
	// if good, set as identified, otherwise return error data
}

func (h *Hub) handlePingACK(client *Client, evt *PingData) {

}

// UserReady ActionCode = iota
// UserJoin
// UserLeave
// UserMessage

type DispatchEvent func(conn *Client, data interface{}) error

func (h *Hub) dispatchEvent(ac ActionCode, conn *Client, data interface{}) error {
	var dpe DispatchEvent
	switch ac {
	case UserReady:
		dpe = h.userReady
	case UserJoin:
		dpe = h.userJoin
	case UserLeave:
		dpe = h.userLeave
	case UserMessage:
		dpe = h.userMessage
	}
	return dpe(conn, data)
}

func (h *Hub) registerClient(client *Client) {
	h.Clients = append(h.Clients, client)
}

func (h *Hub) removeClient(client *Client) {
	for i, c := range h.Clients {
		if c == client {
			h.Clients = append(h.Clients[:i], h.Clients[i+1:]...)
			break
		}
	}
}

func (h *Hub) broadcast(msg interface{}) error {
	// TODO: add subscription policy
	for _, client := range h.Clients {
		if client.Identified {
			client.Conn.WriteJSON(msg)
		}
	}
	return nil
}

// UserReadyData is data that is sent to the user when the server has loaded their data
type UserReadyData struct {
	Code     ActionCode
	Messages []*api.Message `json:"messages"`
	Users    []*api.User    `json:"users"`
}

var (
	ErrInvalidData = errors.New("Invalid data")
)

func (h *Hub) userReady(c *Client, data interface{}) error {
	/*
		msgs := h.Messages[len(h.Messages)-util.Min(len(h.Messages), 50):]
		users := []*api.User{}
		for _, client := range h.Clients {
			if client.Identified {
				users = append(users, client.User)
			}
		}
	*/
	// send all clients and last 50 messages
	// some opcode
	// after this let everyone else know the user joined
	//c.Conn.WriteJSON()

	d, ok := data.(*UserReadyData)
	if !ok {
		return ErrInvalidData
	}

	d2 := Event{
		Op:   Action,
		Data: d,
	}
	return c.Conn.WriteJSON(d2)
}

// UserJoinData is the data to be sent when a user joins
type UserJoinData struct {
	Code ActionCode
	User *api.User `json:"user"`
}

func (h *Hub) userJoin(conn *Client, data interface{}) error {
	d, ok := data.(*UserJoinData)
	if !ok {
		return ErrInvalidData
	}

	d2 := Event{
		Op:   Action,
		Data: d,
	}
	return h.broadcast(d2)
}

// UserLeaveData is the data to be sent when a user leaves
type UserLeaveData struct {
	Code ActionCode
	User *api.User `json:"user"`
}

func (h *Hub) userLeave(conn *Client, data interface{}) error {
	d, ok := data.(*UserLeaveData)
	if !ok {
		return ErrInvalidData
	}

	d2 := Event{
		Op:   Action,
		Data: d,
	}
	return h.broadcast(d2)
}

// UserMessageData is the data to be sent when a user sends a message
type UserMessageData struct {
	Message *api.Message `json:"message"`
}

func (h *Hub) userMessage(conn *Client, data interface{}) error {
	d, ok := data.(*UserMessageData)
	if !ok {
		return ErrInvalidData
	}

	d2 := Event{
		Op:   Action,
		Data: d,
	}
	return h.broadcast(d2)
}
