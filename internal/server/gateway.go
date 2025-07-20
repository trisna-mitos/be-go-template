package server

import (
	"context"
	"database/sql"
	"log"

	"go-backend-service/pkg/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartHTTPGateway(ctx context.Context, mux *runtime.ServeMux, db *sql.DB, grpcAddr string) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register all handlers
	if err := pb.RegisterProductServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		log.Fatalf("Failed to register product gateway: %v", err)
	}

	if err := pb.RegisterDipanTypeServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		log.Fatalf("Failed to register dipan gateway: %v", err)
	}

	// Tambahkan handler service lainnya di sini
}
