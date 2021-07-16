package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb"
)

func (srv ApplicationServer) addVehicleState(c *gin.Context) {
	var data vehicleState
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := addVehicleState(srv.logger, srv.db, data.Position.Geometry().(orb.Point), data.Timestamp)
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

func (srv ApplicationServer) deleteVehicleState(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = deleteVehicleState(srv.logger, srv.db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (srv ApplicationServer) getVehicleState(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := getVehicleState(srv.logger, srv.db, id)
	if err == ErrorNotFound {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := struct {
		VehicleState vehicleState `json:"vehicleState"`
	}{
		VehicleState: data,
	}
	c.JSON(http.StatusOK, res)
}

func (srv ApplicationServer) getVehicleStates(c *gin.Context) {
	data, err := getVehicleStates(srv.logger, srv.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := struct {
		VehicleStates []vehicleState `json:"vehicleStates"`
	}{
		VehicleStates: data,
	}
	c.JSON(http.StatusOK, res)
}
