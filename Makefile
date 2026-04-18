.PHONY: build clean package migrate-up migrate-down migrate-create migrate-force

# ---------------------------------------------------------------------------
# Build targets (AWS Lambda)
# ---------------------------------------------------------------------------

# Build a statically-linked binary for AWS Lambda (Amazon Linux 2023 / Graviton).
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -tags lambda.norpc \
		-o bootstrap cmd/server/main.go

# Package the binary into a deployment zip for Lambda upload.
package: build
	zip deployment.zip bootstrap

clean:
	rm -f bootstrap deployment.zip

# ---------------------------------------------------------------------------
# Database Migration targets (golang-migrate)
# ---------------------------------------------------------------------------
# Construct DATABASE_URL from individual POSTGRES_* env vars if not set.
# Usage: source .env && make migrate-up

DATABASE_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)
MIGRATIONS_DIR = migrations

# Create a new migration pair. Usage: make migrate-create NAME=add_some_table
migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=description"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq -digits 6 $(NAME)

# Apply all pending migrations.
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose up

# Roll back the last migration.
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose down 1

# Force-set the migration version (use to recover from dirty state).
# Usage: make migrate-force VERSION=1
migrate-force:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make migrate-force VERSION=N"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(VERSION)
