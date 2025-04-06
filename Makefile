# Variables
BINARY_NAME = nrun
BUILD_DIR = bin

# Default target
all: build

# Variables
BINARY_NAME = nrun
BUILD_DIR = bin
VERSION ?= "0.0.0"

# Default target
all: build

# Build for all architectures
build: build-linux-amd64 build-linux-arm64 build-mac-amd64 build-mac-arm64
	@echo "Finished building $(BINARY_NAME)@$(VERSION)."

build-linux-amd64:
	@$(MAKE) build-platform GOOS=linux GOARCH=amd64

build-linux-arm64:
	@$(MAKE) build-platform GOOS=linux GOARCH=arm64

build-mac-amd64:
	@$(MAKE) build-platform GOOS=darwin GOARCH=amd64

build-mac-arm64:
	@$(MAKE) build-platform GOOS=darwin GOARCH=arm64 

# e.g. make build-platform goos=OS goarch=architecture
build-platform:
	@echo "🚀 Building for $(GOOS) $(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.VERSION=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH) ./cmd/node-task-runner $(ARGS) 

# Run the binary
run: 
	go run ./cmd/node-task-runner $(ARGS)

# Clean the build directory
clean:
	rm -rf $(BUILD_DIR)

# Run all tests - e.g. make test ARGS="-run TestSpecificFunction"
test:
	go test ./... -v $(ARGS) 

# List the make targets
help:                                                                                                                    
	@grep '^[^#[:space:]].*:' Makefile
		
# Bump the version and push to github
version:
	bash ./bump-version.sh

