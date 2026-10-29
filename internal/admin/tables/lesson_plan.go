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
	Generators["lesson_plans"] = GetLessonPlanTable
}

func formatLessonPlanStatus(model types.FieldModel) interface{} {
	switch model.Value {
	case "SUBMITTED":
		return template.HTML(`<span class="label label-warning"><i class="fa fa-clock-o"></i> Chờ duyệt (Submitted)</span>`)
	case "DEPT_APPROVED":
		return template.HTML(`<span class="label label-info"><i class="fa fa-check-circle"></i> Tổ trưởng đã duyệt</span>`)
	case "APPROVED":
		return template.HTML(`<span class="label label-success"><i class="fa fa-check"></i> BGH đã duyệt (Approved)</span>`)
	case "REJECTED":
		return template.HTML(`<span class="label label-danger"><i class="fa fa-ban"></i> Yêu cầu chỉnh sửa</span>`)
	default:
		return model.Value
	}
}

// GetLessonPlanTable returns the GoAdmin Table generator for Lesson Plans management.
// Compliant with Official Dispatch 2345/BGDĐT-GDTH and STEM Circular 86/TT-BGDĐT.
func GetLessonPlanTable(ctx *context.Context) (planTable table.Table) {
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

	planTable = table.NewDefaultTable(ctx, cfg)

	// ==========================================
	// 1. Info Panel (List View & Filter Form)
	// ==========================================
	info := planTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("lesson_plans").
		SetTitle("Kế hoạch Bài dạy (Giáo án CV 2345 & STEM TT 86)").
		SetDescription("Quản lý kế hoạch bài dạy, tiến độ nộp và quy trình phê duyệt phân cấp sư phạm (Giáo viên -> Tổ trưởng -> BGH)")

	info.AddField("ID", "id", db.UUID).FieldSortable()

	info.AddField("Tên Bài dạy", "title", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Môn học", "subject_id", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "TOAN", Text: "Toán"},
			{Value: "TIENG_VIET", Text: "Tiếng Việt"},
			{Value: "TIENG_ANH", Text: "Tiếng Anh"},
			{Value: "STEM", Text: "STEM"},
			{Value: "TIN_HOC", Text: "Tin học"},
		})

	info.AddField("Khối", "grade_level", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return fmt.Sprintf("Khối %s", model.Value)
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

	info.AddField("Tuần", "week_number", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return fmt.Sprintf("Tuần %s", model.Value)
		}).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorEqual})

	info.AddField("Chuẩn Khung", "curriculum_standard", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			switch model.Value {
			case "CV_2345":
				return template.HTML(`<span class="label label-primary"><i class="fa fa-book"></i> CV 2345</span>`)
			case "STEM_TT_86", "STEM_TT86":
				return template.HTML(`<span class="label label-info"><i class="fa fa-flask"></i> STEM TT 86</span>`)
			default:
				return model.Value
			}
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "CV_2345", Text: "Công văn 2345/BGDĐT-GDTH"},
			{Value: "STEM_TT_86", Text: "STEM (Thông tư 86 / CV 909)"},
		})

	info.AddField("Trạng thái Duyệt", "current_status", db.Varchar).
		FieldDisplay(formatLessonPlanStatus).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "SUBMITTED", Text: "Chờ duyệt (SUBMITTED)"},
			{Value: "DEPT_APPROVED", Text: "Tổ trưởng đã duyệt (DEPT_APPROVED)"},
			{Value: "APPROVED", Text: "BGH đã duyệt (APPROVED)"},
			{Value: "REJECTED", Text: "Yêu cầu chỉnh sửa (REJECTED)"},
		})

	info.AddField("Số tiết", "duration_slots", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return model.Value + " tiết"
		})

	info.AddField("Thời điểm nộp", "created_at", db.Timestamptz).FieldSortable()

	// ==========================================
	// 2. Detail Panel (View Record)
	// ==========================================
	detail := planTable.GetDetail()
	detail.SetTable("lesson_plans").
		SetTitle("Chi tiết Kế hoạch Bài dạy").
		SetDescription("Toàn văn mục tiêu, thiết bị dạy học, các hoạt động và lịch sử duyệt giáo án")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Tên Bài dạy", "title", db.Varchar)
	detail.AddField("Môn học", "subject_id", db.Varchar)
	detail.AddField("Khối Lớp", "grade_level", db.Int)
	detail.AddField("Tuần dạy học", "week_number", db.Int)
	detail.AddField("Thời lượng (tiết)", "duration_slots", db.Int)
	detail.AddField("Chuẩn khung bài dạy", "curriculum_standard", db.Varchar)
	detail.AddField("Trạng thái Phê duyệt", "current_status", db.Varchar).FieldDisplay(formatLessonPlanStatus)
	detail.AddField("Thiết bị dạy học", "equipment", db.Text)
	detail.AddField("Đường dẫn file đính kèm", "attachment_url", db.Text)
	detail.AddField("Giáo viên soạn (teacher_id)", "teacher_id", db.UUID)
	detail.AddField("Lớp học (class_id)", "class_id", db.UUID)
	detail.AddField("Mã Trường (school_id)", "school_id", db.UUID)
	detail.AddField("Thời điểm nộp", "created_at", db.Timestamptz)
	detail.AddField("Cập nhật lần cuối", "updated_at", db.Timestamptz)

	// ==========================================
	// 3. Form Panel (Create / Edit Form)
	// ==========================================
	formList := planTable.GetForm()
	formList.SetTable("lesson_plans").
		SetTitle("Soạn Kế hoạch Bài dạy").
		SetDescription("Tạo mới hoặc chỉnh sửa kế hoạch bài dạy theo chuẩn quy định")

	formList.AddField("ID", "id", db.UUID, form.Text).FieldNotAllowEdit()
	formList.AddField("Tên Bài dạy", "title", db.Varchar, form.Text).FieldMust()
	formList.AddField("Môn học", "subject_id", db.Varchar, form.Text).FieldDefault("TOAN").FieldMust()
	formList.AddField("Khối Lớp (1-5)", "grade_level", db.Int, form.Number).FieldDefault("5").FieldMust()
	formList.AddField("Tuần dạy học (1-35)", "week_number", db.Int, form.Number).FieldDefault("1").FieldMust()
	formList.AddField("Thời lượng số tiết", "duration_slots", db.Int, form.Number).FieldDefault("1")
	formList.AddField("Khung chuẩn bài dạy", "curriculum_standard", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Công văn 2345/BGDĐT-GDTH", Value: "CV_2345"},
			{Text: "STEM (Thông tư 86/2026/TT-BGDĐT)", Value: "STEM_TT_86"},
		}).FieldDefault("CV_2345").FieldMust()
	formList.AddField("Trạng thái", "current_status", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Chờ duyệt (SUBMITTED)", Value: "SUBMITTED"},
			{Text: "Tổ trưởng duyệt (DEPT_APPROVED)", Value: "DEPT_APPROVED"},
			{Text: "BGH duyệt (APPROVED)", Value: "APPROVED"},
			{Text: "Yêu cầu sửa (REJECTED)", Value: "REJECTED"},
		}).FieldDefault("SUBMITTED")
	formList.AddField("Thiết bị và học liệu dạy học", "equipment", db.Text, form.TextArea)
	formList.AddField("Mã Giáo viên", "teacher_id", db.UUID, form.Text)
	formList.AddField("Mã Lớp", "class_id", db.UUID, form.Text)
	formList.AddField("Mã Trường", "school_id", db.UUID, form.Text)

	return planTable
}
