# Go service boilerplate

The boilerplate consists of:
* Go code with
    * [chi](https://github.com/go-chi/chi) http router.
    * [slog](https://pkg.go.dev/log/slog) logging,
    * [knadh/koanf](github.com/knadh/koanf) config library.
    * [testify](github.com/stretchr/testify).
    * Basic code infra with http setup and request and response logger (configurable), logging setup, config setup and graceful shutdown.
    * Optional verbose logging of http traffic, with payload,
* Makefile targets for local installation of external tools (`staticcheck`, `golangci-lint`, `goimports`, `gofumpt`, `gojq`, `air`).
* Makefile targets for linting (installing all necessary tools locally), building, testing, and local run (docker compose or native).
* [GitHub Actions](https://github.com/features/actions) for linter checks and tests on each pr.
* [Dependabot](https://docs.github.com/en/code-security/dependabot/dependabot-security-updates/configuring-dependabot-security-updates) setup for updating dependencies.
* [Visual Studio Code](https://code.visualstudio.com/) settings, tasks and launch configuration for building, debugging and testing the code.
* Docker compose and development Dockerfile with hot-rebuild (using [air](https://github.com/air-verse/air)) and dlv debug server.
* Production Dockerfile.

## How to use
1. Click `Use this template` from this repo github page, and choose your destination repo/name (e.g. `github.com/glendale/service`)
2. Clone **your** repo locally and run `./scripts/rename` giving your new package name and main cmd name. E.g. `./scripts/rename github.com/glendale/service service`.
3. Commit the renaming and you are ready to start.


## Go package and file structure
The project structure attempts to be on a par with [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

> Every local package that collides with a golang std package is named with `x` postfix. For example `.../internal/httpx` or `.../internal/logx`. By doing this cheat you can avoid putting aliases when you need both of those and you can take advantage of package autocomplete/autoimport features more easily than having colliding names.

| Package            | Description    |
|--------------------|-----------------------------------------|
| `cmd/*/main.go`    | main function that is intended to glue together the most core level components like: configuration, http server and router, logs initialization, db connections (if any), the `App` and finally the `Main` struct (`/internal/exec.go`) that handles the long running / signal handling / graceful shutdown of the service. |
| `internal/exec.go` | the `Main` struct that wraps the long running / signal handling / graceful shutdown of the service |
| `internal/config`  | this package contain the initialization of `koanf` config. |
| `internal/httpx`   | this package contains: the setup of the `chi` router, some helpers functions for parsing/writing http request and http response. |
| `internal/logx`   | this package contains the setup/init function for `slog` logger. |

| Folder/File        | Description    |
|--------------------|-----------------------------------------|
| `deployments`      | this folder is intended to hold everything regarding deployment (e.g. helm/kubernetes etc). Inside `local` directory a `docker-compose.yml` file is included that is intended for local development. |
| `build/docker`     | this folder contains production-like docker file as well as a local development one. |

## Makefile targets
Makefile targets can be found in [docs/makefile_targets.md](docs/makefile_targets.md) file.
