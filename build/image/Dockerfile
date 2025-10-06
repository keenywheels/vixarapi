FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

COPY . /vixarapi

WORKDIR /vixarapi

RUN go build -ldflags="-s -w" -o /vixarapi/bin/vixarapi ./cmd/vixarapi/main.go

FROM alpine:latest

ENV TZ=Europe/Moscow

EXPOSE 8000

WORKDIR /vixarapi

RUN adduser -D -g '' appuser

COPY --from=builder /usr/share/zoneinfo/${TZ} /etc/localtime
RUN echo $TZ > /etc/timezone

COPY --from=builder /vixarapi/bin/vixarapi .
