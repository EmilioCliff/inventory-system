# Build stage
FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /app
COPY . .

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk update

RUN apk add --no-cache git-secret gnupg git
RUN git init
ARG GPG_PRIVATE_KEY
RUN echo "$GPG_PRIVATE_KEY" > /app/private_key
RUN gpg --import /app/private_key
RUN git secret reveal -f

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN cd /app/go-api
RUN go build -o main main.go
RUN cd ..

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/go-api/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/go-api/app.env .
COPY go-api/db/migration ./migration
COPY go-api/db/sqlc/Kokomed-Logo.png .
COPY go-api/start.sh .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]