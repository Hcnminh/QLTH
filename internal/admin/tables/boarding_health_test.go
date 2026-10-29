package tables

import (
	"testing"

	"github.com/GoAdminGroup/go-admin/context"
)

func TestBoardingMenuTableGeneratorRegistered(t *testing.T) {
	generator, exists := Generators["boarding_menus"]
	if !exists {
		t.Fatalf("Expected 'boarding_menus' table generator to be registered")
	}

	ctx := context.NewContext(nil)
	tbl := generator(ctx)
	if tbl == nil {
		t.Fatalf("Expected non-nil table instance for 'boarding_menus'")
	}

	info := tbl.GetInfo()
	if info.Table != "boarding_menus" {
		t.Errorf("Expected table name 'boarding_menus', got '%s'", info.Table)
	}
}

func TestStudentHealthRecordTableGeneratorRegistered(t *testing.T) {
	generator, exists := Generators["student_health_records"]
	if !exists {
		t.Fatalf("Expected 'student_health_records' table generator to be registered")
	}

	ctx := context.NewContext(nil)
	tbl := generator(ctx)
	if tbl == nil {
		t.Fatalf("Expected non-nil table instance for 'student_health_records'")
	}

	info := tbl.GetInfo()
	if info.Table != "student_health_records" {
		t.Errorf("Expected table name 'student_health_records', got '%s'", info.Table)
	}
}
