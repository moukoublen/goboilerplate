# syntax=docker/dockerfile:1.6
### https://hub.docker.com/r/docker/dockerfile

# https://hub.docker.com/_/golang
FROM golang:1.22.1-alpine3.19

WORKDIR /wd

RUN <<EOS

apk update
apk add --no-cache make git bash ca-certificates

go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/cosmtrek/air@latest

git config --global --add safe.directory /wd

EOS
