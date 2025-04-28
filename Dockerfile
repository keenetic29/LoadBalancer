# Билд стадии
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o loadbalancer ./cmd/main.go

# Финальная стадия
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и конфиг
COPY --from=builder /app/loadbalancer .
COPY config.json .

# Создаем файл clients.json и устанавливаем права 
RUN touch /app/clients.json && chmod 666 /app/clients.json

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./loadbalancer"]