#!/usr/bin/env bash

set -e

WAIT_VERSION="2.9.0"
MIGRATE_VERSION="4.14.1"
STATICCHECK_VERSION="2020.2.4"
GOLANGCI_LINT_VERSION="1.40.1"
DOCKER_VERSION="20.10.6"
DOCKER_COMPOSE_VERSION="1.29.2"

echo "Installing goimports and golint ..."
go get -u golang.org/x/tools/cmd/goimports
go get -u golang.org/x/lint/golint

echo "Installing wait ..."
curl --fail --silent --show-error --location https://github.com/ufoscout/docker-compose-wait/releases/download/${WAIT_VERSION}/wait -o /usr/local/bin/wait
chmod +x /usr/local/bin/wait

echo "Installing migrate ..."
curl --fail --silent --show-error --location https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz -o /tmp/migrate.tar.gz
tar --extract --gzip --file=/tmp/migrate.tar.gz -C /usr/local/bin
mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate
chmod +x /usr/local/bin/migrate

echo "Installing staticcheck ..."
curl --fail --silent --show-error --location https://github.com/dominikh/go-tools/releases/download/${STATICCHECK_VERSION}/staticcheck_linux_amd64.tar.gz -o /tmp/staticcheck.tar.gz
tar --extract --gzip --strip-components=1 --file=/tmp/staticcheck.tar.gz -C /tmp
mv /tmp/staticcheck /usr/local/bin
chmod +x /usr/local/bin/staticcheck

echo "Installing golangci-lint ..."
curl --fail --silent --show-error --location https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz -o /tmp/golangci-lint.tar.gz
tar --extract --gzip --strip-components=1 --file=/tmp/golangci-lint.tar.gz -C /tmp
mv /tmp/golangci-lint /usr/local/bin
chmod +x /usr/local/bin/golangci-lint

echo "Installing docker cli ..."
curl --fail --silent --show-error --location https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz -o /tmp/docker.tgz
mkdir -p /tmp/dockerin
tar --extract --gzip --file=/tmp/docker.tgz -C /tmp/dockerin
cp /tmp/dockerin/docker/* /usr/local/bin

echo "Installing docker-compose ..."
curl --fail --silent --show-error --location https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-Linux-x86_64 -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

echo "All tools are installed"
rm -rf /tmp/*