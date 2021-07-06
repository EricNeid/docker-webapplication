package server

import (
	"net/http"
	"strconv"

	"github.com/EricNeid/go-webserver/database"
	"github.com/EricNeid/go-webserver/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func welcome(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}

func (srv ApplicationServer) addUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := database.AddUser(srv.Logger, srv.Db, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := model.ResponseUserId{UserId: id}
	c.JSON(http.StatusOK, res)
}

func (srv ApplicationServer) deleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = database.DeleteUser(srv.Logger, srv.Db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (srv ApplicationServer) getUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := database.GetUser(srv.Logger, srv.Db, id)
	if err == pgx.ErrNoRows {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseUser{User: user})
}

func (srv ApplicationServer) getUsers(c *gin.Context) {
	users, err := database.GetUsers(srv.Logger, srv.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ResponseUsers{Users: users})
}
