#  Drone Delivery Order API

A production-ready RESTful API built with **Go (Golang)** and **Gin** for managing drone-based delivery operations — register drones, create orders, assign deliveries, and track status in real time.

> Built as a portfolio project to demonstrate backend engineering skills relevant to drone logistics systems.

---

## 🛠 Tech Stack

| Tool                                    | Purpose                           |
| --------------------------------------- | --------------------------------- |
| [Go](https://golang.org/)               | Core language                     |
| [Gin](https://github.com/gin-gonic/gin) | HTTP web framework                |
| [UUID](https://github.com/google/uuid)  | Unique ID generation              |
| `sync.RWMutex`                          | Concurrent-safe in-memory storage |

---

##  Features

- ✅ Register and manage a fleet of drones
- ✅ Create delivery orders with pickup/dropoff locations
- ✅ Assign orders to available drones (with business rule validation)
- ✅ Track and update order/drone status through the delivery lifecycle
- ✅ Battery-level guard — drones below 20% cannot be assigned
- ✅ Auto-release drone back to `available` when delivery is completed or failed
- ✅ Thread-safe in-memory data store using `sync.RWMutex`
- ✅ Clean separation of concerns: models, store, handlers, routes

---

## 🗂 Project Structure

```
drone-delivery-api/
├── main.go             # Entry point
├── models/
│   └── models.go       # Drone and Order data types
├── store/
│   └── store.go        # Thread-safe in-memory data store
├── handlers/
│   ├── drone.go        # Drone HTTP handlers
│   └── order.go        # Order HTTP handlers
└── routes/
    └── routes.go       # Route definitions
```

---

##  Getting Started

### Prerequisites

- Go 1.21+ installed → [Download Go](https://golang.org/dl/)

### Installation

```bash
# Clone the repo
git clone https://github.com/CookingApps/drone-delivery-api.git
cd drone-delivery-api

# Install dependencies
go mod tidy

# Run the server
go run main.go
```

Server starts at: `http://localhost:8080`

---

##  API Endpoints

### Health Check

| Method | Endpoint  | Description                    |
| ------ | --------- | ------------------------------ |
| GET    | `/health` | Check if the server is running |

---

###  Drones

| Method | Endpoint                    | Description                    |
| ------ | --------------------------- | ------------------------------ |
| POST   | `/api/v1/drones`            | Register a new drone           |
| GET    | `/api/v1/drones`            | Get all drones                 |
| GET    | `/api/v1/drones/:id`        | Get a specific drone           |
| PATCH  | `/api/v1/drones/:id/status` | Update drone status or battery |

---

###  Orders

| Method | Endpoint                    | Description                  |
| ------ | --------------------------- | ---------------------------- |
| POST   | `/api/v1/orders`            | Create a new delivery order  |
| GET    | `/api/v1/orders`            | Get all orders               |
| GET    | `/api/v1/orders/:id`        | Get a specific order         |
| POST   | `/api/v1/orders/:id/assign` | Assign an order to a drone   |
| PATCH  | `/api/v1/orders/:id/status` | Update order delivery status |

---

##  Request & Response Examples

### Register a Drone

**POST** `/api/v1/drones`

```json
// Request
{
  "name": "Eagle-1",
  "model": "DJI Matrice 300",
  "battery_level": 95
}

// Response 201
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "name": "Eagle-1",
  "model": "DJI Matrice 300",
  "battery_level": 95,
  "status": "available",
  "created_at": "2025-04-18T10:00:00Z"
}
```

---

### Create a Delivery Order

**POST** `/api/v1/orders`

```json
// Request
{
  "package_details": "Medicine 500g",
  "pickup_location": "Wuse II, Abuja",
  "dropoff_location": "Garki, Abuja",
  "recipient_name": "Chidi Okafor",
  "recipient_phone": "08012345678"
}

// Response 201
{
  "id": "a1b2c3d4-...",
  "package_details": "Medicine 500g",
  "pickup_location": "Wuse II, Abuja",
  "dropoff_location": "Garki, Abuja",
  "recipient_name": "Chidi Okafor",
  "recipient_phone": "08012345678",
  "status": "pending",
  "created_at": "2025-04-18T10:05:00Z",
  "updated_at": "2025-04-18T10:05:00Z"
}
```

---

### Assign Order to Drone

**POST** `/api/v1/orders/:id/assign`

```json
// Request
{
  "drone_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
}

// Response 200
{
  "message": "order successfully assigned to drone",
  "order": { "status": "assigned", "assigned_drone_id": "f47ac10b-..." },
  "drone": { "status": "busy" }
}
```

> Assignment will fail if:
>
> - Drone is not `available`
> - Drone battery is below 20%
> - Order is not in `pending` status

---

### Update Order Status

**PATCH** `/api/v1/orders/:id/status`

```json
// Request
{
  "status": "in_flight"
}
```

**Valid order statuses (in order):**
`pending` → `assigned` → `in_flight` → `delivered` / `failed`

> When status is set to `delivered` or `failed`, the assigned drone is automatically released back to `available`.

---

### Update Drone Status

**PATCH** `/api/v1/drones/:id/status`

```json
{
  "status": "maintenance",
  "battery_level": 30
}
```

**Valid drone statuses:** `available` | `busy` | `maintenance`

---

##  Delivery Lifecycle

```
[Create Order]     →   status: pending
      ↓
[Assign to Drone]  →   status: assigned   |   drone: busy
      ↓
[Drone Takes Off]  →   status: in_flight
      ↓
[Delivery Done]    →   status: delivered  |   drone: available  ✅
      ↓ (if fails)
[Delivery Failed]  →   status: failed     |   drone: available  ❌
```

---

##  Testing with cURL

```bash
# Health check
curl http://localhost:8080/health

# Register a drone
curl -X POST http://localhost:8080/api/v1/drones \
  -H "Content-Type: application/json" \
  -d '{"name":"Eagle-1","model":"DJI Matrice 300","battery_level":95}'

# Create an order
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"package_details":"Medicine 500g","pickup_location":"Wuse II, Abuja","dropoff_location":"Garki, Abuja","recipient_name":"Chidi Okafor","recipient_phone":"08012345678"}'

# Assign order to drone (replace IDs)
curl -X POST http://localhost:8080/api/v1/orders/<ORDER_ID>/assign \
  -H "Content-Type: application/json" \
  -d '{"drone_id":"<DRONE_ID>"}'

# Update to in_flight
curl -X PATCH http://localhost:8080/api/v1/orders/<ORDER_ID>/status \
  -H "Content-Type: application/json" \
  -d '{"status":"in_flight"}'

# Mark as delivered
curl -X PATCH http://localhost:8080/api/v1/orders/<ORDER_ID>/status \
  -H "Content-Type: application/json" \
  -d '{"status":"delivered"}'
```

---

##  Future Improvements

- [ ] PostgreSQL / Redis persistence
- [ ] JWT authentication & role-based access
- [ ] WebSocket support for live drone telemetry
- [ ] GPS coordinate tracking per order
- [ ] Swagger / OpenAPI documentation
- [ ] Docker + docker-compose setup
- [ ] Unit & integration tests

---

##  Author

**Ayobami Masterpiece**
Software Engineer | Drone Tech Enthusiast

---


