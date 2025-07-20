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
