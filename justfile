# This!
default:
  @just -lu --justfile {{justfile()}}

# Build the app
build:
  @go build -o ./bin/ ./...

# Run the app
run:
	@go run ./...

# Run the app in debug mode
dev:
	@DEBUG=1 go run ./...

# Run the app without making calls to OpenAI API
dry-dev:
	@DEBUG=1 DRY=1 go run ./...

alias dd := dry-dev