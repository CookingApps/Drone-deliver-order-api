package routes

import (
	"drone-delivery-api/handlers"
	"drone-delivery-api/store"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, s *store.Store) {
	droneHandler := handlers.NewDroneHandler(s)
	orderHandler := handlers.NewOrderHandler(s)

	api := r.Group("/api/v1")
	{
		// Drone routes
		drones := api.Group("/drones")
		{
			drones.POST("", droneHandler.RegisterDrone)
			drones.GET("", droneHandler.GetAllDrones)
			drones.GET("/:id", droneHandler.GetDrone)
			drones.PATCH("/:id/status", droneHandler.UpdateDroneStatus)
		}

		// Order routes
		orders := api.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetAllOrders)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.POST("/:id/assign", orderHandler.AssignOrderToDrone)
			orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
		}
	}
}
