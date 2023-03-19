# bmc

BareMetalContainer: running containers without a daemon

## Overview

docker 容器在我们的开发运维过程中提供了很大的便利:
1. 可以使用 docker 配合 docker-compose 搭建开发环境，保证开发环境与生产环境一致
2. 运维中，docker 容器运行在 k8s 中，容器可以伸缩，提高服务吞吐
3. docker 容器使用资源少，镜像轻量

运维过程中，使用 ansible 安装 docker、k8s 集群。安装过程纯离线，所以需要将所有的依赖做离线处理。而 ansible 运行环境依赖与 python 以及一些 python 包。从前我们的方法是，将安装包装到 docker 容器中。部署的时候，申请两台机器，一台用做部署机，一台用于生产。部署完后，部署机回收。使用体验偏差。

docker 提供的运行时灵活，无束缚，如果将当前的 docker 完全迁移，我们需要做更多的工作在适配不同的运行环境。而 docker 运行时又是 ansible 装起来的，一个鸡生蛋还蛋蛋生鸡的问题。

基于对 docker 和 linux namespace 的理解，可以将 docker daemon 舍弃掉，选择一种更清量的方式启动一个 container。container不需要严格的网络隔离。

## Usages

```
vagrant up
curl --output bin/runc --location https://github.com/opencontainers/runc/releases/download/v1.1.3/runc.amd64
make build
docker save -o alpine.tar alpine:edge

/vagrant/bin/bmc-Linux run --build-rootfs -f /vagrant/tmp/tmp.tar --mount "type=bind,source=$(pwd)/snap,target=/host"
```

## deps

```
https://github.com/opencontainers/runc/releases/tag/v1.1.3
curl --output bin/runc.amd64 --location https://github.com/opencontainers/runc/releases/download/v1.1.3/runc.amd64
```

## TODO

- [x] 运行起程序
- [x] 支持加载环境变量
- [x] bind mount 数据目录
- [x] CI
- [x] arm64 支持
- [ ] 使用 github.com/opencontainers/runtime-spec/specs-go 重构 spec 解析
- [ ] 使用 https://0xcf9.org/2021/06/22/embed-and-execute-from-memory-with-golang/ 重构 runc 执行
