run: stop up

mod:
	go mod tidy
	go mod vendor

up:
	docker compose -f docker-compose.yml up -d --build

stop:
	docker compose -f docker-compose.yml stop

down:
	docker compose -f docker-compose.yml down

# test:
# 	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
# 	docker-compose -f docker-compose.test.yml down --volumes
