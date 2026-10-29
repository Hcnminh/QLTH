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
	Generators["staffs"] = GetStaffTable
}

// GetStaffTable returns the GoAdmin Table generator for Staff & Teaching Quota management.
// It displays teacher details, homeroom status, and workload quota under Circular 05/2025/TT-BGDĐT.
func GetStaffTable(ctx *context.Context) (staffTable table.Table) {
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

	staffTable = table.NewDefaultTable(ctx, cfg)

	// ==========================================
	// 1. Info Panel (List View & Filter Form)
	// ==========================================
	info := staffTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("staffs").
		SetTitle("Giáo viên & Cán bộ").
		SetDescription("Quản lý danh sách giáo viên, công tác chủ nhiệm và định mức tiết dạy (TT 05/2025/TT-BGDĐT)")

	info.AddField("ID", "id", db.UUID).FieldSortable()

	info.AddField("Teacher Code", "employee_code", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Full Name", "full_name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Gender", "gender", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			switch model.Value {
			case "MALE":
				return "Nam"
			case "FEMALE":
				return "Nữ"
			default:
				return model.Value
			}
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "MALE", Text: "Nam (Male)"},
			{Value: "FEMALE", Text: "Nữ (Female)"},
			{Value: "OTHER", Text: "Khác (Other)"},
		})

	info.AddField("Job Position", "job_position", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "PRIMARY_TEACHER", Text: "Giáo viên Tiểu học"},
			{Value: "PRINCIPAL", Text: "Hiệu trưởng"},
			{Value: "VICE_PRINCIPAL", Text: "Phó Hiệu trưởng"},
			{Value: "TPT_DOI", Text: "Tổng phụ trách Đội"},
		})

	info.AddField("Teacher Rank", "teacher_rank", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "RANK_I", Text: "Hạng I"},
			{Value: "RANK_II", Text: "Hạng II"},
			{Value: "RANK_III", Text: "Hạng III"},
		})

	info.AddField("Homeroom status", "is_homeroom", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-success"><i class="fa fa-check"></i> GVCN (Homeroom)</span>`)
			}
			return template.HTML(`<span class="label label-default"><i class="fa fa-minus"></i> Bộ môn</span>`)
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "GVCN (Homeroom)"},
			{Value: "false", Text: "Bộ môn (Subject)"},
		})

	info.AddField("Quota", "standard_annual_hours", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-blue">` + model.Value + ` tiết/năm</span>`)
		}).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorEqual})

	info.AddField("Weekly Teaching Hours", "weekly_teaching_hours", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return model.Value + " tiết/tuần"
		}).
		FieldSortable()

	info.AddField("Weekly Reduced Hours", "weekly_reduced_hours", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return model.Value + " tiết/tuần"
		}).
		FieldSortable()

	info.AddField("Active Status", "is_active", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-success">Đang công tác</span>`)
			}
			return template.HTML(`<span class="label label-danger">Đã nghỉ</span>`)
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Đang công tác"},
			{Value: "false", Text: "Đã nghỉ"},
		})

	// ==========================================
	// 2. Detail Panel (View Record)
	// ==========================================
	detail := staffTable.GetDetail()
	detail.SetTable("staffs").
		SetTitle("Chi tiết Giáo viên & Định mức").
		SetDescription("Hồ sơ viên chức và định mức giảng dạy theo TT 05/2025/TT-BGDĐT")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Teacher Code", "employee_code", db.Varchar)
	detail.AddField("Full Name", "full_name", db.Varchar)
	detail.AddField("National ID", "national_id", db.Varchar)
	detail.AddField("Gender", "gender", db.Varchar)
	detail.AddField("Birth Date", "birth_date", db.Date)
	detail.AddField("Phone Number", "phone_number", db.Varchar)
	detail.AddField("Email", "email", db.Varchar)
	detail.AddField("Job Position", "job_position", db.Varchar)
	detail.AddField("Teacher Rank", "teacher_rank", db.Varchar)
	detail.AddField("Homeroom status", "is_homeroom", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-success">GVCN (Homeroom)</span>`)
			}
			return template.HTML(`<span class="label label-default">Bộ môn (Subject)</span>`)
		})
	detail.AddField("Quota", "standard_annual_hours", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return model.Value + " tiết/năm"
		})
	detail.AddField("Weekly Teaching Hours", "weekly_teaching_hours", db.Int)
	detail.AddField("Weekly Reduced Hours", "weekly_reduced_hours", db.Int)
	detail.AddField("Active Status", "is_active", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-success">Đang công tác</span>`)
			}
			return template.HTML(`<span class="label label-danger">Đã nghỉ</span>`)
		})

	// ==========================================
	// 3. Form Panel (Create / Edit Form)
	// ==========================================
	formList := staffTable.GetForm()
	formList.SetTable("staffs").
		SetTitle("Thông tin Giáo viên & Định mức").
		SetDescription("Cập nhật thông tin cán bộ giáo viên và định mức tiết dạy TT 05/2025")

	formList.AddField("ID", "id", db.UUID, form.Default).
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisableWhenCreate()

	formList.AddField("Teacher Code", "employee_code", db.Varchar, form.Text).FieldMust()
	formList.AddField("Full Name", "full_name", db.Varchar, form.Text).FieldMust()
	formList.AddField("National ID", "national_id", db.Varchar, form.Text).FieldMust()
	formList.AddField("Gender", "gender", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Nam (Male)", Value: "MALE"},
			{Text: "Nữ (Female)", Value: "FEMALE"},
			{Text: "Khác (Other)", Value: "OTHER"},
		}).FieldDefault("MALE")
	formList.AddField("Birth Date", "birth_date", db.Date, form.Date).FieldMust()
	formList.AddField("Phone Number", "phone_number", db.Varchar, form.Text)
	formList.AddField("Email", "email", db.Varchar, form.Email)
	formList.AddField("Job Position", "job_position", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Giáo viên Tiểu học (PRIMARY_TEACHER)", Value: "PRIMARY_TEACHER"},
			{Text: "Hiệu trưởng (PRINCIPAL)", Value: "PRINCIPAL"},
			{Text: "Phó Hiệu trưởng (VICE_PRINCIPAL)", Value: "VICE_PRINCIPAL"},
			{Text: "Tổng phụ trách Đội (TPT_DOI)", Value: "TPT_DOI"},
		}).FieldDefault("PRIMARY_TEACHER").FieldMust()
	formList.AddField("Teacher Rank", "teacher_rank", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Hạng I (RANK_I)", Value: "RANK_I"},
			{Text: "Hạng II (RANK_II)", Value: "RANK_II"},
			{Text: "Hạng III (RANK_III)", Value: "RANK_III"},
		}).FieldDefault("RANK_II")
	formList.AddField("Homeroom status", "is_homeroom", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "GVCN", Value: "true"},
			{Text: "Không", Value: "false"},
		}).FieldDefault("false").
		FieldHelpMsg("Giáo viên chủ nhiệm tiểu học được giảm 4 tiết/tuần theo TT 05/2025/TT-BGDĐT")
	formList.AddField("Quota", "standard_annual_hours", db.Int, form.Number).
		FieldDefault("805").
		FieldHelpMsg("Định mức chuẩn giáo viên tiểu học: 805 tiết/năm (TT 05/2025/TT-BGDĐT)")
	formList.AddField("Weekly Teaching Hours", "weekly_teaching_hours", db.Int, form.Number).
		FieldDefault("23").
		FieldHelpMsg("Định mức chuẩn: 23 tiết/tuần (tiểu học)")
	formList.AddField("Weekly Reduced Hours", "weekly_reduced_hours", db.Int, form.Number).
		FieldDefault("0").
		FieldHelpMsg("Số tiết giảm trừ hàng tuần (GVCN: -4, Tổ trưởng: -3, Nuôi con nhỏ: -3)")
	formList.AddField("Active Status", "is_active", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Đang công tác", Value: "true"},
			{Text: "Nghỉ việc", Value: "false"},
		}).FieldDefault("true")

	return staffTable
}
