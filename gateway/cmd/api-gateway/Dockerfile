FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY gateway/go.mod gateway/go.sum ./
RUN go mod download

COPY gateway/ .

RUN go build -o /api-gateway ./cmd/api-gateway

FROM alpine:latest
WORKDIR .
COPY --from=builder /api-gateway /app/
EXPOSE 8080
CMD ["./app/api-gateway"]