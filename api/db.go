package api

type DB interface {
	CreateUser(u *User) (*User, error)
	FindUserByID(id string) *User
	FindUserByUsername(username string) *User

	CreateMessage(m *Message) (*Message, error)
	GetMessages() ([]*Message, error)

	CreateReaction(messageID string, emoji rune) error
	DeleteReaction(messageID string, emoji rune) error
}
