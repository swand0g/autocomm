dev:
	@DEBUG=1 go run ./...
build:
	@go build -o ./bin/ ./...
run:
	@go run ./...