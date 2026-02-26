GO = /opt/homebrew/bin/go
DOCKER = /usr/local/bin/docker
SERVICES = auth user product cart order payment tax shipping return search review media notification cms promotion loyalty affiliate chat ai
REGISTRY ?= ghcr.io/southern-martin/ecommerce
TAG ?= latest
PROTO_DIR := proto
GEN_DIR := proto/gen

# ─── Database defaults ──────────────────────────────────────────
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 15432
POSTGRES_USER ?= ecommerce
POSTGRES_PASSWORD ?= ecommerce_secret

.PHONY: help build test lint vet fmt run-infra stop-infra clean \
        docker-build docker-push deploy-dev deploy-staging deploy-prod \
        proto run-all stop migrate-up test-coverage integration-test \
        e2e-test load-test-smoke load-test verify-all

# ─── Help ───────────────────────────────────────────────────────
help: ## Show this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Proto Generation ──────────────────────────────────────────
proto: ## Generate protobuf code (placeholder)
	@echo "Generating protobuf code..."
	@rm -rf $(GEN_DIR)
	@mkdir -p $(GEN_DIR)
	@for dir in $(PROTO_DIR)/*/; do \
		name=$$(basename $$dir); \
		mkdir -p $(GEN_DIR)/$$name; \
		protoc \
			--proto_path=$(PROTO_DIR) \
			--go_out=$(GEN_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_out=$(GEN_DIR) \
			--go-grpc_opt=paths=source_relative \
			$$dir*.proto 2>/dev/null || true; \
	done
	@echo "Proto generation complete."

# ─── Build ──────────────────────────────────────────────────────
build: ## Build all Go services
	@mkdir -p bin
	@for svc in $(SERVICES); do \
		echo "Building $$svc..."; \
		cd services/$$svc && $(GO) build -o ../../bin/$$svc ./cmd/ && cd ../..; \
	done
	@echo "All services built successfully."

# ─── Test ───────────────────────────────────────────────────────
test: ## Run all tests with race detection
	@for svc in $(SERVICES); do \
		echo "Testing $$svc..."; \
		cd services/$$svc && $(GO) test ./... -race -v && cd ../..; \
	done
	@echo "Testing pkg..."
	@cd pkg && $(GO) test ./... -race -v
	@echo "All tests passed."

test-coverage: ## Run tests with coverage report
	@for svc in $(SERVICES); do \
		echo "Testing $$svc with coverage..."; \
		cd services/$$svc && $(GO) test ./... -race -coverprofile=coverage.out -covermode=atomic && cd ../..; \
	done
	@echo "Testing pkg with coverage..."
	@cd pkg && $(GO) test ./... -race -coverprofile=coverage.out -covermode=atomic
	@echo "Coverage reports generated."

# ─── Lint ───────────────────────────────────────────────────────
lint: ## Run golangci-lint on all services
	@for svc in $(SERVICES); do \
		echo "Linting $$svc..."; \
		cd services/$$svc && golangci-lint run ./... && cd ../..; \
	done
	@echo "Linting pkg..."
	@cd pkg && golangci-lint run ./...
	@echo "All linting passed."

# ─── Vet ────────────────────────────────────────────────────────
vet: ## Run go vet on all services
	@for svc in $(SERVICES); do \
		echo "Vetting $$svc..."; \
		cd services/$$svc && $(GO) vet ./... && cd ../..; \
	done
	@echo "Vetting pkg..."
	@cd pkg && $(GO) vet ./...
	@echo "All vet checks passed."

# ─── Format ─────────────────────────────────────────────────────
fmt: ## Format all Go code
	@for svc in $(SERVICES); do \
		echo "Formatting $$svc..."; \
		cd services/$$svc && gofmt -s -w . && cd ../..; \
	done
	@echo "Formatting pkg..."
	@cd pkg && gofmt -s -w .
	@echo "All code formatted."

# ─── Infrastructure ────────────────────────────────────────────
run-infra: ## Start infrastructure (postgres, redis, nats, elasticsearch, minio)
	$(DOCKER) compose up -d postgres redis nats elasticsearch minio
	@echo "Infrastructure services started."

stop-infra: ## Stop infrastructure containers
	$(DOCKER) compose down
	@echo "Infrastructure services stopped."

run-all: ## Start all services via docker-compose
	$(DOCKER) compose up -d
	@echo "All services started."

stop: ## Stop all services
	$(DOCKER) compose down
	@echo "All services stopped."

# ─── Migrations ─────────────────────────────────────────────────
migrate-up: ## Run database migrations for all services
	@for svc in $(SERVICES); do \
		echo "Migrating $$svc..."; \
		for f in services/$$svc/migrations/*.up.sql; do \
			[ -f "$$f" ] && PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d ecommerce_$$svc -f $$f || true; \
		done; \
	done
	@echo "All migrations applied."

# ─── Docker ─────────────────────────────────────────────────────
docker-build: ## Build Docker images for all services
	@for svc in $(SERVICES); do \
		echo "Building Docker image for $$svc..."; \
		$(DOCKER) build -t $(REGISTRY)/$$svc:$(TAG) -f services/$$svc/Dockerfile .; \
	done
	@echo "All Docker images built."

docker-push: ## Push Docker images to registry
	@for svc in $(SERVICES); do \
		echo "Pushing Docker image for $$svc..."; \
		$(DOCKER) push $(REGISTRY)/$$svc:$(TAG); \
	done
	@echo "All Docker images pushed."

# ─── Deploy ─────────────────────────────────────────────────────
deploy-dev: ## Deploy to dev environment (Kubernetes)
	kubectl apply -k deploy/k8s/overlays/dev
	@echo "Deployed to dev environment."

deploy-staging: ## Deploy to staging environment (Kubernetes)
	kubectl apply -k deploy/k8s/overlays/staging
	@echo "Deployed to staging environment."

deploy-prod: ## Deploy to production environment (Kubernetes)
	@echo "WARNING: You are about to deploy to PRODUCTION."
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	kubectl apply -k deploy/k8s/overlays/prod
	@echo "Deployed to production environment."

# ─── Testing ────────────────────────────────────────────────────
integration-test: ## Run integration tests
	cd tests/integration && $(GO) test ./... -v -tags=integration

e2e-test: ## Run E2E tests
	cd tests/e2e && $(GO) test ./... -v -tags=e2e

load-test-smoke: ## Run k6 smoke test
	k6 run tests/load/smoke.js

load-test: ## Run k6 load test
	k6 run tests/load/load_test.js

# ─── Verify ─────────────────────────────────────────────────────
verify-all: ## Build + vet + lint + test all services
	@echo "=== Formatting all code ===" && $(MAKE) fmt
	@echo "=== Building all services ===" && $(MAKE) build
	@echo "=== Running go vet ===" && $(MAKE) vet
	@echo "=== Running linter ===" && $(MAKE) lint
	@echo "=== Running tests ===" && $(MAKE) test
	@echo "=== All checks passed ==="

# ─── Clean ──────────────────────────────────────────────────────
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf $(GEN_DIR)
	@for svc in $(SERVICES); do \
		rm -f services/$$svc/coverage.out; \
		cd services/$$svc && $(GO) clean && cd ../..; \
	done
	rm -f pkg/coverage.out
	@echo "Cleaned build artifacts."
