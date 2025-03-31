# Variables
APP_NAME = assessmentApi
BUILD_DIR = ./bin
MAIN_PATH = ./cmd/assessmentApi/main.go
GO = go

# Default build architecture
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Define standard colors
BLUE = \033[0;34m
NC = \033[0m # No color

.PHONY: all build clean run test lint fmt help

# Default target
all: clean fmt lint test build

# Build the application
build:
	@echo "${BLUE}Building $(APP_NAME) for $(GOOS)/$(GOARCH)...${NC}"
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "${BLUE}Build complete: $(BUILD_DIR)/$(APP_NAME)${NC}"

# Clean build artifacts
clean:
	@echo "${BLUE}Cleaning build artifacts...${NC}"
	@rm -rf $(BUILD_DIR)
	@echo "${BLUE}Clean complete${NC}"

# Run the application
run:
	@echo "${BLUE}Running $(APP_NAME)...${NC}"
	@$(GO) run $(MAIN_PATH)

# Format code
fmt:
	@echo "${BLUE}Formatting code...${NC}"
	@$(GO) fmt ./...

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "${BLUE}Building for Linux...${NC}"
	@GOOS=linux GOARCH=amd64 $(MAKE) build

build-darwin:
	@echo "${BLUE}Building for macOS...${NC}"
	@GOOS=darwin GOARCH=amd64 $(MAKE) build

build-windows:
	@echo "${BLUE}Building for Windows...${NC}"
	@GOOS=windows GOARCH=amd64 $(MAKE) build

# Show help
help:
	@echo "Available targets:"
	@echo "  all         - Clean, format, lint, test, and build the application"
	@echo "  build       - Build the application for the current platform"
	@echo "  clean       - Remove build artifacts"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  build-all   - Build for Linux, macOS, and Windows"
	@echo "  build-linux - Build for Linux"
	@echo "  build-darwin - Build for macOS"
	@echo "  build-windows - Build for Windows"
	@echo "  help        - Show this help message"