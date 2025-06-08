FROM golang:1.23.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/auth-service .
COPY --from=builder /app/config ./config

EXPOSE 5711 2111

CMD ["./auth-service", "--config=./config/config.yaml"]