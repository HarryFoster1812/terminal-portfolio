# Simple Makefile for Terminal Portfolio

BINARY_NAME=terminal-portfolio
BUILD_DIR=build

.PHONY: build run deploy-local deploy-remote help ssh-keygen

help: ## Show available commands
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make deploy-local  - Deploy locally with Docker"
	@echo "  make deploy-remote - Deploy remotely with Docker"
	@echo "  make ssh-keygen    - Generate SSH key"

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	go run .

ssh-keygen: ## Generate SSH key
	@echo "Generating SSH key..."
	@mkdir -p .ssh
	ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ""

deploy-local: ## Deploy locally with Docker
	@echo "Building and running with Docker (local)..."
	@if [ ! -f .ssh/id_ed25519 ]; then \
		echo "Generating SSH key..."; \
		mkdir -p .ssh; \
		ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ""; \
	fi
	docker build -t $(BINARY_NAME) .
	docker run -p 2222:2222 -v $(PWD)/.ssh:/app/.ssh:ro $(BINARY_NAME)

deploy-remote: ## Deploy remotely with Docker
	@echo "Building and running with Docker (remote)..."
	docker build -t $(BINARY_NAME) .
	docker run -p 22:2222 $(BINARY_NAME)