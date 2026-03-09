#!/bin/bash
# Database Migration Runner for Ecommerce Microservices
# Iterates over all services, creates databases if needed, and applies .up.sql migrations in order.
# Tracks applied migrations in a schema_migrations table per database.

set -e

POSTGRES_HOST="${POSTGRES_HOST:-postgres}"
POSTGRES_PORT="${POSTGRES_PORT:-5432}"
POSTGRES_USER="${POSTGRES_USER:-ecommerce}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-ecommerce_password}"

export PGPASSWORD="$POSTGRES_PASSWORD"

SERVICES=(
  auth user product cart order payment search review
  notification chat media cms promotion loyalty affiliate
  shipping return tax ai
)

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

psql_cmd() {
  local db="$1"
  shift
  psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$db" -v ON_ERROR_STOP=1 "$@"
}

psql_admin() {
  psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d postgres -v ON_ERROR_STOP=1 "$@"
}

ensure_database() {
  local service="$1"
  local db_name="${service}_db"

  local exists
  exists=$(psql_admin -tAc "SELECT 1 FROM pg_database WHERE datname = '${db_name}'" 2>/dev/null || true)

  if [ "$exists" != "1" ]; then
    log "Creating database: ${db_name}"
    psql_admin -c "CREATE DATABASE ${db_name};"
  fi
}

ensure_schema_migrations_table() {
  local db_name="$1"
  psql_cmd "$db_name" -c "
    CREATE TABLE IF NOT EXISTS schema_migrations (
      id SERIAL PRIMARY KEY,
      filename VARCHAR(255) NOT NULL UNIQUE,
      applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );
  " > /dev/null 2>&1
}

is_migration_applied() {
  local db_name="$1"
  local filename="$2"
  local result
  result=$(psql_cmd "$db_name" -tAc "SELECT 1 FROM schema_migrations WHERE filename = '${filename}'" 2>/dev/null || true)
  [ "$result" = "1" ]
}

record_migration() {
  local db_name="$1"
  local filename="$2"
  psql_cmd "$db_name" -c "INSERT INTO schema_migrations (filename) VALUES ('${filename}');" > /dev/null 2>&1
}

run_service_migrations() {
  local service="$1"
  local db_name="${service}_db"
  local migrations_dir="/migrations/${service}"

  if [ ! -d "$migrations_dir" ]; then
    log "No migrations directory for ${service}, skipping"
    return 0
  fi

  local sql_files
  sql_files=$(find "$migrations_dir" -name "*.up.sql" -type f 2>/dev/null | sort)

  if [ -z "$sql_files" ]; then
    log "No .up.sql files for ${service}, skipping"
    return 0
  fi

  ensure_database "$service"
  ensure_schema_migrations_table "$db_name"

  local applied=0
  local skipped=0

  while IFS= read -r sql_file; do
    local filename
    filename=$(basename "$sql_file")

    if is_migration_applied "$db_name" "$filename"; then
      skipped=$((skipped + 1))
      continue
    fi

    log "Applying ${service}/${filename}..."
    if psql_cmd "$db_name" -f "$sql_file"; then
      record_migration "$db_name" "$filename"
      applied=$((applied + 1))
    else
      log "ERROR: Failed to apply ${service}/${filename}"
      return 1
    fi
  done <<< "$sql_files"

  log "${service}: ${applied} applied, ${skipped} already applied"
}

# Main
log "========================================="
log "Starting database migrations"
log "Host: ${POSTGRES_HOST}:${POSTGRES_PORT}"
log "========================================="

# Wait for PostgreSQL to be ready
log "Waiting for PostgreSQL to be ready..."
for i in $(seq 1 30); do
  if psql_admin -c "SELECT 1" > /dev/null 2>&1; then
    log "PostgreSQL is ready"
    break
  fi
  if [ "$i" -eq 30 ]; then
    log "ERROR: PostgreSQL not ready after 30 attempts"
    exit 1
  fi
  sleep 2
done

FAILED=0
for service in "${SERVICES[@]}"; do
  log "-----------------------------------------"
  log "Processing service: ${service}"
  if ! run_service_migrations "$service"; then
    log "ERROR: Migrations failed for ${service}"
    FAILED=$((FAILED + 1))
  fi
done

log "========================================="
if [ "$FAILED" -gt 0 ]; then
  log "ERROR: ${FAILED} service(s) had migration failures"
  exit 1
fi

log "All migrations completed successfully"
exit 0
