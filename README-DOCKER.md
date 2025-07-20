# Docker Development Setup

This guide helps you set up the Go gRPC backend service for local development using Docker and Docker Compose.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) (version 20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 2.0+)
- Git

## Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-backend-service
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   Review and modify `.env` if needed for your local setup.

3. **Start the development environment**
   ```bash
   # Option 1: Using the convenience script
   ./docker-dev.sh
   
   # Option 2: Using make command
   make docker-dev
   
   # Option 3: Using docker-compose directly
   docker-compose up --build
   ```

4. **Access the services**
   - **Swagger UI**: http://localhost:8080/swagger-ui/
   - **HTTP API**: http://localhost:8080/v1/
   - **gRPC API**: localhost:50051
   - **PostgreSQL**: localhost:5432

## Docker Commands

### Basic Operations
```bash
# Start services (with rebuild)
make docker-up

# Start services in background
make docker-up-detached

# Stop services
make docker-down

# Stop services and remove volumes (full reset)
make docker-down-volumes

# Build images only
make docker-build
```

### Development & Debugging
```bash
# View application logs
make docker-logs

# View database logs
make docker-logs-db

# Open shell in app container
make docker-shell

# Open PostgreSQL shell
make docker-db-shell

# Reset database with fresh data
make docker-db-reset

# Clean up all Docker resources
make docker-clean
```

## Development Workflow

### Hot Reload
The Docker setup includes [Air](https://github.com/air-verse/air) for hot reload. When you save Go files, the application automatically rebuilds and restarts.

### Database Changes
1. **Creating Migrations**
   ```bash
   # Create new migration
   make migrate-create name=add_new_table
   ```

2. **Running Migrations**
   Migrations run automatically when containers start. To run manually:
   ```bash
   # In Docker environment, migrations run automatically
   # For manual execution, use the migrate service
   docker-compose run --rm migrate up
   ```

### Protocol Buffer Changes
When you modify `.proto` files:

1. **Generate Go code**
   ```bash
   make proto
   ```

2. **Restart services**
   ```bash
   make docker-down
   make docker-up
   ```

## Environment Configuration

### Environment Variables (.env file)
```bash
# Database Configuration
DB_HOST=postgres          # Docker service name
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=passDblocal
DB_NAME=grpc_product

# Server Configuration
HTTP_PORT=8080
GRPC_PORT=50051
GRPC_HOST=localhost       # For internal gRPC communication

# Development Environment
GO_ENV=development
```

### Port Mappings
| Service | Container Port | Host Port | Description |
|---------|---------------|-----------|-------------|
| HTTP API | 8080 | 8080 | REST API and Swagger UI |
| gRPC API | 50051 | 50051 | gRPC endpoint |
| PostgreSQL | 5432 | 5432 | Database connection |

## Docker Services

### Application (`app`)
- **Base Image**: golang:1.24.1-alpine
- **Hot Reload**: Enabled with Air
- **Volumes**: Source code mounted for development
- **Dependencies**: PostgreSQL, Migration service

### Database (`postgres`)
- **Image**: postgres:15-alpine
- **Persistent Storage**: Named volume `postgres_data`
- **Initialization**: Automatic schema and seed data loading
- **Health Checks**: Ensures database readiness

### Migration (`migrate`)
- **Image**: migrate/migrate
- **Purpose**: Runs database migrations on startup
- **Dependencies**: PostgreSQL health check

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Check what's using the port
   lsof -i :8080
   
   # Kill the process or change port in .env
   ```

2. **Database Connection Failed**
   ```bash
   # Check if PostgreSQL is running
   make docker-logs-db
   
   # Reset database
   make docker-db-reset
   ```

3. **Permission Denied on Scripts**
   ```bash
   # Make scripts executable
   chmod +x docker-dev.sh
   chmod +x scripts/wait-for-db.sh
   ```

4. **Hot Reload Not Working**
   ```bash
   # Check Air configuration
   docker-compose logs app
   
   # Rebuild container
   make docker-build
   make docker-up
   ```

### Reset Everything
```bash
# Complete reset (removes all data)
make docker-clean
make docker-up
```

### Logs and Debugging
```bash
# View all service logs
docker-compose logs

# View specific service logs
docker-compose logs app
docker-compose logs postgres
docker-compose logs migrate

# Follow logs in real-time
docker-compose logs -f app
```

## Production Deployment

The Dockerfile includes a production stage. To build for production:

```bash
# Build production image
docker build --target production -t go-backend-service:prod .

# Run production container
docker run -p 8080:8080 -p 50051:50051 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-db-password \
  go-backend-service:prod
```

## Next Steps

1. **API Testing**: Use Swagger UI at http://localhost:8080/swagger-ui/
2. **gRPC Testing**: Use tools like [grpcurl](https://github.com/fullstorydev/grpcurl) or [BloomRPC](https://github.com/bloomrpc/bloomrpc)
3. **Database Management**: Access PostgreSQL through the exposed port or container shell
4. **Development**: Modify code and see changes automatically reload

## Support

- Check Docker logs: `make docker-logs`
- Verify services: `docker-compose ps`
- Reset environment: `make docker-clean && make docker-up`