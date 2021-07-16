package server

import (
	"time"

	"github.com/paulmach/orb/geojson"
)

type vehicleState struct {
	Position  geojson.Geometry `json:"position"`
	Timestamp time.Time        `json:"timestamp"`
}

type user struct {
	Name string `json:"name"`
}
