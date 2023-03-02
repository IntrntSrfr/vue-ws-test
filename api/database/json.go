package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/intrntsrfr/vue-ws-test/structs"
)

type JsonDB struct {
	path  string
	state *state
}

type state struct {
	sync.Mutex
	Users    map[string]*structs.User    `json:"users"`
	Messages map[string]*structs.Message `json:"messages"`
}

func Open(path string) (*JsonDB, error) {
	var (
		db  *JsonDB
		err error
	)
	db = &JsonDB{
		path: path,
		state: &state{
			Users:    make(map[string]*structs.User, 0),
			Messages: make(map[string]*structs.Message, 0),
		},
	}
	if path != "" {
		err = db.load(path)
	}
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

	state := &state{}
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

func (j *JsonDB) CreateUser(u *structs.User) (*structs.User, error) {
	j.state.Lock()
	defer j.state.Unlock()
	j.state.Users[u.ID.String()] = u
	return u, nil
}

func (j *JsonDB) FindUserByID(id string) *structs.User {
	j.state.Lock()
	defer j.state.Unlock()
	return j.state.Users[id]
}

func (j *JsonDB) FindUserByUsername(username string) *structs.User {
	j.state.Lock()
	defer j.state.Unlock()
	for _, u := range j.state.Users {
		if u.Username == username {
			return u
		}
	}
	return nil
}

func (j *JsonDB) CreateMessage(message *structs.Message) (*structs.Message, error) {
	j.state.Lock()
	defer j.state.Unlock()
	j.state.Messages[message.ID.String()] = message
	return message, nil
}

func (j *JsonDB) GetRecentMessages(limit int) []*structs.Message {
	j.state.Lock()
	defer j.state.Unlock()

	var messages []*structs.Message
	for _, msg := range j.state.Messages {
		messages = append(messages, msg)
	}
	sort.Sort(structs.ByTime{Messages: messages})

	recent := make([]*structs.Message, 0)
	for i, msg := range messages {
		if i >= limit {
			break
		}
		recent = append(recent, msg)
	}
	return recent
}

func (j *JsonDB) CreateReaction(messageID string, emoji rune) error {
	//TODO implement me
	panic("implement me")
}

func (j *JsonDB) DeleteReaction(messageID string, emoji rune) error {
	//TODO implement me
	panic("implement me")
}
