package tables

import (
	"testing"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
)

func TestKPIEvaluationTableGeneratorRegistered(t *testing.T) {
	gen, exists := Generators["kpi_evaluations"]
	if !exists {
		t.Fatal("expected 'kpi_evaluations' generator to be registered in Generators map")
	}
	if gen == nil {
		t.Fatal("expected 'kpi_evaluations' generator function not to be nil")
	}
}

func TestGetKPIEvaluationTableStructure(t *testing.T) {
	ctx := context.NewContext(nil)
	tb := GetKPIEvaluationTable(ctx)
	if tb == nil {
		t.Fatal("expected GetKPIEvaluationTable not to return nil")
	}

	// Verify Primary Key
	pk := tb.GetPrimaryKey()
	if pk.Name != "id" {
		t.Errorf("expected primary key name 'id', got '%s'", pk.Name)
	}
	if pk.Type != db.UUID {
		t.Errorf("expected primary key type db.UUID, got '%v'", pk.Type)
	}

	// Verify Info Panel
	info := tb.GetInfo()
	if info.Table != "kpi_evaluations" {
		t.Errorf("expected table name 'kpi_evaluations', got '%s'", info.Table)
	}

	requiredColumns := map[string]bool{
		"Staff Name":          false,
		"Cycle":               false,
		"System Score (40%)":  false,
		"Manager Score (30%)": false,
		"Peer Score (15%)":    false,
		"Parent Score (15%)":  false,
		"Final Score":         false,
		"Final Rank":          false,
		"Honors Status":       false,
		"Fatal Flaws":         false,
		"Tutoring Violation":  false,
		"20% Cap Status":      false,
	}

	for _, field := range info.FieldList {
		if _, ok := requiredColumns[field.Head]; ok {
			requiredColumns[field.Head] = true
		}
	}

	for col, found := range requiredColumns {
		if !found {
			t.Errorf("expected Info panel column '%s' not found in FieldList", col)
		}
	}

	// Verify Form Panel
	formList := tb.GetForm()
	if formList.Table != "kpi_evaluations" {
		t.Errorf("expected form table name 'kpi_evaluations', got '%s'", formList.Table)
	}

	formFields := map[string]bool{
		"staff_name":                 false,
		"cycle_name":                 false,
		"automated_system_score":     false,
		"manager_assessment_score":   false,
		"peer_review_score":          false,
		"parent_satisfaction_score":  false,
		"final_score":                false,
		"official_rank":              false,
		"has_innovative_solution":    false,
		"has_disciplinary_action":    false,
		"violated_tutoring_rules":    false,
		"has_severe_safety_incident": false,
		"is_capped_by_20pct":         false,
		"reason":                     false,
	}

	for _, field := range formList.FieldList {
		if _, ok := formFields[field.Field]; ok {
			formFields[field.Field] = true
		}
	}

	for fieldName, found := range formFields {
		if !found {
			t.Errorf("expected Form field '%s' not found in FieldList", fieldName)
		}
	}

	// Verify Detail Panel
	detail := tb.GetDetail()
	if detail.Table != "kpi_evaluations" {
		t.Errorf("expected detail table name 'kpi_evaluations', got '%s'", detail.Table)
	}

	detailFields := map[string]bool{
		"staff_name":                 false,
		"cycle_name":                 false,
		"automated_system_score":     false,
		"manager_assessment_score":   false,
		"peer_review_score":          false,
		"parent_satisfaction_score":  false,
		"final_score":                false,
		"official_rank":              false,
		"has_innovative_solution":    false,
		"has_disciplinary_action":    false,
		"violated_tutoring_rules":    false,
		"has_severe_safety_incident": false,
		"is_capped_by_20pct":         false,
		"reason":                     false,
	}

	for _, field := range detail.FieldList {
		if _, ok := detailFields[field.Field]; ok {
			detailFields[field.Field] = true
		}
	}

	for fieldName, found := range detailFields {
		if !found {
			t.Errorf("expected Detail field '%s' not found in FieldList", fieldName)
		}
	}
}
