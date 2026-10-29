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
	Generators["students"] = GetStudentTable
}

func formatStudentGender(model types.FieldModel) interface{} {
	switch model.Value {
	case "MALE":
		return "Nam"
	case "FEMALE":
		return "Nữ"
	default:
		return model.Value
	}
}

func formatStudentStatus(model types.FieldModel) interface{} {
	switch model.Value {
	case "STUDYING":
		return template.HTML(`<span class="label label-primary">Đang học</span>`)
	case "TRANSFERRED":
		return template.HTML(`<span class="label label-info">Chuyển trường</span>`)
	case "DROPPED_OUT":
		return template.HTML(`<span class="label label-warning">Nghỉ học</span>`)
	case "GRADUATED":
		return template.HTML(`<span class="label label-default">Đã tốt nghiệp</span>`)
	default:
		return model.Value
	}
}

func formatInclusiveEducation(model types.FieldModel) interface{} {
	if model.Value == "true" || model.Value == "t" || model.Value == "1" {
		return template.HTML(`<span class="label label-info"><i class="fa fa-heart"></i> Hòa nhập</span>`)
	}
	return template.HTML(`<span class="label label-default">Bình thường</span>`)
}

func formatDigitalConsent(model types.FieldModel) interface{} {
	if model.Value == "true" || model.Value == "t" || model.Value == "1" {
		return template.HTML(`<span class="label label-success"><i class="fa fa-check"></i> Đã đồng thuận</span>`)
	}
	return template.HTML(`<span class="label label-danger"><i class="fa fa-times"></i> Chưa đồng thuận</span>`)
}

func formatConsentVerified(model types.FieldModel) interface{} {
	if model.Value == "true" || model.Value == "t" || model.Value == "1" {
		return template.HTML(`<span class="label label-success"><i class="fa fa-shield"></i> Tuân thủ NĐ 13 (Đã xác thực)</span>`)
	}
	return template.HTML(`<span class="label label-danger"><i class="fa fa-warning"></i> Chưa tuân thủ NĐ 13 (Chưa xác thực)</span>`)
}

func formatConsentMethod(model types.FieldModel) interface{} {
	switch model.Value {
	case "VNEID_SSO":
		return template.HTML(`<span class="badge bg-purple"><i class="fa fa-id-card"></i> VNeID SSO</span>`)
	case "PORTAL_SIGN":
		return template.HTML(`<span class="badge bg-blue"><i class="fa fa-pencil-square-o"></i> Cổng thông tin</span>`)
	case "PAPER_CONSENT":
		return template.HTML(`<span class="badge bg-gray"><i class="fa fa-file-text-o"></i> Biên bản giấy</span>`)
	default:
		if model.Value == "" {
			return template.HTML(`<span class="text-muted"><i class="fa fa-minus"></i> Chưa có</span>`)
		}
		return model.Value
	}
}

// GetStudentTable returns the GoAdmin Table generator for Students & Digital Consent management.
// It manages student profiles and child data privacy compliance under Article 20 of Decree 13/2023/NĐ-CP.
func GetStudentTable(ctx *context.Context) (studentTable table.Table) {
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

	studentTable = table.NewDefaultTable(ctx, cfg)

	setupStudentInfo(studentTable)
	setupStudentDetail(studentTable)
	setupStudentForm(studentTable)

	return studentTable
}

func setupStudentInfo(studentTable table.Table) {
	info := studentTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("students").
		SetTitle("Học sinh & Consent NĐ13").
		SetDescription("Quản lý danh sách học sinh và trạng thái đồng thuận kỹ thuật số bảo vệ dữ liệu trẻ em (Nghị định 13/2023/NĐ-CP)")

	info.AddField("ID", "id", db.UUID).FieldSortable()

	info.AddField("Student Code", "student_code", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Full Name", "full_name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Gender", "gender", db.Varchar).
		FieldDisplay(formatStudentGender).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "MALE", Text: "Nam (Male)"},
			{Value: "FEMALE", Text: "Nữ (Female)"},
		})

	info.AddField("Birth Date", "birth_date", db.Date).
		FieldSortable()

	info.AddField("Status", "status", db.Varchar).
		FieldDisplay(formatStudentStatus).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "STUDYING", Text: "Đang học (Studying)"},
			{Value: "TRANSFERRED", Text: "Chuyển trường (Transferred)"},
			{Value: "DROPPED_OUT", Text: "Nghỉ học (Dropped Out)"},
			{Value: "GRADUATED", Text: "Đã tốt nghiệp (Graduated)"},
		})

	info.AddField("Inclusive Education", "is_inclusive_education", db.Boolean).
		FieldDisplay(formatInclusiveEducation).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Hòa nhập (Inclusive)"},
			{Value: "false", Text: "Bình thường (Standard)"},
		})

	info.AddField("Guardian Name", "guardian_name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Digital Consent", "has_digital_consent", db.Boolean).
		FieldDisplay(formatDigitalConsent).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Đã đồng thuận (Consented)"},
			{Value: "false", Text: "Chưa đồng thuận (Missing)"},
		})

	info.AddField("Consent Verified", "digital_consent_verified", db.Boolean).
		FieldDisplay(formatConsentVerified).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Đã xác thực (Verified / Compliant)"},
			{Value: "false", Text: "Chưa xác thực (Unverified / Non-compliant)"},
		})

	info.AddField("Consent Method", "consent_method", db.Varchar).
		FieldDisplay(formatConsentMethod).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "VNEID_SSO", Text: "VNeID SSO"},
			{Value: "PORTAL_SIGN", Text: "Cổng thông tin (Portal)"},
			{Value: "PAPER_CONSENT", Text: "Biên bản giấy (Paper)"},
		})
}

func setupStudentDetail(studentTable table.Table) {
	detail := studentTable.GetDetail()
	detail.SetTable("students").
		SetTitle("Chi tiết Hồ sơ Học sinh & Consent NĐ 13").
		SetDescription("Hồ sơ học sinh và trạng thái bảo vệ dữ liệu trẻ em theo Nghị định 13/2023/NĐ-CP")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Student Code", "student_code", db.Varchar)
	detail.AddField("Full Name", "full_name", db.Varchar)
	detail.AddField("National ID", "national_id", db.Varchar)
	detail.AddField("Gender", "gender", db.Varchar).FieldDisplay(formatStudentGender)
	detail.AddField("Birth Date", "birth_date", db.Date)
	detail.AddField("Ethnicity", "ethnicity", db.Varchar)
	detail.AddField("Status", "status", db.Varchar).FieldDisplay(formatStudentStatus)
	detail.AddField("Inclusive Education", "is_inclusive_education", db.Boolean).FieldDisplay(formatInclusiveEducation)
	detail.AddField("Guardian Name", "guardian_name", db.Varchar)
	detail.AddField("Guardian National ID", "guardian_national_id", db.Varchar)
	detail.AddField("Guardian Phone", "guardian_phone", db.Varchar)
	detail.AddField("Guardian Relationship", "guardian_relationship", db.Varchar)
	detail.AddField("Digital Consent", "has_digital_consent", db.Boolean).FieldDisplay(formatDigitalConsent)
	detail.AddField("Consent Verified", "digital_consent_verified", db.Boolean).FieldDisplay(formatConsentVerified)
	detail.AddField("Consent Method", "consent_method", db.Varchar)
	detail.AddField("Consent Signed At", "consent_signed_at", db.Datetime)
}

func setupStudentForm(studentTable table.Table) {
	formList := studentTable.GetForm()
	formList.SetTable("students").
		SetTitle("Hồ sơ Học sinh & Đồng thuận NĐ 13").
		SetDescription("Cập nhật thông tin học sinh và cam kết bảo vệ dữ liệu trẻ em theo Điều 20 Nghị định 13/2023/NĐ-CP")

	formList.AddField("ID", "id", db.UUID, form.Default).
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisableWhenCreate()

	formList.AddField("Student Code", "student_code", db.Varchar, form.Text).FieldMust()
	formList.AddField("Full Name", "full_name", db.Varchar, form.Text).FieldMust()
	formList.AddField("National ID", "national_id", db.Varchar, form.Text).
		FieldHelpMsg("Mã định danh cá nhân / Số CCCD của học sinh")
	formList.AddField("Gender", "gender", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Nam (Male)", Value: "MALE"},
			{Text: "Nữ (Female)", Value: "FEMALE"},
		}).FieldDefault("MALE").FieldMust()
	formList.AddField("Birth Date", "birth_date", db.Date, form.Date).FieldMust()
	formList.AddField("Ethnicity", "ethnicity", db.Varchar, form.Text).FieldDefault("Kinh")
	formList.AddField("Status", "status", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Đang học (STUDYING)", Value: "STUDYING"},
			{Text: "Chuyển trường (TRANSFERRED)", Value: "TRANSFERRED"},
			{Text: "Nghỉ học (DROPPED_OUT)", Value: "DROPPED_OUT"},
			{Text: "Đã tốt nghiệp (GRADUATED)", Value: "GRADUATED"},
		}).FieldDefault("STUDYING").FieldMust()
	formList.AddField("Inclusive Education", "is_inclusive_education", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Hòa nhập", Value: "true"},
			{Text: "Bình thường", Value: "false"},
		}).FieldDefault("false").
		FieldHelpMsg("Học sinh khuyết tật / giáo dục hòa nhập theo TT 05/2025/TT-BGDĐT")
	formList.AddField("Guardian Name", "guardian_name", db.Varchar, form.Text).FieldMust()
	formList.AddField("Guardian National ID", "guardian_national_id", db.Varchar, form.Text).
		FieldMust().
		FieldHelpMsg("Số CCCD của cha/mẹ hoặc người giám hộ theo Điều 20 NĐ 13/2023")
	formList.AddField("Guardian Phone", "guardian_phone", db.Varchar, form.Text)
	formList.AddField("Guardian Relationship", "guardian_relationship", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Cha/Mẹ (PARENT)", Value: "PARENT"},
			{Text: "Người giám hộ hợp pháp (GUARDIAN)", Value: "GUARDIAN"},
		}).FieldDefault("PARENT")
	formList.AddField("Digital Consent", "has_digital_consent", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Đã đồng thuận", Value: "true"},
			{Text: "Chưa đồng thuận", Value: "false"},
		}).FieldDefault("false").
		FieldHelpMsg("Đồng thuận của cha mẹ/người giám hộ thu thập và xử lý dữ liệu trẻ em (Điều 20 NĐ 13/2023)")
	formList.AddField("Consent Verified", "digital_consent_verified", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Đã xác thực", Value: "true"},
			{Text: "Chưa xác thực", Value: "false"},
		}).FieldDefault("false").
		FieldHelpMsg("Trạng thái xác thực điện tử qua VNeID hoặc chữ ký số")
	formList.AddField("Consent Method", "consent_method", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "VNeID SSO (VNEID_SSO)", Value: "VNEID_SSO"},
			{Text: "Cổng thông tin điện tử (PORTAL_SIGN)", Value: "PORTAL_SIGN"},
			{Text: "Biên bản giấy (PAPER_CONSENT)", Value: "PAPER_CONSENT"},
		}).FieldDefault("VNEID_SSO")
}
