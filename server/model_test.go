package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/EricNeid/go-webserver/internal/verify"
	"github.com/paulmach/orb/encoding/wkt"
)

func TestVehicleStateToJson(t *testing.T) {
	// arrange
	unit := vehicleState{
		Position:  [2]float64{20, 30},
		Timestamp: time.Date(2021, 6, 15, 9, 0, 0, 0, time.UTC),
	}
	// action
	result, err := json.Marshal(unit)
	// verify
	verify.Ok(t, err)
	verify.Equals(t, "{\"position\":[20,30],\"timestamp\":\"2021-06-15T09:00:00Z\"}", string(result))
}

func TestJsonToVehicleState(t *testing.T) {
	// arrange
	testdata := "{\"position\":[20,30],\"timestamp\":\"2021-06-15T09:00:00Z\"}"
	// action
	var result vehicleState
	err := json.Unmarshal([]byte(testdata), &result)
	// verify
	verify.Ok(t, err)
	verify.Assert(t, result.Position.X()-20.0 < 0.1, wkt.MarshalString(result.Position))
	verify.Assert(t, result.Position.Y()-30.0 < 0.1, wkt.MarshalString(result.Position))
	verify.Assert(t, result.Timestamp.Year() == 2021, result.Timestamp.String())
	verify.Assert(t, result.Timestamp.Month() == 6, result.Timestamp.String())
	verify.Assert(t, result.Timestamp.Day() == 15, result.Timestamp.String())
	verify.Assert(t, result.Timestamp.Hour() == 9, result.Timestamp.String())
}
