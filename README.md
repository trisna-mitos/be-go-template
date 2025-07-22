# Go gRPC Backend Service

A Boilerplate Go backend service built with gRPC and HTTP REST API support, featuring clean architecture, Docker development environment, and comprehensive API documentation.

## 🚀 Quick Start

**Prerequisites:** Docker and Docker Compose

```bash
# Clone and setup
git clone <repository-url>
cd go-backend-service
cp .env.example .env

# Start development environment (one command!)
make docker-dev

# Access your APIs
# 📖 Swagger UI: http://localhost:8080/swagger-ui/
# 🌐 REST API: http://localhost:8080/v1/
# 🔧 gRPC API: localhost:50051
```

## 📋 Table of Contents

- [Project Overview](#-project-overview)
- [Architecture](#-architecture)
- [Technology Stack](#-technology-stack)
- [Complete API Development Guide](#-complete-api-development-guide)
- [Example: Building a Categories API](#-example-building-a-categories-api)
- [Development Workflow](#-development-workflow)
- [Best Practices](#-best-practices)
- [Troubleshooting](#-troubleshooting)

## 🎯 Project Overview

This is a modern Go backend service that demonstrates:

- **Dual API Support**: Both gRPC and HTTP REST endpoints from the same codebase
- **Clean Architecture**: Separation of concerns with domain, usecase, repository, and delivery layers
- **Auto-Generated Documentation**: Swagger UI generated from Protocol Buffer definitions
- **Hot Reload Development**: Automatic rebuilds and restarts during development
- **Production Ready**: Docker containerization with multi-stage builds
- **Database Migrations**: Version-controlled database schema changes

### Current APIs
- **Products**: Full CRUD operations for product management
- **Dipan Types**: Category management system

## 🏗️ Architecture

This service follows **Clean Architecture** principles:

```
📁 Project Structure
├── cmd/server/           # Application entry point
├── internal/
│   ├── domain/          # 🔵 Business entities and interfaces
│   ├── usecase/         # 🟢 Business logic layer
│   ├── repository/      # 🟡 Data access layer
│   ├── delivery/grpc/   # 🔴 API handlers (gRPC)
│   └── server/          # Server setup and middleware
├── pkg/pb/              # Protocol Buffer definitions & generated code
├── migrations/          # Database schema changes
├── docs/                # Auto-generated Swagger documentation
└── docker-compose.yml   # Development environment
```

### Architecture Layers

| Layer | Responsibility | Directory |
|-------|----------------|-----------|
| **Domain** | Business entities, interfaces | `internal/domain/` |
| **Use Case** | Business logic, orchestration | `internal/usecase/` |
| **Repository** | Data access, database operations | `internal/repository/` |
| **Delivery** | API handlers, request/response | `internal/delivery/` |
| **Infrastructure** | Database, external services | `internal/server/` |

### Data Flow
```
HTTP/gRPC Request → Delivery → Use Case → Repository → Database
                        ↓
Response ← JSON/Proto ← Domain Entity ← Query Result
```

## 🛠️ Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Language** | Go 1.24+ | High-performance backend development |
| **API** | gRPC + HTTP REST | Dual protocol support |
| **Documentation** | Protocol Buffers → Swagger | Auto-generated API docs |
| **Database** | PostgreSQL 15 | Primary data store |
| **Migrations** | golang-migrate | Database version control |
| **Development** | Docker + Air | Containerized hot-reload development |
| **Gateway** | gRPC-Gateway | gRPC to HTTP/JSON translation |

## 📚 Complete API Development Guide

Follow this step-by-step guide to create a new API from database design to working endpoints.

### Step 1: Database Migration

Create database schema changes using migrations:

```bash
# Create new migration files
make migrate-create name=create_categories_table

# This creates two files:
# migrations/000003_create_categories_table.up.sql
# migrations/000003_create_categories_table.down.sql
```

**Example Migration (`000003_create_categories_table.up.sql`):**
```sql
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id INT REFERENCES categories(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_active ON categories(is_active);
```

**Down Migration (`000003_create_categories_table.down.sql`):**
```sql
DROP TABLE IF EXISTS categories;
```

### Step 2: Protocol Buffer Definition

Define your gRPC service and messages:

**Create `pkg/pb/category.proto`:**
```protobuf
syntax = "proto3";

package pb;
option go_package = "go-backend-service/pkg/pb";

import "google/api/annotations.proto";

service CategoryService {
  rpc CreateCategory(CreateCategoryRequest) returns (Category) {
    option (google.api.http) = {
      post: "/v1/categories"
      body: "*"
    };
  }
  
  rpc GetCategory(GetCategoryRequest) returns (Category) {
    option (google.api.http) = {
      get: "/v1/categories/{id}"
    };
  }
  
  rpc ListCategories(ListCategoriesRequest) returns (ListCategoriesResponse) {
    option (google.api.http) = {
      get: "/v1/categories"
    };
  }
  
  rpc UpdateCategory(UpdateCategoryRequest) returns (Category) {
    option (google.api.http) = {
      put: "/v1/categories/{id}"
      body: "*"
    };
  }
  
  rpc DeleteCategory(DeleteCategoryRequest) returns (DeleteCategoryResponse) {
    option (google.api.http) = {
      delete: "/v1/categories/{id}"
    };
  }
}

message Category {
  int32 id = 1;
  string name = 2;
  string description = 3;
  int32 parent_id = 4;
  bool is_active = 5;
  string created_at = 6;
  string updated_at = 7;
}

message CreateCategoryRequest {
  string name = 1;
  string description = 2;
  int32 parent_id = 3;
}

message GetCategoryRequest {
  int32 id = 1;
}

message ListCategoriesRequest {
  int32 page = 1;
  int32 limit = 2;
  int32 parent_id = 3; // Optional: filter by parent
}

message ListCategoriesResponse {
  repeated Category categories = 1;
  int32 total = 2;
}

message UpdateCategoryRequest {
  int32 id = 1;
  string name = 2;
  string description = 3;
  int32 parent_id = 4;
  bool is_active = 5;
}

message DeleteCategoryRequest {
  int32 id = 1;
}

message DeleteCategoryResponse {
  bool success = 1;
  string message = 2;
}
```

**Generate Go code:**
```bash
make proto
```

This generates:
- `pkg/pb/category.pb.go` - Go structs
- `pkg/pb/category_grpc.pb.go` - gRPC server/client code  
- `pkg/pb/category.pb.gw.go` - HTTP gateway handlers
- Updates Swagger documentation in `docs/`

### Step 3: Domain Layer

Create business entities and interfaces:

**Create `internal/domain/category.go`:**
```go
package domain

import (
    "context"
    "time"
)

// Category represents the business entity
type Category struct {
    ID          int32     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    ParentID    *int32    `json:"parent_id,omitempty"` // Pointer for nullable field
    IsActive    bool      `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryRepository defines data access interface
type CategoryRepository interface {
    Create(ctx context.Context, category *Category) error
    GetByID(ctx context.Context, id int32) (*Category, error)
    List(ctx context.Context, page, limit int32, parentID *int32) ([]*Category, int32, error)
    Update(ctx context.Context, category *Category) error
    Delete(ctx context.Context, id int32) error
}

// CategoryUsecase defines business logic interface
type CategoryUsecase interface {
    CreateCategory(ctx context.Context, category *Category) error
    GetCategory(ctx context.Context, id int32) (*Category, error)
    ListCategories(ctx context.Context, page, limit int32, parentID *int32) ([]*Category, int32, error)
    UpdateCategory(ctx context.Context, category *Category) error
    DeleteCategory(ctx context.Context, id int32) error
}
```

### Step 4: Repository Layer

Implement data access logic:

**Create `internal/repository/category.go`:**
```go
package repository

import (
    "context"
    "database/sql"
    "time"

    "go-backend-service/internal/domain"
)

type postgresCategoryRepository struct {
    db *sql.DB
}

func NewPostgresCategoryRepository(db *sql.DB) domain.CategoryRepository {
    return &postgresCategoryRepository{db: db}
}

func (r *postgresCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
    query := `
        INSERT INTO categories (name, description, parent_id, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `

    now := time.Now()
    err := r.db.QueryRowContext(ctx, query,
        category.Name,
        category.Description,
        category.ParentID,
        category.IsActive,
        now,
        now,
    ).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

    return err
}

func (r *postgresCategoryRepository) GetByID(ctx context.Context, id int32) (*domain.Category, error) {
    query := `
        SELECT id, name, description, parent_id, is_active, created_at, updated_at
        FROM categories
        WHERE id = $1
    `

    category := &domain.Category{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &category.ID,
        &category.Name,
        &category.Description,
        &category.ParentID,
        &category.IsActive,
        &category.CreatedAt,
        &category.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    return category, nil
}

func (r *postgresCategoryRepository) List(ctx context.Context, page, limit int32, parentID *int32) ([]*domain.Category, int32, error) {
    // Count total records
    countQuery := `SELECT COUNT(*) FROM categories WHERE ($1::int IS NULL OR parent_id = $1)`
    var total int32
    err := r.db.QueryRowContext(ctx, countQuery, parentID).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    // Get paginated results
    offset := (page - 1) * limit
    query := `
        SELECT id, name, description, parent_id, is_active, created_at, updated_at
        FROM categories
        WHERE ($1::int IS NULL OR parent_id = $1)
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

    rows, err := r.db.QueryContext(ctx, query, parentID, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var categories []*domain.Category
    for rows.Next() {
        category := &domain.Category{}
        err := rows.Scan(
            &category.ID,
            &category.Name,
            &category.Description,
            &category.ParentID,
            &category.IsActive,
            &category.CreatedAt,
            &category.UpdatedAt,
        )
        if err != nil {
            return nil, 0, err
        }
        categories = append(categories, category)
    }

    return categories, total, nil
}

func (r *postgresCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
    query := `
        UPDATE categories
        SET name = $1, description = $2, parent_id = $3, is_active = $4, updated_at = $5
        WHERE id = $6
        RETURNING updated_at
    `

    now := time.Now()
    err := r.db.QueryRowContext(ctx, query,
        category.Name,
        category.Description,
        category.ParentID,
        category.IsActive,
        now,
        category.ID,
    ).Scan(&category.UpdatedAt)

    return err
}

func (r *postgresCategoryRepository) Delete(ctx context.Context, id int32) error {
    query := `DELETE FROM categories WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}
```

### Step 5: Use Case Layer

Implement business logic:

**Create `internal/usecase/category.go`:**
```go
package usecase

import (
    "context"
    "fmt"

    "go-backend-service/internal/domain"
)

type categoryUsecase struct {
    categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepo domain.CategoryRepository) domain.CategoryUsecase {
    return &categoryUsecase{
        categoryRepo: categoryRepo,
    }
}

func (u *categoryUsecase) CreateCategory(ctx context.Context, category *domain.Category) error {
    // Business logic validation
    if category.Name == "" {
        return fmt.Errorf("category name is required")
    }

    // Validate parent category exists if parent_id is provided
    if category.ParentID != nil {
        _, err := u.categoryRepo.GetByID(ctx, *category.ParentID)
        if err != nil {
            return fmt.Errorf("parent category not found")
        }
    }

    // Set default values
    category.IsActive = true

    return u.categoryRepo.Create(ctx, category)
}

func (u *categoryUsecase) GetCategory(ctx context.Context, id int32) (*domain.Category, error) {
    return u.categoryRepo.GetByID(ctx, id)
}

func (u *categoryUsecase) ListCategories(ctx context.Context, page, limit int32, parentID *int32) ([]*domain.Category, int32, error) {
    // Set default pagination
    if page <= 0 {
        page = 1
    }
    if limit <= 0 {
        limit = 10
    }
    if limit > 100 {
        limit = 100 // Prevent excessive data retrieval
    }

    return u.categoryRepo.List(ctx, page, limit, parentID)
}

func (u *categoryUsecase) UpdateCategory(ctx context.Context, category *domain.Category) error {
    // Validate category exists
    existing, err := u.categoryRepo.GetByID(ctx, category.ID)
    if err != nil {
        return fmt.Errorf("category not found")
    }

    // Business logic validation
    if category.Name == "" {
        return fmt.Errorf("category name is required")
    }

    // Prevent circular references
    if category.ParentID != nil && *category.ParentID == category.ID {
        return fmt.Errorf("category cannot be its own parent")
    }

    // Preserve created_at
    category.CreatedAt = existing.CreatedAt

    return u.categoryRepo.Update(ctx, category)
}

func (u *categoryUsecase) DeleteCategory(ctx context.Context, id int32) error {
    // Check if category exists
    _, err := u.categoryRepo.GetByID(ctx, id)
    if err != nil {
        return fmt.Errorf("category not found")
    }

    // TODO: Add business logic to check if category has children or is used by products
    
    return u.categoryRepo.Delete(ctx, id)
}
```

### Step 6: Delivery Layer (gRPC Handlers)

Create API handlers:

**Create `internal/delivery/grpc/category.go`:**
```go
package grpc

import (
    "context"

    "go-backend-service/internal/domain"
    "go-backend-service/pkg/pb"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type CategoryHandler struct {
    pb.UnimplementedCategoryServiceServer
    categoryUsecase domain.CategoryUsecase
}

func NewCategoryHandler(categoryUsecase domain.CategoryUsecase) *CategoryHandler {
    return &CategoryHandler{
        categoryUsecase: categoryUsecase,
    }
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
    // Convert from proto to domain
    category := &domain.Category{
        Name:        req.Name,
        Description: req.Description,
        IsActive:    true,
    }

    if req.ParentId > 0 {
        category.ParentID = &req.ParentId
    }

    // Execute use case
    err := h.categoryUsecase.CreateCategory(ctx, category)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
    }

    // Convert from domain to proto
    return h.domainToProto(category), nil
}

func (h *CategoryHandler) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
    category, err := h.categoryUsecase.GetCategory(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "category not found: %v", err)
    }

    return h.domainToProto(category), nil
}

func (h *CategoryHandler) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
    page := req.Page
    limit := req.Limit
    var parentID *int32
    
    if req.ParentId > 0 {
        parentID = &req.ParentId
    }

    categories, total, err := h.categoryUsecase.ListCategories(ctx, page, limit, parentID)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to list categories: %v", err)
    }

    // Convert domain slice to proto slice
    protoCategories := make([]*pb.Category, len(categories))
    for i, cat := range categories {
        protoCategories[i] = h.domainToProto(cat)
    }

    return &pb.ListCategoriesResponse{
        Categories: protoCategories,
        Total:      total,
    }, nil
}

func (h *CategoryHandler) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error) {
    category := &domain.Category{
        ID:          req.Id,
        Name:        req.Name,
        Description: req.Description,
        IsActive:    req.IsActive,
    }

    if req.ParentId > 0 {
        category.ParentID = &req.ParentId
    }

    err := h.categoryUsecase.UpdateCategory(ctx, category)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to update category: %v", err)
    }

    return h.domainToProto(category), nil
}

func (h *CategoryHandler) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.DeleteCategoryResponse, error) {
    err := h.categoryUsecase.DeleteCategory(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to delete category: %v", err)
    }

    return &pb.DeleteCategoryResponse{
        Success: true,
        Message: "Category deleted successfully",
    }, nil
}

// Helper function to convert domain to proto
func (h *CategoryHandler) domainToProto(cat *domain.Category) *pb.Category {
    proto := &pb.Category{
        Id:          cat.ID,
        Name:        cat.Name,
        Description: cat.Description,
        IsActive:    cat.IsActive,
        CreatedAt:   cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt:   cat.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }

    if cat.ParentID != nil {
        proto.ParentId = *cat.ParentID
    }

    return proto
}
```

### Step 7: Server Registration

Register your new service with the server:

**Update `internal/server/grpc.go`:**
```go
// Add import
import (
    // ... existing imports
    "go-backend-service/internal/delivery/grpc"
    "go-backend-service/internal/repository"
    "go-backend-service/internal/usecase"
    "go-backend-service/pkg/pb"
)

func StartGRPCServer(db *sql.DB, port string) {
    // ... existing code ...

    // Register existing services
    productRepo := repository.NewPostgresProductRepository(db)
    productUsecase := usecase.NewProductUsecase(productRepo)
    productHandler := grpc.NewProductHandler(productUsecase)
    pb.RegisterProductServiceServer(grpcServer, productHandler)

    // Register new Category service
    categoryRepo := repository.NewPostgresCategoryRepository(db)
    categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
    categoryHandler := grpc.NewCategoryHandler(categoryUsecase)
    pb.RegisterCategoryServiceServer(grpcServer, categoryHandler)

    // ... rest of the code ...
}
```

**Update `internal/server/gateway.go` to register HTTP gateway:**
```go
func StartHTTPGateway(ctx context.Context, mux *runtime.ServeMux, db *sql.DB, grpcEndpoint string) {
    // ... existing registrations ...

    // Register Category service HTTP gateway
    if err := pb.RegisterCategoryServiceHandlerFromEndpoint(
        ctx,
        mux,
        grpcEndpoint,
        opts,
    ); err != nil {
        log.Fatalf("Failed to register Category gateway: %v", err)
    }

    // ... rest of the code ...
}
```

### Step 8: Testing Your New API

1. **Apply migrations:**
```bash
# Migrations run automatically when containers restart
make docker-up
```

2. **Regenerate documentation:**
```bash
make proto
```

3. **Test via Swagger UI:**
   - Visit http://localhost:8080/swagger-ui/
   - Look for "CategoryService" section
   - Try the endpoints interactively

4. **Test via curl:**
```bash
# Create a category
curl -X POST http://localhost:8080/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Electronics", "description": "Electronic devices and gadgets"}'

# List categories
curl http://localhost:8080/v1/categories

# Get specific category
curl http://localhost:8080/v1/categories/1

# Update category
curl -X PUT http://localhost:8080/v1/categories/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Electronics & Gadgets", "description": "Updated description", "is_active": true}'

# Delete category
curl -X DELETE http://localhost:8080/v1/categories/1
```

## 🔄 Development Workflow

### Daily Development Process

1. **Start Environment:**
```bash
cd go-backend-service
make docker-dev
```

2. **Monitor Logs:**
```bash
# In a separate terminal
export PATH=$PATH:~/.local/bin && docker-compose logs -f app
```

3. **Make Changes:** 
   - Edit Go files (automatic reload)
   - Update .proto files → run `make proto`
   - Add migrations → run `make migrate-create name=feature_name`

4. **Test Changes:**
   - Swagger UI: http://localhost:8080/swagger-ui/
   - Direct API calls: `curl http://localhost:8080/v1/...`

5. **Database Operations:**
```bash
# Access database
export PATH=$PATH:~/.local/bin && docker-compose exec postgres psql -U postgres -d grpc_product

# Reset database (if needed)
make docker-db-reset
```

### Project Commands Reference

| Command | Description |
|---------|-------------|
| `make docker-dev` | Start development environment |
| `make docker-logs` | View application logs |
| `make docker-down` | Stop all services |
| `make proto` | Generate code from .proto files |
| `make migrate-create name=X` | Create new migration |
| `docker ps` | Check container status |

## 📖 Best Practices

### 1. Database Design
- Use appropriate data types (SERIAL for auto-increment IDs)
- Add proper indexes for query performance
- Include created_at/updated_at timestamps
- Use foreign key constraints for referential integrity
- Write both UP and DOWN migrations

### 2. API Design
- Use consistent naming conventions (snake_case for proto fields)
- Include proper HTTP status codes in error responses
- Implement pagination for list endpoints
- Validate input data at the use case layer
- Provide meaningful error messages

### 3. Code Organization
- Follow Clean Architecture layers strictly
- Keep business logic in use cases, not handlers
- Use interfaces for dependency injection
- Handle context cancellation properly
- Implement proper error handling

### 4. Protocol Buffers
- Use semantic field numbers (don't reuse)
- Add HTTP annotations for REST API generation
- Include proper validation rules
- Document your services and messages
- Version your APIs when making breaking changes

### 5. Error Handling
```go
// Good: Specific error types
return status.Errorf(codes.NotFound, "category with id %d not found", id)

// Good: Wrap internal errors
return status.Errorf(codes.Internal, "failed to create category: %v", err)

// Bad: Generic errors
return status.Errorf(codes.Unknown, "something went wrong")
```

### 6. Testing Strategy
- Test each layer independently
- Use dependency injection for mocking
- Test both happy path and error cases
- Include integration tests for APIs
- Use Docker for consistent test environments

## 🔧 Troubleshooting

### Common Issues

**1. Migration Errors**
```bash
# Check migration status
make migrate-status

# Force migration version
docker-compose run --rm migrate force <version>

# Reset database completely
make docker-db-reset
```

**2. Proto Generation Issues**
```bash
# Clean and regenerate
rm -rf third_party/googleapis
make proto
```

**3. Hot Reload Not Working**
```bash
# Check Air logs
docker-compose logs app

# Rebuild container
make docker-build
make docker-up
```

**4. Database Connection Issues**
```bash
# Check database logs
make docker-logs-db

# Verify database is healthy
docker ps  # Should show healthy status
```

**5. gRPC Service Not Found**
- Ensure service is registered in `internal/server/grpc.go`
- Check that proto generated files are up to date
- Verify import paths are correct

### Debug Commands
```bash
# Container shell access
docker-compose exec app sh

# Database shell access
docker-compose exec postgres psql -U postgres -d grpc_product

# Check running processes
docker-compose ps

# View all logs
docker-compose logs
```

## 🚀 Production Deployment

### Building Production Image
```bash
# Build production image
docker build --target production -t go-backend-service:prod .

# Run production container
docker run -p 8080:8080 -p 50051:50051 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-db-password \
  go-backend-service:prod
```

### Environment Variables
```bash
# Required for production
DB_HOST=your-database-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=your-database-name

# Optional
HTTP_PORT=:8080
GRPC_PORT=:50051
GO_ENV=production
```

## 📚 Additional Resources

- **Docker Development Guide**: See `README-DOCKER.md`
- **Development Setup**: See `CLAUDE.md` 
- **Protocol Buffers**: https://protobuf.dev/
- **gRPC-Gateway**: https://grpc-ecosystem.github.io/grpc-gateway/
- **Clean Architecture**: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

## 🤝 Contributing

1. Follow the established patterns in this README
2. 

---

**🎉 You now have everything needed to build robust APIs with this Go gRPC backend service!**

For questions or issues, check the troubleshooting section or review the existing implementations in the codebase.