package transport

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/unknownqazaq/shipment-service/gen"
	"github.com/unknownqazaq/shipment-service/internal/application"
	"github.com/unknownqazaq/shipment-service/internal/domain"
)

type ShipmentHandler struct {
	gen.UnimplementedShipmentServiceServer
	service application.Service
}

func NewShipmentHandler(service application.Service) *ShipmentHandler {
	return &ShipmentHandler{service: service}
}

func (h *ShipmentHandler) CreateShipment(ctx context.Context, req *gen.CreateShipmentRequest) (*gen.Shipment, error) {
	shipment, err := h.service.CreateShipment(
		ctx,
		req.ReferenceNumber,
		req.Origin,
		req.Destination,
		req.DriverName,
		req.UnitNumber,
		req.Amount,
		req.DriverRevenue,
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return domainToProto(shipment), nil
}

func (h *ShipmentHandler) GetShipment(ctx context.Context, req *gen.GetShipmentRequest) (*gen.Shipment, error) {
	shipment, err := h.service.GetShipment(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return domainToProto(shipment), nil
}

func (h *ShipmentHandler) AddStatusEvent(ctx context.Context, req *gen.AddStatusEventRequest) (*gen.StatusEvent, error) {
	newStatus := protoStatusToDomain(req.Status)

	event, err := h.service.AddStatusEvent(ctx, req.ShipmentId, newStatus, req.Comment)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &gen.StatusEvent{
		Id:         event.ID,
		ShipmentId: event.ShipmentID,
		Status:     req.Status,
		Comment:    event.Comment,
		CreatedAt:  timestamppb.New(event.CreatedAt),
	}, nil
}

func (h *ShipmentHandler) GetShipmentHistory(ctx context.Context, req *gen.GetShipmentHistoryRequest) (*gen.GetShipmentHistoryResponse, error) {
	events, err := h.service.GetShipmentHistory(ctx, req.ShipmentId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoEvents []*gen.StatusEvent
	for _, e := range events {
		protoEvents = append(protoEvents, &gen.StatusEvent{
			Id:         e.ID,
			ShipmentId: e.ShipmentID,
			Status:     domainStatusToProto(e.Status),
			Comment:    e.Comment,
			CreatedAt:  timestamppb.New(e.CreatedAt),
		})
	}

	return &gen.GetShipmentHistoryResponse{Events: protoEvents}, nil
}

//Конвертеры domain на proto

func domainToProto(s *domain.Shipment) *gen.Shipment {
	return &gen.Shipment{
		Id:              s.ID,
		ReferenceNumber: s.ReferenceNumber,
		Origin:          s.Origin,
		Destination:     s.Destination,
		Status:          domainStatusToProto(s.Status),
		DriverName:      s.DriverName,
		UnitNumber:      s.UnitNumber,
		Amount:          s.Amount,
		DriverRevenue:   s.DriverRevenue,
		CreatedAt:       timestamppb.New(s.CreatedAt),
	}
}

func domainStatusToProto(s domain.Status) gen.ShipmentStatus {
	switch s {
	case domain.StatusPending:
		return gen.ShipmentStatus_PENDING
	case domain.StatusPickedUp:
		return gen.ShipmentStatus_PICKED_UP
	case domain.StatusInTransit:
		return gen.ShipmentStatus_IN_TRANSIT
	case domain.StatusDelivered:
		return gen.ShipmentStatus_DELIVERED
	case domain.StatusCancelled:
		return gen.ShipmentStatus_CANCELLED
	default:
		return gen.ShipmentStatus_PENDING
	}
}

func protoStatusToDomain(s gen.ShipmentStatus) domain.Status {
	switch s {
	case gen.ShipmentStatus_PICKED_UP:
		return domain.StatusPickedUp
	case gen.ShipmentStatus_IN_TRANSIT:
		return domain.StatusInTransit
	case gen.ShipmentStatus_DELIVERED:
		return domain.StatusDelivered
	case gen.ShipmentStatus_CANCELLED:
		return domain.StatusCancelled
	default:
		return domain.StatusPending
	}
}
