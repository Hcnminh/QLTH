package tables

import (
	"testing"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
)

func TestAllSubsystemGeneratorsRegistered(t *testing.T) {
	requiredGenerators := []string{
		"classes",
		"subjects",
		"student_subject_evaluations",
		"student_behavior_evaluations",
		"boarding_menus",
		"student_health_records",
		"lesson_plans",
		"academic_transcripts_snapshots",
		"staffs",
		"students",
		"kpi_evaluations",
	}

	for _, name := range requiredGenerators {
		gen, exists := Generators[name]
		if !exists {
			t.Errorf("Generator for '%s' is not registered in Generators map", name)
			continue
		}
		if gen == nil {
			t.Errorf("Generator function for '%s' is nil", name)
		}
	}
}

func TestClassTableStructure(t *testing.T) {
	ctx := context.NewContext(nil)
	tb := GetClassTable(ctx)
	if tb == nil {
		t.Fatal("expected GetClassTable not to return nil")
	}

	pk := tb.GetPrimaryKey()
	if pk.Name != "id" || pk.Type != db.UUID {
		t.Errorf("expected UUID PK 'id', got name=%s, type=%v", pk.Name, pk.Type)
	}

	info := tb.GetInfo()
	if info.Table != "classes" {
		t.Errorf("expected table 'classes', got '%s'", info.Table)
	}

	columns := map[string]bool{"Tên Lớp": false, "Khối Lớp": false, "Phòng học": false, "Sĩ số": false}
	for _, f := range info.FieldList {
		if _, ok := columns[f.Head]; ok {
			columns[f.Head] = true
		}
	}
	for col, found := range columns {
		if !found {
			t.Errorf("expected Info column '%s' in classes table", col)
		}
	}

	detail := tb.GetDetail()
	if detail.Table != "classes" {
		t.Errorf("expected detail table 'classes', got '%s'", detail.Table)
	}
}

func TestSubjectTableStructure(t *testing.T) {
	ctx := context.NewContext(nil)
	tb := GetSubjectTable(ctx)
	if tb == nil {
		t.Fatal("expected GetSubjectTable not to return nil")
	}

	pk := tb.GetPrimaryKey()
	if pk.Name != "id" || pk.Type != db.Varchar {
		t.Errorf("expected Varchar PK 'id', got name=%s, type=%v", pk.Name, pk.Type)
	}

	info := tb.GetInfo()
	if info.Table != "subjects" {
		t.Errorf("expected table 'subjects', got '%s'", info.Table)
	}

	columns := map[string]bool{"Mã Môn": false, "Tên Môn học": false, "Hình thức Đánh giá": false}
	for _, f := range info.FieldList {
		if _, ok := columns[f.Head]; ok {
			columns[f.Head] = true
		}
	}
	for col, found := range columns {
		if !found {
			t.Errorf("expected Info column '%s' in subjects table", col)
		}
	}
}

func TestStudentSubjectEvaluationTableStructure(t *testing.T) {
	ctx := context.NewContext(nil)
	tb := GetStudentSubjectEvaluationTable(ctx)
	if tb == nil {
		t.Fatal("expected GetStudentSubjectEvaluationTable not to return nil")
	}

	pk := tb.GetPrimaryKey()
	if pk.Name != "id" || pk.Type != db.UUID {
		t.Errorf("expected UUID PK 'id', got name=%s, type=%v", pk.Name, pk.Type)
	}

	info := tb.GetInfo()
	if info.Table != "student_subject_evaluations" {
		t.Errorf("expected table 'student_subject_evaluations', got '%s'", info.Table)
	}

	columns := map[string]bool{"Môn học": false, "Học kỳ": false, "Mức Đạt": false, "Điểm KT": false, "Thi Máy tính": false}
	for _, f := range info.FieldList {
		if _, ok := columns[f.Head]; ok {
			columns[f.Head] = true
		}
	}
	for col, found := range columns {
		if !found {
			t.Errorf("expected Info column '%s' in subject evaluations table", col)
		}
	}

	detail := tb.GetDetail()
	if detail.Table != "student_subject_evaluations" {
		t.Errorf("expected detail table 'student_subject_evaluations', got '%s'", detail.Table)
	}
}

func TestStudentBehaviorEvaluationTableStructure(t *testing.T) {
	ctx := context.NewContext(nil)
	tb := GetStudentBehaviorEvaluationTable(ctx)
	if tb == nil {
		t.Fatal("expected GetStudentBehaviorEvaluationTable not to return nil")
	}

	pk := tb.GetPrimaryKey()
	if pk.Name != "id" || pk.Type != db.UUID {
		t.Errorf("expected UUID PK 'id', got name=%s, type=%v", pk.Name, pk.Type)
	}

	info := tb.GetInfo()
	if info.Table != "student_behavior_evaluations" {
		t.Errorf("expected table 'student_behavior_evaluations', got '%s'", info.Table)
	}

	columns := map[string]bool{"Yêu nước": false, "Nhân ái": false, "Chăm chỉ": false, "Trung thực": false, "Trách nhiệm": false}
	for _, f := range info.FieldList {
		if _, ok := columns[f.Head]; ok {
			columns[f.Head] = true
		}
	}
	for col, found := range columns {
		if !found {
			t.Errorf("expected Info column '%s' in behavior evaluations table", col)
		}
	}

	detail := tb.GetDetail()
	if detail.Table != "student_behavior_evaluations" {
		t.Errorf("expected detail table 'student_behavior_evaluations', got '%s'", detail.Table)
	}
}

func TestLessonPlanTableStructure(t *testing.T) {
	ctx := context.NewContext(nil)
	tb := GetLessonPlanTable(ctx)
	if tb == nil {
		t.Fatal("expected GetLessonPlanTable not to return nil")
	}

	pk := tb.GetPrimaryKey()
	if pk.Name != "id" || pk.Type != db.UUID {
		t.Errorf("expected UUID PK 'id', got name=%s, type=%v", pk.Name, pk.Type)
	}

	info := tb.GetInfo()
	if info.Table != "lesson_plans" {
		t.Errorf("expected table 'lesson_plans', got '%s'", info.Table)
	}

	columns := map[string]bool{"Tên Bài dạy": false, "Môn học": false, "Khối": false, "Tuần": false, "Chuẩn Khung": false, "Trạng thái Duyệt": false}
	for _, f := range info.FieldList {
		if _, ok := columns[f.Head]; ok {
			columns[f.Head] = true
		}
	}
	for col, found := range columns {
		if !found {
			t.Errorf("expected Info column '%s' in lesson plans table", col)
		}
	}

	detail := tb.GetDetail()
	if detail.Table != "lesson_plans" {
		t.Errorf("expected detail table 'lesson_plans', got '%s'", detail.Table)
	}
}

func TestAcademicTranscriptSnapshotTableStructure(t *testing.T) {
	ctx := context.NewContext(nil)
	tb := GetAcademicTranscriptSnapshotTable(ctx)
	if tb == nil {
		t.Fatal("expected GetAcademicTranscriptSnapshotTable not to return nil")
	}

	pk := tb.GetPrimaryKey()
	if pk.Name != "id" || pk.Type != db.UUID {
		t.Errorf("expected UUID PK 'id', got name=%s, type=%v", pk.Name, pk.Type)
	}

	info := tb.GetInfo()
	if info.Table != "academic_transcripts_snapshots" {
		t.Errorf("expected table 'academic_transcripts_snapshots', got '%s'", info.Table)
	}

	columns := map[string]bool{"Năm học": false, "Khối Lớp": false, "Thông tư Áp dụng": false, "Khen thưởng Năm học": false, "Kết quả Lên lớp": false, "Khóa Bất biến": false}
	for _, f := range info.FieldList {
		if _, ok := columns[f.Head]; ok {
			columns[f.Head] = true
		}
	}
	for col, found := range columns {
		if !found {
			t.Errorf("expected Info column '%s' in transcript snapshots table", col)
		}
	}

	detail := tb.GetDetail()
	if detail.Table != "academic_transcripts_snapshots" {
		t.Errorf("expected detail table 'academic_transcripts_snapshots', got '%s'", detail.Table)
	}
}
