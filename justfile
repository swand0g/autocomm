set positional-arguments
config-file := "~/.config/autocomm/autocomm.toml"
install_file := "./install.sh"

_default:
	@just -lu --justfile {{justfile()}}

# Print app config
config:
	@cat {{config-file}}

# Build the app
build *args='':
	@go build -o ./bin/ $@ ./...

# Freshly build and run the app
build-run:
	@go build -o ./bin/ ./... && ./bin/autocomm

# Install the app with the shell script
install *args='':
	@chmod +x {{install_file}} && {{install_file}} $@

# Run the app
run *args='':
	@go run ./... $@

# Run the app in debug mode
dev:
	@go run ./... --debug

# Run the app in dry mode
dry:
	go run ./... --dry

# Run the app without making calls to OpenAI API
dry-dev:
	go run ./... --debug --dry

film:
	@vhs < ./cinema/movie.tape

# Aliases
alias br := build-run
alias dd := dry-dev
alias d := dev
alias i := install
alias r := run
