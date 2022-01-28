package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (srv ApplicationServer) addUser(c *gin.Context) {
	var user user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := addUser(srv.logger, srv.db, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := struct {
		UserId int64 `json:"userId"`
	}{
		UserId: id,
	}
	c.JSON(http.StatusCreated, res)
}

func (srv ApplicationServer) deleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = deleteUser(srv.logger, srv.db, id)
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
	retrievedUser, err := getUser(srv.logger, srv.db, id)
	if err == ErrorNotFound {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := struct {
		User user `json:"user"`
	}{
		User: retrievedUser,
	}
	c.JSON(http.StatusOK, res)
}

func (srv ApplicationServer) getUsers(c *gin.Context) {
	users, err := getUsers(srv.logger, srv.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := struct {
		Users []user `json:"users"`
	}{
		Users: users,
	}
	c.JSON(http.StatusOK, res)
}
