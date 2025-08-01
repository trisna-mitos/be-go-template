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

### Docker Development (Recommended)
**For new developers, use Docker for the easiest setup:**

```bash
# Quick start with Docker
cp .env.example .env
make docker-dev

# Alternative Docker commands
make docker-up              # Start all services with rebuild
make docker-up-detached     # Start in background
make docker-down            # Stop services
make docker-down-volumes    # Stop and remove data
make docker-logs            # View application logs
make docker-logs-db         # View PostgreSQL logs
make docker-db-reset        # Reset database with fresh data

# Real-time log monitoring (recommended for development)
export PATH=$PATH:~/.local/bin && docker-compose logs -f app
```

**Docker Services Access:**
- HTTP API: http://localhost:8080/v1/
- Swagger UI: http://localhost:8080/swagger-ui/
- gRPC API: localhost:50051
- PostgreSQL: localhost:5432

**Important Notes:**
- If `docker-compose` command is not found, install it locally: `curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o ~/.local/bin/docker-compose && chmod +x ~/.local/bin/docker-compose && export PATH=$PATH:~/.local/bin`
- VCS build errors are already fixed with `-buildvcs=false` flag in `.air.toml`
- Port configuration is handled automatically by the Go application (supports both `:8080` and `8080` formats)

### Native Development (Local Setup)

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

**Docker (Recommended):**
```bash
make docker-up
```

**Native (requires local PostgreSQL):**
```bash
go run cmd/server/main.go
```

The server runs:
- gRPC server on `:50051`
- HTTP gateway on `:8080`
- Swagger UI at `http://localhost:8080/swagger-ui/`

## Database Configuration

**Docker Environment (Automatic):**
The application automatically connects to the PostgreSQL container using environment variables from `.env` file.

**Native Environment:**
- Host: localhost:5432  
- Database: grpc_product
- User: postgres
- Password: passDblocal

**Environment Variables:**
The application supports these environment variables:
- `DB_HOST` (default: postgres for Docker, localhost for native)
- `DB_PORT` (default: 5432)
- `DB_USER` (default: postgres)
- `DB_PASSWORD` (default: passDblocal)
- `DB_NAME` (default: grpc_product)

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
- `docker-compose.yml` defines the Docker development environment
- `Dockerfile` contains multi-stage build for development and production
- `.air.toml` configures hot reload for development
- `.env.example` template for environment variables

## Protocol Buffer Workflow

When modifying APIs:
1. Edit `.proto` files in `pkg/pb/`
2. Run `make proto` to regenerate Go code and documentation
3. Update repository/usecase/delivery layers as needed
4. Create database migrations if schema changes are required
5. Test via gRPC clients or HTTP endpoints

## Getting Started for New Developers

1. **Clone and setup environment:**
   ```bash
   git clone <repository-url>
   cd go-backend-service
   cp .env.example .env
   ```

2. **Start with Docker (recommended):**
   ```bash
   make docker-dev
   ```

3. **Access the application:**
   - Visit http://localhost:8080/swagger-ui/ for API documentation
   - API endpoints available at http://localhost:8080/v1/

4. **For detailed Docker setup:** See `README-DOCKER.md`
5. **For complete development guide:** See `README.md` for step-by-step API creation tutorial

## Daily Development Workflow

### **Starting Work (After Laptop Restart)**
```bash
# Navigate to project directory
cd /home/trisna/personal-project/prima-utama/go-backend-service

# Check if containers are running
docker ps

# Start development environment (if containers stopped)
export PATH=$PATH:~/.local/bin && make docker-dev
# OR start in background: make docker-up-detached

# Monitor logs during development
export PATH=$PATH:~/.local/bin && docker-compose logs -f app
```

### **Adding New Features**
1. **Plan the feature** - Define API endpoints, database changes
2. **Database changes** (if needed):
   ```bash
   make migrate-create name=add_new_feature
   # Edit migration files in migrations/ directory
   ```
3. **Update Protocol Buffers** (if needed):
   ```bash
   # Edit .proto files in pkg/pb/
   make proto  # Regenerates Go code and Swagger docs
   ```
4. **Implement feature** following clean architecture:
   - `internal/domain/` - Business entities and interfaces
   - `internal/repository/` - Data access implementations
   - `internal/usecase/` - Business logic
   - `internal/delivery/grpc/` - gRPC handlers
   - Hot reload will restart app automatically on file saves
5. **Test via APIs**:
   - Swagger UI: http://localhost:8080/swagger-ui/
   - REST API: http://localhost:8080/v1/
   - Use curl or tools like Postman for testing

### **Ending Work Session**
```bash
# Option 1: Leave running (containers will stop on laptop shutdown)
# Your data persists in Docker volumes

# Option 2: Stop gracefully
export PATH=$PATH:~/.local/bin && docker-compose down

# Option 3: Stop and remove data (if you want fresh start tomorrow)
export PATH=$PATH:~/.local/bin && docker-compose down -v
```

### **Data Persistence**
- ✅ **Database data**: Persists in Docker volume `go-backend-service_postgres_data`
- ✅ **Code changes**: Saved to your filesystem
- ✅ **Migrations**: Applied automatically on container restart
- ❌ **Container state**: Stops on laptop shutdown (but data remains)

### **Common Development Commands**
```bash
# View real-time application logs
export PATH=$PATH:~/.local/bin && docker-compose logs -f app

# Access database shell
export PATH=$PATH:~/.local/bin && docker-compose exec postgres psql -U postgres -d grpc_product

# Access application container shell
export PATH=$PATH:~/.local/bin && docker-compose exec app sh

# Reset database (if needed)
make docker-db-reset

# Rebuild containers (after major changes)
export PATH=$PATH:~/.local/bin && docker-compose up --build
```