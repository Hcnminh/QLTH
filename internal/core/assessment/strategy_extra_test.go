package assessment

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"qlth/internal/core/domain"
)

// ValidateSubjectScore phải từ chối mức đạt ngoài tập HTT/HT/CHT
func TestValidateSubjectScore_InvalidAchievementLevel(t *testing.T) {
	strat := NewStrategyTT27Primary()

	err := strat.ValidateSubjectScore(domain.StudentSubjectEvaluation{
		AchievementLevel: "GIOI", // không hợp lệ theo TT 27/2020
	})
	if err == nil {
		t.Fatal("expected error for invalid achievement level, got nil")
	}
}

// Điểm kiểm tra định kỳ chỉ hợp lệ trong khoảng 1.0 - 10.0
func TestValidateSubjectScore_ExamScoreOutOfRange(t *testing.T) {
	strat := NewStrategyTT27Primary()

	zero := 0.0
	if err := strat.ValidateSubjectScore(domain.StudentSubjectEvaluation{
		AchievementLevel: "HT",
		ExamScore:        &zero,
	}); err == nil {
		t.Fatal("expected error for exam score 0 (below 1.0)")
	}

	ten := 10.5
	if err := strat.ValidateSubjectScore(domain.StudentSubjectEvaluation{
		AchievementLevel: "HT",
		ExamScore:        &ten,
	}); err == nil {
		t.Fatal("expected error for exam score 10.5 (above 10.0)")
	}
}

// TT 89/2026: học sinh khuyết tật hòa nhập được đánh giá theo KHGD cá nhân
func TestStrategyTT89_InclusiveEducationStudent(t *testing.T) {
	strat := NewStrategyTT89Primary()

	student := domain.Student{
		ID:                   uuid.New(),
		FullName:             "Nguyễn Văn Hòa",
		IsInclusiveEducation: true, // học sinh khuyết tật học hòa nhập
	}

	// Không có dữ liệu môn học — học sinh hòa nhập đánh giá theo IEP
	rank, reason, err := strat.EvaluateYearEnd(student, nil, domain.StudentBehaviorEvaluation{})
	if err != nil {
		t.Fatalf("unexpected error for inclusive student: %v", err)
	}
	if rank != RankStudentSatisfactory {
		t.Fatalf("inclusive student should be ranked %s per IEP, got %s", RankStudentSatisfactory, rank)
	}
	if reason == "" {
		t.Fatal("expected a reason referencing kế hoạch giáo dục cá nhân")
	}
}

// TT 27/2020: học sinh thiếu dữ liệu môn học phải trả lỗi, không xếp loại bừa
func TestStrategyTT27_NoSubjectData(t *testing.T) {
	strat := NewStrategyTT27Primary()

	_, _, err := strat.EvaluateYearEnd(domain.Student{ID: uuid.New()}, nil, domain.StudentBehaviorEvaluation{})
	if err == nil {
		t.Fatal("expected error when no subject evaluations provided")
	}
}

// TT 27/2020: xếp loại Hoàn thành khi môn chỉ đạt HT và điểm >= 5
func TestStrategyTT27_SatisfactoryRank(t *testing.T) {
	strat := NewStrategyTT27Primary()

	student := domain.Student{ID: uuid.New(), FullName: "Lê Bình An"}
	score6 := 6.0
	subjects := []domain.StudentSubjectEvaluation{
		{SubjectID: "TOAN", AchievementLevel: "HT", ExamScore: &score6},
	}
	behavior := domain.StudentBehaviorEvaluation{
		Patriotism: "D", Compassion: "D", Diligence: "D", Honesty: "D", Responsibility: "D",
		AutonomySelfLearning: "D", CommunicationCooperation: "D", ProblemSolvingCreativity: "D",
	}

	rank, _, err := strat.EvaluateYearEnd(student, subjects, behavior)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rank != RankStudentSatisfactory {
		t.Fatalf("expected %s, got %s", RankStudentSatisfactory, rank)
	}
}

// TT 27/2020: snapshot học bạ phải chứa đầy đủ trường bắt buộc để ký số
func TestStrategyTT27_GenerateTranscriptSnapshot(t *testing.T) {
	strat := NewStrategyTT27Primary()

	student := domain.Student{
		ID:          uuid.New(),
		StudentCode: "HS003",
		FullName:    "Phạm Thu Hà",
	}
	score10 := 10.0
	subjects := []domain.StudentSubjectEvaluation{
		{SubjectID: "TIENG_VIET", AchievementLevel: "HTT", ExamScore: &score10},
	}
	behavior := domain.StudentBehaviorEvaluation{Patriotism: "T"}

	raw, err := strat.GenerateTranscriptSnapshot(student, "2026-2027", RankStudentExcellent, subjects, behavior)
	if err != nil {
		t.Fatalf("snapshot generation failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatalf("snapshot is not valid JSON: %v", err)
	}

	if data["regulation"] != "TT27_2020" {
		t.Fatalf("expected regulation TT27_2020, got %v", data["regulation"])
	}
	if data["annual_rank"] != string(RankStudentExcellent) {
		t.Fatalf("expected annual_rank %s, got %v", RankStudentExcellent, data["annual_rank"])
	}
	if data["student_code"] != "HS003" {
		t.Fatalf("expected student_code HS003, got %v", data["student_code"])
	}
	if _, ok := data["snapshot_created"]; !ok {
		t.Fatal("snapshot must include snapshot_created timestamp for signing")
	}
}

// TT 89/2026: học sinh thường (không hòa nhập) được đánh giá như TT 27 qua base strategy
func TestStrategyTT89_RegularStudentFallsThroughToBase(t *testing.T) {
	strat := NewStrategyTT89Primary()

	student := domain.Student{
		ID:                   uuid.New(),
		FullName:             "Vũ Đức Long",
		IsInclusiveEducation: false,
	}

	score10 := 10.0
	subjects := []domain.StudentSubjectEvaluation{
		{SubjectID: "TOAN", AchievementLevel: "HTT", ExamScore: &score10, IsComputerBasedTest: true},
		{SubjectID: "TIENG_VIET", AchievementLevel: "HTT", ExamScore: &score10, IsCentralizedExam: true},
	}
	behavior := domain.StudentBehaviorEvaluation{
		Patriotism: "T", Compassion: "T", Diligence: "T", Honesty: "T", Responsibility: "T",
		AutonomySelfLearning: "T", CommunicationCooperation: "T", ProblemSolvingCreativity: "T",
	}

	rank, _, err := strat.EvaluateYearEnd(student, subjects, behavior)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rank != RankStudentExcellent {
		t.Fatalf("regular student with all HTT and scores 10 must be %s, got %s", RankStudentExcellent, rank)
	}
}

// TT 89/2026: ValidateSubjectScore ủy quyền cho base, vẫn chặn giá trị ngoài chuẩn
func TestStrategyTT89_ValidateSubjectScore(t *testing.T) {
	strat := NewStrategyTT89Primary()

	// Hợp lệ
	if err := strat.ValidateSubjectScore(domain.StudentSubjectEvaluation{
		AchievementLevel:    "HTT",
		IsComputerBasedTest: true,
	}); err != nil {
		t.Fatalf("valid input rejected: %v", err)
	}

	// Mức đạt sai chuẩn TT 27/89
	if err := strat.ValidateSubjectScore(domain.StudentSubjectEvaluation{
		AchievementLevel: "DAT",
	}); err == nil {
		t.Fatal("expected error for invalid achievement level")
	}

	// Điểm kiểm tra vượt thang 10
	over := 11.0
	if err := strat.ValidateSubjectScore(domain.StudentSubjectEvaluation{
		AchievementLevel: "HT",
		ExamScore:        &over,
	}); err == nil {
		t.Fatal("expected error for exam score above 10.0")
	}
}

// Registry phải phân giải TT27 trước mốc 01/12/2026 cho cả 2 ngày biên
func TestRegistry_BoundaryDate(t *testing.T) {
	registry := NewRegistry()

	// Ngày cuối cùng TT27 còn hiệu lực: 30/11/2026
	lastTT27 := time.Date(2026, time.November, 30, 23, 59, 59, 0, time.Local)
	if got := registry.ResolveStrategy(lastTT27).GetRegulationCode(); got != "TT27_2020" {
		t.Fatalf("30/11/2026 must use TT27_2020, got %s", got)
	}

	// Đúng 00:00 ngày 01/12/2026 phải chuyển sang TT89
	exactTT89 := time.Date(2026, time.December, 1, 0, 0, 0, 0, time.Local)
	if got := registry.ResolveStrategy(exactTT89).GetRegulationCode(); got != "TT89_2026" {
		t.Fatalf("01/12/2026 00:00 must use TT89_2026, got %s", got)
	}
}
