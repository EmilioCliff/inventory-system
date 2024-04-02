# Build stage
FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /app
COPY . .

RUN apt-get update && apt-get install -y git-secret gnupg git
RUN git init
ARG GPG_PRIVATE_KEY
RUN echo "$GPG_PRIVATE_KEY" > /app/private_key
RUN gpg --import /app/private_key
RUN git secret reveal -f

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN go build -o main app/go-api/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/go-api/app.env .
COPY go-api/db/migration ./migration
COPY go-api/db/sqlc/Kokomed-Logo.png .
COPY go-api/start.sh .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]