# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN go build -ldflags "-X main.Version=${VERSION}" -o monocloud-config .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/monocloud-config .

# CONFIG_PATH points to the shared volume with mihomo
ENV CONFIG_PATH=/etc/mihomo/config.yaml
# UPDATE_INTERVAL controls how often the config is refreshed (e.g. 1h, 30m)
ENV UPDATE_INTERVAL=1h

ENTRYPOINT ["./monocloud-config"]
