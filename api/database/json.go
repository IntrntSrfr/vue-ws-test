package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	api "github.com/intrntsrfr/vue-ws-test"
)

type JsonDB struct {
	path  string
	state *state
}

type state struct {
	sync.Mutex
	Users    []*api.User    `json:"users"`
	Messages []*api.Message `json:"messages"`
}

func Open(path string) (*JsonDB, error) {
	db := &JsonDB{
		path: path,
		state: &state{
			Users:    make([]*api.User, 0),
			Messages: make([]*api.Message, 0),
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

func (j *JsonDB) CreateUser(u *api.User) (*api.User, error) {
	j.state.Lock()
	defer j.state.Unlock()
	j.state.Users = append(j.state.Users, u)
	return u, nil
}

func (j *JsonDB) FindUserByID(id string) *api.User {
	j.state.Lock()
	defer j.state.Unlock()
	for _, u := range j.state.Users {
		if u.ID.String() == id {
			return u
		}
	}
	return nil
}

func (j *JsonDB) FindUserByUsername(username string) *api.User {
	j.state.Lock()
	defer j.state.Unlock()
	for _, u := range j.state.Users {
		if u.Username == username {
			return u
		}
	}
	return nil
}
