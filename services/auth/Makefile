# Makefile

# Default environment is local
ENV ?= local
ENV_FILE := config/env/$(ENV).env

# Export environment variables from .env file
include $(ENV_FILE)
export

# Directories
CMD_DIR := cmd/main.go
MIGRATIONS_DIR := migrations
DB_URL := $(DATABASE_URL)

run:
	go run $(CMD_DIR)

test:
	go test ./... -v

test-unit:
	go test ./... -short -v

test-e2e:
	go test test/e2e/... -v

lint:
	golangci-lint run

seed:
	go run integration/scripts/db/seed_db.go --env=local

integration-test:
	go test integration/... -v


reset-db:
	go run scripts/reset_db.go

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force 1

migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

gen-mock:
	mockery --all --keeptree --output=mocks

.PHONY: run test lint seed reset-db migrate-up migrate-down migrate-force migrate-status gen-mock
