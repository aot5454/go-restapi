run:
	ENV=local go run main.go

run-dev:
	ENV=dev go run main.go

run-prod:
	ENV=production go run main.go

test:
	go test --cover ./...

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