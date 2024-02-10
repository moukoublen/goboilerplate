# Go service boilerplate

The boilerplate consists of:
* Go code with [chi](https://github.com/go-chi/chi) http router, [zerolog](https://github.com/rs/zerolog) logging, [knadh/koanf](github.com/knadh/koanf) config library and [testify](github.com/stretchr/testify). Basic code infra with http setup and request and response logger (configurable), logging setup, config setup and graceful shutdown.
* Makefile targets for local installation of external tools (`staticcheck`, `golangci-lint`, `goimports`, `gofumpt`, `gojq`, `air`).
* Makefile targets for linting (installing all necessary tools locally), building, testing, and local run (docker compose or native).
* [GitHub Actions](https://github.com/features/actions) for linter checks and tests on each pr.
* [Dependabot](https://docs.github.com/en/code-security/dependabot/dependabot-security-updates/configuring-dependabot-security-updates) setup for updating dependencies.
* [Visual Studio Code](https://code.visualstudio.com/) settings, tasks and launch configuration for building, debugging and testing the code.
* Docker compose and development Dockerfile with hot-rebuild (using [air](https://github.com/cosmtrek/air)) and dlv debug server.
* Production Dockerfile

## How to use
1. Click `Use this template` from this repo github page, and choose your destination repo/name (e.g. `github.com/glendale/service`)
2. Clone **your** repo locally and run `./scripts/rename` giving your new package name and main cmd name. E.g. `./scripts/rename github.com/glendale/service service`.
3. Commit the renaming and you are ready to start.


## Go package and file structure

> Every local package that collides with a golang std package is named with `x` postfix. For example `.../internal/httpx` or `.../internal/logx`. By doing this cheat you can avoid putting aliases when you need both of those and you can take advantage of package autocomplete/autoimport features more easily than having colliding names.

| Package            | Description    |
|--------------------|-----------------------------------------|
| `cmd/*/main.go`    | main function that is intended to glue together the most core level components like: configuration, http server and router, logs initialization, db connections (if any), the `App` and finally the `Main` struct (`/internal/exec.go`) that handles the long running / signal handling / graceful shutdown of the service. |
| `internal/exec.go` | the `Main` struct that wraps the long running / signal handling / graceful shutdown of the service |
| `internal/config`  | this package contain the initialization of `koanf` config. |
| `internal/httpx`   | this package contains: the setup of the `chi` router, some helpers functions for parsing/writing http request and http response and some middlewares like for logging each request/response. |
| `internal/logx`   | this package contains the setup/init function for `zerolog` logger. |

| Folder/File        | Description    |
|--------------------|-----------------------------------------|
| `deployments`      | this folder is intended to hold everything regarding deployment (e.g. helm/kubernetes etc). Inside `compose` directory a `docker-compose.yml` file is included that is intended for local development. |
| `build/docker`     | this folder contains production-like docker file as well as a local development one. |




## Make file targets


### Build, Debug and Local Deploy Targets

#### `build`
Builds all binaries under `./cmd` (one for each directory) and outputs the binaries into `./output` folder using this patter `./cmd/<directory name> -> ./output/<directory name>`.
___

#### `cmd.<directory name>`
Builds a specific cmd directory `<directory name>` and outputs the binary into `./output` folder using this pattern `./cmd/<directory name> -> ./output/<directory name>`.
___

#### `clean`
Deletes `./output` folder.
___

#### `compose-up`
Deploys (rebuilds and recreates) the docker compose file `deployments/compose/docker-compose.yml` in attached mode. This docker compose file is intended for local development.<br>
It uses the docker file `build/docker/debug.Dockerfile` and runs `air` target (see below), including hot-reload (rebuild on changes) and debug server.
___

#### `compose-down`
Stops (if started) the containers specified by the docker compose file `deployments/compose/docker-compose.yml` removes containers, cleans up the volumes and delete local docker images.
___

#### `air`
Installs (if needed) [air](https://github.com/cosmtrek/air) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `air -c .air.toml`. <br>
Air watches for code file changes and rebuilds the binary according to the configuration `.air.toml`. <br>
Current air configuration executes `cmd.goboilerplate` target (on each file change) and then runs the `./build/dlv` that starts the debug server (`dlv exec`) with the produced binary.
___

#### `image`
Builds the docker image using the docker file `./build/docker/Dockerfile`. <br>
This docker file is intended to be used as a production image. <br>
Image name and tag are specified by `IMAGE_NAME` and `IMAGE_TAG`. Default values can be overwritten during execution (eg `make IMAGE_NAME=myimage IMAGE_TAG=1.0.0 image`).
___


### Util Targets

#### `mod`
Runs go mod [tidy](https://go.dev/ref/mod#go-mod-tidy) and [verify](https://go.dev/ref/mod#go-mod-verify).
___

#### `vendor`
Runs [go mod vendor](https://go.dev/ref/mod#go-mod-vendor) that downloads all dependencies into `./vendor` folder.
___

#### `env`
Prints information regarding the local environment. More specifically: go env, all go packages and folders, and the specified tools bin directory.
___

#### `git-reset`
Full hard reset to HEAD. Cleans up all untracked files and restore all staged and un-staged changes.
___


### Tests and Linters Targets

#### `test`
Runs go [test](https://pkg.go.dev/cmd/go/internal/test) with race conditions and prints cover report.
___

#### `tools`
Installs (if needed) all tools (goimports, staticcheck, gofumpt, etc) under `TOOLS_BIN` folder (default is `./.tools/bin`).
___

#### `checks`
Runs all default checks (`vet`, `staticcheck`, `gofumpt`, `goimports`, `golangci-lint`).
___

#### `vet`
Runs go [vet](https://pkg.go.dev/cmd/vet).
___

#### `gofumpt`
Installs (if needed) [gofumpt](https://github.com/mvdan/gofumpt) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `gofumpt` for each folder. <br>
If `gofumpt` finds files that need fix/format it displays those files as list and the target fails (_exit 1_).
___

#### `gofumpt.display`
In case of `gofumpt` fails this target can be used to display the necessary changes. It runs `gofumpt -d` to display the needed changes that can be applied with `make gofumpt.fix`.
___

#### `gofumpt.fix`
Runs `gofumpt -w` to fix/format files (that need fix).
___

#### `goimports`
Installs (if needed) [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `goimports` for each folder. <br>
If `goimports` finds files that need fix/format it displays those files as list and the target fails (_exit 1_).
___

#### `goimports.display`
In case of `goimports` fails this target can be used to display the necessary changes. It runs `goimports -d` to display the needed changes that can be applied with `make goimports.fix`.
___

#### `goimports.fix`
Runs `goimports -w` to fix files (that need fix).
___

#### `staticcheck`
Installs (if needed) [staticcheck](https://staticcheck.io/) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `staticcheck` using all [checks](https://staticcheck.io/docs/checks) excluding [ST1000](https://staticcheck.io/docs/checks/#ST1000). <br>
If issues are found the target fails (due to `staticcheck` exit status).
___

#### `golangci-lint`
Installs (if needed) [golangci-lint](https://golangci-lint.run/) under `TOOLS_BIN` folder (default is `./.tools/bin`) and runs `golangci-lint`. <br>
If issues are found the target fails (due to `golangci-lint` exit status).
___

