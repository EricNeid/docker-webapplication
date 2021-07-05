package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/EricNeid/go-webserver/database"
	"github.com/EricNeid/go-webserver/model"
	"github.com/gin-gonic/gin"
)

func welcome(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}

func (srv ApplicationServer) addUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, "Could not create user")
		return
	}
	id, err := database.AddUser(srv.Logger, srv.Db, user)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Could not add user to datbase: %v", err))
		return
	}
	res := model.ResponseUserId{UserId: id}
	c.JSON(http.StatusOK, res)
}

func (srv ApplicationServer) deleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Could not parse id")
		return
	}
	err = database.DeleteUser(srv.Logger, srv.Db, id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not delete user")
		return
	}
	c.Status(http.StatusNoContent)
}
