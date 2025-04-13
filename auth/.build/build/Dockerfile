FROM golang:1.24-alpine AS builder

WORKDIR /auth

COPY auth/go.mod .

COPY auth/go.sum .

RUN go mod download

COPY /auth .

RUN go build -o="./bin/app/auth" cmd/auth/auth.go

RUN apk update --no-cache

RUN apk add --no-cache tzdata

RUN apk add --no-cache wget

RUN wget https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.11/grpc_health_probe-linux-amd64
RUN chmod +x grpc_health_probe-linux-amd64
RUN mv grpc_health_probe-linux-amd64 /usr/local/bin/grpc_health_probe

FROM alpine:latest

WORKDIR .

COPY --from=builder ./auth/bin/app/auth /app/bin/

COPY /auth/migrations/ /app/migrations/

COPY /auth/.build/config/ /app/config/

COPY --from=builder /usr/local/bin/grpc_health_probe /usr/local/bin/

COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow

ENV TZ=Europe/Moscow

CMD [ "./app/bin/auth" ]