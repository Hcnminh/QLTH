#!/usr/bin/env bash
set -e

PGHOST="${PGHOST:-localhost}"
PGPORT="${PGPORT:-5432}"
PGUSER="${PGUSER:-qlth_user}"
PGDATABASE="${PGDATABASE:-qlth_eduprimary}"
export PGPASSWORD="${PGPASSWORD:-qlth_password_secure}"

echo "=== VALIDATING VAL-ADM-001 & VAL-DOGFOOD-001: 16-TABLE POSTGRESQL SCHEMA ==="

# 1. Check all 11 GoAdmin system tables
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

echo "1. Checking 11 GoAdmin system tables..."
for tbl in "${GOADMIN_TABLES[@]}"; do
    EXISTS=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -tAc "SELECT to_regclass('public.$tbl');")
    if [ "$EXISTS" = "$tbl" ]; then
        echo "  [OK] GoAdmin table: $tbl"
    else
        echo "  [FAIL] Missing GoAdmin table: $tbl"
        exit 1
    fi
done

# 2. Check complete 16 Business Tables from DAC_TA_THIET_KE_CHI_TIET_HE_THONG.md
BUSINESS_16_TABLES=(
    "schools"
    "academic_years"
    "terms"
    "classes"
    "staffs"
    "students"
    "subjects"
    "student_subject_evaluations"
    "student_behavior_evaluations"
    "academic_transcripts_snapshots"
    "boarding_menus"
    "boarding_daily_attendances"
    "student_health_records"
    "kpi_criteria"
    "kpi_cycles"
    "kpi_evaluation_results"
)

echo "2. Checking existence and row counts (>0) for all 16 Business DDL tables..."
for tbl in "${BUSINESS_16_TABLES[@]}"; do
    EXISTS=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -tAc "SELECT to_regclass('public.$tbl');")
    if [ "$EXISTS" != "$tbl" ]; then
        echo "  [FAIL] Missing 16-table business table: $tbl"
        exit 1
    fi

    COUNT=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -tAc "SELECT count(*) FROM $tbl;")
    if [ "$COUNT" -gt 0 ]; then
        echo "  [OK] Table $tbl exists and has $COUNT rows"
    else
        echo "  [FAIL] Table $tbl has 0 rows (missing seed data)"
        exit 1
    fi
done

# 3. Check GIN Indexes
GIN_INDEXES=(
    "idx_staffs_concurrent_roles"
    "idx_transcript_snapshot"
    "idx_boarding_menus_dishes"
    "idx_student_health_events"
    "idx_kpi_eval_metrics"
)

echo "3. Checking GIN indexes on JSONB fields..."
for idx in "${GIN_INDEXES[@]}"; do
    EXISTS=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -tAc "SELECT count(*) FROM pg_indexes WHERE indexname = '$idx' AND indexdef LIKE '%USING gin%';")
    if [ "$EXISTS" -ge 1 ]; then
        echo "  [OK] GIN index: $idx"
    else
        echo "  [FAIL] Missing GIN index: $idx"
        exit 1
    fi
done

# 4. Check Key Foreign Key constraints
echo "4. Checking Foreign Key constraints..."
FK_PAIRS=(
    "classes:schools"
    "classes:academic_years"
    "classes:staffs"
    "students:schools"
    "students:classes"
    "student_subject_evaluations:students"
    "student_subject_evaluations:terms"
    "student_subject_evaluations:subjects"
    "student_behavior_evaluations:students"
    "student_behavior_evaluations:terms"
    "academic_transcripts_snapshots:students"
    "boarding_menus:schools"
    "boarding_daily_attendances:students"
    "student_health_records:students"
    "kpi_criteria:schools"
    "kpi_cycles:schools"
    "kpi_evaluation_results:kpi_cycles"
    "kpi_evaluation_results:staffs"
)

for pair in "${FK_PAIRS[@]}"; do
    SRC="${pair%%:*}"
    DST="${pair##*:}"
    FK_EXISTS=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -tAc "
        SELECT count(*) FROM pg_constraint 
        WHERE contype = 'f' 
          AND conrelid = 'public.$SRC'::regclass 
          AND confrelid = 'public.$DST'::regclass;
    ")
    if [ "$FK_EXISTS" -ge 1 ]; then
        echo "  [OK] Foreign Key: $SRC -> $DST"
    else
        echo "  [FAIL] Missing Foreign Key: $SRC -> $DST"
        exit 1
    fi
done

# 5. Check default admin user
echo "5. Checking default admin user in goadmin_users..."
ADMIN_COUNT=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -tAc "SELECT count(*) FROM goadmin_users WHERE username = 'admin';")
if [ "$ADMIN_COUNT" -ge 1 ]; then
    echo "  [OK] Default admin user present in goadmin_users (count: $ADMIN_COUNT)"
else
    echo "  [FAIL] Default admin user missing in goadmin_users"
    exit 1
fi

echo "=== VAL-ADM-001 & VAL-DOGFOOD-001: ALL 16 TABLES AND ASSERTIONS PASSED ==="
