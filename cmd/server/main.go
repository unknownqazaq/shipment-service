package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	gen "github.com/unknownqazaq/shipment-service/gen"
	"github.com/unknownqazaq/shipment-service/internal/application"
	"github.com/unknownqazaq/shipment-service/internal/infrastructure/repository"
	"github.com/unknownqazaq/shipment-service/internal/transport"
)

func main() {
	// Создаём репозиторий
	repo := repository.NewInMemoryShipmentRepository()

	// Создаём сервис (use cases)
	service := application.NewShipmentService(repo)

	// Создаём gRPC handler
	handler := transport.NewShipmentHandler(service)

	// Создаём gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервис
	gen.RegisterShipmentServiceServer(grpcServer, handler)

	// Открываем порт
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("🚀 Shipment gRPC server running on :50051")

	// 7. Запускаем
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
