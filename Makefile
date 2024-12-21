# Define the application name
APP_NAME = "gonas"

# Detect the OS and architecture of the host system
ifeq ($(OS),Windows_NT)
    GOOS = windows
    GOARCH = amd64
    EXT = .exe
else
    UNAME_S := $(shell uname -s)
    UNAME_M := $(shell uname -m)

    ifeq ($(UNAME_S),Linux)
        GOOS = linux
    else ifeq ($(UNAME_S),Darwin)
        GOOS = darwin
    endif

    ifeq ($(UNAME_M),x86_64)
        GOARCH = amd64
    else ifeq ($(UNAME_M),arm64)
        GOARCH = arm64
    else ifeq ($(UNAME_M),arm)
        GOARCH = arm
    else
        GOARCH = 386
    endif
    EXT =
endif

# Define the build output directory
BUILD_DIR = bin

# Default target
build-all: build-linux build-windows build-darwin

# Test target
test:
	@echo "Running tests..."
	go test ./...

# Build target
build: test
	@echo "Building $(APP_NAME) for OS: $(GOOS), Architecture: $(GOARCH)"
	go build -ldflags="-s -w" -trimpath -o $(BUILD_DIR)/$(APP_NAME)$(EXT)

# Clean target
clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)

# Run target
run:
	@echo "Running $(APP_NAME)..."
	./$(BUILD_DIR)/$(APP_NAME)$(EXT)

# Cross-compile targets for other OS and architectures
build-linux: test
	set GOOS=linux 
	set GOARCH=amd64 
	go build -ldflags="-s -w" -trimpath -o $(BUILD_DIR)/$(APP_NAME)-linux 
	upx --best --lzma $(BUILD_DIR)/$(APP_NAME)-linux

build-darwin: test
	set GOOS=darwin 
	set GOARCH=amd64 
	go build -ldflags="-s -w" -trimpath -o $(BUILD_DIR)/$(APP_NAME)-darwin 
	upx --best --lzma $(BUILD_DIR)/$(APP_NAME)-darwin

build-windows: test
	set GOOS=windows 
	set GOARCH=amd64 
	go build -ldflags="-s -w" -trimpath -o $(BUILD_DIR)/$(APP_NAME)-windows.exe 
	upx --best --lzma $(BUILD_DIR)/$(APP_NAME)-windows.exe

# Help target to list all commands
help:
	@echo "Makefile commands:"
	@echo "  make build         - Build for the local OS and architecture"
	@echo "  make run           - Build and run the application"
	@echo "  make clean         - Remove the build directory"
	@echo "  make test          - Run all tests"
	@echo "  make build-linux   - Build for Linux"
	@echo "  make build-darwin  - Build for macOS"
	@echo "  make build-windows - Build for Windows"
	@echo "  make help          - Show this help message"