#!/usr/bin/env bash
set -e

PGHOST="${PGHOST:-localhost}"
PGPORT="${PGPORT:-5432}"
PGUSER="${PGUSER:-qlth_user}"
PGDATABASE="${PGDATABASE:-qlth_eduprimary}"
export PGPASSWORD="${PGPASSWORD:-qlth_password_secure}"

echo "=== VALIDATING VAL-ADM-001: GoAdmin & Business Database Schema ==="

# 1. Check all 11 GoAdmin tables
GOADMIN_TABLES=(
    "goadmin_users"
    "goadmin_roles"
    "goadmin_permissions"
    "goadmin_menu"
    "goadmin_operation_log"
    "goadmin_site"
    "goadmin_session"
    "goadmin_role_menu"
    "goadmin_role_permissions"
    "goadmin_role_users"
    "goadmin_user_permissions"
)

echo "Checking 11 GoAdmin system tables..."
for tbl in "${GOADMIN_TABLES[@]}"; do
    EXISTS=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -tAc "SELECT to_regclass('public.$tbl');")
    if [ "$EXISTS" = "$tbl" ]; then
        echo "  [OK] GoAdmin table: $tbl"
    else
        echo "  [FAIL] Missing GoAdmin table: $tbl"
        exit 1
    fi
done

# 2. Check 5 Business DDL tables
BUSINESS_TABLES=(
    "schools"
    "staffs"
    "students"
    "kpi_evaluations"
    "teaching_quotas"
)

echo "Checking 5 Business DDL tables..."
for tbl in "${BUSINESS_TABLES[@]}"; do
    EXISTS=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -tAc "SELECT to_regclass('public.$tbl');")
    if [ "$EXISTS" = "$tbl" ]; then
        echo "  [OK] Business table: $tbl"
    else
        echo "  [FAIL] Missing Business table: $tbl"
        exit 1
    fi
done

# 3. Check default admin user
echo "Checking default admin user in goadmin_users..."
ADMIN_COUNT=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -tAc "SELECT count(*) FROM goadmin_users WHERE username = 'admin';")
if [ "$ADMIN_COUNT" -ge 1 ]; then
    echo "  [OK] Default admin user present in goadmin_users (count: $ADMIN_COUNT)"
else
    echo "  [FAIL] Default admin user missing in goadmin_users"
    exit 1
fi

echo "=== VAL-ADM-001: ALL ASSERTIONS PASSED ==="
