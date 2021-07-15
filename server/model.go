package server

import (
	"time"

	"github.com/paulmach/orb"
)

type vehicleState struct {
	Position  orb.Point `json:"position"`
	Timestamp time.Time `json:"timestamp"`
}

type user struct {
	Name string `json:"name"`
}
