ARG GO_VERSION=1.12.7
# Build Stage
FROM golang:${GO_VERSION}-alpine as build_phase

WORKDIR /go/cache-server

RUN apk add git gcc libc6-compat musl-dev

COPY . .

ENV GO111MODULE=on
ENV CGO=0
RUN go build -ldflags="-s -w" -o cache-bin

# Exec Stage
FROM alpine:3.9

EXPOSE 3001

COPY --from=build_phase /go/cache-server/cache-bin .
RUN chmod +x cache-bin

ENTRYPOINT ./cache-bin -dsn=${SENTRY_DSN} -addr=":3001"
