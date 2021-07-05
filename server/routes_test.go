package server

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/EricNeid/go-webserver/internal/verify"
)

func TestWelcome(t *testing.T) {
	// arrange
	request := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	// action
	welcome(responseRecorder, request)
	// verify
	verify.Assert(t, responseRecorder.Code == 200, fmt.Sprintf("Status code is %d\n", responseRecorder.Code))
	verify.Assert(t, responseRecorder.Body.String() == "Hello, World!", fmt.Sprintf("Body is %s\n", responseRecorder.Body.String()))
}
