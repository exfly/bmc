# BMC

[bmc](https://github.com/exfly/bmc) is a lightweight container runtime environment based on runc. It allows you to run Docker images in tar format without using dockerd. The main goal of this project is to provide convenience for offline software installation and runtime, serving the needs of both development and operations.

## Features

- Runs Docker images in tar format without using dockerd.
- Supports custom environment variables and data directories.
- Based on Docker and Linux namespace, it doesn't require strict network isolation.
- Lightweight images and small runtime environment, with advantages in resource consumption and image size.

## Usage

1. Install dependencies: https://github.com/exfly/bmc/releases
2. `mv bmc-* /bin/bmc && chmod +x /bin/bmc`
3. `docker save -o dockerimg.tar alpine:edge`
3. `bmc run --build-rootfs -f /path/to/dockerimg.tar --mount "type=bind,source=$(pwd)/snap,target=/host"`

## Build

1. Install dependencies: [earthly](https://earthly.dev/get-earthly) [docker](https://docs.docker.com/desktop/install/mac-install/)
2. `make earthfile`

## TODO List

- [x] Implement program features.
- [x] Support loading custom environment variables.
- [x] support bind mount data directory.
- [x] Support CI testing.
- [x] Support arm64.
- [ ] Refactor spec parsing using https://github.com/opencontainers/runtime-spec/specs-go.
- [ ] Refactor runc execution using https://0xcf9.org/2021/06/22/embed-and-execute-from-memory-with-golang/.

## Contributing

Contributions are welcome through pull requests or issues. If you find a bug, please open a new issue. If you want to add new features or fix issues, please fork the project and submit a pull request.

## License

This project is licensed under the Apache License. See the [LICENSE](LICENSE) file for details.
