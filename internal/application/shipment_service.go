package application

import (
	"context"
	"errors"

	"github.com/unknownqazaq/shipment-service/internal/domain"
)

type ShipmentService struct {
	repo domain.ShipmentRepository
}

func NewShipmentService(repo domain.ShipmentRepository) *ShipmentService {
	return &ShipmentService{repo: repo}
}

// CreateShipment — создать новое отправление
func (s *ShipmentService) CreateShipment(
	ctx context.Context,
	refNumber, origin, destination, driverName, unitNumber string,
	amount, driverRevenue float64,
) (*domain.Shipment, error) {
	shipment, err := domain.NewShipment(
		refNumber, origin, destination,
		driverName, unitNumber,
		amount, driverRevenue,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, shipment); err != nil {
		return nil, err
	}

	return shipment, nil
}

// GetShipment — получить отправление по ID
func (s *ShipmentService) GetShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	if id == "" {
		return nil, errors.New("shipment id is required")
	}
	return s.repo.FindByID(ctx, id)
}

// AddStatusEvent — добавить событие смены статуса
func (s *ShipmentService) AddStatusEvent(
	ctx context.Context,
	shipmentID string,
	newStatus domain.Status,
	comment string,
) (*domain.StatusEvent, error) {
	shipment, err := s.repo.FindByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}

	event, err := shipment.Transition(newStatus)
	if err != nil {
		return nil, err
	}

	event.Comment = comment

	if err := s.repo.Save(ctx, shipment); err != nil {
		return nil, err
	}

	if err := s.repo.SaveEvent(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

// GetShipmentHistory — получить историю событий
func (s *ShipmentService) GetShipmentHistory(ctx context.Context, shipmentID string) ([]*domain.StatusEvent, error) {
	if shipmentID == "" {
		return nil, errors.New("shipment id is required")
	}
	return s.repo.FindEventsByShipmentID(ctx, shipmentID)
}
