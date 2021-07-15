package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (srv ApplicationServer) addVehicleState(c *gin.Context) {
	var data vehicleState
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := addVehicleState(srv.logger, srv.db, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := struct {
		VehicleStateId int64 `json:"vehicleStateId"`
	}{
		VehicleStateId: id,
	}
	c.JSON(http.StatusCreated, res)
}
