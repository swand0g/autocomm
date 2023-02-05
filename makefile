dev:
	@DEBUG=1 go run ./...

dry-dev:
	@DEBUG=1 DRY=1 go run ./...

run:
	@go run ./...

list:
	@grep '^[^#[:space:]].*:' makefile

build:
	@go build -o ./bin/ ./...

build-linux:
	@GOOS=linux GOARCH=amd64 go build -o ./bin/ ./...

buiild-darwin:
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/ ./...