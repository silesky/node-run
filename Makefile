# Variables
BINARY_NAME = nrun
BUILD_DIR = bin
VERSION ?= "0.0.0"

help: ## Lists all available make tasks and some short documentation about them
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'

build: build-linux-amd64 build-linux-arm64 build-mac-amd64 build-mac-arm64 ## Build all binaries
	@echo "Finished building $(BINARY_NAME)@$(VERSION)."

build-linux-amd64:
	@$(MAKE) build-platform GOOS=linux GOARCH=amd64

build-linux-arm64:
	@$(MAKE) build-platform GOOS=linux GOARCH=arm64

build-mac-amd64:
	@$(MAKE) build-platform GOOS=darwin GOARCH=amd64

build-mac-arm64:
	@$(MAKE) build-platform GOOS=darwin GOARCH=arm64 

build-platform: ## Build for a specific platform. e.g. make build-platform goos=OS goarch=architecture
	@echo "🚀 Building for $(GOOS) $(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.VERSION=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH) ./cmd/node-task-runner $(ARGS) 

run: ## Run node-task-runner without building
	go run ./cmd/node-task-runner $(ARGS)

clean: ## Clean the build directory
	rm -rf $(BUILD_DIR)

test: ## Run all tests - e.g. make test ARGS="-run TestSpecificFunction"
	go test ./... -v $(ARGS) 
		
version: ## Bump the version and push to github
	bash ./bump-version.sh

