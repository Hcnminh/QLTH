package kpi

import (
	"errors"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

// ClassificationResult defines rank according to Decree 90/2020 & Decree 48/2023
type ClassificationResult string

const (
	RankExcellent      ClassificationResult = "HOÀN THÀNH XUẤT SẮC NHIỆM VỤ"
	RankGood           ClassificationResult = "HOÀN THÀNH TỐT NHIỆM VỤ"
	RankSatisfactory   ClassificationResult = "HOÀN THÀNH NHIỆM VỤ"
	RankUnsatisfactory ClassificationResult = "KHÔNG HOÀN THÀNH NHIỆM VỤ"
)

// KPICriteria represents a measurable metric in the KPI tree
type KPICriteria struct {
	ID           uuid.UUID `json:"id"`
	Code         string    `json:"code"`
	Title        string    `json:"title"`
	Weight       float64   `json:"weight"`       // e.g. 0.20 for 20%
	MaxScore     float64   `json:"max_score"`    // 100
	SourceType   string    `json:"source_type"`  // 'SYSTEM_AUTOMATED', 'MANAGER_REVIEW', 'PEER_REVIEW', 'PARENT_CSAT'
	ActualScore  float64   `json:"actual_score"` // 0 - 100
}

// StaffKPISubmission is the assessment input for a specific teacher/staff
type StaffKPISubmission struct {
	StaffID                 uuid.UUID      `json:"staff_id"`
	StaffName               string         `json:"staff_name"`
	JobPosition             string         `json:"job_position"` // 'HOMEROOM_TEACHER', 'SUBJECT_TEACHER', 'PRINCIPAL'
	AutomatedSystemScore    float64        `json:"automated_system_score"`    // 40%
	ManagerAssessmentScore  float64        `json:"manager_assessment_score"`  // 30%
	PeerReviewScore         float64        `json:"peer_review_score"`         // 15%
	ParentSatisfactionScore float64        `json:"parent_satisfaction_score"` // 15%
	CriteriaDetails         []KPICriteria  `json:"criteria_details"`

	// Legal Constraints & Disciplinary conditions
	HasDisciplinaryAction   bool `json:"has_disciplinary_action"`   // Bị kỷ luật trong năm (NĐ 48/2023)
	ViolatedTutoringRules   bool `json:"violated_tutoring_rules"`   // Vi phạm cấm dạy thêm tiểu học (TT 29/2024 & TT 19/2026)
	HasSevereSafetyIncident bool `json:"has_severe_safety_incident"` // Mất an toàn học đường/ngộ độc
	HasInnovativeSolution   bool `json:"has_innovative_solution"`   // Có giải pháp sáng tạo/SKKN (Luật TĐKT 2022)
}

// FinalKPIEvaluation holds calculated scores and official rank
type FinalKPIEvaluation struct {
	StaffID            uuid.UUID            `json:"staff_id"`
	StaffName          string               `json:"staff_name"`
	FinalScore         float64              `json:"final_score"` // 0 - 100
	OfficialRank       ClassificationResult `json:"official_rank"`
	Reason             string               `json:"reason"`
	CalculatedAt       time.Time            `json:"calculated_at"`
}

// Engine implements the KPI calculation service
type Engine struct{}

// NewEngine creates an instance of the KPI Engine
func NewEngine() *Engine {
	return &Engine{}
}

// CalculateStaffScore calculates the 360-degree weighted final score
func (e *Engine) CalculateStaffScore(sub StaffKPISubmission) float64 {
	// Weights: 40% System + 30% Manager + 15% Peer + 15% Parent
	score := (sub.AutomatedSystemScore * 0.40) +
		(sub.ManagerAssessmentScore * 0.30) +
		(sub.PeerReviewScore * 0.15) +
		(sub.ParentSatisfactionScore * 0.15)

	return math.Round(score*100) / 100
}

// EvaluateIndividual determines the preliminary rank for one person
func (e *Engine) EvaluateIndividual(sub StaffKPISubmission) FinalKPIEvaluation {
	score := e.CalculateStaffScore(sub)

	// 1. Check disqualifying/fatal flaws (Điểm liệt kỷ luật)
	if sub.HasDisciplinaryAction {
		return FinalKPIEvaluation{
			StaffID:      sub.StaffID,
			StaffName:    sub.StaffName,
			FinalScore:   score,
			OfficialRank: RankUnsatisfactory,
			Reason:       "Bị kỷ luật Đảng hoặc chính quyền trong năm theo Nghị định 48/2023/NĐ-CP",
			CalculatedAt: time.Now(),
		}
	}

	if sub.ViolatedTutoringRules {
		return FinalKPIEvaluation{
			StaffID:      sub.StaffID,
			StaffName:    sub.StaffName,
			FinalScore:   score,
			OfficialRank: RankUnsatisfactory,
			Reason:       "Vi phạm quy định nghiêm cấm dạy thêm tiểu học (Thông tư 29/2024 & TT 19/2026)",
			CalculatedAt: time.Now(),
		}
	}

	if sub.HasSevereSafetyIncident {
		return FinalKPIEvaluation{
			StaffID:      sub.StaffID,
			StaffName:    sub.StaffName,
			FinalScore:   score,
			OfficialRank: RankUnsatisfactory,
			Reason:       "Để xảy ra tai nạn thương tích hoặc mất an toàn thực phẩm nghiêm trọng",
			CalculatedAt: time.Now(),
		}
	}

	// 2. Score >= 90: Check Excellence criteria
	if score >= 90.0 {
		if sub.HasInnovativeSolution {
			return FinalKPIEvaluation{
				StaffID:      sub.StaffID,
				StaffName:    sub.StaffName,
				FinalScore:   score,
				OfficialRank: RankExcellent,
				Reason:       "Điểm KPI >= 90 và có sáng kiến / giải pháp đổi mới sáng tạo được công nhận",
				CalculatedAt: time.Now(),
			}
		}
		// Without innovative solution, capped at Good (Hoàn thành tốt)
		return FinalKPIEvaluation{
			StaffID:      sub.StaffID,
			StaffName:    sub.StaffName,
			FinalScore:   score,
			OfficialRank: RankGood,
			Reason:       "Điểm KPI >= 90 nhưng chưa có giải pháp đổi mới sáng tạo được nghiệm thu",
			CalculatedAt: time.Now(),
		}
	}

	// 3. Score >= 75: Good
	if score >= 75.0 {
		return FinalKPIEvaluation{
			StaffID:      sub.StaffID,
			StaffName:    sub.StaffName,
			FinalScore:   score,
			OfficialRank: RankGood,
			Reason:       "Điểm KPI từ 75 đến dưới 90 điểm, hoàn thành tốt nhiệm vụ",
			CalculatedAt: time.Now(),
		}
	}

	// 4. Score >= 50: Satisfactory
	if score >= 50.0 {
		return FinalKPIEvaluation{
			StaffID:      sub.StaffID,
			StaffName:    sub.StaffName,
			FinalScore:   score,
			OfficialRank: RankSatisfactory,
			Reason:       "Điểm KPI từ 50 đến dưới 75 điểm",
			CalculatedAt: time.Now(),
		}
	}

	// 5. Score < 50: Unsatisfactory
	return FinalKPIEvaluation{
		StaffID:      sub.StaffID,
		StaffName:    sub.StaffName,
		FinalScore:   score,
		OfficialRank: RankUnsatisfactory,
		Reason:       "Điểm KPI dưới 50 điểm",
		CalculatedAt: time.Now(),
	}
}

// EvaluateSchoolBatch evaluates all staff in a school and applies the 20% cap on RankExcellent
// according to Decree 90/2020 & Decree 48/2023.
func (e *Engine) EvaluateSchoolBatch(submissions []StaffKPISubmission) ([]FinalKPIEvaluation, error) {
	if len(submissions) == 0 {
		return nil, errors.New("empty staff submissions list")
	}

	results := make([]FinalKPIEvaluation, len(submissions))
	for i, sub := range submissions {
		results[i] = e.EvaluateIndividual(sub)
	}

	// Calculate maximum number of Excellent ratings allowed (max 20% of Good + Excellent)
	// Decree 90/2020 Article 2: Tỷ lệ không vượt quá 20% tổng số viên chức được xếp loại hoàn thành tốt nhiệm vụ
	qualifiedGoodOrAbove := 0
	for _, res := range results {
		if res.OfficialRank == RankExcellent || res.OfficialRank == RankGood {
			qualifiedGoodOrAbove++
		}
	}

	maxExcellentAllowed := int(math.Floor(float64(qualifiedGoodOrAbove) * 0.20))
	if maxExcellentAllowed < 1 && qualifiedGoodOrAbove > 0 {
		maxExcellentAllowed = 1 // Trường hợp quy mô nhỏ, đơn vị được ít nhất 1 xuất sắc nếu đủ điều kiện
	}

	// Collect indexes of those with RankExcellent
	type candidate struct {
		index int
		score float64
	}
	var excellentCandidates []candidate
	for i, res := range results {
		if res.OfficialRank == RankExcellent {
			excellentCandidates = append(excellentCandidates, candidate{index: i, score: res.FinalScore})
		}
	}

	// If candidates exceed quota, keep the highest scores as Excellent, reclassify others to Good
	if len(excellentCandidates) > maxExcellentAllowed {
		sort.Slice(excellentCandidates, func(i, j int) bool {
			return excellentCandidates[i].score > excellentCandidates[j].score
		})

		for i := maxExcellentAllowed; i < len(excellentCandidates); i++ {
			targetIdx := excellentCandidates[i].index
			results[targetIdx].OfficialRank = RankGood
			results[targetIdx].Reason = "Đạt chuẩn xuất sắc nhưng điều chỉnh sang mức Tốt do khống chế trần 20% theo Nghị định 90/2020"
		}
	}

	return results, nil
}
