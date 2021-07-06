package server

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/EricNeid/go-webserver/internal/integrationtest"
	"github.com/EricNeid/go-webserver/internal/verify"
	"github.com/gin-gonic/gin"
)

func TestWelcome(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	request := httptest.NewRequest("GET", "/", nil)
	recoder := httptest.NewRecorder()
	unit := NewApplicationServer(log.New(os.Stdout, "test: ", log.LstdFlags), nil, ":5001")
	// action
	unit.Router.ServeHTTP(recoder, request)
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
	unit := NewApplicationServer(log.New(os.Stdout, "test: ", log.LstdFlags), db, ":5001")
	CreateTableUsers(unit.Logger, unit.Db)

	var id int64
	t.Run("Adding user", func(t *testing.T) {
		// arrange
		testdata := User{Name: "testuser"}
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(testdata.ToJson()))
		// action
		unit.Router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, 200, res.Code)
		result, err := NewResponseUserId(res.Result().Body)
		verify.Ok(t, err)
		id = result.UserId
	})

	t.Run("Getting user by id", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/user/%d", id), nil)
		// action
		unit.Router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, 200, res.Code)
		result, err := NewResponseUser(res.Result().Body)
		verify.Ok(t, err)
		verify.Equals(t, "testuser", result.User.Name)
	})

	t.Run("Getting all users", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/user", nil)
		// action
		unit.Router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, 200, res.Code)
		result, err := NewResponseUsers(res.Result().Body)
		verify.Ok(t, err)
		verify.Equals(t, 1, len(result.Users))
	})

	t.Run("Deleting user by id", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/user/%d", id), nil)
		// action
		unit.Router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, 204, res.Code)
	})

	t.Run("Getting user by id should return 404", func(t *testing.T) {
		// arrange
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/user/%d", id), nil)
		// action
		unit.Router.ServeHTTP(res, req)
		// verify
		verify.Equals(t, 404, res.Code)
	})
}
