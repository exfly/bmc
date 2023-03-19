SHELL := /bin/bash

export GO111MODULE=on
export GOPROXY=https://goproxy.io,direct
export GOSUMDB=sum.golang.google.cn
export GOPRIVATE=github.com/exfly
export GONOPROXY=github.com/exfly

export CGO_ENABLED := 0
export CI_COMMIT_SHA=$(shell git log -1 --pretty=%H)
export CI_COMMIT_REF_NAME=$(shell git branch --show-current)
export VERSION="${CI_COMMIT_REF_NAME}-$(shell echo ${CI_COMMIT_SHA} | cut -c 1-8)"


.PHONY: generate
generate:
	go generate -v ./...

.PHONY: build
build:
	mkdir -p bin
	GOOS=linux go build -x -o bin/bmc-Linux -ldflags "-X github.com/exfly/bmc.Version=${VERSION}" cmd/bmc/main.go
	GOOS=darwin go build -x -o bin/bmc-Darwin -ldflags "-X github.com/exfly/bmc.Version=${VERSION}" cmd/bmc/main.go

.PHONY: release
release:
	docker build -t exfly/bmc:dev --build-arg VERSION=${VERSION} .
	docker build -t exfly/skopeo:dev -f Dockerfile.example --build-arg VERSION=${VERSION} .

.PHONY: builddemoimg
builddemoimg:
	docker build --squash -t tmp:dev -f Dockerfile.example . && docker save -o tmp/tmp.tar tmp:dev
	#Then run in vm: /vagrant/bin/bmc-Linux run --build-rootfs -f /vagrant/tmp/tmp.tar

.PHONY: earthfile
earthfile:
	docker run -t -v ${PWD}:/workspace -v /var/run/docker.sock:/var/run/docker.sock -e NO_BUILDKIT=1 earthly/earthly:v0.7.2 +build-all-platforms
