# syntax=docker/dockerfile:1.12
### https://hub.docker.com/r/docker/dockerfile

# https://docs.docker.com/build/guide/mounts/

# https://hub.docker.com/_/golang
FROM --platform=$BUILDPLATFORM golang:1.24.2-alpine3.21 AS builder
RUN apk update && apk add --no-cache make git bash ca-certificates
WORKDIR /app
ARG TARGETOS TARGETARCH
RUN --mount=type=bind,target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH BUILD_OUTPUT=/tmp/out make build

#FROM scratch
#COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /tmp/out/* /app/
USER nobody
CMD ["./goboilerplate"]
