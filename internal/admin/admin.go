package admin

import (
	"html/template"
	"os"
	"sync"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	_ "github.com/GoAdminGroup/themes/adminlte"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gin-gonic/gin"

	"qlth/internal/admin/tables"
	qlthdb "qlth/internal/db"
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

var (
	adminEngine   *engine.Engine
	adminInitLock sync.Mutex
)

// InitAdmin initializes GoAdmin engine with PostgreSQL and mounts it to Gin router.
func InitAdmin(r *gin.Engine) (*engine.Engine, error) {
	adminInitLock.Lock()
	defer adminInitLock.Unlock()

	if adminEngine != nil {
		if len(tables.Generators) > 0 {
			adminEngine.AdminPlugin().AddGenerators(tables.Generators)
		}
		// Already initialized globally; bind routes to the new gin router instance
		if err := adminEngine.Adapter.Use(r, adminEngine.PluginList); err != nil {
			return adminEngine, err
		}
		adminEngine.HTML("GET", "/admin", DashboardPanel)
		adminEngine.HTML("GET", "/admin/timetable", TimetablePanel)
		adminEngine.HTML("GET", "/admin/integrations", IntegrationsPanel)
		return adminEngine, nil
	}

	dbCfg := config.Database{
		Host:         getEnv("PGHOST", "127.0.0.1"),
		Port:         getEnv("PGPORT", "5432"),
		User:         getEnv("PGUSER", "qlth_user"),
		Pwd:          getEnv("PGPASSWORD", "qlth_password_secure"),
		Name:         getEnv("PGDATABASE", "qlth_eduprimary"),
		Driver:       config.DriverPostgresql,
		MaxIdleConns: 25,
		MaxOpenConns: 50,
		Params: map[string]string{
			"sslmode": "disable",
		},
	}

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		dbCfg.Dsn = dsn
	}

	cfg := &config.Config{
		Databases: config.DatabaseList{
			"default": dbCfg,
		},
		UrlPrefix:    "admin",
		IndexUrl:     "/",
		Theme:        "adminlte",
		Title:        "EduPrimary-KPI Quản Trị Hệ Thống",
		LoginTitle:   "EduPrimary-KPI Quản Trị",
		Logo:         template.HTML("<b>EduPrimary</b>-KPI"),
		MiniLogo:     template.HTML("<b>EP</b>"),
		ColorScheme:  "skin-black",
		Debug:        os.Getenv("GIN_MODE") != "release",
		Language:     language.EN,
		AccessLogOff: true,
		InfoLogOff:   true,
		ErrorLogOff:  false,
	}

	eng := engine.Default()

	if len(tables.Generators) > 0 {
		eng.AddGenerators(tables.Generators)
	}

	if err := eng.AddConfig(cfg).Use(r); err != nil {
		return eng, err
	}

	eng.HTML("GET", "/admin", DashboardPanel)
	eng.HTML("GET", "/admin/timetable", TimetablePanel)
	eng.HTML("GET", "/admin/integrations", IntegrationsPanel)

	ensureAdminMenus()

	adminEngine = eng
	return eng, nil
}

func ensureAdminMenus() {
	conn, err := qlthdb.GetDB()
	if err != nil || conn == nil {
		return
	}
	query := `
INSERT INTO goadmin_menu (id, parent_id, type, "order", title, header, plugin_name, icon, uri)
VALUES 
    (1, 0, 1, 1, 'Dashboard', NULL, '', 'fa-bar-chart', '/'),
    (2, 0, 1, 2, 'Admin', NULL, '', 'fa-tasks', ''),
    (3, 2, 1, 1, 'Users', NULL, '', 'fa-users', '/info/manager'),
    (4, 2, 1, 2, 'Roles', NULL, '', 'fa-user', '/info/roles'),
    (5, 2, 1, 3, 'Permission', NULL, '', 'fa-ban', '/info/permission'),
    (6, 2, 1, 4, 'Menu', NULL, '', 'fa-bars', '/menu'),
    (7, 2, 1, 5, 'Operation log', NULL, '', 'fa-history', '/info/op'),
    (8, 0, 1, 3, 'Giáo viên & Định mức', 'Quản lý Trường học & Sư phạm', '', 'fa-briefcase', '/info/staffs'),
    (11, 0, 1, 4, 'Lớp học', NULL, '', 'fa-sitemap', '/info/classes'),
    (19, 0, 1, 5, 'Thời khóa biểu & Tiết dạy', NULL, '', 'fa-calendar', '/timetable'),
    (9, 0, 1, 6, 'Học sinh & Consent NĐ13', NULL, '', 'fa-graduation-cap', '/info/students'),
    (12, 0, 1, 7, 'Kế hoạch bài dạy CV 2345', NULL, '', 'fa-book', '/info/lesson_plans'),
    (13, 0, 1, 8, 'Danh mục Môn học', 'Học tập & Đánh giá (TT27 / TT89)', '', 'fa-list-ul', '/info/subjects'),
    (14, 0, 1, 9, 'Đánh giá Môn học', NULL, '', 'fa-check-square-o', '/info/student_subject_evaluations'),
    (15, 0, 1, 10, 'Đánh giá Năng lực - Phẩm chất', NULL, '', 'fa-heart', '/info/student_behavior_evaluations'),
    (16, 0, 1, 11, 'Học bạ số & Chữ ký số', NULL, '', 'fa-certificate', '/info/academic_transcripts_snapshots'),
    (17, 0, 1, 12, 'Thực đơn Bán trú', 'Bán trú & Y tế Học đường', '', 'fa-cutlery', '/info/boarding_menus'),
    (18, 0, 1, 13, 'Hồ sơ Sức khỏe Học sinh', NULL, '', 'fa-medkit', '/info/student_health_records'),
    (10, 0, 1, 14, 'Đánh giá KPI 360', 'Thi đua & Hiệu suất', '', 'fa-star', '/info/kpi_evaluations'),
    (20, 0, 1, 15, 'Cổng Tích hợp & Sandbox', 'Tích hợp & Dịch vụ số', '', 'fa-exchange', '/integrations')
ON CONFLICT (id) DO UPDATE SET 
    title = EXCLUDED.title,
    uri = EXCLUDED.uri,
    header = EXCLUDED.header,
    icon = EXCLUDED.icon,
    "order" = EXCLUDED."order";

INSERT INTO goadmin_role_menu (role_id, menu_id)
SELECT 1, m.id FROM goadmin_menu m
WHERE NOT EXISTS (SELECT 1 FROM goadmin_role_menu rm WHERE rm.role_id = 1 AND rm.menu_id = m.id);

INSERT INTO goadmin_role_menu (role_id, menu_id)
SELECT 2, m.id FROM goadmin_menu m WHERE m.id IN (1, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)
AND NOT EXISTS (SELECT 1 FROM goadmin_role_menu rm WHERE rm.role_id = 2 AND rm.menu_id = m.id);
`
	_, _ = conn.Exec(query)
}

// DashboardPanel renders the main admin dashboard panel.
func DashboardPanel(ctx *context.Context) (types.Panel, error) {
	html := `
<div class="row">
    <div class="col-md-12">
        <div class="callout callout-info" style="margin-bottom: 15px;">
            <h4><i class="fa fa-graduation-cap"></i> Hệ thống Quản trị Trường Tiểu học EduPrimary-KPI</h4>
            <p>Hệ thống quản lý tuân thủ Thông tư 05/2025/TT-BGDĐT (Định mức 805 tiết/năm, giảm 4 tiết GVCN), Nghị định 90/2020 & 48/2023/NĐ-CP (Khống chế trần 20% xuất sắc, điểm liệt), Thông tư 27/2020 & 89/2026/TT-BGDĐT, và Điều 20 Nghị định 13/2023/NĐ-CP (Bảo vệ dữ liệu trẻ em).</p>
        </div>
    </div>
</div>
<div class="row">
    <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-aqua">
            <div class="inner">
                <h3>805<sup style="font-size: 20px">h</sup></h3>
                <p>Định mức Tiết dạy Chuẩn (TT 05/2025)</p>
            </div>
            <div class="icon">
                <i class="fa fa-book"></i>
            </div>
            <a href="/admin/info/staffs" class="small-box-footer">Quản lý Cán bộ & GV <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
    <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-green">
            <div class="inner">
                <h3>20<sup style="font-size: 20px">%</sup></h3>
                <p>Khống chế Trần Xuất sắc (NĐ 90/48)</p>
            </div>
            <div class="icon">
                <i class="fa fa-trophy"></i>
            </div>
            <a href="/admin/info/kpi_evaluations" class="small-box-footer">Đánh giá KPI 360 <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
    <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-yellow">
            <div class="inner">
                <h3>NĐ 13</h3>
                <p>Digital Consent & VNeID</p>
            </div>
            <div class="icon">
                <i class="fa fa-child"></i>
            </div>
            <a href="/admin/info/students" class="small-box-footer">Quản lý Học sinh <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
    <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-red">
            <div class="inner">
                <h3>RBAC</h3>
                <p>Quản trị Phân quyền Hệ thống</p>
            </div>
            <div class="icon">
                <i class="fa fa-users"></i>
            </div>
            <a href="/admin/info/manager" class="small-box-footer">Quản lý Người dùng <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
</div>
<div class="row">
    <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-purple">
            <div class="inner">
                <h3>CV 2345</h3>
                <p>Giáo án & Phê duyệt Phân cấp</p>
            </div>
            <div class="icon">
                <i class="fa fa-file-text"></i>
            </div>
            <a href="/admin/info/lesson_plans" class="small-box-footer">Kế hoạch Bài dạy <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
    <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-blue">
            <div class="inner">
                <h3>TT 89/26</h3>
                <p>Đánh giá & Học bạ số SmartCA</p>
            </div>
            <div class="icon">
                <i class="fa fa-certificate"></i>
            </div>
            <a href="/admin/info/academic_transcripts_snapshots" class="small-box-footer">Học bạ số Snapshot <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
    <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-orange">
            <div class="inner">
                <h3>Kcal/P:L:G</h3>
                <p>Thực đơn & Lưu mẫu 24h</p>
            </div>
            <div class="icon">
                <i class="fa fa-cutlery"></i>
            </div>
            <a href="/admin/info/boarding_menus" class="small-box-footer">Dinh dưỡng Bán trú <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
    <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-teal">
            <div class="inner">
                <h3>BMI/Nha</h3>
                <p>Hồ sơ Sức khỏe & Y tế Học đường</p>
            </div>
            <div class="icon">
                <i class="fa fa-medkit"></i>
            </div>
            <a href="/admin/info/student_health_records" class="small-box-footer">Y tế Học sinh <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
</div>
<div class="row">
    <div class="col-lg-6 col-xs-12">
        <div class="small-box bg-maroon">
            <div class="inner">
                <h3>TT 05/2025</h3>
                <p>Thời Khóa Biểu & Phân Công Tiết Dạy Chuẩn Tiểu Học</p>
            </div>
            <div class="icon">
                <i class="fa fa-calendar"></i>
            </div>
            <a href="/admin/timetable" class="small-box-footer">Xem Thời Khóa Biểu & Chống Trùng Tiết <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
    <div class="col-lg-6 col-xs-12">
        <div class="small-box bg-navy">
            <div class="inner">
                <h3>Sandbox 5+</h3>
                <p>Cổng Tích Hợp Ngoại Vi & Dịch Vụ Số Quốc Gia (VNeID, SmartCA, MOET)</p>
            </div>
            <div class="icon">
                <i class="fa fa-plug"></i>
            </div>
            <a href="/admin/integrations" class="small-box-footer">Giám Sát Cổng Kết Nối Sandbox <i class="fa fa-arrow-circle-right"></i></a>
        </div>
    </div>
</div>
`
	return types.Panel{
		Title:       "EduPrimary-KPI Tổng Quan Quản Trị",
		Description: "Bảng điều khiển trung tâm quản lý hoạt động sư phạm, định mức tiết dạy, học sinh và KPI",
		Content:     template.HTML(html),
	}, nil
}

// TimetablePanel renders the interactive timetable & slot management panel
func TimetablePanel(ctx *context.Context) (types.Panel, error) {
	html := `
<div class="row">
    <div class="col-md-12">
        <div class="callout callout-success" style="margin-bottom: 15px;">
            <h4><i class="fa fa-calendar-check-o"></i> Thời Khóa Biểu & Phân Công Tiết Dạy — Thông tư 05/2025/TT-BGDĐT</h4>
            <p>Hệ thống tự động xếp lịch giảng dạy, kiểm soát định mức chuẩn <strong>805 tiết/năm</strong> (35 tuần thực học, 23 tiết/tuần), áp dụng chế độ giảm trừ <strong>4 tiết/tuần cho Giáo viên Chủ nhiệm</strong>, và kích hoạt thuật toán chống trùng tiết giáo viên giữa các lớp.</p>
        </div>
    </div>
</div>
<div class="row">
    <div class="col-md-3">
        <div class="info-box">
            <span class="info-box-icon bg-aqua"><i class="fa fa-calendar"></i></span>
            <div class="info-box-content">
                <span class="info-box-text">Khối Lớp & Lớp Học</span>
                <span class="info-box-number">Lớp 1A1 (Khối 1)</span>
                <span class="progress-description">Học kỳ I • Năm học 2026-2027</span>
            </div>
        </div>
    </div>
    <div class="col-md-3">
        <div class="info-box">
            <span class="info-box-icon bg-green"><i class="fa fa-check-circle"></i></span>
            <div class="info-box-content">
                <span class="info-box-text">Kiểm Tra Trùng Tiết</span>
                <span class="info-box-number">HỢP LỆ (0 Xung Đột)</span>
                <span class="progress-description">Thuật toán TT 05/2025 xác thực</span>
            </div>
        </div>
    </div>
    <div class="col-md-3">
        <div class="info-box">
            <span class="info-box-icon bg-yellow"><i class="fa fa-clock-o"></i></span>
            <div class="info-box-content">
                <span class="info-box-text">Số Tiết / Tuần</span>
                <span class="info-box-number">30 Tiết / 35 Tuần</span>
                <span class="progress-description">2 buổi/ngày (GDPT 2018)</span>
            </div>
        </div>
    </div>
    <div class="col-md-3">
        <div class="info-box">
            <span class="info-box-icon bg-purple"><i class="fa fa-user-secret"></i></span>
            <div class="info-box-content">
                <span class="info-box-text">GVCN Phụ Trách</span>
                <span class="info-box-number">Thầy Nguyễn Văn An</span>
                <span class="progress-description">Giảm 4 tiết/tuần định mức</span>
            </div>
        </div>
    </div>
</div>
<div class="row">
    <div class="col-md-12">
        <div class="box box-primary">
            <div class="box-header with-border">
                <h3 class="box-title"><i class="fa fa-table"></i> Bảng Phân Bổ Thời Khóa Biểu Tuần Học Chuẩn — Lớp 1A1</h3>
                <div class="box-tools pull-right">
                    <span class="label label-primary">Phòng học 101</span>
                    <span class="label label-success">Sĩ số: 35 HS</span>
                </div>
            </div>
            <div class="box-body table-responsive no-padding">
                <table class="table table-bordered table-striped text-center">
                    <thead>
                        <tr class="bg-gray-active">
                            <th style="width: 100px;">Buổi</th>
                            <th style="width: 80px;">Tiết</th>
                            <th>Thứ Hai</th>
                            <th>Thứ Ba</th>
                            <th>Thứ Tư</th>
                            <th>Thứ Năm</th>
                            <th>Thứ Sáu</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td rowspan="4" class="text-bold" style="vertical-align: middle; background-color: #f9f9f9;">BUỔI SÁNG<br><small class="text-muted">(07:30 - 10:45)</small></td>
                            <td><strong>1</strong></td>
                            <td><span class="badge bg-purple">Chào cờ / SHDC</span><br><small>Toàn trường</small></td>
                            <td><span class="badge bg-blue">Toán học</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-green">Tiếng Việt</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-blue">Toán học</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-green">Tiếng Việt</span><br><small>GVCN Nguyễn Văn An</small></td>
                        </tr>
                        <tr>
                            <td><strong>2</strong></td>
                            <td><span class="badge bg-green">Tiếng Việt</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-green">Tiếng Việt</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-blue">Toán học</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-green">Tiếng Việt</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-blue">Toán học</span><br><small>GVCN Nguyễn Văn An</small></td>
                        </tr>
                        <tr>
                            <td><strong>3</strong></td>
                            <td><span class="badge bg-maroon">Tiếng Anh</span><br><small>GVBM Trần Thị Mai</small></td>
                            <td><span class="badge bg-orange">Đạo đức</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-maroon">Tiếng Anh</span><br><small>GVBM Trần Thị Mai</small></td>
                            <td><span class="badge bg-teal">Tự nhiên & Xã hội</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-maroon">Tiếng Anh</span><br><small>GVBM Trần Thị Mai</small></td>
                        </tr>
                        <tr>
                            <td><strong>4</strong></td>
                            <td><span class="badge bg-red">Giáo dục thể chất</span><br><small>GVBM Lê Văn Cường</small></td>
                            <td><span class="badge bg-navy">Âm nhạc</span><br><small>GVBM Hoàng Thị Oanh</small></td>
                            <td><span class="badge bg-olive">Mĩ thuật</span><br><small>GVBM Phạm Văn Long</small></td>
                            <td><span class="badge bg-red">Giáo dục thể chất</span><br><small>GVBM Lê Văn Cường</small></td>
                            <td><span class="badge bg-purple">HĐ Trải nghiệm</span><br><small>GVCN Nguyễn Văn An</small></td>
                        </tr>
                        <tr class="bg-gray" style="font-weight: bold;">
                            <td colspan="7">NGHỈ TRƯA — ĂN BÁN TRÚ & NGHỈ NGƠI THEO QUY CHUẨN (11:00 - 13:45)</td>
                        </tr>
                        <tr>
                            <td rowspan="3" class="text-bold" style="vertical-align: middle; background-color: #f9f9f9;">BUỔI CHIỀU<br><small class="text-muted">(14:00 - 16:30)</small></td>
                            <td><strong>5</strong></td>
                            <td><span class="badge bg-light-blue">Tự học có hướng dẫn</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-yellow">Kỹ năng sống</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-purple">Bài học STEM (CV 2345)</span><br><small>GVCN Nguyễn Văn An</small></td>
                            <td><span class="badge bg-light-blue">Đọc sách thư viện</span><br><small>GV Phụ trách Thư viện</small></td>
                            <td><span class="badge bg-maroon">Sinh hoạt lớp</span><br><small>GVCN Nguyễn Văn An</small></td>
                        </tr>
                        <tr>
                            <td><strong>6</strong></td>
                            <td><span class="badge bg-teal">CLB Năng khiếu</span><br><small>Giáo viên CLB</small></td>
                            <td><span class="badge bg-olive">CLB Cờ vua</span><br><small>Giáo viên CLB</small></td>
                            <td><span class="badge bg-aqua">CLB Tin học</span><br><small>GVBM Tin học</small></td>
                            <td><span class="badge bg-orange">CLB Mĩ thuật sáng tạo</span><br><small>GVBM Mĩ thuật</small></td>
                            <td><span class="badge bg-green">Tổng kết tuần & Khen thưởng</span><br><small>GVCN & Lớp trưởng</small></td>
                        </tr>
                        <tr>
                            <td><strong>7</strong></td>
                            <td><span class="badge bg-gray">Hoạt động cuối ngày</span><br><small>Trả trẻ bán trú</small></td>
                            <td><span class="badge bg-gray">Hoạt động cuối ngày</span><br><small>Trả trẻ bán trú</small></td>
                            <td><span class="badge bg-gray">Hoạt động cuối ngày</span><br><small>Trả trẻ bán trú</small></td>
                            <td><span class="badge bg-gray">Hoạt động cuối ngày</span><br><small>Trả trẻ bán trú</small></td>
                            <td><span class="badge bg-gray">Hoạt động cuối ngày</span><br><small>Trả trẻ bán trú</small></td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="box-footer">
                <i class="fa fa-info-circle text-info"></i> Hệ thống tự động khóa việc xếp trùng giáo viên (HTTP 409 Conflict) khi phát hiện cùng một giáo viên được phân công 2 lớp trong cùng một khung giờ học theo quy chuẩn Thông tư 05/2025/TT-BGDĐT.
            </div>
        </div>
    </div>
</div>
`
	return types.Panel{
		Title:       "EduPrimary-KPI Phân Bổ Thời Khóa Biểu & Tiết Dạy",
		Description: "Quản lý xếp thời khóa biểu tự động, kiểm soát định mức tiết dạy và thuật toán chống xung đột lịch giáo viên",
		Content:     template.HTML(html),
	}, nil
}

// IntegrationsPanel renders the sandbox stubs overview panel
func IntegrationsPanel(ctx *context.Context) (types.Panel, error) {
	html := `
<div class="row">
    <div class="col-md-12">
        <div class="callout callout-info" style="margin-bottom: 15px;">
            <h4><i class="fa fa-plug"></i> Cổng Tích Hợp Ngoại Vi & Sandbox Dịch Vụ Số Quốc Gia</h4>
            <p>Hệ sinh thái kết nối liên thông đa nền tảng tuân thủ quy chuẩn dữ liệu chính phủ điện tử: VNeID (Bộ Công An), Chữ ký số SmartCA kèm TSA Timestamp, CSDL Ngành Giáo Dục MOET (Bộ GD&ĐT), Cổng thanh toán VietQR EMVCo (NAPAS), và Hóa đơn điện tử (Nghị định 123/2020/NĐ-CP).</p>
        </div>
    </div>
</div>
<div class="row">
    <div class="col-md-4">
        <div class="box box-primary">
            <div class="box-header with-border">
                <h3 class="box-title"><i class="fa fa-id-card"></i> VNeID Mức 2 (Bộ Công An)</h3>
                <span class="pull-right badge bg-green">SANDBOX ACTIVE</span>
            </div>
            <div class="box-body">
                <p>Xác thực sinh trắc học và căn cước công dân gắn chip của phụ huynh / người giám hộ, đảm bảo tuân thủ Điều 20 Nghị định 13/2023/NĐ-CP về bảo vệ dữ liệu trẻ em.</p>
                <ul class="list-group list-group-unbordered">
                    <li class="list-group-item"><b>Cơ quan cấp:</b> <span class="pull-right">C06 - Bộ Công An</span></li>
                    <li class="list-group-item"><b>Cấp độ xác thực:</b> <span class="pull-right">Level 2 (Biometric + Chip)</span></li>
                    <li class="list-group-item"><b>Trạng thái kết nối:</b> <span class="pull-right label label-success">Sẵn sàng (200 OK)</span></li>
                </ul>
            </div>
        </div>
    </div>
    <div class="col-md-4">
        <div class="box box-success">
            <div class="box-header with-border">
                <h3 class="box-title"><i class="fa fa-certificate"></i> SmartCA & National TSA</h3>
                <span class="pull-right badge bg-green">READY</span>
            </div>
            <div class="box-body">
                <p>Ký số học bạ điện tử bất biến bằng định dạng chuẩn CMS/PKCS#7, tích hợp dấu thời gian Time-Stamp Authority (TSA) quốc gia từ VNPT/Viettel đảm bảo giá trị pháp lý.</p>
                <ul class="list-group list-group-unbordered">
                    <li class="list-group-item"><b>Định dạng ký:</b> <span class="pull-right">CMS / PKCS#7 Detached</span></li>
                    <li class="list-group-item"><b>TSA Authority:</b> <span class="pull-right">VNPT-TSA Sandbox Mock</span></li>
                    <li class="list-group-item"><b>Thuật toán băm:</b> <span class="pull-right">SHA-256 with RSA 2048</span></li>
                </ul>
            </div>
        </div>
    </div>
    <div class="col-md-4">
        <div class="box box-warning">
            <div class="box-header with-border">
                <h3 class="box-title"><i class="fa fa-cloud-upload"></i> CSDL Ngành MOET</h3>
                <span class="pull-right badge bg-blue">SYNCING</span>
            </div>
            <div class="box-body">
                <p>Cơ chế Outbox Pattern đồng bộ định kỳ danh sách học sinh, điểm đánh giá định kỳ, kết quả rèn luyện lên Cơ sở dữ liệu ngành Giáo dục & Đào tạo.</p>
                <ul class="list-group list-group-unbordered">
                    <li class="list-group-item"><b>Giao thức:</b> <span class="pull-right">REST API HTTPS / JSON</span></li>
                    <li class="list-group-item"><b>Cơ chế đồng bộ:</b> <span class="pull-right">Transactional Outbox Queue</span></li>
                    <li class="list-group-item"><b>Trạng thái hàng đợi:</b> <span class="pull-right label label-info">Ổn định (0 lỗi chờ)</span></li>
                </ul>
            </div>
        </div>
    </div>
</div>
<div class="row">
    <div class="col-md-6">
        <div class="box box-danger">
            <div class="box-header with-border">
                <h3 class="box-title"><i class="fa fa-qrcode"></i> VietQR Động Chuẩn EMVCo (NAPAS 247)</h3>
                <span class="pull-right badge bg-green">OPERATIONAL</span>
            </div>
            <div class="box-body">
                <p>Sinh mã phản hồi nhanh VietQR động chứa chính xác số tài khoản nhà trường, mã định danh học sinh và số tiền suất ăn bán trú/học phí theo chuẩn EMVCo Quốc Tế.</p>
                <div class="well well-sm">
                    <strong>Cấu trúc EMVCo Payload mẫu:</strong><br>
                    <code>00020101021238540010A000000727012600069704220112888800102026530370454067000005802VN...</code>
                </div>
            </div>
        </div>
    </div>
    <div class="col-md-6">
        <div class="box box-info">
            <div class="box-header with-border">
                <h3 class="box-title"><i class="fa fa-file-text-o"></i> Hóa Đơn Điện Tử (NĐ 123/2020 & TT 78/2021)</h3>
                <span class="pull-right badge bg-green">ISSUING</span>
            </div>
            <div class="box-body">
                <p>Phát hành hóa đơn điện tử tự động cho các khoản thu bán trú và đóng góp tự nguyện, tích hợp ký số nhà trường và mã xác thực của Cơ quan Thuế.</p>
                <div class="well well-sm">
                    <strong>Ký hiệu mẫu hóa đơn:</strong> <code>1C26TAA</code> &nbsp;|&nbsp; 
                    <strong>Trạng thái:</strong> <code>Đã ký số phát hành tự động</code> &nbsp;|&nbsp; 
                    <strong>Đồng bộ cơ quan thuế:</strong> <code>Hợp lệ</code>
                </div>
            </div>
        </div>
    </div>
</div>
<div class="row">
    <div class="col-md-12">
        <div class="box box-default">
            <div class="box-header with-border">
                <h3 class="box-title"><i class="fa fa-history"></i> Nhật Ký Giao Dịch Sandbox Gần Nhất (Recent Sandbox Transactions)</h3>
            </div>
            <div class="box-body table-responsive no-padding">
                <table class="table table-hover">
                    <thead>
                        <tr class="bg-gray-light">
                            <th>Thời Gian</th>
                            <th>Dịch Vụ Tích Hợp</th>
                            <th>Mã Giao Dịch</th>
                            <th>Persona Role</th>
                            <th>Trạng Thái</th>
                            <th>Chi Tiết Payload</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>2026-10-29 12:08:40</td>
                            <td><span class="label label-primary">MOET CSDL Outbox</span></td>
                            <td><code>TX-MOET-CSDL-20261029120840-05e2fc6b</code></td>
                            <td><span class="badge bg-navy">BGH / VAN_THU</span></td>
                            <td><span class="label label-success">200 OK</span></td>
                            <td>Đồng bộ hồ sơ 35 học sinh Lớp 1A1 lên Cơ sở dữ liệu Quốc gia ngành GD&ĐT</td>
                        </tr>
                        <tr>
                            <td>2026-10-29 12:08:41</td>
                            <td><span class="label label-success">SmartCA & TSA</span></td>
                            <td><code>TX-SMARTCA-20261029120841-89e4c1</code></td>
                            <td><span class="badge bg-purple">GVCN</span></td>
                            <td><span class="label label-success">200 OK</span></td>
                            <td>Ký số CMS/PKCS#7 detached học bạ số Nguyễn Hoàng Nam, gắn dấu thời gian TSA VNPT</td>
                        </tr>
                        <tr>
                            <td>2026-10-29 12:08:42</td>
                            <td><span class="label label-danger">VietQR EMVCo</span></td>
                            <td><code>TX-VIETQR-20261029120842-771a0b</code></td>
                            <td><span class="badge bg-orange">KE_TOAN</span></td>
                            <td><span class="label label-success">200 OK</span></td>
                            <td>Sinh QR động thanh toán suất ăn bán trú tháng 10: 700,000 VND (Học sinh HS2026001)</td>
                        </tr>
                        <tr>
                            <td>2026-10-29 12:08:43</td>
                            <td><span class="label label-warning">VNeID Mức 2</span></td>
                            <td><code>TX-VNEID-20261029120843-c06a91</code></td>
                            <td><span class="badge bg-teal">VAN_THU</span></td>
                            <td><span class="label label-success">200 OK</span></td>
                            <td>Xác thực CCCD 001202012345 & Digital Consent người giám hộ theo Điều 20 NĐ 13/2023</td>
                        </tr>
                        <tr>
                            <td>2026-10-29 12:08:44</td>
                            <td><span class="label label-info">Hóa Đơn Điện Tử NĐ 123</span></td>
                            <td><code>TX-EINV-20261029120844-55bf28</code></td>
                            <td><span class="badge bg-orange">KE_TOAN</span></td>
                            <td><span class="label label-success">201 Created</span></td>
                            <td>Phát hành hóa đơn điện tử Ký hiệu 1C26TAA số 00000042 cho Phụ huynh Nguyễn Hoàng Nam</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</div>
`
	return types.Panel{
		Title:       "EduPrimary-KPI Cổng Tích Hợp Ngoại Vi & Sandbox",
		Description: "Giám sát trạng thái 5 cổng kết nối: VNeID Mức 2, SmartCA Chữ ký số TSA, MOET CSDL Outbox, VietQR EMVCo và Hóa đơn điện tử NĐ 123",
		Content:     template.HTML(html),
	}, nil
}
