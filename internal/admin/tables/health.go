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
	Generators["student_health_records"] = GetStudentHealthRecordTable
}

func formatBMICategory(model types.FieldModel) interface{} {
	switch model.Value {
	case "Bình thường":
		return template.HTML(`<span class="label label-success"><i class="fa fa-heart"></i> Bình thường</span>`)
	case "Suy dinh dưỡng":
		return template.HTML(`<span class="label label-warning"><i class="fa fa-warning"></i> Suy dinh dưỡng</span>`)
	case "Thừa cân":
		return template.HTML(`<span class="label label-info"><i class="fa fa-arrow-up"></i> Thừa cân</span>`)
	case "Béo phì":
		return template.HTML(`<span class="label label-danger"><i class="fa fa-exclamation-triangle"></i> Béo phì</span>`)
	default:
		return model.Value
	}
}

// GetStudentHealthRecordTable returns the GoAdmin Table generator for Student Health Records & BMI
func GetStudentHealthRecordTable(ctx *context.Context) (healthTable table.Table) {
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

	healthTable = table.NewDefaultTable(ctx, cfg)

	info := healthTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("student_health_records").
		SetTitle("Hồ sơ Sức khỏe Học sinh").
		SetDescription("Theo dõi chỉ số BMI, thị lực học đường, nha khoa và nhật ký y tế trường học")

	info.AddField("ID", "id", db.UUID).FieldSortable()
	info.AddField("Student ID", "student_id", db.UUID).FieldSortable()

	info.AddField("Ngày khám", "checkup_date", db.Date).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.Date})

	info.AddField("Chiều cao (cm)", "height_cm", db.Numeric).FieldSortable()
	info.AddField("Cân nặng (kg)", "weight_kg", db.Numeric).FieldSortable()

	info.AddField("BMI", "bmi", db.Numeric).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<strong>` + model.Value + `</strong>`)
		})

	info.AddField("Phân loại thể trạng", "bmi_category", db.Varchar).
		FieldDisplay(formatBMICategory).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "Bình thường", Text: "Bình thường"},
			{Value: "Suy dinh dưỡng", Text: "Suy dinh dưỡng"},
			{Value: "Thừa cân", Text: "Thừa cân"},
			{Value: "Béo phì", Text: "Béo phì"},
		})

	info.AddField("Mắt trái", "left_eye_vision", db.Varchar)
	info.AddField("Mắt phải", "right_eye_vision", db.Varchar)
	info.AddField("Nha học đường", "dental_condition", db.Varchar)
	info.AddField("Bệnh mãn tính / Dị ứng", "chronic_conditions", db.Text)

	// ==========================================
	// Detail Panel (View Record)
	// ==========================================
	detail := healthTable.GetDetail()
	detail.SetTable("student_health_records").
		SetTitle("Chi tiết Hồ sơ Sức khỏe Học sinh").
		SetDescription("Chi tiết chỉ số nhân trắc học BMI, thị lực học đường, nha khoa và tiền sử bệnh lý")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Mã Học sinh (student_id)", "student_id", db.UUID)
	detail.AddField("Ngày khám", "checkup_date", db.Date)
	detail.AddField("Chiều cao (cm)", "height_cm", db.Numeric)
	detail.AddField("Cân nặng (kg)", "weight_kg", db.Numeric)
	detail.AddField("Chỉ số BMI", "bmi", db.Numeric)
	detail.AddField("Phân loại thể trạng BMI", "bmi_category", db.Varchar).FieldDisplay(formatBMICategory)
	detail.AddField("Thị lực mắt trái", "left_eye_vision", db.Varchar)
	detail.AddField("Thị lực mắt phải", "right_eye_vision", db.Varchar)
	detail.AddField("Nha học đường", "dental_condition", db.Varchar)
	detail.AddField("Bệnh mãn tính / Dị ứng", "chronic_conditions", db.Text)
	detail.AddField("Thời điểm cập nhật", "updated_at", db.Timestamptz)

	formList := healthTable.GetForm()
	formList.SetTable("student_health_records").SetTitle("Hồ sơ Sức khỏe").SetDescription("Cập nhật thông tin khám sức khỏe học sinh")
	formList.AddField("ID", "id", db.UUID, form.Text).FieldNotAllowEdit()
	formList.AddField("Student ID", "student_id", db.UUID, form.Text).FieldMust()
	formList.AddField("Ngày khám", "checkup_date", db.Date, form.Date).FieldMust()
	formList.AddField("Chiều cao (cm)", "height_cm", db.Numeric, form.Number).FieldMust()
	formList.AddField("Cân nặng (kg)", "weight_kg", db.Numeric, form.Number).FieldMust()
	formList.AddField("Chỉ số BMI", "bmi", db.Numeric, form.Number).FieldMust()
	formList.AddField("Xếp loại BMI", "bmi_category", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "Bình thường", Text: "Bình thường"},
			{Value: "Suy dinh dưỡng", Text: "Suy dinh dưỡng"},
			{Value: "Thừa cân", Text: "Thừa cân"},
			{Value: "Béo phì", Text: "Béo phì"},
		})
	formList.AddField("Thị lực mắt trái", "left_eye_vision", db.Varchar, form.Text)
	formList.AddField("Thị lực mắt phải", "right_eye_vision", db.Varchar, form.Text)
	formList.AddField("Nha học đường", "dental_condition", db.Varchar, form.Text)
	formList.AddField("Bệnh mãn tính", "chronic_conditions", db.Text, form.TextArea)

	return healthTable
}
