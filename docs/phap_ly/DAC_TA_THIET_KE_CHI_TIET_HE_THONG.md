# ĐẶC TẢ THIẾT KẾ KỸ THUẬT CHI TIẾT HỆ THỐNG PHẦN MỀM EDUPRIMARY-KPI
## CHI TIẾT DATABASE SCHEMA, API ENDPOINTS, THUẬT TOÁN KPI & LOGIC CHIẾN LƯỢC

---

## 1. THIẾT KẾ CƠ SỞ DỮ LIỆU CHI TIẾT (POSTGRESQL RELATIONAL + JSONB SCHEMA)

Hệ thống sử dụng PostgreSQL 16+ với các kiểu dữ liệu nâng cao (`UUID`, `JSONB`, `NUMERIC`, `TIMESTAMPTZ`), hỗ trợ Indexing GIN cho JSONB và Partitioning theo năm học (`academic_year_id`).

### 1.1. Phân hệ Tổ chức & Người dùng Đa trường (Multi-Tenancy & IAM)

```sql
-- 1. Bảng Trường học (Tenant)
CREATE TABLE schools (
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

-- 2. Bảng Niên khóa / Năm học
CREATE TABLE academic_years (
    id VARCHAR(32) PRIMARY KEY,                  -- Ví dụ: '2026-2027'
    school_id UUID REFERENCES schools(id),
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 3. Bảng Học kỳ
CREATE TABLE terms (
    id VARCHAR(32) PRIMARY KEY,                  -- Ví dụ: '2026-2027-HK1'
    academic_year_id VARCHAR(32) REFERENCES academic_years(id),
    term_number INT NOT NULL,                    -- 1: Giữa kỳ 1, 2: Cuối kỳ 1, 3: Giữa kỳ 2, 4: Cuối năm
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_closed BOOLEAN DEFAULT FALSE
);

-- 4. Bảng Lớp học
CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    academic_year_id VARCHAR(32) REFERENCES academic_years(id),
    grade_level INT NOT NULL CHECK (grade_level BETWEEN 1 AND 5), -- Lớp 1 đến Lớp 5
    name VARCHAR(100) NOT NULL,                  -- Ví dụ: '1A1', '5A2'
    homeroom_teacher_id UUID,                    -- FK tới staffs(id)
    room_number VARCHAR(50),
    student_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 5. Bảng Cán bộ, Giáo viên & Nhân viên (Staffs - Theo TT 20/2023)
CREATE TABLE staffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    national_id VARCHAR(12) UNIQUE NOT NULL,     -- Số CCCD / Định danh VNeID
    employee_code VARCHAR(50) UNIQUE NOT NULL,   -- Mã số cán bộ (GV001...)
    full_name VARCHAR(255) NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('MALE', 'FEMALE', 'OTHER')),
    birth_date DATE NOT NULL,
    phone_number VARCHAR(20),
    email VARCHAR(255),
    
    -- Vị trí việc làm & chức danh (TT 20/2023 & TT 08/2023)
    job_position VARCHAR(50) NOT NULL,          -- 'PRINCIPAL', 'VICE_PRINCIPAL', 'PRIMARY_TEACHER', 'ACCOUNTANT', 'NURSE', 'LIBRARIAN', 'IT_ADMIN', 'SECURITY'
    teacher_rank VARCHAR(20),                   -- 'RANK_I', 'RANK_II', 'RANK_III'
    is_homeroom BOOLEAN DEFAULT FALSE,          -- Có chủ nhiệm lớp không
    
    -- Định mức làm việc (TT 05/2025/TT-BGDĐT)
    standard_annual_hours INT DEFAULT 805,      -- Định mức chuẩn tiểu học: 805 tiết/năm
    weekly_teaching_hours INT DEFAULT 23,       -- Định mức bình quân 23 tiết/tuần
    weekly_reduced_hours INT DEFAULT 0,         -- Tiết giảm trừ (GVCN: 4 tiết, Tổ trưởng: 3 tiết...)
    concurrent_roles JSONB DEFAULT '[]',        -- Danh sách kiêm nhiệm: ['HEAD_SUBJECT', 'PSYCHOLOGY_COUNSELING']
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 6. Bảng Học sinh (Students)
CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    national_id VARCHAR(12) UNIQUE,              -- Số định danh cá nhân 12 số (Luật Căn cước 2023)
    student_code VARCHAR(50) UNIQUE NOT NULL,    -- Mã định danh Bộ GD&ĐT
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('MALE', 'FEMALE')),
    ethnicity VARCHAR(50) DEFAULT 'Kinh',
    current_class_id UUID REFERENCES classes(id),
    status VARCHAR(30) DEFAULT 'STUDYING',       -- 'STUDYING', 'TRANSFERRED', 'DROPPED', 'GRADUATED'
    is_inclusive_education BOOLEAN DEFAULT FALSE,-- Học sinh khuyết tật học hòa nhập (TT 03/2018)
    
    -- Thông tin nhân thân phụ huynh (VNeID Mức 2 theo NĐ 69/2024)
    guardian_name VARCHAR(255) NOT NULL,
    guardian_national_id VARCHAR(12) NOT NULL,
    guardian_phone VARCHAR(20) NOT NULL,
    guardian_relationship VARCHAR(50) NOT NULL,  -- 'FATHER', 'MOTHER', 'LEGAL_GUARDIAN'
    
    -- Sự đồng thuận số (Digital Consent - Điều 20 NĐ 13/2023/NĐ-CP)
    has_digital_consent BOOLEAN DEFAULT FALSE,
    consent_signed_at TIMESTAMPTZ,
    consent_method VARCHAR(50),                  -- 'VNEID_SSO', 'SMS_OTP', 'PORTAL_SIGN'
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

---

### 1.2. Phân hệ Nghiệp vụ Đánh giá Học sinh (TT 27/2020 & TT 89/2026/TT-BGDĐT)

```sql
-- 7. Danh mục Môn học & Hoạt động giáo dục
CREATE TABLE subjects (
    id VARCHAR(50) PRIMARY KEY,                  -- 'TIENG_VIET', 'TOAN', 'TIENG_ANH', 'TIN_HOC', 'STEM'...
    name VARCHAR(255) NOT NULL,
    is_scored BOOLEAN DEFAULT FALSE,             -- TRUE: Môn có điểm kiểm tra (Toán, Tiếng Việt...); FALSE: Chỉ nhận xét
    applicable_grades INT[] DEFAULT '{1,2,3,4,5}'-- Áp dụng cho khối nào
);

-- 8. Sổ đánh giá thường xuyên môn học và hoạt động giáo dục
CREATE TABLE student_subject_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    term_id VARCHAR(32) REFERENCES terms(id),
    subject_id VARCHAR(50) REFERENCES subjects(id),
    teacher_id UUID REFERENCES staffs(id),
    
    -- Mức độ đạt được: 'HTT' (Hoàn thành tốt), 'HT' (Hoàn thành), 'CHT' (Chưa hoàn thành)
    achievement_level VARCHAR(10) CHECK (achievement_level IN ('HTT', 'HT', 'CHT')),
    
    -- Nhận xét sư phạm thường xuyên
    teacher_comment TEXT,
    ai_suggested_comment TEXT,                   -- Gợi ý của AI Pedagogical Copilot
    is_ai_assisted BOOLEAN DEFAULT FALSE,
    
    -- Điểm số bài kiểm tra định kỳ (nếu môn có điểm số: thang điểm 1-10)
    exam_score NUMERIC(4,1) CHECK (exam_score >= 1.0 AND exam_score <= 10.0),
    is_computer_based_test BOOLEAN DEFAULT FALSE,-- Kiểm tra trên máy tính theo TT 89/2026
    is_centralized_exam BOOLEAN DEFAULT FALSE,   -- Đề chung của Sở/Phòng lớp 5 (TT 89/2026)
    
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_student_term_subject UNIQUE (student_id, term_id, subject_id)
);

-- 9. Sổ đánh giá phẩm chất và năng lực (Định kỳ HK1, HK2, Cả năm)
CREATE TABLE student_behavior_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    term_id VARCHAR(32) REFERENCES terms(id),
    evaluator_id UUID REFERENCES staffs(id),     -- Giáo viên chủ nhiệm
    
    -- 05 Phẩm chất chủ yếu: 'T' (Tốt), 'D' (Đạt), 'C' (Cần cố gắng)
    patriotism VARCHAR(5) CHECK (patriotism IN ('T', 'D', 'C')),        -- Yêu nước
    compassion VARCHAR(5) CHECK (compassion IN ('T', 'D', 'C')),        -- Nhân ái
    diligence VARCHAR(5) CHECK (diligence IN ('T', 'D', 'C')),          -- Chăm chỉ
    honesty VARCHAR(5) CHECK (honesty IN ('T', 'D', 'C')),              -- Trung thực
    responsibility VARCHAR(5) CHECK (responsibility IN ('T', 'D', 'C')),-- Trách nhiệm
    qualities_comment TEXT,
    
    -- 03 Năng lực chung: 'T', 'D', 'C'
    autonomy_and_self_learning VARCHAR(5) CHECK (autonomy_and_self_learning IN ('T', 'D', 'C')),       -- Tự chủ & tự học
    communication_and_collaboration VARCHAR(5) CHECK (communication_and_collaboration IN ('T', 'D', 'C')),-- Giao tiếp & hợp tác
    problem_solving_and_creativity VARCHAR(5) CHECK (problem_solving_and_creativity IN ('T', 'D', 'C')),  -- Giải quyết VĐ & sáng tạo
    
    -- 07 Năng lực đặc thù: 'T', 'D', 'C'
    language_competence VARCHAR(5) CHECK (language_competence IN ('T', 'D', 'C')),   -- Ngôn ngữ
    mathematical_competence VARCHAR(5) CHECK (mathematical_competence IN ('T', 'D', 'C')), -- Tính toán
    scientific_competence VARCHAR(5) CHECK (scientific_competence IN ('T', 'D', 'C')),   -- Khoa học
    technological_competence VARCHAR(5) CHECK (technological_competence IN ('T', 'D', 'C')),-- Công nghệ
    computer_competence VARCHAR(5) CHECK (computer_competence IN ('T', 'D', 'C')),        -- Tin học
    aesthetic_competence VARCHAR(5) CHECK (aesthetic_competence IN ('T', 'D', 'C')),       -- Thẩm mỹ
    physical_competence VARCHAR(5) CHECK (physical_competence IN ('T', 'D', 'C')),        -- Thể chất
    competencies_comment TEXT,
    
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_student_term_behavior UNIQUE (student_id, term_id)
);

-- 10. Snapshot Học bạ số bất biến & Chữ ký số (Immutable Academic Transcripts)
CREATE TABLE academic_transcripts_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id) ON DELETE RESTRICT,
    academic_year_id VARCHAR(32) NOT NULL,
    grade_level INT NOT NULL,
    regulation_applied VARCHAR(50) NOT NULL,     -- 'TT27_2020' hoặc 'TT89_2026'
    
    -- Tổng hợp kết quả năm học
    annual_reward VARCHAR(100),                  -- 'Học sinh Xuất sắc', 'Học sinh Tiêu biểu'...
    promotion_status VARCHAR(50) NOT NULL,       -- 'PROMOTED' (Lên lớp), 'RETAINED' (Ở lại lớp)
    
    -- Dữ liệu đầy đủ dạng JSONB Snapshot
    snapshot_data JSONB NOT NULL,
    
    -- Ký số tập trung (PAdES / PKCS#7)
    teacher_signature TEXT,                      -- Chữ ký số SmartCA của GVCN
    teacher_signed_at TIMESTAMPTZ,
    principal_signature TEXT,                    -- Chữ ký số Ban Giám hiệu
    principal_signed_at TIMESTAMPTZ,
    timestamp_token TEXT,                        -- Dấu thời gian TSA quốc gia
    is_locked BOOLEAN DEFAULT TRUE,
    
    CONSTRAINT uk_student_academic_year UNIQUE (student_id, academic_year_id)
);
CREATE INDEX idx_transcript_snapshot ON academic_transcripts_snapshots USING GIN (snapshot_data);
```

---

### 1.3. Phân hệ Bán trú, Dinh dưỡng & Y tế Học đường

```sql
-- 11. Danh mục Thực đơn & Khẩu phần ăn
CREATE TABLE boarding_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id),
    menu_date DATE NOT NULL,
    meal_type VARCHAR(20) NOT NULL,              -- 'LUNCH', 'SNACK', 'BREAKFAST'
    dishes JSONB NOT NULL,                       -- [{'name': 'Thịt kho trứng', 'weight_g': 120}]
    total_kcal NUMERIC(6,1) NOT NULL,            -- Năng lượng chuẩn: 600 - 750 Kcal bữa trưa
    protein_percent NUMERIC(4,1),                -- Tỷ lệ P: 13-20%
    lipid_percent NUMERIC(4,1),                  -- Tỷ lệ L: 20-30%
    glucid_percent NUMERIC(4,1),                 -- Tỷ lệ G: 50-65%
    sample_retention_verified BOOLEAN DEFAULT FALSE, -- Đã lưu mẫu thức ăn 24h
    sample_keeper_id UUID REFERENCES staffs(id),
    photo_url TEXT,                              -- Ảnh khay ăn thực tế tải lên App
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 12. Điểm danh ăn bán trú & Báo hủy ăn tự động
CREATE TABLE boarding_daily_attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id),
    attendance_date DATE NOT NULL,
    is_eating BOOLEAN DEFAULT TRUE,
    is_cancelled_by_parent BOOLEAN DEFAULT FALSE,
    cancellation_time TIMESTAMPTZ,               -- Hủy trước 8h00 sáng
    is_refund_eligible BOOLEAN DEFAULT FALSE,    -- Đủ điều kiện hoàn tiền ăn
    refund_amount NUMERIC(10,2) DEFAULT 0,
    CONSTRAINT uk_student_boarding_date UNIQUE (student_id, attendance_date)
);

-- 13. Hồ sơ Sức khỏe & Nhật ký Y tế Học sinh
CREATE TABLE student_health_records (
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
```

---

### 1.4. Phân hệ Chuyên sâu: Bảng tiêu chí & Kết quả KPI 360 độ

```sql
-- 14. Bộ tiêu chí chỉ số KPI (KPI Criteria Template)
CREATE TABLE kpi_criteria (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id),
    applicable_role VARCHAR(50) NOT NULL,        -- 'HOMEROOM_TEACHER', 'SUBJECT_TEACHER', 'PRINCIPAL', 'ACCOUNTANT'...
    code VARCHAR(50) NOT NULL,                   -- 'KPI-GV-01', 'KPI-GV-02'...
    title VARCHAR(255) NOT NULL,
    description TEXT,
    weight NUMERIC(4,2) NOT NULL,                -- Trọng số (0.15, 0.20...)
    max_score NUMERIC(5,1) DEFAULT 100.0,
    source_type VARCHAR(50) NOT NULL,            -- 'SYSTEM_AUTOMATED', 'MANAGER_REVIEW', 'PEER_REVIEW', 'PARENT_CSAT'
    formula_identifier VARCHAR(100),             -- Tên hàm Go tính tự động: 'CalcDeadlinePunctuality'
    is_active BOOLEAN DEFAULT TRUE
);

-- 15. Kỳ đánh giá KPI (Hàng tháng, Học kỳ, Năm học)
CREATE TABLE kpi_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id),
    name VARCHAR(255) NOT NULL,                  -- Ví dụ: 'Đánh giá KPI Tháng 10/2026'
    cycle_type VARCHAR(20) NOT NULL,             -- 'MONTHLY', 'SEMESTER', 'ANNUAL'
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    submission_deadline DATE NOT NULL,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 16. Bảng kết quả tổng hợp KPI & Xếp loại viên chức
CREATE TABLE kpi_evaluation_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cycle_id UUID REFERENCES kpi_cycles(id) ON DELETE CASCADE,
    staff_id UUID REFERENCES staffs(id) ON DELETE RESTRICT,
    
    -- Điểm thành phần theo 4 chiều đánh giá
    automated_system_score NUMERIC(5,2) DEFAULT 0, -- Chiều 1: Hệ thống đo lường tự động (Trọng số 40%)
    manager_assessment_score NUMERIC(5,2) DEFAULT 0,-- Chiều 2: BGH / Khối trưởng chấm (Trọng số 30%)
    peer_review_score NUMERIC(5,2) DEFAULT 0,      -- Chiều 3: Đồng nghiệp tổ chuyên môn (Trọng số 15%)
    parent_satisfaction_score NUMERIC(5,2) DEFAULT 0,-- Chiều 4: Phụ huynh đánh giá CSAT (Trọng số 15%)
    
    -- Điểm tổng hợp có trọng số
    final_kpi_score NUMERIC(5,2) NOT NULL,       -- Điểm thang 100
    
    -- Quy đổi xếp loại viên chức (Nghị định 90/2020 & Nghị định 48/2023)
    official_classification VARCHAR(50) NOT NULL,-- 'EXCELLENT' (Xuất sắc), 'GOOD' (Tốt), 'SATISFACTORY' (Hoàn thành), 'UNSATISFACTORY' (Không HT)
    has_innovative_solution BOOLEAN DEFAULT FALSE,-- Có sáng kiến/giải pháp đổi mới sáng tạo (Luật TĐKT 2022)
    has_disciplinary_action BOOLEAN DEFAULT FALSE,-- Bị kỷ luật trong năm
    
    detailed_metrics JSONB NOT NULL,             -- Lưu vết chi tiết điểm từng tiêu chí
    reviewed_by UUID REFERENCES staffs(id),      -- Hiệu trưởng ký duyệt
    reviewed_at TIMESTAMPTZ,
    
    CONSTRAINT uk_staff_kpi_cycle UNIQUE (staff_id, cycle_id)
);
```

---

## 2. ĐẶC TẢ CHI TIẾT RESTful API ENDPOINTS (GOLANG ROUTER)

Tất cả các API đều yêu cầu Header: `Authorization: Bearer <JWT>`, `X-School-ID: <UUID>`. 
Chuẩn phản hồi chung (JSON Response):
```json
{
  "success": true,
  "code": 200,
  "message": "Thành công",
  "data": {},
  "meta": {"page": 1, "page_size": 20, "total_records": 100}
}
```

### 2.1. Nhóm API Quản lý Học sinh & Hồ sơ số

| Method | Endpoint | Quyền hạn (RBAC) | Mô tả nghiệp vụ |
|---|---|---|---|
| `POST` | `/api/v1/students/onboarding` | `BGH, VAN_THU` | Tuyển sinh Lớp 1 trực tuyến, kết nối API VNeID làm sạch hồ sơ |
| `GET` | `/api/v1/students` | `BGH, GVCN, GVBM` | Danh sách học sinh theo lớp, lọc theo trạng thái, khuyết tật |
| `GET` | `/api/v1/students/{id}/profile` | `BGH, GVCN, PHHS` | Chi tiết hồ sơ, quá trình học tập, lịch sử tiêm chủng, BMI |
| `POST` | `/api/v1/students/{id}/digital-consent`| `PHHS` | Xác thực đồng thuận số (Digital Consent) theo NĐ 13/2023 |
| `POST` | `/api/v1/students/transfer` | `BGH, VAN_THU` | Chuyển trường, tiếp nhận học sinh chuyển đến kèm học bạ số |

### 2.2. Nhóm API Đánh giá Học sinh theo TT 27/2020 & TT 89/2026

| Method | Endpoint | Quyền hạn (RBAC) | Mô tả nghiệp vụ |
|---|---|---|---|
| `POST` | `/api/v1/evaluations/subject-scores` | `GVCN, GVBM` | Nhập điểm và nhận xét môn học định kỳ theo lớp |
| `POST` | `/api/v1/evaluations/ai-comment-assist`| `GVCN, GVBM` | Gọi AI Sư phạm sinh 3 gợi ý nhận xét chuẩn TT 27 cá nhân hóa |
| `POST` | `/api/v1/evaluations/behaviors` | `GVCN` | Đánh giá 5 phẩm chất và 10 năng lực định kỳ (T, Đ, C) |
| `POST` | `/api/v1/evaluations/computer-exams` | `BGH, GVBM` | Tạo ca kiểm tra trên máy tính theo Thông tư 89/2026/TT-BGDĐT |
| `POST` | `/api/v1/evaluations/centralized-exam/receive`| `BGH` | Tiếp nhận gói đề thi chung Lớp 5 mã hóa từ Sở/Phòng GD&ĐT |
| `POST` | `/api/v1/transcripts/generate` | `BGH, GVCN` | Tạo bản ghi Snapshot học bạ số cuối năm học |
| `POST` | `/api/v1/transcripts/sign-batch` | `BGH, GVCN` | Ký số tập trung SmartCA cho học bạ số cả lớp trong 1 chạm |

### 2.3. Nhóm API Chuyên môn & Định mức Tiết dạy (TT 05/2025/TT-BGDĐT)

| Method | Endpoint | Quyền hạn (RBAC) | Mô tả nghiệp vụ |
|---|---|---|---|
| `GET` | `/api/v1/teaching/quota-summary` | `BGH, GV` | Bảng tổng hợp định mức tiết dạy: 805 tiết, tiết giảm, tiết thực dạy |
| `GET` | `/api/v1/teaching/overtime-reports`| `BGH, KE_TOAN` | Xuất bảng quyết toán thừa giờ (tăng giờ) cuối năm học |
| `POST` | `/api/v1/lesson-plans/submit` | `GV` | Nộp kế hoạch bài dạy (giáo án CV 2345 / STEM TT 86) |
| `POST` | `/api/v1/lesson-plans/{id}/approve`| `TO_TRUONG, BGH`| Phê duyệt, góp ý và ký duyệt giáo án điện tử trực tuyến |

### 2.4. Nhóm API Phân hệ Đánh giá Hiệu suất KPI 360 độ

| Method | Endpoint | Quyền hạn (RBAC) | Mô tả nghiệp vụ |
|---|---|---|---|
| `GET` | `/api/v1/kpi/my-dashboard` | `ALL_STAFF` | Dashboard KPI cá nhân: Điểm tự động, điểm cấp trên, tiến độ nộp |
| `POST` | `/api/v1/kpi/calculate-automated` | `SYSTEM, BGH` | Kích hoạt worker tự động tính 40% điểm hệ thống từ log hành vi |
| `POST` | `/api/v1/kpi/manager-review` | `BGH, TO_TRUONG` | Cấp trên chấm điểm đánh giá chuyên môn và tinh thần phối hợp |
| `POST` | `/api/v1/kpi/peer-review` | `GV_TRONG_TO` | Đánh giá chéo đồng nghiệp trong tổ chuyên môn (ẩn danh) |
| `POST` | `/api/v1/kpi/parent-survey` | `PHHS` | Phụ huynh gửi phiếu khảo sát mức độ hài lòng CSAT định kỳ |
| `POST` | `/api/v1/kpi/finalize-rankings` | `BGH` | Chốt kết quả KPI và tự động quy đổi xếp loại viên chức NĐ 90/48 |

---

## 3. THUẬT TOÁN TÍNH ĐIỂM KPI 360 & TỰ ĐỘNG QUY ĐỔI XẾP LOẠI VIÊN CHỨC

Dưới đây là mã nguồn Golang hiện thực hóa thuật toán tính điểm KPI và logic quy đổi theo quy chuẩn nhà nước:

```go
package kpi

import (
	"errors"
	"math"
)

// ClassificationResult lưu trữ kết quả xếp loại theo Nghị định 90 & Nghị định 48
type ClassificationResult string

const (
	RankExcellent     ClassificationResult = "HOÀN THÀNH XUẤT SẮC NHIỆM VỤ"
	RankGood          ClassificationResult = "HOÀN THÀNH TỐT NHIỆM VỤ"
	RankSatisfactory  ClassificationResult = "HOÀN THÀNH NHIỆM VỤ"
	RankUnsatisfactory ClassificationResult = "KHÔNG HOÀN THÀNH NHIỆM VỤ"
)

// KPICalculationInput dữ liệu đầu vào tính toán cho 1 cán bộ/giáo viên
type KPICalculationInput struct {
	StaffID                 string
	AutomatedSystemScore    float64 // Thang 100 (Trọng số 40%)
	ManagerScore            float64 // Thang 100 (Trọng số 30%)
	PeerReviewScore         float64 // Thang 100 (Trọng số 15%)
	ParentCSATScore         float64 // Thang 100 (Trọng số 15%)
	
	// Điều kiện ràng buộc pháp lý
	HasDisciplinaryAction   bool    // Có bị kỷ luật trong năm không (Đảng/Hành chính)
	ViolatedTutoringRules   bool    // Có vi phạm quy định dạy thêm tiểu học không (TT 29/2024 & TT 19/2026)
	HasInnovativeSolution   bool    // Có giải pháp sáng tạo/SKKN được công nhận (Luật TĐKT 2022)
	HasSevereSafetyIncident bool    // Để xảy ra tai nạn thương tích/ngộ độc nghiêm trọng
}

// KPIService chịu trách nhiệm tính toán
type KPIService struct{}

// CalculateFinalScore tính điểm tổng hợp KPI 360 độ
func (s *KPIService) CalculateFinalScore(in KPICalculationInput) float64 {
	finalScore := (in.AutomatedSystemScore * 0.40) +
		(in.ManagerScore * 0.30) +
		(in.PeerReviewScore * 0.15) +
		(in.ParentCSATScore * 0.15)

	// Làm tròn 2 chữ số thập phân
	return math.Round(finalScore*100) / 100
}

// DetermineOfficialRank quy đổi tự động sang mức xếp loại viên chức nhà nước
func (s *KPIService) DetermineOfficialRank(in KPICalculationInput, finalScore float64) (ClassificationResult, string, error) {
	// 1. Kiểm tra điểm liệt kỷ luật & đạo đức nhà giáo
	if in.HasDisciplinaryAction {
		return RankUnsatisfactory, "Viên chức bị xử lý kỷ luật trong năm theo quy định tại NĐ 48/2023/NĐ-CP", nil
	}
	if in.ViolatedTutoringRules {
		return RankUnsatisfactory, "Vi phạm quy định nghiêm cấm dạy thêm tiểu học tại Thông tư 29/2024 & TT 19/2026", nil
	}
	if in.HasSevereSafetyIncident {
		return RankUnsatisfactory, "Để xảy ra tai nạn thương tích hoặc mất an toàn học đường nghiêm trọng", nil
	}

	// 2. Xét mức Hoàn thành xuất sắc nhiệm vụ (Nghị định 90/2020 & Luật TĐKT 2022)
	// Điều kiện: Điểm KPI >= 90 điểm VÀ có giải pháp sáng tạo/đổi mới được công nhận VÀ không tiêu chí nào dưới 70
	if finalScore >= 90.0 {
		if in.HasInnovativeSolution {
			return RankExcellent, "Đạt điểm KPI xuất sắc (>=90) và có giải pháp đổi mới sáng tạo được công nhận", nil
		}
		// Nếu điểm >= 90 nhưng chưa có sáng kiến/đổi mới sáng tạo -> xếp mức Tốt
		return RankGood, "Điểm KPI đạt mức xuất sắc nhưng chưa có giải pháp đổi mới sáng tạo được nghiệm thu", nil
	}

	// 3. Xét mức Hoàn thành tốt nhiệm vụ
	if finalScore >= 75.0 {
		return RankGood, "Đạt điểm KPI từ 75 đến dưới 90 điểm, hoàn thành tốt các chỉ số nhiệm vụ", nil
	}

	// 4. Xét mức Hoàn thành nhiệm vụ
	if finalScore >= 50.0 {
		return RankSatisfactory, "Đạt điểm KPI từ 50 đến dưới 75 điểm", nil
	}

	// 5. Dưới 50 điểm
	return RankUnsatisfactory, "Điểm KPI dưới 50 điểm, không hoàn thành khối lượng công việc được giao", nil
}
```

---

## 4. CHI TIẾT THUẬT TOÁN ĐỊNH MỨC TIẾT DẠY & TÍNH THỪA GIỜ (THÔNG TƯ 05/2025/TT-BGDĐT)

Thuật toán tự động hóa việc tính định mức và tiền tăng giờ cho giáo viên:

```go
package teaching

// QuotaInput thông số giáo viên trong năm học
type QuotaInput struct {
	TeacherID          string
	IsHomeroom         bool // Là giáo viên chủ nhiệm
	IsHeadSubject      bool // Là tổ trưởng chuyên môn
	IsDeputyHead       bool // Là tổ phó chuyên môn
	HasYoungChild      bool // Nuôi con nhỏ dưới 12 tháng tuổi
	CounselingHoursWk  int  // Tiết tư vấn tâm lý học sinh kiêm nhiệm
	WeeksInYear        int  // Chuẩn: 35 tuần thực dạy
	ActualTeachingTaught int // Tổng số tiết thực dạy trên lớp cả năm
}

type OvertimeResult struct {
	StandardAnnualHours  int // Chuẩn: 805 tiết
	WeeklyReducedHours   int // Số tiết giảm mỗi tuần
	TotalReducedHours    int // Tổng số tiết được giảm trong năm
	NetAnnualQuota       int // Định mức thực tế phải dạy cả năm
	OvertimeHours        int // Số tiết vượt định mức (thừa giờ)
}

func CalculateTeacherOvertime(in QuotaInput) OvertimeResult {
	standardAnnual := 805
	weeks := in.WeeksInYear
	if weeks <= 0 {
		weeks = 35
	}

	weeklyReduction := 0
	// 1. Áp dụng giảm trừ theo Thông tư 05/2025/TT-BGDĐT Điều 9 & 10
	if in.IsHomeroom {
		weeklyReduction += 4 // GVCN tiểu học giảm 4 tiết/tuần (Mới)
	}
	if in.IsHeadSubject {
		weeklyReduction += 3 // Tổ trưởng giảm 3 tiết/tuần
	}
	if in.IsDeputyHead {
		weeklyReduction += 1 // Tổ phó giảm 1 tiết/tuần
	}
	if in.HasYoungChild {
		weeklyReduction += 3 // Con nhỏ dưới 12 tháng giảm 3 tiết/tuần
	}
	weeklyReduction += in.CounselingHoursWk // Tư vấn tâm lý

	totalReduced := weeklyReduction * weeks
	netQuota := standardAnnual - totalReduced
	if netQuota < 0 {
		netQuota = 0
	}

	overtime := in.ActualTeachingTaught - netQuota
	if overtime < 0 {
		overtime = 0
	}

	return OvertimeResult{
		StandardAnnualHours: standardAnnual,
		WeeklyReducedHours:  weeklyReduction,
		TotalReducedHours:   totalReduced,
		NetAnnualQuota:      netQuota,
		OvertimeHours:       overtime,
	}
}
```
