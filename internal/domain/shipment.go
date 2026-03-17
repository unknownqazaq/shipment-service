package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Статусы
type Status string

const (
	StatusPending   Status = "pending"
	StatusPickedUp  Status = "picked_up"
	StatusInTransit Status = "in_transit"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
)

// Валидные переходы статусов
var validTransitions = map[Status][]Status{
	StatusPending:   {StatusPickedUp, StatusCancelled},
	StatusPickedUp:  {StatusInTransit, StatusCancelled},
	StatusInTransit: {StatusDelivered, StatusCancelled},
	StatusDelivered: {},
	StatusCancelled: {},
}

// Модель отправления
type Shipment struct {
	ID              string
	ReferenceNumber string
	Origin          string
	Destination     string
	Status          Status
	DriverName      string
	UnitNumber      string
	Amount          float64
	DriverRevenue   float64
	CreatedAt       time.Time
}

// Событие смены статуса
type StatusEvent struct {
	ID         string
	ShipmentID string
	Status     Status
	Comment    string
	CreatedAt  time.Time
}

// Создать новое отправление
func NewShipment(refNumber, origin, destination, driverName, unitNumber string, amount, driverRevenue float64) (*Shipment, error) {
	if refNumber == "" {
		return nil, errors.New("reference number is required")
	}
	if origin == "" {
		return nil, errors.New("origin is required")
	}
	if destination == "" {
		return nil, errors.New("destination is required")
	}

	return &Shipment{
		ID:              uuid.New().String(),
		ReferenceNumber: refNumber,
		Origin:          origin,
		Destination:     destination,
		Status:          StatusPending,
		DriverName:      driverName,
		UnitNumber:      unitNumber,
		Amount:          amount,
		DriverRevenue:   driverRevenue,
		CreatedAt:       time.Now(),
	}, nil
}

// Проверить можно ли перейти в новый статус
func (s *Shipment) CanTransition(newStatus Status) bool {
	allowed, ok := validTransitions[s.Status]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == newStatus {
			return true
		}
	}
	return false
}

// Перейти в новый статус
func (s *Shipment) Transition(newStatus Status) (*StatusEvent, error) {
	if !s.CanTransition(newStatus) {
		return nil, errors.New("invalid status transition from " + string(s.Status) + " to " + string(newStatus))
	}

	s.Status = newStatus

	event := &StatusEvent{
		ID:         uuid.New().String(),
		ShipmentID: s.ID,
		Status:     newStatus,
		CreatedAt:  time.Now(),
	}

	return event, nil
}
