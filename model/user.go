package model

import (
	"encoding/json"
)

type User struct {
	Name string
}

type ResponseUserId struct {
	UserId int64
}
