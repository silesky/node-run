# Variables
BINARY_NAME = nrun
BUILD_DIR = bin

# Default target
all: build


# Variables
BINARY_NAME = nrun
BUILD_DIR = bin

# Default target
all: build

# Build for all architectures
build: build-linux-amd64 build-linux-arm64 build-mac-amd64 build-mac-arm64

build-linux-amd64:
		@echo "🚀 Building for Linux amd64...\n"
		GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/node-task-runner $(ARGS)

build-linux-arm64:
		@echo "🚀 Building for Linux arm64...\n"
		GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/node-task-runner $(ARGS)

build-mac-amd64:
		@echo "🍏 Building for macOS amd64...\n"
		GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/node-task-runner $(ARGS)

build-mac-arm64:
		@echo "🍏 Building for macOS arm64 (Apple Silicon)...\n"
		GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/node-task-runner $(ARGS)

# Run the binary
run: 
		go run ./cmd/node-task-runner $(ARGS)

# Clean the build directory
clean:
		rm -rf $(BUILD_DIR)

# Run all tests - make test ARGS="-run TestSpecificFunction"
test:
		go test ./... -v $(ARGS) 

# List the make targets
help:                                                                                                                    
		@grep '^[^#[:space:]].*:' Makefile
		
# Bump the version and push to github
version:
		bash ./bump-version.sh

