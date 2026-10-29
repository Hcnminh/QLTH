package tables

import (
	"fmt"
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func init() {
	Generators["classes"] = GetClassTable
}

// GetClassTable returns the GoAdmin Table generator for Classes management.
func GetClassTable(ctx *context.Context) (classTable table.Table) {
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

	classTable = table.NewDefaultTable(ctx, cfg)

	// ==========================================
	// 1. Info Panel (List View & Filter Form)
	// ==========================================
	info := classTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("classes").
		SetTitle("Danh sách Lớp học").
		SetDescription("Quản lý danh sách lớp học theo khối 1-5, phân công GVCN và sĩ số")

	info.AddField("ID", "id", db.UUID).FieldSortable()

	info.AddField("Tên Lớp", "name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Khối Lớp", "grade_level", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(fmt.Sprintf(`<span class="label label-info">Khối %s</span>`, model.Value))
		}).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "1", Text: "Khối 1"},
			{Value: "2", Text: "Khối 2"},
			{Value: "3", Text: "Khối 3"},
			{Value: "4", Text: "Khối 4"},
			{Value: "5", Text: "Khối 5"},
		})

	info.AddField("Phòng học", "room_number", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Sĩ số", "student_count", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(fmt.Sprintf(`<span class="badge bg-aqua">%s HS</span>`, model.Value))
		}).
		FieldSortable()

	info.AddField("Năm học", "academic_year_id", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "2026-2027", Text: "Năm học 2026-2027"},
			{Value: "2025-2026", Text: "Năm học 2025-2026"},
		})

	info.AddField("Mã GVCN", "homeroom_teacher_id", db.UUID)

	info.AddField("Ngày tạo", "created_at", db.Timestamptz).FieldSortable()

	// ==========================================
	// 2. Detail Panel (View Record)
	// ==========================================
	detail := classTable.GetDetail()
	detail.SetTable("classes").
		SetTitle("Chi tiết Lớp học").
		SetDescription("Thông tin chi tiết lớp học, phòng học, sĩ số và giáo viên chủ nhiệm")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Tên Lớp", "name", db.Varchar)
	detail.AddField("Khối Lớp", "grade_level", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return fmt.Sprintf("Khối %s", model.Value)
		})
	detail.AddField("Phòng học", "room_number", db.Varchar)
	detail.AddField("Sĩ số học sinh", "student_count", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return fmt.Sprintf("%s học sinh", model.Value)
		})
	detail.AddField("Năm học", "academic_year_id", db.Varchar)
	detail.AddField("Mã GVCN (homeroom_teacher_id)", "homeroom_teacher_id", db.UUID)
	detail.AddField("Mã Trường (school_id)", "school_id", db.UUID)
	detail.AddField("Thời điểm tạo", "created_at", db.Timestamptz)

	// ==========================================
	// 3. Form Panel (Create / Edit Form)
	// ==========================================
	formList := classTable.GetForm()
	formList.SetTable("classes").
		SetTitle("Thông tin Lớp học").
		SetDescription("Cập nhật thông tin lớp học và phân công GVCN")

	formList.AddField("ID", "id", db.UUID, form.Text).FieldNotAllowEdit()
	formList.AddField("Tên Lớp", "name", db.Varchar, form.Text).FieldMust()
	formList.AddField("Khối Lớp (1-5)", "grade_level", db.Int, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Khối 1", Value: "1"},
			{Text: "Khối 2", Value: "2"},
			{Text: "Khối 3", Value: "3"},
			{Text: "Khối 4", Value: "4"},
			{Text: "Khối 5", Value: "5"},
		}).FieldDefault("1").FieldMust()
	formList.AddField("Phòng học", "room_number", db.Varchar, form.Text).FieldDefault("P.101")
	formList.AddField("Sĩ số", "student_count", db.Int, form.Number).FieldDefault("30")
	formList.AddField("Năm học", "academic_year_id", db.Varchar, form.Text).FieldDefault("2026-2027").FieldMust()
	formList.AddField("Mã GVCN", "homeroom_teacher_id", db.UUID, form.Text)
	formList.AddField("Mã Trường", "school_id", db.UUID, form.Text)

	return classTable
}
