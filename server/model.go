package server

import (
	"github.com/paulmach/orb"
)

type vehicleState struct {
	Position orb.Point
}

type user struct {
	Name string
}

type responseUserId struct {
	UserId int64
}

type responseUser struct {
	User user
}

type responseUsers struct {
	Users []user
}
