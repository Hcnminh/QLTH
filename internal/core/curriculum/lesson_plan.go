package curriculum

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Standard constants for curriculum lesson plans
const (
	StandardCV2345   = "CV_2345"    // Công văn 2345/BGDĐT-GDTH về Kế hoạch bài dạy
	StandardSTEMTT86 = "STEM_TT_86" // Thông tư 86/BGDĐT & Công văn 909/BGDĐT-GDTH về STEM
)

// Status constants for lesson plan approval workflow
type PlanStatus string

const (
	StatusDraft                  PlanStatus = "DRAFT"
	StatusSubmitted              PlanStatus = "SUBMITTED"
	StatusDepartmentHeadApproved PlanStatus = "DEPARTMENT_HEAD_APPROVED"
	StatusApproved               PlanStatus = "APPROVED"
	StatusRejected               PlanStatus = "REJECTED"
)

// ReviewAction represents the action performed by a reviewer
type ReviewAction string

const (
	ActionSubmit          ReviewAction = "SUBMIT"
	ActionApprove         ReviewAction = "APPROVE"
	ActionReject          ReviewAction = "REJECT"
	ActionRequestRevision ReviewAction = "REQUEST_REVISION"
)

// ReviewerRole constants
const (
	RoleTeacher        = "GV"
	RoleHomeroom       = "GVCN"
	RoleSubjectTeacher = "GVBM"
	RoleDepartmentHead = "TO_TRUONG"
	RolePrincipal      = "BGH"
	RoleSuperAdmin     = "SUPER_ADMIN"
)

var (
	ErrInvalidGradeLevel    = errors.New("khối lớp không hợp lệ (phải từ lớp 1 đến lớp 5)")
	ErrInvalidWeekNumber    = errors.New("tuần học không hợp lệ (phải từ tuần 1 đến tuần 35 theo TT 05/2025)")
	ErrEmptyTitle           = errors.New("tên bài học / kế hoạch bài dạy không được để trống")
	ErrEmptySubject         = errors.New("mã môn học không được để trống")
	ErrInvalidStandard      = errors.New("khung kế hoạch bài dạy phải là CV_2345 hoặc STEM_TT_86")
	ErrAlreadyApproved      = errors.New("kế hoạch bài dạy đã được phê duyệt chính thức, không thể sửa đổi trực tiếp")
	ErrInvalidTransition    = errors.New("chuyển đổi trạng thái không hợp lệ trong quy trình phê duyệt")
	ErrUnauthorizedReviewer = errors.New("người dùng không có thẩm quyền phê duyệt kế hoạch bài dạy")
	ErrMissingReviewComment = errors.New("nhận xét/góp ý của cấp phê duyệt là bắt buộc khi từ chối hoặc yêu cầu chỉnh sửa")
	ErrPlanNotFound         = errors.New("không tìm thấy kế hoạch bài dạy")
)

// LessonObjectives captures knowledge, competencies, and qualities (CV 2345)
type LessonObjectives struct {
	Knowledge    string `json:"knowledge"`    // Yêu cầu cần đạt về Kiến thức
	Competencies string `json:"competencies"` // Yêu cầu cần đạt về Năng lực
	Qualities    string `json:"qualities"`    // Yêu cầu cần đạt về Phẩm chất
}

// LessonActivity captures teaching and learning activities according to CV 2345 4-phase model
type LessonActivity struct {
	Phase           string `json:"phase"`                      // 'Khởi động', 'Khám phá', 'Luyện tập', 'Vận dụng' (hoặc các bước STEM)
	DurationMins    int    `json:"duration_mins,omitempty"`    // Thời lượng (phút)
	Objective       string `json:"objective,omitempty"`        // Mục tiêu hoạt động
	Content         string `json:"content"`                    // Nội dung hoạt động
	ExpectedProduct string `json:"expected_product,omitempty"` // Sản phẩm học tập của học sinh
	Implementation  string `json:"implementation,omitempty"`   // Tổ chức thực hiện (Giao nhiệm vụ, Thực hiện, Báo cáo, Đánh giá)
}

// LessonPlanReview records an approval or feedback event in the workflow
type LessonPlanReview struct {
	ID               uuid.UUID    `json:"id"`
	LessonPlanID     uuid.UUID    `json:"lesson_plan_id"`
	ReviewerID       uuid.UUID    `json:"reviewer_id"`
	ReviewerRole     string       `json:"reviewer_role"` // TO_TRUONG, BGH, SUPER_ADMIN
	ReviewerName     string       `json:"reviewer_name"`
	Action           ReviewAction `json:"action"` // APPROVE, REJECT, REQUEST_REVISION
	Comments         string       `json:"comments"`
	DigitalSignature string       `json:"digital_signature,omitempty"`
	ReviewedAt       time.Time    `json:"reviewed_at"`
}

// LessonPlan represents a pedagogical lesson plan (CV 2345 / STEM TT 86)
type LessonPlan struct {
	ID                 uuid.UUID          `json:"id"`
	SchoolID           uuid.UUID          `json:"school_id"`
	TeacherID          uuid.UUID          `json:"teacher_id"`
	TeacherName        string             `json:"teacher_name,omitempty"`
	ClassID            *uuid.UUID         `json:"class_id,omitempty"`
	GradeLevel         int                `json:"grade_level"`
	SubjectID          string             `json:"subject_id"`
	Title              string             `json:"title"`
	WeekNumber         int                `json:"week_number"`
	DurationSlots      int                `json:"duration_slots"`
	CurriculumStandard string             `json:"curriculum_standard"`
	Objectives         LessonObjectives   `json:"objectives"`
	Equipment          string             `json:"equipment"`
	Activities         []LessonActivity   `json:"activities"`
	AttachmentURL      string             `json:"attachment_url,omitempty"`
	CurrentStatus      PlanStatus         `json:"current_status"`
	Reviews            []LessonPlanReview `json:"reviews,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

// LessonPlanEngine manages validation and multi-level approval workflow
type LessonPlanEngine struct{}

// NewLessonPlanEngine initializes a new engine
func NewLessonPlanEngine() *LessonPlanEngine {
	return &LessonPlanEngine{}
}

// ValidateLessonPlan checks legal and structural requirements for a lesson plan
func (e *LessonPlanEngine) ValidateLessonPlan(plan *LessonPlan) error {
	if plan == nil {
		return errors.New("kế hoạch bài dạy không được nil")
	}
	if plan.GradeLevel < 1 || plan.GradeLevel > 5 {
		return ErrInvalidGradeLevel
	}
	if plan.WeekNumber < 1 || plan.WeekNumber > 35 {
		return ErrInvalidWeekNumber
	}
	if strings.TrimSpace(plan.Title) == "" {
		return ErrEmptyTitle
	}
	if strings.TrimSpace(plan.SubjectID) == "" {
		return ErrEmptySubject
	}
	standard := strings.ToUpper(strings.TrimSpace(plan.CurriculumStandard))
	if standard != "" && standard != StandardCV2345 && standard != StandardSTEMTT86 {
		return ErrInvalidStandard
	}
	if plan.DurationSlots <= 0 {
		plan.DurationSlots = 1
	}
	return nil
}

// SubmitLessonPlan validates and transitions a lesson plan into SUBMITTED status
func (e *LessonPlanEngine) SubmitLessonPlan(plan *LessonPlan) error {
	if err := e.ValidateLessonPlan(plan); err != nil {
		return err
	}

	if plan.ID == uuid.Nil {
		plan.ID = uuid.New()
	}

	standard := strings.ToUpper(strings.TrimSpace(plan.CurriculumStandard))
	if standard == "" {
		plan.CurriculumStandard = StandardCV2345
	} else {
		plan.CurriculumStandard = standard
	}

	now := time.Now()
	if plan.CreatedAt.IsZero() {
		plan.CreatedAt = now
	}
	plan.UpdatedAt = now
	plan.CurrentStatus = StatusSubmitted

	return nil
}

// ReviewLessonPlan processes an approval or rejection step in the multi-level workflow
// Level 1: TO_TRUONG (Department Head) approves -> DEPARTMENT_HEAD_APPROVED
// Level 2: BGH (Principal) approves -> APPROVED
// Any level rejects -> REJECTED (requires comments)
func (e *LessonPlanEngine) ReviewLessonPlan(plan *LessonPlan, review LessonPlanReview) error {
	if plan == nil {
		return ErrPlanNotFound
	}

	// Cannot modify already approved plans without re-submission
	if plan.CurrentStatus == StatusApproved {
		return ErrAlreadyApproved
	}

	role := strings.ToUpper(strings.TrimSpace(review.ReviewerRole))
	action := ReviewAction(strings.ToUpper(strings.TrimSpace(string(review.Action))))

	// Authorization check
	switch role {
	case RoleDepartmentHead, RolePrincipal, RoleSuperAdmin:
		// authorized reviewers
	default:
		return fmt.Errorf("%w: vai trò '%s' không được phép duyệt giáo án", ErrUnauthorizedReviewer, role)
	}

	// Action validation
	switch action {
	case ActionApprove:
		if role == RoleDepartmentHead {
			// Tổ trưởng duyệt bước 1
			plan.CurrentStatus = StatusDepartmentHeadApproved
		} else {
			// BGH hoặc SUPER_ADMIN duyệt hoàn tất
			plan.CurrentStatus = StatusApproved
		}
	case ActionReject, ActionRequestRevision:
		if strings.TrimSpace(review.Comments) == "" {
			return ErrMissingReviewComment
		}
		plan.CurrentStatus = StatusRejected
	default:
		return fmt.Errorf("%w: hành động duyệt '%s' không được hỗ trợ", ErrInvalidTransition, action)
	}

	now := time.Now()
	review.ID = uuid.New()
	review.LessonPlanID = plan.ID
	review.Action = action
	review.ReviewerRole = role
	if review.ReviewedAt.IsZero() {
		review.ReviewedAt = now
	}

	plan.Reviews = append(plan.Reviews, review)
	plan.UpdatedAt = now

	return nil
}
