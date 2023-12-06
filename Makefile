SHELL:=/bin/bash

.SILENT:
.DEFAULT_GOAL := run

compose_file := ./cmd/shortlyd/deploy/local/docker-compose.yml

run: down up

up:
	docker compose -f ${compose_file} up --build

start: 
	docker compose -f ${compose_file} start

stop:
	docker compose -f ${compose_file} stop

down:
	docker compose -f ${compose_file} down

# test:
# 	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
# 	docker-compose -f docker-compose.test.yml down --volumes

.PHONY: run up start stop down