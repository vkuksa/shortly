# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20-alpine AS build-stage

# Set the working directory inside the container
WORKDIR /build

# Copy the entire project directory to the container
COPY . /build/

# Build an application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./shortlyd ./cmd/shortlyd

# Deploy the application binary into a lean image
FROM gcr.io/distroless/static-debian11 AS build-release-stage 

WORKDIR /

# Copy the binary
COPY --from=build-stage /build/shortlyd .

# Copy the config
COPY --from=build-stage /build/shortly.conf .

# Expose the desired port for the application
EXPOSE 8080
EXPOSE 6060

USER nonroot:nonroot

ENTRYPOINT ["./shortlyd"]