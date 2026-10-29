-- =============================================================================
-- EduPrimary-KPI Database Initialization Script
-- Database: qlth_eduprimary
-- Standard: PostgreSQL 14+
-- Includes: 11 GoAdmin system tables + 16 Core business tables + seed data
-- Reference: DAC_TA_THIET_KE_CHI_TIET_HE_THONG.md
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
-- PART 3: 16 CORE BUSINESS TABLES (FROM DAC_TA_THIET_KE_CHI_TIET_HE_THONG.MD)
-- =============================================================================

-- 1. schools (Tenant)
CREATE TABLE IF NOT EXISTS schools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,             -- Mã trường định danh Bộ GD&ĐT (VD: 'TH_LE_VAN_TAM_01')
    name VARCHAR(255) NOT NULL,
    type VARCHAR(30) NOT NULL DEFAULT 'PRIMARY',  -- 'PRIMARY', 'SECONDARY', 'KINDERGARTEN', 'INTER_LEVEL'
    address TEXT,
    province_code VARCHAR(20) NOT NULL,
    district_code VARCHAR(20) NOT NULL,
    ward_code VARCHAR(20),
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 2. academic_years (Niên khóa / Năm học)
CREATE TABLE IF NOT EXISTS academic_years (
    id VARCHAR(32) PRIMARY KEY,                  -- Ví dụ: '2026-2027'
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 3. terms (Học kỳ)
CREATE TABLE IF NOT EXISTS terms (
    id VARCHAR(32) PRIMARY KEY,                  -- Ví dụ: '2026-2027-HK1'
    academic_year_id VARCHAR(32) REFERENCES academic_years(id) ON DELETE CASCADE,
    term_number INT NOT NULL,                    -- 1: Giữa kỳ 1, 2: Cuối kỳ 1, 3: Giữa kỳ 2, 4: Cuối năm
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_closed BOOLEAN DEFAULT FALSE
);

-- 4. staffs (Cán bộ, Giáo viên & Nhân viên - Theo TT 20/2023 & TT 05/2025)
CREATE TABLE IF NOT EXISTS staffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    national_id VARCHAR(12) UNIQUE NOT NULL,     -- Số CCCD / Định danh VNeID
    employee_code VARCHAR(50) UNIQUE NOT NULL,   -- Mã số cán bộ (GV001...)
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
CREATE INDEX IF NOT EXISTS idx_staffs_concurrent_roles ON staffs USING GIN (concurrent_roles);

-- 5. classes (Lớp học)
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    academic_year_id VARCHAR(32) REFERENCES academic_years(id) ON DELETE CASCADE,
    grade_level INT NOT NULL CHECK (grade_level BETWEEN 1 AND 5), -- Lớp 1 đến Lớp 5
    name VARCHAR(100) NOT NULL,                  -- Ví dụ: '1A1', '5A2'
    homeroom_teacher_id UUID REFERENCES staffs(id) ON DELETE SET NULL, -- FK tới staffs(id)
    room_number VARCHAR(50),
    student_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 6. students (Học sinh - Decree 13/2023 Child Privacy Protection)
CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    national_id VARCHAR(12) UNIQUE,              -- Số định danh cá nhân 12 số (Luật Căn cước 2023)
    student_code VARCHAR(50) UNIQUE NOT NULL,    -- Mã định danh Bộ GD&ĐT
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('MALE', 'FEMALE')),
    ethnicity VARCHAR(50) DEFAULT 'Kinh',
    current_class_id UUID REFERENCES classes(id) ON DELETE SET NULL,
    status VARCHAR(30) DEFAULT 'STUDYING',       -- 'STUDYING', 'TRANSFERRED', 'DROPPED', 'GRADUATED'
    is_inclusive_education BOOLEAN DEFAULT FALSE,-- Học sinh khuyết tật học hòa nhập (TT 03/2018)
    guardian_name VARCHAR(255) NOT NULL,
    guardian_national_id VARCHAR(12) NOT NULL,
    guardian_phone VARCHAR(20),
    guardian_relationship VARCHAR(50) DEFAULT 'PARENT',
    has_digital_consent BOOLEAN DEFAULT FALSE,
    digital_consent_verified BOOLEAN DEFAULT FALSE,
    consent_signed_at TIMESTAMPTZ,
    consent_method VARCHAR(50),                  -- 'VNEID_SSO', 'SMS_OTP', 'PORTAL_SIGN'
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key constraint to students.current_class_id if table was created previously without it
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_students_current_class'
    ) THEN
        ALTER TABLE students 
        ADD CONSTRAINT fk_students_current_class 
        FOREIGN KEY (current_class_id) REFERENCES classes(id) ON DELETE SET NULL;
    END IF;
END $$;

-- 7. subjects (Danh mục Môn học & Hoạt động giáo dục)
CREATE TABLE IF NOT EXISTS subjects (
    id VARCHAR(50) PRIMARY KEY,                  -- 'TIENG_VIET', 'TOAN', 'TIENG_ANH', 'TIN_HOC', 'STEM'...
    name VARCHAR(255) NOT NULL,
    is_scored BOOLEAN DEFAULT FALSE,             -- TRUE: Môn có điểm kiểm tra; FALSE: Chỉ nhận xét
    applicable_grades INT[] DEFAULT '{1,2,3,4,5}'-- Áp dụng cho khối nào
);

-- 8. student_subject_evaluations (Sổ đánh giá thường xuyên môn học)
CREATE TABLE IF NOT EXISTS student_subject_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    term_id VARCHAR(32) REFERENCES terms(id) ON DELETE CASCADE,
    subject_id VARCHAR(50) REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID REFERENCES staffs(id) ON DELETE SET NULL,
    achievement_level VARCHAR(10) CHECK (achievement_level IN ('HTT', 'HT', 'CHT')),
    teacher_comment TEXT,
    ai_suggested_comment TEXT,                   -- Gợi ý của AI Pedagogical Copilot
    is_ai_assisted BOOLEAN DEFAULT FALSE,
    exam_score NUMERIC(4,1) CHECK (exam_score >= 1.0 AND exam_score <= 10.0),
    is_computer_based_test BOOLEAN DEFAULT FALSE,-- Kiểm tra trên máy tính theo TT 89/2026
    is_centralized_exam BOOLEAN DEFAULT FALSE,   -- Đề chung của Sở/Phòng lớp 5 (TT 89/2026)
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_student_term_subject UNIQUE (student_id, term_id, subject_id)
);

-- 9. student_behavior_evaluations (Sổ đánh giá phẩm chất và năng lực)
CREATE TABLE IF NOT EXISTS student_behavior_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    term_id VARCHAR(32) REFERENCES terms(id) ON DELETE CASCADE,
    evaluator_id UUID REFERENCES staffs(id) ON DELETE SET NULL,     -- Giáo viên chủ nhiệm
    
    -- 05 Phẩm chất chủ yếu: 'T' (Tốt), 'D' (Đạt), 'C' (Cần cố gắng)
    patriotism VARCHAR(5) CHECK (patriotism IN ('T', 'D', 'C')),        -- Yêu nước
    compassion VARCHAR(5) CHECK (compassion IN ('T', 'D', 'C')),        -- Nhân ái
    diligence VARCHAR(5) CHECK (diligence IN ('T', 'D', 'C')),          -- Chăm chỉ
    honesty VARCHAR(5) CHECK (honesty IN ('T', 'D', 'C')),              -- Trung thực
    responsibility VARCHAR(5) CHECK (responsibility IN ('T', 'D', 'C')),-- Trách nhiệm
    qualities_comment TEXT,
    
    -- 03 Năng lực chung: 'T', 'D', 'C'
    autonomy_and_self_learning VARCHAR(5) CHECK (autonomy_and_self_learning IN ('T', 'D', 'C')),
    communication_and_collaboration VARCHAR(5) CHECK (communication_and_collaboration IN ('T', 'D', 'C')),
    problem_solving_and_creativity VARCHAR(5) CHECK (problem_solving_and_creativity IN ('T', 'D', 'C')),
    
    -- 07 Năng lực đặc thù: 'T', 'D', 'C'
    language_competence VARCHAR(5) CHECK (language_competence IN ('T', 'D', 'C')),
    mathematical_competence VARCHAR(5) CHECK (mathematical_competence IN ('T', 'D', 'C')),
    scientific_competence VARCHAR(5) CHECK (scientific_competence IN ('T', 'D', 'C')),
    technological_competence VARCHAR(5) CHECK (technological_competence IN ('T', 'D', 'C')),
    computer_competence VARCHAR(5) CHECK (computer_competence IN ('T', 'D', 'C')),
    aesthetic_competence VARCHAR(5) CHECK (aesthetic_competence IN ('T', 'D', 'C')),
    physical_competence VARCHAR(5) CHECK (physical_competence IN ('T', 'D', 'C')),
    competencies_comment TEXT,
    
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_student_term_behavior UNIQUE (student_id, term_id)
);

-- 10. academic_transcripts_snapshots (Snapshot Học bạ số bất biến & Chữ ký số)
CREATE TABLE IF NOT EXISTS academic_transcripts_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id) ON DELETE RESTRICT,
    academic_year_id VARCHAR(32) NOT NULL,
    grade_level INT NOT NULL,
    regulation_applied VARCHAR(50) NOT NULL,     -- 'TT27_2020' hoặc 'TT89_2026'
    annual_reward VARCHAR(100),                  -- 'Học sinh Xuất sắc', 'Học sinh Tiêu biểu'...
    promotion_status VARCHAR(50) NOT NULL,       -- 'PROMOTED' (Lên lớp), 'RETAINED' (Ở lại lớp)
    snapshot_data JSONB NOT NULL,
    teacher_signature TEXT,                      -- Chữ ký số SmartCA của GVCN
    teacher_signed_at TIMESTAMPTZ,
    principal_signature TEXT,                    -- Chữ ký số Ban Giám hiệu
    principal_signed_at TIMESTAMPTZ,
    timestamp_token TEXT,                        -- Dấu thời gian TSA quốc gia
    is_locked BOOLEAN DEFAULT TRUE,
    CONSTRAINT uk_student_academic_year UNIQUE (student_id, academic_year_id)
);
CREATE INDEX IF NOT EXISTS idx_transcript_snapshot ON academic_transcripts_snapshots USING GIN (snapshot_data);

-- 11. boarding_menus (Danh mục Thực đơn & Khẩu phần ăn)
CREATE TABLE IF NOT EXISTS boarding_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    menu_date DATE NOT NULL,
    meal_type VARCHAR(20) NOT NULL,              -- 'LUNCH', 'SNACK', 'BREAKFAST'
    dishes JSONB NOT NULL,                       -- [{'name': 'Thịt kho trứng', 'weight_g': 120}]
    total_kcal NUMERIC(6,1) NOT NULL,            -- Năng lượng chuẩn: 600 - 750 Kcal bữa trưa
    protein_percent NUMERIC(4,1),                -- Tỷ lệ P: 13-20%
    lipid_percent NUMERIC(4,1),                  -- Tỷ lệ L: 20-30%
    glucid_percent NUMERIC(4,1),                 -- Tỷ lệ G: 50-65%
    sample_retention_verified BOOLEAN DEFAULT FALSE, -- Đã lưu mẫu thức ăn 24h
    sample_keeper_id UUID REFERENCES staffs(id) ON DELETE SET NULL,
    photo_url TEXT,                              -- Ảnh khay ăn thực tế tải lên App
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_boarding_menus_dishes ON boarding_menus USING GIN (dishes);

-- 12. boarding_daily_attendances (Điểm danh ăn bán trú & Báo hủy ăn tự động)
CREATE TABLE IF NOT EXISTS boarding_daily_attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    is_eating BOOLEAN DEFAULT TRUE,
    is_cancelled_by_parent BOOLEAN DEFAULT FALSE,
    cancellation_time TIMESTAMPTZ,               -- Hủy trước 8h00 sáng
    is_refund_eligible BOOLEAN DEFAULT FALSE,    -- Đủ điều kiện hoàn tiền ăn
    refund_amount NUMERIC(10,2) DEFAULT 0,
    CONSTRAINT uk_student_boarding_date UNIQUE (student_id, attendance_date)
);

-- 13. student_health_records (Hồ sơ Sức khỏe & Nhật ký Y tế Học sinh)
CREATE TABLE IF NOT EXISTS student_health_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    checkup_date DATE NOT NULL,
    height_cm NUMERIC(5,1) NOT NULL,
    weight_kg NUMERIC(5,1) NOT NULL,
    bmi NUMERIC(4,1) NOT NULL,
    bmi_category VARCHAR(50),                    -- 'Suy dinh dưỡng', 'Bình thường', 'Thừa cân', 'Béo phì'
    left_eye_vision VARCHAR(10),                 -- Thị lực: '10/10', '7/10'
    right_eye_vision VARCHAR(10),
    dental_condition VARCHAR(100),               -- Nha học đường
    chronic_conditions TEXT,                     -- Bệnh bẩm sinh/dị ứng
    medical_events JSONB DEFAULT '[]',           -- Nhật ký sơ cấp cứu tại phòng y tế
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_student_health_events ON student_health_records USING GIN (medical_events);

-- 14. kpi_criteria (Bộ tiêu chí chỉ số KPI)
CREATE TABLE IF NOT EXISTS kpi_criteria (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    applicable_role VARCHAR(50) NOT NULL,        -- 'HOMEROOM_TEACHER', 'PRIMARY_TEACHER', 'PRINCIPAL', 'ACCOUNTANT'...
    code VARCHAR(50) NOT NULL,                   -- 'KPI-GV-01', 'KPI-GV-02'...
    title VARCHAR(255) NOT NULL,
    description TEXT,
    weight NUMERIC(4,2) NOT NULL,                -- Trọng số (0.15, 0.20...)
    max_score NUMERIC(5,1) DEFAULT 100.0,
    source_type VARCHAR(50) NOT NULL,            -- 'SYSTEM_AUTOMATED', 'MANAGER_REVIEW', 'PEER_REVIEW', 'PARENT_CSAT'
    formula_identifier VARCHAR(100),             -- Tên hàm Go tính tự động: 'CalcDeadlinePunctuality'
    is_active BOOLEAN DEFAULT TRUE
);

-- 15. kpi_cycles (Kỳ đánh giá KPI)
CREATE TABLE IF NOT EXISTS kpi_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,                  -- Ví dụ: 'Đánh giá KPI Học kỳ 1 Năm học 2026-2027'
    cycle_type VARCHAR(20) NOT NULL,             -- 'MONTHLY', 'SEMESTER', 'ANNUAL'
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    submission_deadline DATE NOT NULL,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 16. kpi_evaluation_results (Bảng kết quả tổng hợp KPI & Xếp loại viên chức)
CREATE TABLE IF NOT EXISTS kpi_evaluation_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cycle_id UUID REFERENCES kpi_cycles(id) ON DELETE CASCADE,
    staff_id UUID REFERENCES staffs(id) ON DELETE RESTRICT,
    automated_system_score NUMERIC(5,2) DEFAULT 0, -- Chiều 1: Hệ thống tự động (40%)
    manager_assessment_score NUMERIC(5,2) DEFAULT 0,-- Chiều 2: BGH / Khối trưởng (30%)
    peer_review_score NUMERIC(5,2) DEFAULT 0,      -- Chiều 3: Đồng nghiệp (15%)
    parent_satisfaction_score NUMERIC(5,2) DEFAULT 0,-- Chiều 4: Phụ huynh (15%)
    final_kpi_score NUMERIC(5,2) NOT NULL,       -- Điểm thang 100
    official_classification VARCHAR(50) NOT NULL,-- 'EXCELLENT', 'GOOD', 'SATISFACTORY', 'UNSATISFACTORY'
    has_innovative_solution BOOLEAN DEFAULT FALSE,
    has_disciplinary_action BOOLEAN DEFAULT FALSE,
    detailed_metrics JSONB NOT NULL,
    reviewed_by UUID REFERENCES staffs(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    CONSTRAINT uk_staff_kpi_cycle UNIQUE (staff_id, cycle_id)
);
CREATE INDEX IF NOT EXISTS idx_kpi_eval_metrics ON kpi_evaluation_results USING GIN (detailed_metrics);

-- -----------------------------------------------------------------------------
-- Additional Compatibility Business Tables (for existing Admin table generators)
-- -----------------------------------------------------------------------------

-- kpi_evaluations (Used by internal/admin/tables/kpi.go)
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

-- teaching_quotas (Circular 05/2025/TT-BGDĐT)
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
-- PART 4: RICH BUSINESS SEED DATA (FOR ALL 16 TABLES)
-- =============================================================================

-- 1. Seed schools
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

-- 2. Seed academic_years
INSERT INTO academic_years (id, school_id, name, start_date, end_date, is_active)
VALUES 
(
    '2026-2027',
    'a0000000-0000-0000-0000-000000000001',
    'Năm học 2026-2027',
    '2026-09-01',
    '2027-05-31',
    TRUE
),
(
    '2025-2026',
    'a0000000-0000-0000-0000-000000000001',
    'Năm học 2025-2026',
    '2025-09-01',
    '2026-05-31',
    FALSE
) ON CONFLICT (id) DO NOTHING;

-- 3. Seed terms
INSERT INTO terms (id, academic_year_id, term_number, name, start_date, end_date, is_closed)
VALUES 
(
    '2026-2027-HK1',
    '2026-2027',
    1,
    'Học kỳ 1 Năm học 2026-2027',
    '2026-09-01',
    '2027-01-15',
    FALSE
),
(
    '2026-2027-HK2',
    '2026-2027',
    2,
    'Học kỳ 2 Năm học 2026-2027',
    '2027-01-16',
    '2027-05-31',
    FALSE
),
(
    '2025-2026-HK1',
    '2025-2026',
    1,
    'Học kỳ 1 Năm học 2025-2026',
    '2025-09-01',
    '2026-01-15',
    TRUE
),
(
    '2025-2026-HK2',
    '2025-2026',
    2,
    'Học kỳ 2 Năm học 2025-2026',
    '2026-01-16',
    '2026-05-31',
    TRUE
) ON CONFLICT (id) DO NOTHING;

-- 4. Seed staffs
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
    'nguyenvanan@thlevantam.edu.vn',
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
    'tranthibinh@thlevantam.edu.vn',
    'PRIMARY_TEACHER',
    'RANK_II',
    TRUE,
    805,
    23,
    7,
    '["GVCN", "HEAD_SUBJECT"]'::jsonb
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
    'lehoangcuong@thlevantam.edu.vn',
    'PRINCIPAL',
    'RANK_I',
    FALSE,
    805,
    4,
    19,
    '["PRINCIPAL"]'::jsonb
),
(
    'b0000000-0000-0000-0000-000000000004',
    'a0000000-0000-0000-0000-000000000001',
    '079192000004',
    'NV001',
    'Phạm Thị Dung',
    'FEMALE',
    '1992-08-25',
    '0904567890',
    'phamthidung@thlevantam.edu.vn',
    'NURSE',
    'RANK_III',
    FALSE,
    0,
    0,
    0,
    '["NURSE", "BOARDING_SAMPLE_KEEPER"]'::jsonb
) ON CONFLICT (employee_code) DO NOTHING;

-- 5. Seed classes
INSERT INTO classes (
    id, school_id, academic_year_id, grade_level, name, homeroom_teacher_id, room_number, student_count
) VALUES 
(
    'f0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    '2026-2027',
    1,
    '1A1',
    'b0000000-0000-0000-0000-000000000001',
    'P.101',
    32
),
(
    'f0000000-0000-0000-0000-000000000002',
    'a0000000-0000-0000-0000-000000000001',
    '2026-2027',
    5,
    '5A1',
    'b0000000-0000-0000-0000-000000000002',
    'P.301',
    30
) ON CONFLICT (id) DO NOTHING;

-- 6. Seed students
INSERT INTO students (
    id, school_id, national_id, student_code, full_name, birth_date, gender, ethnicity,
    current_class_id, status, is_inclusive_education, guardian_name, guardian_national_id,
    guardian_phone, guardian_relationship, has_digital_consent, digital_consent_verified,
    consent_signed_at, consent_method
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
    'f0000000-0000-0000-0000-000000000001',
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
    'f0000000-0000-0000-0000-000000000001',
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
    '2015-11-05',
    'MALE',
    'Kinh',
    'f0000000-0000-0000-0000-000000000002',
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
),
(
    'c0000000-0000-0000-0000-000000000004',
    'a0000000-0000-0000-0000-000000000001',
    '079215000004',
    'HS2026004',
    'Vũ Gia Huy',
    '2015-04-18',
    'MALE',
    'Kinh',
    'f0000000-0000-0000-0000-000000000002',
    'STUDYING',
    TRUE,
    'Vũ Tuấn Kiệt',
    '079183000006',
    '0906789012',
    'FATHER',
    TRUE,
    TRUE,
    NOW(),
    'VNEID_SSO'
) ON CONFLICT (student_code) DO UPDATE SET 
    current_class_id = EXCLUDED.current_class_id;

-- 7. Seed subjects
INSERT INTO subjects (id, name, is_scored, applicable_grades)
VALUES 
    ('TOAN', 'Toán', TRUE, '{1,2,3,4,5}'),
    ('TIENG_VIET', 'Tiếng Việt', TRUE, '{1,2,3,4,5}'),
    ('TIENG_ANH', 'Tiếng Anh', TRUE, '{1,2,3,4,5}'),
    ('TIN_HOC', 'Tin học và Công nghệ', TRUE, '{3,4,5}'),
    ('LICH_SU_DIA_LI', 'Lịch sử và Địa lí', TRUE, '{4,5}'),
    ('KHOA_HOC', 'Khoa học', TRUE, '{4,5}'),
    ('DAO_DUC', 'Đạo đức', FALSE, '{1,2,3,4,5}'),
    ('STEM', 'Hoạt động trải nghiệm / STEM', FALSE, '{1,2,3,4,5}'),
    ('THE_CHAT', 'Giáo dục thể chất', FALSE, '{1,2,3,4,5}'),
    ('MY_THUAT', 'Mĩ thuật', FALSE, '{1,2,3,4,5}'),
    ('AM_NHAC', 'Âm nhạc', FALSE, '{1,2,3,4,5}')
ON CONFLICT (id) DO NOTHING;

-- 8. Seed student_subject_evaluations
INSERT INTO student_subject_evaluations (
    id, student_id, term_id, subject_id, teacher_id,
    achievement_level, teacher_comment, ai_suggested_comment, is_ai_assisted,
    exam_score, is_computer_based_test, is_centralized_exam
) VALUES 
(
    '10000000-0000-0000-0000-000000000001',
    'c0000000-0000-0000-0000-000000000001',
    '2026-2027-HK1',
    'TOAN',
    'b0000000-0000-0000-0000-000000000001',
    'HTT',
    'Em tính toán nhanh nhẹn, tư duy tốt các bài toán hình học trực quan.',
    'Em nắm vững kiến thức toán học và tư duy logic tốt.',
    TRUE,
    9.5,
    FALSE,
    FALSE
),
(
    '10000000-0000-0000-0000-000000000002',
    'c0000000-0000-0000-0000-000000000001',
    '2026-2027-HK1',
    'TIENG_VIET',
    'b0000000-0000-0000-0000-000000000001',
    'HTT',
    'Chữ viết sạch đẹp, đọc trôi chảy, trả lời đúng trọng tâm câu hỏi.',
    'Em có vốn từ ngữ phong phú và diễn đạt cảm xúc tốt.',
    TRUE,
    9.0,
    FALSE,
    FALSE
),
(
    '10000000-0000-0000-0000-000000000003',
    'c0000000-0000-0000-0000-000000000002',
    '2026-2027-HK1',
    'TOAN',
    'b0000000-0000-0000-0000-000000000001',
    'HT',
    'Nắm chắc các phép tính trong phạm vi 10, cần cẩn thận hơn khi giải toán có lời văn.',
    'Em có tiến bộ trong học tập.',
    FALSE,
    8.0,
    FALSE,
    FALSE
),
(
    '10000000-0000-0000-0000-000000000004',
    'c0000000-0000-0000-0000-000000000003',
    '2026-2027-HK1',
    'TOAN',
    'b0000000-0000-0000-0000-000000000002',
    'HTT',
    'Làm bài kiểm tra trên máy tính theo chuẩn TT 89 đạt điểm tối đa phần trắc nghiệm.',
    'Em có tư duy phản biện xuất sắc và làm chủ công nghệ.',
    TRUE,
    9.5,
    TRUE,
    TRUE
) ON CONFLICT (student_id, term_id, subject_id) DO NOTHING;

-- 9. Seed student_behavior_evaluations
INSERT INTO student_behavior_evaluations (
    id, student_id, term_id, evaluator_id,
    patriotism, compassion, diligence, honesty, responsibility, qualities_comment,
    autonomy_and_self_learning, communication_and_collaboration, problem_solving_and_creativity,
    language_competence, mathematical_competence, scientific_competence,
    technological_competence, computer_competence, aesthetic_competence, physical_competence,
    competencies_comment
) VALUES 
(
    '20000000-0000-0000-0000-000000000001',
    'c0000000-0000-0000-0000-000000000001',
    '2026-2027-HK1',
    'b0000000-0000-0000-0000-000000000001',
    'T', 'T', 'T', 'T', 'T',
    'Em ngoan ngoãn, lễ phép với thầy cô, chan hòa với bạn bè trong lớp.',
    'T', 'T', 'T',
    'T', 'T', 'T', 'T', 'T', 'T', 'T',
    'Năng lực tự chủ và giải quyết vấn đề xuất sắc, tích cực phát biểu.'
),
(
    '20000000-0000-0000-0000-000000000002',
    'c0000000-0000-0000-0000-000000000002',
    '2026-2027-HK1',
    'b0000000-0000-0000-0000-000000000001',
    'T', 'T', 'T', 'T', 'T',
    'Tích cực tham gia các phong trào văn nghệ và giữ gìn vệ sinh lớp học.',
    'T', 'T', 'D',
    'T', 'D', 'D', 'D', 'T', 'T', 'T',
    'Hợp tác tốt trong nhóm, cần chủ động hơn khi làm bài tập sáng tạo.'
),
(
    '20000000-0000-0000-0000-000000000003',
    'c0000000-0000-0000-0000-000000000003',
    '2026-2027-HK1',
    'b0000000-0000-0000-0000-000000000002',
    'T', 'T', 'T', 'T', 'T',
    'Có tinh thần trách nhiệm cao, luôn hoàn thành nhiệm vụ cán sự lớp.',
    'T', 'T', 'T',
    'T', 'T', 'T', 'T', 'T', 'T', 'T',
    'Năng lực toàn diện, có khả năng lãnh đạo nhóm xuất sắc.'
) ON CONFLICT (student_id, term_id) DO NOTHING;

-- 10. Seed academic_transcripts_snapshots
INSERT INTO academic_transcripts_snapshots (
    id, student_id, academic_year_id, grade_level, regulation_applied,
    annual_reward, promotion_status, snapshot_data,
    teacher_signature, teacher_signed_at, principal_signature, principal_signed_at,
    timestamp_token, is_locked
) VALUES 
(
    '30000000-0000-0000-0000-000000000001',
    'c0000000-0000-0000-0000-000000000001',
    '2026-2027',
    1,
    'TT27_2020',
    'Học sinh Xuất sắc',
    'PROMOTED',
    '{
        "student_code": "HS2026001",
        "full_name": "Nguyễn Hoàng Nam",
        "grade_level": 1,
        "class_name": "1A1",
        "regulation": "TT27_2020",
        "subjects": {
            "TOAN": {"score": 9.5, "level": "HTT"},
            "TIENG_VIET": {"score": 9.0, "level": "HTT"}
        },
        "behavior": {
            "qualities": "T",
            "competencies": "T"
        },
        "summary": "Hoàn thành xuất sắc các nội dung học tập và rèn luyện"
    }'::jsonb,
    'SMARTCA_SIG_GVCN_079185000001',
    NOW(),
    'SMARTCA_SIG_BGH_079180000003',
    NOW(),
    'TSA_VNPT_TIMESTAMP_TOKEN_20270531_001',
    TRUE
),
(
    '30000000-0000-0000-0000-000000000002',
    'c0000000-0000-0000-0000-000000000003',
    '2026-2027',
    5,
    'TT89_2026',
    'Học sinh Xuất sắc',
    'PROMOTED',
    '{
        "student_code": "HS2026003",
        "full_name": "Lê Quốc Bảo",
        "grade_level": 5,
        "class_name": "5A1",
        "regulation": "TT89_2026",
        "computer_exam_participated": true,
        "centralized_exam_passed": true,
        "subjects": {
            "TOAN": {"score": 9.5, "level": "HTT"},
            "TIENG_VIET": {"score": 9.0, "level": "HTT"}
        },
        "behavior": {
            "qualities": "T",
            "competencies": "T"
        },
        "summary": "Hoàn thành chương trình Tiểu học - Đạt danh hiệu Học sinh Xuất sắc"
    }'::jsonb,
    'SMARTCA_SIG_GVCN_079188000002',
    NOW(),
    'SMARTCA_SIG_BGH_079180000003',
    NOW(),
    'TSA_VNPT_TIMESTAMP_TOKEN_20270531_002',
    TRUE
) ON CONFLICT (student_id, academic_year_id) DO NOTHING;

-- 11. Seed boarding_menus
INSERT INTO boarding_menus (
    id, school_id, menu_date, meal_type, dishes, total_kcal,
    protein_percent, lipid_percent, glucid_percent,
    sample_retention_verified, sample_keeper_id, photo_url
) VALUES 
(
    '40000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    '2026-10-28',
    'LUNCH',
    '[
        {"name": "Thịt ba chỉ kho trứng cút", "weight_g": 120},
        {"name": "Canh bí đỏ nấu tôm", "weight_g": 150},
        {"name": "Rau bắp cải xào", "weight_g": 80},
        {"name": "Cơm trắng", "weight_g": 180},
        {"name": "Tráng miệng chuối cau", "weight_g": 60}
    ]'::jsonb,
    680.5,
    16.5,
    24.0,
    59.5,
    TRUE,
    'b0000000-0000-0000-0000-000000000004',
    '/uploads/menus/20261028_lunch.jpg'
),
(
    '40000000-0000-0000-0000-000000000002',
    'a0000000-0000-0000-0000-000000000001',
    '2026-10-29',
    'LUNCH',
    '[
        {"name": "Gà om nấm đông cô", "weight_g": 130},
        {"name": "Canh chua cá lóc rau muống", "weight_g": 160},
        {"name": "Đậu que xào thịt bò", "weight_g": 80},
        {"name": "Cơm gạo tám thơm", "weight_g": 180},
        {"name": "Tráng miệng thanh long", "weight_g": 80}
    ]'::jsonb,
    710.0,
    17.0,
    23.5,
    59.5,
    TRUE,
    'b0000000-0000-0000-0000-000000000004',
    '/uploads/menus/20261029_lunch.jpg'
) ON CONFLICT (id) DO NOTHING;

-- 12. Seed boarding_daily_attendances
INSERT INTO boarding_daily_attendances (
    id, student_id, attendance_date, is_eating, is_cancelled_by_parent,
    cancellation_time, is_refund_eligible, refund_amount
) VALUES 
(
    'a1000000-0000-0000-0000-000000000001',
    'c0000000-0000-0000-0000-000000000001',
    '2026-10-29',
    TRUE,
    FALSE,
    NULL,
    FALSE,
    0.00
),
(
    'a1000000-0000-0000-0000-000000000002',
    'c0000000-0000-0000-0000-000000000002',
    '2026-10-29',
    TRUE,
    FALSE,
    NULL,
    FALSE,
    0.00
),
(
    'a1000000-0000-0000-0000-000000000003',
    'c0000000-0000-0000-0000-000000000003',
    '2026-10-29',
    FALSE,
    TRUE,
    '2026-10-29 07:30:00+07',
    TRUE,
    35000.00
) ON CONFLICT (student_id, attendance_date) DO NOTHING;

-- 13. Seed student_health_records
INSERT INTO student_health_records (
    id, student_id, checkup_date, height_cm, weight_kg, bmi, bmi_category,
    left_eye_vision, right_eye_vision, dental_condition, chronic_conditions, medical_events
) VALUES 
(
    '60000000-0000-0000-0000-000000000001',
    'c0000000-0000-0000-0000-000000000001',
    '2026-09-15',
    118.5,
    21.0,
    14.9,
    'Bình thường',
    '10/10',
    '10/10',
    'Răng tốt, không sâu',
    'Không có tiền sử bệnh lý mãn tính',
    '[
        {"date": "2026-10-10", "type": "Trầy xước nhẹ", "treatment": "Sát khuẩn bằng Betadine và băng gạc cá nhân"}
    ]'::jsonb
),
(
    '60000000-0000-0000-0000-000000000002',
    'c0000000-0000-0000-0000-000000000002',
    '2026-09-15',
    116.0,
    20.5,
    15.2,
    'Bình thường',
    '10/10',
    '10/10',
    'Hàn 01 răng sâu số 4',
    'Dị ứng nhẹ với phấn hoa',
    '[]'::jsonb
),
(
    '60000000-0000-0000-0000-000000000003',
    'c0000000-0000-0000-0000-000000000003',
    '2026-09-15',
    142.0,
    36.0,
    17.9,
    'Bình thường',
    '8/10',
    '9/10',
    'Răng vĩnh viễn phát triển tốt',
    'Cận thị nhẹ độ 0.5',
    '[]'::jsonb
) ON CONFLICT (id) DO NOTHING;

-- 14. Seed kpi_criteria
INSERT INTO kpi_criteria (
    id, school_id, applicable_role, code, title, description, weight, max_score, source_type, formula_identifier, is_active
) VALUES 
(
    '70000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    'PRIMARY_TEACHER',
    'KPI-GV-01',
    'Chấp hành quy chế chuyên môn và tiến độ vào sổ điểm',
    'Đo lường tự động tỷ lệ nộp giáo án, cập nhật điểm số đúng thời hạn trên hệ thống',
    0.20,
    100.0,
    'SYSTEM_AUTOMATED',
    'CalcDeadlinePunctuality',
    TRUE
),
(
    '70000000-0000-0000-0000-000000000002',
    'a0000000-0000-0000-0000-000000000001',
    'PRIMARY_TEACHER',
    'KPI-GV-02',
    'Chất lượng hồ sơ bài dạy (CV 2345/STEM TT 86) & Dự giờ',
    'Đánh giá chất lượng chuyên môn qua các tiết dạy thao giảng và giáo án',
    0.20,
    100.0,
    'MANAGER_REVIEW',
    'CalcLessonPlanQuality',
    TRUE
),
(
    '70000000-0000-0000-0000-000000000003',
    'a0000000-0000-0000-0000-000000000001',
    'PRIMARY_TEACHER',
    'KPI-GV-03',
    'Sáng kiến kinh nghiệm & Đổi mới phương pháp dạy học',
    'Sáng kiến kinh nghiệm được hội đồng khoa học cấp trường hoặc cấp quận/huyện công nhận',
    0.15,
    100.0,
    'SYSTEM_AUTOMATED',
    'CalcInnovativeSolutionScore',
    TRUE
),
(
    '70000000-0000-0000-0000-000000000004',
    'a0000000-0000-0000-0000-000000000001',
    'PRIMARY_TEACHER',
    'KPI-GV-04',
    'Mức độ hài lòng của phụ huynh học sinh (CSAT)',
    'Khảo sát phụ huynh học sinh định kỳ qua cổng thông tin điện tử',
    0.15,
    100.0,
    'PARENT_CSAT',
    'CalcParentCSATScore',
    TRUE
) ON CONFLICT (id) DO NOTHING;

-- 15. Seed kpi_cycles
INSERT INTO kpi_cycles (
    id, school_id, name, cycle_type, start_date, end_date, submission_deadline, is_locked
) VALUES 
(
    'c1000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    'Đánh giá KPI Học kỳ 1 Năm học 2026-2027',
    'SEMESTER',
    '2026-09-01',
    '2027-01-15',
    '2027-01-10',
    FALSE
),
(
    'c1000000-0000-0000-0000-000000000002',
    'a0000000-0000-0000-0000-000000000001',
    'Đánh giá KPI Năm học 2026-2027',
    'ANNUAL',
    '2026-09-01',
    '2027-05-31',
    '2027-05-25',
    FALSE
) ON CONFLICT (id) DO NOTHING;

-- 16. Seed kpi_evaluation_results
INSERT INTO kpi_evaluation_results (
    id, cycle_id, staff_id,
    automated_system_score, manager_assessment_score, peer_review_score, parent_satisfaction_score,
    final_kpi_score, official_classification, has_innovative_solution, has_disciplinary_action,
    detailed_metrics, reviewed_by, reviewed_at
) VALUES 
(
    '90000000-0000-0000-0000-000000000001',
    'c1000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000001',
    95.00,
    92.00,
    90.00,
    94.00,
    93.20,
    'HOÀN THÀNH XUẤT SẮC NHIỆM VỤ',
    TRUE,
    FALSE,
    '{"KPI-GV-01": 95, "KPI-GV-02": 92, "KPI-GV-03": 90, "KPI-GV-04": 94}'::jsonb,
    'b0000000-0000-0000-0000-000000000003',
    NOW()
),
(
    '90000000-0000-0000-0000-000000000002',
    'c1000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000002',
    85.00,
    88.00,
    84.00,
    86.00,
    85.90,
    'HOÀN THÀNH TỐT NHIỆM VỤ',
    FALSE,
    FALSE,
    '{"KPI-GV-01": 85, "KPI-GV-02": 88, "KPI-GV-03": 84, "KPI-GV-04": 86}'::jsonb,
    'b0000000-0000-0000-0000-000000000003',
    NOW()
),
(
    '90000000-0000-0000-0000-000000000003',
    'c1000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000003',
    91.00,
    90.00,
    88.00,
    92.00,
    90.40,
    'HOÀN THÀNH TỐT NHIỆM VỤ',
    TRUE,
    FALSE,
    '{"KPI-GV-01": 91, "KPI-GV-02": 90, "KPI-GV-03": 88, "KPI-GV-04": 92, "capped_by_20pct": true}'::jsonb,
    'b0000000-0000-0000-0000-000000000003',
    NOW()
) ON CONFLICT (staff_id, cycle_id) DO NOTHING;

-- Seed teaching_quotas (for backward compatibility)
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

-- Seed kpi_evaluations (for backward compatibility)
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
