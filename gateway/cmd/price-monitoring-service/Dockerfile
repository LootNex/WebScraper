# Этап сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY gateway/go.mod gateway/go.sum ./
RUN go mod download

# Копируем исходный код
COPY gateway/ .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o price-monitoring-service ./cmd/price-monitoring-service

# Этап запуска
FROM alpine:latest

WORKDIR /root/

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/price-monitoring-service .
COPY --from=builder /app/internal/scrapping ./internal/scrapping
COPY /gateway/.env /root/

# Копируем файлы шаблонов и статики, если они есть
# COPY --from=builder /app/templates ./templates
# COPY --from=builder /app/static ./static

# Устанавливаем зависимости для скрапинга (если нужны)
RUN apk --no-cache add ca-certificates

# Порт, который будет слушать сервис
EXPOSE 50052

# Команда для запуска сервиса
CMD ["./price-monitoring-service"]