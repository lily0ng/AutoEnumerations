.PHONY: build install clean test run help

BINARY_NAME=autoenum
INSTALL_PATH=/usr/local/bin

help:
	@echo "AutoEnumeration - Makefile Commands"
	@echo ""
	@echo "  make build      - Build the binary"
	@echo "  make install    - Build and install to $(INSTALL_PATH)"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make run        - Build and run with example"
	@echo "  make deps       - Download dependencies"
	@echo "  make tools      - Install required Go tools"
	@echo ""

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) main.go
	@echo "Build complete: ./$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo mv $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Installation complete!"

clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf output/
	@echo "Clean complete!"

test:
	@echo "Running tests..."
	@go test -v ./...

run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME) --help

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated!"

tools:
	@echo "Installing Go tools..."
	@./scripts/install.sh
	@echo "Tools installed!"

update:
	@echo "Updating tools..."
	@./scripts/update.sh
	@echo "Tools updated!"
