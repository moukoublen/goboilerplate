FROM golang:1.20.3-alpine3.17 as builder
RUN apk update && apk add --no-cache make git bash ca-certificates
WORKDIR /app
COPY . .
RUN set -xe && \
    make build

#FROM scratch
#COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
FROM alpine:3.17
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/output/* /app/
USER nobody
CMD ["./goboilerplate"]
