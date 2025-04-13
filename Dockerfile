# Этап сборки
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем ВЕСЬ проект (включая internal/, cmd/, и другие папки)
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/cryptoapp/main.go

# Этап запуска
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["/app/main"]
