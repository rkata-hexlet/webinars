FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o server main.go

FROM alpine:latest
WORKDIR /root/
# Добавляем установку curl
# --no-cache позволяет не сохранять индекс пакетов внутри образа, уменьшая его размер
RUN apk add --no-cache curl

COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]