#!/bin/bash

for goos in "linux" "darwin"
do
  for arch in "arm64" "amd64"
  do
    docker run --rm --env GOPROXY=https://goproxy.cn --env GOOS=$goos --env GOARCH=$arch -v "$(dirname $PWD)":/root -w /root/fvpn golang:1.20.6  go build -o /root/fvpn/bin/$goos/$arch/fvpn -v main.go
  done
done
echo "Build Success!"
