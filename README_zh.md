# BMC

[bmc](https://github.com/exfly/bmc) 是基于 runc 的轻量级容器运行时环境。它允许您在不使用 dockerd 的情况下以 tar 格式运行 Docker 镜像。该项目的主要目标是为离线软件安装和运行提供便利，满足开发和运维的需求。

## 特点

- 以 tar 格式运行 Docker 镜像，无需使用 dockerd。
- 支持自定义环境变量和数据目录。
- 基于 Docker 和 Linux 命名空间，不需要严格的网络隔离。
- 轻量级镜像和小的运行时环境，具有在资源消耗和镜像大小方面的优势。

## 使用方法

1. 安装依赖 https://github.com/exfly/bmc/releases
2. 在终端运行 `mv bmc-* /bin/bmc && chmod +x /bin/bmc`
3. 运行 `docker save -o dockerimg.tar alpine:edge`
4. 在终端运行 `bmc run --build-rootfs -f /path/to/dockerimg.tar --mount "type=bind,source=$(pwd)/snap,target=/host"`

## 构建

1. 安装依赖：[earthly](https://earthly.dev/get-earthly) [docker](https://docs.docker.com/desktop/install/mac-install/)
2. 在终端运行 `make earthfile` 

## 待办事项

- [x] 实现程序功能。
- [x] 支持加载自定义环境变量。
- [x] 支持绑定挂载数据目录。
- [x] 支持 CI 测试。
- [x] 支持 arm64 体系结构。
- [ ] 使用 https://github.com/opencontainers/runtime-spec/specs-go 重构规范解析。
- [ ] 使用 https://0xcf9.org/2021/06/22/embed-and-execute-from-memory-with-golang/ 重构 runc 执行。

## 贡献

欢迎通过拉取请求或问题进行贡献。如果您发现错误，请提交新问题。如果您想添加新功能或修复问题，请 fork 该项目并提交拉取请求。

## 许可证

该项目基于 Apache 许可证。有关详细信息，请参见 [LICENSE](LICENSE) 文件。
