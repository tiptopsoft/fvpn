VERSION ?= latest

build:
	docker run --rm --env GOPROXY=https://goproxy.cn -v "$(shell pwd)":/root/fvpn -w /root/fvpn golang:1.20.6 bash /root/fvpn/hack/build_all.sh
build-m1:
	GOPROXY=https://goproxy.cn,direct go build -v -o bin/fvpn main.go

image: build
	docker build -t registry.cn-hangzhou.aliyuncs.com/fvpn/fvpn:${tags} -f ${shell pwd}/docker/Dockerfile ${shell pwd}/bin/linux/amd64

dist:
	bash ${shell pwd}/hack/dist.sh

image-push: image
	docker push registry.cn-hangzhou.aliyuncs.com/fvpn/fvpn:${tags}