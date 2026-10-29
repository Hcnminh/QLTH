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
	Generators["student_behavior_evaluations"] = GetStudentBehaviorEvaluationTable
}

func formatTDC(model types.FieldModel) interface{} {
	switch model.Value {
	case "T":
		return template.HTML(`<span class="label label-success">T (Tốt)</span>`)
	case "D":
		return template.HTML(`<span class="label label-primary">D (Đạt)</span>`)
	case "C":
		return template.HTML(`<span class="label label-danger">C (Cần cố gắng)</span>`)
	default:
		return model.Value
	}
}

// GetStudentBehaviorEvaluationTable returns the GoAdmin Table generator for Student Behavior Evaluations.
// Tracks 5 core qualities and 10 competencies evaluated as T (Tốt), D (Đạt), C (Cần cố gắng).
func GetStudentBehaviorEvaluationTable(ctx *context.Context) (behTable table.Table) {
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

	behTable = table.NewDefaultTable(ctx, cfg)

	// ==========================================
	// 1. Info Panel (List View & Filter Form)
	// ==========================================
	info := behTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("student_behavior_evaluations").
		SetTitle("Đánh giá Năng lực & Phẩm chất Học sinh").
		SetDescription("Đánh giá 5 phẩm chất chủ yếu và 10 năng lực cốt lõi theo Thông tư 27/2020 và Thông tư 89/2026/TT-BGDĐT")

	info.AddField("ID", "id", db.UUID).FieldSortable()

	info.AddField("Mã Học sinh", "student_id", db.UUID).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorEqual})

	info.AddField("Học kỳ", "term_id", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "2026-2027-HK1", Text: "Học kỳ 1 (2026-2027)"},
			{Value: "2026-2027-HK2", Text: "Học kỳ 2 (2026-2027)"},
		})

	info.AddField("Yêu nước", "patriotism", db.Varchar).FieldDisplay(formatTDC)
	info.AddField("Nhân ái", "compassion", db.Varchar).FieldDisplay(formatTDC)
	info.AddField("Chăm chỉ", "diligence", db.Varchar).FieldDisplay(formatTDC)
	info.AddField("Trung thực", "honesty", db.Varchar).FieldDisplay(formatTDC)
	info.AddField("Trách nhiệm", "responsibility", db.Varchar).FieldDisplay(formatTDC)

	info.AddField("Tự chủ & Tự học", "autonomy_and_self_learning", db.Varchar).FieldDisplay(formatTDC)
	info.AddField("Giao tiếp & Hợp tác", "communication_and_collaboration", db.Varchar).FieldDisplay(formatTDC)
	info.AddField("Giải quyết vấn đề", "problem_solving_and_creativity", db.Varchar).FieldDisplay(formatTDC)

	info.AddField("Nhận xét Phẩm chất", "qualities_comment", db.Text)
	info.AddField("Nhận xét Năng lực", "competencies_comment", db.Text)

	// ==========================================
	// 2. Detail Panel (View Record)
	// ==========================================
	detail := behTable.GetDetail()
	detail.SetTable("student_behavior_evaluations").
		SetTitle("Chi tiết Đánh giá Phẩm chất & Năng lực").
		SetDescription("Kết quả đánh giá 5 phẩm chất chủ yếu và 10 năng lực (3 năng lực chung + 7 năng lực đặc thù)")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Mã Học sinh (student_id)", "student_id", db.UUID)
	detail.AddField("Học kỳ (term_id)", "term_id", db.Varchar)
	detail.AddField("Giáo viên chủ nhiệm (evaluator_id)", "evaluator_id", db.UUID)

	detail.AddField("Yêu nước", "patriotism", db.Varchar)
	detail.AddField("Nhân ái", "compassion", db.Varchar)
	detail.AddField("Chăm chỉ", "diligence", db.Varchar)
	detail.AddField("Trung thực", "honesty", db.Varchar)
	detail.AddField("Trách nhiệm", "responsibility", db.Varchar)
	detail.AddField("Nhận xét Phẩm chất", "qualities_comment", db.Text)

	detail.AddField("Tự chủ và tự học", "autonomy_and_self_learning", db.Varchar)
	detail.AddField("Giao tiếp và hợp tác", "communication_and_collaboration", db.Varchar)
	detail.AddField("Giải quyết vấn đề và sáng tạo", "problem_solving_and_creativity", db.Varchar)

	detail.AddField("Năng lực ngôn ngữ", "language_competence", db.Varchar)
	detail.AddField("Năng lực tính toán", "mathematical_competence", db.Varchar)
	detail.AddField("Năng lực khoa học", "scientific_competence", db.Varchar)
	detail.AddField("Năng lực công nghệ", "technological_competence", db.Varchar)
	detail.AddField("Năng lực tin học", "computer_competence", db.Varchar)
	detail.AddField("Năng lực thẩm mĩ", "aesthetic_competence", db.Varchar)
	detail.AddField("Năng lực thể chất", "physical_competence", db.Varchar)
	detail.AddField("Nhận xét Năng lực", "competencies_comment", db.Text)

	detail.AddField("Thời điểm cập nhật", "updated_at", db.Timestamptz)

	// ==========================================
	// 3. Form Panel (Create / Edit Form)
	// ==========================================
	formList := behTable.GetForm()
	formList.SetTable("student_behavior_evaluations").
		SetTitle("Đánh giá Phẩm chất & Năng lực").
		SetDescription("Đánh giá T/D/C cho học sinh theo quy định")

	formList.AddField("ID", "id", db.UUID, form.Text).FieldNotAllowEdit()
	formList.AddField("Mã Học sinh", "student_id", db.UUID, form.Text).FieldMust()
	formList.AddField("Học kỳ", "term_id", db.Varchar, form.Text).FieldDefault("2026-2027-HK1").FieldMust()
	formList.AddField("Mã GVCN", "evaluator_id", db.UUID, form.Text)

	tdcOptions := types.FieldOptions{
		{Text: "T (Tốt)", Value: "T"},
		{Text: "D (Đạt)", Value: "D"},
		{Text: "C (Cần cố gắng)", Value: "C"},
	}

	formList.AddField("Yêu nước", "patriotism", db.Varchar, form.SelectSingle).FieldOptions(tdcOptions).FieldDefault("T")
	formList.AddField("Nhân ái", "compassion", db.Varchar, form.SelectSingle).FieldOptions(tdcOptions).FieldDefault("T")
	formList.AddField("Chăm chỉ", "diligence", db.Varchar, form.SelectSingle).FieldOptions(tdcOptions).FieldDefault("T")
	formList.AddField("Trung thực", "honesty", db.Varchar, form.SelectSingle).FieldOptions(tdcOptions).FieldDefault("T")
	formList.AddField("Trách nhiệm", "responsibility", db.Varchar, form.SelectSingle).FieldOptions(tdcOptions).FieldDefault("T")
	formList.AddField("Nhận xét Phẩm chất", "qualities_comment", db.Text, form.TextArea)

	formList.AddField("Tự chủ và tự học", "autonomy_and_self_learning", db.Varchar, form.SelectSingle).FieldOptions(tdcOptions).FieldDefault("T")
	formList.AddField("Giao tiếp và hợp tác", "communication_and_collaboration", db.Varchar, form.SelectSingle).FieldOptions(tdcOptions).FieldDefault("T")
	formList.AddField("Giải quyết vấn đề và sáng tạo", "problem_solving_and_creativity", db.Varchar, form.SelectSingle).FieldOptions(tdcOptions).FieldDefault("T")
	formList.AddField("Nhận xét Năng lực", "competencies_comment", db.Text, form.TextArea)

	return behTable
}
