package kpi

import (
	"testing"

	"github.com/google/uuid"
)

func TestCalculateStaffScore(t *testing.T) {
	engine := NewEngine()

	sub := StaffKPISubmission{
		AutomatedSystemScore:    90.0, // 40% = 36.0
		ManagerAssessmentScore:  80.0, // 30% = 24.0
		PeerReviewScore:         90.0, // 15% = 13.5
		ParentSatisfactionScore: 90.0, // 15% = 13.5
	}
	// Total = 36 + 24 + 13.5 + 13.5 = 87.0
	score := engine.CalculateStaffScore(sub)
	if score != 87.0 {
		t.Fatalf("expected score 87.0, got %f", score)
	}
}

func TestEvaluateIndividual_DisciplinaryFatalFlaw(t *testing.T) {
	engine := NewEngine()

	// High score (95) but penalized for tutoring ban violation (Circular 29/2024 & TT 19/2026)
	sub := StaffKPISubmission{
		StaffID:                 uuid.New(),
		StaffName:               "Nguyễn Văn A",
		AutomatedSystemScore:    95.0,
		ManagerAssessmentScore:  95.0,
		PeerReviewScore:         95.0,
		ParentSatisfactionScore: 95.0,
		ViolatedTutoringRules:   true,
	}

	result := engine.EvaluateIndividual(sub)
	if result.OfficialRank != RankUnsatisfactory {
		t.Fatalf("expected rank %s, got %s", RankUnsatisfactory, result.OfficialRank)
	}

	// Penalized for official disciplinary action under Decree 48/2023
	sub.ViolatedTutoringRules = false
	sub.HasDisciplinaryAction = true
	result2 := engine.EvaluateIndividual(sub)
	if result2.OfficialRank != RankUnsatisfactory {
		t.Fatalf("expected rank %s, got %s", RankUnsatisfactory, result2.OfficialRank)
	}
}

func TestEvaluateIndividual_InnovativeSolution(t *testing.T) {
	engine := NewEngine()

	// High score (92) but NO innovative solution -> Capped at RankGood
	subWithoutInnovation := StaffKPISubmission{
		StaffID:                 uuid.New(),
		StaffName:               "Trần Thị B",
		AutomatedSystemScore:    92.0,
		ManagerAssessmentScore:  92.0,
		PeerReviewScore:         92.0,
		ParentSatisfactionScore: 92.0,
		HasInnovativeSolution:   false,
	}
	res1 := engine.EvaluateIndividual(subWithoutInnovation)
	if res1.OfficialRank != RankGood {
		t.Fatalf("expected rank %s without innovation, got %s", RankGood, res1.OfficialRank)
	}

	// High score (92) WITH innovative solution -> RankExcellent
	subWithoutInnovation.HasInnovativeSolution = true
	res2 := engine.EvaluateIndividual(subWithoutInnovation)
	if res2.OfficialRank != RankExcellent {
		t.Fatalf("expected rank %s with innovation, got %s", RankExcellent, res2.OfficialRank)
	}
}

func TestEvaluateSchoolBatch_Enforce20PercentCap(t *testing.T) {
	engine := NewEngine()

	// Create 10 teachers with high scores and innovative solutions
	submissions := make([]StaffKPISubmission, 10)
	for i := 0; i < 10; i++ {
		score := 90.0 + float64(i) // 90 to 99
		submissions[i] = StaffKPISubmission{
			StaffID:                 uuid.New(),
			StaffName:               "Teacher",
			AutomatedSystemScore:    score,
			ManagerAssessmentScore:  score,
			PeerReviewScore:         score,
			ParentSatisfactionScore: score,
			HasInnovativeSolution:   true,
		}
	}

	results, err := engine.EvaluateSchoolBatch(submissions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	excellentCount := 0
	for _, res := range results {
		if res.OfficialRank == RankExcellent {
			excellentCount++
		}
	}

	// In a cohort of 10 qualified candidates, 20% cap = 2 teachers maximum
	if excellentCount > 2 {
		t.Fatalf("expected max 2 excellent rankings (20%% cap), got %d", excellentCount)
	}
	if excellentCount != 2 {
		t.Fatalf("expected exactly 2 excellent rankings, got %d", excellentCount)
	}
}
