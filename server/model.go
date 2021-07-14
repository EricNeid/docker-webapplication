package server

import (
	"github.com/paulmach/orb"
)

type vehicleState struct {
	Position  orb.Point `json:"position"`
	Timestamp string    `json:"timestamp"`
}

type user struct {
	Name string `json:"name"`
}
