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
	"github.com/EricNeid/go-webserver/model"
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
	// arrange
	integrationtest.Setup()
	defer integrationtest.Cleanup()
	db, _ := integrationtest.GetDbConnectionPool()
	gin.SetMode(gin.TestMode)
	unit := NewApplicationServer(log.New(os.Stdout, "test: ", log.LstdFlags), db, ":5001")
	recoder := httptest.NewRecorder()

	// action
	t.Run("Adding user", func(t *testing.T) {
		// arrange
		testdata := model.User{Name: "testuser"}
		req := httptest.NewRequest("POST", "/user", strings.NewReader(testdata.ToJson()))
		// action
		unit.Router.ServeHTTP(recoder, req)
		// verify
		verify.Equals(t, 200, recoder.Code)
		userId = recoder.Result().Body

	})
}
