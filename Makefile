# const
PROJECT_NAME := vixarapi
MIGRATIONS_PATH := ./migrations
POSTGRES_DB_URL := postgres://user:password@localhost:5432/db?sslmode=disable

# paths
ENV_PATH := ./build/.dev.env
DOCKER_COMPOSE_PATH = ./build/compose.yaml

# exec
GO := go
DOCKER_COMPOSE := docker compose -f $(DOCKER_COMPOSE_PATH) --env-file $(ENV_PATH) -p $(PROJECT_NAME)

# ======================================================================
# APP MANAGEMENT
# ======================================================================
.PHONY: build
build: clean build/bin/vixarapi build/bin/processor

.PHONY: build/bin/vixarapi
build/bin/vixarapi:
	$(GO) build -o $(@) ./cmd/vixarapi/main.go

.PHONY: build/bin/processor
build/bin/processor:
	$(GO) build -o $(@) ./cmd/processor/main.go

.PHONY: clean
clean:
	rm -rf build/bin/*

.PHONY: tidyvendor
tidyvendor:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: generate
generate:
	$(GO) generate ./...

# ======================================================================
# DOCKER-COMPOSE
# ======================================================================
.PHONY: docker-build
docker-build:
	$(DOCKER_COMPOSE) build

.PHONY: docker-up
docker-up:
	$(DOCKER_COMPOSE) up -d

.PHONY: docker-stop
docker-stop:
	$(DOCKER_COMPOSE) stop

.PHONY: docker-down
docker-down:
	$(DOCKER_COMPOSE) down -v

# ======================================================================
# MIGRATIONS
# ======================================================================
.PHONY: postgres-new
postgres-new:
	@if [ -z "$(name)" ]; then \
		echo "Error: 'name' is not set. Usage: make postgres-new name=migration_name"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) $(name)
	@echo "[INFO] create $(name) migrations in $(MIGRATIONS_PATH)"

.PHONY: postgres-up
postgres-up:
	@migrate -database "${POSTGRES_DB_URL}" -path $(MIGRATIONS_PATH) -verbose up
	@echo "[INFO] migrations up done"

.PHONY: postgres-down
postgres-down:
	@migrate -database "${POSTGRES_DB_URL}" -path $(MIGRATIONS_PATH) -verbose down
	@echo "[INFO] migrations down done"
