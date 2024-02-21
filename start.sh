#!/bin/sh

set -e

echo "running migrations"
/app/migrate -path /app/migration -database postgresql://root:secret@postgres:5432/inventorydb?sslmode=disable -verbose up

echo "run go app"
exec "$@"