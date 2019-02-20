# Build Stage
FROM golang:1.11.5-alpine as build_phase

WORKDIR /go/src/github.com/lucaslollobrigida/cache-server
RUN apk add --no-cache git gcc libc6-compat musl-dev
RUN go get github.com/valyala/fasthttp
COPY . .
ENV GO111MODULE=on
RUN go build -o cache-bin

# Exec Stage
FROM alpine:3.9

EXPOSE 3001

COPY --from=build_phase /go/src/github.com/lucaslollobrigida/cache-server/cache-bin .
RUN chmod +x cache-bin
CMD ["./cache-bin"]

