package handlers

import (
	"drone-delivery-api/models"
	"drone-delivery-api/store"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DroneHandler struct {
	store *store.Store
}

func NewDroneHandler(s *store.Store) *DroneHandler {
	return &DroneHandler{store: s}
}

// RegisterDrone - POST /drones
func (h *DroneHandler) RegisterDrone(c *gin.Context) {
	var input struct {
		Name         string `json:"name" binding:"required"`
		Model        string `json:"model" binding:"required"`
		BatteryLevel int    `json:"battery_level" binding:"required,min=0,max=100"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	drone := &models.Drone{
		ID:           uuid.NewString(),
		Name:         input.Name,
		Model:        input.Model,
		BatteryLevel: input.BatteryLevel,
		Status:       models.DroneAvailable,
		CreatedAt:    time.Now(),
	}

	h.store.AddDrone(drone)
	c.JSON(http.StatusCreated, drone)
}

// GetAllDrones - GET /drones
func (h *DroneHandler) GetAllDrones(c *gin.Context) {
	drones := h.store.GetAllDrones()
	c.JSON(http.StatusOK, gin.H{
		"count":  len(drones),
		"drones": drones,
	})
}

// GetDrone - GET /drones/:id
func (h *DroneHandler) GetDrone(c *gin.Context) {
	id := c.Param("id")
	drone, err := h.store.GetDrone(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, drone)
}

// UpdateDroneStatus - PATCH /drones/:id/status
func (h *DroneHandler) UpdateDroneStatus(c *gin.Context) {
	id := c.Param("id")
	drone, err := h.store.GetDrone(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Status       models.DroneStatus `json:"status" binding:"required"`
		BatteryLevel *int               `json:"battery_level"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	drone.Status = input.Status
	if input.BatteryLevel != nil {
		drone.BatteryLevel = *input.BatteryLevel
	}

	h.store.UpdateDrone(drone)
	c.JSON(http.StatusOK, drone)
}
