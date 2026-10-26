#!/usr/bin/env bash
# ==============================================================================
# BN NEAR — Automated Database Migration & Seeding Script
# ==============================================================================
set -euo pipefail

DB_URL="${DATABASE_URL:-postgres://bn_near:bn_near_secret_password@localhost:5432/bn_near_db?sslmode=disable}"

echo "[BN-NEAR-DB] Starting database migration..."
echo "[BN-NEAR-DB] Target DB: ${DB_URL}"

# Apply initial schema migration
psql "${DB_URL}" -f "$(dirname "$0")/../migrations/001_initial_schema.sql"

echo "[BN-NEAR-DB] Migration completed successfully."

# Seed demo data if SEED_DATA=true
if [ "${SEED_DATA:-false}" = "true" ]; then
  echo "[BN-NEAR-DB] Seeding initial data..."
  psql "${DB_URL}" -f "$(dirname "$0")/../fixtures/seed_data.sql"
  echo "[BN-NEAR-DB] Seeding completed."
fi
