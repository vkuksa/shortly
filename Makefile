SHELL:=/bin/bash

.SILENT:
.DEFAULT_GOAL := run

export GOCACHE := /go/cache
export GOMODCACHE := /go/pkg/mod

run: down up

up:
	docker compose build --build-arg GOCACHE=${GOCACHE} --build-arg GOMODCACHE=${GOMODCACHE}
	docker compose up --build -d
	docker compose logs -f shortlyd

start: 
	docker compose start

stop:
	docker compose stop

down:
	docker compose down

test:
	docker compose -f test/docker-compose.yml build --build-arg GOCACHE=${GOCACHE} --build-arg GOMODCACHE=${GOMODCACHE}
	docker compose -f test/docker-compose.yml up --build --abort-on-container-exit
	docker compose -f test/docker-compose.yml stop shortly-svc
	docker compose -f test/docker-compose.yml down --volumes

lint:
	golangci-lint run 

.PHONY: run up start stop down test lint