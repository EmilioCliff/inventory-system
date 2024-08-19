#!/bin/sh

set -e

# echo "running migrate down"
# /app/migrate -path /app/migration -database "$DB_SOURCE" -verbose down 2

echo "running migrations"
/app/migrate -path /app/migration -database postgresql://root:secret@inventorydb:5432/inventorydb?sslmode=disable -verbose up

exec "$@"