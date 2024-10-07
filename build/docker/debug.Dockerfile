# syntax=docker/dockerfile:1.10
### https://hub.docker.com/r/docker/dockerfile

# https://hub.docker.com/_/golang
FROM golang:1.23.2-alpine3.20

WORKDIR /wd

RUN <<EOS

apk update
apk add --no-cache make git bash ca-certificates

go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/air-verse/air@latest

git config --global --add safe.directory /wd

EOS
