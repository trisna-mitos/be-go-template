package server

import (
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"

	delivery "go-backend-service/internal/delivery/grpc"
	"go-backend-service/internal/repository"
	"go-backend-service/internal/usecase"
	"go-backend-service/pkg/pb"
)

func StartGRPCServer(db *sql.DB, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register all services
	registerServices(grpcServer, db)

	log.Printf("gRPC server listening at %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
func registerServices(server *grpc.Server, db *sql.DB) {
	// Product
	productRepo := repository.NewPostgresProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := delivery.NewProductHandler(productUsecase)
	pb.RegisterProductServiceServer(server, productHandler)

	// DipanType
	dipanRepo := repository.NewPostgresDipanTypeRepository(db)
	dipanUsecase := usecase.NewDipanTypeUsecase(dipanRepo)
	dipanHandler := delivery.NewDipanTypeHandler(dipanUsecase)
	pb.RegisterDipanTypeServiceServer(server, dipanHandler)
}
