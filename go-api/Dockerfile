# Build stage
FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /app
COPY . /app
RUN pwd
RUN ls -la
ARG DECRYPTION_KEY
RUN apk add --no-cache git-crypt
RUN git init
RUN echo "$DECRYPTION_KEY" > /tmp/git-crypt-key
RUN git-crypt migrate-key /tmp/git-crypt-key /tmp/new-git-crypt-key
RUN cat /tmp/new-git-crypt-key
RUN git-crypt unlock /tmp/new-git-crypt-key
RUN rm -r /tmp/new-git-crypt-key
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/go-api/db/migration ./migration
COPY --from=builder /app/go-api/db/sqlc/Kokomed-Logo.png .
COPY --from=builder /app/go-api/start.sh .
COPY --from=builder /app/go-api/app.env .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]