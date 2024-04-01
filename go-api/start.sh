#!/bin/sh

set -e

. /app/.env

echo "running migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE"   -verbose up

echo "run go app"
exec "$@"