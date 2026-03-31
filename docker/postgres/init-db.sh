#!/bin/sh
set -eu

DB_HOST="${POSTGRES_HOST:-postgres}"
DB_PORT="${POSTGRES_PORT_INTERNAL:-5432}"

export PGPASSWORD="${POSTGRES_PASSWORD:-}"

psql_cmd() {
  psql \
    -v ON_ERROR_STOP=1 \
    -h "$DB_HOST" \
    -p "$DB_PORT" \
    -U "$POSTGRES_USER" \
    -d "$POSTGRES_DB" \
    "$@"
}

echo "Waiting for PostgreSQL at ${DB_HOST}:${DB_PORT}..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" >/dev/null 2>&1; do
  sleep 2
done

has_core_schema="$(psql_cmd -tAc "SELECT CASE WHEN to_regclass('public.\"User\"') IS NULL THEN '0' ELSE '1' END")"

if [ "$has_core_schema" = "0" ]; then
  echo "Applying base schema dump..."
  psql_cmd -f /dumps/01-schema.sql

  echo "Applying base data dump..."
  psql_cmd -f /dumps/data-only-dump.sql
else
  echo "Base schema already present, skipping schema/data dump import."
fi

echo "Applying incremental dump scripts..."
for sql_file in /dumps/[0-9][0-9]-*.sql; do
  [ -f "$sql_file" ] || continue
  case "$(basename "$sql_file")" in
    01-schema.sql|data-only-dump.sql)
      continue
      ;;
  esac
  echo "Running $(basename "$sql_file")..."
  psql_cmd -f "$sql_file"
done

echo "Database initialization finished."
