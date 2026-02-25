# Load environment variables from .env file (if it exists)
-include .env
export

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)
GOOSE_MIGRATION_DIR=migrations

# Build the application
build: generate-api
	go build -o bin/main ./main.go

# Run the application
run:
	go run main.go

# Run all pending migrations
migrate-up:
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up

# Roll back the last migration
migrate-down:
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down

# Show current migration status
migrate-status:
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" status

# Create a new migration (usage: make migrate-create name=create_users)
migrate-create:
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" create $(name) sql

# Generate API code from remote OpenAPI spec
# Usage: make generate-api url=https://example.com/api.json
generate-api:
ifdef url
	oapi-codegen --config oapi-codegen.yaml $(url)
else
	@echo "Skipping API generation: url is not set"
endif

# Run golangci-lint static analysis
lint:
	golangci-lint run ./...

.PHONY: build run migrate-up migrate-down migrate-status migrate-create generate-api lint