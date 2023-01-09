package api

import "github.com/google/uuid"

// Message represents a message sent over the websocket
type Message struct {
	Username  string      `json:"username"`
	Text      string      `json:"text"`
	Timestamp string      `json:"timestamp"`
	Reactions []*Reaction `json:"reactions"`
}

// Reaction represents a message reaction
type Reaction struct {
	Emoji rune    `json:"emoji"`
	Users []*User `json:"users"`
}

// User represents a user
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Created  string    `json:"created"`
}
