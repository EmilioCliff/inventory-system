postgres:
	docker run --name postgres3 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -p 5432:5432 -d postgres:alpine3.19

migrateup:
	migrate -path go-api/db/migration -database "$(DB_SOURCE_DEVELOPMENT)" -verbose up

migrateup1:
	migrate -path go-api/db/migration -database "$(DB_SOURCE_DEVELOPMENT)" -verbose up 1

migratedown:
	migrate -path go-api/db/migration -database "$(DB_SOURCE_DEVELOPMENT)" -verbose down

migratedown1:
	migrate -path go-api/db/migration -database "$(DB_SOURCE_DEVELOPMENT)" -verbose down 1

createdb:
	docker exec -it postgres3 createdb --username=root --owner=root inventorydb

dropdb:
	docker exec -it postgres3 dropdb inventorydb

sqlc:
	sqlc generate

test:
	cd go-api && go test -v -cover ./...

server:
	cd go-api && go run main.go

python-image:
	cd python-flask && docker build -t python-flask:latest .

python:
	docker run --name pythonapp -d -p 3000:3000 python-flask:latest 

go-image:
	cd go-api && docker build -t go-api:latest .

go:
	docker run --name goapi -d -p 8080:8080 go-api:latest 

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2.4-alpine3.19

connecRedis:
	asynq dash --uri="$(REDIS_URI)" --password="$(REDIS_PASSWORD)"

.PHONY: postgres createdb dropdb migratedown migrateup migratedown1 migrateup1 sqlc test server python-image python go-image go redis connecRedis