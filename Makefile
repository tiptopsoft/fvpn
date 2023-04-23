VERSION ?= latest
OUT_DIR = ${M}/bin
BINARY = ${M}

RELEASE_BIN = ${M}-$(VERSION)-bin
RELEASE_SRC = ${M}-$(VERSION)-src

OS = $(shell uname)

GO = go
GO_PATH = $$($(GO) env GOPATH)
GO_BUILD = $(GO) build
GO_GET = $(GO) get
GO_CLEAN = $(GO) clean
GO_TEST = $(GO) test
GO_LINT = $(GO_PATH)/bin/golangci-lint
GO_LICENSER = $(GO_PATH)/bin/go-licenser
GO_BUILD_FLAGS = -v

build:
	GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/fvpn main.go

image: build
	cd ${shell pwd}/bin/ && docker buildx build  -t fvpn:${tags} -f ${shell pwd}/docker/Dockerfile .
