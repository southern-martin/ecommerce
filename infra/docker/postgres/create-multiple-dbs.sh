#!/bin/bash
set -e

# ─────────────────────────────────────────────────────────────
# Create multiple databases for the ecommerce microservices.
# This script runs during PostgreSQL container initialization.
# ─────────────────────────────────────────────────────────────

DATABASES=(
    ecommerce_auth
    ecommerce_users
    ecommerce_products
    ecommerce_orders
    ecommerce_payments
    ecommerce_cart
    ecommerce_search
    ecommerce_reviews
    ecommerce_notifications
    ecommerce_chat
    ecommerce_media
    ecommerce_ai
    ecommerce_promotions
    ecommerce_returns
    ecommerce_shipping
    ecommerce_loyalty
    ecommerce_affiliates
    ecommerce_tax
    ecommerce_cms
)

# Databases that require the pgvector extension
VECTOR_DATABASES=(
    ecommerce_products
    ecommerce_ai
)

echo "=== Creating ecommerce databases ==="

for db in "${DATABASES[@]}"; do
    echo "Creating database: $db"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        SELECT 'CREATE DATABASE $db'
        WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$db')\gexec
EOSQL

    echo "Enabling uuid-ossp extension on: $db"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$db" <<-EOSQL
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
EOSQL
done

for db in "${VECTOR_DATABASES[@]}"; do
    echo "Enabling vector extension on: $db (if available)"
    psql --username "$POSTGRES_USER" --dbname "$db" <<-EOSQL || echo "WARNING: pgvector not available, skipping for $db"
        CREATE EXTENSION IF NOT EXISTS vector;
EOSQL
done

echo "=== All databases created successfully ==="
