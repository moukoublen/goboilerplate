#!/usr/bin/env bash
set -e


__log_installing() {
  echo -e "Installing \e[1;36m${1}\e[0m"
}

__install() {
  install $* -o 0 -g 0 -m 0755 -t /usr/local/bin/
}

__install_tool() {
  local NAME=$1
  local URL=$2
  local DWN_FILE=$3
  local TMP_DIRECTORY=$(mktemp -d -t "${NAME}_XXXXXXX")

  __log_installing "${NAME}"

  echo -e "    Downloading    \e[34m${URL}\e[0m into \e[2;34m${TMP_DIRECTORY}/${DWN_FILE}\e[0m"
  curl --fail --silent --show-error --location "${URL}" --output "${TMP_DIRECTORY}/${DWN_FILE}"

  printf "    Installing...  "
  "__install_${NAME}" "${TMP_DIRECTORY}"

  rm -rf $TMP_DIRECTORY
  printf "\e[1;32mDone!\e[0m\n\n"
}



__VERSION_GOIMPORTS="0.1.4"
__VERSION_GOLINT="0.0.0-20210508222113-6edffad5e616"
__VERSION_WAIT="2.9.0"
__VERSION_MIGRATE="4.14.1"
__VERSION_STATICCHECK="2021.1"
__VERSION_GOLANGCI_LINT="1.41.1"
__VERSION_DOCKER="20.10.7"
__VERSION_DOCKER_COMPOSE="1.29.2"


__log_installing "goimports"
go install "golang.org/x/tools/cmd/goimports@v${__VERSION_GOIMPORTS}"

__log_installing "golint"
go install "golang.org/x/lint/golint@v${__VERSION_GOLINT}"


__install_wait() {
  __install "${1}/wait"
}
__install_tool "wait" \
  "https://github.com/ufoscout/docker-compose-wait/releases/download/${__VERSION_WAIT}/wait" \
  "wait"


__install_migrate() {
  tar --extract --gzip --file="${1}/migrate.tar.gz" -C "${1}"
  mv "${1}/migrate.linux-amd64" "${1}/migrate"
  __install "${1}/migrate"
}
__install_tool "migrate" \
  "https://github.com/golang-migrate/migrate/releases/download/v${__VERSION_MIGRATE}/migrate.linux-amd64.tar.gz" \
  "migrate.tar.gz"


__install_staticcheck() {
  tar --extract --gzip --strip-components=1 --file="${1}/staticcheck.tar.gz" -C "${1}"
  __install "${1}/staticcheck"
}
__install_tool "staticcheck" \
  "https://github.com/dominikh/go-tools/releases/download/${__VERSION_STATICCHECK}/staticcheck_linux_amd64.tar.gz" \
  "staticcheck.tar.gz"


__install_golangci-lint() {
  tar --extract --gzip --strip-components=1 --file="${1}/golangci-lint.tar.gz" -C "${1}"
  __install "${1}/golangci-lint"
}
__install_tool "golangci-lint" \
  "https://github.com/golangci/golangci-lint/releases/download/v${__VERSION_GOLANGCI_LINT}/golangci-lint-${__VERSION_GOLANGCI_LINT}-linux-amd64.tar.gz" \
  "golangci-lint.tar.gz"


__install_docker-cli() {
  mkdir -p "${1}/extr"
  tar --extract --gzip --file="${1}/docker.tgz" -C "${1}/extr"
  __install ${1}/extr/docker/*
}
__install_tool "docker-cli" \
  "https://download.docker.com/linux/static/stable/x86_64/docker-${__VERSION_DOCKER}.tgz" \
  "docker.tgz"


__install_docker-compose() {
  __install "${1}/docker-compose"
}
__install_tool "docker-compose" \
  "https://github.com/docker/compose/releases/download/${__VERSION_DOCKER_COMPOSE}/docker-compose-Linux-x86_64" \
  "docker-compose"


echo "Done"