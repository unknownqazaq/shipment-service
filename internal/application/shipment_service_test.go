package application

import (
	"context"
	"errors"
	"testing"

	"github.com/unknownqazaq/shipment-service/internal/domain"
)

// Mock репозиторий
type mockRepo struct {
	shipments map[string]*domain.Shipment
	events    map[string][]*domain.StatusEvent
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		shipments: make(map[string]*domain.Shipment),
		events:    make(map[string][]*domain.StatusEvent),
	}
}

func (m *mockRepo) Save(ctx context.Context, s *domain.Shipment) error {
	m.shipments[s.ID] = s
	return nil
}

func (m *mockRepo) FindByID(ctx context.Context, id string) (*domain.Shipment, error) {
	s, ok := m.shipments[id]
	if !ok {
		return nil, errors.New("shipment not found")
	}
	return s, nil
}

func (m *mockRepo) SaveEvent(ctx context.Context, e *domain.StatusEvent) error {
	m.events[e.ShipmentID] = append(m.events[e.ShipmentID], e)
	return nil
}

func (m *mockRepo) FindEventsByShipmentID(ctx context.Context, shipmentID string) ([]*domain.StatusEvent, error) {
	return m.events[shipmentID], nil
}

// Тест 1 — создание shipment
func TestShipmentService_CreateShipment(t *testing.T) {
	svc := NewShipmentService(newMockRepo())

	shipment, err := svc.CreateShipment(
		context.Background(),
		"REF-001", "Almaty", "Astana",
		"John", "TRUCK-1",
		1000, 500,
	)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if shipment.ID == "" {
		t.Error("expected ID to be set")
	}
	if shipment.Status != domain.StatusPending {
		t.Errorf("expected pending, got: %v", shipment.Status)
	}
}

// Тест 2 — получение shipment
func TestShipmentService_GetShipment(t *testing.T) {
	svc := NewShipmentService(newMockRepo())

	created, _ := svc.CreateShipment(
		context.Background(),
		"REF-001", "Almaty", "Astana",
		"John", "TRUCK-1",
		1000, 500,
	)

	found, err := svc.GetShipment(context.Background(), created.ID)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if found.ID != created.ID {
		t.Errorf("expected ID %v, got %v", created.ID, found.ID)
	}
}

// Тест 3 — shipment не найден
func TestShipmentService_GetShipment_NotFound(t *testing.T) {
	svc := NewShipmentService(newMockRepo())

	_, err := svc.GetShipment(context.Background(), "non-existent-id")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// Тест 4 — добавление статусного события
func TestShipmentService_AddStatusEvent(t *testing.T) {
	svc := NewShipmentService(newMockRepo())

	created, _ := svc.CreateShipment(
		context.Background(),
		"REF-001", "Almaty", "Astana",
		"John", "TRUCK-1",
		1000, 500,
	)

	event, err := svc.AddStatusEvent(
		context.Background(),
		created.ID,
		domain.StatusPickedUp,
		"Driver picked up",
	)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if event.Status != domain.StatusPickedUp {
		t.Errorf("expected picked_up, got: %v", event.Status)
	}
}

// Тест 5 — история событий
func TestShipmentService_GetShipmentHistory(t *testing.T) {
	svc := NewShipmentService(newMockRepo())

	created, _ := svc.CreateShipment(
		context.Background(),
		"REF-001", "Almaty", "Astana",
		"John", "TRUCK-1",
		1000, 500,
	)

	svc.AddStatusEvent(context.Background(), created.ID, domain.StatusPickedUp, "picked up")
	svc.AddStatusEvent(context.Background(), created.ID, domain.StatusInTransit, "in transit")

	events, err := svc.GetShipmentHistory(context.Background(), created.ID)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got: %v", len(events))
	}
}
