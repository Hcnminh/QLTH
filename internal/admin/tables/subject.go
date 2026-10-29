package tables

import (
	"html/template"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func init() {
	Generators["subjects"] = GetSubjectTable
}

// GetSubjectTable returns the GoAdmin Table generator for Subjects management.
func GetSubjectTable(ctx *context.Context) (subjectTable table.Table) {
	driver := config.GetDatabases().GetDefault().Driver
	if driver == "" {
		driver = config.DriverPostgresql
	}

	cfg := table.DefaultConfigWithDriver(driver)
	cfg.PrimaryKey.Type = db.Varchar
	cfg.PrimaryKey.Name = "id"
	cfg.CanAdd = true
	cfg.Editable = true
	cfg.Deletable = true
	cfg.Exportable = true

	subjectTable = table.NewDefaultTable(ctx, cfg)

	// ==========================================
	// 1. Info Panel (List View & Filter Form)
	// ==========================================
	info := subjectTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("subjects").
		SetTitle("Danh mục Môn học & Hoạt động Giáo dục").
		SetDescription("Quản lý danh mục môn học, hình thức đánh giá (cho điểm kiểm tra / nhận xét) theo CT GDPT 2018")

	info.AddField("Mã Môn", "id", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Tên Môn học", "name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Hình thức Đánh giá", "is_scored", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-success"><i class="fa fa-calculator"></i> Có điểm kiểm tra (1-10)</span>`)
			}
			return template.HTML(`<span class="label label-info"><i class="fa fa-comment-o"></i> Đánh giá bằng nhận xét</span>`)
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Có điểm kiểm tra định kỳ"},
			{Value: "false", Text: "Chỉ đánh giá bằng nhận xét"},
		})

	info.AddField("Khối lớp Áp dụng", "applicable_grades", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			val := strings.Trim(model.Value, "{}")
			if val == "" {
				return "Tất cả"
			}
			return "Khối " + strings.ReplaceAll(val, ",", ", ")
		})

	// ==========================================
	// 2. Detail Panel (View Record)
	// ==========================================
	detail := subjectTable.GetDetail()
	detail.SetTable("subjects").
		SetTitle("Chi tiết Môn học").
		SetDescription("Thông tin môn học, hình thức đánh giá và phạm vi áp dụng các khối")

	detail.AddField("Mã Môn học", "id", db.Varchar)
	detail.AddField("Tên Môn học", "name", db.Varchar)
	detail.AddField("Hình thức Đánh giá", "is_scored", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return "Đánh giá bằng điểm số kết hợp nhận xét (Thang 10)"
			}
			return "Đánh giá bằng nhận xét (HTT / HT / CHT)"
		})
	detail.AddField("Khối lớp Áp dụng", "applicable_grades", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Trim(model.Value, "{}")
		})

	// ==========================================
	// 3. Form Panel (Create / Edit Form)
	// ==========================================
	formList := subjectTable.GetForm()
	formList.SetTable("subjects").
		SetTitle("Thông tin Môn học").
		SetDescription("Thêm mới hoặc cập nhật môn học trong chương trình giáo dục")

	formList.AddField("Mã Môn học", "id", db.Varchar, form.Text).FieldMust()
	formList.AddField("Tên Môn học", "name", db.Varchar, form.Text).FieldMust()
	formList.AddField("Có điểm kiểm tra", "is_scored", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Có", Value: "true"},
			{Text: "Không", Value: "false"},
		}).FieldDefault("true")

	return subjectTable
}
