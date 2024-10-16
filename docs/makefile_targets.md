# Makefile targets

| Target               | Description |
| -------------------- | ----------- |
| `mod`                | Runs go mod [tidy](https://go.dev/ref/mod#go-mod-tidy) and [verify](https://go.dev/ref/mod#go-mod-verify). |
| `vendor`             | Runs [go mod vendor](https://go.dev/ref/mod#go-mod-vendor) that downloads all dependencies into `./vendor` folder. |
| `go-deps-upgrade`    | Runs updates all `go.mod` dependencies. |
| `env`                | Prints information regarding the local environment. More specifically: go env, all go packages and folders, and the specified tools bin directory. |
| `git-reset`          | Full hard reset to HEAD. Cleans up all untracked files and restore all staged and un-staged changes. |
| `build`              | Builds all binaries under `./cmd` (one for each directory) and outputs the binaries into `./output` folder using this pattern: `./cmd/<directory name> -> ./output/<directory name>`. |
| `cmd.<cmd dir name>` | Builds a specific cmd directory `<directory name>` and outputs the binary into `./output` folder using this pattern: `./cmd/<directory name> -> ./output/<directory name>`.
| `dbg.<cmd dir name>` | Builds a specific cmd directory `<directory name>` (keeping debug symbols) and outputs the binary into `./output` folder using this pattern: `./cmd/<directory name> -> ./output/<directory name>`.
| `clean`              | Deletes `./output` folder.
| `compose-up`         | Deploys (rebuilds and recreates) the docker compose file `deployments/compose/docker-compose.yml` in attached mode. This docker compose file is intended for local development.<br>It uses the docker file `build/docker/debug.Dockerfile` and runs `air` target (see below), including hot-reload (rebuild on changes) and debug server.
| `compose-down`       | Stops (if started) the containers specified by the docker compose file `deployments/compose/docker-compose.yml` removes containers, cleans up the volumes and delete local docker images.
| `air`                | Installs (if needed) [air](https://github.com/air-verse/air) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `air -c .air.toml`.<br>Air watches for code file changes and rebuilds the binary according to the configuration `.air.toml`.<br>Current air configuration executes `cmd.goboilerplate` target (on each file change) and then runs the `./build/dlv` that starts the debug server (`dlv exec`) with the produced binary.
| `build-image`        | Builds the docker image using the docker file `./build/docker/Dockerfile`.<br>This docker file is intended to be used as a production image.<br>Image name and tag are specified by `IMAGE_NAME` and `IMAGE_TAG`. Default values can be overwritten during execution (eg `make IMAGE_NAME=myimage IMAGE_TAG=1.0.0 image`).
| `test`               | Runs go [test](https://pkg.go.dev/cmd/go/internal/test) with race conditions and prints cover report.
| `tools`              | Installs (if needed) all tools (goimports, staticcheck, gofumpt, etc) under `TOOLS_BIN` folder (default is `./.tools/bin`).
| `checks`             | Runs all default checks (`vet`, `staticcheck`, `gofumpt`, `goimports`, `golangci-lint`).
| `vet`                | Runs go [vet](https://pkg.go.dev/cmd/vet).
| `gofumpt`            | Installs (if needed) [gofumpt](https://github.com/mvdan/gofumpt) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `gofumpt` for each folder.<br>If `gofumpt` finds files that need fix/format it displays those files as list and the target fails (_exit 1_).
| `gofumpt.fix`        | Runs `gofumpt -w` to fix/format files (that need fix).
| `goimports`          | Installs (if needed) [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `goimports` for each folder.<br>If `goimports` finds files that need fix/format it displays those files as list and the target fails (_exit 1_).
| `goimports.fix`      | Runs `goimports -w` to fix files (that need fix).
| `staticcheck`        | Installs (if needed) [staticcheck](https://staticcheck.io/) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `staticcheck` using all [checks](https://staticcheck.io/docs/checks) excluding [ST1000](https://staticcheck.io/docs/checks/#ST1000).<br>If issues are found the target fails (due to `staticcheck` exit status).
| `golangci-lint`      | Installs (if needed) [golangci-lint](https://golangci-lint.run/) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `golangci-lint`. <br>If issues are found the target fails (due to `golangci-lint` exit status).

