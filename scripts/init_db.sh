#!/usr/bin/env bash
set -e

PGHOST="${PGHOST:-localhost}"
PGPORT="${PGPORT:-5432}"
PGUSER="${PGUSER:-qlth_user}"
PGDATABASE="${PGDATABASE:-qlth_eduprimary}"
export PGPASSWORD="${PGPASSWORD:-qlth_password_secure}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "Initializing PostgreSQL database: $PGDATABASE on $PGHOST:$PGPORT as $PGUSER..."
psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -f "$SCRIPT_DIR/init_db.sql"
echo "Database initialization complete."
