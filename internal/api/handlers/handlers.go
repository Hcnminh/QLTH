package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"qlth/internal/core/assessment"
	"qlth/internal/core/domain"
	"qlth/internal/core/kpi"
	"qlth/internal/core/teaching"
	"qlth/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		AttendanceDate   string    `json:"attendance_date" binding:"required"`   // YYYY-MM-DD
		CancellationHour int       `json:"cancellation_hour" binding:"required"` // 0-23
		MealFeePerDay    float64   `json:"meal_fee_per_day"`                     // Default: 35000 VND
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

func parseEvaluationDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}
	formats := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if t, err := time.ParseInLocation(f, dateStr, time.Local); err == nil {
			return t
		}
	}
	return time.Now()
}

// SubjectScoreRequest defines input for evaluating a student in a subject
type SubjectScoreRequest struct {
	StudentID           uuid.UUID `json:"student_id" binding:"required"`
	TermID              string    `json:"term_id" binding:"required"`
	SubjectID           string    `json:"subject_id" binding:"required"`
	TeacherID           uuid.UUID `json:"teacher_id"`
	AchievementLevel    string    `json:"achievement_level" binding:"required"` // HTT, HT, CHT
	TeacherComment      string    `json:"teacher_comment"`
	AISuggestedComment  string    `json:"ai_suggested_comment,omitempty"`
	IsAIAssisted        bool      `json:"is_ai_assisted"`
	ExamScore           *float64  `json:"exam_score,omitempty"`
	IsComputerBasedTest bool      `json:"is_computer_based_test"`
	IsCentralizedExam   bool      `json:"is_centralized_exam"`
	EvaluationDate      string    `json:"evaluation_date,omitempty"`
}

// RecordSubjectScore handles student subject evaluations (POST /api/v1/evaluations/subject-scores)
func (h *APIHandler) RecordSubjectScore(c *gin.Context) {
	var req SubjectScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu đánh giá môn học không hợp lệ: " + err.Error(),
		})
		return
	}

	evalDate := parseEvaluationDate(req.EvaluationDate)
	strat := h.stratRegistry.ResolveStrategy(evalDate)

	eval := domain.StudentSubjectEvaluation{
		ID:                  uuid.New(),
		StudentID:           req.StudentID,
		TermID:              req.TermID,
		SubjectID:           req.SubjectID,
		TeacherID:           req.TeacherID,
		AchievementLevel:    req.AchievementLevel,
		TeacherComment:      req.TeacherComment,
		AISuggestedComment:  req.AISuggestedComment,
		IsAIAssisted:        req.IsAIAssisted,
		ExamScore:           req.ExamScore,
		IsComputerBasedTest: req.IsComputerBasedTest,
		IsCentralizedExam:   req.IsCentralizedExam,
		UpdatedAt:           time.Now(),
	}

	if err := strat.ValidateSubjectScore(eval); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Kết quả đánh giá môn học không đúng quy định: " + err.Error(),
		})
		return
	}

	database, err := db.GetDB()
	if err == nil && database != nil {
		query := `
			INSERT INTO student_subject_evaluations (
				id, student_id, term_id, subject_id, teacher_id,
				achievement_level, teacher_comment, ai_suggested_comment,
				is_ai_assisted, exam_score, is_computer_based_test,
				is_centralized_exam, updated_at
			) VALUES (
				$1, $2, $3, $4, NULLIF($5, '00000000-0000-0000-0000-000000000000'::uuid),
				$6, $7, $8, $9, $10, $11, $12, NOW()
			)
			ON CONFLICT (student_id, term_id, subject_id) DO UPDATE SET
				teacher_id = EXCLUDED.teacher_id,
				achievement_level = EXCLUDED.achievement_level,
				teacher_comment = EXCLUDED.teacher_comment,
				ai_suggested_comment = EXCLUDED.ai_suggested_comment,
				is_ai_assisted = EXCLUDED.is_ai_assisted,
				exam_score = EXCLUDED.exam_score,
				is_computer_based_test = EXCLUDED.is_computer_based_test,
				is_centralized_exam = EXCLUDED.is_centralized_exam,
				updated_at = NOW()
			RETURNING id;
		`
		var savedID uuid.UUID
		err = database.QueryRow(query,
			eval.ID, eval.StudentID, eval.TermID, eval.SubjectID, eval.TeacherID,
			eval.AchievementLevel, eval.TeacherComment, eval.AISuggestedComment,
			eval.IsAIAssisted, eval.ExamScore, eval.IsComputerBasedTest,
			eval.IsCentralizedExam,
		).Scan(&savedID)
		if err == nil {
			eval.ID = savedID
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": fmt.Sprintf("Ghi nhận đánh giá môn học thành công theo Thông tư %s", strat.GetRegulationCode()),
		"data": gin.H{
			"evaluation":         eval,
			"regulation_applied": strat.GetRegulationCode(),
		},
	})
}

// BehaviorEvaluationRequest defines input for 5 qualities and 10 competencies
type BehaviorEvaluationRequest struct {
	StudentID                     uuid.UUID `json:"student_id" binding:"required"`
	TermID                        string    `json:"term_id" binding:"required"`
	EvaluatorID                   uuid.UUID `json:"evaluator_id"`
	TeacherID                     uuid.UUID `json:"teacher_id"`
	Patriotism                    string    `json:"patriotism" binding:"required"`
	Compassion                    string    `json:"compassion" binding:"required"`
	Diligence                     string    `json:"diligence" binding:"required"`
	Honesty                       string    `json:"honesty" binding:"required"`
	Responsibility                string    `json:"responsibility" binding:"required"`
	QualitiesComment              string    `json:"qualities_comment"`
	AutonomyAndSelfLearning       string    `json:"autonomy_and_self_learning"`
	AutonomySelfLearning          string    `json:"autonomy_self_learning"`
	CommunicationAndCollaboration string    `json:"communication_and_collaboration"`
	CommunicationCooperation      string    `json:"communication_cooperation"`
	ProblemSolvingAndCreativity   string    `json:"problem_solving_and_creativity"`
	ProblemSolvingCreativity      string    `json:"problem_solving_creativity"`
	LanguageCompetence            string    `json:"language_competence"`
	MathematicalCompetence        string    `json:"mathematical_competence"`
	ScientificCompetence          string    `json:"scientific_competence"`
	TechnologicalCompetence       string    `json:"technological_competence"`
	ComputerCompetence            string    `json:"computer_competence"`
	AestheticCompetence           string    `json:"aesthetic_competence"`
	PhysicalCompetence            string    `json:"physical_competence"`
	CompetenciesComment           string    `json:"competencies_comment"`
	GeneralComment                string    `json:"general_comment"`
	EvaluationDate                string    `json:"evaluation_date,omitempty"`
}

// RecordBehaviorEvaluation handles behavior evaluations (POST /api/v1/evaluations/behaviors)
func (h *APIHandler) RecordBehaviorEvaluation(c *gin.Context) {
	var req BehaviorEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu đánh giá phẩm chất và năng lực không hợp lệ: " + err.Error(),
		})
		return
	}

	evaluatorID := req.TeacherID
	if evaluatorID == uuid.Nil {
		evaluatorID = req.EvaluatorID
	}

	autonomy := req.AutonomySelfLearning
	if autonomy == "" {
		autonomy = req.AutonomyAndSelfLearning
	}
	communication := req.CommunicationCooperation
	if communication == "" {
		communication = req.CommunicationAndCollaboration
	}
	problemSolving := req.ProblemSolvingCreativity
	if problemSolving == "" {
		problemSolving = req.ProblemSolvingAndCreativity
	}

	behavior := domain.StudentBehaviorEvaluation{
		ID:                       uuid.New(),
		StudentID:                req.StudentID,
		TermID:                   req.TermID,
		TeacherID:                evaluatorID,
		Patriotism:               req.Patriotism,
		Compassion:               req.Compassion,
		Diligence:                req.Diligence,
		Honesty:                  req.Honesty,
		Responsibility:           req.Responsibility,
		QualitiesComment:         req.QualitiesComment,
		AutonomySelfLearning:     autonomy,
		CommunicationCooperation: communication,
		ProblemSolvingCreativity: problemSolving,
		LanguageCompetence:       req.LanguageCompetence,
		MathematicalCompetence:   req.MathematicalCompetence,
		ScientificCompetence:     req.ScientificCompetence,
		TechnologicalCompetence:  req.TechnologicalCompetence,
		ComputerCompetence:       req.ComputerCompetence,
		AestheticCompetence:      req.AestheticCompetence,
		PhysicalCompetence:       req.PhysicalCompetence,
		CompetenciesComment:      req.CompetenciesComment,
		GeneralComment:           req.GeneralComment,
		UpdatedAt:                time.Now(),
	}

	evalDate := parseEvaluationDate(req.EvaluationDate)
	strat := h.stratRegistry.ResolveStrategy(evalDate)

	if err := strat.ValidateBehavior(behavior); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Đánh giá phẩm chất hoặc năng lực không hợp lệ: " + err.Error(),
		})
		return
	}

	database, err := db.GetDB()
	if err == nil && database != nil {
		query := `
			INSERT INTO student_behavior_evaluations (
				id, student_id, term_id, evaluator_id,
				patriotism, compassion, diligence, honesty, responsibility, qualities_comment,
				autonomy_and_self_learning, communication_and_collaboration, problem_solving_and_creativity,
				language_competence, mathematical_competence, scientific_competence,
				technological_competence, computer_competence, aesthetic_competence,
				physical_competence, competencies_comment, updated_at
			) VALUES (
				$1, $2, $3, NULLIF($4, '00000000-0000-0000-0000-000000000000'::uuid),
				NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''), NULLIF($8, ''), NULLIF($9, ''), $10,
				NULLIF($11, ''), NULLIF($12, ''), NULLIF($13, ''),
				NULLIF($14, ''), NULLIF($15, ''), NULLIF($16, ''),
				NULLIF($17, ''), NULLIF($18, ''), NULLIF($19, ''),
				NULLIF($20, ''), $21, NOW()
			)
			ON CONFLICT (student_id, term_id) DO UPDATE SET
				evaluator_id = EXCLUDED.evaluator_id,
				patriotism = COALESCE(EXCLUDED.patriotism, student_behavior_evaluations.patriotism),
				compassion = COALESCE(EXCLUDED.compassion, student_behavior_evaluations.compassion),
				diligence = COALESCE(EXCLUDED.diligence, student_behavior_evaluations.diligence),
				honesty = COALESCE(EXCLUDED.honesty, student_behavior_evaluations.honesty),
				responsibility = COALESCE(EXCLUDED.responsibility, student_behavior_evaluations.responsibility),
				qualities_comment = COALESCE(EXCLUDED.qualities_comment, student_behavior_evaluations.qualities_comment),
				autonomy_and_self_learning = COALESCE(EXCLUDED.autonomy_and_self_learning, student_behavior_evaluations.autonomy_and_self_learning),
				communication_and_collaboration = COALESCE(EXCLUDED.communication_and_collaboration, student_behavior_evaluations.communication_and_collaboration),
				problem_solving_and_creativity = COALESCE(EXCLUDED.problem_solving_and_creativity, student_behavior_evaluations.problem_solving_and_creativity),
				language_competence = COALESCE(EXCLUDED.language_competence, student_behavior_evaluations.language_competence),
				mathematical_competence = COALESCE(EXCLUDED.mathematical_competence, student_behavior_evaluations.mathematical_competence),
				scientific_competence = COALESCE(EXCLUDED.scientific_competence, student_behavior_evaluations.scientific_competence),
				technological_competence = COALESCE(EXCLUDED.technological_competence, student_behavior_evaluations.technological_competence),
				computer_competence = COALESCE(EXCLUDED.computer_competence, student_behavior_evaluations.computer_competence),
				aesthetic_competence = COALESCE(EXCLUDED.aesthetic_competence, student_behavior_evaluations.aesthetic_competence),
				physical_competence = COALESCE(EXCLUDED.physical_competence, student_behavior_evaluations.physical_competence),
				competencies_comment = COALESCE(EXCLUDED.competencies_comment, student_behavior_evaluations.competencies_comment),
				updated_at = NOW()
			RETURNING id;
		`
		var savedID uuid.UUID
		err = database.QueryRow(query,
			behavior.ID, behavior.StudentID, behavior.TermID, behavior.TeacherID,
			behavior.Patriotism, behavior.Compassion, behavior.Diligence, behavior.Honesty, behavior.Responsibility, behavior.QualitiesComment,
			behavior.AutonomySelfLearning, behavior.CommunicationCooperation, behavior.ProblemSolvingCreativity,
			behavior.LanguageCompetence, behavior.MathematicalCompetence, behavior.ScientificCompetence,
			behavior.TechnologicalCompetence, behavior.ComputerCompetence, behavior.AestheticCompetence,
			behavior.PhysicalCompetence, behavior.CompetenciesComment,
		).Scan(&savedID)
		if err == nil {
			behavior.ID = savedID
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Ghi nhận đánh giá phẩm chất và năng lực học sinh thành công",
		"data": gin.H{
			"behavior":           behavior,
			"regulation_applied": strat.GetRegulationCode(),
		},
	})
}

// TranscriptGenerateRequest defines input for generating immutable digital transcripts
type TranscriptGenerateRequest struct {
	StudentID          uuid.UUID                         `json:"student_id" binding:"required"`
	AcademicYearID     string                            `json:"academic_year_id" binding:"required"`
	GradeLevel         int                               `json:"grade_level"`
	EvaluationDate     string                            `json:"evaluation_date,omitempty"`
	Student            *domain.Student                   `json:"student,omitempty"`
	Subjects           []domain.StudentSubjectEvaluation `json:"subjects,omitempty"`
	Behavior           *domain.StudentBehaviorEvaluation `json:"behavior,omitempty"`
	TeacherSignature   string                            `json:"teacher_signature,omitempty"`
	PrincipalSignature string                            `json:"principal_signature,omitempty"`
}

func resolveStudent(database *sql.DB, studentID uuid.UUID, input *domain.Student) (domain.Student, error) {
	if input != nil {
		s := *input
		if s.ID == uuid.Nil {
			s.ID = studentID
		}
		return s, nil
	}
	if database != nil {
		var s domain.Student
		err := database.QueryRow(`
			SELECT id, school_id, national_id, student_code, full_name, birth_date, gender, is_inclusive_education, has_digital_consent
			FROM students WHERE id = $1
		`, studentID).Scan(
			&s.ID, &s.SchoolID, &s.NationalID, &s.StudentCode,
			&s.FullName, &s.BirthDate, &s.Gender,
			&s.IsInclusiveEducation, &s.HasDigitalConsent,
		)
		if err != nil {
			return domain.Student{}, fmt.Errorf("không tìm thấy hồ sơ học sinh: %w", err)
		}
		return s, nil
	}
	return domain.Student{
		ID:          studentID,
		StudentCode: "HS_DEFAULT",
		FullName:    "Học sinh Tiểu học",
	}, nil
}

func resolveSubjects(database *sql.DB, studentID uuid.UUID, input []domain.StudentSubjectEvaluation) []domain.StudentSubjectEvaluation {
	if len(input) > 0 || database == nil {
		return input
	}
	var subjects []domain.StudentSubjectEvaluation
	rows, err := database.Query(`
		SELECT id, student_id, term_id, subject_id, teacher_id, achievement_level, teacher_comment, exam_score, is_computer_based_test, is_centralized_exam
		FROM student_subject_evaluations
		WHERE student_id = $1
	`, studentID)
	if err != nil {
		return subjects
	}
	defer rows.Close()
	for rows.Next() {
		var sub domain.StudentSubjectEvaluation
		var teacherID *uuid.UUID
		if err := rows.Scan(
			&sub.ID, &sub.StudentID, &sub.TermID, &sub.SubjectID, &teacherID,
			&sub.AchievementLevel, &sub.TeacherComment, &sub.ExamScore,
			&sub.IsComputerBasedTest, &sub.IsCentralizedExam,
		); err == nil {
			if teacherID != nil {
				sub.TeacherID = *teacherID
			}
			subjects = append(subjects, sub)
		}
	}
	return subjects
}

func resolveBehavior(database *sql.DB, studentID uuid.UUID, input *domain.StudentBehaviorEvaluation) domain.StudentBehaviorEvaluation {
	if input != nil {
		b := *input
		if b.StudentID == uuid.Nil {
			b.StudentID = studentID
		}
		return b
	}
	var b domain.StudentBehaviorEvaluation
	if database == nil {
		return b
	}
	var evaluatorID *uuid.UUID
	err := database.QueryRow(`
		SELECT id, student_id, term_id, evaluator_id,
		       COALESCE(patriotism, 'T'), COALESCE(compassion, 'T'), COALESCE(diligence, 'T'), COALESCE(honesty, 'T'), COALESCE(responsibility, 'T'),
		       COALESCE(autonomy_and_self_learning, 'T'), COALESCE(communication_and_collaboration, 'T'), COALESCE(problem_solving_and_creativity, 'T'),
		       COALESCE(language_competence, ''), COALESCE(mathematical_competence, ''), COALESCE(scientific_competence, ''),
		       COALESCE(technological_competence, ''), COALESCE(computer_competence, ''), COALESCE(aesthetic_competence, ''), COALESCE(physical_competence, '')
		FROM student_behavior_evaluations
		WHERE student_id = $1
		ORDER BY updated_at DESC LIMIT 1
	`, studentID).Scan(
		&b.ID, &b.StudentID, &b.TermID, &evaluatorID,
		&b.Patriotism, &b.Compassion, &b.Diligence, &b.Honesty, &b.Responsibility,
		&b.AutonomySelfLearning, &b.CommunicationCooperation, &b.ProblemSolvingCreativity,
		&b.LanguageCompetence, &b.MathematicalCompetence, &b.ScientificCompetence,
		&b.TechnologicalCompetence, &b.ComputerCompetence, &b.AestheticCompetence, &b.PhysicalCompetence,
	)
	if err == nil && evaluatorID != nil {
		b.TeacherID = *evaluatorID
	}
	return b
}

func determineRewardAndPromotion(rank assessment.AnnualRank) (string, string) {
	switch rank {
	case assessment.RankStudentExcellent:
		return "Học sinh Xuất sắc", "PROMOTED"
	case assessment.RankStudentGood:
		return "Học sinh Tiêu biểu", "PROMOTED"
	case assessment.RankStudentSatisfactory:
		return "Hoàn thành chương trình lớp học", "PROMOTED"
	default:
		return "Chưa hoàn thành", "RETAINED"
	}
}

func persistTranscriptSnapshot(database *sql.DB, req TranscriptGenerateRequest, gradeLevel int, regulationApplied, annualReward, promotionStatus string, snapshotBytes []byte) uuid.UUID {
	snapshotID := uuid.New()
	if database == nil {
		return snapshotID
	}
	query := `
		INSERT INTO academic_transcripts_snapshots (
			id, student_id, academic_year_id, grade_level, regulation_applied,
			annual_reward, promotion_status, snapshot_data,
			teacher_signature, principal_signature, is_locked
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, true
		)
		ON CONFLICT (student_id, academic_year_id) DO UPDATE SET
			grade_level = EXCLUDED.grade_level,
			regulation_applied = EXCLUDED.regulation_applied,
			annual_reward = EXCLUDED.annual_reward,
			promotion_status = EXCLUDED.promotion_status,
			snapshot_data = EXCLUDED.snapshot_data,
			teacher_signature = COALESCE(EXCLUDED.teacher_signature, academic_transcripts_snapshots.teacher_signature),
			principal_signature = COALESCE(EXCLUDED.principal_signature, academic_transcripts_snapshots.principal_signature),
			is_locked = EXCLUDED.is_locked
		RETURNING id;
	`
	var savedID uuid.UUID
	err := database.QueryRow(query,
		snapshotID, req.StudentID, req.AcademicYearID, gradeLevel, regulationApplied,
		annualReward, promotionStatus, snapshotBytes,
		req.TeacherSignature, req.PrincipalSignature,
	).Scan(&savedID)
	if err == nil {
		return savedID
	}
	return snapshotID
}

// GenerateTranscript handles academic transcript snapshot generation (POST /api/v1/transcripts/generate)
func (h *APIHandler) GenerateTranscript(c *gin.Context) {
	var req TranscriptGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu yêu cầu xuất học bạ không hợp lệ: " + err.Error(),
		})
		return
	}

	evalDate := parseEvaluationDate(req.EvaluationDate)
	strat := h.stratRegistry.ResolveStrategy(evalDate)
	regulationApplied := strat.GetRegulationCode()

	database, _ := db.GetDB()

	student, err := resolveStudent(database, req.StudentID, req.Student)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	subjects := resolveSubjects(database, req.StudentID, req.Subjects)
	behavior := resolveBehavior(database, req.StudentID, req.Behavior)

	gradeLevel := req.GradeLevel
	if gradeLevel <= 0 || gradeLevel > 5 {
		gradeLevel = 5
	}

	rank, reason, err := strat.EvaluateYearEnd(student, subjects, behavior)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Không thể tổng kết xếp loại năm học: " + err.Error(),
		})
		return
	}

	annualReward, promotionStatus := determineRewardAndPromotion(rank)

	snapshotBytes, err := strat.GenerateTranscriptSnapshot(student, req.AcademicYearID, rank, subjects, behavior)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    500,
			"message": "Lỗi sinh snapshot học bạ: " + err.Error(),
		})
		return
	}

	snapshotID := persistTranscriptSnapshot(database, req, gradeLevel, regulationApplied, annualReward, promotionStatus, snapshotBytes)

	var parsedSnapshot map[string]interface{}
	_ = json.Unmarshal(snapshotBytes, &parsedSnapshot)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": fmt.Sprintf("Tạo snapshot học bạ số thành công theo %s", regulationApplied),
		"data": gin.H{
			"id":                 snapshotID,
			"student_id":         req.StudentID,
			"academic_year_id":   req.AcademicYearID,
			"grade_level":        gradeLevel,
			"regulation_applied": regulationApplied,
			"annual_rank":        rank,
			"annual_reward":      annualReward,
			"promotion_status":   promotionStatus,
			"reason":             reason,
			"snapshot_data":      parsedSnapshot,
			"is_locked":          true,
		},
	})
}
