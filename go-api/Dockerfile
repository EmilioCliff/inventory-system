# Build stage
FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY db/sqlc/logi.png .
COPY start.sh .
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]