#!/bin/sh

set -e

echo "running migrations"
/app/migrate -path /app/migration -database postgresql://postgres:D*fd6EDcc5bdA6AG1c4FF6DafFeG33b*@viaduct.proxy.rlwy.net:26578/railway -verbose up

echo "run go app"
exec "$@"