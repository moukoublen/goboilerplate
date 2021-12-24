FROM golang:1.17.5-alpine3.15

COPY ./scripts/install-tools /install-tools

RUN mkdir /src && \
    apk update && \
    apk add --no-cache make git curl bash ca-certificates && \
    /install-tools

WORKDIR /src
