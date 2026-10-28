# BÁO CÁO ĐÁNH GIÁ TOÀN DIỆN HỆ THỐNG TÀI LIỆU DỰ ÁN & RÀ SOÁT HIỆU LỰC VĂN BẢN PHÁP LÝ CĂN CỨ
## DỰ ÁN: HỆ THỐNG PHẦN MỀM QUẢN LÝ TRƯỜNG TIỂU HỌC TÍCH HỢP KPI (EDUPRIMARY-KPI)

---

## I. TỔNG QUAN ĐÁNH GIÁ 5 TÀI LIỆU TRONG HỒ SƠ DỰ ÁN

Bộ hồ sơ gồm **5 tài liệu chuyên sâu** bao quát toàn bộ các khía cạnh từ nghiệp vụ, kiến trúc công nghệ, phân tích thị trường đến mô hình tài chính:
1. `DE_XUAT_CHUC_NANG_QLTH_KPI.md` (Đề xuất tính năng, phân hệ quản trị, bộ chỉ số KPI giáo dục, kết nối dữ liệu quốc gia).
2. `KIEN_TRUC_CORE_GO_VA_ADMIN.md` (Kiến trúc lõi Golang, Domain-Driven Design, Clean Architecture, mở rộng đa cấp học K-12).
3. `KIEN_TRUC_THICH_UNG_THAY_DOI_QUY_DINH.md` (Thiết kế kiến trúc Temporal Policy & Strategy Versioning thích ứng thay đổi văn bản pháp luật).
4. `PHAN_TICH_DOI_THU_VA_DE_XUAT_TINH_NANG_DOT_PHA.md` (Phân tích vnEdu, SMAS, MISA và 5 tính năng đột phá thị trường).
5. `PHAN_TICH_HIEU_QUA_DAU_TU_VA_MO_HINH_GIA_SAAS.md` (Hiệu quả tài chính, mô hình giá SaaS B2B/B2G, P&L 3 năm và chiến lược GTM).

### 1. Đánh giá chất lượng và điểm mạnh nổi bật
- **Am hiểu thực tiễn quản trị giáo dục Việt Nam**: Bộ tài liệu nắm trúng các nỗi đau (pain points) nhức nhối nhất của các trường học hiện nay: áp lực gõ nhận xét học sinh định kỳ theo Thông tư 27, tính toán định mức tiết dạy và giờ thừa theo quy định mới, xung đột xếp loại thi đua viên chức theo Nghị định 90 & Nghị định 48, và sự lúng túng trong triển khai Học bạ số.
- **Tư duy kiến trúc hệ thống hiện đại, gắn chặt với thể chế**: Điểm sáng giá nhất là tài liệu `KIEN_TRUC_THICH_UNG_THAY_DOI_QUY_DINH.md` đã giải quyết triệt để rủi ro "vỡ hệ thống khi văn bản luật thay đổi". Việc tách biệt logic nghiệp vụ thành *Temporal Policy Registry* và đóng băng kết quả học bạ bằng *Immutable Historical Snapshot* đạt chuẩn kiến trúc phần mềm cấp doanh nghiệp (Enterprise Architecture).
- **Tính liên kết và tính thực chiến cao**: Cả 5 tài liệu kết nối chặt chẽ từ căn cứ pháp lý $\rightarrow$ tính năng nghiệp vụ $\rightarrow$ thiết kế kỹ thuật backend $\rightarrow$ định vị cạnh tranh $\rightarrow$ bài toán đầu tư hoàn vốn.

### 2. Các điểm cần hoàn thiện bổ sung
- **Chuẩn hóa thời điểm hiệu lực của các văn bản mới ban hành**: Có văn bản đã ban hành nhưng chưa tới ngày có hiệu lực (như Thông tư 89/2026/TT-BGDĐT có hiệu lực từ ngày 01/12/2026), hệ thống cần phản ánh rõ trạng thái chuyển tiếp này.
- **Bổ sung văn bản sửa đổi mới phát sinh**: Thông tư 29/2024/TT-BGDĐT về dạy thêm học thêm đã được sửa đổi, bổ sung bởi **Thông tư 19/2026/TT-BGDĐT (hiệu lực 15/05/2026)**, cần cập nhật vào tiêu chí KPI kiểm soát đạo đức nhà giáo.
- **Hoàn thiện căn cứ bảo mật & an toàn dữ liệu trẻ em**: Bổ sung cam kết tuân thủ Điều 20 Nghị định 13/2023/NĐ-CP và các văn bản quy định về chữ ký số công vụ (Nghị định 130/2018/NĐ-CP, Thông tư 185/2019/TT-BQP).

---

## II. BẢNG TRA CỨU & ĐÁNH GIÁ CHI TIẾT HIỆU LỰC VĂN BẢN PHÁP LÝ

*(Kết quả tra cứu, đối chiếu từ Cổng Thông tin điện tử Bộ GD&ĐT, Cổng TTĐT Chính phủ và Hệ thống Cơ sở dữ liệu Quốc gia về Văn bản Pháp luật)*

### 1. Nhóm Nghiệp vụ Giáo dục Tiểu học, Đánh giá Học sinh & Chương trình

| STT | Số hiệu văn bản | Ngày ban hành / Hiệu lực | Trạng thái hiệu lực thực tế | Ý nghĩa & Lưu ý nghiệp vụ áp dụng |
|:---:|---|:---:|:---:|---|
| 1 | **Thông tư số 89/2026/TT-BGDĐT** | BH: **16/10/2026**<br>HL: **01/12/2026** | ⏳ **ĐÃ BAN HÀNH - CHƯA CÓ HIỆU LỰC** *(Hiệu lực từ 01/12/2026)* | **Văn bản sửa đổi tối quan trọng**: Sửa đổi, bổ sung Quy định đánh giá học sinh tiểu học (kèm TT 27/2020) và đánh giá học sinh THCS/THPT (kèm TT 22/2021). Do Thứ trưởng Phạm Ngọc Thưởng ký ban hành.<br>• *Điểm mới 1*: Chính thức pháp lý hóa Học bạ số trên toàn quốc; quy định ghi "Miễn" trong Học bạ số đối với môn học học sinh khuyết tật được miễn.<br>• *Điểm mới 2*: Cho phép tổ chức kiểm tra định kỳ trên máy tính / trực tuyến.<br>• *Điểm mới 3*: Đề kiểm tra chung cuối năm đối với Lớp 5.<br>• **Lưu ý kiến trúc**: Trước ngày 01/12/2026 hệ thống chạy chiến lược TT 27/2020; từ 01/12/2026 tự động kích hoạt chiến lược TT 89/2026. |
| 2 | **Thông tư số 27/2020/TT-BGDĐT** | BH: 04/09/2020<br>HL: 20/10/2020 | **ĐANG CÓ HIỆU LỰC** | Quy định đánh giá học sinh tiểu học. Lộ trình cuốn chiếu đã hoàn tất trong năm học 2024-2025 (toàn bộ Lớp 1 đến Lớp 5). Tiếp tục áp dụng cho đến khi các nội dung sửa đổi theo Thông tư 89/2026 chính thức có hiệu lực từ 01/12/2026. Thay thế hoàn toàn Thông tư 30/2014 và Thông tư 22/2016. |
| 3 | **Thông tư số 86/2026/TT-BGDĐT** | HL: **30/09/2026** | **ĐANG CÓ HIỆU LỰC** | Quy định về tổ chức hoạt động giáo dục tích hợp khoa học, công nghệ, kĩ thuật và toán học (**Giáo dục STEM**) trong giáo dục phổ thông. Tích hợp module theo dõi kế hoạch bài dạy STEM và tiêu chí đổi mới phương pháp giảng dạy trong KPI giáo viên. |
| 4 | **Thông tư số 28/2020/TT-BGDĐT** | BH: 04/09/2020<br>HL: 20/10/2020 | **ĐANG CÓ HIỆU LỰC** | Điều lệ trường tiểu học. Căn cứ nền tảng quy định về cơ cấu tổ chức nhà trường, nhiệm vụ, quyền hạn của Hiệu trưởng, Tổ chuyên môn, Giáo viên chủ nhiệm, Giáo viên bộ môn và danh mục hồ sơ quản lý giáo dục. |
| 5 | **Quyết định số 391/QĐ-BGDĐT** | BH: 24/01/2024<br>HL: 24/01/2024 | **ĐANG CÓ HIỆU LỰC** | Kế hoạch triển khai Học bạ số cấp tiểu học trên toàn quốc. Căn cứ chuẩn định dạng dữ liệu XML/JSON và quy trình ký số học bạ. |
| 6 | **Quyết định số 131/QĐ-TTg** | BH: 25/01/2022<br>HL: 25/01/2022 | **ĐANG CÓ HIỆU LỰC** | Đề án "Tăng cường ứng dụng CNTT và chuyển đổi số trong GD&ĐT giai đoạn 2022 - 2025, định hướng đến năm 2030" của Thủ tướng Chính phủ. |
| 7 | **Công văn số 2345/BGDĐT-GDTH** | BH: 07/06/2021 | **Đang áp dụng** | Hướng dẫn xây dựng kế hoạch giáo dục nhà trường và kế hoạch bài dạy (giáo án) cấp tiểu học. Căn cứ thiết kế module soạn giảng điện tử và ký duyệt giáo án. |

---

### 2. Nhóm Chế độ Làm việc Giáo viên, Vị trí Việc làm & Định mức Nhân sự

| STT | Số hiệu văn bản | Ngày ban hành / Hiệu lực | Trạng thái hiệu lực thực tế | Ý nghĩa & Lưu ý nghiệp vụ áp dụng |
|:---:|---|:---:|:---:|---|
| 8 | **Thông tư số 05/2025/TT-BGDĐT** | BH: **07/03/2025**<br>HL: **22/04/2025** | **ĐANG CÓ HIỆU LỰC** | Quy định chế độ làm việc đối với giáo viên phổ thông, dự bị đại học.<br>• **Thay thế hoàn toàn**: Thông tư số 28/2009/TT-BGDĐT và Thông tư số 15/2017/TT-BGDĐT.<br>• **Quy định đột phá**: Chuyển định mức giảng dạy sang tính theo năm học với **805 tiết/năm** đối với cấp tiểu học.<br>• **Tăng mức giảm tiết cho GVCN tiểu học**: Giảm từ **3 tiết/tuần lên 4 tiết/tuần**.<br>• Bổ sung định mức giảm trừ cho GV kiêm nhiệm tư vấn tâm lý, tổ phó chuyên môn, giáo viên dạy hòa nhập người khuyết tật. |
| 9 | **Thông tư số 20/2023/TT-BGDĐT** | BH: 30/10/2023<br>HL: **16/12/2023** | **ĐANG CÓ HIỆU LỰC** | Hướng dẫn về vị trí việc làm, cơ cấu viên chức theo chức danh nghề nghiệp và định mức số lượng người làm việc trong các cơ sở GDPT.<br>• **Thay thế hoàn toàn**: Thông tư số 16/2017/TT-BGDĐT.<br>• Căn cứ pháp lý chuẩn xác để phân quyền chức danh và định mức công việc cho khối hành chính - hỗ trợ: Văn thư, Kế toán, Thủ quỹ, Thư viện, Thiết bị, Y tế trường học, CNTT. |
| 10 | **Thông tư 02/2021/TT-BGDĐT & Thông tư 08/2023/TT-BGDĐT** | TT 08: HL **30/05/2023** | **ĐANG CÓ HIỆU LỰC** | Quy định mã số, tiêu chuẩn chức danh nghề nghiệp và bổ nhiệm, xếp lương viên chức giảng dạy tiểu học công lập (Hạng I, II, III). Thông tư 08/2023 đã bãi bỏ quy định chứng chỉ chức danh theo từng hạng, chỉ cần 1 chứng chỉ chung. |
| 11 | **Thông tư số 20/2018/TT-BGDĐT & Thông tư 14/2018/TT-BGDĐT** | HL: 2018 | **ĐANG CÓ HIỆU LỰC** | Chuẩn nghề nghiệp giáo viên cơ sở GDPT (5 tiêu chuẩn, 15 tiêu chí) và Chuẩn Hiệu trưởng GDPT. Căn cứ cốt lõi xây dựng cây chỉ số KPI chuyên môn cho BGH và Giáo viên. |

---

### 3. Nhóm Đánh giá Xếp loại Viên chức, Thi đua Khen thưởng & CSDL Nhân sự

| STT | Số hiệu văn bản | Ngày ban hành / Hiệu lực | Trạng thái hiệu lực thực tế | Ý nghĩa & Lưu ý nghiệp vụ áp dụng |
|:---:|---|:---:|:---:|---|
| 12 | **Nghị định 90/2020/NĐ-CP & Nghị định 48/2023/NĐ-CP** | NĐ 48: HL **15/09/2023** | **ĐANG CÓ HIỆU LỰC** | Đánh giá, xếp loại chất lượng cán bộ, công chức, viên chức. Căn cứ thuật toán quy đổi tự động điểm KPI sang 4 mức xếp loại nhà nước (Hoàn thành xuất sắc, Tốt, Hoàn thành, Không hoàn thành). Nghị định 48/2023 hoàn thiện cơ chế lượng hóa và loại trừ tiêu chí định tính. |
| 13 | **Luật Thi đua, Khen thưởng 06/2022/QH15 & Nghị định 98/2023/NĐ-CP** | HL: **01/01/2024** | **ĐANG CÓ HIỆU LỰC** | **Thay đổi mang tính bước ngoặt**: Bỏ quy định cứng bắt buộc phải có Sáng kiến kinh nghiệm (SKKN) mới được xét danh hiệu "Chiến sĩ thi đua cơ sở", mở rộng cho các cá nhân hoàn thành xuất sắc nhiệm vụ hoặc có mô hình đổi mới sáng tạo, giải pháp ứng dụng công nghệ hiệu quả. |
| 14 | **Thông tư số 06/2023/TT-BNV** | BH: 04/05/2023<br>HL: **01/07/2023** | **ĐANG CÓ HIỆU LỰC** | Quy chế cập nhật, sử dụng, khai thác CSDL quốc gia về cán bộ, công chức, viên chức trong các cơ quan nhà nước. Căn cứ kết nối API đẩy dữ liệu xếp loại viên chức về Bộ Nội vụ. |

---

### 4. Nhóm Pháp lý Số: Dữ liệu, Định danh, Đầu tư CNTT & An toàn Thông tin

| STT | Số hiệu văn bản | Ngày ban hành / Hiệu lực | Trạng thái hiệu lực thực tế | Ý nghĩa & Lưu ý nghiệp vụ áp dụng |
|:---:|---|:---:|:---:|---|
| 15 | **Luật Căn cước số 26/2023/QH15** | BH: 27/11/2023<br>HL: **01/07/2024** | **ĐANG CÓ HIỆU LỰC** | **Thay thế Luật Căn cước công dân 2014**. Chính thức cấp Căn cước và mã số định danh cá nhân cho người dưới 14 tuổi (học sinh tiểu học) $\rightarrow$ Căn cứ pháp lý then chốt giúp làm sạch dữ liệu hồ sơ học sinh và tuyển sinh đầu cấp. |
| 16 | **Nghị định số 69/2024/NĐ-CP** | BH: 25/06/2024<br>HL: **01/07/2024** | **ĐANG CÓ HIỆU LỰC** | **Thay thế hoàn toàn Nghị định 59/2022/NĐ-CP**. Căn cứ tích hợp đăng nhập một lần (SSO) qua tài khoản định danh điện tử **VNeID**. Cho phép phụ huynh dùng tài khoản định danh Mức 2 đại diện xác thực nhân thân cho con em dưới 14 tuổi (Điều 9). |
| 17 | **Luật Giao dịch điện tử 20/2023/QH15** | BH: 22/06/2023<br>HL: **01/07/2024** | **ĐANG CÓ HIỆU LỰC** | **Thay thế Luật Giao dịch điện tử 2005**. Khẳng định giá trị pháp lý của thông điệp dữ liệu, chứng thư điện tử, học bạ số và chữ ký số. |
| 18 | **Nghị định số 13/2023/NĐ-CP** | BH: 17/04/2023<br>HL: **01/07/2023** | **ĐANG CÓ HIỆU LỰC** | **Bảo vệ dữ liệu cá nhân**: Đặc biệt Điều 20 quy định bảo vệ dữ liệu cá nhân trẻ em. Bắt buộc hệ thống phải có sự đồng ý của cha mẹ/người giám hộ trước khi xử lý dữ liệu và mã hóa cấp cột dữ liệu nhạy cảm. |
| 19 | **Nghị định 73/2019/NĐ-CP & Nghị định 82/2024/NĐ-CP** | NĐ 82: HL **10/07/2024** | **ĐANG CÓ HIỆU LỰC** | Quản lý đầu tư ứng dụng CNTT sử dụng nguồn vốn ngân sách nhà nước. Nghị định 82/2024 sửa đổi, bổ sung NĐ 73/2019 là căn cứ bắt buộc khi triển khai dự án hoặc bán dịch vụ phần mềm cho khối trường công lập. |
| 20 | **Nghị định 85/2016/NĐ-CP & Thông tư 12/2022/TT-BTTTT** | Đang thi hành | **ĐANG CÓ HIỆU LỰC** | Bảo đảm an toàn hệ thống thông tin theo cấp độ (Hồ sơ hệ thống quản lý trường học phải đạt tiêu chuẩn **Cấp độ 3**). |
| 21 | **Quyết định 2568/QĐ-BTTTT** | BH: 29/12/2023 | **ĐANG CÓ HIỆU LỰC** | Khung Kiến trúc Chính phủ điện tử Việt Nam, phiên bản 3.0. Căn cứ thiết kế API, kết nối Trục liên thông NDXP/LGSP. |
| 22 | **Thông tư 42/2021/TT-BGDĐT & Quyết định 4998/QĐ-BGDĐT** | Đang thi hành | **ĐANG CÓ HIỆU LỰC** | Quy định về cơ sở dữ liệu và cấu trúc dữ liệu kỹ thuật CSDL ngành GD&ĐT (`csdl.moet.gov.vn`). |

---

### 5. Nhóm Hoạt động An toàn Trường học, Y tế, Thư viện, Dạy thêm & Chuyển đổi số

| STT | Số hiệu văn bản | Ngày ban hành / Hiệu lực | Trạng thái hiệu lực thực tế | Ý nghĩa & Lưu ý nghiệp vụ áp dụng |
|:---:|---|:---:|:---:|---|
| 23 | **Thông tư 29/2024/TT-BGDĐT & Thông tư 19/2026/TT-BGDĐT** | TT 29: HL 14/02/2025<br>TT 19: HL **15/05/2026** | **ĐANG CÓ HIỆU LỰC** | Quy định về dạy thêm, học thêm (thay thế TT 17/2012/TT-BGDĐT). **Thông tư 19/2026/TT-BGDĐT ban hành ngày 31/03/2026 đã sửa đổi Thông tư 29** về việc báo cáo, công khai hoạt động dạy thêm. Khoản 1 Điều 4 quy định **nghiêm cấm dạy thêm đối với học sinh tiểu học** (tiêu chí điểm liệt đạo đức nhà giáo trong bộ KPI). |
| 24 | **Thông tư số 18/2023/TT-BGDĐT** | BH: 26/10/2023<br>HL: **12/12/2023** | **ĐANG CÓ HIỆU LỰC** | Hướng dẫn xây dựng trường học an toàn, phòng, chống tai nạn thương tích trong cơ sở GDPT.<br>• **Thay thế hoàn toàn**: Thông tư 07/2004/TT-BGDĐT và Quyết định 4458/QĐ-BGDĐT. Căn cứ phân hệ Y tế, kiểm tra cơ sở vật chất và an toàn bán trú. |
| 25 | **Thông tư số 16/2022/TT-BGDĐT** | BH: 22/11/2022<br>HL: **07/01/2023** | **ĐANG CÓ HIỆU LỰC** | Tiêu chuẩn thư viện cơ sở giáo dục mầm non và phổ thông.<br>• **Thay thế hoàn toàn**: Quyết định 01/2003/QĐ-BGDĐT. Căn cứ cấu hình tiêu chuẩn thư viện mức 1, mức 2 trong phân hệ Quản lý thiết bị - Thư viện. |
| 26 | **Quyết định số 4725/QĐ-BGDĐT** | BH: 30/12/2022 | **ĐANG CÓ HIỆU LỰC** | Ban hành Bộ chỉ số đánh giá mức độ chuyển đổi số của cơ sở GDPT và GDTX. Căn cứ thiết kế thuật toán tự động tính điểm chỉ số chuyển đổi số trường học cho Ban Giám hiệu. |

---

## III. 4 ĐIỂM LƯU Ý CHIẾN LƯỢC KHI THỰC THI HỆ THỐNG

### 1. Quản trị thời điểm chuyển tiếp hiệu lực của Thông tư 89/2026/TT-BGDĐT (Mốc 01/12/2026)
- **Hiện tại**: Toàn bộ các trường tiểu học vẫn áp dụng cách đánh giá thường xuyên và định kỳ theo nguyên bản **Thông tư 27/2020/TT-BGDĐT**.
- **Kể từ ngày 01/12/2026**: Cơ chế `Strategy Versioning` trong kiến trúc Golang phải tự động kích hoạt chiến lược thực thi theo **Thông tư 89/2026/TT-BGDĐT**:
  + Bổ sung chức năng tổ chức bài kiểm tra định kỳ trên máy tính.
  + Mở cổng tiếp nhận đề kiểm tra chung môn Toán, Tiếng Việt Lớp 5 từ Phòng/Sở GD&ĐT.
  + Ghi trường giá trị `"Miễn"` vào Học bạ số đối với môn học học sinh khuyết tật được miễn trong kế hoạch giáo dục cá nhân.
- Cơ chế `Immutable Historical Snapshots` bảo đảm toàn bộ học bạ các năm học 2020 - 2025 được đóng băng nguyên vẹn theo chuẩn Thông tư 27/2020 cũ, không bị tính toán lại.

### 2. Cập nhật bổ sung Thông tư số 19/2026/TT-BGDĐT về Dạy thêm, học thêm
- Tài liệu trong hồ sơ hiện tại mới chỉ viện dẫn Thông tư 29/2024/TT-BGDĐT.
- Cần bổ sung ngay **Thông tư số 19/2026/TT-BGDĐT (ban hành 31/03/2026, có hiệu lực từ 15/05/2026)** sửa đổi Thông tư 29/2024 để đảm bảo tính cập nhật pháp lý tuyệt đối về trách nhiệm giải trình và kê khai dạy thêm ngoài nhà trường của giáo viên.

### 3. Tận dụng triệt để lợi thế cạnh tranh của Thông tư 05/2025/TT-BGDĐT
- Thông tư 05/2025/TT-BGDĐT đã chính thức có hiệu lực từ **22/04/2025**, bãi bỏ Thông tư 28/2009 và Thông tư 15/2017. 
- Việc phần mềm tự động tính toán chính xác định mức **805 tiết/năm** và tự động giảm **4 tiết/tuần cho GVCN tiểu học** là "vũ khí bán hàng" áp đảo so với các đối thủ (vnEdu, SMAS, MISA) vốn vẫn đang loay hoay với công thức tuần cũ hoặc phải tính thủ công bằng Excel.

### 4. Tuân thủ nghiêm ngặt Điều 20 Nghị định 13/2023/NĐ-CP (Bảo vệ dữ liệu trẻ em)
- Dữ liệu học sinh tiểu học là dữ liệu trẻ em dưới 14 tuổi.
- Bắt buộc phải triển khai quy trình **Thu thập sự đồng ý số (Digital Consent)** từ cha mẹ/người giám hộ thông qua tài khoản định danh VNeID trước khi hệ thống kích hoạt hồ sơ học sinh.
- Cam kết không sử dụng dữ liệu hành vi, học tập của học sinh cho mục đích quảng cáo hoặc chia sẻ với bên thứ ba không có thẩm quyền.

---

## IV. BẢNG KHUYẾN NGHỊ CẬP NHẬT TRỰC TIẾP VÀO CÁC FILE TRONG DỰ ÁN

| Tập tin | Vị trí cần bổ sung / chỉnh sửa | Nội dung cập nhật chi tiết |
|---|---|---|
| `DE_XUAT_CHUC_NANG_QLTH_KPI.md` | **Mục 1.1** (Căn cứ pháp lý) | Nêu rõ Thông tư 89/2026/TT-BGDĐT ban hành ngày 16/10/2026, có hiệu lực từ ngày **01/12/2026**; làm rõ mốc thời gian chuyển tiếp từ TT 27 sang TT 89. |
| `DE_XUAT_CHUC_NANG_QLTH_KPI.md` | **Mục 1.6** (Dạy thêm, học thêm) | Bổ sung **Thông tư số 19/2026/TT-BGDĐT ngày 31/03/2026 (hiệu lực 15/05/2026)** sửa đổi Thông tư 29/2024/TT-BGDĐT. |
| `DE_XUAT_CHUC_NANG_QLTH_KPI.md` | **Mục 6.3 & 7.1** (Bảo mật & Ký số) | Bổ sung căn cứ **Nghị định số 130/2018/NĐ-CP** và **Thông tư số 185/2019/TT-BQP** về chữ ký số chuyên dùng công vụ; làm rõ luồng Digital Consent theo Điều 20 Nghị định 13/2023/NĐ-CP. |
| `KIEN_TRUC_THICH_UNG_THAY_DOI_QUY_DINH.md` | **Mục 1 & Mục 3** (Strategy Versioning) | Bổ sung mã ví dụ cấu hình chính sách: Policy `TT_27_2020` áp dụng đến 30/11/2026; Policy `TT_89_2026` kích hoạt tự động từ 01/12/2026. |
| `PHAN_TICH_HIEU_QUA_DAU_TU_VA_MO_HINH_GIA_SAAS.md` | **Mục 5.2** (Lợi thế cạnh tranh) | Bổ sung câu chuyện bán hàng: *Cập nhật Thông tư 05/2025, Thông tư 89/2026 và Thông tư 19/2026* giúp nhà trường hoàn toàn an tâm trước các đợt thanh kiểm tra của Sở/Phòng GD&ĐT. |

---
*Báo cáo được lập trên cơ sở rà soát đối chiếu thực tế với hệ thống Văn bản Quy phạm Pháp luật Quốc gia tại ngày 27/10/2026.*