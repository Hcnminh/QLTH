# BÁO CÁO NGHIÊN CỨU & ĐỀ XUẤT CHỨC NĂNG
## HỆ THỐNG PHẦN MỀM QUẢN LÝ TRƯỜNG TIỂU HỌC TÍCH HỢP BỘ CÔNG CỤ ĐÁNH GIÁ KPI (EDUPRIMARY-KPI)

---

## MỤC LỤC
1. [CĂN CỨ PHÁP LÝ & TIÊU CHUẨN KỸ THUẬT QUỐC GIA](#1-căn-cứ-pháp-lý--tiêu-chuẩn-kỹ-thuật-quốc-gia)
2. [ĐỐI TƯỢNG SỬ DỤNG & MA TRẬN PHÂN QUYỀN](#2-đối-tượng-sử-dụng--ma-trận-phân-quyền)
3. [KIẾN TRÚC TỔNG THỂ HỆ THỐNG](#3-kiến-trúc-tổng-thể-hệ-thống)
4. [CHI TIẾT CÁC PHÂN HỆ NGHIỆP VỤ QUẢN TRỊ TRƯỜNG TIỂU HỌC](#4-chi-tiết-các-phân-hệ-nghiệp-vụ-quản-trị-trường-tiểu-học)
5. [PHÂN HỆ CHUYÊN SÂU: ĐÁNH GIÁ HIỆU SUẤT & BỘ CHỈ SỐ KPI GIÁO DỤC](#5-phân-hệ-chuyên-sâu-đánh-giá-hiệu-suất--bộ-chỉ-số-kpi-giáo-dục)
6. [KHẢ NĂNG TÍCH HỢP & LIÊN THÔNG DỮ LIỆU](#6-khả-năng-tích-hợp--liên-thông-dữ-liệu)
7. [YÊU CẦU PHI CHỨC NĂNG & BẢO MẬT DỮ LIỆU CÁ NHÂN TRẺ EM](#7-yêu-cầu-phi-chức-năng--bảo-mật-dữ-liệu-cá-nhân-trẻ-em)
8. [LỘ TRÌNH TRIỂN KHAI & CÔNG NGHỆ KHUYẾN NGHỊ](#8-lộ-trình-triển-khai--công-nghệ-khuyến-nghị)
9. [QUẢN TRỊ RỦI RO & PHƯƠNG ÁN AN TOÀN THÔNG TIN PHỤC VỤ THẨM ĐỊNH](#9-quản-trị-rủi-ro--phương-án-an-toàn-thông-tin-phục-vụ-thẩm-định)
10. [DANH MỤC TÀI LIỆU TRONG HỒ SƠ THẨM ĐỊNH DỰ ÁN](#10-danh-mục-tài-liệu-trong-hồ-sơ-thẩm-định-dự-án)

---

## 1. CĂN CỨ PHÁP LÝ & TIÊU CHUẨN KỸ THUẬT QUỐC GIA

Phần mềm được thiết kế tuân thủ nghiêm ngặt hệ thống văn bản quy phạm pháp luật hiện hành của Bộ Giáo dục & Đào tạo (BGDĐT), Bộ Nội vụ (BNV), và các tiêu chuẩn bảo vệ dữ liệu:

### 1.1. Nghiệp vụ chuyên môn tiểu học & học sinh
- **Thông tư số 28/2020/TT-BGDĐT** ban hành Điều lệ trường tiểu học (Cơ cấu tổ chức, quyền và nghĩa vụ của giáo viên/học sinh, hồ sơ quản lý hoạt động giáo dục, sách giáo khoa, thiết bị dạy học) - *Đang có hiệu lực thi hành*.
- **Thông tư số 27/2020/TT-BGDĐT** ban hành Quy định đánh giá học sinh tiểu học (áp dụng trọn vẹn cho Lớp 1 đến Lớp 5).
  - *Lưu ý lộ trình chuyển tiếp pháp lý*: Đã được sửa đổi, bổ sung toàn diện bởi **Thông tư số 89/2026/TT-BGDĐT** (Thứ trưởng Phạm Ngọc Thưởng ký ban hành ngày **16/10/2026**, chính thức có hiệu lực từ ngày **01/12/2026**). Hệ thống áp dụng nguyên bản TT 27/2020 đến hết ngày 30/11/2026 và tự động chuyển sang cơ chế TT 89/2026 từ 01/12/2026:
    + *Điểm mới 1 (Chính thức hóa Học bạ số)*: Pháp lý hóa việc sử dụng Học bạ số trên toàn quốc, có giá trị pháp lý thay thế hoàn toàn học bạ giấy truyền thống; quy định ghi rõ "Miễn" trong Học bạ số đối với môn học học sinh khuyết tật được miễn trong kế hoạch giáo dục cá nhân.
    + *Điểm mới 2 (Kiểm tra định kỳ trên máy tính)*: Cho phép các cơ sở giáo dục tổ chức kiểm tra định kỳ bằng hình thức trực tuyến hoặc trên phòng máy tính của trường.
    + *Điểm mới 3 (Đề kiểm tra chung cho học sinh cuối cấp lớp 5)*: Quy định việc tổ chức kiểm tra định kỳ cuối năm đối với học sinh lớp 5 theo đề kiểm tra chung của Phòng GD&ĐT / Sở GD&ĐT đối với môn Toán và Tiếng Việt nhằm đánh giá chất lượng đầu ra thực chất và chống bệnh thành tích.
    + *Đánh giá thường xuyên*: Nhận xét sự tiến bộ về kiến thức các môn học/hoạt động giáo dục; hình thành và phát triển phẩm chất chủ yếu (Yêu nước, Nhân ái, Chăm chỉ, Trung thực, Trách nhiệm) và năng lực cốt lõi.
    + *Đánh giá định kỳ*: Mức Đạt tốt (T), Hoàn thành tốt (HTT) / Đạt (Đ), Hoàn thành (HT) / Cần cố gắng (C), Chưa hoàn thành (CHT); bài kiểm tra định kỳ cho điểm số (1-10).
- **Thông tư số 86/2026/TT-BGDĐT** (có hiệu lực từ ngày 30/09/2026): Quy định về tổ chức giáo dục tích hợp khoa học, công nghệ, kỹ thuật và toán học (Giáo dục STEM) trong các cơ sở giáo dục phổ thông.
- **Quyết định số 131/QĐ-TTg** của Thủ tướng Chính phủ phê duyệt Đề án "Tăng cường ứng dụng CNTT và chuyển đổi số trong giáo dục và đào tạo giai đoạn 2022 - 2025, định hướng đến năm 2030".
- **Quyết định số 391/QĐ-BGDĐT ngày 24/01/2024** của Bộ trưởng Bộ GD&ĐT ban hành Kế hoạch triển khai Học bạ số cấp tiểu học và nhân rộng toàn quốc.

### 1.2. Nghiệp vụ cán bộ quản lý, giáo viên & nhân viên
- **Thông tư số 05/2025/TT-BGDĐT ngày 07/03/2025** (có hiệu lực từ ngày 22/04/2025) quy định chế độ làm việc đối với giáo viên phổ thông, dự bị đại học (*Cập nhật mới nhất - Thay thế hoàn toàn Thông tư 28/2009/TT-BGDĐT và Thông tư 15/2017/TT-BGDĐT*):
  - Định mức tiết dạy hàng năm của giáo viên trường tiểu học là **805 tiết/năm học** (trung bình 23 tiết/tuần $\times$ 35 tuần giảng dạy).
  - **Tăng số tiết giảm cho giáo viên chủ nhiệm (GVCN) cấp tiểu học từ 3 tiết/tuần lên 4 tiết/tuần** (Điều chỉnh quan trọng tác động trực tiếp đến định mức và tính giờ làm việc).
  - Bổ sung quy định giảm định mức tiết dạy cho giáo viên kiêm nhiệm tư vấn tâm lý học sinh, tổ phó chuyên môn, giáo viên dạy người khuyết tật hòa nhập, dạy tiếng Việt cho trẻ em dân tộc thiểu số trước khi vào lớp Một.
  - Chuẩn hóa công thức tính làm thêm giờ (tăng giờ) khi vượt định mức 805 tiết/năm.
- **Thông tư số 20/2018/TT-BGDĐT**: Quy định chuẩn nghề nghiệp giáo viên cơ sở giáo dục phổ thông (5 tiêu chuẩn, 15 tiêu chí: Phẩm chất nhà giáo, Phát triển chuyên môn/nghiệp vụ, Xây dựng môi trường giáo dục, Phát triển mối quan hệ nhà trường - gia đình - xã hội, Sử dụng ngoại ngữ/tin học) - *Đang có hiệu lực thi hành*.
- **Thông tư số 14/2018/TT-BGDĐT**: Quy định chuẩn hiệu trưởng cơ sở giáo dục phổ thông - *Đang có hiệu lực thi hành*.
- **Thông tư số 02/2021/TT-BGDĐT** (đã được sửa đổi, bổ sung bởi **Thông tư số 08/2023/TT-BGDĐT** ngày 14/04/2023 có hiệu lực từ 30/05/2023): Quy định mã số, tiêu chuẩn chức danh nghề nghiệp và bổ nhiệm, xếp lương viên chức giảng dạy trong các trường tiểu học công lập (Hạng I, II, III; bãi bỏ quy định chứng chỉ chức danh theo từng hạng, chỉ cần 1 chứng chỉ chung).
- **Thông tư số 20/2023/TT-BGDĐT ngày 30/10/2023** (có hiệu lực từ 16/12/2023, thay thế Thông tư 16/2017/TT-BGDĐT): Hướng dẫn về vị trí việc làm, cơ cấu viên chức theo chức danh nghề nghiệp và định mức số lượng người làm việc trong các cơ sở GDPT công lập (Giáo viên, Giám thị/Tổng phụ trách Đội, Thiết bị - thí nghiệm, Thư viện, Văn thư, Kế toán, Thủ quỹ, Y tế trường học, CNTT...).

### 1.3. Đánh giá xếp loại viên chức, thi đua khen thưởng & cơ chế KPI
- **Luật Thi đua, Khen thưởng số 06/2022/QH15** (có hiệu lực từ 01/01/2024) và **Nghị định số 98/2023/NĐ-CP**:
  - *Điểm mới then chốt*: **Không còn bắt buộc phải có Sáng kiến kinh nghiệm (SKKN)** mới được xét danh hiệu "Chiến sĩ thi đua cơ sở", mà mở rộng cho các cá nhân hoàn thành xuất sắc nhiệm vụ hoặc có mô hình đổi mới sáng tạo, giải pháp ứng dụng hiệu quả.
- **Nghị định số 90/2020/NĐ-CP** (đã được sửa đổi, bổ sung bởi **Nghị định số 48/2023/NĐ-CP** ngày 17/07/2023 có hiệu lực từ 15/09/2023): Quy định về đánh giá, xếp loại chất lượng cán bộ, công chức, viên chức (Hoàn thành xuất sắc nhiệm vụ, Hoàn thành tốt nhiệm vụ, Hoàn thành nhiệm vụ, Không hoàn thành nhiệm vụ; cơ chế lượng hóa và loại trừ tiêu chí định tính; hướng dẫn rõ quy định xếp loại đối với trường hợp bị xử lý kỷ luật).
- **Mô hình BSC - KPI & OKR trong giáo dục**: Chuẩn hóa đo lường mục tiêu chất lượng giảng dạy, tỷ lệ học sinh hoàn thành môn học, chỉ số hài lòng của phụ huynh, mức độ chuyển đổi số cá nhân, kỷ luật lao động và phối hợp tổ bộ môn.

### 1.4. Tiêu chuẩn kết nối dữ liệu dùng chung Quốc gia, Bộ GD&ĐT & An toàn thông tin
- **Luật Giao dịch điện tử số 20/2023/QH15** (có hiệu lực từ ngày 01/07/2024): Quy định về thông điệp dữ liệu, chứng thư điện tử, giá trị pháp lý của học bạ điện tử, chữ ký số và kết nối chia sẻ dữ liệu trong cơ quan nhà nước.
- **Luật Căn cước số 26/2023/QH15** (có hiệu lực từ ngày 01/07/2024): Quy định về cấp thẻ căn cước và số định danh cá nhân cho người dưới 14 tuổi (học sinh tiểu học).
- **Nghị định số 69/2024/NĐ-CP ngày 25/06/2024** (có hiệu lực từ ngày 01/07/2024, *thay thế hoàn toàn Nghị định 59/2022/NĐ-CP*): Quy định về định danh và xác thực điện tử (sử dụng tài khoản VNeID của phụ huynh để đại diện xác thực cho trẻ em dưới 14 tuổi).
- **Quyết định số 2568/QĐ-BTTTT ngày 29/12/2023**: Ban hành Khung Kiến trúc Chính phủ điện tử Việt Nam, phiên bản 3.0.
- **Quyết định số 142/QĐ-TTg ngày 02/02/2024**: Phê duyệt Chiến lược dữ liệu quốc gia đến năm 2030.
- **Nghị định số 47/2020/NĐ-CP**: Quy định về quản lý, kết nối và chia sẻ dữ liệu số của cơ quan nhà nước (kết nối qua Trục liên thông dữ liệu quốc gia NDXP/LGSP).
- **Thông tư số 42/2021/TT-BGDĐT** và **Quyết định số 4998/QĐ-BGDĐT**: Quy định về cơ sở dữ liệu và cấu trúc dữ liệu kỹ thuật CSDL ngành GD&ĐT (`csdl.moet.gov.vn`).
- **Thông tư số 06/2023/TT-BNV**: Quy chế cập nhật, sử dụng, khai thác CSDL quốc gia về cán bộ, công chức, viên chức trong các cơ quan nhà nước (Bộ Nội vụ).
- **Nghị định số 73/2019/NĐ-CP** và **Nghị định số 82/2024/NĐ-CP**: Quản lý đầu tư ứng dụng CNTT sử dụng nguồn vốn ngân sách nhà nước.
- **Nghị định số 13/2023/NĐ-CP** về Bảo vệ dữ liệu cá nhân (Đặc biệt Điều 20 bảo vệ dữ liệu trẻ em).
- **Nghị định số 85/2016/NĐ-CP** và **Thông tư số 12/2022/TT-BTTTT**: Bảo đảm an toàn hệ thống thông tin theo cấp độ (Cấp độ 3).

### 1.5. Đào tạo, bồi dưỡng thường xuyên & Đánh giá chuyển đổi số trường học
- **Quyết định số 4725/QĐ-BGDĐT ngày 30/12/2022**: Ban hành Bộ chỉ số đánh giá mức độ chuyển đổi số của cơ sở giáo dục phổ thông và giáo dục thường xuyên (Nhóm 1: Chuyển đổi số trong dạy, học; Nhóm 2: Chuyển đổi số trong quản trị nhà trường - Tích hợp tính điểm tự động trên phần mềm).
- **Thông tư số 19/2019/TT-BGDĐT** (đã được sửa đổi, bổ sung bởi **Thông tư số 17/2022/TT-BGDĐT**): Ban hành Quy chế bồi dưỡng thường xuyên giáo viên, cán bộ quản lý cơ sở giáo dục mầm non, cơ sở giáo dục phổ thông.
- **Thông tư số 17/2019/TT-BGDĐT** và **Thông tư số 18/2019/TT-BGDĐT**: Ban hành Chương trình bồi dưỡng thường xuyên giáo viên và cán bộ quản lý cơ sở giáo dục phổ thông (Tích hợp theo dõi tiến độ hoàn thành các mô-đun bồi dưỡng thường xuyên vào KPI-GV-06).

### 1.6. Kỷ luật nghề nghiệp, An toàn trường học, Y tế, Dinh dưỡng & Kiểm định chất lượng
- **Thông tư số 29/2024/TT-BGDĐT ngày 30/12/2024** (có hiệu lực từ ngày 14/02/2025, thay thế Thông tư 17/2012/TT-BGDĐT) và **Thông tư số 19/2026/TT-BGDĐT ngày 31/03/2026** (có hiệu lực từ ngày **15/05/2026** sửa đổi, bổ sung Thông tư 29): Quy định về dạy thêm, học thêm.
  - *Khoản 1 Điều 4*: **Nghiêm cấm dạy thêm đối với học sinh tiểu học** (trừ bồi dưỡng nghệ thuật, TDTT, rèn luyện kỹ năng sống).
  - *Quy định sửa đổi của TT 19/2026*: Chuẩn hóa quy trình báo cáo, công khai hoạt động dạy học ngoài nhà trường và trách nhiệm giải trình của viên chức giảng dạy. Đây là tiêu chí kiểm soát đạo đức nhà giáo và điểm liệt bắt buộc trong bộ chỉ số KPI.
- **Thông tư số 18/2023/TT-BGDĐT ngày 26/10/2023**: Hướng dẫn xây dựng trường học an toàn, phòng, chống tai nạn thương tích trong cơ sở giáo dục phổ thông (thay thế Thông tư 07/2004/TT-BGDĐT và Quyết định 4458/QĐ-BGDĐT).
- **Quyết định số 85/QĐ-TTg ngày 17/01/2022**: Phê duyệt Chương trình y tế trường học trong các cơ sở giáo dục mầm non và phổ thông gắn với y tế cơ sở giai đoạn 2021 - 2025.
- **Quyết định số 41/QĐ-TTg ngày 08/01/2019**: Phê duyệt Đề án bảo đảm dinh dưỡng hợp lý và tăng cường hoạt động thể lực cho trẻ em, học sinh.
- **Thông tư số 17/2018/TT-BGDĐT**: Quy định về kiểm định chất lượng giáo dục và công nhận đạt chuẩn quốc gia đối với trường tiểu học (5 tiêu chuẩn phục vụ công tác thanh kiểm tra của BGH).
- **Thông tư số 31/2017/TT-BGDĐT**: Hướng dẫn công tác tư vấn tâm lý học sinh trong trường phổ thông (đóng vai trò căn cứ giảm trừ định mức tiết dạy cho GV kiêm nhiệm theo TT 05/2025/TT-BGDĐT).
- **Thông tư số 03/2018/TT-BGDĐT**: Quy định về giáo dục hòa nhập đối với người khuyết tật (căn cứ xây dựng kế hoạch giáo dục cá nhân và tính phụ cấp/hệ số tiết dạy cho giáo viên).

---

## 2. ĐỐI TƯỢNG SỬ DỤNG & MA TRẬN PHÂN QUYỀN

| Đối tượng | Vai trò chính | Các chức năng tiếp cận chính |
|---|---|---|
| **Ban Giám hiệu (Hiệu trưởng, P. Hiệu trưởng)** | Quản trị điều hành toàn diện | Dashboard chỉ số trường học, duyệt kế hoạch GD, phê duyệt KPI toàn trường, ký số học bạ/báo cáo, thanh tra giáo án, đối soát báo cáo CSDL ngành. |
| **Tổ trưởng chuyên môn / Khối trưởng** | Quản lý học thuật & nghiệp vụ khối | Phân phối chương trình, sinh hoạt chuyên môn theo NCBH, sơ duyệt giáo án, đánh giá KPI thành viên trong tổ. |
| **Giáo viên chủ nhiệm (GVCN)** | Quản lý học sinh & tương tác phụ huynh | Sổ điểm, nhận xét định kỳ theo TT27, điểm danh, bàn giao học bạ số, liên lạc phụ huynh, tự đánh giá KPI. |
| **Giáo viên bộ môn (Tin, Anh, Mỹ thuật, Âm nhạc...)** | Giảng dạy & đánh giá chuyên môn | Nhập nhận xét định kỳ môn học, giáo án điện tử, chấm điểm kiểm tra theo TT27, tự đánh giá KPI. |
| **Tổng phụ trách Đội** | Công tác Đội & phong trào Sao Nhi Đồng | Quản lý nề nếp Sao nhi đồng, phong trào Thiếu nhi, chấm điểm thi đua các lớp, đóng góp chỉ số KPI nề nếp. |
| **Nhân viên Hành chính - Hỗ trợ (Kế toán, Y tế, Thư viện, Văn thư)** | Nghiệp vụ chuyên trách trường học | Thu học phí/bán trú, hồ sơ sức khỏe/tiêm chủng, mượn trả thiết bị/sách, văn thư lưu trữ, KPI khối hỗ trợ. |
| **Phụ huynh học sinh (PHHS)** | Giám sát & đồng hành | Xem kết quả học tập (nhận xét + điểm), thực đơn/bán trú, đóng phí trực tuyến, gửi phản hồi/đánh giá mức độ hài lòng (CSAT đóng góp vào KPI trường). |
| **Cơ quan quản lý (Phòng GD&ĐT / Sở GD&ĐT)** | Giám sát vĩ mô & kiểm tra | Khai thác API đồng bộ dữ liệu, kiểm tra thống kê chất lượng giáo dục, kiểm tra chuẩn phổ cập. |

---

## 3. KIẾN TRÚC TỔNG THỂ HỆ THỐNG

Hệ thống được thiết kế theo kiến trúc **Cloud-native Microservices**, đảm bảo tính sẵn sàng cao, bảo mật nhiều lớp và khả năng tích hợp mở rộng:

```
+---------------------------------------------------------------------------------------------------------+
|                                     LỚP TRÌNH DIỄN (FRONTEND & ADMIN)                                   |
|  +--------------------+  +----------------------+  +---------------------+  +-------------+  +--------+ |
|  | Web Portal Quản lý |  | Web Giáo viên/Cán bộ |  | Go Auto-Admin Portal|  | Mobile App  |  | Mobile | |
|  | (BGH, Khối trưởng) |  | (Nhận xét/Giáo án)   |  | (Cấu hình danh mục) |  | PHHS        |  | App GV | |
|  +--------------------+  +----------------------+  +---------------------+  +-------------+  +--------+ |
+---------------------------------------------------------------------------------------------------------+
                                              | API Gateway (OIDC/OAuth2, RBAC, WAF)
+------------------------------------------------------------------------------------------+
|                                 LỚP DỊCH VỤ NGHIỆP VỤ (MICROSERVICES)                     |
|  +------------------------+  +--------------------------+  +---------------------------+ |
|  | 1. Quản lý Học sinh    |  | 2. Đánh giá TT 27/2020   |  | 3. Quản lý Kế hoạch &     | |
|  |    & Hồ sơ Học bạ số   |  |    (Nhận xét & Điểm số)  |  |    Giáo án điện tử        | |
|  +------------------------+  +--------------------------+  +---------------------------+ |
|  +------------------------+  +--------------------------+  +---------------------------+ |
|  | 4. Bán trú, Dinh dưỡng |  | 5. Thiết bị, Thư viện,   |  | 6. Tài chính, Học phí     | |
|  |    & Y tế học đường    |  |    Y tế & Văn thư số     |  |    & Cổng thanh toán      | |
|  +------------------------+  +--------------------------+  +---------------------------+ |
|  +-------------------------------------------------------------------------------------+ |
|  | 7. PHÂN HỆ ĐÁNH GIÁ HIỆU SUẤT & BỘ CHỈ SỐ KPI GIÁO VIÊN / NHÂN VIÊN (KPI ENGINE)    | |
|  |    (Thiết lập trọng số - Đo lường tự động - Đánh giá 360 - Xếp loại NĐ 90/NĐ 48)   | |
|  +-------------------------------------------------------------------------------------+ |
|  +------------------------+  +--------------------------+  +---------------------------+ |
|  | 8. Truyền thông & Sổ   |  | 9. Ký số điện tử         |  | 10. Đồng bộ CSDL Ngành    | |
|  |    liên lạc điện tử    |  |    (SmartCA/USB Token)   |  |     (MOET API / Phòng GD) | |
|  +------------------------+  +--------------------------+  +---------------------------+ |
+------------------------------------------------------------------------------------------+
                                              |
+------------------------------------------------------------------------------------------+
|                                  LỚP DỮ LIỆU & NỀN TẢNG (DATA & INFRA)                    |
|  PostgreSQL (Relational/Audit) | Redis (Cache/Session) | MinIO/S3 (Hồ sơ/Giáo án/Học bạ)  |
|  Elasticsearch (Tìm kiếm nhanh)| Kafka (Event-driven)  | Docker & Kubernetes (K8s)        |
+------------------------------------------------------------------------------------------+
```

---

## 4. CHI TIẾT CÁC PHÂN HỆ NGHIỆP VỤ QUẢN TRỊ TRƯỜNG TIỂU HỌC

### Phân hệ 1: Quản lý Hồ sơ Học sinh & Học bạ số
1. **Quản lý lý lịch học sinh**:
   - Quản lý mã định danh cá nhân (kết nối CSDL quốc gia về Dân cư VNeID).
   - Quản lý quá trình chuyển trường, tiếp nhận, thôi học, phân lớp tự động hoặc theo nguyện vọng.
   - Quản lý đối tượng chính sách: Miễn giảm, học sinh khuyết tật học hòa nhập (có kế hoạch giáo dục cá nhân riêng theo quy định).
2. **Học bạ điện tử (Học bạ số)**:
   - Đáp ứng định dạng chuẩn XML/JSON theo đề án Học bạ số của Bộ GD&ĐT (Quyết định 391/QĐ-BGDĐT).
   - Tích hợp ký số tập trung (SmartCA/Cloud CA) hoặc ký số token phần cứng cho BGH và Giáo viên chủ nhiệm.
   - Cơ chế khóa dữ liệu (Locking mechanism) và lưu vết chỉnh sửa (Audit Trail) để đảm bảo không thể sửa đổi sau khi đã ký số đóng sổ cuối năm học.
   - Bàn giao kết quả giáo dục trực tuyến giữa các khối lớp và chuyển cấp lên THCS.

### Phân hệ 2: Đánh giá Học sinh theo Thông tư 27/2020/TT-BGDĐT & Thông tư 89/2026/TT-BGDĐT
1. **Đánh giá thường xuyên**:
   - **Các môn học & hoạt động giáo dục**: Ghi nhận nhận xét định kỳ về sự tiến bộ, những điểm nổi bật và biện pháp khắc phục. Hỗ trợ kho ngân hàng nhận xét mẫu chuẩn mực sư phạm, tự động gợi ý nhận xét thông minh (tránh nhận xét rập khuôn, vô cảm).
   - **Sự hình thành và phát triển phẩm chất**: Đánh giá 5 phẩm chất (Yêu nước, Nhân ái, Chăm chỉ, Trung thực, Trách nhiệm) với 3 mức: Tốt (T), Đạt (Đ), Cần cố gắng (C).
   - **Sự hình thành và phát triển năng lực**: Đánh giá 3 năng lực chung và 7 năng lực đặc thù theo 3 mức (T, Đ, C).
2. **Đánh giá định kỳ & Tổ chức thi trên máy tính (Cập nhật TT 89/2026/TT-BGDĐT)**:
   - Ghi nhận điểm số bài kiểm tra định kỳ (giữa kỳ 1, cuối kỳ 1, giữa kỳ 2, cuối kỳ 2) cho các môn có điểm số: Tiếng Việt, Toán (Lớp 1-5); Ngoại ngữ 1, Lịch sử & Địa lý, Khoa học, Tin học và Công nghệ (Lớp 4-5).
   - **Tổ chức kiểm tra trên máy tính**: Cho phép nhà trường tạo ca thi, giao đề trắc nghiệm/tự luận kết hợp trên máy tính phòng học hoặc trực tuyến có giám sát, tự động chấm điểm và đồng bộ vào sổ điểm.
   - **Tiếp nhận ngân hàng đề thi chung cho Lớp 5**: Kết nối module bảo mật để tiếp nhận và phân phối đề kiểm tra định kỳ cuối năm từ Phòng GD&ĐT / Sở GD&ĐT đối với môn Toán và Tiếng Việt cho học sinh cuối cấp lớp 5.
   - Hỗ trợ công thức tính toán và tổng hợp xếp loại tự động: Hoàn thành xuất sắc, Hoàn thành tốt, Hoàn thành, Chưa hoàn thành.
   - Xét khen thưởng cuối năm: Khen thưởng học sinh Xuất sắc, Khen thưởng học sinh Tiêu biểu hoàn thành tốt trong học tập và rèn luyện; khen thưởng đột xuất.
3. **Sổ theo dõi và đánh giá học sinh điện tử & Học bạ số**:
   - Xuất file PDF/XML chuẩn biểu mẫu của BGDĐT, ký số tập trung và đồng bộ lên Kho lưu trữ quốc gia.

### Phân hệ 3: Quản lý Hoạt động Chuyên môn, Giáo án & Định mức Tiết dạy (Thông tư 05/2025/TT-BGDĐT)
1. **Phân phối chương trình, Thời khóa biểu & Quản lý định mức tiết dạy**:
   - Quản lý khung phân phối chương trình GDPT 2018 (35 tuần/năm học, 2 buổi/ngày).
   - Xếp thời khóa biểu thông minh: Tự động tránh trùng tiết, cân bằng số tiết dạy/ngày của giáo viên, ưu tiên môn học buổi sáng và buổi chiều phù hợp tâm sinh lý học sinh tiểu học.
   - **Tự động hóa quản lý định mức giảng dạy theo Thông tư 05/2025/TT-BGDĐT**:
     + Theo dõi định mức chuẩn: **805 tiết/năm học** của giáo viên tiểu học.
     + Tự động áp dụng giảm trừ định mức: Giảm **4 tiết/tuần** cho Giáo viên chủ nhiệm lớp; giảm trừ cho tổ trưởng/tổ phó chuyên môn, giáo viên kiêm nhiệm tư vấn tâm lý, Tổng phụ trách Đội...
     + Tính hệ số quy đổi tiết dạy (dạy trẻ khuyết tật hòa nhập, dạy tiếng Việt cho học sinh DTTS trước lớp 1).
     + Tự động tính toán số tiết dạy thêm giờ (tăng giờ) và xuất bảng đối soát thừa giờ cuối học kỳ/năm học.
2. **Giáo án & Kế hoạch bài dạy (KHBD) điện tử**:
   - Soạn thảo và tải lên kế hoạch bài dạy theo mẫu Công văn 2345/BGDĐT-GDTH.
   - Quy trình nộp - duyệt giáo án phân tầng: Giáo viên -> Tổ trưởng chuyên môn duyệt -> Ban giám hiệu ký duyệt kiểm tra.
   - Ghi nhận lịch sử duyệt, nhận xét, ký số vào giáo án trực tuyến.
3. **Sinh hoạt chuyên môn theo Nghiên cứu bài học**:
   - Lập kế hoạch sinh hoạt tổ chuyên môn định kỳ (2 tuần/lần theo Điều lệ).
   - Lưu trữ biên bản, hình ảnh, video clip minh họa tiết dạy, phản hồi chia sẻ kinh nghiệm của đồng nghiệp.

### Phân hệ 4: Quản lý Công tác Bán trú & Dinh dưỡng Học đường
1. **Quản lý xuất ăn & Thực đơn**:
   - Lập thực đơn tuần/tháng, tự động tính toán cân bằng năng lượng (Kcal), tỷ lệ P:L:G (Protein - Lipid - Glucid) theo tháp dinh dưỡng lứa tuổi tiểu học của Viện Dinh dưỡng Quốc gia.
   - Tự động trừ kho nguyên liệu thực phẩm, kiểm soát lưu mẫu thức ăn 24h theo chuẩn an toàn vệ sinh thực phẩm (ATVSTP).
2. **Điểm danh bán trú & Chăm sóc giấc ngủ trưa**:
   - Điểm danh ăn sáng, ăn trưa, ăn xế bằng nhận diện khuôn mặt hoặc GVCN xác nhận qua Mobile App trước 8h30 hàng ngày để báo bếp chuẩn bị suất ăn chính xác, chống lãng phí.
   - Ghi nhận theo dõi sức khỏe đột xuất trong giờ bán trú.

### Phân hệ 5: Quản lý Y tế học đường & Sức khỏe học sinh
1. **Hồ sơ sức khỏe điện tử**:
   - Theo dõi chỉ số BMI (chiều cao, cân nặng), tầm soát thị lực (cận thị học đường), nha học đường, bệnh lý bẩm sinh.
   - Lịch sử tiêm chủng mở rộng và các đợt tiêm chủng quốc gia.
2. **Nhật ký phòng y tế**:
   - Cấp phát thuốc theo đơn phụ huynh gửi qua App, ghi nhận sơ cấp cứu ban đầu tại trường, tích hợp thông báo đẩy khẩn cấp cho phụ huynh khi học sinh có biểu hiện sốt/ốm.

### Phân hệ 6: Quản lý Thiết bị, Phòng chức năng & Thư viện số
1. **Thiết bị & Phòng học chuyên dùng**:
   - Quản lý danh mục thiết bị dạy học tối thiểu theo Thông tư của Bộ GD&ĐT.
   - Đăng ký mượn/trả đồ dùng dạy học trực tuyến; mượn phòng Tin học, Ngoại ngữ, Âm nhạc, Thể chất.
2. **Thư viện chuẩn quốc gia (Thông tư 16/2022/TT-BGDĐT)**:
   - Quản lý mã vạch sách, thẻ đọc học sinh/giáo viên.
   - Triển khai góc đọc sách điện tử, số hóa tư liệu tham khảo.

### Phân hệ 7: Quản lý Tài chính, Khoản thu Bán trú & Cổng thanh toán
1. **Quản lý danh mục khoản thu**:
   - Thiết lập cấu hình khoản thu thỏa thuận (tiền ăn bán trú, chăm sóc bán trú, nước uống...) và khoản thu hộ - chi hộ tuân thủ theo Nghị quyết HĐND cấp tỉnh/thành phố.
   - Tính toán bù trừ tiền ăn tự động dựa trên số ngày nghỉ học có phép/báo hủy ăn hợp lệ của học sinh.
2. **Hóa đơn điện tử & Cổng thanh toán không dùng tiền mặt**:
   - Tích hợp VietQR, VNPay, MoMo, ngân hàng số.
   - Tự động gạch nợ tức thì khi phụ huynh chuyển khoản. Xuất hóa đơn điện tử tự động kết nối Tổng cục Thuế.

---

## 5. PHÂN HỆ CHUYÊN SÂU: ĐÁNH GIÁ HIỆU SUẤT & BỘ CHỈ SỐ KPI GIÁO DỤC

Khác với doanh nghiệp thuần túy chỉ đo lường doanh số/lợi nhuận, bộ công cụ KPI trong trường tiểu học phải kết hợp hài hòa giữa **Quy chuẩn chức danh nghề nghiệp nhà nước (Nghị định 90/2020, NĐ 48/2023, Thông tư 20/2018)** với **Hiệu quả thực tế (Performance), Chất lượng sư phạm, Mức độ ứng dụng công nghệ và Kỷ luật lao động**.

### 5.1. Cơ chế vận hành bộ công cụ KPI (KPI Engine)
1. **Phương pháp đa chiều (Multi-source 360-degree Feedback)**:
   - **Chỉ số hệ thống đo lường tự động (System Automated Metrics - 40%)**: Tỷ lệ vào điểm/nhận xét đúng hạn, số tiết thao giảng, tỷ lệ duyệt giáo án đúng hạn, số lượng học sinh duy trì chuyên cần, sử dụng thiết bị dạy học.
   - **Đánh giá của cấp quản lý (BGH, Khối trưởng - 30%)**: Chất lượng chuyên môn tiết dạy, tham gia hoạt động tổ, thực hiện nhiệm vụ phân công đột xuất.
   - **Tự đánh giá & Đồng nghiệp đánh giá (Self & Peer Review - 15%)**: Tinh thần phối hợp, đóng góp vào bài giảng chung.
   - **Khảo sát phản hồi từ Phụ huynh học sinh (Parent CSAT - 15%)**: Mức độ nhiệt tình hỗ trợ, thông tin kịp thời về tình hình học sinh, ứng xử văn hóa sư phạm (ẩn danh để đảm bảo khách quan).
2. **Chu kỳ đánh giá linh hoạt**:
   - Theo Tháng (dùng để bình bầu thi đua, chi trả thu nhập tăng thêm/thưởng hiệu suất theo Nghị quyết/quy chế chi tiêu nội bộ).
   - Theo Học kỳ & Năm học (tổng hợp trực tiếp thành bảng Xếp loại viên chức theo Nghị định 90/2020/NĐ-CP và Chuẩn nghề nghiệp GV theo TT 20/2018/TT-BGDĐT).

---

### 5.2. Khung chỉ số KPI chi tiết theo từng nhóm vị trí công tác

#### A. Nhóm Giáo viên Chủ nhiệm (GVCN) tiểu học

| Mã KPI | Tiêu chí đánh giá | Cách đo lường / Công thức | Trọng số | Nguồn thu thập dữ liệu |
|---|---|---|---|---|
| **KPI-GV-01** | **Chất lượng tiến độ hồ sơ & sổ sách điện tử** | Tỷ lệ nhận xét học sinh (TT27), nộp giáo án, cập nhật điểm danh đúng hạn trên phần mềm | 20% | Hệ thống tự động ghi log thời gian submit |
| **KPI-GV-02** | **Chất lượng giáo dục & sự tiến bộ của học sinh** | Tỷ lệ học sinh hoàn thành chương trình lớp học; mức độ cải thiện của nhóm học sinh cần hỗ trợ đặc biệt | 20% | Dữ liệu tổng hợp đánh giá định kỳ TT27 |
| **KPI-GV-03** | **Đổi mới phương pháp & chất lượng giờ dạy** | Điểm trung bình các tiết dự giờ, thao giảng; ứng dụng CNTT, đồ dùng dạy học số hóa | 15% | Phiếu chấm tiết dạy của BGH & Tổ chuyên môn |
| **KPI-GV-04** | **Quản lý nề nếp lớp & phong trào Sao Nhi Đồng** | Điểm thi đua tuần của lớp do Tổng phụ trách chấm; tỷ lệ chuyên cần của học sinh | 15% | Phân hệ Công tác Đội & Điểm danh |
| **KPI-GV-05** | **Tương tác & Mức độ hài lòng của Phụ huynh** | Tỷ lệ phản hồi tin nhắn sổ liên lạc trong 24h; Điểm CSAT từ đợt khảo sát ẩn danh định kỳ | 15% | Module Liên lạc & Khảo sát PHHS |
| **KPI-GV-06** | **Nghiên cứu khoa học & Nâng cao nghiệp vụ** | Hoàn thành sáng kiến/đề tài khoa học/giải pháp đổi mới sáng tạo hoặc SKKN được công nhận; hoàn thành các mô-đun bồi dưỡng thường xuyên của BGDĐT; tham gia hội thi GVDG | 10% | Module Hồ sơ nhân sự / Nghiệm thu sáng kiến |
| **KPI-GV-07** | **Kỷ luật lao động & Văn hóa học đường** | Đảm bảo giờ giấc lên lớp, chấp hành quy chế chi tiêu nội bộ, tham gia hoạt động đoàn thể | 5% | Chấm công khuôn mặt & Ghi nhận từ Công đoàn |

---

#### B. Nhóm Giáo viên Bộ môn (Tiếng Anh, Tin học, Thể dục, Âm nhạc, Mỹ thuật)

| Mã KPI | Tiêu chí đánh giá | Cách đo lường / Công thức | Trọng số | Nguồn dữ liệu |
|---|---|---|---|---|
| **KPI-BM-01** | **Kỷ luật chuyên môn & Hoàn thành tiến độ** | Tỷ lệ nhập nhận xét, điểm số đúng thời hạn quy định trên toàn bộ các lớp phụ trách | 25% | Hệ thống tự động kiểm tra Deadline |
| **KPI-BM-02** | **Tỷ lệ học sinh đạt chuẩn đầu ra môn học** | Tỷ lệ học sinh đạt mức Hoàn thành tốt (HTT) và Hoàn thành (HT); phát hiện bồi dưỡng HS có năng khiếu | 25% | Sổ đánh giá môn học TT27 |
| **KPI-BM-03** | **Tổ chức hoạt động trải nghiệm / Câu lạc bộ** | Số lượng CLB (tiếng Anh, cờ vua, vẽ, bơi lội...), chất lượng các sự kiện ngoại khóa | 20% | BGH nghiệm thu kế hoạch CLB |
| **KPI-BM-04** | **Chất lượng tiết dạy & Thiết bị dạy học** | Điểm dự giờ chuyên đề, tần suất khai thác phòng chức năng/thiết bị chuyên dùng | 20% | Nhật ký mượn phòng & Đánh giá của Tổ trưởng |
| **KPI-BM-05** | **Bồi dưỡng chuyên môn & Đóng góp phong trào** | Tham gia hội đồng bộ môn cấp quận/huyện, bồi dưỡng học sinh thi đấu phong trào thể thao/nghệ thuật | 10% | Hồ sơ thành tích trường học |

---

#### C. Nhóm Cán bộ Quản lý (Hiệu trưởng, Phó Hiệu trưởng)
Được cấu hình bám sát **Thông tư 14/2018/TT-BGDĐT** về Chuẩn hiệu trưởng:
1. **Tiêu chuẩn 1: Phẩm chất nghề nghiệp & Đạo đức lối sống (Trọng số 15%)**.
2. **Tiêu chuẩn 2: Quản trị nhà trường (Trọng số 35%)**:
   - Quản trị hoạt động dạy học và giáo dục học sinh (kết quả chung toàn trường).
   - Quản trị nhân sự: Mức độ hoàn thành KPI trung bình của toàn thể giáo viên, nhân viên.
   - Quản trị tài chính, cơ sở vật chất và ứng dụng CNTT (tỷ lệ chuyển đổi số, an toàn trường học, chất lượng kiểm định giáo dục).
3. **Tiêu chuẩn 3: Xây dựng môi trường giáo dục an toàn, lành mạnh, thân thiện (Trọng số 20%)**: Không xảy ra tai nạn thương tích, bạo lực học đường, ngộ độc thực phẩm bán trú.
4. **Tiêu chuẩn 4: Phát triển mối quan hệ giữa nhà trường, gia đình và xã hội (Trọng số 15%)**: Tỷ lệ hài lòng chung của phụ huynh toàn trường $\ge 90\%$.
5. **Tiêu chuẩn 5: Sử dụng ngoại ngữ và công nghệ thông tin trong quản lý (Trọng số 15%)**.

---

#### D. Nhóm Nhân viên Hành chính - Hỗ trợ (Vị trí việc làm theo TT 20/2023/TT-BGDĐT)

##### 1. Kế toán & Quản lý Thu chi
- **KPI-KT-01**: Tỷ lệ thu học phí/tiền ăn đúng hạn qua cổng thanh toán không dùng tiền mặt ($\ge 95\%$).
- **KPI-KT-02**: Độ chính xác và thời gian quyết toán các khoản chi trả phụ cấp, tăng thêm, bảo hiểm xã hội.
- **KPI-KT-03**: Báo cáo tài chính, quyết toán ngân sách nhà nước đúng hạn, không sai sót số liệu.
- **KPI-KT-04**: Giải quyết thắc mắc về tiền học của phụ huynh trong vòng 24 giờ làm việc.

##### 2. Y tế học đường
- **KPI-YT-01**: 100% học sinh được khám sức khỏe định kỳ và cập nhật hồ sơ BMI trên hệ thống.
- **KPI-YT-02**: Kiểm soát an toàn vệ sinh thực phẩm bếp ăn bán trú (100% lưu mẫu thức ăn đúng quy chuẩn, 0 vụ ngộ độc).
- **KPI-YT-03**: Tốc độ xử lý và nhật ký sơ cấp cứu kịp thời cho học sinh tại trường.
- **KPI-YT-04**: Tuyên truyền phòng chống dịch bệnh theo mùa cho học sinh và giáo viên.

##### 3. Văn thư & CNTT
- **KPI-VT-01**: Tỷ lệ văn bản đi/đến được số hóa, gắn mã QR và luân chuyển xử lý trong ngày.
- **KPI-VT-02**: Tỷ lệ đồng bộ dữ liệu thông suốt lên CSDL ngành của Bộ GD&ĐT (100% đúng hạn kỳ đầu năm, giữa năm, cuối năm).
- **KPI-VT-03**: Duy trì hệ thống mạng nội bộ, phòng máy tính, hệ thống camera học đường hoạt động ổn định ($\ge 99.5\%$).

##### 4. Cán bộ Thư viện - Thiết bị
- **KPI-TV-01**: Tỷ lệ sách được số hóa, dán mã vạch quản lý trên hệ thống.
- **KPI-TV-02**: Số lượt mượn trả sách của học sinh và giáo viên trong tháng (đo lường văn hóa đọc).
- **KPI-TV-03**: Chuẩn bị đầy đủ thiết bị, dụng cụ dạy học cho giáo viên theo phiếu đăng ký trước 24h.

---

### 5.3. Quy trình tính toán & Liên kết xếp loại thi đua tự động

Hệ thống cung cấp thuật toán chuyển đổi tự động từ **Điểm số KPI thực tế** sang **Mức xếp loại theo Nghị định 90/2020/NĐ-CP & NĐ 48/2023/NĐ-CP**:

$$\text{Điểm Tổng Hợp KPI} = \sum_{i=1}^{n} (\text{Điểm Chỉ Số}_i \times \text{Trọng Số}_i)$$

**Quy tắc xếp loại:**
- **Hoàn thành xuất sắc nhiệm vụ**: Điểm KPI $\ge 90$ điểm; đồng thời không có chỉ số bắt buộc nào dưới điểm liệt; có ít nhất 01 sáng kiến/đề tài khoa học/giải pháp đổi mới sáng tạo áp dụng hiệu quả hoặc có thành tích xuất sắc tiêu biểu được cấp có thẩm quyền công nhận (theo Luật Thi đua Khen thưởng 2022 & NĐ 98/2023/NĐ-CP, không bắt buộc chỉ duy nhất SKKN). Tỷ lệ không vượt quá mức tối đa theo quy định của cơ quan chủ quản.
- **Hoàn thành tốt nhiệm vụ**: Điểm KPI từ $75$ đến dưới $90$ điểm; hoàn thành đầy đủ các chỉ số nghiệp vụ chuyên môn.
- **Hoàn thành nhiệm vụ**: Điểm KPI từ $50$ đến dưới $75$ điểm; có tối đa 2 chỉ số không đạt nhưng không gây hậu quả nghiêm trọng.
- **Không hoàn thành nhiệm vụ**: Điểm KPI $< 50$ điểm hoặc vi phạm đạo đức nhà giáo, có hành vi bạo lực học đường, để xảy ra tai nạn thương tích nghiêm trọng trong giờ quản lý.

---

## 6. KHẢ NĂNG TÍCH HỢP & LIÊN THÔNG DỮ LIỆU DÙNG CHUNG QUỐC GIA VÀ BỘ GD&ĐT

Hệ thống được thiết kế tuân thủ **Khung Kiến trúc Chính phủ điện tử Việt Nam phiên bản 3.0 (Quyết định số 2568/QĐ-BTTTT ngày 29/12/2023)**, **Luật Giao dịch điện tử số 20/2023/QH15** và **Chiến lược dữ liệu quốc gia đến năm 2030 (Quyết định số 142/QĐ-TTg ngày 02/02/2024)**.

```
                                  +-------------------------------------------------------------+
                                  |   NỀN TẢNG TÍCH HỢP, CHIA SẺ DỮ LIỆU QUỐC GIA (NDXP / LGSP)  |
                                  |              (Nghị định số 47/2020/NĐ-CP)                   |
                                  +-------------------------------------------------------------+
                                           |                        |                       |
                 +-------------------------+                        |                       +-------------------------+
                 v                                                  v                                                 v
+----------------------------------+     +-----------------------------------+     +----------------------------------+
|  CSDL QUỐC GIA VỀ DÂN CƯ (C06)   |     |     CSDL NGÀNH GIÁO DỤC (MOET)     |     |   CSDLQG VỀ CÁN BỘ VIÊN CHỨC     |
| (Luật Căn cước 2023 & NĐ 69/2024)|     |  (TT 42/2021 & QĐ 4998/QĐ-BGDĐT)  |     |      (Thông tư 06/2023/TT-BNV)   |
| • Xác thực định danh VNeID Mức 2 |     | • CSDL Trường, Lớp, Học sinh      |     | • Hồ sơ chức danh nghề nghiệp    |
| • Căn cước trẻ em dưới 14 tuổi   |     | • Kho Học bạ số quốc gia (QĐ 391) |     | • Kết quả đánh giá xếp loại KPI  |
| • Làm sạch hồ sơ tuyển sinh lớp 1|     | • Báo cáo thống kê chỉ số EMIS/IOC|     |   viên chức (Nghị định 90 & 48)  |
+----------------------------------+     +-----------------------------------+     +----------------------------------+
                 ^                                                  ^                                                 ^
                 |                                                  |                                                 |
                 +--------------------------------------------------+-------------------------------------------------+
                                                                    | Trục API Gateway an toàn (TLS 1.3, mTLS, OAuth2)
                                                +---------------------------------------+
                                                |  HỆ THỐNG PHẦN MỀM EDUPRIMARY-KPI GO  |
                                                +---------------------------------------+
                                                                    |
                                 +----------------------------------+----------------------------------+
                                 v                                                                     v
+----------------------------------------------------+       +----------------------------------------------------+
|  CỔNG DVC QUỐC GIA & THANH TOÁN KHÔNG DÙNG TIỀN MẶT |       |         HỆ THỐNG KÝ SỐ & DẤU THỜI GIAN QUỐC GIA    |
| (Quyết định 1813/QĐ-TTg & Nghị định 123/2020/NĐ-CP) |       |            (Luật Giao dịch điện tử 2023)           |
| • Cổng Dịch vụ công Quốc gia (DVCQG)               |       | • Ký số tập trung SmartCA / Cloud CA BGH & GV      |
| • VietQR động định danh học sinh / Cổng Napas      |       | • Đóng dấu thời gian Timestamp Authority (TSA)     |
| • Kết nối Hóa đơn điện tử e-Invoice Tổng cục Thuế  |       | • Xác thực tính toàn vẹn học bạ XML/JSON chuẩn QĐ 391|
+----------------------------------------------------+       +----------------------------------------------------+
```

### 6.1. Kết nối Cơ sở dữ liệu quốc gia về Dân cư & Định danh điện tử (Đề án 06)
- **Căn cứ pháp lý mới**: **Luật Căn cước số 26/2023/QH15** (cấp căn cước cho trẻ dưới 14 tuổi) và **Nghị định số 69/2024/NĐ-CP ngày 25/06/2024** quy định về định danh và xác thực điện tử (*thay thế Nghị định 59/2022/NĐ-CP*).
- **Cơ chế triển khai**:
  - Tích hợp cổng đăng nhập một lần (SSO) qua tài khoản định danh điện tử **VNeID**. Phụ huynh sử dụng tài khoản định danh Mức độ 2 để xác thực mối quan hệ nhân thân với con em học tiểu học theo quy định tại Điều 9 Nghị định 69/2024.
  - Tự động làm sạch dữ liệu tuyển sinh đầu cấp Lớp 1: Kiểm tra đối khớp Số định danh cá nhân (12 số) của học sinh với CSDLQG về Dân cư, loại bỏ hoàn toàn tình trạng sai lệch ngày tháng năm sinh hoặc trùng lặp hồ sơ.

### 6.2. Kết nối Trục Cơ sở dữ liệu ngành Giáo dục và Đào tạo (MOET API)
- **Căn cứ pháp lý**: **Thông tư số 42/2021/TT-BGDĐT** quy định về CSDL GD&ĐT và **Quyết định số 4998/QĐ-BGDĐT** về quy định kỹ thuật cấu trúc dữ liệu.
- **Cơ chế đồng bộ hai chiều**:
  - Chuẩn hóa danh mục dùng chung của Bộ GD&ĐT: Mã cơ sở giáo dục, mã lớp học, mã khối, mã môn học, mã diện chính sách, mã khuyết tật.
  - Đồng bộ kỳ đầu năm (tháng 9), giữa năm (tháng 1) và cuối năm học (tháng 5): Tự động đẩy danh sách học sinh, phổ cập giáo dục tiểu học, hồ sơ cán bộ giáo viên, kết quả đánh giá theo Thông tư 27/2020/TT-BGDĐT lên hệ thống `csdl.moet.gov.vn`.
  - Kết nối dữ liệu vào Trung tâm điều hành thông minh giáo dục (IOC) cấp Sở/Phòng GD&ĐT.

### 6.3. Kết nối Kho Học bạ số quốc gia & Hệ thống Chữ ký số Chuyên dùng (Bộ GD&ĐT)
- **Căn cứ pháp lý**: 
  - **Quyết định số 391/QĐ-BGDĐT ngày 24/01/2024** về Kế hoạch triển khai Học bạ số cấp tiểu học.
  - **Thông tư số 89/2026/TT-BGDĐT ngày 16/10/2026** (chính thức pháp lý hóa Học bạ số thay thế học bạ giấy).
  - **Nghị định số 130/2018/NĐ-CP** và **Thông tư số 185/2019/TT-BQP** về quản lý, sử dụng chữ ký số chuyên dùng công vụ và dịch vụ chứng thực chữ ký số công cộng.
- **Tiêu chuẩn kỹ thuật**:
  - Xuất dữ liệu học bạ theo cấu trúc XML/JSON chuẩn được Bộ GD&ĐT phê duyệt.
  - Tích hợp chữ ký số cá nhân của Giáo viên chủ nhiệm và chữ ký số cơ quan của Ban Giám hiệu theo chuẩn PKCS#7 / PAdES (Luật Giao dịch điện tử số 20/2023/QH15). Hỗ trợ cả 2 luồng: Chữ ký số chuyên dùng công vụ (Ban Cơ yếu Chính phủ) và Chữ ký số từ xa (Remote Signing / SmartCA công cộng).
  - Tự động gửi gói tin học bạ số đã ký lên Kho lưu trữ học bạ số quốc gia của Bộ GD&ĐT, phục vụ liên thông chuyển trường toàn quốc và tuyển sinh chuyển cấp THCS.

### 6.4. Kết nối CSDL quốc gia về Cán bộ, công chức, viên chức (Bộ Nội vụ)
- **Căn cứ pháp lý**: **Quyết định số 893/QĐ-TTg** và **Thông tư số 06/2023/TT-BNV** ban hành Quy chế cập nhật, sử dụng, khai thác CSDLQG về CBCCVC.
- **Cơ chế liên thông**:
  - Đồng bộ 109 trường thông tin hồ sơ viên chức ngành giáo dục (vị trí việc làm theo TT 20/2023, hệ số lương, chứng chỉ chức danh theo TT 08/2023, đào tạo bồi dưỡng thường xuyên).
  - Kết xuất và đồng bộ kết quả đánh giá, xếp loại thi đua viên chức hàng năm từ hệ thống KPI sang CSDLQG về CBCCVC của Bộ Nội vụ.

### 6.5. Kết nối Cổng Dịch vụ công Quốc gia & Nền tảng Thanh toán không dùng tiền mặt
- **Căn cứ pháp lý**: **Quyết định số 1813/QĐ-TTg** về thanh toán không dùng tiền mặt; **Nghị định số 123/2020/NĐ-CP** và **Thông tư số 78/2021/TT-BTC** về hóa đơn điện tử.
- **Cơ chế triển khai**:
  - Tích hợp thanh toán dịch vụ công trực tuyến trên Cổng Dịch vụ công Quốc gia (DVCQG) cho các khoản lệ phí, học phí.
  - Tạo mã **VietQR động** cho từng học sinh và từng đợt thu (tiền ăn bán trú, chăm sóc bán trú, đồng phục...). Khi phụ huynh chuyển khoản qua bất kỳ app ngân hàng nào, hệ thống nhận Webhook qua Open API và thực hiện gạch nợ tức thì ($< 3\text{ giây}$).
  - Tự động phát hành Hóa đơn điện tử có mã của cơ quan Thuế khởi tạo từ hệ thống, gửi hóa đơn điện tử định dạng PDF/XML qua tin nhắn App cho phụ huynh.

---

## 7. YÊU CẦU PHI CHỨC NĂNG & BẢO MẬT DỮ LIỆU TRẺ EM

1. **Bảo mật Dữ liệu Trẻ em (Nghị định 13/2023/NĐ-CP - Điều 20) & Quy trình Đồng thuận số (Digital Consent)**:
   - Dữ liệu học sinh tiểu học là dữ liệu trẻ em dưới 14 tuổi theo Luật Căn cước 2023. Hệ thống bắt buộc triển khai quy trình **Digital Consent**: Thu thập và ghi nhật ký sự đồng ý rõ ràng của cha, mẹ hoặc người giám hộ (thông qua tài khoản định danh VNeID hoặc xác thực OTP SMS/App) trước khi kích hoạt hồ sơ học sinh, chụp ảnh nhận diện bán trú hoặc số hóa thông tin sức khỏe.
   - Toàn bộ dữ liệu nhạy cảm của trẻ em (ảnh chân dung, hồ sơ bệnh án y tế, nhận xét năng lực, hoàn cảnh gia đình) được mã hóa tại tầng lưu trữ (Encryption at Rest - AES-256) và trên đường truyền (TLS 1.3).
   - Cam kết tuyệt đối không sử dụng dữ liệu học sinh vào mục đích thương mại, quảng cáo hoặc chia sẻ với bên thứ ba khi chưa có sự đồng ý của cha mẹ.
   - Cơ chế ẩn danh hóa dữ liệu (Anonymization) khi trích xuất số liệu phân tích thống kê báo cáo lên cấp quản lý.
2. **Kiểm soát phân quyền & Nhật ký vết (Audit Logging)**:
   - Áp dụng mô hình RBAC (Role-Based Access Control) kết hợp ABAC (Attribute-Based Access Control - Giáo viên chỉ xem được lớp mình dạy).
   - Ghi nhật ký bất biến (Immutable Logs) cho mọi hành vi sửa điểm, sửa nhận xét, thay đổi kết quả KPI.
3. **Hiệu năng & Khả năng chịu tải**:
   - Phản hồi tác vụ nhập nhận xét, tra cứu $< 1.5$ giây.
   - Chịu tải cao điểm vào các đợt cuối kỳ khi toàn bộ giáo viên vào điểm và phụ huynh truy cập xem học bạ đồng thời ($\ge 10.000$ concurrent users).

---

## 8. LỘ TRÌNH TRIỂN KHAI & CÔNG NGHỆ KHUYẾN NGHỊ

### 8.1. Tech Stack khuyến nghị
- **Backend Core**: **Go (Golang)** áp dụng kiến trúc Clean Architecture & Domain-Driven Design (DDD). Siêu nhẹ, hiệu năng cao, tiết kiệm tài nguyên (RAM chỉ $\sim 30 - 60\text{ MB}$), bảo mật biên dịch nhị phân.
- **Admin Engine quản trị hệ thống**: Tích hợp **GoAdmin Engine (`github.com/GoAdminGroup/go-admin`)** trực tiếp trong binary của Go, cung cấp giao diện quản trị cấu hình bảng biểu, phân quyền RBAC, xuất nhập dữ liệu tự động tương tự Django Admin.
- **Cơ chế tái sử dụng đa cấp học (Edu-Core Engine)**: Sử dụng mô hình *Strategy Pattern*, tách biệt Core dùng chung (Nhân sự, Học sinh, Bán trú, Thu phí, KPI Engine) và các Plugin chiến lược theo cấp học (Tiểu học theo TT 27 & TT 05/2025; THCS/THPT theo TT 22 & TT 05/2025; Mầm non theo TT 51).
- **Frontend Portal**: React / Next.js (Dành cho Portal Giáo viên, Khối trưởng, Ban Giám hiệu tra cứu & nhận xét).
- **Mobile App**: Flutter (Đa nền tảng iOS & Android cho Phụ huynh và Giáo viên điểm danh, nhận xét nhanh).
- **Database**: PostgreSQL (CSDL chính ACID), Redis (Caching & Queue), MinIO/S3 (Hồ sơ học bạ số, giáo án).
- **Hạ tầng**: Docker container trên nền tảng Cloud đạt chuẩn an toàn thông tin cấp độ 3 (Viettel IDC, VNPT Cloud, FPT Cloud). Chi tiết thiết kế xem tại file: `KIEN_TRUC_CORE_GO_VA_ADMIN.md`.

### 8.2. Kế hoạch triển khai theo 4 giai đoạn

```
Giai đoạn 1 (Tháng 1-2)    Giai đoạn 2 (Tháng 3-4)    Giai đoạn 3 (Tháng 5-6)    Giai đoạn 4 (Tháng 7+)
+----------------------+  +----------------------+  +----------------------+  +----------------------+
| - Khảo sát nghiệp vụ |  | - Phân hệ TT27/2020  |  | - Phát triển KPI     |  | - Mở rộng toàn hệ    |
| - Thiết kế Database  |  | - Giáo án & Hồ sơ    |  |   Engine đa chiều    |  |   thống trường học   |
| - Quản lý Học sinh & |  | - Bán trú, Dinh dưỡng|  | - Ký số học bạ số    |  | - Tinh chỉnh chỉ số  |
|   Cán bộ giáo viên   |  | - App Phụ huynh      |  | - Kết nối CSDL ngành |  |   KPI theo thực tế   |
+----------------------+  +----------------------+  +----------------------+  +----------------------+
```

---

## 9. QUẢN TRỊ RỦI RO & PHƯƠNG ÁN AN TOÀN THÔNG TIN PHỤC VỤ THẨM ĐỊNH

Theo quy định quản lý đầu tư CNTT (Nghị định 73/2019/NĐ-CP, Nghị định 82/2024/NĐ-CP) và Tiêu chuẩn An toàn thông tin Cấp độ 3 (Nghị định 85/2016/NĐ-CP, Thông tư 12/2022/TT-BTTTT), dự án thiết lập ma trận quản trị rủi ro như sau:

| STT | Rủi ro tiềm ẩn | Mức độ | Biện pháp kiểm soát & Giải pháp kỹ thuật |
|:---:|---|:---:|---|
| 1 | **Rủi ro thay đổi quy định pháp luật (Thông tư, Nghị định mới)** | Cao | Áp dụng kiến trúc **Temporal Policy & Strategy Versioning** (chi tiết tại `KIEN_TRUC_THICH_UNG_THAY_DOI_QUY_DINH.md`). Dữ liệu lịch sử được đóng băng bất biến bằng Snapshot có ký số, không bị ảnh hưởng khi luật đổi. |
| 2 | **Rủi ro rò rỉ dữ liệu trẻ em (Nghị định 13/2023/NĐ-CP)** | Cao | Mã hóa CSDL cấp độ cột (AES-256), mã hóa đường truyền TLS 1.3; kiểm soát truy cập thuộc tính (ABAC); ẩn danh hóa dữ liệu khi xuất báo cáo thống kê. |
| 3 | **Rủi ro nghẽn mạng giờ cao điểm (Cuối kỳ vào điểm & xem học bạ)** | Trung bình | Sử dụng Go Backend hiệu năng cao, kết hợp Redis Cache cho dữ liệu tra cứu và Message Queue (Kafka/RabbitMQ) cho các tác vụ ký số, xuất file PDF. |
| 4 | **Rủi ro giáo viên khó tiếp cận công nghệ mới** | Trung bình | Thiết kế UI/UX tinh gọn, hỗ trợ công cụ gợi ý nhận xét thông minh theo TT 27, cung cấp ứng dụng Mobile thân thiện và tổ chức đào tạo chuyển giao theo khối/tổ chuyên môn. |
| 5 | **Rủi ro mất kết nối đồng bộ CSDL Ngành (MOET API)** | Thấp | Xây dựng cơ chế *Outbox Pattern & Retry Queue*; tự động lưu trữ tạm cục bộ và kích hoạt đồng bộ lại khi kết nối internet/API khôi phục. |

---

## 10. DANH MỤC TÀI LIỆU TRONG HỒ SƠ THẨM ĐỊNH DỰ ÁN

Bộ hồ sơ hoàn chỉnh trình Hội đồng thẩm định gồm 6 tài liệu đồng bộ và thư viện văn bản pháp lý gốc:
1. **Tập 1 - Báo cáo đề xuất tính năng & nghiệp vụ**: `DE_XUAT_CHUC_NANG_QLTH_KPI.md` (Tài liệu này).
2. **Tập 2 - Thiết kế kiến trúc Go Backend, Admin Engine & Core đa cấp học**: `KIEN_TRUC_CORE_GO_VA_ADMIN.md`.
3. **Tập 3 - Phương án kỹ thuật bảo toàn dữ liệu & thích ứng thay đổi thể chế**: `KIEN_TRUC_THICH_UNG_THAY_DOI_QUY_DINH.md`.
4. **Tập 4 - Phân tích đối thủ cạnh tranh & Đề xuất tính năng đột phá**: `PHAN_TICH_DOI_THU_VA_DE_XUAT_TINH_NANG_DOT_PHA.md`.
5. **Tập 5 - Phân tích tiềm năng thị trường, chiến lược giá SaaS & hiệu quả đầu tư**: `PHAN_TICH_HIEU_QUA_DAU_TU_VA_MO_HINH_GIA_SAAS.md`.
6. **Tập 6 - Đặc tả thiết kế kỹ thuật chi tiết (Database Schema, API Endpoints & Thuật toán Golang)**: `docs/phap_ly/DAC_TA_THIET_KE_CHI_TIET_HE_THONG.md`.
7. **Thư viện pháp lý gốc đính kèm**: Thư mục `docs/phap_ly/` lưu trữ toàn văn các tệp PDF gốc của Thông tư 05/2025/TT-BGDĐT, Thông tư 27/2020/TT-BGDĐT, Thông tư 28/2020/TT-BGDĐT, Thông tư 29/2024/TT-BGDĐT, Nghị định 69/2024/NĐ-CP và Báo cáo phân tích chuyên sâu từng điều khoản (`TONG_HOP_QUY_DINH_PHAP_LY_CHI_TIET.md`).
