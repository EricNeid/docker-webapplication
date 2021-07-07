package server

import (
	"encoding/json"
	"io"
)

func newResponseUserId(in io.Reader) (responseUserId, error) {
	var model responseUserId
	err := json.NewDecoder(in).Decode(&model)
	return model, err
}

func newResponseUser(in io.Reader) (responseUser, error) {
	var model responseUser
	err := json.NewDecoder(in).Decode(&model)
	return model, err
}

func newResponseUsers(in io.Reader) (responseUsers, error) {
	var model responseUsers
	err := json.NewDecoder(in).Decode(&model)
	return model, err
}

func (model user) toJson() string {
	json, _ := json.Marshal(model)
	return string(json)
}
