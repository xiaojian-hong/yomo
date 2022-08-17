GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /example/)
CLI_VERSION ?= $(shell git describe --match "cli*" --tags 2>/dev/null | cut -c 5- || git rev-parse --short HEAD)
GO_LDFLAGS ?= -X $(shell $(GO) list -m)/cmd.Version=$(CLI_VERSION)
# The tag pattern for CLI: cli-v0.0.1, cut the first 4 characters to extract the version of CLI
VER ?= $(shell git describe --match "cli*" --tags --abbrev=0 | cut -c 5-)

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

vet:
	$(GO) vet $(VETPACKAGES)

lint:
	revive -exclude example/... -exclude cli/... -formatter friendly ./...

build-cli:
	$(GO) build -o bin/yomo -ldflags "-s -w ${GO_LDFLAGS}" ./cmd/yomo/main.go

archive-release-cli:
	rm -rf bin/yomo
	GOARCH=arm64 GOOS=darwin $(GO) build -o bin/yomo -ldflags "-s -w ${GO_LDFLAGS}" ./cmd/yomo/main.go
	tar -C ./bin -czf bin/yomo-${VER}-arm64-Darwin.tar.gz yomo
	rm -rf bin/yomo
	GOARCH=amd64 GOOS=darwin $(GO) build -o bin/yomo -ldflags "-s -w ${GO_LDFLAGS}" ./cmd/yomo/main.go
	tar -C ./bin -czf bin/yomo-${VER}-x86_64-Darwin.tar.gz yomo
	rm -rf bin/yomo
	GOARCH=arm64 GOOS=linux $(GO) build -o bin/yomo -ldflags "-s -w ${GO_LDFLAGS}" ./cmd/yomo/main.go
	tar -C ./bin -czf bin/yomo-${VER}-arm64-Linux.tar.gz yomo
	rm -rf bin/yomo
	GOARCH=amd64 GOOS=linux $(GO) build -o bin/yomo -ldflags "-s -w ${GO_LDFLAGS}" ./cmd/yomo/main.go
	tar -C ./bin -czf bin/yomo-${VER}-x86_64-Linux.tar.gz yomo
	rm -rf bin/yomo
	make bina

bina_json = '{"platforms": { "darwin-arm64": { "asset": "yomo-${VER}-arm64-Darwin.tar.gz", "file": "yomo" }, "darwin-amd64": { "asset": "yomo-${VER}-x86_64-Darwin.tar.gz", "file": "yomo" }, "linux-arm64": { "asset": "yomo-${VER}-arm64-Linux.tar.gz", "file": "yomo" }, "linux-amd64": { "asset": "yomo-${VER}-x86_64-Linux.tar.gz", "file": "yomo" } } }'
bina:
	@echo ${bina_json} > ./bin/bina.json

cli_version:
	@echo ${VER} > ./cli/VERSION

build-w-sym:
	GOARCH=amd64 GOOS=linux $(GO) build -o bin/yomo -ldflags "${GO_LDFLAGS}" -gcflags=-l ./cmd/yomo/main.go
