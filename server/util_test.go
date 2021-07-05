package server

import (
	"encoding/json"
	"io"

	"github.com/EricNeid/go-webserver/model"
)

func NewUser(in io.Reader) (model.User, error) {
	var user model.User
	err := json.NewDecoder(in).Decode(&user)
	return user, err
}

func ToJson(res model.ResponseUserId) []byte {
	json, _ := json.Marshal(res)
	return json
}

func ToJson(user model.User) string {
	json, _ := json.Marshal(user)
	return string(json)
}
