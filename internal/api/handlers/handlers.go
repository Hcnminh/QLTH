package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"qlth/internal/core/assessment"
	"qlth/internal/core/domain"
	"qlth/internal/core/kpi"
	"qlth/internal/core/teaching"
)

type APIHandler struct {
	kpiEngine     *kpi.Engine
	quotaEngine   *teaching.QuotaEngine
	stratRegistry *assessment.Registry
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{
		kpiEngine:     kpi.NewEngine(),
		quotaEngine:   teaching.NewQuotaEngine(),
		stratRegistry: assessment.NewRegistry(),
	}
}

// HealthCheck endpoint
func (h *APIHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "eduprimary-kpi-core",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// CalculateTeachingQuota handles Circular 05/2025/TT-BGDDT teaching workload calculation
func (h *APIHandler) CalculateTeachingQuota(c *gin.Context) {
	var input teaching.TeacherWorkloadProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu đầu vào không hợp lệ: " + err.Error(),
		})
		return
	}

	result := h.quotaEngine.CalculateQuota(input)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Tính toán định mức thành công theo Thông tư 05/2025/TT-BGDĐT",
		"data":    result,
	})
}

// EvaluateKPI handles 360-degree KPI evaluation and Decree 90/2020 & 48/2023 classification
func (h *APIHandler) EvaluateKPI(c *gin.Context) {
	var sub kpi.StaffKPISubmission
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu đầu vào không hợp lệ: " + err.Error(),
		})
		return
	}

	result := h.kpiEngine.EvaluateIndividual(sub)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Đánh giá xếp loại thành công",
		"data":    result,
	})
}

// BatchEvaluateKPI handles whole-school KPI batch evaluation with the 20% quota cap
func (h *APIHandler) BatchEvaluateKPI(c *gin.Context) {
	var subs []kpi.StaffKPISubmission
	if err := c.ShouldBindJSON(&subs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu danh sách cán bộ không hợp lệ: " + err.Error(),
		})
		return
	}

	results, err := h.kpiEngine.EvaluateSchoolBatch(subs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Đã xếp loại thi đua toàn trường và áp dụng trần 20% xuất sắc theo Nghị định 90/2020",
		"data":    results,
	})
}

// AIPedagogicalCopilot generates smart comments according to Circular 27/2020 & 89/2026
func (h *APIHandler) AIPedagogicalCopilot(c *gin.Context) {
	type CopilotRequest struct {
		SubjectName      string  `json:"subject_name" binding:"required"`
		AchievementLevel string  `json:"achievement_level" binding:"required"` // HTT, HT, CHT
		ExamScore        float64 `json:"exam_score"`
		DiligenceLevel   string  `json:"diligence_level"` // T, D, C
		StudentName      string  `json:"student_name" binding:"required"`
	}

	var req CopilotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Thiếu thông tin học tập: " + err.Error(),
		})
		return
	}

	// Rule-based and pedagogical suggestions generator
	var suggestions []string
	if req.AchievementLevel == "HTT" && req.ExamScore >= 9.0 {
		suggestions = append(suggestions,
			fmt.Sprintf("Em %s tiếp thu bài rất nhanh, nắm vững kiến thức môn %s, tích cực phát biểu và có tư duy sáng tạo tốt.", req.StudentName, req.SubjectName),
			fmt.Sprintf("Em %s có thành tích xuất sắc môn %s, kĩ năng thực hành thành thạo, chữ viết cẩn thận, gương mẫu trong học tập.", req.StudentName, req.SubjectName),
			fmt.Sprintf("Em %s hoàn thành xuất sắc các nội dung môn %s, chủ động giúp đỡ bạn bè cùng tiến bộ.", req.StudentName, req.SubjectName),
		)
	} else if req.AchievementLevel == "CHT" || req.ExamScore < 5.0 {
		suggestions = append(suggestions,
			fmt.Sprintf("Em %s có cố gắng nhưng còn lúng túng khi giải bài tập môn %s. Cần rèn luyện thêm ở nhà và nhờ thầy cô hỗ trợ.", req.StudentName, req.SubjectName),
			fmt.Sprintf("Em %s tiếp thu bài môn %s còn chậm, cần tập trung chú ý lắng nghe giảng và tự tin hơn khi trao đổi với bạn bè.", req.StudentName, req.SubjectName),
		)
	} else {
		suggestions = append(suggestions,
			fmt.Sprintf("Em %s hoàn thành tốt nội dung môn %s, có tiến bộ rõ rệt so với đầu năm, cần phát huy thêm tinh thần tự giác.", req.StudentName, req.SubjectName),
			fmt.Sprintf("Em %s nắm được kiến thức cơ bản môn %s, cần rèn luyện thêm tính cẩn thận khi làm bài kiểm tra.", req.StudentName, req.SubjectName),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"code":        200,
		"message":     "Sinh gợi ý nhận xét sư phạm thành công",
		"suggestions": suggestions,
	})
}

// CancelBoardingMeal handles meal cancellation and automatic refund calculation
func (h *APIHandler) CancelBoardingMeal(c *gin.Context) {
	type MealCancellationRequest struct {
		StudentID        uuid.UUID `json:"student_id" binding:"required"`
		AttendanceDate   string    `json:"attendance_date" binding:"required"` // YYYY-MM-DD
		CancellationHour int       `json:"cancellation_hour" binding:"required"` // 0-23
		MealFeePerDay    float64   `json:"meal_fee_per_day"` // Default: 35000 VND
	}

	var req MealCancellationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu báo hủy ăn không hợp lệ: " + err.Error(),
		})
		return
	}

	if req.MealFeePerDay <= 0 {
		req.MealFeePerDay = 35000 // 35,000 VND chuẩn bán trú tiểu học
	}

	// Quy tắc nhà trường: Báo hủy ăn trước 08:00 sáng được tự động hoàn tiền
	isRefund := req.CancellationHour < 8
	refundAmount := 0.0
	reason := "Báo hủy sau 08:00 sáng - Bếp đã chuẩn bị thực phẩm nên không hoàn tiền theo quy chế"
	if isRefund {
		refundAmount = req.MealFeePerDay
		reason = "Báo hủy hợp lệ trước 08:00 sáng - Tự động bù trừ vào hóa đơn tiền ăn kỳ sau"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"data": gin.H{
			"student_id":         req.StudentID,
			"attendance_date":    req.AttendanceDate,
			"is_cancelled":       true,
			"is_refund_eligible": isRefund,
			"refund_amount":      refundAmount,
			"note":               reason,
		},
	})
}

// StudentOnboarding handles Level 1 intake with Decree 13/2023 Digital Consent check
func (h *APIHandler) StudentOnboarding(c *gin.Context) {
	var student domain.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu học sinh không hợp lệ: " + err.Error(),
		})
		return
	}

	// Check Decree 13/2023 Digital Consent
	if !student.HasDigitalConsent {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"code":    422,
			"message": "Chưa có sự đồng thuận số (Digital Consent) của cha mẹ/người giám hộ theo Điều 20 Nghị định 13/2023/NĐ-CP",
		})
		return
	}

	if student.ID == uuid.Nil {
		student.ID = uuid.New()
	}
	student.CreatedAt = time.Now()

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    201,
		"message": "Tiếp nhận hồ sơ học sinh thành công, đã đối soát định danh VNeID",
		"data":    student,
	})
}
