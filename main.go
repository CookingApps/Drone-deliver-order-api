package main

import (
	"drone-delivery-api/routes"
	"drone-delivery-api/store"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create our in-memory store
	s := store.New()

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "Drone Delivery API",
			"version": "1.0.0",
		})
	})

	// Setup all routes
	routes.Setup(r, s)

	log.Println("🚁 Drone Delivery API running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
