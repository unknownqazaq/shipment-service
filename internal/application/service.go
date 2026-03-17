package application

import (
	"context"

	"github.com/unknownqazaq/shipment-service/internal/domain"
)

type Service interface {
	CreateShipment(ctx context.Context, refNumber, origin, destination, driverName, unitNumber string, amount, driverRevenue float64) (*domain.Shipment, error)
	GetShipment(ctx context.Context, id string) (*domain.Shipment, error)
	AddStatusEvent(ctx context.Context, shipmentID string, newStatus domain.Status, comment string) (*domain.StatusEvent, error)
	GetShipmentHistory(ctx context.Context, shipmentID string) ([]*domain.StatusEvent, error)
}
