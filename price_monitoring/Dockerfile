FROM golang:alpine AS builder

WORKDIR /price_tracker

RUN apk add --no-cache git

COPY price_monitoring/go.mod price_monitoring/go.sum ./
RUN go mod download

COPY price_monitoring/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -x -o price_tracker ./cmd

FROM alpine:latest

RUN apk add --no-cache python3 py3-pip ca-certificates

RUN apk add --no-cache py3-requests && \
    ln -sf /usr/bin/python3 /usr/bin/python

WORKDIR /root/

COPY --from=builder /price_tracker .

EXPOSE 50051

CMD ["./price_tracker"]