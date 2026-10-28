package assessment

import (
	"encoding/json"
	"errors"
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

// qualityLevels trả về danh sách mức đánh giá của 5 phẩm chất + 3 năng lực chung (TT 27/2020 Điều 7, 8)
func qualityLevels(b domain.StudentBehaviorEvaluation) []string {
	return []string{
		b.Patriotism, b.Compassion, b.Diligence, b.Honesty, b.Responsibility,
		b.AutonomySelfLearning, b.CommunicationCooperation, b.ProblemSolvingCreativity,
	}
}

func allLevelsEqual(levels []string, level string) bool {
	for _, l := range levels {
		if l != level {
			return false
		}
	}
	return true
}

func anyLevelEquals(levels []string, level string) bool {
	for _, l := range levels {
		if l == level {
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
	tt89EffectiveDate := time.Date(2026, time.December, 1, 0, 0, 0, 0, time.Local)

	if evalDate.After(tt89EffectiveDate) || evalDate.Equal(tt89EffectiveDate) {
		return r.strategies["TT89_2026"]
	}
	return r.strategies["TT27_2020"]
}
