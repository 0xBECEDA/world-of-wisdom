
build:
	docker build -f ./docker/Dockerfile.server -t server .
	docker build -f ./docker/Dockerfile.client -t client .

run: build
	docker compose -f docker/docker-compose.yml up