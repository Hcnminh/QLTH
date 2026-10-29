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

	adminEngine = eng
	return eng, nil
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
`
	return types.Panel{
		Title:       "EduPrimary-KPI Tổng Quan Quản Trị",
		Description: "Bảng điều khiển trung tâm quản lý hoạt động sư phạm, định mức tiết dạy, học sinh và KPI",
		Content:     template.HTML(html),
	}, nil
}
