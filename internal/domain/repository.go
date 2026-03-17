package domain

import "context"

type ShipmentRepository interface {
	Save(ctx context.Context, shipment *Shipment) error
	FindByID(ctx context.Context, id string) (*Shipment, error)
	SaveEvent(ctx context.Context, event *StatusEvent) error
	FindEventsByShipmentID(ctx context.Context, shipmentID string) ([]*StatusEvent, error)
}
