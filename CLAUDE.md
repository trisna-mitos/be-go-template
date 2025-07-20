# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a Go-based gRPC backend service with HTTP gateway support for a product management system. The codebase follows clean architecture principles with clear separation of concerns:

- **Domain Layer** (`internal/domain/`): Contains business entities and interfaces (Product domain model)
- **Repository Layer** (`internal/repository/`): Data access implementations for PostgreSQL
- **Use Case Layer** (`internal/usecase/`): Business logic and orchestration
- **Delivery Layer** (`internal/delivery/`): gRPC and HTTP handlers
- **Server Layer** (`internal/server/`): Server setup for both gRPC and HTTP gateway

The service uses Protocol Buffers for API definition and gRPC-Gateway to provide REST API access alongside native gRPC.

## Key Technologies
- **gRPC** with Protocol Buffers for API definitions
- **gRPC-Gateway** for HTTP/REST API exposure
- **PostgreSQL** as the primary database
- **golang-migrate** for database migrations
- **Swagger UI** for API documentation (served at `/swagger-ui/`)

## Development Commands

### Database Operations
```bash
# Create database
make create-db

# Initialize database schema
make init-db

# Run migrations up
make migrate-up

# Run migrations down
make migrate-down

# Create new migration (requires name parameter)
make migrate-create name=<migration_name>

# Check migration status
make migrate-status

# Seed database with test data
make seed-db

# Full database setup (reset + seed)
make setup-full

# Reset database (drop, create, init)
make reset-db
```

### Protocol Buffer Generation
```bash
# Generate Go code from .proto files
make proto
```

This command:
- Exports Google APIs to `third_party/googleapis`
- Generates Go structs, gRPC server/client code, and gRPC-Gateway handlers
- Creates OpenAPI/Swagger documentation in `docs/`

### Running the Service
```bash
# Start the server (from cmd/server/)
go run cmd/server/main.go
```

The server runs:
- gRPC server on `:50051`
- HTTP gateway on `:8080`
- Swagger UI at `http://localhost:8080/swagger-ui/`

## Database Configuration

Default PostgreSQL connection:
- Host: localhost:5432
- Database: grpc_product
- User: postgres
- Password: passDblocal

## API Structure

The service manages two main entities:
1. **Products** - Core product management with CRUD operations
2. **Dipan Types** - Supporting entity type system

Both entities have:
- Protocol buffer definitions in `pkg/pb/`
- Generated Go code for gRPC and HTTP
- Swagger documentation
- Database migrations

## File Structure Notes

- `pkg/pb/` contains Protocol Buffer definitions and generated code
- `docs/` contains generated Swagger/OpenAPI specifications
- `migrations/` contains database migration files
- `scripts/` contains database initialization and seed scripts
- `static/swagger-ui/` contains Swagger UI assets for API documentation
- `third_party/googleapis/` contains Google API dependencies (generated)

## Protocol Buffer Workflow

When modifying APIs:
1. Edit `.proto` files in `pkg/pb/`
2. Run `make proto` to regenerate Go code and documentation
3. Update repository/usecase/delivery layers as needed
4. Create database migrations if schema changes are required
5. Test via gRPC clients or HTTP endpoints