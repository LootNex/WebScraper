# Этап сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY gateway/go.mod gateway/go.sum ./
RUN go mod download

# Копируем исходный код
COPY gateway/ .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o telegram-bot ./cmd/telegram-bot

# Этап запуска
FROM alpine:latest

WORKDIR /root/

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/telegram-bot .

# Устанавливаем зависимости (если нужны)
RUN apk --no-cache add ca-certificates

# Команда для запуска бота
CMD ["./telegram-bot"]