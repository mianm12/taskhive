# 把你 go.mod 里的模块路径填这里
MODULE  := github.com/mianm12/taskhive
VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# 把版本信息拼成 ldflags,避免每个 target 重复写
LDFLAGS := -X '$(MODULE)/internal/version.Version=$(VERSION)' \
           -X '$(MODULE)/internal/version.Commit=$(COMMIT)' \
           -X '$(MODULE)/internal/version.BuildDate=$(DATE)'

.PHONY: build test lint fmt run clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/taskhive ./cmd/taskhive

test:
	go test -race ./...

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

run: build
	./bin/taskhive version

clean:
	rm -rf bin/