
build-srv:
	docker build -f ./docker/Dockerfile.server -t server .

build-client:
	docker build -f ./docker/Dockerfile.client -t client .

start: build-srv build-client
	docker compose -f docker/docker-compose.yml up