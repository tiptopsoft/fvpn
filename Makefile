VERSION ?= latest

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
PLATFORM = linux
ARCH = amd64


build:
	bash ${shell pwd}/hack/build.sh
build-m1:
	GOPROXY=https://goproxy.cn,direct go build -v -o bin/fvpn main.go

image: build
	docker build -t registry.cn-hangzhou.aliyuncs.com/fvpn/fvpn:${tags} -f ${shell pwd}/docker/Dockerfile ${shell pwd}/bin/linux/amd64
image-push: image
	docker push registry.cn-hangzhou.aliyuncs.com/fvpn/fvpn:${tags}