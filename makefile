dev:
	@DEBUG=1 go run ./...

dry-dev:
	@DEBUG=1 DRY=1 go run ./...

build:
	@go build -o ./bin/ ./...

run:
	@go run ./...

list:
	@grep '^[^#[:space:]].*:' makefile
