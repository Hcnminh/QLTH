package tables

import (
	"testing"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
)

func TestStaffTableGeneratorRegistered(t *testing.T) {
	gen, exists := Generators["staffs"]
	if !exists {
		t.Fatal("expected 'staffs' generator to be registered in Generators map")
	}
	if gen == nil {
		t.Fatal("expected 'staffs' generator function not to be nil")
	}
}

func TestGetStaffTableStructure(t *testing.T) {
	ctx := context.NewContext(nil)
	tb := GetStaffTable(ctx)
	if tb == nil {
		t.Fatal("expected GetStaffTable not to return nil")
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
	if info.Table != "staffs" {
		t.Errorf("expected table name 'staffs', got '%s'", info.Table)
	}

	requiredColumns := map[string]bool{
		"Teacher Code":    false,
		"Full Name":       false,
		"Homeroom status": false,
		"Quota":           false,
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
	if formList.Table != "staffs" {
		t.Errorf("expected form table name 'staffs', got '%s'", formList.Table)
	}

	formFields := map[string]bool{
		"employee_code":         false,
		"full_name":             false,
		"is_homeroom":           false,
		"standard_annual_hours": false,
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
	if detail.Table != "staffs" {
		t.Errorf("expected detail table name 'staffs', got '%s'", detail.Table)
	}
}
