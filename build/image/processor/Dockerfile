FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

COPY go.mod go.sum /processor/

WORKDIR /processor

RUN go mod download
COPY . /processor
RUN go build -ldflags="-s -w" -o /processor/bin/processor ./cmd/processor/main.go

FROM alpine:latest

ENV TZ=Europe/Moscow

EXPOSE 8000

WORKDIR /processor

COPY --from=builder /usr/share/zoneinfo/${TZ} /etc/localtime
RUN echo $TZ > /etc/timezone

COPY --from=builder /processor/bin/processor .
