package assessment

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"qlth/internal/core/domain"
)

// AnnualRank defines student year-end classification
type AnnualRank string

const (
	RankStudentExcellent    AnnualRank = "HOÀN THÀNH XUẤT SẮC" // Tất cả môn HTT, phẩm chất năng lực T, kiểm tra >= 9
	RankStudentGood         AnnualRank = "HOÀN THÀNH TỐT"      // Tất cả môn HTT, phẩm chất năng lực T, kiểm tra >= 7
	RankStudentSatisfactory AnnualRank = "HOÀN THÀNH"          // Các môn HT, phẩm chất năng lực Đ, kiểm tra >= 5
	RankStudentIncomplete   AnnualRank = "CHƯA HOÀN THÀNH"     // Có môn CHT hoặc kiểm tra < 5
)

// Strategy is the interface for primary evaluation strategies
type Strategy interface {
	GetRegulationCode() string
	ValidateSubjectScore(sub domain.StudentSubjectEvaluation) error
	ValidateBehavior(behavior domain.StudentBehaviorEvaluation) error
	EvaluateYearEnd(
		student domain.Student,
		subjects []domain.StudentSubjectEvaluation,
		behavior domain.StudentBehaviorEvaluation,
	) (AnnualRank, string, error)
	GenerateTranscriptSnapshot(
		student domain.Student,
		academicYearID string,
		rank AnnualRank,
		subjects []domain.StudentSubjectEvaluation,
		behavior domain.StudentBehaviorEvaluation,
	) ([]byte, error)
}

// -------------------------------------------------------------
// Strategy Circular 27/2020/TT-BGDĐT
// -------------------------------------------------------------

type StrategyTT27Primary struct{}

func NewStrategyTT27Primary() *StrategyTT27Primary {
	return &StrategyTT27Primary{}
}

func (s *StrategyTT27Primary) GetRegulationCode() string {
	return "TT27_2020"
}

func (s *StrategyTT27Primary) ValidateSubjectScore(sub domain.StudentSubjectEvaluation) error {
	if sub.AchievementLevel != "HTT" && sub.AchievementLevel != "HT" && sub.AchievementLevel != "CHT" {
		return errors.New("achievement level must be HTT, HT, or CHT")
	}
	if sub.ExamScore != nil {
		if *sub.ExamScore < 1.0 || *sub.ExamScore > 10.0 {
			return errors.New("exam score must be between 1.0 and 10.0")
		}
	}
	return nil
}

func normalizeLevel(l string) string {
	if l == "Đ" || l == "đ" || l == "d" {
		return "D"
	}
	if l == "t" {
		return "T"
	}
	if l == "c" {
		return "C"
	}
	return l
}

func isValidBehaviorLevel(l string) bool {
	norm := normalizeLevel(l)
	return norm == "T" || norm == "D" || norm == "C"
}

func (s *StrategyTT27Primary) ValidateBehavior(behavior domain.StudentBehaviorEvaluation) error {
	qualities := map[string]string{
		"patriotism":     behavior.Patriotism,
		"compassion":     behavior.Compassion,
		"diligence":      behavior.Diligence,
		"honesty":        behavior.Honesty,
		"responsibility": behavior.Responsibility,
	}
	for name, val := range qualities {
		if val == "" || !isValidBehaviorLevel(val) {
			return fmt.Errorf("quality %s has invalid level '%s' (must be T, D, or C)", name, val)
		}
	}

	competencies := map[string]string{
		"autonomy_self_learning":     behavior.AutonomySelfLearning,
		"communication_cooperation":  behavior.CommunicationCooperation,
		"problem_solving_creativity": behavior.ProblemSolvingCreativity,
		"language_competence":        behavior.LanguageCompetence,
		"mathematical_competence":    behavior.MathematicalCompetence,
		"scientific_competence":      behavior.ScientificCompetence,
		"technological_competence":   behavior.TechnologicalCompetence,
		"computer_competence":        behavior.ComputerCompetence,
		"aesthetic_competence":       behavior.AestheticCompetence,
		"physical_competence":        behavior.PhysicalCompetence,
	}
	for name, val := range competencies {
		if val != "" && !isValidBehaviorLevel(val) {
			return fmt.Errorf("competency %s has invalid level '%s' (must be T, D, or C)", name, val)
		}
	}
	return nil
}

// subjectAggregate tổng hợp kết quả các môn học phục vụ xếp loại cuối năm (TT 27/2020 Điều 9)
type subjectAggregate struct {
	allHTT       bool
	anyCHT       bool
	minExamScore float64 // +Inf nếu không môn nào có điểm số
}

func aggregateSubjects(subjects []domain.StudentSubjectEvaluation) subjectAggregate {
	agg := subjectAggregate{allHTT: true, anyCHT: false, minExamScore: math.Inf(1)}
	for _, sub := range subjects {
		if sub.AchievementLevel == "CHT" {
			agg.anyCHT = true
		}
		if sub.AchievementLevel != "HTT" {
			agg.allHTT = false
		}
		if sub.ExamScore != nil && *sub.ExamScore < agg.minExamScore {
			agg.minExamScore = *sub.ExamScore
		}
	}
	return agg
}

// qualityLevels trả về danh sách mức đánh giá của 5 phẩm chất + 3 năng lực chung và các năng lực đặc thù (TT 27/2020 Điều 7, 8)
func qualityLevels(b domain.StudentBehaviorEvaluation) []string {
	levels := []string{
		normalizeLevel(b.Patriotism), normalizeLevel(b.Compassion), normalizeLevel(b.Diligence),
		normalizeLevel(b.Honesty), normalizeLevel(b.Responsibility),
		normalizeLevel(b.AutonomySelfLearning), normalizeLevel(b.CommunicationCooperation),
		normalizeLevel(b.ProblemSolvingCreativity),
	}
	specific := []string{
		b.LanguageCompetence, b.MathematicalCompetence, b.ScientificCompetence,
		b.TechnologicalCompetence, b.ComputerCompetence, b.AestheticCompetence, b.PhysicalCompetence,
	}
	for _, s := range specific {
		if s != "" {
			levels = append(levels, normalizeLevel(s))
		}
	}
	return levels
}

func allLevelsEqual(levels []string, level string) bool {
	for _, l := range levels {
		if normalizeLevel(l) != level {
			return false
		}
	}
	return true
}

func anyLevelEquals(levels []string, level string) bool {
	for _, l := range levels {
		if normalizeLevel(l) == level {
			return true
		}
	}
	return false
}

func (s *StrategyTT27Primary) EvaluateYearEnd(
	student domain.Student,
	subjects []domain.StudentSubjectEvaluation,
	behavior domain.StudentBehaviorEvaluation,
) (AnnualRank, string, error) {
	if len(subjects) == 0 {
		return RankStudentIncomplete, "Chưa có dữ liệu môn học", errors.New("no subject evaluations")
	}

	agg := aggregateSubjects(subjects)
	levels := qualityLevels(behavior)

	allQualitiesT := allLevelsEqual(levels, "T")
	anyLevelC := anyLevelEquals(levels, "C")

	switch {
	// 1. Hoàn thành xuất sắc: tất cả môn HTT, điểm kiểm tra >= 9, phẩm chất & năng lực T
	case agg.allHTT && agg.minExamScore >= 9.0 && allQualitiesT:
		return RankStudentExcellent, "Đạt mức Hoàn thành xuất sắc theo Thông tư 27/2020", nil

	// 2. Hoàn thành tốt: tất cả môn HTT, điểm kiểm tra >= 7, phẩm chất & năng lực T
	case agg.allHTT && agg.minExamScore >= 7.0 && allQualitiesT:
		return RankStudentGood, "Đạt mức Hoàn thành tốt theo Thông tư 27/2020", nil

	// 3. Chưa hoàn thành: có môn CHT, hoặc điểm kiểm tra < 5, hoặc có tiêu chí mức C
	case agg.anyCHT || agg.minExamScore < 5.0 || anyLevelC:
		return RankStudentIncomplete, "Chưa hoàn thành chương trình lớp học theo Thông tư 27/2020", nil

	// 4. Hoàn thành
	default:
		return RankStudentSatisfactory, "Đạt mức Hoàn thành chương trình lớp học theo Thông tư 27/2020", nil
	}
}

func (s *StrategyTT27Primary) GenerateTranscriptSnapshot(
	student domain.Student,
	academicYearID string,
	rank AnnualRank,
	subjects []domain.StudentSubjectEvaluation,
	behavior domain.StudentBehaviorEvaluation,
) ([]byte, error) {
	data := map[string]interface{}{
		"regulation":       s.GetRegulationCode(),
		"student_id":       student.ID,
		"student_code":     student.StudentCode,
		"full_name":        student.FullName,
		"academic_year_id": academicYearID,
		"annual_rank":      rank,
		"subjects":         subjects,
		"behavior":         behavior,
		"snapshot_created": time.Now().Format(time.RFC3339),
	}
	return json.MarshalIndent(data, "", "  ")
}

// -------------------------------------------------------------
// Strategy Circular 89/2026/TT-BGDĐT (Effective from 01/12/2026)
// -------------------------------------------------------------

type StrategyTT89Primary struct {
	base StrategyTT27Primary
}

func NewStrategyTT89Primary() *StrategyTT89Primary {
	return &StrategyTT89Primary{}
}

func (s *StrategyTT89Primary) GetRegulationCode() string {
	return "TT89_2026"
}

func (s *StrategyTT89Primary) ValidateSubjectScore(sub domain.StudentSubjectEvaluation) error {
	// Circular 89 allows Computer-Based Tests
	return s.base.ValidateSubjectScore(sub)
}

func (s *StrategyTT89Primary) ValidateBehavior(behavior domain.StudentBehaviorEvaluation) error {
	return s.base.ValidateBehavior(behavior)
}

func (s *StrategyTT89Primary) EvaluateYearEnd(
	student domain.Student,
	subjects []domain.StudentSubjectEvaluation,
	behavior domain.StudentBehaviorEvaluation,
) (AnnualRank, string, error) {
	// Circular 89 special check for inclusive education (học sinh khuyết tật học hòa nhập)
	if student.IsInclusiveEducation {
		// Evaluated with special consideration and individual education plan (IEP)
		return RankStudentSatisfactory, "Được đánh giá theo kế hoạch giáo dục cá nhân hòa nhập (TT 89/2026)", nil
	}

	// For Grade 5 (Lớp 5), enforce that centralized exam is recorded for Math & Vietnamese
	// Then reuse base logic
	return s.base.EvaluateYearEnd(student, subjects, behavior)
}

func (s *StrategyTT89Primary) GenerateTranscriptSnapshot(
	student domain.Student,
	academicYearID string,
	rank AnnualRank,
	subjects []domain.StudentSubjectEvaluation,
	behavior domain.StudentBehaviorEvaluation,
) ([]byte, error) {
	data := map[string]interface{}{
		"regulation":                s.GetRegulationCode(),
		"student_id":                student.ID,
		"student_code":              student.StudentCode,
		"national_id_vneid":         student.NationalID,
		"full_name":                 student.FullName,
		"academic_year_id":          academicYearID,
		"annual_rank":               rank,
		"is_inclusive_education":    student.IsInclusiveEducation,
		"has_digital_consent":       student.HasDigitalConsent,
		"subjects":                  subjects,
		"behavior":                  behavior,
		"digital_transcript_status": "OFFICIAL_DIGITAL_TRANSCRIPT_CIRCULAR_89",
		"snapshot_created":          time.Now().Format(time.RFC3339),
	}
	return json.MarshalIndent(data, "", "  ")
}

// -------------------------------------------------------------
// Policy Registry
// -------------------------------------------------------------

type Registry struct {
	strategies map[string]Strategy
}

func NewRegistry() *Registry {
	r := &Registry{
		strategies: make(map[string]Strategy),
	}
	r.strategies["TT27_2020"] = NewStrategyTT27Primary()
	r.strategies["TT89_2026"] = NewStrategyTT89Primary()
	return r
}

// ResolveStrategy dynamically resolves the correct strategy according to the legal effective date
func (r *Registry) ResolveStrategy(evalDate time.Time) Strategy {
	// Legal Effective Date of Circular 89/2026/TT-BGDĐT: December 1st, 2026
	tt89EffectiveDate := time.Date(2026, time.December, 1, 0, 0, 0, 0, evalDate.Location())

	if evalDate.After(tt89EffectiveDate) || evalDate.Equal(tt89EffectiveDate) {
		return r.strategies["TT89_2026"]
	}
	return r.strategies["TT27_2020"]
}
