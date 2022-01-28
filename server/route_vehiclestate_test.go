package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EricNeid/go-webserver/internal/integrationtest"
	"github.com/EricNeid/go-webserver/internal/verify"
	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb"
)

func TestCrudVehicleStateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test")
	}

	// arrange
	integrationtest.Setup()
	defer integrationtest.Cleanup()
	db, _ := integrationtest.GetDbConnectionPool()
	gin.SetMode(gin.TestMode)
	unit := NewApplicationServer(db, ":5001")
	createTableVehicleState(unit.logger, unit.db)

	var id int64
	t.Run("Add", func(t *testing.T) {
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
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/vehicleStates", strings.NewReader(testdata))
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusCreated, res.Code)
		result := struct {
			VehicleStateId int64 `json:"vehicleStateId"`
		}{}
		err := json.NewDecoder(res.Body).Decode(&result)
		verify.Ok(t, err)
		id = result.VehicleStateId
	})

	t.Run("Get by id", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/vehicleStates/%d", id), nil)
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusOK, res.Code)
		result := struct {
			VehicleState vehicleState `json:"vehicleState"`
		}{}
		err := json.NewDecoder(res.Body).Decode(&result)
		verify.Ok(t, err)
		point := result.VehicleState.Position.Geometry().(orb.Point)
		verify.Condition(t, point.X()-20.0 < 0.1)
		verify.Condition(t, point.Y()-30.0 < 0.1)
		verify.Condition(t, result.VehicleState.Timestamp.Year() == 2021)
		verify.Condition(t, result.VehicleState.Timestamp.Month() == 6)
		verify.Condition(t, result.VehicleState.Timestamp.Day() == 15)
		verify.Condition(t, result.VehicleState.Timestamp.Hour() == 9)
	})

	t.Run("Get all", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/vehicleStates", nil)
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusOK, res.Code)
		result := struct {
			VehicleStates []vehicleState `json:"vehicleStates"`
		}{}
		err := json.NewDecoder(res.Body).Decode(&result)
		verify.Ok(t, err)
		verify.Equals(t, 1, len(result.VehicleStates))
	})

	t.Run("Delete by id", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/vehicleStates/%d", id), nil)
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusNoContent, res.Code)
	})

	t.Run("Get by id should return 404", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/vehicleStates/%d", id), nil)
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusNotFound, res.Code)
	})
}
