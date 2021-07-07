package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func welcome(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}

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
	res := responseUserId{UserId: id}
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
	user, err := getUser(srv.logger, srv.db, id)
	if err == pgx.ErrNoRows {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, responseUser{User: user})
}

func (srv ApplicationServer) getUsers(c *gin.Context) {
	users, err := getUsers(srv.logger, srv.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, responseUsers{Users: users})
}
