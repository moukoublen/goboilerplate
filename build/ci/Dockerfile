FROM golang:1.16.5-alpine3.14

COPY ./install-tools.bash /install-tools.bash

RUN mkdir /src && \
    apk update && \
    apk add --no-cache make git curl bash ca-certificates && \
    /install-tools.bash

WORKDIR /src
