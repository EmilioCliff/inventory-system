# Build stage
FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN pwd
RUN ls -la
RUN apk add --no-cache git-crypt
ARG DECRYPTION_KEY
RUN echo "$DECRYPTION_KEY" > /tmp/decryption.key && git-crypt unlock /tmp/decryption.key && rm /tmp/decryption.key
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY db/sqlc/Kokomed-Logo.png .
COPY start.sh .
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]