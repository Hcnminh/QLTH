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
	Generators["student_subject_evaluations"] = GetStudentSubjectEvaluationTable
}

// GetStudentSubjectEvaluationTable returns the GoAdmin Table generator for Student Subject Evaluations.
// Covers Circular 27/2020 and Circular 89/2026 evaluation rules (HTT/HT/CHT, 1-10 scores, computer testing).
func GetStudentSubjectEvaluationTable(ctx *context.Context) (evalTable table.Table) {
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

	evalTable = table.NewDefaultTable(ctx, cfg)

	// ==========================================
	// 1. Info Panel (List View & Filter Form)
	// ==========================================
	info := evalTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("student_subject_evaluations").
		SetTitle("Sổ Đánh giá Môn học Học sinh").
		SetDescription("Đánh giá thường xuyên và định kỳ môn học theo Thông tư 27/2020 và Thông tư 89/2026/TT-BGDĐT")

	info.AddField("ID", "id", db.UUID).FieldSortable()

	info.AddField("Mã Học sinh", "student_id", db.UUID).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorEqual})

	info.AddField("Môn học", "subject_id", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "TOAN", Text: "Toán"},
			{Value: "TIENG_VIET", Text: "Tiếng Việt"},
			{Value: "TIENG_ANH", Text: "Tiếng Anh"},
			{Value: "TIN_HOC", Text: "Tin học"},
			{Value: "STEM", Text: "STEM"},
		})

	info.AddField("Học kỳ", "term_id", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "2026-2027-HK1", Text: "Học kỳ 1 (2026-2027)"},
			{Value: "2026-2027-HK2", Text: "Học kỳ 2 (2026-2027)"},
		})

	info.AddField("Mức Đạt", "achievement_level", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			switch model.Value {
			case "HTT":
				return template.HTML(`<span class="label label-success"><i class="fa fa-star"></i> HTT (Hoàn thành tốt)</span>`)
			case "HT":
				return template.HTML(`<span class="label label-primary"><i class="fa fa-check"></i> HT (Hoàn thành)</span>`)
			case "CHT":
				return template.HTML(`<span class="label label-danger"><i class="fa fa-times"></i> CHT (Chưa hoàn thành)</span>`)
			default:
				return model.Value
			}
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "HTT", Text: "HTT (Hoàn thành tốt)"},
			{Value: "HT", Text: "HT (Hoàn thành)"},
			{Value: "CHT", Text: "CHT (Chưa hoàn thành)"},
		})

	info.AddField("Điểm KT", "exam_score", db.Numeric).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "" {
				return template.HTML(`<span class="text-muted">-</span>`)
			}
			return template.HTML(fmt.Sprintf(`<span class="badge bg-purple" style="font-size: 13px;">%s đ</span>`, model.Value))
		}).
		FieldSortable()

	info.AddField("Thi Máy tính", "is_computer_based_test", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-info"><i class="fa fa-desktop"></i> Máy tính (TT 89)</span>`)
			}
			return template.HTML(`<span class="label label-default">Truyền thống</span>`)
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Kiểm tra trên máy tính (TT 89/2026)"},
			{Value: "false", Text: "Kiểm tra trên giấy"},
		})

	info.AddField("Đề chung", "is_centralized_exam", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-warning"><i class="fa fa-file-text"></i> Đề chung Sở/Phòng</span>`)
			}
			return template.HTML(`<span class="label label-default">Đề trường</span>`)
		})

	info.AddField("AI Trợ lý", "is_ai_assisted", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-success"><i class="fa fa-magic"></i> AI Copilot</span>`)
			}
			return template.HTML(`<span class="text-muted">-</span>`)
		})

	info.AddField("Nhận xét Giáo viên", "teacher_comment", db.Text)

	// ==========================================
	// 2. Detail Panel (View Record)
	// ==========================================
	detail := evalTable.GetDetail()
	detail.SetTable("student_subject_evaluations").
		SetTitle("Chi tiết Đánh giá Môn học").
		SetDescription("Chi tiết mức đạt, điểm số kiểm tra định kỳ và nhận xét sư phạm của giáo viên")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Mã Học sinh (student_id)", "student_id", db.UUID)
	detail.AddField("Học kỳ (term_id)", "term_id", db.Varchar)
	detail.AddField("Môn học (subject_id)", "subject_id", db.Varchar)
	detail.AddField("Giáo viên đánh giá (teacher_id)", "teacher_id", db.UUID)
	detail.AddField("Mức Đạt", "achievement_level", db.Varchar)
	detail.AddField("Điểm kiểm tra định kỳ (1.0 - 10.0)", "exam_score", db.Numeric)
	detail.AddField("Kiểm tra trên máy tính (TT 89/2026)", "is_computer_based_test", db.Boolean)
	detail.AddField("Đề chung Phòng/Sở GD&ĐT", "is_centralized_exam", db.Boolean)
	detail.AddField("Nhận xét của Giáo viên", "teacher_comment", db.Text)
	detail.AddField("Gợi ý từ AI Pedagogical Copilot", "ai_suggested_comment", db.Text)
	detail.AddField("Sử dụng AI hỗ trợ", "is_ai_assisted", db.Boolean)
	detail.AddField("Thời điểm cập nhật", "updated_at", db.Timestamptz)

	// ==========================================
	// 3. Form Panel (Create / Edit Form)
	// ==========================================
	formList := evalTable.GetForm()
	formList.SetTable("student_subject_evaluations").
		SetTitle("Đánh giá Môn học").
		SetDescription("Nhập mức đạt và điểm kiểm tra định kỳ")

	formList.AddField("ID", "id", db.UUID, form.Text).FieldNotAllowEdit()
	formList.AddField("Mã Học sinh", "student_id", db.UUID, form.Text).FieldMust()
	formList.AddField("Học kỳ", "term_id", db.Varchar, form.Text).FieldDefault("2026-2027-HK1").FieldMust()
	formList.AddField("Môn học", "subject_id", db.Varchar, form.Text).FieldMust()
	formList.AddField("Mã Giáo viên", "teacher_id", db.UUID, form.Text)
	formList.AddField("Mức đạt", "achievement_level", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "HTT (Hoàn thành tốt)", Value: "HTT"},
			{Text: "HT (Hoàn thành)", Value: "HT"},
			{Text: "CHT (Chưa hoàn thành)", Value: "CHT"},
		}).FieldDefault("HTT").FieldMust()
	formList.AddField("Điểm kiểm tra (1.0 - 10.0)", "exam_score", db.Numeric, form.Number).
		FieldHelpMsg("Áp dụng cho các môn có bài kiểm tra định kỳ theo TT 27/TT 89")
	formList.AddField("Nhận xét của Giáo viên", "teacher_comment", db.Text, form.TextArea)
	formList.AddField("Kiểm tra trên máy tính", "is_computer_based_test", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Có", Value: "true"},
			{Text: "Không", Value: "false"},
		}).FieldDefault("false")
	formList.AddField("Đề kiểm tra chung", "is_centralized_exam", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Có", Value: "true"},
			{Text: "Không", Value: "false"},
		}).FieldDefault("false")

	return evalTable
}
