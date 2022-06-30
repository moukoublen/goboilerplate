FROM golang:1.18.3-alpine3.16 as builder
RUN apk update && apk add --no-cache make git bash ca-certificates
WORKDIR /app
COPY . .
RUN set -xe && \
        make build

#FROM scratch
#COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
FROM alpine:3.15
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/goboilerplate /app/goboilerplate
USER nobody
CMD ["./goboilerplate"]
