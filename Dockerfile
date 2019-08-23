ARG GO_VERSION=1.12.9
# Build Stage
FROM golang:${GO_VERSION}-alpine as build_phase

WORKDIR /go/cache-server

RUN adduser -D -g '' appuser

RUN apk update && apk upgrade \
        && apk add --no-cache git gcc libc6-compat musl-dev ca-certificates tzdata \
        && update-ca-certificates

COPY . .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o cache-bin
# RUN chmod +x cache-bin

# Exec Stage
FROM scratch

COPY --from=build_phase /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=Build_phase /etc/passwd /etc/passwd
COPY --from=build_phase /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build_phase /go/cache-server/cache-bin /app/cache-bin

USER appuser

EXPOSE 3001

ENTRYPOINT ["/app/cache-bin"]
