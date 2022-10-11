#!/bin/bash

for goos in "linux" "darwin"
do
  for arch in "arm64" "amd64"
  do
    GOOS=$goos GOARCH=$arch go build -o bin/$goos/$arch/fvpn main.go
  done
done
echo "Build Success!"
