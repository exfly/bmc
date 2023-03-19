VERSION --shell-out-anywhere --use-copy-link 0.6

FROM golang:1.20-alpine
WORKDIR /build
RUN apk add --no-cache curl
ARG TARGETARCH
ARG RUNC_VERSION v1.1.3

build:
    COPY . .
	ENV VERSION dev
    RUN mkdir -p bin
	# RUN curl --output bin/runc --location https://github.com/opencontainers/runc/releases/download/${RUNC_VERSION}/runc.${TARGETARCH}
	RUN curl --output bin/runc --location https://cdn.jsdelivr.net/gh/opencontainers/runc@releases/download/${RUNC_VERSION}/runc.${TARGETARCH}
	RUN GOPROXY=https://goproxy.io,direct GOSUMDB=sum.golang.google.cn GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -x -o /build/bmc -ldflags "-X github.com/exfly/bmc.Version=${VERSION}" cmd/bmc/main.go
    SAVE ARTIFACT /build/bmc AS LOCAL bin/bmc-${TARGETARCH}

build-all-platforms:
	BUILD --platform=linux/amd64 --platform=linux/arm64 +build
