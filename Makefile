.PHONY: build install clean run test help

# Binary name
BINARY_NAME=termius-from-walmart
INSTALL_PATH=/usr/local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
	@echo "Build complete! Binary: ./$(BINARY_NAME)"

install: build ## Build and install to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo mv $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Installation complete! Run with: $(BINARY_NAME)"

clean: ## Remove build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@echo "Clean complete!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies ready!"

run: build ## Build and run the application
	@./$(BINARY_NAME)

uninstall: ## Uninstall from /usr/local/bin
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstall complete!"

test: ## Run tests
	$(GOTEST) -v ./...

all: clean build ## Clean and build
