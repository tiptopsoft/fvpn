#!/bin/bash
set -x
set -e
export GOPROXY=https://goproxy.cn,direct
export CGO_ENABLED=0

function build() {
#    var goos = $1
#    var arch = $2
  if [ $# -lt 2 ]; then
    echo "should provide 2 params";
    exit 1;
  fi
  GOOS=$1 GOARCH=$2 go build -o bin/$1/$2/fvpn main.go
}

function buildImage() {
  docker buildx -t fpvn.cc/fvpn/fvpn:${tags} -f ${shell pwd}/docker/Dockerfile .
}

## build linux amd64
build linux amd64
## build linux arm64
build linux arm64

## build darwin amd64
build darwin amd64
## build darwin arm64
build darwin arm64

## build other platform