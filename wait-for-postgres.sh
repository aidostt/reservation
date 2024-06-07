#!/bin/sh

# wait-for-postgres.sh

set -e

host="${POSTGRES_HOST}"
port="${POSTGRES_PORT}"

until pg_isready -h "$host" -p "$port"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - running migrations"
migrate -path /app/migrations -database postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable up

>&2 echo "Migrations completed - starting application"
exec ./api

