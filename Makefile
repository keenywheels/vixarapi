# paths
ENV_PATH := .env.dev
DOCKER_COMPOSE_PATH = ./docker-compose.yaml

# exec
GO := go
DOCKER_COMPOSE := docker compose -f $(DOCKER_COMPOSE_PATH) --env-file $(ENV_PATH)

# ======================================================================
# APP MANAGEMENT
.PHONY: build
build:
	@$(GO) build -o bin/app cmd/app/main.go
	@echo "[MAKE] done building app"

.PHONY: generate-api
generate-api:
	@ogen --target internal/api/v1/ --clean api/openapi.yaml
	@echo "[MAKE] done generating user api"

# ======================================================================
# DOCKER-COMPOSE
.PHONY: docker-build
docker-build:
	@$(DOCKER_COMPOSE) build
	@echo "[MAKE] build docker images"

.PHONY: docker-up
docker-up:
	@$(DOCKER_COMPOSE) up -d
	@echo "[MAKE] done starting containers"

.PHONY: docker-stop
docker-stop:
	@$(DOCKER_COMPOSE) stop
	@echo "[MAKE] done stopping containers"

.PHONY: docker-down
docker-down:
	@$(DOCKER_COMPOSE) down -v
	@echo "[MAKE] done removing containers"