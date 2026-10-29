package assessment

import (
	"testing"

	"qlth/internal/core/domain"

	"github.com/google/uuid"
)

func TestValidateBehavior_Valid(t *testing.T) {
	strat := NewStrategyTT27Primary()

	b := domain.StudentBehaviorEvaluation{
		Patriotism:               "T",
		Compassion:               "T",
		Diligence:                "D",
		Honesty:                  "T",
		Responsibility:           "D",
		AutonomySelfLearning:     "T",
		CommunicationCooperation: "D",
		ProblemSolvingCreativity: "T",
		LanguageCompetence:       "T",
		MathematicalCompetence:   "T",
	}

	if err := strat.ValidateBehavior(b); err != nil {
		t.Fatalf("expected valid behavior, got error: %v", err)
	}
}

func TestValidateBehavior_InvalidQuality(t *testing.T) {
	strat := NewStrategyTT27Primary()

	b := domain.StudentBehaviorEvaluation{
		Patriotism:     "INVALID",
		Compassion:     "T",
		Diligence:      "D",
		Honesty:        "T",
		Responsibility: "T",
	}

	if err := strat.ValidateBehavior(b); err == nil {
		t.Fatal("expected error for invalid quality level, got nil")
	}
}

func TestValidateBehavior_MissingQuality(t *testing.T) {
	strat := NewStrategyTT27Primary()

	b := domain.StudentBehaviorEvaluation{
		Patriotism:     "",
		Compassion:     "T",
		Diligence:      "D",
		Honesty:        "T",
		Responsibility: "T",
	}

	if err := strat.ValidateBehavior(b); err == nil {
		t.Fatal("expected error for empty patriotism, got nil")
	}
}

func TestValidateBehavior_InvalidCompetency(t *testing.T) {
	strat := NewStrategyTT27Primary()

	b := domain.StudentBehaviorEvaluation{
		Patriotism:             "T",
		Compassion:             "T",
		Diligence:              "D",
		Honesty:                "T",
		Responsibility:         "T",
		MathematicalCompetence: "KEM",
	}

	if err := strat.ValidateBehavior(b); err == nil {
		t.Fatal("expected error for invalid mathematical competence level, got nil")
	}
}

func TestStrategyTT89_ValidateBehavior(t *testing.T) {
	strat := NewStrategyTT89Primary()

	b := domain.StudentBehaviorEvaluation{
		Patriotism:               "T",
		Compassion:               "T",
		Diligence:                "T",
		Honesty:                  "T",
		Responsibility:           "T",
		AutonomySelfLearning:     "T",
		CommunicationCooperation: "T",
		ProblemSolvingCreativity: "T",
	}

	if err := strat.ValidateBehavior(b); err != nil {
		t.Fatalf("expected valid behavior under TT89, got %v", err)
	}
}

func TestEvaluateYearEnd_WithSpecificCompetencyC(t *testing.T) {
	strat := NewStrategyTT27Primary()

	student := domain.Student{
		ID:          uuid.New(),
		StudentCode: "HS_SPEC",
		FullName:    "Đặng Thu Thảo",
	}
	score10 := 10.0
	subjects := []domain.StudentSubjectEvaluation{
		{SubjectID: "TOAN", AchievementLevel: "HTT", ExamScore: &score10},
	}
	behavior := domain.StudentBehaviorEvaluation{
		Patriotism:               "T",
		Compassion:               "T",
		Diligence:                "T",
		Honesty:                  "T",
		Responsibility:           "T",
		AutonomySelfLearning:     "T",
		CommunicationCooperation: "T",
		ProblemSolvingCreativity: "T",
		PhysicalCompetence:       "C", // Cần cố gắng -> Phải bị xếp Chưa hoàn thành
	}

	rank, _, err := strat.EvaluateYearEnd(student, subjects, behavior)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rank != RankStudentIncomplete {
		t.Fatalf("expected Incomplete due to PhysicalCompetence 'C', got %s", rank)
	}
}
