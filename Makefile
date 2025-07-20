.PHONY: create-db
create-db:
	@echo "Creating database..."
	@PGPASSWORD=passDblocal psql -h localhost -U postgres -c "CREATE DATABASE grpc_product;"

.PHONY: init-db
init-db:
	@echo "Initializing database schema..."
	@PGPASSWORD=passDblocal psql -h localhost -U postgres -d grpc_product -f scripts/init_db.sql

.PHONY: drop-db
drop-db:
	@echo "Dropping database..."
	@PGPASSWORD=passDblocal psql -h localhost -U postgres -c "DROP DATABASE IF EXISTS grpc_product;"

.PHONY: reset-db
reset-db: drop-db create-db init-db
	@echo "Database reset completed."

.PHONY: proto
PROTO_FILES := $(shell find pkg/pb -name '*.proto')
proto:
	@buf export buf.build/googleapis/googleapis --output third_party/googleapis
	@protoc -I. -Ithird_party/googleapis \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=docs --openapiv2_opt=logtostderr=true,allow_merge=true,merge_file_name=backend_api \
		$(PROTO_FILES)

.PHONY: migrate-up
migrate-up:
	@migrate -database "postgres://postgres:passDblocal@localhost:5432/grpc_product?sslmode=disable" -path migrations up

.PHONY: migrate-down
migrate-down:
	@migrate -database "postgres://postgres:passDblocal@localhost:5432/grpc_product?sslmode=disable" -path migrations down

.PHONY: migrate-create
migrate-create:
	@migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-status
migrate-status:
	@migrate -database "postgres://postgres:passDblocal@localhost:5432/grpc_product?sslmode=disable" -path migrations version

.PHONY: seed-db
seed-db:
	@echo "Seeding database with test data..."
	@PGPASSWORD=passDblocal psql -h localhost -U postgres -d grpc_product -f scripts/seed_data.sql

.PHONY: setup-full
setup-full: reset-db seed-db
	@echo "Database setup completed with test data."

# ==========================================
# Docker Development Commands
# ==========================================

.PHONY: docker-dev
docker-dev:
	@echo "🐳 Starting Docker development environment..."
	@./docker-dev.sh

.PHONY: docker-up
docker-up:
	@echo "🚀 Starting Docker services..."
	@docker-compose up --build

.PHONY: docker-up-detached
docker-up-detached:
	@echo "🚀 Starting Docker services in background..."
	@docker-compose up --build -d

.PHONY: docker-down
docker-down:
	@echo "🛑 Stopping Docker services..."
	@docker-compose down

.PHONY: docker-down-volumes
docker-down-volumes:
	@echo "🗑️  Stopping Docker services and removing volumes..."
	@docker-compose down -v

.PHONY: docker-build
docker-build:
	@echo "🔨 Building Docker images..."
	@docker-compose build

.PHONY: docker-logs
docker-logs:
	@echo "📋 Showing Docker logs..."
	@docker-compose logs -f app

.PHONY: docker-logs-db
docker-logs-db:
	@echo "📋 Showing PostgreSQL logs..."
	@docker-compose logs -f postgres

.PHONY: docker-db-reset
docker-db-reset: docker-down-volumes docker-up-detached
	@echo "🔄 Database reset completed in Docker."

.PHONY: docker-shell
docker-shell:
	@echo "🐚 Opening shell in app container..."
	@docker-compose exec app sh

.PHONY: docker-db-shell
docker-db-shell:
	@echo "🐚 Opening PostgreSQL shell..."
	@docker-compose exec postgres psql -U postgres -d grpc_product

.PHONY: docker-clean
docker-clean:
	@echo "🧹 Cleaning up Docker resources..."
	@docker-compose down -v
	@docker system prune -f
