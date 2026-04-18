package handlers

import (
	"drone-delivery-api/models"
	"drone-delivery-api/store"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	store *store.Store
}

func NewOrderHandler(s *store.Store) *OrderHandler {
	return &OrderHandler{store: s}
}

// CreateOrder - POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input struct {
		PackageDetails  string `json:"package_details" binding:"required"`
		PickupLocation  string `json:"pickup_location" binding:"required"`
		DropoffLocation string `json:"dropoff_location" binding:"required"`
		RecipientName   string `json:"recipient_name" binding:"required"`
		RecipientPhone  string `json:"recipient_phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := &models.Order{
		ID:              uuid.NewString(),
		PackageDetails:  input.PackageDetails,
		PickupLocation:  input.PickupLocation,
		DropoffLocation: input.DropoffLocation,
		RecipientName:   input.RecipientName,
		RecipientPhone:  input.RecipientPhone,
		Status:          models.OrderPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	h.store.AddOrder(order)
	c.JSON(http.StatusCreated, order)
}

// GetAllOrders - GET /orders
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders := h.store.GetAllOrders()
	c.JSON(http.StatusOK, gin.H{
		"count":  len(orders),
		"orders": orders,
	})
}

// GetOrder - GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.store.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// AssignOrderToDrone - POST /orders/:id/assign
func (h *OrderHandler) AssignOrderToDrone(c *gin.Context) {
	orderID := c.Param("id")
	order, err := h.store.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if order.Status != models.OrderPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only pending orders can be assigned"})
		return
	}

	var input struct {
		DroneID string `json:"drone_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	drone, err := h.store.GetDrone(input.DroneID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "drone not found"})
		return
	}

	if drone.Status != models.DroneAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "drone is not available"})
		return
	}

	if drone.BatteryLevel < 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "drone battery too low for delivery"})
		return
	}

	// Update both drone and order
	drone.Status = models.DroneBusy
	order.Status = models.OrderAssigned
	order.AssignedDroneID = drone.ID
	order.UpdatedAt = time.Now()

	h.store.UpdateDrone(drone)
	h.store.UpdateOrder(order)

	c.JSON(http.StatusOK, gin.H{
		"message": "order successfully assigned to drone",
		"order":   order,
		"drone":   drone,
	})
}

// UpdateOrderStatus - PATCH /orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	order, err := h.store.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	var input struct {
		Status models.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If order is completed or failed, free up the drone
	if input.Status == models.OrderDelivered || input.Status == models.OrderFailed {
		if order.AssignedDroneID != "" {
			drone, err := h.store.GetDrone(order.AssignedDroneID)
			if err == nil {
				drone.Status = models.DroneAvailable
				h.store.UpdateDrone(drone)
			}
		}
	}

	order.Status = input.Status
	order.UpdatedAt = time.Now()
	h.store.UpdateOrder(order)

	c.JSON(http.StatusOK, order)
}
