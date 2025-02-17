# Variables
BINARY_NAME = nr
BUILD_DIR = bin

# Default target
all: build

build: build-amd64 build-arm64

build-amd64:
		@echo "Building for amd64...\n"
		mkdir -p $(BUILD_DIR)/amd64
		GOARCH=amd64 go build -o $(BUILD_DIR)/amd64/$(BINARY_NAME) ./cmd/node-task-runner $(ARGS)
		chmod +x $(BUILD_DIR)/amd64/$(BINARY_NAME)

build-arm64:
		@echo "Building for arm64...\n"
		mkdir -p $(BUILD_DIR)/arm64
		GOARCH=arm64 go build -o $(BUILD_DIR)/arm64/$(BINARY_NAME) ./cmd/node-task-runner $(ARGS)
		chmod +x $(BUILD_DIR)/arm64/$(BINARY_NAME)

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
