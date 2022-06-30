FROM golang:1.18.3-alpine3.16

COPY ./scripts/install-tools /install-tools

RUN mkdir /src && \
    apk update && \
    apk add --no-cache make git curl bash ca-certificates && \
    /install-tools

WORKDIR /src
