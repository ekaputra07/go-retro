deps-get:
	go mod download ./...

deps-tidy:
	go mod tidy

start:
	go run .

test:
	go test ./... -cover