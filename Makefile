# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o bin/sekai-core cmd/app/main.go

# Run the application
run:
	@go run cmd/app/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bin/sekai-core
