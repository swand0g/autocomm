set positional-arguments
config-file := "~/.autocomm"

_default:
  @just -lu --justfile {{justfile()}}

# Print the OpenAI API key being used
check-key:
  @cat {{config-file}}

# Build the app
build *args='':
  @go build -o ./bin/ $@ ./...

# Run the app
run:
  @go run ./...

# Run the app in debug mode
dev:
  @DEBUG=1 go run ./...

# Run the app in dry mode
dry:
  @DRY=1 go run ./...

# Run the app without making calls to OpenAI API
dry-dev:
  @DEBUG=1 DRY=1 go run ./...

alias dd := dry-dev