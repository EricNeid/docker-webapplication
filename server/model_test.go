package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/EricNeid/go-webserver/internal/verify"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func TestVehicleStateToJson(t *testing.T) {
	// arrange
	unit := vehicleState{
		Position:  *geojson.NewGeometry(orb.Point([2]float64{20, 30})),
		Timestamp: time.Date(2021, 6, 15, 9, 0, 0, 0, time.UTC),
	}
	// action
	result, err := json.Marshal(unit)
	// verify
	verify.Ok(t, err)
	verify.Equals(t, "{\"position\":{\"type\":\"Point\",\"coordinates\":[20,30]},\"timestamp\":\"2021-06-15T09:00:00Z\"}", string(result))
}

func TestJsonToVehicleState(t *testing.T) {
	// arrange
	testdata := `
	{
		"timestamp": "2021-06-15T09:00:00Z",
		"position": {
			"type": "Point",
			"coordinates": [
				20,
				30
			]
		}
	}
	`
	// action
	var result vehicleState
	err := json.Unmarshal([]byte(testdata), &result)
	// verify
	verify.Ok(t, err)
	point := result.Position.Geometry().(orb.Point)
	verify.Condition(t, point.X()-20.0 < 0.1)
	verify.Condition(t, point.Y()-30.0 < 0.1)
	verify.Condition(t, result.Timestamp.Year() == 2021)
	verify.Condition(t, result.Timestamp.Month() == 6)
	verify.Condition(t, result.Timestamp.Day() == 15)
	verify.Condition(t, result.Timestamp.Hour() == 9)
}
