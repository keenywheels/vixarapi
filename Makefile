# const
PROJECT_NAME := vixarapi

# paths
ENV_PATH := .env.dev
DOCKER_COMPOSE_PATH = ./build/compose.yaml

# exec
GO := go
DOCKER_COMPOSE := docker compose -f $(DOCKER_COMPOSE_PATH) --env-file $(ENV_PATH) -p $(PROJECT_NAME)

# ======================================================================
# APP MANAGEMENT
# ======================================================================
.PHONY: build
build: clean build/bin/vixarapi

build/bin/vixarapi:
	$(GO) build -o $(@) ./cmd/vixarapi/main.go

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
