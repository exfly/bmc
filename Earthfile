VERSION --shell-out-anywhere --use-copy-link 0.7

FROM golang:1.20-alpine
WORKDIR /build
RUN apk add --no-cache curl

build:
    COPY . .
	ARG TARGETARCH
	ARG RUNC_VERSION=v1.1.3
	ENV VERSION dev
	RUN env
    RUN mkdir -p bin
	RUN curl --show-error --fail -v --output bin/runc --location https://github.com/opencontainers/runc/releases/download/${RUNC_VERSION}/runc.${TARGETARCH}
	RUN GOPROXY=https://goproxy.io,direct GOSUMDB=sum.golang.google.cn GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /build/bmc -ldflags "-X github.com/exfly/bmc.Version=${VERSION}" cmd/bmc/main.go
    SAVE ARTIFACT /build/bmc AS LOCAL bmc-${TARGETARCH}

build-amd64:
	BUILD --platform=linux/amd64 +build

build-all-platforms:
	BUILD --platform=linux/amd64 --platform=linux/arm64 +build
