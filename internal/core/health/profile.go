package health

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Categories for Primary School BMI (5-11 years old)
const (
	BMICategoryUnderweight = "Suy dinh dưỡng"
	BMICategoryNormal      = "Bình thường"
	BMICategoryOverweight  = "Thừa cân"
	BMICategoryObese       = "Béo phì"

	// Thresholds for Vietnamese Primary School Children
	BMIThresholdUnderweight = 14.0
	BMIThresholdNormal      = 18.5
	BMIThresholdOverweight  = 21.0
)

// MedicalEvent represents an entry in the school infirmary medical diary
type MedicalEvent struct {
	ID              string `json:"id,omitempty"`
	Date            string `json:"date"`
	EventType       string `json:"event_type"` // SƠ_CẤP_CỨU, SỐT, ĐAU_BỤNG, TRẦY_XƯỚC, DỊ_ỨNG, KHÁM_ĐỊNH_KỲ
	Symptoms        string `json:"symptoms,omitempty"`
	Treatment       string `json:"treatment"`
	MedicationGiven string `json:"medication_given,omitempty"`
	NurseName       string `json:"nurse_name,omitempty"`
	ParentNotified  bool   `json:"parent_notified"`
}

// PrescriptionItem represents a parent-requested medication delivery at school
type PrescriptionItem struct {
	ID                 string     `json:"id"`
	MedicationName     string     `json:"medication_name"`
	Dosage             string     `json:"dosage"`
	AdministrationTime string     `json:"administration_time"` // e.g. "12:30 sau ăn trưa"
	Instructions       string     `json:"instructions,omitempty"`
	ParentNotes        string     `json:"parent_notes,omitempty"`
	Status             string     `json:"status"` // PENDING, ADMINISTERED, COMPLETED, REJECTED
	AdministeredAt     *time.Time `json:"administered_at,omitempty"`
	AdministeredBy     string     `json:"administered_by,omitempty"`
}

// StudentHealthProfile represents a comprehensive health snapshot
type StudentHealthProfile struct {
	StudentID              uuid.UUID          `json:"student_id"`
	CheckupDate            string             `json:"checkup_date"`
	HeightCm               float64            `json:"height_cm"`
	WeightKg               float64            `json:"weight_kg"`
	BMI                    float64            `json:"bmi"`
	BMICategory            string             `json:"bmi_category"`
	BMIAdvice              string             `json:"bmi_advice"`
	LeftEyeVision          string             `json:"left_eye_vision"`
	RightEyeVision         string             `json:"right_eye_vision"`
	EyesightStatus         string             `json:"eyesight_status"`
	EyesightRecommendation string             `json:"eyesight_recommendation"`
	DentalCondition        string             `json:"dental_condition"`
	ChronicConditions      string             `json:"chronic_conditions"`
	MedicalEvents          []MedicalEvent     `json:"medical_events"`
	Prescriptions          []PrescriptionItem `json:"prescriptions"`
	UpdatedAt              time.Time          `json:"updated_at"`
}

// Engine encapsulates health metrics calculation and evaluations
type Engine struct{}

// NewEngine creates a new Health Engine instance
func NewEngine() *Engine {
	return &Engine{}
}

// round1 rounds float64 to 1 decimal place
func round1(val float64) float64 {
	return math.Round(val*10.0) / 10.0
}

// CalculateBMI computes Body Mass Index = weight (kg) / (height (m))^2
func (e *Engine) CalculateBMI(heightCm, weightKg float64) float64 {
	if heightCm <= 0 || weightKg <= 0 {
		return 0.0
	}
	heightM := heightCm / 100.0
	bmi := weightKg / (heightM * heightM)
	return round1(bmi)
}

// ClassifyBMICategory determines nutritional status according to primary school standards
func (e *Engine) ClassifyBMICategory(bmi float64) (category string, advice string) {
	if bmi <= 0 {
		return "Chưa xác định", "Cần đo lại chiều cao và cân nặng hợp lệ."
	}

	if bmi < BMIThresholdUnderweight {
		return BMICategoryUnderweight, "Học sinh có dấu hiệu nhẹ cân, suy dinh dưỡng. Khuyến nghị gia đình và nhà trường bổ sung vi chất và theo dõi chế độ ăn bán trú."
	} else if bmi < BMIThresholdNormal {
		return BMICategoryNormal, "Thể trạng và dinh dưỡng phát triển tốt, đạt chuẩn thể lực học đường."
	} else if bmi < BMIThresholdOverweight {
		return BMICategoryOverweight, "Có dấu hiệu thừa cân. Khuyến nghị tăng cường hoạt động thể dục thể thao và giảm đồ ngọt/dầu mỡ."
	}

	return BMICategoryObese, "Nguy cơ béo phì học đường cao. Cần bác sĩ dinh dưỡng tư vấn khẩu phần và lộ trình vận động phù hợp."
}

// parseVisionToDecimal converts vision strings like "10/10", "7/10", "0.7", "10" to float fraction 0.0 - 1.0
func parseVisionToDecimal(v string) float64 {
	clean := strings.TrimSpace(v)
	if clean == "" {
		return 1.0
	}
	if strings.Contains(clean, "/") {
		parts := strings.Split(clean, "/")
		if len(parts) == 2 {
			numerator, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			denominator, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil && denominator > 0 {
				return numerator / denominator
			}
		}
	}
	val, err := strconv.ParseFloat(clean, 64)
	if err == nil {
		if val > 1.0 && val <= 10.0 {
			return val / 10.0
		}
		return val
	}
	return 1.0
}

// EvaluateEyesight evaluates left and right eye vision for myopia/refractive disorders
func (e *Engine) EvaluateEyesight(leftVision, rightVision string) (status string, hasDefect bool, recommendation string) {
	leftDec := parseVisionToDecimal(leftVision)
	rightDec := parseVisionToDecimal(rightVision)

	minDec := math.Min(leftDec, rightDec)

	if minDec >= 1.0 {
		return "Thị lực tốt, đạt chuẩn học đường", false, "Duy trì thói quen giữ khoảng cách đọc sách và hạn chế màn hình điện tử."
	} else if minDec >= 0.7 {
		return "Thị lực giảm nhẹ, nghi ngờ khúc xạ học đường", true, "Khuyến nghị phụ huynh cho học sinh đi đo khám thị lực tại cơ sở nhãn khoa và điều chỉnh khoảng cách bàn học."
	}

	return "Cận thị / Tật khúc xạ", true, "Thị lực suy giảm dưới 7/10. Cần khám chuyên khoa mắt để cắt kính phù hợp và bố trí vị trí ngồi ở hàng ghế đầu trong lớp."
}

// BuildHealthProfile constructs a complete health profile with BMI and eyesight evaluations
func (e *Engine) BuildHealthProfile(
	studentID uuid.UUID,
	checkupDate time.Time,
	heightCm, weightKg float64,
	leftVision, rightVision string,
	dental, chronic string,
	events []MedicalEvent,
	prescriptions []PrescriptionItem,
) StudentHealthProfile {
	bmi := e.CalculateBMI(heightCm, weightKg)
	cat, advice := e.ClassifyBMICategory(bmi)
	eyeStatus, _, eyeRec := e.EvaluateEyesight(leftVision, rightVision)

	dateStr := checkupDate.Format("2006-01-02")
	if checkupDate.IsZero() {
		dateStr = time.Now().Format("2006-01-02")
	}

	return StudentHealthProfile{
		StudentID:              studentID,
		CheckupDate:            dateStr,
		HeightCm:               round1(heightCm),
		WeightKg:               round1(weightKg),
		BMI:                    bmi,
		BMICategory:            cat,
		BMIAdvice:              advice,
		LeftEyeVision:          leftVision,
		RightEyeVision:         rightVision,
		EyesightStatus:         eyeStatus,
		EyesightRecommendation: eyeRec,
		DentalCondition:        dental,
		ChronicConditions:      chronic,
		MedicalEvents:          events,
		Prescriptions:          prescriptions,
		UpdatedAt:              time.Now(),
	}
}

// AdministerPrescription simulates school nurse administering prescribed medicine to a student
func (e *Engine) AdministerPrescription(
	item *PrescriptionItem,
	nurseName string,
) error {
	if item == nil {
		return fmt.Errorf("đơn thuốc không tồn tại")
	}
	if item.Status == "ADMINISTERED" || item.Status == "COMPLETED" {
		return fmt.Errorf("đơn thuốc này đã được cấp phát trước đó")
	}

	now := time.Now()
	item.Status = "ADMINISTERED"
	item.AdministeredAt = &now
	item.AdministeredBy = nurseName
	return nil
}
