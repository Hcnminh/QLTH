package teaching

import "testing"

// Trường PTDT bán trú tiểu học: định mức 21 tiết/tuần (735 tiết/năm)
func TestCalculateQuota_EthnicBoardingSchool(t *testing.T) {
	engine := NewQuotaEngine()

	profile := TeacherWorkloadProfile{
		TeacherID:         "GV010",
		SchoolType:        "ETHNIC_BOARDING_PRIMARY",
		WeeksInSchoolYear: 35,
		ActualTaughtHours: 700,
	}

	res := engine.CalculateQuota(profile)

	if res.StandardAnnualHours != 735 { // 21 tiết/tuần x 35 tuần
		t.Fatalf("expected ethnic boarding standard 735, got %d", res.StandardAnnualHours)
	}
	if res.OvertimeHours != 0 { // 700 < 735 nên không có tiết thừa
		t.Fatalf("expected no overtime, got %f", res.OvertimeHours)
	}
}

// weeks_in_school_year <= 0 phải fallback về chuẩn 35 tuần
func TestCalculateQuota_DefaultWeeksFallback(t *testing.T) {
	engine := NewQuotaEngine()

	profile := TeacherWorkloadProfile{
		TeacherID:         "GV011",
		SchoolType:        "PRIMARY",
		WeeksInSchoolYear: 0, // không khai báo
		ActualTaughtHours: 700,
	}

	res := engine.CalculateQuota(profile)

	if res.StandardAnnualHours != 805 { // 23 x 35 mặc định
		t.Fatalf("expected default standard 805, got %d", res.StandardAnnualHours)
	}
	if res.NetAnnualQuota != 805 {
		t.Fatalf("expected net quota 805 without reductions, got %d", res.NetAnnualQuota)
	}
	if res.OvertimeHours != 0 { // 700 < 805
		t.Fatalf("expected no overtime, got %f", res.OvertimeHours)
	}
}

// Quy đổi tiết dạy hè tiếng Việt cho trẻ DTTS trước lớp 1: hệ số 1.5
func TestCalculateQuota_EthnicSummerTeachingBonus(t *testing.T) {
	engine := NewQuotaEngine()

	profile := TeacherWorkloadProfile{
		TeacherID:            "GV012",
		SchoolType:           "PRIMARY",
		WeeksInSchoolYear:    35,
		ActualTaughtHours:    800,
		PreGrade1EthnicHours: 20, // dạy hè 20 tiết, quy đổi 20 x 1.5 = 30
	}

	res := engine.CalculateQuota(profile)

	// Tổng tiết quy đổi = 800 + 30 = 830; định mức 805 → thừa 25
	if res.TotalEffectiveHours != 830.0 {
		t.Fatalf("expected total effective hours 830, got %f", res.TotalEffectiveHours)
	}
	if res.OvertimeHours != 25.0 {
		t.Fatalf("expected overtime 25, got %f", res.OvertimeHours)
	}
}

// Dạy hòa nhập học sinh khuyết tật được quy đổi hệ số 0.2 trên tổng tiết cả năm
func TestCalculateQuota_InclusiveTeachingBonus(t *testing.T) {
	engine := NewQuotaEngine()

	profile := TeacherWorkloadProfile{
		TeacherID:           "GV013",
		SchoolType:          "PRIMARY",
		WeeksInSchoolYear:   35,
		ActualTaughtHours:   800,
		InclusiveTeachingWk: 2, // 2 tiết/tuần x 35 tuần = 70 tiết, quy đổi 70 x 0.2 = 14
	}

	res := engine.CalculateQuota(profile)

	if res.InclusiveBonusHours != 14.0 {
		t.Fatalf("expected inclusive bonus 14.0, got %f", res.InclusiveBonusHours)
	}
	// Tổng quy đổi = 800 + 14 = 814; định mức 805 → thừa 9
	if res.OvertimeHours != 9.0 {
		t.Fatalf("expected overtime 9, got %f", res.OvertimeHours)
	}
}

// Kiêm nhiệm tư vấn tâm lý học sinh được giảm định mức theo số tiết khai báo
func TestCalculateQuota_PsychologyCounselingReduction(t *testing.T) {
	engine := NewQuotaEngine()

	profile := TeacherWorkloadProfile{
		TeacherID:           "GV014",
		SchoolType:          "PRIMARY",
		WeeksInSchoolYear:   35,
		CounselingHoursWeek: 4, // tư vấn tâm lý 4 tiết/tuần
		ActualTaughtHours:   665,
	}

	res := engine.CalculateQuota(profile)

	if res.WeeklyReducedHours != 4 {
		t.Fatalf("expected weekly reduction 4, got %d", res.WeeklyReducedHours)
	}
	if res.NetAnnualQuota != 665 { // 805 - 140
		t.Fatalf("expected net quota 665, got %d", res.NetAnnualQuota)
	}
	if res.OvertimeHours != 0 {
		t.Fatalf("expected no overtime, got %f", res.OvertimeHours)
	}
}

// Tổng các khoản giảm không được vượt định mức chuẩn (net quota không âm)
func TestCalculateQuota_ReductionNeverNegativeQuota(t *testing.T) {
	engine := NewQuotaEngine()

	profile := TeacherWorkloadProfile{
		TeacherID:           "GV015",
		SchoolType:          "PRIMARY",
		WeeksInSchoolYear:   35,
		IsHomeroom:          true, // -4
		IsHeadSubject:       true, // -3
		IsDeputyHeadSubject: true, // -1
		HasYoungChild:       true, // -3
		CounselingHoursWeek: 8,    // -8 (tối đa)
		ActualTaughtHours:   100,
	}

	res := engine.CalculateQuota(profile)

	// Tổng giảm 19 tiết/tuần x 35 = 665; 805 - 665 = 140 (vẫn dương)
	if res.WeeklyReducedHours != 19 {
		t.Fatalf("expected weekly reduction 19, got %d", res.WeeklyReducedHours)
	}
	if res.NetAnnualQuota != 140 {
		t.Fatalf("expected net quota 140, got %d", res.NetAnnualQuota)
	}
	if res.OvertimeHours != 0 { // 100 < 140
		t.Fatalf("expected no overtime, got %f", res.OvertimeHours)
	}
}
