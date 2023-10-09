
build-srv:
	docker build -f ./docker/Dockerfile.server -t server .

build-cli:
	docker build -f ./docker/Dockerfile.client -t client .

start-server-local:
	go run ./cmd/server

start-cli-local:
	go run ./cmd/client

start-app-docker: build-srv build-cli
	cd docker && docker compose up