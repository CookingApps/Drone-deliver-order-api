package store

import (
	"drone-delivery-api/models"
	"errors"
	"sync"
)

// Store holds all our data in memory
type Store struct {
	mu     sync.RWMutex
	drones map[string]*models.Drone
	orders map[string]*models.Order
}

// New creates a fresh store
func New() *Store {
	return &Store{
		drones: make(map[string]*models.Drone),
		orders: make(map[string]*models.Order),
	}
}

// --- Drone Methods ---

func (s *Store) AddDrone(d *models.Drone) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.drones[d.ID] = d
}

func (s *Store) GetDrone(id string) (*models.Drone, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.drones[id]
	if !ok {
		return nil, errors.New("drone not found")
	}
	return d, nil
}

func (s *Store) GetAllDrones() []*models.Drone {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.Drone, 0, len(s.drones))
	for _, d := range s.drones {
		list = append(list, d)
	}
	return list
}

func (s *Store) UpdateDrone(d *models.Drone) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.drones[d.ID] = d
}

// --- Order Methods ---

func (s *Store) AddOrder(o *models.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[o.ID] = o
}

func (s *Store) GetOrder(id string) (*models.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return o, nil
}

func (s *Store) GetAllOrders() []*models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.Order, 0, len(s.orders))
	for _, o := range s.orders {
		list = append(list, o)
	}
	return list
}

func (s *Store) UpdateOrder(o *models.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[o.ID] = o
}
