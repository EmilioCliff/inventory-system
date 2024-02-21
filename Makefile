postgres:
	docker run --name postgres3 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -p 5432:5432 -d postgres:alpine3.19

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@postgres3:5432/inventorydb?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@postgres3:5432/inventorydb?sslmode=disable" -verbose down

createdb:
	docker exec -it postgres3 createdb --username=root --owner=root inventorydb

dropdb:
	docker exec -it postgres3 dropdb inventorydb

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migratedown migrateup sqlc test server