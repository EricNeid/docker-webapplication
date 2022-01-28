package server

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/EricNeid/go-webserver/internal/verify"
	"github.com/gin-gonic/gin"
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
