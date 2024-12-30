# Variables
BINARY_NAME = ntk
BUILD_DIR = bin

# Default target
all: build

# Build the binary
build:
		@echo "Building...\n"
		mkdir -p $(BUILD_DIR)
		go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/node-task-runner $(ARGS)
		chmod +x $(BUILD_DIR)/$(BINARY_NAME)

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
