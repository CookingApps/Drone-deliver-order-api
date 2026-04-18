package models

import "time"

// Drone statuses
type DroneStatus string

const (
	DroneAvailable   DroneStatus = "available"
	DroneBusy        DroneStatus = "busy"
	DroneMaintenance DroneStatus = "maintenance"
)

// Order statuses
type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderAssigned  OrderStatus = "assigned"
	OrderInFlight  OrderStatus = "in_flight"
	OrderDelivered OrderStatus = "delivered"
	OrderFailed    OrderStatus = "failed"
)

// Drone represents a delivery drone
type Drone struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Model        string      `json:"model"`
	BatteryLevel int         `json:"battery_level"` // 0-100
	Status       DroneStatus `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
}

// Order represents a delivery order
type Order struct {
	ID              string      `json:"id"`
	PackageDetails  string      `json:"package_details"`
	PickupLocation  string      `json:"pickup_location"`
	DropoffLocation string      `json:"dropoff_location"`
	RecipientName   string      `json:"recipient_name"`
	RecipientPhone  string      `json:"recipient_phone"`
	Status          OrderStatus `json:"status"`
	AssignedDroneID string      `json:"assigned_drone_id,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
