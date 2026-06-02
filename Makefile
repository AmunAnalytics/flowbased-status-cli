.DEFAULT_GOAL := run

update:
	@go mod tidy

build: update
	@go build -o bin/fb-status-cli main.go

run: build
	@bin/fb-status-cli

