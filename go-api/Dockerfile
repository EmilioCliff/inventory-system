# Build stage
FROM golang:1.22.0-alpine3.19 AS builder
RUN apk add --no-cache git-crypt
WORKDIR /app
RUN git init
COPY . .
ARG DECRYPTION_KEY
RUN git config user.email "clifftest33@gmail.com"
RUN git config user.name "Inventory System"
RUN git add .
RUN git commit -m "Add app.env file"
RUN echo "$DECRYPTION_KEY" | base64 -d > git-crypt-key
RUN git-crypt unlock git-crypt-key
RUN rm git-crypt-key

# RUN cp /app/app.env /app/decrypted
# RUN cat /app/decrypted/app.env

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/app.env .
COPY db/migration ./migration
COPY db/sqlc/Kokomed-Logo.png .
COPY start.sh .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]