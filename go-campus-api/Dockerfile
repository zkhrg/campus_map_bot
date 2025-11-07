# --- Stage 1: Builder ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN echo "https://mirror.yandex.ru/mirrors/alpine/v3.22/main" > /etc/apk/repositories \
 && echo "https://mirror.yandex.ru/mirrors/alpine/v3.22/community" >> /etc/apk/repositories \
 && apk update \
 && apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/main.go


# --- Stage 2: Packing ---
FROM alpine:3.22

RUN echo "https://mirror.yandex.ru/mirrors/alpine/v3.22/main" > /etc/apk/repositories \
 && echo "https://mirror.yandex.ru/mirrors/alpine/v3.22/community" >> /etc/apk/repositories \
 && apk update \
 && apk add --no-cache ca-certificates \
 && update-ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/app .

EXPOSE 8080

CMD ["./app"]
