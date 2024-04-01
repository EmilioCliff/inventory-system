# Build stage
FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /app
COPY . /app
ARG DECRYPTION_KEY
RUN apk add --no-cache git-crypt
RUN git init
RUN echo "$DECRYPTION_KEY" | base64 -d > git-crypt-key
RUN git-crypt unlock git-crypt-key
RUN rm git-crypt-key

RUN mkdir /decrypted
RUN cp -r /app/app.env /decrypted

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /decrypted .
COPY db/migration ./migration
COPY db/sqlc/Kokomed-Logo.png .
COPY start.sh .


EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]