#!/bin/bash

for goos in "linux" "darwin" "windows"
do
  for arch in "arm64" "amd64"
  do
    go build -o /root/fvpn/bin/$goos/$arch/fvpn -v main.go
  done
done
echo "Build Success!"
