package api

type DB interface {
	CreateUser() *User
	GetUser() *User

	CreateMessage() *Message
	GetMessages() []*Message

	CreateReaction()
	DeleteReaction()
}
