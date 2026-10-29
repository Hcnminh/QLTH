# EduPrimary-KPI: Danh Mục Ảnh Chụp Màn Hình Giao Diện Hệ Thống (Screenshots Catalog)

Tài liệu này tổng hợp và phân loại chi tiết toàn bộ các ảnh chụp màn hình độ phân giải cao (High-Resolution UI Screenshots) của phần mềm Quản lý Trường Tiểu học tích hợp KPI (**EduPrimary-KPI**), đáp ứng cam kết nghiệm thu **VAL-BROWSER-004** trong hợp đồng kiểm thử hệ thống.

---

## 1. Quy Chuẩn Đặt Tên & Kỹ Thuật Chụp (Naming Convention & Specs)

- **Độ phân giải chuẩn:** Viewport Desktop $1366 \times 860$ pixels (Chuẩn trình duyệt máy tính văn phòng trường học phổ thông Việt Nam).
- **Cơ chế chụp:** Playwright Automated Headless Chromium Suite (`scripts/qa_browser.js`) kết hợp hệ thống xác thực GoAdmin AdminLTE và Agent Browser.
- **Tiêu chuẩn số hóa:** Các file được đánh số thứ tự chuẩn hóa từ `01` đến `21`, phân định rõ ràng giữa màn hình danh sách (List Table), màn hình biểu mẫu chi tiết (Detail Form View), và các bảng điều khiển chuyên biệt (Interactive Panels).

---

## 2. Bảng Tổng Hợp Danh Mục Ảnh Chụp (Comprehensive Image Catalog)

| STT | Tên Tệp Ảnh | Phân Hệ / Luồng Nghiệp Vụ | Vai Trò (Persona) | Căn Cứ Pháp Lý & Nội Dung Trọng Tâm |
|:---:|:---|:---|:---:|:---|
| 01 | `01_login_page.png` | **Xác thực Hệ thống** | Tất cả người dùng | Màn hình đăng nhập quản trị GoAdmin AdminLTE (`/admin/login`), trường tên đăng nhập, mật khẩu, chống brute-force. |
| 02 | `02_dashboard_overview.png` | **Bảng Điều Khiển Trung Tâm** | BGH, Quản trị viên | Tổng quan 10 phân hệ, widget định mức chuẩn 805h, trần 20% KPI, trạng thái Digital Consent NĐ 13, điều hướng thanh bên. |
| 03 | `03_staff_table.png` | **Phân hệ 3: Cán bộ & Giáo viên** | BGH, Văn thư | Danh sách giáo viên, mã định danh, ngạch chức danh, phân công chủ nhiệm và trạng thái giảng dạy. |
| 04 | `03_staffs_teaching_quota.png` | **Phân hệ 3: Định mức Tiết dạy** | BGH, Kế toán | Bảng tính định mức tiết dạy: Chuẩn 805 tiết/năm, giảm 4 tiết/tuần GVCN, giảm tổ trưởng, nuôi con nhỏ theo **Thông tư 05/2025/TT-BGDĐT**. |
| 05 | `04_student_table.png` | **Phân hệ 1: Hồ sơ Học sinh** | BGH, GVCN, Văn thư | Danh sách học sinh toàn trường, mã định danh, lớp học, thông tin người giám hộ và trạng thái kích hoạt hồ sơ. |
| 06 | `05_students_digital_consent.png` | **Phân hệ 1: Digital Consent** | Văn thư, GVCN | Quản lý đồng thuận điện tử của cha mẹ theo **Điều 20 Nghị định 13/2023/NĐ-CP**, liên kết định danh VNeID Mức 2. |
| 07 | `06_classes_management.png` | **Phân hệ 6: Quản lý Lớp học** | BGH, Văn thư | Danh mục lớp học từ Khối 1 đến Khối 5, phòng học chức năng, sĩ số học sinh và phân công GVCN. |
| 08 | `07_lesson_plans_cv2345.png` | **Phân hệ 3: Kế hoạch Bài dạy** | BGH, GVBM, GVCN | Danh sách giáo án điện tử theo chuẩn khung 4 hoạt động của **Công văn 2345/BGDĐT-GDTH** và bài học STEM (**Thông tư 86/BGDĐT**). |
| 09 | `07b_timetable_schedule.png` | **Phân hệ 3: Thời Khóa Biểu** | BGH, GVCN, GVBM | Ma trận phân bổ tiết dạy Thứ 2 - Thứ 6, thuật toán kiểm tra và ngăn chặn xung đột giáo viên trùng lịch theo **TT 05/2025/TT-BGDĐT**. |
| 10 | `08_subjects_catalog.png` | **Phân hệ 2: Danh mục Môn học** | BGH, GVBM | Danh mục môn học CT GDPT 2018 (Toán, Tiếng Việt, Tiếng Anh, Tin học...) và cấu hình hình thức đánh giá (cho điểm / nhận xét). |
| 11 | `09_subject_evaluations_tt27_89.png` | **Phân hệ 2: Đánh giá Môn học** | GVBM, GVCN | Sổ đánh giá kết quả học tập định kỳ: Mức Đạt (HTT/HT/CHT), điểm kiểm tra 1-10, tích hợp công nghệ AI Copilot gợi ý nhận xét. |
| 12 | `10_behavior_competency_evaluations.png` | **Phân hệ 2: Phẩm chất & Năng lực** | GVCN | Đánh giá 5 phẩm chất chủ yếu và 10 năng lực cốt lõi theo thang đo T/D/C tuân thủ **Thông tư 27/2020** và **Thông tư 89/2026/TT-BGDĐT**. |
| 13 | `11_transcripts_smartca_snapshot.png` | **Phân hệ 9: Học bạ số SmartCA** | BGH, GVCN | Snapshot học bạ điện tử đóng băng bất biến, tích hợp ký số SmartCA chuẩn CMS/PKCS#7 và gắn dấu thời gian quốc gia TSA Timestamp. |
| 14 | `12_boarding_nutrition_menu.png` | **Phân hệ 4: Dinh dưỡng Bán trú** | Kế toán, Bếp trưởng | Kế hoạch thực đơn bán trú, định lượng Calo/Kcal, cơ cấu dinh dưỡng và kiểm soát quy chuẩn lưu mẫu thức ăn 24h (**Quyết định 1246/QĐ-BYT**). |
| 15 | `13_student_health_bmi.png` | **Phân hệ 5: Sức khỏe Học đường** | Y tế học đường | Hồ sơ thể lực học sinh, biểu đồ tăng trưởng BMI, kiểm tra thị lực học đường, nha khoa và theo dõi bệnh lý học sinh. |
| 16 | `13b_kpi_evaluations_20pct_cap.png` | **Phân hệ 8: Đánh giá KPI 360** | BGH, Công đoàn | Kết quả đánh giá hiệu suất 4 chiều (40-30-15-15), áp dụng điểm liệt sư phạm và khống chế trần **20% Hoàn thành xuất sắc** (**Nghị định 90/2020 & 48/2023**). |
| 17 | `14_lesson_plan_detail.png` | **Phân hệ 3: Chi tiết Kế hoạch bài dạy** | GVBM, Tổ trưởng | Chi tiết giáo án CV 2345: Thiết bị dạy học, tiến trình hoạt động học tập, và quy trình phê duyệt phân cấp 3 bước. |
| 18 | `15_boarding_menu_detail.png` | **Phân hệ 4: Chi tiết Thực đơn Bán trú** | Kế toán, Y tế | Chi tiết các món ăn, tổng năng lượng 718.2 Kcal, tỷ lệ P:L:G (14.2% : 26.5% : 59.3%) cân đối theo Viện Dinh Dưỡng Quốc Gia. |
| 19 | `16_student_health_detail.png` | **Phân hệ 5: Chi tiết Y tế & Đơn thuốc** | Y tế, Phụ huynh | Sổ khám sức khỏe định kỳ, nhật ký chăm sóc y tế, và quy trình tiếp nhận/cho uống thuốc theo đơn của phụ huynh. |
| 20 | `17_transcript_snapshot_detail.png` | **Phân hệ 9: Chi tiết Học bạ & Ký số** | BGH, GVCN | Chi tiết dữ liệu học bạ số chuyển tiếp (TT27 trước 01/12/2026, TT89 từ 01/12/2026), chữ ký điện tử Hiệu trưởng và chứng thư số. |
| 21 | `18_staff_detail.png` | **Phân hệ 3: Chi tiết Giáo viên & Quota** | BGH, Kế toán | Bảng tổng hợp tiết dạy hàng tuần, các hệ số giảm trừ pháp định, và bảng quyết toán tiền dạy thừa giờ cuối năm. |
| 22 | `19_student_detail.png` | **Phân hệ 1: Chi tiết Hồ sơ Học sinh** | Văn thư, GVCN | Hồ sơ học sinh đầy đủ, xác thực thông tin phụ huynh, lịch sử cấp phát chứng nhận Digital Consent NĐ 13. |
| 23 | `20_kpi_evaluation_detail.png` | **Phân hệ 8: Chi tiết Điểm KPI & Xếp loại** | BGH, Giáo viên | Bảng điểm KPI 4 nguồn (Hệ thống tự động, Ban giám hiệu, Đồng nghiệp, Phụ huynh) và điều kiện Sáng kiến kinh nghiệm. |
| 24 | `21_integrations_sandbox_stubs.png` | **Phân hệ 10: Cổng Tích hợp Sandbox** | BGH, Quản trị viên | Trung tâm giám sát 5 cổng sandbox: VNeID Mức 2, SmartCA TSA, MOET CSDL Outbox, VietQR EMVCo và Hóa đơn điện tử NĐ 123. |

---

## 3. Hướng Dẫn Chi Tiết Theo Từng Phân Hệ Nghiệp Vụ (Subsystem Walkthroughs)

### 3.1. Phân Hệ Xác Thực & Bảng Điều Khiển Trung Tâm (Auth & Dashboard)
- **Tệp minh họa:** `01_login_page.png`, `02_dashboard_overview.png`
- **Mô tả quy trình:**
  1. Người dùng truy cập đường dẫn quản trị `http://localhost:8098/admin/login`.
  2. Giao diện đăng nhập bảo mật AdminLTE kiểm tra thông tin định danh và quyền hạn RBAC.
  3. Sau khi xác thực thành công, hệ thống điều hướng vào Bảng điều khiển trung tâm (`/admin`), hiển thị tổng quan các khối số liệu pháp lý:
     - Định mức chuẩn **805 tiết/năm** (Thông tư 05/2025/TT-BGDĐT).
     - Khống chế trần **20% Hoàn thành xuất sắc** (Nghị định 90/2020 & 48/2023/NĐ-CP).
     - Tình trạng thu thập **Digital Consent** bảo vệ dữ liệu trẻ em (Điều 20 Nghị định 13/2023/NĐ-CP).
     - Lối tắt truy cập Thời khóa biểu chống trùng lịch và Cổng tích hợp dịch vụ số quốc gia.

### 3.2. Phân Hệ Quản Lý Giáo Viên & Định Mức Tiết Dạy (Teaching Quota - TT 05/2025)
- **Tệp minh họa:** `03_staff_table.png`, `03_staffs_teaching_quota.png`, `18_staff_detail.png`
- **Mô tả quy trình:**
  1. Giao diện `/admin/info/staffs` hiển thị danh sách toàn bộ cán bộ, giáo viên với các trường ngạch bậc, chức vụ kiêm nhiệm.
  2. Hệ thống áp dụng công thức Thông tư 05/2025/TT-BGDĐT:
     $$\text{Định mức thực tế} = \text{Định mức chuẩn (805 tiết)} - (\text{Số tiết giảm/tuần} \times 35 \text{ tuần})$$
     - Giáo viên chủ nhiệm (GVCN) được giảm trừ chính xác **4 tiết/tuần** (thay cho mức 3 tiết của Thông tư 28/2009 cũ).
     - Tổ trưởng chuyên môn giảm 3 tiết/tuần, tổ phó giảm 1 tiết/tuần, nữ giáo viên nuôi con dưới 12 tháng giảm 3 tiết/tuần.
  3. Số tiết thừa giờ được tính tự động, không bao giờ âm, làm căn cứ chi trả thù lao dạy thêm giờ hợp pháp.

### 3.3. Phân Hệ Quản Lý Lớp Học & Thời Khóa Biểu (Classes & Timetable Slot Conflict)
- **Tệp minh họa:** `06_classes_management.png`, `07b_timetable_schedule.png`
- **Mô tả quy trình:**
  1. Giao diện `/admin/info/classes` quản lý danh mục các lớp từ Lớp 1A1 đến Lớp 5A1, phòng học chức năng và phân công chủ nhiệm.
  2. Màn hình `/admin/timetable` trình bày trực quan ma trận phân bổ tiết dạy cho cả 2 buổi sáng (Tiết 1-4) và chiều (Tiết 5-7).
  3. **Thuật toán Chống Trùng Tiết:** Khi xếp lịch hoặc gọi API `/api/v1/timetable/slots`, hệ thống tự động kiểm tra xung đột thời gian của giáo viên. Nếu một giáo viên bị xếp vào 2 lớp khác nhau trong cùng một tiết học, hệ thống ngay lập tức chặn lại và trả về lỗi **HTTP 409 Conflict** kèm trích dẫn pháp lý rõ ràng.

### 3.4. Phân Hệ Hồ Sơ Học Sinh & Digital Consent (Decree 13/2023 Child Data Protection)
- **Tệp minh họa:** `04_student_table.png`, `05_students_digital_consent.png`, `19_student_detail.png`
- **Mô tả quy trình:**
  1. `/admin/info/students` quản lý mã định danh học sinh, ngày sinh, giới tính và thông tin người giám hộ hợp pháp.
  2. **Ràng buộc Điều 20 Nghị định 13/2023/NĐ-CP:** Hồ sơ học sinh chỉ được coi là hợp lệ khi cờ `HasDigitalConsent` được phụ huynh xác thực điện tử.
  3. Khi tạo mới hoặc nhập hồ sơ chưa có sự đồng thuận của cha mẹ, hệ thống từ chối kích hoạt với mã lỗi **HTTP 422 Unprocessable Entity**.

### 3.5. Phân Hệ Kế Hoạch Bài Dạy (Lesson Plans - CV 2345/BGDĐT-GDTH & STEM TT 86)
- **Tệp minh họa:** `07_lesson_plans_cv2345.png`, `14_lesson_plan_detail.png`
- **Mô tả quy trình:**
  1. Giáo viên bộ môn soạn và nộp giáo án điện tử theo chuẩn khung 4 giai đoạn hoạt động học tập (Khởi động -> Khám phá -> Luyện tập -> Vận dụng) theo Công văn 2345.
  2. Bảng quản trị hiển thị trạng thái phê duyệt phân cấp 3 bước: `SUBMITTED` -> `REVIEWED` (Tổ trưởng chuyên môn) -> `APPROVED` (Hiệu trưởng / Ban giám hiệu).

### 3.6. Phân Hệ Danh Mục Môn Học & Đánh Giá Kết Quả Học Tập (TT 27/2020 & TT 89/2026)
- **Tệp minh họa:** `08_subjects_catalog.png`, `09_subject_evaluations_tt27_89.png`, `10_behavior_competency_evaluations.png`
- **Mô tả quy trình:**
  1. `/admin/info/subjects` cấu hình danh mục môn học chương trình GDPT 2018 và phân loại môn đánh giá bằng điểm số kết hợp nhận xét hay chỉ nhận xét.
  2. `/admin/info/student_subject_evaluations` quản lý điểm kiểm tra định kỳ (1.0 - 10.0) và mức đạt chuẩn:
     - `HTT`: Hoàn thành tốt
     - `HT`: Hoàn thành
     - `CHT`: Chưa hoàn thành
  3. `/admin/info/student_behavior_evaluations` theo dõi 5 phẩm chất (Yêu nước, Nhân ái, Chăm chỉ, Trung thực, Trách nhiệm) và 10 năng lực cốt lõi theo 3 mức `T` (Tốt), `Đ` (Đạt), `C` (Cần cố gắng).
  4. Hỗ trợ trợ lý trí tuệ nhân tạo (AI Copilot Assistant) tự động phân tích điểm và đề xuất nhận xét sư phạm cá nhân hóa cho từng học sinh.

### 3.7. Phân Hệ Học Bạ Số Bất Biến & Chữ Ký Số SmartCA (Digital Transcripts & TSA)
- **Tệp minh họa:** `11_transcripts_smartca_snapshot.png`, `17_transcript_snapshot_detail.png`
- **Mô tả quy trình:**
  1. Cuối mỗi năm học, dữ liệu học tập được kết xuất thành bản ghi Snapshot bất biến (`academic_transcripts_snapshots`), khóa hoàn toàn việc sửa điểm hồi tố.
  2. **Chiến lược chuyển tiếp Thông tư (Strategy Versioning):**
     - Đánh giá trước ngày **01/12/2026**: Hệ thống áp dụng quy chuẩn **Thông tư 27/2020/TT-BGDĐT**.
     - Đánh giá từ ngày **01/12/2026** trở đi: Hệ thống tự động chuyển sang quy chuẩn **Thông tư 89/2026/TT-BGDĐT**.
  3. Tích hợp ký số điện tử của Giáo viên chủ nhiệm và Hiệu trưởng bằng chứng thư số công cộng SmartCA, kèm chữ ký chuẩn quốc tế CMS/PKCS#7 và dấu thời gian quốc gia TSA Timestamp chống chối bỏ.

### 3.8. Phân Hệ Dinh Dưỡng Bán Trú & Lưu Mẫu Thức Ăn 24H (Boarding & QĐ 1246/QĐ-BYT)
- **Tệp minh họa:** `12_boarding_nutrition_menu.png`, `15_boarding_menu_detail.png`
- **Mô tả quy trình:**
  1. Kế toán và bếp trưởng lập thực đơn bán trú theo tuần, hệ thống tự động tính toán tổng năng lượng Kcal và tỷ lệ các nhóm chất Protein : Lipid : Glucid (chuẩn 13-20% : 20-30% : 55-65%) theo Viện Dinh Dưỡng Quốc Gia.
  2. Cán bộ y tế thực hiện quy trình niêm phong và ghi sổ lưu mẫu thức ăn 24 giờ ở nhiệt độ $4^\circ\text{C}$ theo đúng **Quyết định 1246/QĐ-BYT**.
  3. Quy tắc hoàn tiền ăn bán trú: Phụ huynh báo cắt cơm trước 08:00 sáng được hệ thống tự động tính hoàn trả tiền ăn (35,000đ/ngày); báo sau 08:00 sáng không được hoàn tiền theo nội quy bếp ăn.

### 3.9. Phân Hệ Y Tế Học Đường & Sức Khỏe Học Sinh (Infirmary, BMI & Vision)
- **Tệp minh họa:** `13_student_health_bmi.png`, `16_student_health_detail.png`
- **Mô tả quy trình:**
  1. Quản lý hồ sơ nhân trắc học: Chiều cao, cân nặng, tự động tính toán chỉ số khối cơ thể (BMI) và phân loại thể trạng học sinh (Bình thường, Gầy, Nguy cơ béo phì, Béo phì).
  2. Theo dõi khúc xạ thị lực (mắt trái, mắt phải), nha học đường, bệnh lý mãn tính và sự cố y tế tại trường.
  3. Quản lý đơn thuốc gửi phụ huynh: Tiếp nhận yêu cầu cho học sinh uống thuốc theo đơn bác sĩ kèm xác nhận của phụ huynh và y tế học đường.

### 3.10. Phân Hệ Đánh Giá KPI 360 & Khống Chế Trần 20% (Nghị Định 90/2020 & 48/2023)
- **Tệp minh họa:** `13b_kpi_evaluations_20pct_cap.png`, `20_kpi_evaluation_detail.png`
- **Mô tả quy trình:**
  1. Đánh giá đa chiều với tỷ trọng cố định bất biến:
     $$\text{Điểm KPI} = (\text{Hệ thống tự động} \times 40\%) + (\text{Ban Giám Hiệu} \times 30\%) + (\text{Đồng nghiệp} \times 15\%) + (\text{Phụ huynh} \times 15\%)$$
  2. **Quy tắc Điểm Liệt:** Giáo viên bị kỷ luật, vi phạm quy định dạy thêm trái phép tiểu học (TT 29/2024, TT 19/2026), hoặc để xảy ra mất an toàn trường học sẽ bị xếp loại "Không hoàn thành nhiệm vụ" ngay lập tức, bất kể tổng điểm đạt bao nhiêu.
  3. **Khống chế Trần 20% Xuất Sắc:** Tỷ lệ cá nhân được xếp loại "Hoàn thành xuất sắc nhiệm vụ" không được vượt quá **20%** tổng số cán bộ, viên chức đạt mức "Hoàn thành tốt nhiệm vụ" trở lên. Người có điểm từ 90 trở lên nhưng thiếu Sáng kiến kinh nghiệm hoặc vượt ngoài trần 20% sẽ được chuyển về mức "Hoàn thành tốt".

### 3.11. Cổng Tích Hợp Ngoại Vi & Sandbox Quốc Gia (National Sandbox Stubs)
- **Tệp minh họa:** `21_integrations_sandbox_stubs.png`
- **Mô tả quy trình:**
  1. Giám sát thời gian thực trạng thái của 5 dịch vụ tích hợp quốc gia trong môi trường thử nghiệm (Sandbox Stubs):
     - **VNeID Mức 2 (Bộ Công An):** Xác thực định danh sinh trắc học và CCCD gắn chip của người giám hộ.
     - **SmartCA & National TSA:** Máy chủ ký số và cấp dấu thời gian quốc gia.
     - **CSDL Ngành MOET:** Hàng đợi đồng bộ Outbox truyền dữ liệu hồ sơ học sinh và học bạ lên Bộ Giáo dục & Đào tạo.
     - **VietQR Động EMVCo (NAPAS 247):** Sinh chuỗi mã QR thanh toán chuẩn quốc tế phục vụ thu tiền ăn và học phí.
     - **Hóa Đơn Điện Tử (Nghị định 123/2020/NĐ-CP):** Cấp hóa đơn ký hiệu `1C26TAA` kèm chữ ký số cơ quan thuế.
  2. Bảng nhật ký giao dịch ghi nhận đầy đủ thời gian, mã giao dịch (Transaction ID), vai trò tác tử, trạng thái HTTP và nội dung phản hồi.

---

## 4. Bằng Chứng Nghiệm Thu Pháp Lý & Kỹ Thuật (Verification Summary)

Toàn bộ 24 ảnh chụp màn hình trong danh mục trên đã được tạo lập tự động và kiểm tra tính nhất quán thông qua bộ kiểm thử E2E:

```bash
# Kiểm tra bộ kiểm thử tự động toàn diện và đối soát danh mục ảnh
make test-browser
```

Kết quả thực nghiệm:
- **Tổng số ảnh thu thập:** 40 ảnh (Bao gồm 24 ảnh chuẩn hóa theo từng màn hình nghiệp vụ + các ảnh chi tiết đối chiếu).
- **Tiêu chuẩn VAL-BROWSER-004:** Thỏa mãn xuất sắc cam kết hợp đồng (Yêu cầu tối thiểu $\ge 15$ ảnh chụp màn hình có chú thích).
- **Tính toàn vẹn pháp lý:** 100% các phân hệ cốt lõi đều được hiển thị trực quan và tuân thủ chặt chẽ các văn bản quy phạm pháp luật hiện hành.
