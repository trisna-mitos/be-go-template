package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	gorilla "github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"

	"go-backend-service/internal/server"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDatabaseURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "passDblocal")
	dbname := getEnv("DB_NAME", "grpc_product")
	
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)
}

func main() {
	// Get configuration from environment
	grpcPort := getEnv("GRPC_PORT", "50051")
	httpPort := getEnv("HTTP_PORT", "8080")
	
	// Add colon prefix if not present
	if !strings.HasPrefix(grpcPort, ":") {
		grpcPort = ":" + grpcPort
	}
	if !strings.HasPrefix(httpPort, ":") {
		httpPort = ":" + httpPort
	}
	dbConnStr := getDatabaseURL()
	
	log.Printf("🚀 Starting Go gRPC Backend Server with config:")
	log.Printf("  📡 HTTP Port: %s", httpPort)
	log.Printf("  🔧 gRPC Port: %s", grpcPort)
	log.Printf("  💾 Database: %s", getDatabaseURL())
	
	// Initialize DB
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established successfully")

	go server.StartGRPCServer(db, grpcPort)

	/*
		// Initialize dependencies
		productRepo := repository.NewPostgresProductRepository(db)
		productUsecase := usecase.NewProductUsecase(productRepo)
		productHandler := grpc.NewProductHandler(productUsecase)

		// Start gRPC server
		go func() {
			lis, err := net.Listen("tcp", grpcPort)
			if err != nil {
				log.Fatalf("Failed to listen: %v", err)
			}

			grpcServer := grpcpkg.NewServer()
			pb.RegisterProductServiceServer(grpcServer, productHandler)

			log.Printf("Starting gRPC server on port %s", grpcPort)
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
		}()
	*/
	// Start HTTP server (gRPC-Gateway)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	// In Docker, services communicate using service names
	grpcHost := getEnv("GRPC_HOST", "localhost")
	grpcEndpoint := grpcHost + grpcPort
	server.StartHTTPGateway(ctx, mux, db, grpcEndpoint)
	/*
		opts := []grpcpkg.DialOption{grpcpkg.WithTransportCredentials(insecure.NewCredentials())}

		if err := pb.RegisterProductServiceHandlerFromEndpoint(
			ctx,
			mux,
			"localhost"+grpcPort,
			opts,
		); err != nil {
			log.Fatalf("Failed to register gateway: %v", err)
		}
	*/
	// Router untuk Swagger UI
	router := muxWithSwagger(mux)

	log.Printf("Starting HTTP server on port %s", httpPort)
	if err := http.ListenAndServe(httpPort, router); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}

// muxWithSwagger menambahkan handler Swagger ke router
func muxWithSwagger(mux http.Handler) http.Handler {
	r := gorilla.NewRouter()

	// CORS Middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Apply CORS middleware to all routes
	r.Use(corsMiddleware)

	// Serve gRPC-Gateway API
	r.PathPrefix("/v1/").Handler(mux)

	// Serve Swagger JSON
	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		filePath := "docs/backend_api.swagger.json"
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("Swagger file not found at: %s", filePath)
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading swagger file: %v", err)
			http.Error(w, "Error reading swagger file", http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}).Methods("GET", "OPTIONS")

	// Serve Swagger UI (static)
	fs := http.FileServer(http.Dir("static/swagger-ui/dist"))
	r.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", fs))

	// Redirect root to Swagger UI
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-ui/", http.StatusMovedPermanently)
	})

	return r
}
