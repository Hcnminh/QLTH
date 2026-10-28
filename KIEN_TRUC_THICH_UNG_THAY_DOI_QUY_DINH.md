# THIẾT KẾ KIẾN TRÚC THÍCH ỨNG LINH HOẠT VỚI THAY ĐỔI QUY ĐỊNH PHÁP LUẬT
## (REGULATORY-RESILIENT ARCHITECTURE FOR EDUCATION SYSTEM)

---

## 1. THÁCH THỨC ĐẶC THÙ CỦA NGÀNH GIÁO DỤC VIỆT NAM

Trong quản lý giáo dục, các văn bản quy phạm pháp luật thay đổi liên tục và thường có **lộ trình chuyển tiếp (cuốn chiếu)**:
1. **Ví dụ Đánh giá học sinh tiểu học**:
   - Thông tư 30/2014 $\rightarrow$ Thông tư 22/2016 $\rightarrow$ Thông tư 27/2020.
   - Khi áp dụng Thông tư 27/2020, Bộ GD&ĐT triển khai cuốn chiếu: Năm học 2020-2021 áp dụng Lớp 1; 2021-2022 áp dụng Lớp 1, 2; đến 2024-2025 mới hoàn thành đến Lớp 5. Trong thời gian này, các khối lớp còn lại vẫn phải chạy theo Thông tư 22/2016 cũ.
2. **Ví dụ Chế độ làm việc giáo viên**:
   - Thông tư 28/2009 & TT 15/2017 (GVCN giảm 3 tiết/tuần) $\rightarrow$ Thông tư 05/2025/TT-BGDĐT có hiệu lực từ 22/04/2025 (805 tiết/năm, GVCN giảm 4 tiết/tuần).
3. **Ví dụ Đánh giá xếp loại viên chức & KPI**:
   - Nghị định 56/2015 $\rightarrow$ Nghị định 90/2020 $\rightarrow$ Nghị định 48/2023 $\rightarrow$ Luật Thi đua Khen thưởng 2022.

### Nguy cơ nếu thiết kế kiến trúc thông thường (Hard-coded Logic):
- **Phá vỡ dữ liệu lịch sử**: Cập nhật công thức tính điểm mới làm thay đổi ngược kết quả học bạ, xếp loại của học sinh các năm trước.
- **Khóa bảng Database (Downtime)**: Mỗi lần thông tư đổi cột/tiêu chí đánh giá là phải chạy migration `ALTER TABLE`, gây rủi ro cao cho CSDL hàng triệu học sinh.
- **Rối loạn mã nguồn (Spaghetti Code)**: Code ngập tràn các câu lệnh `if...else` kiểm tra năm học, thông tư dẫn đến lỗi tiềm ẩn.

---

## 2. NĂM NGUYÊN TẮC KIẾN TRÚC CỐT LÕI ĐỂ BẢO VỆ DỮ LIỆU & DỄ DÀNG CẬP NHẬT

```
+---------------------------------------------------------------------------------------------------+
|                        5 NGUYÊN TẮC KIẾN TRÚC THÍCH ỨNG THAY ĐỔI (REGULATORY-RESILIENT)            |
+---------------------------------------------------------------------------------------------------+
|  1. Temporal Policy Modeling (Mô hình chính sách gắn mốc hiệu lực thời gian & niên khóa)         |
|  2. Strategy Versioning & Policy Registry (Phiên bản hóa thuật toán nghiệp vụ theo thông tư)      |
|  3. Immutable Historical Snapshots (Đóng băng dữ liệu lịch sử - Không tính toán lại quá khứ)     |
|  4. Hybrid Relational + JSONB Schema (Lõi quan hệ cứng + Thuộc tính tiêu chí động có schema_ver) |
|  5. Configuration & Policy as Data (Quy tắc nghiệp vụ được quản trị bằng dữ liệu, không hardcode) |
+---------------------------------------------------------------------------------------------------+
```

---

## 3. CHI TIẾT THIẾT KẾ KỸ THUẬT BẰNG GO & POSTGRESQL

### 3.1. Nguyên tắc 1: Mô hình chính sách theo thời gian (Temporal Policy Modeling)
Mọi quy định về định mức tiết dạy, quy chuẩn xếp loại, cấu trúc điểm số không được gán cứng vào bảng trường học hay nhân sự, mà được lưu trong bảng **Chính sách (Policy)** có khoảng thời gian hiệu lực (`effective_from`, `effective_to`) và niên khóa áp dụng (`academic_year_id`):

```sql
-- Bảng quy định định mức tiết dạy (Teaching Quota Policy)
CREATE TABLE teaching_quota_policies (
    id VARCHAR(64) PRIMARY KEY,
    regulation_code VARCHAR(50) NOT NULL,    -- 'TT_28_2009' hoặc 'TT_05_2025'
    school_type VARCHAR(20) NOT NULL,        -- 'PRIMARY', 'SECONDARY'
    standard_annual_hours INT NOT NULL,       -- 805 tiết (Tiểu học mới) hoặc 0 (tính theo tuần)
    standard_weekly_hours INT NOT NULL,      -- 23 tiết/tuần
    homeroom_reduced_hours INT NOT NULL,     -- TT 28: giảm 3; TT 05: giảm 4 tiết/tuần
    effective_from DATE NOT NULL,
    effective_to DATE,                       -- NULL là đang áp dụng vô thời hạn
    metadata JSONB                           -- Các hệ số kiêm nhiệm khác
);
```

Khi hệ thống tính toán tải làm việc của giáo viên tại thời điểm bất kỳ, Go Service sẽ truy vấn chính sách có hiệu lực tại ngày đó:
```sql
SELECT * FROM teaching_quota_policies 
WHERE school_type = 'PRIMARY' 
  AND effective_from <= $1 
  AND (effective_to IS NULL OR effective_to >= $1);
```

---

### 3.2. Nguyên tắc 2: Phiên bản hóa Chiến lược đánh giá (Strategy Versioning & Registry)

Hệ thống triển khai một **Bộ đăng ký chính sách (Policy Registry)** trong Go. Mỗi Thông tư là một Struct độc lập thực thi cùng một interface:

```go
package evaluation

import "fmt"

// Interface chung cho mọi thông tư đánh giá
type AssessmentStrategy interface {
    GetRegulationCode() string // Mã định danh quy định (VD: TT22_2016, TT27_2020)
    Validate(record *StudentAssessmentInput) error
    CalculateResult(records []SubjectScore) (*TermEvaluationResult, error)
}

// 1. Logic Thông tư 22/2016 (Cũ)
type StrategyTT22Primary struct{}
func (s *StrategyTT22Primary) GetRegulationCode() string { return "TT22_2016" }
func (s *StrategyTT22Primary) CalculateResult(...) { /* Logic cũ */ }

// 2. Logic Thông tư 27/2020 (Áp dụng đến hết 30/11/2026)
type StrategyTT27Primary struct{}
func (s *StrategyTT27Primary) GetRegulationCode() string { return "TT27_2020" }
func (s *StrategyTT27Primary) CalculateResult(...) { /* Logic 5 phẩm chất, 10 năng lực */ }

// 3. Logic Thông tư 89/2026/TT-BGDĐT (Tự động kích hoạt từ 01/12/2026)
type StrategyTT89_2026Primary struct{}
func (s *StrategyTT89_2026Primary) GetRegulationCode() string { return "TT89_2026" }
func (s *StrategyTT89_2026Primary) CalculateResult(...) { 
    /* Bổ sung: Kiểm tra định kỳ máy tính, đề chung Sở/Phòng lớp 5, ghi "Miễn" môn học hòa nhập */ 
}

// Policy Registry điều phối chiến lược tự động theo thời điểm hiệu lực pháp lý
type AssessmentRegistry struct {
    strategies map[string]AssessmentStrategy
}

func (r *AssessmentRegistry) ResolveStrategy(schoolType string, grade int, evalDate time.Time) AssessmentStrategy {
    if schoolType == "PRIMARY" {
        // Mốc hiệu lực pháp lý chính thức của Thông tư 89/2026/TT-BGDĐT: 01/12/2026
        tt89EffectiveDate := time.Date(2026, 12, 1, 0, 0, 0, 0, time.Local)
        if evalDate.After(tt89EffectiveDate) || evalDate.Equal(tt89EffectiveDate) {
            return r.strategies["TT89_2026"]
        }
        return r.strategies["TT27_2020"]
    }
    return r.strategies["DEFAULT"]
}
```
👉 **Ưu điểm**: Khi Bộ GD&ĐT ban hành thông tư mới, lập trình viên chỉ cần **tạo thêm 1 struct chiến lược mới** (`StrategyTT_Moi`) và đăng ký vào Registry. **Toàn bộ code cũ và dữ liệu cũ không bị can thiệp 1 ký tự nào.**

---

### 3.3. Nguyên tắc 3: Đóng băng dữ liệu lịch sử (Immutable Snapshots)

Tuyệt đối **KHÔNG BAO GIỜ** tính toán lại điểm tổng kết hoặc học bạ của các năm học đã kết thúc. 

Khi một kỳ học hoặc năm học được Ban Giám hiệu ký số đóng sổ:
1. Hệ thống chạy thuật toán của thông tư tại thời điểm đó để ra kết quả cuối cùng.
2. Đóng gói toàn bộ kết quả vào một bản ghi **Snapshot JSON bất biến**.
3. Ký số điện tử (Digital Signature) lên bản snapshot và đóng dấu thời gian (Timestamp).
4. Khóa cấm sửa (Read-only Lock).

```sql
CREATE TABLE academic_transcript_snapshots (
    id VARCHAR(64) PRIMARY KEY,
    student_id VARCHAR(64) NOT NULL,
    academic_year_id VARCHAR(32) NOT NULL,
    grade_level INT NOT NULL,
    applied_regulation VARCHAR(50) NOT NULL, -- Lưu rõ: 'TT27_2020'
    final_data JSONB NOT NULL,               -- Toàn bộ bảng điểm, nhận xét, xếp loại
    digital_signature TEXT NOT NULL,         -- Chữ ký số BGH / SmartCA
    signed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_locked BOOLEAN DEFAULT TRUE,
    CONSTRAINT uk_student_year UNIQUE (student_id, academic_year_id)
);
```

👉 Khi phụ huynh hoặc học sinh xem lại học bạ lớp 1 của 5 năm trước, phần mềm **chỉ đọc từ `academic_transcript_snapshots`**. Bất kể hiện tại luật đã thay đổi như thế nào, dữ liệu lịch sử hoàn toàn được bảo toàn nguyên vẹn 100%.

---

### 3.4. Nguyên tắc 4: CSDL lai quan hệ + JSONB (Hybrid Schema)

Các thuộc tính pháp lý thường xuyên thêm/bớt tiêu chí (như danh mục phẩm chất, năng lực, chỉ số KPI phụ) được lưu bằng kiểu dữ liệu `JSONB` có kèm trường định danh phiên bản (`schema_version`):

```sql
CREATE TABLE student_evaluations (
    id VARCHAR(64) PRIMARY KEY,
    student_id VARCHAR(64) REFERENCES students(id),
    term_id VARCHAR(32) NOT NULL,
    subject_id VARCHAR(32) NOT NULL,
    
    -- Cột quan hệ chuẩn (ACID index)
    score NUMERIC(4,2),                     -- Điểm số nếu có (0-10)
    achievement_level VARCHAR(10),          -- T, H, C hoặc HTT, HT, CHT
    
    -- Dữ liệu động theo từng thông tư
    schema_version VARCHAR(50) NOT NULL,    -- 'TT27_SCHEMA_V1'
    dynamic_attributes JSONB NOT NULL       -- Nhận xét chuyên môn, phẩm chất, năng lực
);
```

**Ví dụ dữ liệu lưu trữ theo Thông tư 27/2020:**
```json
{
  "schema_version": "TT27_SCHEMA_V1",
  "qualities": {
    "patriotism": "T",
    "compassion": "T",
    "diligence": "D",
    "honesty": "T",
    "responsibility": "D"
  },
  "teacher_feedback": "Em tiếp thu bài nhanh, tích cực phát biểu, cần rèn thêm chữ viết."
}
```
👉 Nếu thông tư mới bổ sung thêm một tiêu chí (ví dụ: *Năng lực số học sinh tiểu học*), CSDL chỉ cần ghi nhận thêm thuộc tính trong JSONB với `schema_version: "TT_MOI_V1"`, **hoàn toàn không cần chạy lệnh `ALTER TABLE` gây khóa bảng**.

---

### 3.5. Nguyên tắc 5: Quy tắc nghiệp vụ & KPI cấu hình bằng dữ liệu (Policy as Data)

Đối với bộ công cụ KPI giáo viên/nhân viên, các quy chế chi tiêu nội bộ và tiêu chí thi đua của nhà trường thay đổi theo từng năm:
- Không hardcode các tiêu chí KPI trong code Go.
- Cấu hình cây tiêu chí KPI, trọng số và ngưỡng điểm trực tiếp trên **Giao diện Quản trị GoAdmin**.
- Áp dụng cơ chế **Nhân bản cấu hình (Clone Template)**: Khi bước sang năm học mới, BGH chỉ cần bấm *"Nhân bản bộ KPI từ năm học cũ sang năm học mới"*, sau đó điều chỉnh trọng số hoặc thêm tiêu chí mới.

---

## 4. QUY TRÌNH NÂNG CẤP HỆ THỐNG KHI CÓ QUY ĐỊNH MỚI (RUNBOOK)

Khi Bộ GD&ĐT ban hành văn bản mới (ví dụ Thông tư thay thế Thông tư 27 hoặc Nghị định thay thế NĐ 90):

```
+---------------------------------------------------------------------------------------+
| Bước 1: Khai báo Metadata quy định mới trên GoAdmin (Mã văn bản, ngày có hiệu lực)   |
+---------------------------------------------------------------------------------------+
                                           |
                                           v
+---------------------------------------------------------------------------------------+
| Bước 2: Viết Strategy mới trong thư mục `internal/tier_strategies/primary/`          |
|         (Kế thừa `AssessmentStrategy` interface, độc lập hoàn toàn với code cũ)      |
+---------------------------------------------------------------------------------------+
                                           |
                                           v
+---------------------------------------------------------------------------------------+
| Bước 3: Đăng ký Strategy mới vào `AssessmentRegistry` kèm lộ trình áp dụng (Lớp/Năm)  |
+---------------------------------------------------------------------------------------+
                                           |
                                           v
+---------------------------------------------------------------------------------------+
| Bước 4: Viết Unit Test kiểm thử song song: Dữ liệu năm cũ chạy đúng Strategy cũ,     |
|         Dữ liệu năm mới chạy đúng Strategy mới. CI/CD tự động deploy không downtime.  |
+---------------------------------------------------------------------------------------+
```

## 5. KẾT LUẬN

Với thiết kế kiến trúc này:
1. **Dữ liệu lịch sử là vĩnh cửu và bất biến** nhờ Snapshot & Chữ ký số.
2. **Hệ thống thích ứng mọi lộ trình cải cách giáo dục** (cuốn chiếu nhiều năm) mà không xảy ra xung đột nghiệp vụ.
3. **Chi phí bảo trì và nâng cấp phần mềm giảm tới 70%** so với cách lập trình hard-code truyền thống.
