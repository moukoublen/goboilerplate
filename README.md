# Go service boilerplate

## How to use
1. Click `Use this template` from this repo github page, and choose your destination repo/name (e.g. `github.com/glendale/service`)
2. Clone **your** repo locally and run `./scripts/rename` giving your new package name and main cmd name. E.g. `./scripts/rename github.com/glendale/service service`.
3. Commit the renaming and you are ready to start.

## Make file targets

| **target**             | **description** |
| ---------------------- | --------------- |
| `mod`                  | Runs go mod [tidy](https://go.dev/ref/mod#go-mod-tidy) and [verify](https://go.dev/ref/mod#go-mod-verify). |
| `vendor`               | Runs [go mod vendor](https://go.dev/ref/mod#go-mod-vendor) that downloads all dependencies into `./vendor` folder. |
| `build`                | Builds all binaries under `./cmd` (one for each directory) and outputs the binaries into `./output` folder using this patter `./cmd/<directory name> -> ./output/<directory name>`. |
| `cmd.<directory name>` | Builds a specific cmd directory `<directory name>` and outputs the binary into `./output` folder using this pattern `./cmd/<directory name> -> ./output/<directory name>`. |
| `clean`                | Deletes `./output` folder. |
| `git-reset`            | Full hard reset to HEAD. Cleans up all untracked files and restore all staged and un-staged changes. |
| `image`                | Builds the docker image using the docker file `./build/docker/Dockerfile`. <br>Image name and tag are specified by `IMAGE_NAME` and `IMAGE_TAG`. Default values can be overwritten during execution (eg `make IMAGE_NAME=myimage IMAGE_TAG=1.0.0 image`). |
| `up`                   | Deploys (rebuilds and recreates) the docker compose file `deployments/compose/docker-compose.yml` in attached mode. |
| `down`                 | Stops (if started) the containers specified by the docker compose file `deployments/compose/docker-compose.yml` removes containers, cleans up the volumes and delete local docker images. |
| `env`                  | Prints information regarding the local environment. More specifically: go env, all go packages and folders, and the specified tools bin directory. |
| `test`                 | Runs go [test](https://pkg.go.dev/cmd/go/internal/test) with race conditions and prints cover report. |
| `checks`               | Runs all default tests (eg vet, lint, etc). |
| `vet`                  | Runs go [vet](https://pkg.go.dev/cmd/vet). |
| `gofumpt`              | Installs (if needed) [gofumpt](https://github.com/mvdan/gofumpt) under `TOOLSBIN` folder (default is `./.bin`) and runs `gofumpt` for each folder. <br>If `gofumpt` finds files that need fix/format it displays those files as list and the target fails (_exit 1_). |
| `gofumpt.display`      | In case of `gofumpt` fails this target can be used to display the necessary changes. It runs `gofumpt -d` to display the needed changes that can be applied with `make gofumpt.fix`. |
| `gofumpt.fix`          | Runs `gofumpt -w` to fix/format files (that need fix). |
| `goimports`            | Installs (if needed) [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) under `TOOLSBIN` folder (default is `./.bin`) and runs `goimports` for each folder. <br>If `goimports` finds files that need fix/format it displays those files as list and the target fails (_exit 1_). |
| `goimports.display`    | In case of `goimports` fails this target can be used to display the necessary changes. It runs `goimports -d` to display the needed changes that can be applied with `make goimports.fix`. |
| `goimports.fix`        | Runs `goimports -w` to fix files (that need fix). |
| `staticcheck`          | Installs (if needed) [staticcheck](https://staticcheck.io/) under `TOOLSBIN` folder (default is `./.bin`) and runs `staticcheck` using all [checks](https://staticcheck.io/docs/checks) excluding [ST1000](https://staticcheck.io/docs/checks/#ST1000). <br>If issues are found the target fails (due to `staticcheck` exit status). |
| `golangci-lint`        | Installs (if needed) [golangci-lint](https://golangci-lint.run/) under `TOOLSBIN` folder (default is `./.bin`) and runs `golangci-lint`. <br>If issues are found the target fails (due to `golangci-lint` exit status). |
| `tools`                | Installs (if needed) all tools (goimports, staticcheck, gofumpt, etc) under `TOOLSBIN` folder (default is `./.bin`). |
