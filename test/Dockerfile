FROM golang:alpine

RUN set -ex; \
    apk update; \
    apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./

ARG GOMODCACHE GOCACHE
RUN --mount=type=cache,target="$GOMODCACHE" go mod download
RUN go mod download

COPY . .

# TODO: add tags
CMD CGO_ENABLED=0 go test -v ./...
