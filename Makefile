.DEFAULT_GOAL := run
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X main/buildinfo.Version=$(VERSION) \
                  -X main/buildinfo.GitCommit=$(GIT_COMMIT)
update:
	@go mod tidy

build: update
	@go build -o bin/fb-status-cli main.go

run: build
	@bin/fb-status-cli

release: update
	@mkdir -p releases/$(VERSION)
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o releases/$(VERSION)/fb-status-mac-silicon
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(VERSION)/fb-status-mac-intel
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(VERSION)/fb-status-linux
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(VERSION)/fb-status-windows.exe
