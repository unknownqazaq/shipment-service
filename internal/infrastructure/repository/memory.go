package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/unknownqazaq/shipment-service/internal/domain"
)

type InMemoryShipmentRepository struct {
	mu        sync.RWMutex
	shipments map[string]*domain.Shipment
	events    map[string][]*domain.StatusEvent
}

func NewInMemoryShipmentRepository() *InMemoryShipmentRepository {
	return &InMemoryShipmentRepository{
		shipments: make(map[string]*domain.Shipment),
		events:    make(map[string][]*domain.StatusEvent),
	}
}

func (r *InMemoryShipmentRepository) Save(ctx context.Context, shipment *domain.Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.shipments[shipment.ID] = shipment
	return nil
}

func (r *InMemoryShipmentRepository) FindByID(ctx context.Context, id string) (*domain.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shipment, ok := r.shipments[id]
	if !ok {
		return nil, errors.New("shipment not found")
	}
	return shipment, nil
}

func (r *InMemoryShipmentRepository) SaveEvent(ctx context.Context, event *domain.StatusEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events[event.ShipmentID] = append(r.events[event.ShipmentID], event)
	return nil
}

func (r *InMemoryShipmentRepository) FindEventsByShipmentID(ctx context.Context, shipmentID string) ([]*domain.StatusEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events, ok := r.events[shipmentID]
	if !ok {
		return []*domain.StatusEvent{}, nil
	}
	return events, nil
}
