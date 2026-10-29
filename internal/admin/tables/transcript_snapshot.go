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
	Generators["academic_transcripts_snapshots"] = GetAcademicTranscriptSnapshotTable
}

func formatRegulation(model types.FieldModel) interface{} {
	switch model.Value {
	case "TT27_2020":
		return template.HTML(`<span class="label label-primary"><i class="fa fa-history"></i> Thông tư 27/2020</span>`)
	case "TT89_2026":
		return template.HTML(`<span class="label label-success"><i class="fa fa-graduation-cap"></i> Thông tư 89/2026 (Mới)</span>`)
	default:
		return model.Value
	}
}

func formatPromotion(model types.FieldModel) interface{} {
	switch model.Value {
	case "PROMOTED":
		return template.HTML(`<span class="label label-success"><i class="fa fa-arrow-circle-up"></i> Được lên lớp</span>`)
	case "RETAINED":
		return template.HTML(`<span class="label label-warning"><i class="fa fa-refresh"></i> Rèn luyện lại / Chưa lên lớp</span>`)
	default:
		return model.Value
	}
}

// GetAcademicTranscriptSnapshotTable returns the GoAdmin Table generator for Academic Transcript Snapshots.
// Encapsulates digital academic records with SmartCA digital signature and national TSA timestamp.
func GetAcademicTranscriptSnapshotTable(ctx *context.Context) (transTable table.Table) {
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

	transTable = table.NewDefaultTable(ctx, cfg)

	// ==========================================
	// 1. Info Panel (List View & Filter Form)
	// ==========================================
	info := transTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("academic_transcripts_snapshots").
		SetTitle("Snapshot Học bạ số Bất biến & Ký số").
		SetDescription("Quản lý học bạ điện tử đóng gói dạng snapshot bất biến, ký số SmartCA và đóng dấu thời gian TSA")

	info.AddField("ID", "id", db.UUID).FieldSortable()

	info.AddField("Mã Học sinh", "student_id", db.UUID).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorEqual})

	info.AddField("Năm học", "academic_year_id", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "2026-2027", Text: "Năm học 2026-2027"},
			{Value: "2025-2026", Text: "Năm học 2025-2026"},
		})

	info.AddField("Khối Lớp", "grade_level", db.Int).
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

	info.AddField("Thông tư Áp dụng", "regulation_applied", db.Varchar).
		FieldDisplay(formatRegulation).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "TT27_2020", Text: "Thông tư 27/2020 (trước 01/12/2026)"},
			{Value: "TT89_2026", Text: "Thông tư 89/2026 (từ 01/12/2026)"},
		})

	info.AddField("Khen thưởng Năm học", "annual_reward", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "" {
				return template.HTML(`<span class="text-muted">-</span>`)
			}
			return template.HTML(`<span class="label label-warning"><i class="fa fa-trophy"></i> ` + model.Value + `</span>`)
		}).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Kết quả Lên lớp", "promotion_status", db.Varchar).
		FieldDisplay(formatPromotion).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "PROMOTED", Text: "Được lên lớp"},
			{Value: "RETAINED", Text: "Chưa lên lớp / Rèn luyện lại"},
		})

	info.AddField("Khóa Bất biến", "is_locked", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-success"><i class="fa fa-lock"></i> Đã khóa & Ký số</span>`)
			}
			return template.HTML(`<span class="label label-default"><i class="fa fa-unlock"></i> Bản thảo</span>`)
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Đã khóa snapshot"},
			{Value: "false", Text: "Bản thảo"},
		})

	info.AddField("Ký GVCN", "teacher_signed_at", db.Timestamptz)
	info.AddField("Ký BGH", "principal_signed_at", db.Timestamptz)

	// ==========================================
	// 2. Detail Panel (View Record)
	// ==========================================
	detail := transTable.GetDetail()
	detail.SetTable("academic_transcripts_snapshots").
		SetTitle("Chi tiết Học bạ số Snapshot").
		SetDescription("Dữ liệu học bạ số đóng băng, thông tin chữ ký số và xác thực timestamp")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Mã Học sinh (student_id)", "student_id", db.UUID)
	detail.AddField("Năm học", "academic_year_id", db.Varchar)
	detail.AddField("Khối Lớp", "grade_level", db.Int)
	detail.AddField("Thông tư Áp dụng", "regulation_applied", db.Varchar).FieldDisplay(formatRegulation)
	detail.AddField("Khen thưởng Năm học", "annual_reward", db.Varchar)
	detail.AddField("Kết quả Lên lớp", "promotion_status", db.Varchar).FieldDisplay(formatPromotion)
	detail.AddField("Trạng thái Khóa Bất biến", "is_locked", db.Boolean)
	detail.AddField("Thời điểm GVCN Ký số", "teacher_signed_at", db.Timestamptz)
	detail.AddField("Thời điểm Hiệu trưởng Ký số", "principal_signed_at", db.Timestamptz)
	detail.AddField("Dấu thời gian Quốc gia (TSA Token)", "timestamp_token", db.Text)
	detail.AddField("Chữ ký GVCN (SmartCA CMS)", "teacher_signature", db.Text)
	detail.AddField("Chữ ký Hiệu trưởng (SmartCA CMS)", "principal_signature", db.Text)

	// ==========================================
	// 3. Form Panel (Create / Edit Form)
	// ==========================================
	formList := transTable.GetForm()
	formList.SetTable("academic_transcripts_snapshots").
		SetTitle("Snapshot Học bạ").
		SetDescription("Quản lý trạng thái học bạ số")

	formList.AddField("ID", "id", db.UUID, form.Text).FieldNotAllowEdit()
	formList.AddField("Mã Học sinh", "student_id", db.UUID, form.Text).FieldMust()
	formList.AddField("Năm học", "academic_year_id", db.Varchar, form.Text).FieldDefault("2026-2027").FieldMust()
	formList.AddField("Khối Lớp (1-5)", "grade_level", db.Int, form.Number).FieldDefault("5").FieldMust()
	formList.AddField("Thông tư Quy định", "regulation_applied", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Thông tư 27/2020/TT-BGDĐT", Value: "TT27_2020"},
			{Text: "Thông tư 89/2026/TT-BGDĐT", Value: "TT89_2026"},
		}).FieldDefault("TT89_2026").FieldMust()
	formList.AddField("Danh hiệu Khen thưởng", "annual_reward", db.Varchar, form.Text).FieldDefault("Học sinh Xuất sắc")
	formList.AddField("Lên lớp", "promotion_status", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Được lên lớp (PROMOTED)", Value: "PROMOTED"},
			{Text: "Chưa lên lớp (RETAINED)", Value: "RETAINED"},
		}).FieldDefault("PROMOTED").FieldMust()
	formList.AddField("Khóa Bất biến", "is_locked", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Đã khóa", Value: "true"},
			{Text: "Mở", Value: "false"},
		}).FieldDefault("true")

	return transTable
}
