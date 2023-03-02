package database

import (
	"github.com/google/uuid"
	"github.com/intrntsrfr/vue-ws-test/structs"
	"reflect"
	"testing"
	"time"
)

func TestJsonDB_CreateMessage(t *testing.T) {
	db, err := Open("")
	if err != nil {
		t.Errorf("encountered error: %v", err)
	}

	msg := &structs.Message{
		ID:        uuid.New(),
		Author:    nil,
		Content:   "message",
		Timestamp: time.Now(),
		Reactions: nil,
	}

	msg, err = db.CreateMessage(msg)
	if err != nil {
		t.Errorf("encountered error: %v", err)
	}

	if len(db.state.Messages) != 1 {
		t.Errorf("len(state.Messages) got %v, wanted %v", len(db.state.Messages), 1)
	}

	msg, ok := db.state.Messages[msg.ID.String()]
	if !ok {
		t.Errorf("state.Messages[msg.ID] was not found")
	}

	if msg.Content != "message" {
		t.Errorf("msg.Content got %v, wanted 'content'", msg.Content)
	}
}

func TestJsonDB_CreateUser(t *testing.T) {
	db, err := Open("")
	if err != nil {
		t.Errorf("encountered error: %v", err)
	}

	user := &structs.User{
		ID:       uuid.New(),
		Username: "jeff",
		Created:  time.Now(),
	}

	user, err = db.CreateUser(user)
	if err != nil {
		t.Errorf("encountered error: %v", err)
	}

	if len(db.state.Users) != 1 {
		t.Errorf("len(state.Users) got %v, wanted %v", len(db.state.Users), 1)
	}

	user, ok := db.state.Users[user.ID.String()]
	if !ok {
		t.Errorf("state.Users[user.ID] was not found")
	}

	if user.Username != "jeff" {
		t.Errorf("user.Content got %v, wanted 'jeff'", user.Username)
	}
}

func TestJsonDB_FindUserByID(t *testing.T) {
	type fields struct {
		path  string
		state *state
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *structs.User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JsonDB{
				path:  tt.fields.path,
				state: tt.fields.state,
			}
			if got := j.FindUserByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonDB_FindUserByUsername(t *testing.T) {
	type fields struct {
		path  string
		state *state
	}
	type args struct {
		username string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *structs.User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JsonDB{
				path:  tt.fields.path,
				state: tt.fields.state,
			}
			if got := j.FindUserByUsername(tt.args.username); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonDB_GetRecentMessages(t *testing.T) {
	type fields struct {
		path  string
		state *state
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*structs.Message
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JsonDB{
				path:  tt.fields.path,
				state: tt.fields.state,
			}
			if got := j.GetRecentMessages(tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecentMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}
