# Этап сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY gateway/go.mod gateway/go.sum ./
RUN go mod download

COPY gateway/ .

RUN go build -o /telegram-bot ./cmd/telegram-bot

FROM alpine:latest

WORKDIR .

COPY --from=builder /telegram-bot /app/

RUN apk --no-cache add ca-certificates

CMD ["./app/telegram-bot"]