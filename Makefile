run:
	ENV=local PORT=8080 DB_USERNAME=root DB_PASSWORD=password DB_HOST=localhost DB_PORT=3306 DB_NAME=restapi go run main.go

run-dev:
	ENV=dev PORT=8080 DB_USERNAME=root DB_PASSWORD=password DB_HOST=localhost DB_PORT=3306 DB_NAME=restapi go run main.go

test:
	go test -v --cover ./...

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

check:
	golangci-lint run --timeout 5m
	go fmt ./...
	make test