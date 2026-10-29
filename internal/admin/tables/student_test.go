package tables

import (
	"testing"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
)

func TestStudentTableGeneratorRegistered(t *testing.T) {
	gen, exists := Generators["students"]
	if !exists {
		t.Fatal("expected 'students' generator to be registered in Generators map")
	}
	if gen == nil {
		t.Fatal("expected 'students' generator function not to be nil")
	}
}

func TestGetStudentTableStructure(t *testing.T) {
	ctx := context.NewContext(nil)
	tb := GetStudentTable(ctx)
	if tb == nil {
		t.Fatal("expected GetStudentTable not to return nil")
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
	if info.Table != "students" {
		t.Errorf("expected table name 'students', got '%s'", info.Table)
	}

	requiredColumns := map[string]bool{
		"Student Code":     false,
		"Full Name":        false,
		"Gender":           false,
		"Status":           false,
		"Guardian Name":    false,
		"Digital Consent":  false,
		"Consent Verified": false,
		"Consent Method":   false,
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
	if formList.Table != "students" {
		t.Errorf("expected form table name 'students', got '%s'", formList.Table)
	}

	formFields := map[string]bool{
		"student_code":             false,
		"full_name":                false,
		"guardian_name":            false,
		"has_digital_consent":      false,
		"digital_consent_verified": false,
		"consent_method":           false,
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
	if detail.Table != "students" {
		t.Errorf("expected detail table name 'students', got '%s'", detail.Table)
	}

	detailFields := map[string]bool{
		"student_code":             false,
		"full_name":                false,
		"has_digital_consent":      false,
		"digital_consent_verified": false,
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
