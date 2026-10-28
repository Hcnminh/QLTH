package assessment

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"qlth/internal/core/domain"
)

func TestRegistry_StrategyResolutionByDate(t *testing.T) {
	registry := NewRegistry()

	// 1. Date before 01/12/2026 (e.g. 20/10/2026) -> should return TT27_2020
	dateBefore := time.Date(2026, time.October, 20, 0, 0, 0, 0, time.Local)
	stratBefore := registry.ResolveStrategy(dateBefore)
	if stratBefore.GetRegulationCode() != "TT27_2020" {
		t.Fatalf("expected TT27_2020, got %s", stratBefore.GetRegulationCode())
	}

	// 2. Date on or after 01/12/2026 (e.g. 05/12/2026) -> should return TT89_2026
	dateAfter := time.Date(2026, time.December, 5, 0, 0, 0, 0, time.Local)
	stratAfter := registry.ResolveStrategy(dateAfter)
	if stratAfter.GetRegulationCode() != "TT89_2026" {
		t.Fatalf("expected TT89_2026, got %s", stratAfter.GetRegulationCode())
	}
}

func TestStrategyTT27_EvaluateYearEnd(t *testing.T) {
	strat := NewStrategyTT27Primary()

	student := domain.Student{
		ID:          uuid.New(),
		StudentCode: "HS001",
		FullName:    "Trần Minh Trí",
	}

	score10 := 10.0
	score9 := 9.0
	score4 := 4.0

	subjectsExcellent := []domain.StudentSubjectEvaluation{
		{SubjectID: "TOAN", AchievementLevel: "HTT", ExamScore: &score10},
		{SubjectID: "TIENG_VIET", AchievementLevel: "HTT", ExamScore: &score9},
		{SubjectID: "TIENG_ANH", AchievementLevel: "HTT"},
	}

	behaviorT := domain.StudentBehaviorEvaluation{
		Patriotism:               "T",
		Compassion:               "T",
		Diligence:                "T",
		Honesty:                  "T",
		Responsibility:           "T",
		AutonomySelfLearning:     "T",
		CommunicationCooperation: "T",
		ProblemSolvingCreativity: "T",
	}

	// Case 1: Excellent
	rank1, _, err := strat.EvaluateYearEnd(student, subjectsExcellent, behaviorT)
	if err != nil || rank1 != RankStudentExcellent {
		t.Fatalf("expected %s, got %s (err: %v)", RankStudentExcellent, rank1, err)
	}

	// Case 2: Incomplete due to exam score < 5.0
	subjectsIncomplete := []domain.StudentSubjectEvaluation{
		{SubjectID: "TOAN", AchievementLevel: "CHT", ExamScore: &score4},
		{SubjectID: "TIENG_VIET", AchievementLevel: "HTT", ExamScore: &score9},
	}
	rank2, _, _ := strat.EvaluateYearEnd(student, subjectsIncomplete, behaviorT)
	if rank2 != RankStudentIncomplete {
		t.Fatalf("expected %s, got %s", RankStudentIncomplete, rank2)
	}
}

func TestStrategyTT89_DigitalTranscriptSnapshot(t *testing.T) {
	strat := NewStrategyTT89Primary()

	student := domain.Student{
		ID:                uuid.New(),
		StudentCode:       "HS002",
		NationalID:        "001206012345",
		FullName:          "Lê Hoàng Nam",
		HasDigitalConsent: true,
	}

	subjects := []domain.StudentSubjectEvaluation{
		{SubjectID: "TOAN", AchievementLevel: "HTT", IsComputerBasedTest: true},
	}

	behavior := domain.StudentBehaviorEvaluation{
		Patriotism: "T",
	}

	bytes, err := strat.GenerateTranscriptSnapshot(student, "2026-2027", RankStudentExcellent, subjects, behavior)
	if err != nil {
		t.Fatalf("snapshot failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		t.Fatalf("unmarshal snapshot failed: %v", err)
	}

	if data["regulation"] != "TT89_2026" {
		t.Fatalf("expected regulation TT89_2026, got %v", data["regulation"])
	}
	if data["digital_transcript_status"] != "OFFICIAL_DIGITAL_TRANSCRIPT_CIRCULAR_89" {
		t.Fatalf("missing digital transcript status: %v", data["digital_transcript_status"])
	}
}
