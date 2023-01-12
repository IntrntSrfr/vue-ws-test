package api

import (
	"fmt"
	"github.com/goccy/go-json"
	"os"
	"sync"
)

type DB interface {
	CreateUser(u *User) (*User, error)
	FindUserByID(id string) *User
	FindUserByUsername(username string) *User

	CreateMessage(m *Message) (*Message, error)
	GetMessages() []*Message

	CreateReaction(messageID string, emoji rune) error
	DeleteReaction(messageID string, emoji rune) error
}

type JsonDB struct {
	path  string
	state *State
}

type State struct {
	sync.Mutex
	Users    []*User    `json:"users"`
	Messages []*Message `json:"messages"`
}

func Open(path string) (*JsonDB, error) {
	db := &JsonDB{
		path: path,
		state: &State{
			Users:    make([]*User, 0),
			Messages: make([]*Message, 0),
		},
	}
	err := db.load(path)
	return db, err
}

func (j *JsonDB) Close() error {
	return j.save()
}

func (j *JsonDB) load(path string) error {
	if _, err := os.Stat(path); err != nil {
		// file does not exist, so use default
		fmt.Println("no data file found, using default")
		return nil
	}

	fmt.Println("data file found")
	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	state := &State{}
	err = json.Unmarshal(d, &state)
	if err != nil {
		return err
	}

	j.state = state
	return nil
}

func (j *JsonDB) save() error {
	d, err := json.Marshal(j.state)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(j.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(d)
	return err
}

func (j *JsonDB) CreateUser(u *User) (*User, error) {
	j.state.Lock()
	defer j.state.Unlock()
	j.state.Users = append(j.state.Users, u)
	return u, nil
}

func (j *JsonDB) FindUserByID(id string) *User {
	j.state.Lock()
	defer j.state.Unlock()
	for _, u := range j.state.Users {
		if u.ID.String() == id {
			return u
		}
	}
	return nil
}

func (j *JsonDB) FindUserByUsername(username string) *User {
	j.state.Lock()
	defer j.state.Unlock()
	for _, u := range j.state.Users {
		if u.Username == username {
			return u
		}
	}
	return nil
}

func (j *JsonDB) CreateMessage(m *Message) (*Message, error) {
	j.state.Lock()
	defer j.state.Unlock()
	j.state.Messages = append(j.state.Messages, m)
	return m, nil
}

func (j *JsonDB) GetMessages() []*Message {
	j.state.Lock()
	defer j.state.Unlock()
	return j.state.Messages[len(j.state.Messages)-min(len(j.state.Messages), 50):]
}

func (j *JsonDB) CreateReaction(messageID string, emoji rune) error {
	//TODO implement me
	panic("implement me")
}

func (j *JsonDB) DeleteReaction(messageID string, emoji rune) error {
	//TODO implement me
	panic("implement me")
}
