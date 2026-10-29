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
    (9, 0, 1, 5, 'Học sinh & Consent NĐ13', NULL, '', 'fa-graduation-cap', '/info/students'),
    (12, 0, 1, 6, 'Kế hoạch bài dạy CV 2345', NULL, '', 'fa-book', '/info/lesson_plans'),
    (13, 0, 1, 7, 'Danh mục Môn học', 'Học tập & Đánh giá (TT27 / TT89)', '', 'fa-list-ul', '/info/subjects'),
    (14, 0, 1, 8, 'Đánh giá Môn học', NULL, '', 'fa-check-square-o', '/info/student_subject_evaluations'),
    (15, 0, 1, 9, 'Đánh giá Năng lực - Phẩm chất', NULL, '', 'fa-heart', '/info/student_behavior_evaluations'),
    (16, 0, 1, 10, 'Học bạ số & Chữ ký số', NULL, '', 'fa-certificate', '/info/academic_transcripts_snapshots'),
    (17, 0, 1, 11, 'Thực đơn Bán trú', 'Bán trú & Y tế Học đường', '', 'fa-cutlery', '/info/boarding_menus'),
    (18, 0, 1, 12, 'Hồ sơ Sức khỏe Học sinh', NULL, '', 'fa-medkit', '/info/student_health_records'),
    (10, 0, 1, 13, 'Đánh giá KPI 360', 'Thi đua & Hiệu suất', '', 'fa-star', '/info/kpi_evaluations')
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
SELECT 2, m.id FROM goadmin_menu m WHERE m.id IN (1, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18)
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
`
	return types.Panel{
		Title:       "EduPrimary-KPI Tổng Quan Quản Trị",
		Description: "Bảng điều khiển trung tâm quản lý hoạt động sư phạm, định mức tiết dạy, học sinh và KPI",
		Content:     template.HTML(html),
	}, nil
}
