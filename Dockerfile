# syntax=docker/dockerfile:1
# docker.io/docker/dockerfile:1

FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN env \
	GOSUMDB=sum.golang.google.cn \
	GOPROXY=https://goproxy.io,direct \
	go mod download -x

COPY . .
ARG VERSION ${VERSION}
RUN env \
	GO111MODULE=on \
	CGO_ENABLED=0 \
	go build -o /bin/bmc -ldflags "-X github.com/exfly/bmc.Version=${VERSION}" cmd/bmc/main.go

FROM alpine:edge
COPY --from=build /bin/bmc /bin/bmc
