# Project name
BINARY_NAME=bobacom

# Installation directory
INSTALL_DIR=$(HOME)/.local/bin

# Build the project and create an executable
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME)

# Install the binary to ~/.local/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	# Ensure the directory exists
	@mkdir -p $(INSTALL_DIR)
	# Move the binary to the install directory
	@mv $(BINARY_NAME) $(INSTALL_DIR)/

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# Convenience target for building, installing, and cleaning up in one step
all: clean build install

