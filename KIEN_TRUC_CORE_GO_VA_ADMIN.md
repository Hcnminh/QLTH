# THIẾT KẾ KIẾN TRÚC BỘ CORE HỆ THỐNG GIÁO DỤC BẰNG GO (EDU-CORE GO)
## TÍCH HỢP GIAO DIỆN QUẢN TRỊ ADMIN TỰ ĐỘNG & TÁI SỬ DỤNG ĐA CẤP HỌC

---

## 1. TỔNG QUAN KIẾN TRÚC HỆ THỐNG BẰNG GO

Hệ thống được thiết kế theo nguyên lý **Domain-Driven Design (DDD)** và **Clean Architecture (Hexagonal Architecture)**. 
Bộ khung **Edu-Core** tách biệt hoàn toàn giữa **Phần lõi dùng chung (Core Engine)** và **Nghiệp vụ đặc thù từng cấp học (Tier Strategy Plugins)**, cho phép tái sử dụng cho:
- **Mầm non** (Thông tư 51/2020/TT-BGDĐT)
- **Tiểu học** (Thông tư 28/2020/TT-BGDĐT, Thông tư 27/2020/TT-BGDĐT, Thông tư 05/2025/TT-BGDĐT)
- **Trung học cơ sở & Trung học phổ thông** (Thông tư 32/2020/TT-BGDĐT, Thông tư 22/2021/TT-BGDĐT, Thông tư 05/2025/TT-BGDĐT)
- **Trường liên cấp (K-12)**

```
+-----------------------------------------------------------------------------------------------+
|                                      LỚP TRÌNH DIỄN (PRESENTATION)                           |
|  +---------------------------+  +---------------------------+  +---------------------------+  |
|  |   GO AUTO-ADMIN PORTAL    |  |  PORTAL GIÁO VIÊN & BGH   |  |   MOBILE APP (FLUTTER)    |  |
|  | (Back-office Quản trị hệ  |  |  (React/Next.js SPA)      |  |   (Phụ huynh & Giáo viên) |  |
|  |  thống, cấu hình dữ liệu) |  |                           |  |                           |  |
|  +---------------------------+  +---------------------------+  +---------------------------+  |
+-----------------------------------------------------------------------------------------------+
                                               |
+-----------------------------------------------------------------------------------------------+
|                                    API GATEWAY & ROUTER (GIN / FIBER)                         |
|                   Middleware: Auth (JWT/OAuth2), RBAC/ABAC, Rate Limiter, Audit Log           |
+-----------------------------------------------------------------------------------------------+
                                               |
+-----------------------------------------------------------------------------------------------+
|                                   EDU-CORE FOUNDATION (GO ENGINE)                             |
|  +---------------------+  +---------------------+  +---------------------+  +---------------+ |
|  | 1. Multi-Tenancy    |  | 2. IAM & RBAC Engine|  | 3. Organization     |  | 4. Person &   | |
|  |    (Trường/Phòng/Sở)|  |    (Phân quyền động)|  |    (Trường/Khối/Lớp)|  |    Hồ sơ gốc   | |
|  +---------------------+  +---------------------+  +---------------------+  +---------------+ |
|  +---------------------+  +---------------------+  +---------------------+  +---------------+ |
|  | 5. KPI & Performance|  | 6. Finance & Billing|  | 7. Boarding & Health|  | 8. Audit Log &| |
|  |    (Đo lường 360)   |  |    (VietQR Engine)  |  |    (Bán trú, Y tế)  |  |    Ký số PKI  | |
|  +---------------------+  +---------------------+  +---------------------+  +---------------+ |
+-----------------------------------------------------------------------------------------------+
                                               |
+-----------------------------------------------------------------------------------------------+
|                          TIER STRATEGY PLUGINS (CHIẾN LƯỢC THEO CẤP HỌC)                      |
|  +-----------------------+  +------------------------+  +-----------------------------------+ |
|  | Strategy: TIỂU HỌC    |  | Strategy: THCS & THPT  |  | Strategy: MẦM NON                 | |
|  | - Đánh giá TT 27 &    |  | - Đánh giá TT 22 &     |  | - Đánh giá sự phát triển 5 lĩnh   | |
|  |   TT 89/2026/TT-BGDĐT  |  |   TT 89/2026/TT-BGDĐT   |  |   vực theo TT 51/2020/TT-BGDĐT    | |
|  | - Định mức 805 tiết/n |  | - Điểm số TB môn (0-10)|  | - Quản lý chăm sóc, cân đo dinh   | |
|  | - Giảm 4 tiết/tuần GVCN|  | - Định mức 665/595 tiết|  |   dưỡng định kỳ                   | |
|  |   (TT 05/2025/TT-BGDĐT)|  |   (TT 05/2025/TT-BGDĐT)|  |                                   | |
|  +-----------------------+  +------------------------+  +-----------------------------------+ |
+-----------------------------------------------------------------------------------------------+
                                               |
+-----------------------------------------------------------------------------------------------+
|                                  LỚP HẠ TẦNG & DỮ LIỆU (INFRASTRUCTURE)                       |
|   PostgreSQL (pgx / GORM)  |  Redis (Cache & Distributed Lock)  |  MinIO/S3 (Hồ sơ, Học bạ)   |
+-----------------------------------------------------------------------------------------------+
```

---

## 2. GIẢI PHÁP GIAO DIỆN QUẢN TRỊ ADMIN (GO ADMIN ENGINE)

Để có được sức mạnh như **Django Admin** trong hệ sinh thái **Go**, chúng tôi đề xuất 2 phương án tối ưu:

### Phương án A: Tích hợp GoAdmin Engine trực tiếp trong Go binary (Khuyên dùng)
- **Thư viện**: Sử dụng **GoAdmin (`github.com/GoAdminGroup/go-admin`)** kết hợp theme chuẩn AdminLTE hoặc Sword.
- **Cơ chế hoạt động**:
  - Tự động sinh bảng dữ liệu (Data Table), form nhập liệu (Create/Edit), bộ lọc nâng cao (Filters), tìm kiếm (Search), xuất Excel từ struct mô hình dữ liệu của Go.
  - Hỗ trợ Role-based Access Control (RBAC) nội bộ cho phép phân cấp quản trị viên hệ thống (Super Admin, IT trường, Trưởng phòng giáo dục).
  - Biên dịch toàn bộ CSS/JS và templates HTML vào thẳng **một file Binary Go duy nhất**, triển khai siêu nhẹ (RAM chỉ tốn $\sim 30\text{ MB}$).

### Phương án B: Headless Meta-Admin (Go Generic API + React Admin / Refine)
- Backend Go cung cấp các `Generic REST/GraphQL Endpoints` trả về siêu dữ liệu (Metadata Schema: danh sách cột, kiểu dữ liệu, quan hệ khóa ngoại).
- Frontend quản trị dùng framework **Refine.dev** hoặc **React-Admin** tự động dựng UI từ JSON Schema.
- Thích hợp khi muốn giao diện quản trị cực kỳ hiện đại, đồng bộ thiết kế với Portal giáo viên.

---

## 3. THIẾT KẾ CẤU TRÚC BỘ CORE TÁI SỬ DỤNG (EDU-CORE DOMAIN)

### 3.1. Phân loại Cấp học và Cấu hình Đa trường (Multi-Tenancy)

```go
package domain

type SchoolType string

const (
    SchoolTypeKindergarten SchoolType = "KINDERGARTEN" // Mầm non
    SchoolTypePrimary      SchoolType = "PRIMARY"      // Tiểu học
    SchoolTypeSecondary    SchoolType = "SECONDARY"    // THCS
    SchoolTypeHighSchool   SchoolType = "HIGH_SCHOOL"  // THPT
    SchoolTypeInterLevel   SchoolType = "INTER_LEVEL"  // Liên cấp (K-12)
)

// Trường học (Tenant)
type School struct {
    ID          string     `json:"id" gorm:"primaryKey"`
    Code        string     `json:"code" gorm:"uniqueIndex"`      // Mã trường Bộ GD&ĐT
    Name        string     `json:"name"`
    Type        SchoolType `json:"type"`                         // Cấp học
    Address     string     `json:"address"`
    ProvinceID  string     `json:"province_id"`
    DistrictID  string     `json:"district_id"`
    IsPublic    bool       `json:"is_public"`                    // Công lập hay Tư thục
    CreatedAt   int64      `json:"created_at"`
}
```

---

### 3.2. Cấu trúc Học sinh & Cán bộ nhân viên cốt lõi

```go
// Hồ sơ gốc học sinh (Dùng chung cho mọi cấp)
type Student struct {
    ID            string    `json:"id" gorm:"primaryKey"`
    SchoolID      string    `json:"school_id" gorm:"index"`
    NationalID    string    `json:"national_id" gorm:"index"`   // Số định danh VNeID
    FullName      string    `json:"full_name"`
    BirthDate     string    `json:"birth_date"`
    Gender        string    `json:"gender"`
    CurrentClassID string   `json:"current_class_id" gorm:"index"`
    Status        string    `json:"status"`                     // Đang học, Chuyển trường, Tốt nghiệp
    Metadata      JSONB     `json:"metadata"`                   // Thuộc tính mở rộng theo cấp học
}

// Cán bộ giáo viên & Nhân viên (Vị trí việc làm TT 20/2023)
type Staff struct {
    ID            string    `json:"id" gorm:"primaryKey"`
    SchoolID      string    `json:"school_id" gorm:"index"`
    EmployeeCode  string    `json:"employee_code" gorm:"uniqueIndex"`
    FullName      string    `json:"full_name"`
    JobTitle      string    `json:"job_title"`                  // GVCN, GV Bộ môn, Hiệu trưởng, Kế toán, Y tế...
    StandardHours int       `json:"standard_hours"`             // Định mức tiết: 805 (Tiểu học), 665 (THCS)...
    ReducedHours  int       `json:"reduced_hours"`              // Tiết giảm trừ (GVCN - 4 tiết, Tổ trưởng...)
    IsHomeroom    bool      `json:"is_homeroom"`                // Có chủ nhiệm lớp không
}
```

---

### 3.3. Core Đánh giá học tập đa hình (Polymorphic Assessment Engine)

Thay vì viết cứng công thức chấm điểm, Core Go định nghĩa một **Interface chiến lược đánh giá (Strategy Pattern)**:

```go
package evaluation

// Giao diện đánh giá dùng chung cho mọi cấp học
type AssessmentStrategy interface {
    GetAssessmentType() string
    ValidateScore(input ScoreInput) error
    EvaluatePeriod(studentID string, termID string) (*AssessmentResult, error)
    GenerateReportCard(studentID string, yearID string) (*ReportCardData, error)
}

// Chiến lược 1: Dành cho Tiểu học (Thông tư 27/2020/TT-BGDĐT)
type PrimarySchoolStrategy struct{}

func (s *PrimarySchoolStrategy) GetAssessmentType() string {
    return "CIRCULAR_27_PRIMARY"
}

func (s *PrimarySchoolStrategy) ValidateScore(input ScoreInput) error {
    // Kiểm tra: Chỉ một số môn được chấm điểm 1-10; còn lại chỉ nhận xét và đánh giá mức Đạt (T/H/C)
    // Validate 5 phẩm chất và 10 năng lực
    return nil
}

// Chiến lược 2: Dành cho THCS & THPT (Thông tư 22/2021/TT-BGDĐT)
type SecondarySchoolStrategy struct{}

func (s *SecondarySchoolStrategy) GetAssessmentType() string {
    return "CIRCULAR_22_SECONDARY"
}

func (s *SecondarySchoolStrategy) ValidateScore(input ScoreInput) error {
    // Kiểm tra: Điểm đánh giá thường xuyên (ĐĐGtx), giữa kỳ (ĐĐGgk), cuối kỳ (ĐĐGck)
    // Tính điểm trung bình môn (ĐTBmhk), xếp loại: Tốt, Khá, Đạt, Chưa đạt
    return nil
}

// Factory khởi tạo chiến lược đánh giá tương ứng với loại trường
func NewAssessmentEngine(schoolType domain.SchoolType) AssessmentStrategy {
    switch schoolType {
    case domain.SchoolTypePrimary:
        return &PrimarySchoolStrategy{}
    case domain.SchoolTypeSecondary, domain.SchoolTypeHighSchool:
        return &SecondarySchoolStrategy{}
    default:
        return &PrimarySchoolStrategy{}
    }
}
```

---

### 3.4. Core Đánh giá hiệu suất & KPI Engine dùng chung

Phân hệ KPI được thiết kế tổng quát hóa để áp dụng cho mọi loại trường (Tiểu học, THCS, THPT, Mầm non, Liên cấp):

```go
package kpi

// Danh mục chỉ số KPI
type KPICriterion struct {
    ID            string   `json:"id" gorm:"primaryKey"`
    SchoolID      string   `json:"school_id" gorm:"index"`
    RoleApply     string   `json:"role_apply"`     // HOMEROOM_TEACHER, SUBJECT_TEACHER, PRINCIPAL, ACCOUNTANT...
    Code          string   `json:"code"`           // Ví dụ: KPI-GV-01
    Title         string   `json:"title"`
    Weight        float64  `json:"weight"`         // Trọng số (ví dụ: 0.20 cho 20%)
    MetricSource  string   `json:"metric_source"`  // SYSTEM_AUTOMATED, MANAGER_REVIEW, PEER_REVIEW, PARENT_CSAT
    Formula       string   `json:"formula"`        // Công thức tính điểm tự động
}

// Kỳ đánh giá KPI
type KPIPeriod struct {
    ID          string `json:"id" gorm:"primaryKey"`
    SchoolID    string `json:"school_id"`
    Title       string `json:"title"`              // Ví dụ: Đánh giá KPI Tháng 10/2026
    Type        string `json:"type"`               // MONTHLY, SEMESTER, ANNUAL
    StartDate   int64  `json:"start_date"`
    EndDate     int64  `json:"end_date"`
    IsClosed    bool   `json:"is_closed"`
}

// Bảng kết quả tổng hợp KPI cá nhân
type KPIAssessmentRecord struct {
    ID            string  `json:"id" gorm:"primaryKey"`
    PeriodID      string  `json:"period_id" gorm:"index"`
    StaffID       string  `json:"staff_id" gorm:"index"`
    AutomatedScore float64 `json:"automated_score"` // Điểm hệ thống tính (40%)
    ManagerScore   float64 `json:"manager_score"`   // Điểm cấp trên chấm (30%)
    PeerScore      float64 `json:"peer_score"`      // Điểm đồng nghiệp (15%)
    ParentCSAT     float64 `json:"parent_csat"`     // Khảo sát phụ huynh (15%)
    FinalScore     float64 `json:"final_score"`     // Tổng điểm có trọng số
    RankClassification string `json:"rank_classification"` // Xuất sắc, Tốt, Hoàn thành, Không hoàn thành (NĐ 90 & NĐ 48)
}
```

---

## 4. CẤU TRÚC THƯ MỤC DỰ ÁN CHUẨN GO (CLEAN ARCHITECTURE)

```
QLTH/
├── cmd/
│   ├── api/                 # Entrypoint API Server chạy Fiber hoặc Gin
│   │   └── main.go
│   └── admin/               # Entrypoint Web Admin (GoAdmin Engine)
│       └── main.go
├── internal/
│   ├── core/                # BỘ LÕI DÙNG CHUNG (TÁI SỬ DỤNG CHO MỌI CẤP HỌC)
│   │   ├── auth/            # JWT, Casbin (RBAC/ABAC phân quyền ma trận)
│   │   ├── organization/    # Quản lý Trường, Cơ sở, Khối, Lớp, Năm học
│   │   ├── person/          # Học sinh, Giáo viên, Nhân viên, Phụ huynh
│   │   ├── kpi/             # Thuật toán KPI Engine, chấm điểm 360, xếp loại NĐ 90/48
│   │   ├── billing/         # Khoản thu, Cổng thanh toán VietQR, Gạch nợ tự động
│   │   ├── boarding/        # Quản lý bán trú, dinh dưỡng, Kcal
│   │   ├── health/          # Y tế học đường, khám sức khỏe, tiêm chủng
│   │   └── digital_sign/    # Ký số học bạ (PKCS#7, SmartCA, Timestamp)
│   │
│   ├── tier_strategies/     # CÁC CHIẾN LƯỢC ĐẶC THÙ THEO TỪNG CẤP HỌC
│   │   ├── primary/         # TIỂU HỌC: TT 27/2020 (Đánh giá), TT 05/2025 (805 tiết/năm, giảm 4 tiết GVCN)
│   │   ├── secondary/       # THCS & THPT: TT 22/2021 (Tính ĐTBmhk), TT 05/2025 (665/595 tiết/năm)
│   │   └── kindergarten/    # MẦM NON: Đánh giá thể chất & phát triển trẻ TT 51/2020
│   │
│   ├── admin/               # CẤU HÌNH GIAO DIỆN QUẢN TRỊ GOADMIN
│   │   ├── tables/          # Định nghĩa giao diện CRUD cho các bảng dữ liệu
│   │   └── pages/           # Dashboard BGH, Quản lý danh mục dùng chung
│   │
│   └── platform/            # KẾT NỐI HẠ TẦNG
│       ├── database/        # PostgreSQL Connection Pool (pgx/gorm)
│       ├── cache/           # Redis Cache & Distributed Lock
│       └── storage/         # MinIO / AWS S3 Upload file giáo án, học bạ
│
├── go.mod
├── go.sum
└── Dockerfile
```

---

## 5. ĐÁNH GIÁ LỢI THẾ CỦA GIẢI PHÁP NÀY

1. **Hiệu năng & Tài nguyên**:
   - Biên dịch thành một file nhị phân duy nhất (`single binary`), không phụ thuộc runtime bên ngoài.
   - Chạy container tốn cực ít RAM ($\sim 30 - 60\text{ MB}$), chịu tải hàng chục nghìn kết nối đồng thời với CPU tối thiểu.
2. **Khả năng mở rộng không giới hạn**:
   - Khi nhà trường mở rộng thêm cấp THCS hay THPT (trường liên cấp), chỉ cần kích hoạt thêm `SecondarySchoolStrategy`, không cần đập đi xây lại cấu trúc Database hay Core Auth/KPI.
3. **Quản trị dễ dàng**:
   - Nhờ tích hợp **GoAdmin**, toàn bộ việc cấu hình tham số, danh mục, phân quyền cán bộ, kiểm tra log hệ thống được thực hiện trên giao diện Web trực quan như Django Admin.
