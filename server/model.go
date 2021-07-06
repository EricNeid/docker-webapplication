package server

import (
	"encoding/json"
	"io"

	"github.com/paulmach/orb"
)

type Position struct {
	Position orb.Point
}

type User struct {
	Name string
}

type ResponseUserId struct {
	UserId int64
}

type ResponseUser struct {
	User User
}

type ResponseUsers struct {
	Users []User
}

func NewUser(in io.Reader) (User, error) {
	var model User
	err := json.NewDecoder(in).Decode(&model)
	return model, err
}

func NewResponseUserId(in io.Reader) (ResponseUserId, error) {
	var model ResponseUserId
	err := json.NewDecoder(in).Decode(&model)
	return model, err
}

func NewResponseUser(in io.Reader) (ResponseUser, error) {
	var model ResponseUser
	err := json.NewDecoder(in).Decode(&model)
	return model, err
}

func NewResponseUsers(in io.Reader) (ResponseUsers, error) {
	var model ResponseUsers
	err := json.NewDecoder(in).Decode(&model)
	return model, err
}

func (model ResponseUserId) ToJson() string {
	json, _ := json.Marshal(model)
	return string(json)
}

func (model User) ToJson() string {
	json, _ := json.Marshal(model)
	return string(json)
}
