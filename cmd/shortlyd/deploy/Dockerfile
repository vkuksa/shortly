FROM golang:1.21.4 AS builder

COPY go.mod go.sum /src/
WORKDIR /src

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o ./bin/app ./cmd/shortlyd

# Runtime image
FROM registry.access.redhat.com/ubi9-minimal:9.3 AS runtime

COPY --from=builder /src/bin/app app

ENTRYPOINT ["./app"]