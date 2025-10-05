GO := go

.PHONY: build
build:
	$(GO) build -o bin/app cmd/app/main.go

.PHONY: generate-api
generate-api:
	@ogen --target internal/api/v1/ --clean api/openapi.yaml
	@echo "[MAKE] done generating user api"
