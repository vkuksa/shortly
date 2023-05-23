# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20-alpine AS build-stage

# Set the working directory inside the container
WORKDIR /build

# Install build dependencies and golangci-lint
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin && \
    mv $(go env GOPATH)/bin/golangci-lint /usr/local/bin/

# Copy the entire project directory to the container
COPY . /build/

# Build an application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./shortly ./cmd/shortlyd

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

# Copy the binary
COPY --from=build-stage /build/bin/shortly .

# Copy the config
COPY --from=build-stage /build/shortly.conf .

# Expose the desired port for the application
EXPOSE 8080
EXPOSE 6060

ENTRYPOINT ["./shortly"]