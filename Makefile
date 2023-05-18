# Go variables
GO := go
GOFLAGS := -v
GOTEST := $(GO) test $(GOFLAGS)
GOBUILD := $(GO) build $(GOFLAGS)
GOCLEAN := $(GO) clean
GOMOD := $(GO) mod

# Default target
all: clean fetch lint test build

# Run tests
test:
	$(GOTEST) -cover ./...

# Build the project
build:
	$(GOBUILD) -o bin/shortly ./cmd/shortly/main.go

# Clean the build artifacts
clean:
	$(GOCLEAN)
	rm -f myapp

# Linting
lint:
	golangci-lint run

# Fetch modules 
fetch:
	$(GOMOD) download
	$(GOMOD) tidy