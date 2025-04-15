# Navigate to the web directory for frontend commands
WEB_DIR := ./web

# Default target
all: build

# Install frontend dependencies
install-frontend:
	@echo ">>> Installing frontend dependencies in $(WEB_DIR)..."
	@cd $(WEB_DIR) && npm install

# Run frontend development server
dev-frontend:
	@echo ">>> Starting frontend development server in $(WEB_DIR)..."
	@cd $(WEB_DIR) && npm run dev

# Build frontend for production
build-frontend: install-frontend
	@echo ">>> Building frontend in $(WEB_DIR)..."
	@cd $(WEB_DIR) && npm run build

# Build Go backend
build-backend:
	@echo ">>> Building Go backend..."
	@go build -o ./bin/umka-server ./cmd/umka-server

# Build both frontend and backend
build: build-frontend build-backend

# Clean build artifacts
clean:
	@echo ">>> Cleaning build artifacts..."
	@rm -rf $(WEB_DIR)/dist
	@rm -rf $(WEB_DIR)/node_modules
	@rm -f ./bin/umka-server

# Run the backend server (assuming it serves the frontend)
run: build
	@echo ">>> Running server..."
	@./bin/umka-server

.PHONY: all install-frontend dev-frontend build-frontend build-backend build clean run 