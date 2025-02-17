# Variables
BINARY_NAME = nr
BUILD_DIR = bin

# Default target
all: build


# Variables
BINARY_NAME = nr
BUILD_DIR = bin

# Default target
all: build

# Build for all architectures
build: build-linux-amd64 build-linux-arm64 build-mac-amd64 build-mac-arm64

build-linux-amd64:
		@echo "🚀 Building for Linux amd64...\n"
		mkdir -p $(BUILD_DIR)/linux/amd64
		GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/amd64/$(BINARY_NAME) ./cmd/node-task-runner $(ARGS)

build-linux-arm64:
		@echo "🚀 Building for Linux arm64...\n"
		mkdir -p $(BUILD_DIR)/linux/arm64
		GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/linux/arm64/$(BINARY_NAME) ./cmd/node-task-runner $(ARGS)

build-mac-amd64:
		@echo "🍏 Building for macOS amd64...\n"
		mkdir -p $(BUILD_DIR)/darwin/amd64
		GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/darwin/amd64/$(BINARY_NAME) ./cmd/node-task-runner $(ARGS)

build-mac-arm64:
		@echo "🍏 Building for macOS arm64 (Apple Silicon)...\n"
		mkdir -p $(BUILD_DIR)/darwin/arm64
		GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/darwin/arm64/$(BINARY_NAME) ./cmd/node-task-runner $(ARGS)

# Run the binary
run: 
		go run ./cmd/node-task-runner $(ARGS)

# Clean the build directory
clean:
		rm -rf $(BUILD_DIR)

# Run all tests - make test ARGS="-run TestSpecificFunction"
test:
		go test ./... -v$(ARGS) 

# List the make targets
help:                                                                                                                    
		@grep '^[^#[:space:]].*:' Makefile
