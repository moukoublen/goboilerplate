FROM golang:1.16.4-alpine3.13

COPY ./install-tools.bash /install-tools.bash

RUN mkdir /src && \
    apk update && apk add --no-cache make git curl bash ca-certificates && \
    /install-tools.bash

WORKDIR /src