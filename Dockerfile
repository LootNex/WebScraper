# FROM golang:alpine AS builder

# WORKDIR /price_tracker

# RUN apk add --no-cache git

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o price_tracker ./cmd

# FROM alpine:latest

# WORKDIR /root/

# RUN apk add --no-cache ca-certificates

# COPY --from=builder /price_tracker .

# EXPOSE 50051

# CMD ["./price_tracker"]

# FROM golang:alpine AS builder

# WORKDIR /price_tracker

# RUN apk add --no-cache git

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o price_tracker ./cmd

# FROM alpine:latest

# WORKDIR /root/

# RUN apk add --no-cache ca-certificates

# COPY --from=builder /price_tracker .

# EXPOSE 50051

# CMD ["./price_tracker"]

FROM golang:1.24-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o price_tracker ./cmd

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
  wget \
  unzip \
  gnupg \
  curl \
  ca-certificates \
  fonts-liberation \
  libappindicator3-1 \
  libasound2 \
  libatk-bridge2.0-0 \
  libatk1.0-0 \
  libcups2 \
  libdbus-1-3 \
  libgdk-pixbuf2.0-0 \
  libnspr4 \
  libnss3 \
  libx11-xcb1 \
  libxcomposite1 \
  libxdamage1 \
  libxrandr2 \
  libgbm1 \
  xdg-utils \
  libu2f-udev \
  && rm -rf /var/lib/apt/lists/*


RUN wget https://storage.googleapis.com/chrome-for-testing-public/135.0.7049.95/linux64/chrome-linux64.zip && \
    unzip chrome-linux64.zip && \
    mv chrome-linux64 /opt/chrome && \
    ln -s /opt/chrome/chrome /usr/bin/google-chrome && \
    rm chrome-linux64.zip


RUN wget https://storage.googleapis.com/chrome-for-testing-public/135.0.7049.95/linux64/chromedriver-linux64.zip && \
    unzip chromedriver-linux64.zip && \
    mv chromedriver-linux64/chromedriver /usr/local/bin/chromedriver && \
    chmod +x /usr/local/bin/chromedriver && \
    rm -rf chromedriver-linux64* 


RUN google-chrome --version && chromedriver --version

ENV CHROME_BIN=/usr/bin/google-chrome \
    CHROMEDRIVER_PATH=/usr/local/bin/chromedriver


WORKDIR /root/
COPY --from=builder /app/price_tracker .

EXPOSE 50051
CMD ["./price_tracker"]