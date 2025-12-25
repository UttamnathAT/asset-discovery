.PHONY: proto proto-install proto-deps proto-clean run build help stop

# Variables
PROTO_DIR := api/protos
GEN_DIR := api/gen
PROTOC := protoc

# Colors
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

help: ## Show help
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

proto-install: ## Install proto tools
	@echo "$(GREEN)Installing proto tools...$(NC)"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

proto-deps: ## Download third_party/googleapis
	@echo "$(GREEN)Downloading googleapis...$(NC)"
	@if [ ! -d "third_party/googleapis" ]; then \
		git clone --depth 1 https://github.com/googleapis/googleapis.git third_party/googleapis; \
		echo "$(GREEN)✓ Downloaded$(NC)"; \
	else \
		echo "$(YELLOW)✓ Already exists$(NC)"; \
	fi

proto: proto-deps ## Generate proto files
	@echo "$(GREEN)Generating proto files...$(NC)"
	@mkdir -p $(GEN_DIR)/health
	@echo "$(YELLOW)Generating health.proto...$(NC)"
	@$(PROTOC) \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		-I=$(PROTO_DIR) \
		-I=third_party/googleapis \
		$(PROTO_DIR)/health/health.proto
	@echo "$(YELLOW)Generating asset_discovery.proto...$(NC)"
	@$(PROTOC) \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GEN_DIR) \
		--grpc-gateway_opt=paths=source_relative \
		-I=$(PROTO_DIR) \
		-I=third_party/googleapis \
		$(PROTO_DIR)/asset_discovery.proto
	@echo "$(GREEN)✓ Proto files generated in $(GEN_DIR)$(NC)"

proto-clean: ## Clean generated files
	@echo "$(YELLOW)Cleaning generated files...$(NC)"
	@rm -rf $(GEN_DIR)
	@echo "$(GREEN)✓ Cleaned$(NC)"

proto-all: proto-install proto ## Install tools and generate proto files

setup: proto-all ## Complete setup (install tools + generate protos)

stop: ## Stop running servers
	@echo "$(YELLOW)Stopping servers...$(NC)"
	@-PID=$$(lsof -ti:50051 2>/dev/null); [ -n "$$PID" ] && kill -9 $$PID 2>/dev/null || true
	@-PID=$$(lsof -ti:8003 2>/dev/null); [ -n "$$PID" ] && kill -9 $$PID 2>/dev/null || true
	@-pkill -f "health-server" 2>/dev/null || true
	@-pkill -f "cmd/server/main.go" 2>/dev/null || true
	@sleep 1
	@echo "$(GREEN)✓ Servers stopped$(NC)"

run: ## Run the server (stops existing servers first)
	@$(MAKE) stop || true
	@echo "$(GREEN)Running server...$(NC)"
	@go run cmd/server/main.go

build: ## Build the server
	@echo "$(GREEN)Building server...$(NC)"
	@go build -o bin/server cmd/server/main.go
	@echo "$(GREEN)✓ Build complete: bin/server$(NC)"

tidy: ## Run go mod tidy
	@go mod tidy
	@echo "$(GREEN)✓ Done$(NC)"

clean: proto-clean ## Clean all generated files

