FROM golang:1.24-alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o /app/bin/vixarapi ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/vixarapi /app/vixarapi
COPY --from=builder /app/config /app/config

ENTRYPOINT [ "/app/vixarapi" ]