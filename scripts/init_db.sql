-- =============================================================================
-- EduPrimary-KPI Database Initialization Script
-- Database: qlth_eduprimary
-- Standard: PostgreSQL 14+
-- Includes: 11 GoAdmin system tables + 5 Core business tables + seed data
-- =============================================================================

-- Ensure pgcrypto extension is available for UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =============================================================================
-- PART 1: GOADMIN SYSTEM TABLES (11 Tables)
-- =============================================================================

-- 1. goadmin_users
CREATE TABLE IF NOT EXISTS goadmin_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    avatar VARCHAR(255) DEFAULT '',
    remember_token VARCHAR(100),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 2. goadmin_roles
CREATE TABLE IF NOT EXISTS goadmin_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 3. goadmin_permissions
CREATE TABLE IF NOT EXISTS goadmin_permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) NOT NULL,
    http_method VARCHAR(255),
    http_path TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 4. goadmin_menu
CREATE TABLE IF NOT EXISTS goadmin_menu (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER DEFAULT 0 NOT NULL,
    type INTEGER DEFAULT 0,
    "order" INTEGER DEFAULT 0 NOT NULL,
    title VARCHAR(50) NOT NULL,
    header VARCHAR(100),
    plugin_name VARCHAR(100) NOT NULL DEFAULT '',
    icon VARCHAR(50) NOT NULL,
    uri VARCHAR(3000) NOT NULL,
    uuid VARCHAR(100),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 5. goadmin_operation_log
CREATE TABLE IF NOT EXISTS goadmin_operation_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    path VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    ip VARCHAR(15) NOT NULL,
    input TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 6. goadmin_site
CREATE TABLE IF NOT EXISTS goadmin_site (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    type INTEGER DEFAULT 0,
    description VARCHAR(3000),
    state INTEGER DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 7. goadmin_session
CREATE TABLE IF NOT EXISTS goadmin_session (
    id SERIAL PRIMARY KEY,
    sid VARCHAR(50) NOT NULL,
    "values" VARCHAR(3000) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 8. goadmin_role_menu
CREATE TABLE IF NOT EXISTS goadmin_role_menu (
    role_id INTEGER NOT NULL,
    menu_id INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 9. goadmin_role_permissions
CREATE TABLE IF NOT EXISTS goadmin_role_permissions (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 10. goadmin_role_users
CREATE TABLE IF NOT EXISTS goadmin_role_users (
    role_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- 11. goadmin_user_permissions
CREATE TABLE IF NOT EXISTS goadmin_user_permissions (
    user_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- =============================================================================
-- PART 2: GOADMIN SEED DATA (Users, Roles, Permissions, Menus)
-- Password for admin is 'admin' (bcrypt: $2a$10$AZKkf6GNUI9KSHy/SW5Mh..hUJIwjE597lUUi7EeOCoXQcgJbmavK)
-- =============================================================================

INSERT INTO goadmin_users (id, username, password, name, avatar, remember_token)
VALUES 
    (1, 'admin', '$2a$10$AZKkf6GNUI9KSHy/SW5Mh..hUJIwjE597lUUi7EeOCoXQcgJbmavK', 'Administrator', '', 'admin_token'),
    (2, 'operator', '$2a$10$AZKkf6GNUI9KSHy/SW5Mh..hUJIwjE597lUUi7EeOCoXQcgJbmavK', 'Operator', '', NULL)
ON CONFLICT (id) DO UPDATE SET 
    password = EXCLUDED.password,
    name = EXCLUDED.name;

SELECT setval('goadmin_users_id_seq', (SELECT MAX(id) FROM goadmin_users));

INSERT INTO goadmin_roles (id, name, slug)
VALUES 
    (1, 'Administrator', 'administrator'),
    (2, 'Operator', 'operator')
ON CONFLICT (id) DO NOTHING;

SELECT setval('goadmin_roles_id_seq', (SELECT MAX(id) FROM goadmin_roles));

INSERT INTO goadmin_permissions (id, name, slug, http_method, http_path)
VALUES 
    (1, 'All permission', '*', '', '*'),
    (2, 'Dashboard', 'dashboard', 'GET,PUT,POST,DELETE', '/')
ON CONFLICT (id) DO NOTHING;

SELECT setval('goadmin_permissions_id_seq', (SELECT MAX(id) FROM goadmin_permissions));

INSERT INTO goadmin_role_permissions (role_id, permission_id)
SELECT 1, 1 WHERE NOT EXISTS (SELECT 1 FROM goadmin_role_permissions WHERE role_id = 1 AND permission_id = 1);

INSERT INTO goadmin_role_permissions (role_id, permission_id)
SELECT 1, 2 WHERE NOT EXISTS (SELECT 1 FROM goadmin_role_permissions WHERE role_id = 1 AND permission_id = 2);

INSERT INTO goadmin_role_permissions (role_id, permission_id)
SELECT 2, 2 WHERE NOT EXISTS (SELECT 1 FROM goadmin_role_permissions WHERE role_id = 2 AND permission_id = 2);

INSERT INTO goadmin_role_users (role_id, user_id)
SELECT 1, 1 WHERE NOT EXISTS (SELECT 1 FROM goadmin_role_users WHERE role_id = 1 AND user_id = 1);

INSERT INTO goadmin_role_users (role_id, user_id)
SELECT 2, 2 WHERE NOT EXISTS (SELECT 1 FROM goadmin_role_users WHERE role_id = 2 AND user_id = 2);

INSERT INTO goadmin_user_permissions (user_id, permission_id)
SELECT 1, 1 WHERE NOT EXISTS (SELECT 1 FROM goadmin_user_permissions WHERE user_id = 1 AND permission_id = 1);

INSERT INTO goadmin_user_permissions (user_id, permission_id)
SELECT 2, 2 WHERE NOT EXISTS (SELECT 1 FROM goadmin_user_permissions WHERE user_id = 2 AND permission_id = 2);

INSERT INTO goadmin_menu (id, parent_id, type, "order", title, header, plugin_name, icon, uri)
VALUES 
    (1, 0, 1, 1, 'Dashboard', NULL, '', 'fa-bar-chart', '/'),
    (2, 0, 1, 2, 'Admin', NULL, '', 'fa-tasks', ''),
    (3, 2, 1, 1, 'Users', NULL, '', 'fa-users', '/info/manager'),
    (4, 2, 1, 2, 'Roles', NULL, '', 'fa-user', '/info/roles'),
    (5, 2, 1, 3, 'Permission', NULL, '', 'fa-ban', '/info/permission'),
    (6, 2, 1, 4, 'Menu', NULL, '', 'fa-bars', '/menu'),
    (7, 2, 1, 5, 'Operation log', NULL, '', 'fa-history', '/info/op'),
    (8, 0, 1, 3, 'Giáo viên & Định mức', 'Nghiệp vụ Giáo dục', '', 'fa-briefcase', '/info/staffs'),
    (9, 0, 1, 4, 'Học sinh & Consent NĐ13', NULL, '', 'fa-graduation-cap', '/info/students'),
    (10, 0, 1, 5, 'Đánh giá KPI 360', NULL, '', 'fa-star', '/info/kpi_evaluations')
ON CONFLICT (id) DO UPDATE SET 
    title = EXCLUDED.title,
    uri = EXCLUDED.uri,
    header = EXCLUDED.header;

SELECT setval('goadmin_menu_id_seq', (SELECT MAX(id) FROM goadmin_menu));

INSERT INTO goadmin_role_menu (role_id, menu_id)
SELECT 1, m.id FROM goadmin_menu m
WHERE NOT EXISTS (SELECT 1 FROM goadmin_role_menu rm WHERE rm.role_id = 1 AND rm.menu_id = m.id);

INSERT INTO goadmin_role_menu (role_id, menu_id)
SELECT 2, m.id FROM goadmin_menu m WHERE m.id IN (1, 8, 9, 10)
AND NOT EXISTS (SELECT 1 FROM goadmin_role_menu rm WHERE rm.role_id = 2 AND rm.menu_id = m.id);

-- =============================================================================
-- PART 3: BUSINESS TABLES (5 Tables)
-- =============================================================================

-- 1. schools
CREATE TABLE IF NOT EXISTS schools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(30) NOT NULL DEFAULT 'PRIMARY',
    address TEXT,
    province_code VARCHAR(20) NOT NULL,
    district_code VARCHAR(20) NOT NULL,
    ward_code VARCHAR(20),
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 2. staffs (Circular 20/2023 & Circular 05/2025)
CREATE TABLE IF NOT EXISTS staffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    national_id VARCHAR(12) UNIQUE NOT NULL,
    employee_code VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('MALE', 'FEMALE', 'OTHER')),
    birth_date DATE NOT NULL,
    phone_number VARCHAR(20),
    email VARCHAR(255),
    job_position VARCHAR(50) NOT NULL,
    teacher_rank VARCHAR(20),
    is_homeroom BOOLEAN DEFAULT FALSE,
    standard_annual_hours INT DEFAULT 805,
    weekly_teaching_hours INT DEFAULT 23,
    weekly_reduced_hours INT DEFAULT 0,
    concurrent_roles JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 3. students (Decree 13/2023 Child Privacy Protection)
CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    national_id VARCHAR(12) UNIQUE,
    student_code VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('MALE', 'FEMALE')),
    ethnicity VARCHAR(50) DEFAULT 'Kinh',
    current_class_id UUID,
    status VARCHAR(30) DEFAULT 'STUDYING',
    is_inclusive_education BOOLEAN DEFAULT FALSE,
    guardian_name VARCHAR(255) NOT NULL,
    guardian_national_id VARCHAR(12) NOT NULL,
    guardian_phone VARCHAR(20),
    guardian_relationship VARCHAR(50) DEFAULT 'PARENT',
    has_digital_consent BOOLEAN DEFAULT FALSE,
    digital_consent_verified BOOLEAN DEFAULT FALSE,
    consent_signed_at TIMESTAMPTZ,
    consent_method VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 4. kpi_evaluations (Decree 90/2020 & Decree 48/2023 & 20% cap)
CREATE TABLE IF NOT EXISTS kpi_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID REFERENCES staffs(id) ON DELETE CASCADE,
    staff_name VARCHAR(255) NOT NULL,
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    cycle_name VARCHAR(100) DEFAULT '2026-2027',
    automated_system_score NUMERIC(5,2) DEFAULT 0,
    manager_assessment_score NUMERIC(5,2) DEFAULT 0,
    peer_review_score NUMERIC(5,2) DEFAULT 0,
    parent_satisfaction_score NUMERIC(5,2) DEFAULT 0,
    final_score NUMERIC(5,2) NOT NULL DEFAULT 0,
    official_rank VARCHAR(100) NOT NULL,
    has_disciplinary_action BOOLEAN DEFAULT FALSE,
    violated_tutoring_rules BOOLEAN DEFAULT FALSE,
    has_severe_safety_incident BOOLEAN DEFAULT FALSE,
    has_innovative_solution BOOLEAN DEFAULT FALSE,
    is_capped_by_20pct BOOLEAN DEFAULT FALSE,
    reason TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 5. teaching_quotas (Circular 05/2025/TT-BGDĐT)
CREATE TABLE IF NOT EXISTS teaching_quotas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID REFERENCES staffs(id) ON DELETE CASCADE,
    teacher_name VARCHAR(255) NOT NULL,
    school_type VARCHAR(50) DEFAULT 'PRIMARY',
    is_homeroom BOOLEAN DEFAULT FALSE,
    is_head_subject BOOLEAN DEFAULT FALSE,
    is_deputy_head_subject BOOLEAN DEFAULT FALSE,
    has_young_child BOOLEAN DEFAULT FALSE,
    counseling_hours_week INT DEFAULT 0,
    inclusive_teaching_wk INT DEFAULT 0,
    pre_grade1_ethnic_hours INT DEFAULT 0,
    weeks_in_school_year INT DEFAULT 35,
    standard_annual_hours INT DEFAULT 805,
    weekly_reduced_hours INT DEFAULT 0,
    annual_reduced_hours INT DEFAULT 0,
    net_annual_quota INT DEFAULT 805,
    actual_taught_hours INT DEFAULT 0,
    inclusive_bonus_hours NUMERIC(6,2) DEFAULT 0,
    total_effective_hours NUMERIC(6,2) DEFAULT 0,
    overtime_hours NUMERIC(6,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- PART 4: BUSINESS SEED DATA
-- =============================================================================

-- Seed default school
INSERT INTO schools (id, code, name, type, address, province_code, district_code, ward_code, is_public)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'TH_LE_VAN_TAM_01',
    'Trường Tiểu học Lê Văn Tám',
    'PRIMARY',
    'Số 123 Đường Hai Bà Trưng, Quận 1, TP. Hồ Chí Minh',
    '79',
    '760',
    '26734',
    TRUE
) ON CONFLICT (code) DO NOTHING;

-- Seed staffs
INSERT INTO staffs (
    id, school_id, national_id, employee_code, full_name, gender, birth_date,
    phone_number, email, job_position, teacher_rank, is_homeroom,
    standard_annual_hours, weekly_teaching_hours, weekly_reduced_hours, concurrent_roles
) VALUES 
(
    'b0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    '079185000001',
    'GV001',
    'Nguyễn Văn An',
    'MALE',
    '1985-05-15',
    '0901234567',
    'an.nv@eduprimary.vn',
    'PRIMARY_TEACHER',
    'RANK_II',
    TRUE,
    805,
    23,
    4,
    '["GVCN"]'::jsonb
),
(
    'b0000000-0000-0000-0000-000000000002',
    'a0000000-0000-0000-0000-000000000001',
    '079188000002',
    'GV002',
    'Trần Thị Bình',
    'FEMALE',
    '1988-10-20',
    '0902345678',
    'binh.tt@eduprimary.vn',
    'PRIMARY_TEACHER',
    'RANK_II',
    FALSE,
    805,
    23,
    3,
    '["HEAD_SUBJECT"]'::jsonb
),
(
    'b0000000-0000-0000-0000-000000000003',
    'a0000000-0000-0000-0000-000000000001',
    '079180000003',
    'GV003',
    'Lê Hoàng Cường',
    'MALE',
    '1980-03-12',
    '0903456789',
    'cuong.lh@eduprimary.vn',
    'PRINCIPAL',
    'RANK_I',
    FALSE,
    805,
    4,
    19,
    '["PRINCIPAL"]'::jsonb
) ON CONFLICT (employee_code) DO NOTHING;

-- Seed students
INSERT INTO students (
    id, school_id, national_id, student_code, full_name, birth_date, gender, ethnicity,
    status, is_inclusive_education, guardian_name, guardian_national_id, guardian_phone,
    guardian_relationship, has_digital_consent, digital_consent_verified, consent_signed_at, consent_method
) VALUES 
(
    'c0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    '079219000001',
    'HS2026001',
    'Nguyễn Hoàng Nam',
    '2019-03-15',
    'MALE',
    'Kinh',
    'STUDYING',
    FALSE,
    'Nguyễn Văn An',
    '079185000001',
    '0901234567',
    'FATHER',
    TRUE,
    TRUE,
    NOW(),
    'VNEID_SSO'
),
(
    'c0000000-0000-0000-0000-000000000002',
    'a0000000-0000-0000-0000-000000000001',
    '079219000002',
    'HS2026002',
    'Trần Mai Linh',
    '2019-07-22',
    'FEMALE',
    'Kinh',
    'STUDYING',
    FALSE,
    'Trần Văn Dũng',
    '079182000005',
    '0905678901',
    'FATHER',
    TRUE,
    TRUE,
    NOW(),
    'PORTAL_SIGN'
),
(
    'c0000000-0000-0000-0000-000000000003',
    'a0000000-0000-0000-0000-000000000001',
    '079219000003',
    'HS2026003',
    'Lê Quốc Bảo',
    '2019-11-05',
    'MALE',
    'Kinh',
    'STUDYING',
    FALSE,
    'Lê Hoàng Cường',
    '079180000003',
    '0903456789',
    'FATHER',
    FALSE,
    FALSE,
    NULL,
    NULL
) ON CONFLICT (student_code) DO NOTHING;

-- Seed teaching_quotas
INSERT INTO teaching_quotas (
    id, staff_id, teacher_name, school_type, is_homeroom, is_head_subject,
    weeks_in_school_year, standard_annual_hours, weekly_reduced_hours,
    annual_reduced_hours, net_annual_quota, actual_taught_hours, overtime_hours
) VALUES 
(
    'd0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000001',
    'Nguyễn Văn An',
    'PRIMARY',
    TRUE,
    FALSE,
    35,
    805,
    4,
    140,
    665,
    720,
    55.00
),
(
    'd0000000-0000-0000-0000-000000000002',
    'b0000000-0000-0000-0000-000000000002',
    'Trần Thị Bình',
    'PRIMARY',
    FALSE,
    TRUE,
    35,
    805,
    3,
    105,
    700,
    700,
    0.00
) ON CONFLICT (id) DO NOTHING;

-- Seed kpi_evaluations
INSERT INTO kpi_evaluations (
    id, staff_id, staff_name, school_id, cycle_name,
    automated_system_score, manager_assessment_score, peer_review_score, parent_satisfaction_score,
    final_score, official_rank, has_innovative_solution, is_capped_by_20pct, reason
) VALUES 
(
    'e0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000001',
    'Nguyễn Văn An',
    'a0000000-0000-0000-0000-000000000001',
    'Năm học 2026-2027',
    95.00,
    92.00,
    90.00,
    94.00,
    93.20,
    'HOÀN THÀNH XUẤT SẮC NHIỆM VỤ',
    TRUE,
    FALSE,
    'Đạt 93.20 điểm và có sáng kiến kinh nghiệm cấp trường'
),
(
    'e0000000-0000-0000-0000-000000000002',
    'b0000000-0000-0000-0000-000000000002',
    'Trần Thị Bình',
    'a0000000-0000-0000-0000-000000000001',
    'Năm học 2026-2027',
    85.00,
    88.00,
    84.00,
    86.00,
    85.90,
    'HOÀN THÀNH TỐT NHIỆM VỤ',
    FALSE,
    FALSE,
    'Hoàn thành tốt nhiệm vụ giảng dạy và công tác tổ trưởng'
),
(
    'e0000000-0000-0000-0000-000000000003',
    'b0000000-0000-0000-0000-000000000003',
    'Lê Hoàng Cường',
    'a0000000-0000-0000-0000-000000000001',
    'Năm học 2026-2027',
    91.00,
    90.00,
    88.00,
    92.00,
    90.40,
    'HOÀN THÀNH TỐT NHIỆM VỤ',
    TRUE,
    TRUE,
    'Đạt 90.40 điểm có SKKN nhưng hạ xuống mức Tốt do khống chế trần 20% theo NĐ 90/2020'
) ON CONFLICT (id) DO NOTHING;
