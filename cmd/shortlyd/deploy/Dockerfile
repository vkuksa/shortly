FROM --platform=$BUILDPLATFORM golang:alpine AS build

COPY go.mod go.sum /src/
WORKDIR /src

ARG GOMODCACHE GOCACHE
RUN --mount=type=cache,target="$GOMODCACHE" go mod download

COPY . .

ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target="$GOMODCACHE" \
    --mount=type=cache,target="$GOCACHE" \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -buildvcs=false -o ./bin/app ./cmd/shortlyd


FROM registry.access.redhat.com/ubi9-minimal:9.3 AS runtime

COPY --from=build /src/bin/app app

ENTRYPOINT ["./app"]