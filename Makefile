VERSION ?= latest

build:
	bash ${shell pwd}/hack/build_all.sh
build-m1:
	GOPROXY=https://goproxy.cn,direct go build -v -o bin/fvpn main.go

image: build
	docker build -t registry.cn-hangzhou.aliyuncs.com/fvpn/fvpn:${tags} -f ${shell pwd}/docker/Dockerfile ${shell pwd}/bin/linux/amd64

image-push: image
	docker push registry.cn-hangzhou.aliyuncs.com/fvpn/fvpn:${tags}