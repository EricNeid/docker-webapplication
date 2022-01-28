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
)

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
