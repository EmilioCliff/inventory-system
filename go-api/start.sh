#!/bin/sh

set -e

# echo "running migrations"
# /app/migrate -path /app/migration -database ${{ secrets.DB_SOURCE }}  -verbose up

echo "run go app"
exec "$@"