#!/bin/bash
for goos in "darwin" "linux"
do
  for arch in "arm64" "amd64"
    do
      if [ ! -d "bin/dist" ];
      then
      mkdir bin/dist
      fi
      tar -cvf bin/dist/fvpn-$goos-$arch-v0.1.tar.gz bin/$goos/$arch
    done
done
