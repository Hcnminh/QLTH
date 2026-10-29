package health

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHealthEngine_CalculateBMI(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name     string
		heightCm float64
		weightKg float64
		wantBMI  float64
	}{
		{
			name:     "Normal primary 1st grader",
			heightCm: 118.5,
			weightKg: 21.0,
			wantBMI:  15.0,
		},
		{
			name:     "5th grader normal",
			heightCm: 140.0,
			weightKg: 32.0,
			wantBMI:  16.3,
		},
		{
			name:     "Invalid height or weight",
			heightCm: 0,
			weightKg: 25.0,
			wantBMI:  0.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := engine.CalculateBMI(tc.heightCm, tc.weightKg)
			if got != tc.wantBMI {
				t.Errorf("CalculateBMI(%f, %f) = %f, want %f", tc.heightCm, tc.weightKg, got, tc.wantBMI)
			}
		})
	}
}

func TestHealthEngine_ClassifyBMICategory(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		bmi     float64
		wantCat string
	}{
		{bmi: 13.5, wantCat: BMICategoryUnderweight},
		{bmi: 15.0, wantCat: BMICategoryNormal},
		{bmi: 19.2, wantCat: BMICategoryOverweight},
		{bmi: 22.5, wantCat: BMICategoryObese},
	}

	for _, tc := range tests {
		cat, advice := engine.ClassifyBMICategory(tc.bmi)
		if cat != tc.wantCat {
			t.Errorf("ClassifyBMICategory(%f) = %s, want %s", tc.bmi, cat, tc.wantCat)
		}
		if advice == "" {
			t.Errorf("Expected advice for BMI %f", tc.bmi)
		}
	}
}

func TestHealthEngine_EvaluateEyesight(t *testing.T) {
	engine := NewEngine()

	// 10/10 both eyes
	status, hasDefect, _ := engine.EvaluateEyesight("10/10", "10/10")
	if hasDefect || status != "Thị lực tốt, đạt chuẩn học đường" {
		t.Errorf("Expected perfect vision, got status=%s, defect=%v", status, hasDefect)
	}

	// 8/10 mild reduction
	status, hasDefect, _ = engine.EvaluateEyesight("10/10", "8/10")
	if !hasDefect || status != "Thị lực giảm nhẹ, nghi ngờ khúc xạ học đường" {
		t.Errorf("Expected mild reduction, got status=%s, defect=%v", status, hasDefect)
	}

	// 5/10 severe myopia
	status, hasDefect, rec := engine.EvaluateEyesight("5/10", "6/10")
	if !hasDefect || status != "Cận thị / Tật khúc xạ" {
		t.Errorf("Expected myopia status, got %s", status)
	}
	if rec == "" {
		t.Errorf("Expected eye recommendation")
	}
}

func TestHealthEngine_BuildHealthProfile(t *testing.T) {
	engine := NewEngine()
	studentID := uuid.New()
	checkupDate := time.Date(2026, 9, 15, 0, 0, 0, 0, time.Local)

	events := []MedicalEvent{
		{
			ID:             "ev-1",
			Date:           "2026-10-10",
			EventType:      "TRẦY_XƯỚC",
			Treatment:      "Sát khuẩn bằng Betadine và băng gạc cá nhân",
			NurseName:      "Nguyễn Y Tá",
			ParentNotified: true,
		},
	}

	prescriptions := []PrescriptionItem{
		{
			ID:                 "rx-1",
			MedicationName:     "Paracetamol 250mg",
			Dosage:             "1 gói khi sốt trên 38.5 độ C",
			AdministrationTime: "Sau ăn trưa",
			Status:             "PENDING",
		},
	}

	profile := engine.BuildHealthProfile(
		studentID,
		checkupDate,
		118.5,
		21.0,
		"10/10",
		"10/10",
		"Răng tốt, không sâu",
		"Không có tiền sử bệnh lý mãn tính",
		events,
		prescriptions,
	)

	if profile.StudentID != studentID {
		t.Errorf("Expected StudentID %s, got %s", studentID, profile.StudentID)
	}
	if profile.BMI != 15.0 {
		t.Errorf("Expected BMI 15.0, got %f", profile.BMI)
	}
	if profile.BMICategory != BMICategoryNormal {
		t.Errorf("Expected normal category, got %s", profile.BMICategory)
	}
	if len(profile.MedicalEvents) != 1 {
		t.Errorf("Expected 1 medical event, got %d", len(profile.MedicalEvents))
	}
	if len(profile.Prescriptions) != 1 {
		t.Errorf("Expected 1 prescription, got %d", len(profile.Prescriptions))
	}
}

func TestHealthEngine_AdministerPrescription(t *testing.T) {
	engine := NewEngine()

	rx := &PrescriptionItem{
		ID:             "rx-test",
		MedicationName: "Oresol 245",
		Status:         "PENDING",
	}

	err := engine.AdministerPrescription(rx, "Trần Y Tá")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if rx.Status != "ADMINISTERED" {
		t.Errorf("Expected ADMINISTERED status, got %s", rx.Status)
	}
	if rx.AdministeredBy != "Trần Y Tá" {
		t.Errorf("Expected nurse Trần Y Tá, got %s", rx.AdministeredBy)
	}
	if rx.AdministeredAt == nil {
		t.Errorf("Expected AdministeredAt timestamp")
	}

	// Repeated administration should fail
	errRepeat := engine.AdministerPrescription(rx, "Trần Y Tá")
	if errRepeat == nil {
		t.Errorf("Expected error on duplicate administration")
	}
}
