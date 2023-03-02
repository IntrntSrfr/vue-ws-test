package structs

import (
	"github.com/google/uuid"
	"time"
)

// Message represents a message sent over the websocket
type Message struct {
	ID        uuid.UUID   `json:"id"`
	Author    *User       `json:"author"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
	Reactions []*Reaction `json:"reactions,omitempty"`
}

type Messages []*Message

func (m Messages) Len() int      { return len(m) }
func (m Messages) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

type ByTime struct{ Messages }

func (m ByTime) Less(i, j int) bool { return m.Messages[i].Timestamp.Before(m.Messages[j].Timestamp) }

// Reaction represents a message reaction
type Reaction struct {
	Emoji rune    `json:"emoji"`
	Users []*User `json:"users"`
}

// User represents a user
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password,omitempty"`
	Created  time.Time `json:"created"`
}
