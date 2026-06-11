.DEFAULT_GOAL := run
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X main/buildinfo.Version=$(VERSION) \
                  -X main/buildinfo.GitCommit=$(GIT_COMMIT) \
                  -X main/telemetry.TelemetryMagicString=$(TELEMETRY_MAGIC_STRING)
update:
	@go mod tidy

download:
	@go mod download

build: update
	@go build -o bin/fb-status-cli main.go

run: build
	@bin/fb-status-cli

release: download
	@mkdir -p releases/$(VERSION)
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o releases/$(VERSION)/fbstatus && \
		zip -j releases/$(VERSION)/fbstatus-$(VERSION)-mac-silicon.zip releases/$(VERSION)/fbstatus && \
		rm releases/$(VERSION)/fbstatus
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(VERSION)/fbstatus && \
		zip -j releases/$(VERSION)/fbstatus-$(VERSION)-mac-intel.zip releases/$(VERSION)/fbstatus && \
		rm releases/$(VERSION)/fbstatus
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(VERSION)/fbstatus && \
		zip -j releases/$(VERSION)/fbstatus-$(VERSION)-linux.zip releases/$(VERSION)/fbstatus && \
		rm releases/$(VERSION)/fbstatus
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(VERSION)/fbstatus.exe && \
		zip -j releases/$(VERSION)/fbstatus-$(VERSION)-windows.zip releases/$(VERSION)/fbstatus.exe && \
		rm releases/$(VERSION)/fbstatus.exe
