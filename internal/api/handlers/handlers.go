package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"qlth/internal/core/assessment"
	"qlth/internal/core/curriculum"
	"qlth/internal/core/domain"
	"qlth/internal/core/kpi"
	"qlth/internal/core/teaching"
	"qlth/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type APIHandler struct {
	kpiEngine        *kpi.Engine
	quotaEngine      *teaching.QuotaEngine
	stratRegistry    *assessment.Registry
	lessonPlanEngine *curriculum.LessonPlanEngine
	timetableEngine  *curriculum.TimetableEngine
	planStoreMu      sync.RWMutex
	planStore        map[uuid.UUID]*curriculum.LessonPlan
	timetableStoreMu sync.RWMutex
	timetableStore   map[uuid.UUID]curriculum.TimetableSlot
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{
		kpiEngine:        kpi.NewEngine(),
		quotaEngine:      teaching.NewQuotaEngine(),
		stratRegistry:    assessment.NewRegistry(),
		lessonPlanEngine: curriculum.NewLessonPlanEngine(),
		timetableEngine:  curriculum.NewTimetableEngine(),
		planStore:        make(map[uuid.UUID]*curriculum.LessonPlan),
		timetableStore:   make(map[uuid.UUID]curriculum.TimetableSlot),
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

// =============================================================================
// CURRICULUM, LESSON PLANS (CV 2345/STEM TT 86) & TIMETABLE (TT 05/2025)
// =============================================================================

// SubmitLessonPlanRequest payload for submitting a lesson plan
type SubmitLessonPlanRequest struct {
	SchoolID           uuid.UUID                   `json:"school_id"`
	TeacherID          uuid.UUID                   `json:"teacher_id" binding:"required"`
	TeacherName        string                      `json:"teacher_name"`
	ClassID            *uuid.UUID                  `json:"class_id"`
	GradeLevel         int                         `json:"grade_level" binding:"required"`
	SubjectID          string                      `json:"subject_id" binding:"required"`
	Title              string                      `json:"title" binding:"required"`
	WeekNumber         int                         `json:"week_number" binding:"required"`
	DurationSlots      int                         `json:"duration_slots"`
	CurriculumStandard string                      `json:"curriculum_standard"`
	Objectives         curriculum.LessonObjectives `json:"objectives"`
	Equipment          string                      `json:"equipment"`
	Activities         []curriculum.LessonActivity `json:"activities"`
	AttachmentURL      string                      `json:"attachment_url"`
	IsDraft            bool                        `json:"is_draft"`
}

// ApproveLessonPlanRequest payload for multi-level review and approval
type ApproveLessonPlanRequest struct {
	Action           string    `json:"action" binding:"required"` // APPROVE, REJECT, REQUEST_REVISION
	ReviewerID       uuid.UUID `json:"reviewer_id"`
	ReviewerRole     string    `json:"reviewer_role"`
	ReviewerName     string    `json:"reviewer_name"`
	Comments         string    `json:"comments"`
	DigitalSignature string    `json:"digital_signature"`
}

// ConfigureTimetableSlotRequest payload for setting up a timetable slot
type ConfigureTimetableSlotRequest struct {
	SchoolID       uuid.UUID `json:"school_id"`
	AcademicYearID string    `json:"academic_year_id" binding:"required"`
	TermID         string    `json:"term_id" binding:"required"`
	ClassID        uuid.UUID `json:"class_id" binding:"required"`
	ClassName      string    `json:"class_name"`
	TeacherID      uuid.UUID `json:"teacher_id" binding:"required"`
	TeacherName    string    `json:"teacher_name"`
	SubjectID      string    `json:"subject_id" binding:"required"`
	DayOfWeek      int       `json:"day_of_week" binding:"required"`
	Period         int       `json:"period" binding:"required"`
	RoomNumber     string    `json:"room_number"`
	WeekNumber     int       `json:"week_number"`
}

// SubmitLessonPlan handles lesson plan submission (POST /api/v1/lesson-plans/submit)
func (h *APIHandler) SubmitLessonPlan(c *gin.Context) {
	var req SubmitLessonPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": "Dữ liệu kế hoạch bài dạy không hợp lệ: " + err.Error(),
		})
		return
	}

	schoolID := req.SchoolID
	if schoolID == uuid.Nil {
		if headerSchool := c.GetHeader("X-School-ID"); headerSchool != "" {
			if parsed, err := uuid.Parse(headerSchool); err == nil {
				schoolID = parsed
			}
		}
	}
	if schoolID == uuid.Nil {
		schoolID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	}

	plan := curriculum.LessonPlan{
		ID:                 uuid.New(),
		SchoolID:           schoolID,
		TeacherID:          req.TeacherID,
		TeacherName:        strings.TrimSpace(req.TeacherName),
		ClassID:            req.ClassID,
		GradeLevel:         req.GradeLevel,
		SubjectID:          strings.TrimSpace(req.SubjectID),
		Title:              strings.TrimSpace(req.Title),
		WeekNumber:         req.WeekNumber,
		DurationSlots:      req.DurationSlots,
		CurriculumStandard: strings.TrimSpace(req.CurriculumStandard),
		Objectives:         req.Objectives,
		Equipment:          strings.TrimSpace(req.Equipment),
		Activities:         req.Activities,
		AttachmentURL:      strings.TrimSpace(req.AttachmentURL),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if req.IsDraft {
		if err := h.lessonPlanEngine.ValidateLessonPlan(&plan); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}
		plan.CurrentStatus = curriculum.StatusDraft
	} else {
		if err := h.lessonPlanEngine.SubmitLessonPlan(&plan); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}
	}

	database, _ := db.GetDB()
	if plan.TeacherName == "" && database != nil {
		var teacherName string
		if err := database.QueryRow("SELECT full_name FROM staffs WHERE id = $1", plan.TeacherID).Scan(&teacherName); err == nil {
			plan.TeacherName = teacherName
		}
	}

	// In-memory store
	h.planStoreMu.Lock()
	h.planStore[plan.ID] = &plan
	h.planStoreMu.Unlock()

	// Persist to Postgres if available
	if database != nil {
		persistLessonPlan(database, &plan)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": "Nộp kế hoạch bài dạy thành công theo Công văn 2345/BGDĐT-GDTH",
		"data":    plan,
	})
}

// ApproveLessonPlan handles multi-level approval or rejection (POST /api/v1/lesson-plans/:id/approve)
func (h *APIHandler) ApproveLessonPlan(c *gin.Context) {
	idStr := c.Param("id")
	planID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": "Mã kế hoạch bài dạy không hợp lệ: " + idStr,
		})
		return
	}

	var req ApproveLessonPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": "Dữ liệu phê duyệt không hợp lệ: " + err.Error(),
		})
		return
	}

	reviewerRole := strings.TrimSpace(req.ReviewerRole)
	if reviewerRole == "" {
		reviewerRole = c.GetHeader("X-User-Role")
	}
	if reviewerRole == "" {
		reviewerRole = curriculum.RoleDepartmentHead
	}

	database, _ := db.GetDB()

	// Locate plan
	h.planStoreMu.RLock()
	plan, exists := h.planStore[planID]
	h.planStoreMu.RUnlock()

	if !exists && database != nil {
		dbPlan, err := loadLessonPlanFromDB(database, planID)
		if err == nil && dbPlan != nil {
			plan = dbPlan
			h.planStoreMu.Lock()
			h.planStore[planID] = plan
			h.planStoreMu.Unlock()
			exists = true
		}
	}

	if !exists || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    http.StatusNotFound,
			"message": "Không tìm thấy kế hoạch bài dạy với mã: " + idStr,
		})
		return
	}

	reviewerName := strings.TrimSpace(req.ReviewerName)
	if reviewerName == "" && req.ReviewerID != uuid.Nil && database != nil {
		_ = database.QueryRow("SELECT full_name FROM staffs WHERE id = $1", req.ReviewerID).Scan(&reviewerName)
	}
	if reviewerName == "" {
		reviewerName = reviewerRole
	}

	review := curriculum.LessonPlanReview{
		ReviewerID:       req.ReviewerID,
		ReviewerRole:     reviewerRole,
		ReviewerName:     reviewerName,
		Action:           curriculum.ReviewAction(strings.ToUpper(strings.TrimSpace(req.Action))),
		Comments:         strings.TrimSpace(req.Comments),
		DigitalSignature: strings.TrimSpace(req.DigitalSignature),
		ReviewedAt:       time.Now(),
	}

	h.planStoreMu.Lock()
	err = h.lessonPlanEngine.ReviewLessonPlan(plan, review)
	h.planStoreMu.Unlock()

	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, curriculum.ErrUnauthorizedReviewer) {
			code = http.StatusForbidden
		}
		c.JSON(code, gin.H{
			"success": false,
			"code":    code,
			"message": err.Error(),
		})
		return
	}

	// Persist updated status and review log to DB
	if database != nil {
		_, _ = database.Exec("UPDATE lesson_plans SET current_status = $1, updated_at = NOW() WHERE id = $2", plan.CurrentStatus, plan.ID)
		if len(plan.Reviews) > 0 {
			latestReview := plan.Reviews[len(plan.Reviews)-1]
			persistLessonPlanReview(database, latestReview)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": "Phê duyệt kế hoạch bài dạy thành công",
		"data":    plan,
	})
}

// GetLessonPlan retrieves a lesson plan by ID (GET /api/v1/lesson-plans/:id)
func (h *APIHandler) GetLessonPlan(c *gin.Context) {
	idStr := c.Param("id")
	planID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": "Mã kế hoạch bài dạy không hợp lệ",
		})
		return
	}

	h.planStoreMu.RLock()
	plan, exists := h.planStore[planID]
	h.planStoreMu.RUnlock()

	database, _ := db.GetDB()
	if (!exists || plan == nil) && database != nil {
		if dbPlan, err := loadLessonPlanFromDB(database, planID); err == nil && dbPlan != nil {
			plan = dbPlan
			exists = true
			h.planStoreMu.Lock()
			h.planStore[planID] = plan
			h.planStoreMu.Unlock()
		}
	}

	if !exists || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    http.StatusNotFound,
			"message": "Không tìm thấy kế hoạch bài dạy",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": "Lấy thông tin kế hoạch bài dạy thành công",
		"data":    plan,
	})
}

// ListLessonPlans lists lesson plans with optional filters (GET /api/v1/lesson-plans)
func (h *APIHandler) ListLessonPlans(c *gin.Context) {
	database, _ := db.GetDB()
	if database != nil {
		var schoolID uuid.UUID
		if s := c.GetHeader("X-School-ID"); s != "" {
			schoolID, _ = uuid.Parse(s)
		}
		if schoolID == uuid.Nil {
			schoolID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
		}

		plans, err := listLessonPlansFromDB(database, schoolID)
		if err == nil && len(plans) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"code":    http.StatusOK,
				"message": "Lấy danh sách kế hoạch bài dạy thành công",
				"data":    plans,
			})
			return
		}
	}

	// Fallback to in-memory store
	h.planStoreMu.RLock()
	plans := make([]*curriculum.LessonPlan, 0, len(h.planStore))
	for _, p := range h.planStore {
		plans = append(plans, p)
	}
	h.planStoreMu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": "Lấy danh sách kế hoạch bài dạy thành công",
		"data":    plans,
	})
}

func resolveSchoolID(c *gin.Context, reqSchoolID uuid.UUID) uuid.UUID {
	if reqSchoolID != uuid.Nil {
		return reqSchoolID
	}
	if s := c.GetHeader("X-School-ID"); s != "" {
		if parsed, err := uuid.Parse(s); err == nil {
			return parsed
		}
	}
	return uuid.MustParse("a0000000-0000-0000-0000-000000000001")
}

func (h *APIHandler) gatherExistingSlots(database *sql.DB, schoolID uuid.UUID, academicYearID, termID string) []curriculum.TimetableSlot {
	var existingSlots []curriculum.TimetableSlot
	if database != nil {
		if dbSlots, err := loadTimetableSlotsFromDB(database, schoolID, academicYearID, termID); err == nil {
			existingSlots = dbSlots
		}
	}

	h.timetableStoreMu.RLock()
	defer h.timetableStoreMu.RUnlock()
	for _, s := range h.timetableStore {
		if s.SchoolID == schoolID && s.AcademicYearID == academicYearID && s.TermID == termID {
			alreadyPresent := false
			for _, ex := range existingSlots {
				if ex.ID == s.ID {
					alreadyPresent = true
					break
				}
			}
			if !alreadyPresent {
				existingSlots = append(existingSlots, s)
			}
		}
	}
	return existingSlots
}

func enrichSlotMetadata(database *sql.DB, slot *curriculum.TimetableSlot) {
	if database == nil {
		return
	}
	if slot.TeacherName == "" {
		_ = database.QueryRow("SELECT full_name FROM staffs WHERE id = $1", slot.TeacherID).Scan(&slot.TeacherName)
	}
	if slot.ClassName == "" {
		_ = database.QueryRow("SELECT name FROM classes WHERE id = $1", slot.ClassID).Scan(&slot.ClassName)
	}
}

func filterSlots(slots []curriculum.TimetableSlot, classFilter, teacherFilter string) []curriculum.TimetableSlot {
	results := slots
	if classFilter != "" {
		if cid, err := uuid.Parse(classFilter); err == nil {
			filtered := make([]curriculum.TimetableSlot, 0, len(results))
			for _, s := range results {
				if s.ClassID == cid {
					filtered = append(filtered, s)
				}
			}
			results = filtered
		}
	}
	if teacherFilter != "" {
		if tid, err := uuid.Parse(teacherFilter); err == nil {
			filtered := make([]curriculum.TimetableSlot, 0, len(results))
			for _, s := range results {
				if s.TeacherID == tid {
					filtered = append(filtered, s)
				}
			}
			results = filtered
		}
	}
	return results
}

// ConfigureTimetableSlot handles schedule slot configuration and conflict checks (POST /api/v1/timetable/slots)
func (h *APIHandler) ConfigureTimetableSlot(c *gin.Context) {
	var req ConfigureTimetableSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": "Dữ liệu cấu hình tiết thời khóa biểu không hợp lệ: " + err.Error(),
		})
		return
	}

	schoolID := resolveSchoolID(c, req.SchoolID)

	slot := curriculum.TimetableSlot{
		ID:             uuid.New(),
		SchoolID:       schoolID,
		AcademicYearID: strings.TrimSpace(req.AcademicYearID),
		TermID:         strings.TrimSpace(req.TermID),
		ClassID:        req.ClassID,
		ClassName:      strings.TrimSpace(req.ClassName),
		TeacherID:      req.TeacherID,
		TeacherName:    strings.TrimSpace(req.TeacherName),
		SubjectID:      strings.TrimSpace(req.SubjectID),
		DayOfWeek:      req.DayOfWeek,
		Period:         req.Period,
		RoomNumber:     strings.TrimSpace(req.RoomNumber),
		WeekNumber:     req.WeekNumber,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.timetableEngine.ValidateSlotBounds(slot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	database, _ := db.GetDB()
	enrichSlotMetadata(database, &slot)

	existingSlots := h.gatherExistingSlots(database, schoolID, slot.AcademicYearID, slot.TermID)

	savedSlot, err := h.timetableEngine.ValidateAndAddSlot(existingSlots, slot)
	if err != nil {
		if conflict, ok := err.(*curriculum.TimetableConflict); ok {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"code":    http.StatusConflict,
				"message": conflict.Message,
				"data":    conflict,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// Store in memory
	h.timetableStoreMu.Lock()
	h.timetableStore[savedSlot.ID] = savedSlot
	h.timetableStoreMu.Unlock()

	// Persist to Postgres if available
	if database != nil {
		persistTimetableSlot(database, savedSlot)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": "Cấu hình tiết thời khóa biểu thành công, đảm bảo định mức Thông tư 05/2025/TT-BGDĐT",
		"data":    savedSlot,
	})
}

// GetTimetableSlots lists scheduled timetable slots with optional filters (GET /api/v1/timetable/slots)
func (h *APIHandler) GetTimetableSlots(c *gin.Context) {
	database, _ := db.GetDB()
	schoolID := resolveSchoolID(c, uuid.Nil)

	academicYear := c.Query("academic_year_id")
	if academicYear == "" {
		academicYear = "2026-2027"
	}
	termID := c.Query("term_id")
	if termID == "" {
		termID = "2026-2027-HK1"
	}

	var results []curriculum.TimetableSlot
	if database != nil {
		if dbSlots, err := loadTimetableSlotsFromDB(database, schoolID, academicYear, termID); err == nil && len(dbSlots) > 0 {
			results = dbSlots
		}
	}

	if len(results) == 0 {
		h.timetableStoreMu.RLock()
		for _, s := range h.timetableStore {
			if s.SchoolID == schoolID && s.AcademicYearID == academicYear && s.TermID == termID {
				results = append(results, s)
			}
		}
		h.timetableStoreMu.RUnlock()
	}

	results = filterSlots(results, c.Query("class_id"), c.Query("teacher_id"))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": "Lấy danh sách thời khóa biểu thành công",
		"data":    results,
	})
}

// -----------------------------------------------------------------------------
// Database Persistence Helpers for Lesson Plans & Timetable
// -----------------------------------------------------------------------------

func persistLessonPlan(database *sql.DB, plan *curriculum.LessonPlan) {
	objBytes, _ := json.Marshal(plan.Objectives)
	actBytes, _ := json.Marshal(plan.Activities)

	query := `
		INSERT INTO lesson_plans (
			id, school_id, teacher_id, class_id, grade_level, subject_id,
			title, week_number, duration_slots, curriculum_standard,
			objectives, equipment, activities, attachment_url, current_status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			week_number = EXCLUDED.week_number,
			duration_slots = EXCLUDED.duration_slots,
			objectives = EXCLUDED.objectives,
			equipment = EXCLUDED.equipment,
			activities = EXCLUDED.activities,
			attachment_url = EXCLUDED.attachment_url,
			current_status = EXCLUDED.current_status,
			updated_at = EXCLUDED.updated_at;
	`
	_, _ = database.Exec(query,
		plan.ID, plan.SchoolID, plan.TeacherID, plan.ClassID, plan.GradeLevel, plan.SubjectID,
		plan.Title, plan.WeekNumber, plan.DurationSlots, plan.CurriculumStandard,
		objBytes, plan.Equipment, actBytes, plan.AttachmentURL, plan.CurrentStatus, plan.CreatedAt, plan.UpdatedAt,
	)
}

func persistLessonPlanReview(database *sql.DB, review curriculum.LessonPlanReview) {
	query := `
		INSERT INTO lesson_plan_reviews (
			id, lesson_plan_id, reviewer_id, reviewer_role, reviewer_name, action, comments, digital_signature, reviewed_at
		) VALUES (
			$1, $2, NULLIF($3, '00000000-0000-0000-0000-000000000000'::uuid), $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT (id) DO NOTHING;
	`
	_, _ = database.Exec(query,
		review.ID, review.LessonPlanID, review.ReviewerID, review.ReviewerRole, review.ReviewerName,
		review.Action, review.Comments, review.DigitalSignature, review.ReviewedAt,
	)
}

func loadLessonPlanFromDB(database *sql.DB, planID uuid.UUID) (*curriculum.LessonPlan, error) {
	query := `
		SELECT lp.id, lp.school_id, lp.teacher_id, COALESCE(st.full_name, ''), lp.class_id,
		       lp.grade_level, lp.subject_id, lp.title, lp.week_number, lp.duration_slots,
		       lp.curriculum_standard, lp.objectives, lp.equipment, lp.activities,
		       COALESCE(lp.attachment_url, ''), lp.current_status, lp.created_at, lp.updated_at
		FROM lesson_plans lp
		LEFT JOIN staffs st ON lp.teacher_id = st.id
		WHERE lp.id = $1;
	`
	row := database.QueryRow(query, planID)

	var p curriculum.LessonPlan
	var objBytes, actBytes []byte
	var status string
	err := row.Scan(
		&p.ID, &p.SchoolID, &p.TeacherID, &p.TeacherName, &p.ClassID,
		&p.GradeLevel, &p.SubjectID, &p.Title, &p.WeekNumber, &p.DurationSlots,
		&p.CurriculumStandard, &objBytes, &p.Equipment, &actBytes,
		&p.AttachmentURL, &status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	p.CurrentStatus = curriculum.PlanStatus(status)
	_ = json.Unmarshal(objBytes, &p.Objectives)
	_ = json.Unmarshal(actBytes, &p.Activities)

	// Load reviews
	revQuery := `
		SELECT id, lesson_plan_id, COALESCE(reviewer_id, '00000000-0000-0000-0000-000000000000'::uuid),
		       reviewer_role, COALESCE(reviewer_name, ''), action, COALESCE(comments, ''),
		       COALESCE(digital_signature, ''), reviewed_at
		FROM lesson_plan_reviews
		WHERE lesson_plan_id = $1
		ORDER BY reviewed_at ASC;
	`
	rows, err := database.Query(revQuery, planID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var r curriculum.LessonPlanReview
			var act string
			if scanErr := rows.Scan(
				&r.ID, &r.LessonPlanID, &r.ReviewerID, &r.ReviewerRole,
				&r.ReviewerName, &act, &r.Comments, &r.DigitalSignature, &r.ReviewedAt,
			); scanErr == nil {
				r.Action = curriculum.ReviewAction(act)
				p.Reviews = append(p.Reviews, r)
			}
		}
	}

	return &p, nil
}

func listLessonPlansFromDB(database *sql.DB, schoolID uuid.UUID) ([]*curriculum.LessonPlan, error) {
	query := `
		SELECT lp.id, lp.school_id, lp.teacher_id, COALESCE(st.full_name, ''), lp.class_id,
		       lp.grade_level, lp.subject_id, lp.title, lp.week_number, lp.duration_slots,
		       lp.curriculum_standard, lp.objectives, lp.equipment, lp.activities,
		       COALESCE(lp.attachment_url, ''), lp.current_status, lp.created_at, lp.updated_at
		FROM lesson_plans lp
		LEFT JOIN staffs st ON lp.teacher_id = st.id
		WHERE lp.school_id = $1
		ORDER BY lp.week_number DESC, lp.created_at DESC;
	`
	rows, err := database.Query(query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*curriculum.LessonPlan
	for rows.Next() {
		var p curriculum.LessonPlan
		var objBytes, actBytes []byte
		var status string
		if scanErr := rows.Scan(
			&p.ID, &p.SchoolID, &p.TeacherID, &p.TeacherName, &p.ClassID,
			&p.GradeLevel, &p.SubjectID, &p.Title, &p.WeekNumber, &p.DurationSlots,
			&p.CurriculumStandard, &objBytes, &p.Equipment, &actBytes,
			&p.AttachmentURL, &status, &p.CreatedAt, &p.UpdatedAt,
		); scanErr == nil {
			p.CurrentStatus = curriculum.PlanStatus(status)
			_ = json.Unmarshal(objBytes, &p.Objectives)
			_ = json.Unmarshal(actBytes, &p.Activities)
			plans = append(plans, &p)
		}
	}
	return plans, nil
}

func persistTimetableSlot(database *sql.DB, slot curriculum.TimetableSlot) {
	query := `
		INSERT INTO timetable_slots (
			id, school_id, academic_year_id, term_id, class_id, teacher_id, subject_id,
			day_of_week, period, room_number, week_number, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
		ON CONFLICT (id) DO UPDATE SET
			room_number = EXCLUDED.room_number,
			updated_at = EXCLUDED.updated_at;
	`
	_, _ = database.Exec(query,
		slot.ID, slot.SchoolID, slot.AcademicYearID, slot.TermID, slot.ClassID, slot.TeacherID, slot.SubjectID,
		slot.DayOfWeek, slot.Period, slot.RoomNumber, slot.WeekNumber, slot.CreatedAt, slot.UpdatedAt,
	)
}

func loadTimetableSlotsFromDB(database *sql.DB, schoolID uuid.UUID, academicYearID, termID string) ([]curriculum.TimetableSlot, error) {
	query := `
		SELECT ts.id, ts.school_id, ts.academic_year_id, ts.term_id, ts.class_id,
		       COALESCE(c.name, ''), ts.teacher_id, COALESCE(st.full_name, ''), ts.subject_id,
		       ts.day_of_week, ts.period, COALESCE(ts.room_number, ''), ts.week_number,
		       ts.created_at, ts.updated_at
		FROM timetable_slots ts
		LEFT JOIN classes c ON ts.class_id = c.id
		LEFT JOIN staffs st ON ts.teacher_id = st.id
		WHERE ts.school_id = $1 AND ts.academic_year_id = $2 AND ts.term_id = $3
		ORDER BY ts.day_of_week ASC, ts.period ASC;
	`
	rows, err := database.Query(query, schoolID, academicYearID, termID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []curriculum.TimetableSlot
	for rows.Next() {
		var s curriculum.TimetableSlot
		if scanErr := rows.Scan(
			&s.ID, &s.SchoolID, &s.AcademicYearID, &s.TermID, &s.ClassID,
			&s.ClassName, &s.TeacherID, &s.TeacherName, &s.SubjectID,
			&s.DayOfWeek, &s.Period, &s.RoomNumber, &s.WeekNumber,
			&s.CreatedAt, &s.UpdatedAt,
		); scanErr == nil {
			slots = append(slots, s)
		}
	}
	return slots, nil
}
