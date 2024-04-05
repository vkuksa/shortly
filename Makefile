SHELL:=/bin/bash

.SILENT:
.DEFAULT_GOAL := run

run: down up

up:
	docker compose up --build

start: 
	docker compose start

stop:
	docker compose stop

down:
	docker compose down

test:
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker compose -f docker-compose.test.yml down --volumes

.PHONY: run up start stop down test