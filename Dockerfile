# сборка проекта: из исходников мы должны получить бинарный файл для запуска сервера
FROM golang:1.23.5 as build-service
# build-service - это название стадии при сборке образа из Dockerfile

WORKDIR /
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o UMKA_SERVER ./cmd/main.go

# Запуск собранной программы 
FROM alpine:latest
COPY --from=build-service /UMKA_SERVER .

CMD /UMKA_SERVER