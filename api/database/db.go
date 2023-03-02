package database

import (
	"github.com/intrntsrfr/vue-ws-test/structs"
)

type DB interface {
	CreateUser(u *structs.User) (*structs.User, error)
	FindUserByID(id string) *structs.User
	FindUserByUsername(username string) *structs.User

	CreateMessage(message *structs.Message) (*structs.Message, error)
	GetRecentMessages(limit int) []*structs.Message

	CreateReaction(messageID string, emoji rune) error
	DeleteReaction(messageID string, emoji rune) error
}
