package tables

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"

	"qlth/internal/core/kpi"
)

func init() {
	Generators["kpi_evaluations"] = GetKPIEvaluationTable
}

func formatOfficialRank(model types.FieldModel) interface{} {
	switch model.Value {
	case string(kpi.RankExcellent), "XUẤT SẮC":
		return template.HTML(`<span class="label label-success rank-badge"><i class="fa fa-trophy"></i> Xuất sắc (NĐ 90/2020)</span>`)
	case string(kpi.RankGood), "TỐT":
		return template.HTML(`<span class="label label-primary rank-badge"><i class="fa fa-thumbs-up"></i> Tốt</span>`)
	case string(kpi.RankSatisfactory), "HOÀN THÀNH":
		return template.HTML(`<span class="label label-info rank-badge"><i class="fa fa-check"></i> Hoàn thành</span>`)
	case string(kpi.RankUnsatisfactory), "KHÔNG HOÀN THÀNH":
		return template.HTML(`<span class="label label-danger rank-badge"><i class="fa fa-ban"></i> Không hoàn thành</span>`)
	default:
		return template.HTML(`<span class="label label-default rank-badge">` + template.HTMLEscapeString(model.Value) + `</span>`)
	}
}

func format20PctCap(model types.FieldModel) interface{} {
	if model.Value == "true" || model.Value == "t" || model.Value == "1" {
		return template.HTML(`<span class="label label-warning"><i class="fa fa-exclamation-triangle"></i> Hạ do trần 20% (NĐ 90/2020)</span>`)
	}
	return template.HTML(`<span class="label label-default"><i class="fa fa-check"></i> Chuẩn định mức (≤20%)</span>`)
}

func formatFatalFlaws(model types.FieldModel) interface{} {
	if model.Value == "true" || model.Value == "t" || model.Value == "1" {
		return template.HTML(`<span class="label label-danger"><i class="fa fa-exclamation-circle"></i> Điểm liệt kỷ luật (NĐ 48/2023)</span>`)
	}
	return template.HTML(`<span class="label label-success"><i class="fa fa-check"></i> Không vi phạm</span>`)
}

func formatTutoringViolation(model types.FieldModel) interface{} {
	if model.Value == "true" || model.Value == "t" || model.Value == "1" {
		return template.HTML(`<span class="label label-danger"><i class="fa fa-ban"></i> Vi phạm dạy thêm (TT 29/2024)</span>`)
	}
	return template.HTML(`<span class="label label-success"><i class="fa fa-check"></i> Tuân thủ quy định</span>`)
}

func formatHonorsStatus(model types.FieldModel) interface{} {
	if model.Value == "true" || model.Value == "t" || model.Value == "1" {
		return template.HTML(`<span class="label label-warning"><i class="fa fa-lightbulb-o"></i> Đủ điều kiện xuất sắc (Có SKKN)</span>`)
	}
	return template.HTML(`<span class="label label-default"><i class="fa fa-minus"></i> Chưa có SKKN</span>`)
}

// GetKPIEvaluationTable returns the GoAdmin Table generator for KPI 360 evaluations and ranking management.
// It manages 4-dimension weighted evaluations (40/30/15/15), fatal flaw detection, and Decree 90/2020 20% cap ranking.
func GetKPIEvaluationTable(ctx *context.Context) (kpiTable table.Table) {
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

	kpiTable = table.NewDefaultTable(ctx, cfg)

	setupKPIEvaluationInfo(kpiTable)
	setupKPIEvaluationDetail(kpiTable)
	setupKPIEvaluationForm(kpiTable)

	return kpiTable
}

func setupKPIEvaluationInfo(kpiTable table.Table) {
	info := kpiTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetTable("kpi_evaluations").
		SetTitle("Đánh giá KPI 360 & Xếp loại Thi đua").
		SetDescription("Đánh giá 4 chiều (40-30-15-15), kiểm soát điểm liệt và khống chế trần 20% xuất sắc theo Nghị định 90/2020/NĐ-CP & Nghị định 48/2023/NĐ-CP")

	info.AddField("ID", "id", db.UUID).FieldSortable()

	info.AddField("Staff Name", "staff_name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Cycle", "cycle_name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("System Score (40%)", "automated_system_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-aqua">` + template.HTMLEscapeString(model.Value) + ` (40%)</span>`)
		}).
		FieldSortable()

	info.AddField("Manager Score (30%)", "manager_assessment_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-blue">` + template.HTMLEscapeString(model.Value) + ` (30%)</span>`)
		}).
		FieldSortable()

	info.AddField("Peer Score (15%)", "peer_review_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-purple">` + template.HTMLEscapeString(model.Value) + ` (15%)</span>`)
		}).
		FieldSortable()

	info.AddField("Parent Score (15%)", "parent_satisfaction_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-teal">` + template.HTMLEscapeString(model.Value) + ` (15%)</span>`)
		}).
		FieldSortable()

	info.AddField("Final Score", "final_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-navy" style="font-size: 13px; font-weight: bold;">` + template.HTMLEscapeString(model.Value) + `</span>`)
		}).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorEqual})

	info.AddField("Final Rank", "official_rank", db.Varchar).
		FieldDisplay(formatOfficialRank).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "HOÀN THÀNH XUẤT SẮC NHIỆM VỤ", Text: "Hoàn thành xuất sắc nhiệm vụ"},
			{Value: "HOÀN THÀNH TỐT NHIỆM VỤ", Text: "Hoàn thành tốt nhiệm vụ"},
			{Value: "HOÀN THÀNH NHIỆM VỤ", Text: "Hoàn thành nhiệm vụ"},
			{Value: "KHÔNG HOÀN THÀNH NHIỆM VỤ", Text: "Không hoàn thành nhiệm vụ"},
		})

	info.AddField("Honors Status", "has_innovative_solution", db.Boolean).
		FieldDisplay(formatHonorsStatus).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Có sáng kiến / Đủ điều kiện xuất sắc"},
			{Value: "false", Text: "Không có sáng kiến"},
		})

	info.AddField("Fatal Flaws", "has_disciplinary_action", db.Boolean).
		FieldDisplay(formatFatalFlaws).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Bị kỷ luật (Điểm liệt)"},
			{Value: "false", Text: "Không vi phạm"},
		})

	info.AddField("Tutoring Violation", "violated_tutoring_rules", db.Boolean).
		FieldDisplay(formatTutoringViolation).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Vi phạm cấm dạy thêm (Điểm liệt)"},
			{Value: "false", Text: "Tuân thủ quy định"},
		})

	info.AddField("20% Cap Status", "is_capped_by_20pct", db.Boolean).
		FieldDisplay(format20PctCap).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Hạ do khống chế trần 20% (NĐ 90/2020)"},
			{Value: "false", Text: "Chuẩn định mức (≤ 20%)"},
		})
}

func setupKPIEvaluationDetail(kpiTable table.Table) {
	detail := kpiTable.GetDetail()
	detail.SetTable("kpi_evaluations").
		SetTitle("Chi tiết Đánh giá KPI & Xếp loại Thi đua").
		SetDescription("Hồ sơ đánh giá thi đua 4 chiều, kiểm soát điểm liệt và khống chế trần 20% theo Nghị định 90/2020 & 48/2023")

	detail.AddField("ID", "id", db.UUID)
	detail.AddField("Staff Name", "staff_name", db.Varchar)
	detail.AddField("Cycle", "cycle_name", db.Varchar)
	detail.AddField("Automated System Score (40%)", "automated_system_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-aqua">` + template.HTMLEscapeString(model.Value) + ` / 100 (Trọng số 40%)</span>`)
		})
	detail.AddField("Manager Assessment Score (30%)", "manager_assessment_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-blue">` + template.HTMLEscapeString(model.Value) + ` / 100 (Trọng số 30%)</span>`)
		})
	detail.AddField("Peer Review Score (15%)", "peer_review_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-purple">` + template.HTMLEscapeString(model.Value) + ` / 100 (Trọng số 15%)</span>`)
		})
	detail.AddField("Parent Satisfaction Score (15%)", "parent_satisfaction_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-teal">` + template.HTMLEscapeString(model.Value) + ` / 100 (Trọng số 15%)</span>`)
		})
	detail.AddField("Final Score", "final_score", db.Decimal).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(`<span class="badge bg-navy" style="font-size: 14px; font-weight: bold; padding: 5px 10px;">` + template.HTMLEscapeString(model.Value) + ` / 100</span>`)
		})
	detail.AddField("Final Rank", "official_rank", db.Varchar).FieldDisplay(formatOfficialRank)
	detail.AddField("Honors Status (Innovative Solution)", "has_innovative_solution", db.Boolean).FieldDisplay(formatHonorsStatus)
	detail.AddField("Fatal Flaws (Disciplinary Action)", "has_disciplinary_action", db.Boolean).FieldDisplay(formatFatalFlaws)
	detail.AddField("Fatal Flaws (Tutoring Violation)", "violated_tutoring_rules", db.Boolean).FieldDisplay(formatTutoringViolation)
	detail.AddField("Fatal Flaws (Safety Incident)", "has_severe_safety_incident", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "true" || model.Value == "t" || model.Value == "1" {
				return template.HTML(`<span class="label label-danger"><i class="fa fa-ambulance"></i> Có sự cố an toàn/ngộ độc</span>`)
			}
			return template.HTML(`<span class="label label-success"><i class="fa fa-check"></i> An toàn tuyệt đối</span>`)
		})
	detail.AddField("20% Cap Status", "is_capped_by_20pct", db.Boolean).FieldDisplay(format20PctCap)
	detail.AddField("Reason / Notes", "reason", db.Text)
	detail.AddField("Created At", "created_at", db.Datetime)
	detail.AddField("Updated At", "updated_at", db.Datetime)
}

func setupKPIEvaluationForm(kpiTable table.Table) {
	formList := kpiTable.GetForm()
	formList.SetTable("kpi_evaluations").
		SetTitle("Đánh giá KPI & Xếp loại Thi đua").
		SetDescription("Cập nhật kết quả đánh giá 4 chiều, điểm liệt và áp dụng trần 20% theo Nghị định 90/2020 & Nghị định 48/2023")

	formList.AddField("ID", "id", db.UUID, form.Default).
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisableWhenCreate()

	formList.AddField("Staff Name", "staff_name", db.Varchar, form.Text).
		FieldMust().
		FieldHelpMsg("Họ và tên cán bộ / giáo viên được đánh giá")

	formList.AddField("Cycle", "cycle_name", db.Varchar, form.Text).
		FieldDefault("Năm học 2026-2027").
		FieldMust().
		FieldHelpMsg("Kỳ đánh giá KPI / Năm học")

	formList.AddField("Automated System Score (40%)", "automated_system_score", db.Decimal, form.Number).
		FieldDefault("0").
		FieldHelpMsg("Điểm chấm tự động từ hệ thống quản lý giảng dạy (Trọng số 40%)")

	formList.AddField("Manager Assessment Score (30%)", "manager_assessment_score", db.Decimal, form.Number).
		FieldDefault("0").
		FieldHelpMsg("Điểm đánh giá từ Ban Giám hiệu / Tổ trưởng chuyên môn (Trọng số 30%)")

	formList.AddField("Peer Review Score (15%)", "peer_review_score", db.Decimal, form.Number).
		FieldDefault("0").
		FieldHelpMsg("Điểm đánh giá từ đồng nghiệp trong tổ sư phạm (Trọng số 15%)")

	formList.AddField("Parent Satisfaction Score (15%)", "parent_satisfaction_score", db.Decimal, form.Number).
		FieldDefault("0").
		FieldHelpMsg("Điểm khảo sát độ hài lòng từ cha mẹ học sinh (Trọng số 15%)")

	formList.AddField("Final Score", "final_score", db.Decimal, form.Number).
		FieldDefault("0").
		FieldMust().
		FieldHelpMsg("Tổng điểm tổng hợp theo công thức 4 chiều (0 - 100)")

	formList.AddField("Final Rank", "official_rank", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "Hoàn thành xuất sắc nhiệm vụ (HOÀN THÀNH XUẤT SẮC NHIỆM VỤ)", Value: "HOÀN THÀNH XUẤT SẮC NHIỆM VỤ"},
			{Text: "Hoàn thành tốt nhiệm vụ (HOÀN THÀNH TỐT NHIỆM VỤ)", Value: "HOÀN THÀNH TỐT NHIỆM VỤ"},
			{Text: "Hoàn thành nhiệm vụ (HOÀN THÀNH NHIỆM VỤ)", Value: "HOÀN THÀNH NHIỆM VỤ"},
			{Text: "Không hoàn thành nhiệm vụ (KHÔNG HOÀN THÀNH NHIỆM VỤ)", Value: "KHÔNG HOÀN THÀNH NHIỆM VỤ"},
		}).
		FieldDefault("HOÀN THÀNH TỐT NHIỆM VỤ").
		FieldMust().
		FieldHelpMsg("Mức xếp loại theo quy định tại Nghị định 90/2020/NĐ-CP và Nghị định 48/2023/NĐ-CP")

	formList.AddField("Honors Status (Innovative Solution)", "has_innovative_solution", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Có sáng kiến", Value: "true"},
			{Text: "Không có", Value: "false"},
		}).
		FieldDefault("false").
		FieldHelpMsg("Có sáng kiến / giải pháp đổi mới sáng tạo được công nhận (Điều kiện bắt buộc để đạt Xuất sắc theo Luật TĐKT 2022)")

	formList.AddField("Fatal Flaws (Disciplinary Action)", "has_disciplinary_action", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Bị kỷ luật", Value: "true"},
			{Text: "Không", Value: "false"},
		}).
		FieldDefault("false").
		FieldHelpMsg("Bị kỷ luật trong năm: Điểm liệt tự động xếp Không hoàn thành nhiệm vụ theo NĐ 48/2023")

	formList.AddField("Fatal Flaws (Tutoring Violation)", "violated_tutoring_rules", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Vi phạm dạy thêm", Value: "true"},
			{Text: "Không", Value: "false"},
		}).
		FieldDefault("false").
		FieldHelpMsg("Vi phạm quy định nghiêm cấm dạy thêm tiểu học (TT 29/2024 & TT 19/2026): Điểm liệt")

	formList.AddField("Fatal Flaws (Safety Incident)", "has_severe_safety_incident", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Có sự cố", Value: "true"},
			{Text: "Không", Value: "false"},
		}).
		FieldDefault("false").
		FieldHelpMsg("Để xảy ra tai nạn thương tích hoặc mất an toàn thực phẩm nghiêm trọng: Điểm liệt")

	formList.AddField("20% Cap Applied", "is_capped_by_20pct", db.Boolean, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Bị khống chế trần 20%", Value: "true"},
			{Text: "Không bị khống chế", Value: "false"},
		}).
		FieldDefault("false").
		FieldHelpMsg("Đánh dấu cán bộ đạt chuẩn xuất sắc nhưng bị điều chỉnh về mức Tốt do khống chế tỷ lệ 20% theo Điều 2 Nghị định 90/2020/NĐ-CP")

	formList.AddField("Reason / Notes", "reason", db.Text, form.TextArea).
		FieldHelpMsg("Chi tiết căn cứ đánh giá, sáng kiến kinh nghiệm hoặc lý do áp dụng trần 20%")
}
