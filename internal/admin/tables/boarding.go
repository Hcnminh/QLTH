package tables

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func init() {
	Generators["boarding_menus"] = GetBoardingMenuTable
}

func formatMealType(model types.FieldModel) interface{} {
	switch model.Value {
	case "LUNCH":
		return template.HTML(`<span class="label label-primary"><i class="fa fa-cutlery"></i> Bữa trưa (Lunch)</span>`)
	case "SNACK":
		return template.HTML(`<span class="label label-warning"><i class="fa fa-coffee"></i> Bữa xế (Snack)</span>`)
	case "BREAKFAST":
		return template.HTML(`<span class="label label-info"><i class="fa fa-sun-o"></i> Bữa sáng (Breakfast)</span>`)
	default:
		return model.Value
	}
}

func formatSampleRetention(model types.FieldModel) interface{} {
	if model.Value == "true" || model.Value == "t" || model.Value == "1" {
		return template.HTML(`<span class="label label-success"><i class="fa fa-shield"></i> Đã lưu mẫu 24h (QĐ 1246/QĐ-BYT)</span>`)
	}
	return template.HTML(`<span class="label label-danger"><i class="fa fa-clock-o"></i> Chưa xác nhận lưu mẫu</span>`)
}

// GetBoardingMenuTable returns the GoAdmin Table generator for Boarding Menus & Nutrition Planning
func GetBoardingMenuTable(ctx *context.Context) (menuTable table.Table) {
	driver := config.GetDatabases().GetDefault().Driver
	if driver == "" {
		driver = config.DriverPostgresql
	}

	cfg := table.DefaultConfigWithDriver(driver)
	cfg.PrimaryKey.Type = db.UUID
	cfg.PrimaryKey.Name = "id"
	cfg.CanAdd = true
	cfg.Editable = true
	cfg.Deletable = true
	cfg.Exportable = true

	menuTable = table.NewDefaultTable(ctx, cfg)

	info := menuTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("boarding_menus").
		SetTitle("Thực đơn & Dinh dưỡng Bán trú").
		SetDescription("Quản lý thực đơn, cân bằng Kcal/P:L:G (Viện Dinh Dưỡng) và lưu mẫu thức ăn 24h (QĐ 1246/QĐ-BYT)")

	info.AddField("ID", "id", db.UUID).FieldSortable()

	info.AddField("Ngày thực đơn", "menu_date", db.Date).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.Date})

	info.AddField("Bữa ăn", "meal_type", db.Varchar).
		FieldDisplay(formatMealType).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "LUNCH", Text: "Bữa trưa (Lunch 600-750 Kcal)"},
			{Value: "SNACK", Text: "Bữa phụ (Snack 150-250 Kcal)"},
			{Value: "BREAKFAST", Text: "Bữa sáng (Breakfast 350-500 Kcal)"},
		})

	info.AddField("Năng lượng (Kcal)", "total_kcal", db.Numeric).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<strong>` + model.Value + ` Kcal</strong>`)
		})

	info.AddField("Protein (%)", "protein_percent", db.Numeric).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return model.Value + "%"
		})

	info.AddField("Lipid (%)", "lipid_percent", db.Numeric).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return model.Value + "%"
		})

	info.AddField("Glucid (%)", "glucid_percent", db.Numeric).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return model.Value + "%"
		})

	info.AddField("Lưu mẫu 24h", "sample_retention_verified", db.Bool).
		FieldDisplay(formatSampleRetention).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Đã lưu mẫu 24h"},
			{Value: "false", Text: "Chưa lưu mẫu"},
		})

	// ==========================================
	// Detail Panel (View Record)
	// ==========================================
	detail := menuTable.GetDetail()
	detail.SetTable("boarding_menus").
		SetTitle("Chi tiết Thực đơn Bán trú").
		SetDescription("Chi tiết năng lượng, tỷ lệ dinh dưỡng P:L:G và thông tin lưu mẫu an toàn thực phẩm 24h")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Ngày thực đơn", "menu_date", db.Date)
	detail.AddField("Bữa ăn", "meal_type", db.Varchar).FieldDisplay(formatMealType)
	detail.AddField("Tổng năng lượng (Kcal)", "total_kcal", db.Numeric)
	detail.AddField("Tỷ lệ Protein (%)", "protein_percent", db.Numeric)
	detail.AddField("Tỷ lệ Lipid (%)", "lipid_percent", db.Numeric)
	detail.AddField("Tỷ lệ Glucid (%)", "glucid_percent", db.Numeric)
	detail.AddField("Xác nhận lưu mẫu 24h", "sample_retention_verified", db.Bool).FieldDisplay(formatSampleRetention)
	detail.AddField("Người lưu mẫu (sample_keeper_id)", "sample_keeper_id", db.UUID)
	detail.AddField("Mã Trường (school_id)", "school_id", db.UUID)
	detail.AddField("Thời điểm tạo", "created_at", db.Timestamptz)

	formList := menuTable.GetForm()
	formList.SetTable("boarding_menus").SetTitle("Thực đơn Bán trú").SetDescription("Lập thực đơn và kiểm soát an toàn thực phẩm")
	formList.AddField("ID", "id", db.UUID, form.Text).FieldNotAllowEdit()
	formList.AddField("Ngày thực đơn", "menu_date", db.Date, form.Date).FieldMust()
	formList.AddField("Bữa ăn", "meal_type", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "LUNCH", Text: "Bữa trưa (LUNCH)"},
			{Value: "SNACK", Text: "Bữa phụ (SNACK)"},
			{Value: "BREAKFAST", Text: "Bữa sáng (BREAKFAST)"},
		}).FieldMust()
	formList.AddField("Tổng Kcal", "total_kcal", db.Numeric, form.Number).FieldMust()
	formList.AddField("Protein (%)", "protein_percent", db.Numeric, form.Number)
	formList.AddField("Lipid (%)", "lipid_percent", db.Numeric, form.Number)
	formList.AddField("Glucid (%)", "glucid_percent", db.Numeric, form.Number)
	formList.AddField("Đã lưu mẫu 24h", "sample_retention_verified", db.Bool, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "true", Text: "Đã lưu mẫu"},
			{Value: "false", Text: "Chưa lưu mẫu"},
		})

	return menuTable
}
