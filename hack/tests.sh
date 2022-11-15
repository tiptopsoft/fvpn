#/bin/sh
make build
make image
docker run --name registry  -p 3001:3000 -d star:v0.1 star registry
docker run --name star -p 3000:3000 -d star:v0.1 star edge
