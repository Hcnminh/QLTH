package kpi

import (
	"testing"

	"github.com/google/uuid"
)

// Điểm < 50 phải xếp Không hoàn thành nhiệm vụ
func TestEvaluateIndividual_BelowFifty(t *testing.T) {
	engine := NewEngine()

	sub := StaffKPISubmission{
		StaffID:                 uuid.New(),
		AutomatedSystemScore:    40.0,
		ManagerAssessmentScore:  40.0,
		PeerReviewScore:         40.0,
		ParentSatisfactionScore: 40.0,
	}

	result := engine.EvaluateIndividual(sub)
	if result.OfficialRank != RankUnsatisfactory {
		t.Fatalf("expected %s for score 40, got %s", RankUnsatisfactory, result.OfficialRank)
	}
	if result.FinalScore != 40.0 {
		t.Fatalf("expected final score 40.0, got %f", result.FinalScore)
	}
}

// Sự cố an toàn học đường nghiêm trọng là điểm liệt bất kể điểm KPI cao
func TestEvaluateIndividual_SevereSafetyIncident(t *testing.T) {
	engine := NewEngine()

	sub := StaffKPISubmission{
		StaffID:                 uuid.New(),
		AutomatedSystemScore:    98.0,
		ManagerAssessmentScore:  98.0,
		PeerReviewScore:         98.0,
		ParentSatisfactionScore: 98.0,
		HasSevereSafetyIncident: true,
	}

	result := engine.EvaluateIndividual(sub)
	if result.OfficialRank != RankUnsatisfactory {
		t.Fatalf("severe safety incident must be fatal flaw, got %s", result.OfficialRank)
	}
}

// Ngưỡng biên: đúng 75.0 → Hoàn thành tốt; đúng 50.0 → Hoàn thành
func TestEvaluateIndividual_ThresholdBoundaries(t *testing.T) {
	engine := NewEngine()

	base := StaffKPISubmission{
		StaffID:                 uuid.New(),
		AutomatedSystemScore:    75.0,
		ManagerAssessmentScore:  75.0,
		PeerReviewScore:         75.0,
		ParentSatisfactionScore: 75.0,
	}
	if got := engine.EvaluateIndividual(base).OfficialRank; got != RankGood {
		t.Fatalf("score exactly 75.0 must be %s, got %s", RankGood, got)
	}

	base.AutomatedSystemScore = 50.0
	base.ManagerAssessmentScore = 50.0
	base.PeerReviewScore = 50.0
	base.ParentSatisfactionScore = 50.0
	if got := engine.EvaluateIndividual(base).OfficialRank; got != RankSatisfactory {
		t.Fatalf("score exactly 50.0 must be %s, got %s", RankSatisfactory, got)
	}

	// 74.99 làm tròn 2 chữ số → vẫn dưới 75 → Hoàn thành
	base.AutomatedSystemScore = 74.99
	base.ManagerAssessmentScore = 74.99
	base.PeerReviewScore = 74.99
	base.ParentSatisfactionScore = 74.99
	if got := engine.EvaluateIndividual(base).OfficialRank; got != RankSatisfactory {
		t.Fatalf("score 74.99 must be %s, got %s", RankSatisfactory, got)
	}
}

// Danh sách rỗng phải trả lỗi thay vì panic
func TestEvaluateSchoolBatch_EmptyList(t *testing.T) {
	engine := NewEngine()

	_, err := engine.EvaluateSchoolBatch(nil)
	if err == nil {
		t.Fatal("expected error for empty submissions list")
	}
}

// Đơn vị nhỏ: 2 người đạt tốt thì vẫn được tối thiểu 1 suất xuất sắc
func TestEvaluateSchoolBatch_SmallUnitMinimumOne(t *testing.T) {
	engine := NewEngine()

	submissions := []StaffKPISubmission{
		{
			StaffID:                 uuid.New(),
			AutomatedSystemScore:    95.0,
			ManagerAssessmentScore:  95.0,
			PeerReviewScore:         95.0,
			ParentSatisfactionScore: 95.0,
			HasInnovativeSolution:   true,
		},
		{
			StaffID:                 uuid.New(),
			AutomatedSystemScore:    80.0,
			ManagerAssessmentScore:  80.0,
			PeerReviewScore:         80.0,
			ParentSatisfactionScore: 80.0,
		},
	}

	results, err := engine.EvaluateSchoolBatch(submissions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	excellent := 0
	for _, r := range results {
		if r.OfficialRank == RankExcellent {
			excellent++
		}
	}
	// floor(2 * 0.20) = 0, nhưng quy tắc đơn vị nhỏ cho tối thiểu 1
	if excellent != 1 {
		t.Fatalf("small unit should allow minimum 1 excellent, got %d", excellent)
	}
}

// Người bị điểm liệt không được tính vào mẫu trần 20% (chỉ đếm người đạt tốt trở lên)
func TestEvaluateSchoolBatch_FatalFlawExcludedFromQuota(t *testing.T) {
	engine := NewEngine()

	submissions := make([]StaffKPISubmission, 0, 6)
	// 5 người đạt xuất sắc hợp lệ (score 90-94, có sáng kiến)
	for i := 0; i < 5; i++ {
		score := 90.0 + float64(i)
		submissions = append(submissions, StaffKPISubmission{
			StaffID:                 uuid.New(),
			AutomatedSystemScore:    score,
			ManagerAssessmentScore:  score,
			PeerReviewScore:         score,
			ParentSatisfactionScore: score,
			HasInnovativeSolution:   true,
		})
	}
	// 1 người điểm rất cao nhưng vi phạm dạy thêm → điểm liệt
	submissions = append(submissions, StaffKPISubmission{
		StaffID:                 uuid.New(),
		AutomatedSystemScore:    99.0,
		ManagerAssessmentScore:  99.0,
		PeerReviewScore:         99.0,
		ParentSatisfactionScore: 99.0,
		HasInnovativeSolution:   true,
		ViolatedTutoringRules:   true,
	})

	results, err := engine.EvaluateSchoolBatch(submissions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Người cuối (index 5) phải là Không hoàn thành dù điểm cao nhất
	if results[5].OfficialRank != RankUnsatisfactory {
		t.Fatalf("tutoring violator must be %s, got %s", RankUnsatisfactory, results[5].OfficialRank)
	}

	// Mẫu trần 20% chỉ gồm 5 người đạt tốt trở lên → tối đa 1 suất xuất sắc
	excellent := 0
	for _, r := range results {
		if r.OfficialRank == RankExcellent {
			excellent++
		}
	}
	if excellent != 1 {
		t.Fatalf("expected exactly 1 excellent (20%% of 5), got %d", excellent)
	}
	// Người điểm cao nhất trong nhóm hợp lệ (94) giữ suất xuất sắc
	if results[4].OfficialRank != RankExcellent {
		t.Fatalf("highest valid scorer must keep excellent rank, got %s", results[4].OfficialRank)
	}
}
