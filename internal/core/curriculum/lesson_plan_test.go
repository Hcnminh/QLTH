package curriculum_test

import (
	"testing"

	"qlth/internal/core/curriculum"

	"github.com/google/uuid"
)

func TestLessonPlanEngine_Submit_SuccessCV2345(t *testing.T) {
	engine := curriculum.NewLessonPlanEngine()
	plan := &curriculum.LessonPlan{
		SchoolID:           uuid.New(),
		TeacherID:          uuid.New(),
		GradeLevel:         5,
		SubjectID:          "TOAN",
		Title:              "Bài 15: Thể tích hình lập phương",
		WeekNumber:         15,
		DurationSlots:      2,
		CurriculumStandard: curriculum.StandardCV2345,
		Objectives: curriculum.LessonObjectives{
			Knowledge:    "Nắm công thức tính thể tích hình lập phương V = a x a x a",
			Competencies: "Năng lực mô hình hóa toán học",
			Qualities:    "Chăm chỉ, kiên trì",
		},
		Equipment: "Mô hình lập phương đơn vị 1cm3",
		Activities: []curriculum.LessonActivity{
			{Phase: "Khởi động", DurationMins: 5, Content: "Trò chơi xếp khối rubik"},
			{Phase: "Khám phá", DurationMins: 15, Content: "Xây dựng công thức"},
			{Phase: "Luyện tập", DurationMins: 15, Content: "Làm bài tập 1, 2 trang 120"},
			{Phase: "Vận dụng", DurationMins: 5, Content: "Đo thể tích hộp phấn"},
		},
	}

	err := engine.SubmitLessonPlan(plan)
	if err != nil {
		t.Fatalf("expected submit to succeed, got %v", err)
	}

	if plan.CurrentStatus != curriculum.StatusSubmitted {
		t.Errorf("expected status SUBMITTED, got %s", plan.CurrentStatus)
	}
	if plan.ID == uuid.Nil {
		t.Errorf("expected plan ID to be populated")
	}
	if plan.CreatedAt.IsZero() || plan.UpdatedAt.IsZero() {
		t.Errorf("expected timestamps to be set")
	}
}

func TestLessonPlanEngine_Submit_SuccessSTEMTT86(t *testing.T) {
	engine := curriculum.NewLessonPlanEngine()
	plan := &curriculum.LessonPlan{
		SchoolID:           uuid.New(),
		TeacherID:          uuid.New(),
		GradeLevel:         4,
		SubjectID:          "STEM",
		Title:              "Chế tạo bình lọc nước mini",
		WeekNumber:         10,
		CurriculumStandard: "stem_tt_86",
	}

	err := engine.SubmitLessonPlan(plan)
	if err != nil {
		t.Fatalf("expected STEM submit to succeed, got %v", err)
	}

	if plan.CurriculumStandard != curriculum.StandardSTEMTT86 {
		t.Errorf("expected curriculum standard normalized to STEM_TT_86, got %s", plan.CurriculumStandard)
	}
	if plan.DurationSlots != 1 {
		t.Errorf("expected default duration slot 1, got %d", plan.DurationSlots)
	}
}

func TestLessonPlanEngine_Submit_ValidationErrors(t *testing.T) {
	engine := curriculum.NewLessonPlanEngine()

	// 1. Invalid grade level
	p1 := &curriculum.LessonPlan{
		GradeLevel: 6,
		WeekNumber: 10,
		Title:      "Test",
		SubjectID:  "TOAN",
	}
	if err := engine.SubmitLessonPlan(p1); err != curriculum.ErrInvalidGradeLevel {
		t.Errorf("expected ErrInvalidGradeLevel, got %v", err)
	}

	// 2. Invalid week number
	p2 := &curriculum.LessonPlan{
		GradeLevel: 3,
		WeekNumber: 36, // Max is 35 weeks per Circular 05/2025
		Title:      "Test",
		SubjectID:  "TOAN",
	}
	if err := engine.SubmitLessonPlan(p2); err != curriculum.ErrInvalidWeekNumber {
		t.Errorf("expected ErrInvalidWeekNumber, got %v", err)
	}

	// 3. Empty title
	p3 := &curriculum.LessonPlan{
		GradeLevel: 3,
		WeekNumber: 5,
		Title:      "   ",
		SubjectID:  "TOAN",
	}
	if err := engine.SubmitLessonPlan(p3); err != curriculum.ErrEmptyTitle {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}

	// 4. Empty subject
	p4 := &curriculum.LessonPlan{
		GradeLevel: 3,
		WeekNumber: 5,
		Title:      "Học hát",
		SubjectID:  "",
	}
	if err := engine.SubmitLessonPlan(p4); err != curriculum.ErrEmptySubject {
		t.Errorf("expected ErrEmptySubject, got %v", err)
	}

	// 5. Invalid standard
	p5 := &curriculum.LessonPlan{
		GradeLevel:         3,
		WeekNumber:         5,
		Title:              "Học hát",
		SubjectID:          "AM_NHAC",
		CurriculumStandard: "UNKNOWN_STANDARD",
	}
	if err := engine.SubmitLessonPlan(p5); err != curriculum.ErrInvalidStandard {
		t.Errorf("expected ErrInvalidStandard, got %v", err)
	}
}

func TestLessonPlanEngine_MultiLevelApprovalWorkflow(t *testing.T) {
	engine := curriculum.NewLessonPlanEngine()
	plan := &curriculum.LessonPlan{
		ID:                 uuid.New(),
		SchoolID:           uuid.New(),
		TeacherID:          uuid.New(),
		TeacherName:        "Nguyễn Văn An",
		GradeLevel:         5,
		SubjectID:          "TOAN",
		Title:              "Bài 20: Diện tích xung quanh hình lập phương",
		WeekNumber:         16,
		DurationSlots:      2,
		CurriculumStandard: curriculum.StandardCV2345,
	}

	if err := engine.SubmitLessonPlan(plan); err != nil {
		t.Fatalf("submit failed: %v", err)
	}
	if plan.CurrentStatus != curriculum.StatusSubmitted {
		t.Fatalf("expected SUBMITTED, got %s", plan.CurrentStatus)
	}

	// Step 1: Department Head (TO_TRUONG) approves
	totruongReview := curriculum.LessonPlanReview{
		ReviewerID:   uuid.New(),
		ReviewerRole: curriculum.RoleDepartmentHead,
		ReviewerName: "Trần Thị Bình",
		Action:       curriculum.ActionApprove,
		Comments:     "Giáo án đúng chuẩn CV 2345, các hoạt động logic. Duyệt cấp tổ.",
	}
	if err := engine.ReviewLessonPlan(plan, totruongReview); err != nil {
		t.Fatalf("department head review failed: %v", err)
	}
	if plan.CurrentStatus != curriculum.StatusDepartmentHeadApproved {
		t.Errorf("expected DEPARTMENT_HEAD_APPROVED, got %s", plan.CurrentStatus)
	}
	if len(plan.Reviews) != 1 {
		t.Fatalf("expected 1 review record, got %d", len(plan.Reviews))
	}

	// Step 2: Principal (BGH) final approval
	bghReview := curriculum.LessonPlanReview{
		ReviewerID:       uuid.New(),
		ReviewerRole:     curriculum.RolePrincipal,
		ReviewerName:     "Lê Hoàng Cường",
		Action:           curriculum.ActionApprove,
		Comments:         "Ký duyệt chính thức đưa vào giảng dạy tuần 16.",
		DigitalSignature: "SMARTCA_BGH_SIGN_TSA_OK",
	}
	if err := engine.ReviewLessonPlan(plan, bghReview); err != nil {
		t.Fatalf("principal review failed: %v", err)
	}
	if plan.CurrentStatus != curriculum.StatusApproved {
		t.Errorf("expected APPROVED, got %s", plan.CurrentStatus)
	}
	if len(plan.Reviews) != 2 {
		t.Fatalf("expected 2 review records, got %d", len(plan.Reviews))
	}

	// Step 3: Attempting to review an already APPROVED plan fails
	attemptAfterApproved := curriculum.LessonPlanReview{
		ReviewerID:   uuid.New(),
		ReviewerRole: curriculum.RolePrincipal,
		Action:       curriculum.ActionReject,
		Comments:     "Thay đổi ý kiến",
	}
	if err := engine.ReviewLessonPlan(plan, attemptAfterApproved); err != curriculum.ErrAlreadyApproved {
		t.Errorf("expected ErrAlreadyApproved, got %v", err)
	}
}

func TestLessonPlanEngine_RejectionWorkflow(t *testing.T) {
	engine := curriculum.NewLessonPlanEngine()
	plan := &curriculum.LessonPlan{
		GradeLevel: 2,
		WeekNumber: 5,
		Title:      "Bài 5: Phép cộng có nhớ",
		SubjectID:  "TOAN",
	}
	_ = engine.SubmitLessonPlan(plan)

	// Rejection without comment should fail
	err := engine.ReviewLessonPlan(plan, curriculum.LessonPlanReview{
		ReviewerID:   uuid.New(),
		ReviewerRole: curriculum.RoleDepartmentHead,
		Action:       curriculum.ActionReject,
		Comments:     "   ",
	})
	if err != curriculum.ErrMissingReviewComment {
		t.Errorf("expected ErrMissingReviewComment, got %v", err)
	}

	// Rejection with comment transitions to REJECTED
	err = engine.ReviewLessonPlan(plan, curriculum.LessonPlanReview{
		ReviewerID:   uuid.New(),
		ReviewerRole: curriculum.RoleDepartmentHead,
		Action:       curriculum.ActionReject,
		Comments:     "Mục tiêu chưa bám sát chuẩn năng lực toán học của chương trình 2018. Đề nghị bổ sung hoạt động trải nghiệm.",
	})
	if err != nil {
		t.Fatalf("expected rejection to succeed, got %v", err)
	}
	if plan.CurrentStatus != curriculum.StatusRejected {
		t.Errorf("expected status REJECTED, got %s", plan.CurrentStatus)
	}
}

func TestLessonPlanEngine_UnauthorizedReviewer(t *testing.T) {
	engine := curriculum.NewLessonPlanEngine()
	plan := &curriculum.LessonPlan{
		GradeLevel: 1,
		WeekNumber: 2,
		Title:      "Tiếng Việt bài 3",
		SubjectID:  "TIENG_VIET",
	}
	_ = engine.SubmitLessonPlan(plan)

	err := engine.ReviewLessonPlan(plan, curriculum.LessonPlanReview{
		ReviewerID:   uuid.New(),
		ReviewerRole: "PHHS", // Parent cannot review lesson plan
		Action:       curriculum.ActionApprove,
		Comments:     "Đồng ý",
	})
	if err == nil {
		t.Errorf("expected error for unauthorized role, got nil")
	}
}
