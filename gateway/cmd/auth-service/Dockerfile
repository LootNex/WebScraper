FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY gateway/go.mod gateway/go.sum ./
RUN go mod download

COPY gateway/ .

RUN go build -o /auth-service ./cmd/auth-service

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /auth-service .
EXPOSE 50051
CMD ["./auth-service"]