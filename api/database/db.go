package database

import (
	"github.com/intrntsrfr/vue-ws-test"
)

type DB interface {
	CreateUser(u *api.User) (*api.User, error)
	FindUserByID(id string) *api.User
	FindUserByUsername(username string) *api.User

	CreateMessage(m *api.Message) (*api.Message, error)
	GetMessages() []*api.Message

	CreateReaction(messageID string, emoji rune) error
	DeleteReaction(messageID string, emoji rune) error
}
