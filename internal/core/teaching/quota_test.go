package teaching

import "testing"

func TestCalculateQuota_PrimaryHomeroomTeacher(t *testing.T) {
	engine := NewQuotaEngine()

	// Teacher who is a homeroom teacher in primary school (reduction = 4 hrs/wk)
	profile := TeacherWorkloadProfile{
		TeacherID:         "GV001",
		TeacherName:       "Lê Thị Mai",
		SchoolType:        "PRIMARY",
		IsHomeroom:        true,
		WeeksInSchoolYear: 35,
		ActualTaughtHours: 700, // Actual teaching done
	}

	res := engine.CalculateQuota(profile)

	// Standard annual = 805
	if res.StandardAnnualHours != 805 {
		t.Fatalf("expected standard annual 805, got %d", res.StandardAnnualHours)
	}

	// Weekly reduction for primary homeroom = 4
	if res.WeeklyReducedHours != 4 {
		t.Fatalf("expected weekly reduction 4, got %d", res.WeeklyReducedHours)
	}

	// Annual reduction = 4 * 35 = 140
	if res.AnnualReducedHours != 140 {
		t.Fatalf("expected annual reduction 140, got %d", res.AnnualReducedHours)
	}

	// Net quota = 805 - 140 = 665
	if res.NetAnnualQuota != 665 {
		t.Fatalf("expected net quota 665, got %d", res.NetAnnualQuota)
	}

	// Overtime = 700 - 665 = 35 hours
	if res.OvertimeHours != 35.0 {
		t.Fatalf("expected overtime 35.0, got %f", res.OvertimeHours)
	}
}

func TestCalculateQuota_HeadSubjectWithYoungChild(t *testing.T) {
	engine := NewQuotaEngine()

	// Head subject (reduced 3) + young child under 12 months (reduced 3) = 6 hrs/wk
	profile := TeacherWorkloadProfile{
		TeacherID:         "GV002",
		TeacherName:       "Trần Phương Lan",
		SchoolType:        "PRIMARY",
		IsHomeroom:        false,
		IsHeadSubject:     true,
		HasYoungChild:     true,
		WeeksInSchoolYear: 35,
		ActualTaughtHours: 650,
	}

	res := engine.CalculateQuota(profile)

	if res.WeeklyReducedHours != 6 {
		t.Fatalf("expected weekly reduction 6, got %d", res.WeeklyReducedHours)
	}

	// Annual reduction = 6 * 35 = 210
	if res.AnnualReducedHours != 210 {
		t.Fatalf("expected annual reduction 210, got %d", res.AnnualReducedHours)
	}

	// Net quota = 805 - 210 = 595
	if res.NetAnnualQuota != 595 {
		t.Fatalf("expected net quota 595, got %d", res.NetAnnualQuota)
	}

	// Overtime = 650 - 595 = 55 hours
	if res.OvertimeHours != 55.0 {
		t.Fatalf("expected overtime 55.0, got %f", res.OvertimeHours)
	}
}
