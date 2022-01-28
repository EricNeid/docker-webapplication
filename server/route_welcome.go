package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func welcome(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}
