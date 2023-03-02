package structs

// UserReady is data that is sent to the user when the server has loaded their data
type UserReady struct {
	Messages []*Message `json:"messages"`
	Users    []*User    `json:"users"`
}

// UserJoin is the data to be sent when a user joins
type UserJoin struct {
	*User `json:"user"`
}

// UserLeave is the data to be sent when a user leaves
type UserLeave struct {
	*User `json:"user"`
}

// UserMessage is the data to be sent when a user sends a message
type UserMessage struct {
	*Message `json:"message"`
}
