# Variables
BINARY_NAME = ntk
BUILD_DIR = bin

# Default target
all: build


# Build the binary
build:
		mkdir -p $(BUILD_DIR)
		go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/node-task-runner $(ARGS)
		chmod +x $(BUILD_DIR)/$(BINARY_NAME)

# Run the binary
run: 
		go run ./cmd/node-task-runner $(ARGS)

# Clean the build directory
clean:
		rm -rf $(BUILD_DIR)

# Run tests
test:
		go test ./... -v$(ARGS)
