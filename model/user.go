package model

import (
	"encoding/json"
	"io"
)

type User struct {
	Name string
}

type ResponseUserId struct {
	UserId int64
}

func NewUser(in io.Reader) (User, error) {
	var user User
	err := json.NewDecoder(in).Decode(&user)
	return user, err
}

func (res ResponseUserId) ToJson() []byte {
	json, _ := json.Marshal(res)
	return json
}
