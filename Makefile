run:
	ENV=local PORT=8080 go run main.go

test:
	go test -v --cover ./...