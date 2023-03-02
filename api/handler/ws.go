package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	api "github.com/intrntsrfr/vue-ws-test"
	"github.com/intrntsrfr/vue-ws-test/structs"

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
	ActionNone ActionCode = iota
	ActionUserReady
	ActionUserJoin
	ActionUserLeave
	ActionUserMessage
)

type WSEvent struct {
	Client *Client
	Event  *Event
}

// Event represents data send over the websocket
type Event struct {
	Operation OpCode          `json:"op"`
	RawData   json.RawMessage `json:"data"`
	Action    ActionCode      `json:"action"`
}

type sendEvent struct {
	Operation OpCode      `json:"op"`
	Data      interface{} `json:"data"`
	Action    ActionCode  `json:"action"`
}

type IdentifyData struct {
	Token string `json:"token"`
}

type PingData struct {
	Sequence int `json:"sequence"`
}

type ErrorData struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorCode int

const (
	UnknownError ErrorCode = iota
	PingTimedOut
	AuthFailed
)

var (
	ErrUnknownError = errors.New("unknown error")
	ErrPingTimedOut = errors.New("no ping for too long")
	ErrInvalidData  = errors.New("invalid data")
	ErrAuthFailed   = errors.New("authentication failed")
)

var ErrNoSuchError = errors.New("no such error")

// Client represents a connected websocket client
type Client struct {
	User       *structs.User
	Conn       *websocket.Conn
	Identified bool
	LastPing   time.Time
}

// Hub maintains a list of connected clients and broadcasts messages to them
type Hub struct {
	Clients    []*Client
	Messages   []*structs.Message
	EventCh    chan *WSEvent
	Register   chan *Client
	Unregister chan *Client
	db         database.DB
	jwt        api.JWTService
}

// NewHub returns a default Hub
func NewHub(db database.DB, jwt api.JWTService) *Hub {
	hub := &Hub{
		Clients:    []*Client{},
		Messages:   []*structs.Message{},
		EventCh:    make(chan *WSEvent),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		db:         db,
		jwt:        jwt,
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

		// &structs.User{ID: uuid.New(), Username: username, Created: time.Now().Format(time.RFC3339)}
		client := &Client{User: nil, Conn: conn, Identified: false, LastPing: time.Now()}
		h.Register <- client

		for {
			var evt Event
			err := conn.ReadJSON(&evt)
			if err != nil {
				fmt.Println(err)
				break
			}

			h.EventCh <- &WSEvent{
				Client: client,
				Event:  &evt,
			}
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
		case evt := <-h.EventCh:
			h.onEvent(evt)
		}
	}
}

func (h *Hub) onEvent(evt *WSEvent) {
	switch evt.Event.Operation {
	case Identify:
		data := IdentifyData{}
		if err := json.Unmarshal(evt.Event.RawData, &data); err != nil {
			return
		}
		h.identifyClient(evt.Client, &data)
	case Ping:
		data := PingData{}
		if err := json.Unmarshal(evt.Event.RawData, &data); err != nil {
			return
		}
		h.handlePing(evt.Client, &data)
	}
}

func (h *Hub) identifyClient(client *Client, evt *IdentifyData) {
	token, err := h.jwt.ParseToken(evt.Token)
	if err != nil {
		_ = h.disconnectClient(client, AuthFailed)
		return
	}
	claims, ok := token.(*api.UserClaims)
	if !ok {
		_ = h.disconnectClient(client, AuthFailed)
		return
	}

	user := h.db.FindUserByID(claims.Subject)
	if user == nil {
		_ = h.disconnectClient(client, AuthFailed)
		return
	}
	userCopy := *user
	userCopy.Password = ""

	client.Identified = true
	client.User = &userCopy
	_ = h.dispatchEvent(ActionUserReady, client, nil)
}

func (h *Hub) handlePing(client *Client, evt *PingData) {

}

type DispatchEvent func(conn *Client, data interface{}) error

func (h *Hub) dispatchEvent(ac ActionCode, conn *Client, data interface{}) error {
	fmt.Println("dispatching event, code:", ac)

	var dpe DispatchEvent
	switch ac {
	case ActionUserReady:
		dpe = h.userReady
	case ActionUserJoin:
		dpe = h.userJoin
	case ActionUserLeave:
		dpe = h.userLeave
	case ActionUserMessage:
		dpe = h.userMessage
	}
	err := dpe(conn, data)
	if err != nil {
		fmt.Println(err)
	}
	return err
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
	h.dispatchEvent(ActionUserLeave, nil, client.User)
}

func getError(code ErrorCode) (error, error) {
	switch code {
	case PingTimedOut:
		return ErrPingTimedOut, nil
	case AuthFailed:
		return ErrAuthFailed, nil
	}
	return nil, ErrNoSuchError
}

func (h *Hub) disconnectClient(client *Client, code ErrorCode) error {
	defer func(client *Client) {
		_ = client.Conn.Close()
	}(client)

	msgErr, err := getError(code)
	if err != nil {
		return err
	}

	data := &sendEvent{
		Operation: Error,
		Data:      &ErrorData{Code: code, Message: msgErr.Error()},
		Action:    ActionNone,
	}

	return client.Conn.WriteJSON(data)
}

func (h *Hub) broadcast(msg interface{}) error {
	// TODO: add subscription policy
	for _, client := range h.Clients {
		/*
			if time.Since(client.LastPing) > time.Second*15 {
				h.disconnectClient(client, PingTimedOut)
			}
		*/
		if client.Identified {
			_ = client.Conn.WriteJSON(msg)
		}
	}
	return nil
}

func (h *Hub) userReady(c *Client, _ interface{}) error {
	msgs := h.db.GetRecentMessages(50)
	users := make([]*structs.User, 0)
	for _, client := range h.Clients {
		if client.Identified && client != c {
			users = append(users, client.User)
		}
	}

	data := &sendEvent{
		Operation: Action,
		Data:      &structs.UserReady{msgs, users},
		Action:    ActionUserReady,
	}
	_ = c.Conn.WriteJSON(data)
	_ = h.dispatchEvent(ActionUserJoin, nil, c.User)
	return nil
}

func (h *Hub) userJoin(_ *Client, data interface{}) error {
	d, ok := data.(*structs.User)
	if !ok {
		return ErrInvalidData
	}

	d2 := &sendEvent{
		Operation: Action,
		Data:      &structs.UserJoin{d},
		Action:    ActionUserJoin,
	}
	return h.broadcast(d2)
}

func (h *Hub) userLeave(_ *Client, data interface{}) error {
	d, ok := data.(*structs.User)
	if !ok {
		return ErrInvalidData
	}

	d2 := &sendEvent{
		Operation: Action,
		Data:      &structs.UserLeave{d},
		Action:    ActionUserLeave,
	}
	return h.broadcast(d2)
}

func (h *Hub) userMessage(_ *Client, data interface{}) error {
	d, ok := data.(*structs.Message)
	if !ok {
		return ErrInvalidData
	}

	d2 := &sendEvent{
		Operation: Action,
		Data:      &structs.UserMessage{d},
		Action:    ActionUserMessage,
	}
	return h.broadcast(d2)
}
