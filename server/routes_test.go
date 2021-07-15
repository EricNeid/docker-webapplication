package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/EricNeid/go-webserver/internal/integrationtest"
	"github.com/EricNeid/go-webserver/internal/verify"
	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb/encoding/wkt"
)

func TestWelcome(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	request := httptest.NewRequest("GET", "/", nil)
	recoder := httptest.NewRecorder()
	unit := NewApplicationServer(nil, ":5001")
	// action
	unit.router.ServeHTTP(recoder, request)
	// verify
	verify.Assert(t, recoder.Code == 200, fmt.Sprintf("Status code is %d\n", recoder.Code))
	verify.Assert(t, recoder.Body.String() == "Hello, World!", fmt.Sprintf("Body is %s\n", recoder.Body.String()))
}

func TestCrudUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test")
	}

	// arrange
	integrationtest.Setup()
	defer integrationtest.Cleanup()
	db, _ := integrationtest.GetDbConnectionPool()
	gin.SetMode(gin.TestMode)
	unit := NewApplicationServer(db, ":5001")
	createTableUsers(unit.logger, unit.db)

	var id int64
	t.Run("Adding user", func(t *testing.T) {
		// arrange
		testdata, _ := json.Marshal(user{Name: "testuser"})
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/users", strings.NewReader(string(testdata)))
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusCreated, res.Code)
		result := struct {
			UserId int64 `json:"userId"`
		}{}
		err := json.NewDecoder(res.Body).Decode(&result)
		verify.Ok(t, err)
		id = result.UserId
	})

	t.Run("Getting user by id", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d", id), nil)
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusOK, res.Code)
		result := struct {
			User user `json:"user"`
		}{}
		err := json.NewDecoder(res.Body).Decode(&result)
		verify.Ok(t, err)
		verify.Equals(t, "testuser", result.User.Name)
	})

	t.Run("Getting all users", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users", nil)
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusOK, res.Code)
		result := struct {
			Users []user `json:"users"`
		}{}
		err := json.NewDecoder(res.Body).Decode(&result)
		verify.Ok(t, err)
		verify.Equals(t, 1, len(result.Users))
	})

	t.Run("Deleting user by id", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/users/%d", id), nil)
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusNoContent, res.Code)
	})

	t.Run("Getting user by id should return 404", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d", id), nil)
		// action
		unit.router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, http.StatusNotFound, res.Code)
	})
}

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
		testdata, _ := json.Marshal(
			vehicleState{
				Position:  [2]float64{20, 30},
				Timestamp: time.Date(2021, 6, 15, 9, 0, 0, 0, time.UTC),
			},
		)
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/vehicleStates", strings.NewReader(string(testdata)))
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
		verify.Assert(t, result.VehicleState.Position.X()-20.0 < 0.1, wkt.MarshalString(result.VehicleState.Position))
		verify.Assert(t, result.VehicleState.Position.Y()-30.0 < 0.1, wkt.MarshalString(result.VehicleState.Position))
		verify.Assert(t, result.VehicleState.Timestamp.Year() == 2021, result.VehicleState.Timestamp.String())
		verify.Assert(t, result.VehicleState.Timestamp.Month() == 6, result.VehicleState.Timestamp.String())
		verify.Assert(t, result.VehicleState.Timestamp.Day() == 15, result.VehicleState.Timestamp.String())
		verify.Assert(t, result.VehicleState.Timestamp.Hour() == 9, result.VehicleState.Timestamp.String())
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
