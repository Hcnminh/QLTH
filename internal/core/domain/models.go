package domain

import (
	"time"

	"github.com/google/uuid"
)

// SchoolType defines the category of school
type SchoolType string

const (
	SchoolTypeKindergarten SchoolType = "KINDERGARTEN"
	SchoolTypePrimary      SchoolType = "PRIMARY"
	SchoolTypeSecondary    SchoolType = "SECONDARY"
	SchoolTypeHighSchool   SchoolType = "HIGH_SCHOOL"
	SchoolTypeInterLevel   SchoolType = "INTER_LEVEL"
)

// School represents an educational institution tenant
type School struct {
	ID           uuid.UUID  `json:"id"`
	Code         string     `json:"code"` // Mã trường MOET (VD: 'TH_LE_VAN_TAM_01')
	Name         string     `json:"name"`
	Type         SchoolType `json:"type"`
	Address      string     `json:"address"`
	ProvinceCode string     `json:"province_code"`
	DistrictCode string     `json:"district_code"`
	WardCode     string     `json:"ward_code"`
	IsPublic     bool       `json:"is_public"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// AcademicYear represents a school academic year (e.g., '2026-2027')
type AcademicYear struct {
	ID        string    `json:"id"` // '2026-2027'
	SchoolID  uuid.UUID `json:"school_id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Term represents an evaluation semester/period
type Term struct {
	ID             string    `json:"id"` // '2026-2027-HK1'
	AcademicYearID string    `json:"academic_year_id"`
	TermNumber     int       `json:"term_number"` // 1: Giữa HK1, 2: Cuối HK1, 3: Giữa HK2, 4: Cuối Năm
	Name           string    `json:"name"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	IsClosed       bool      `json:"is_closed"`
}

// Class represents a classroom unit
type Class struct {
	ID                uuid.UUID `json:"id"`
	SchoolID          uuid.UUID `json:"school_id"`
	AcademicYearID    string    `json:"academic_year_id"`
	GradeLevel        int       `json:"grade_level"` // 1 to 5
	Name              string    `json:"name"`        // '1A1', '5A2'
	HomeroomTeacherID uuid.UUID `json:"homeroom_teacher_id"`
	RoomNumber        string    `json:"room_number"`
	StudentCount      int       `json:"student_count"`
	CreatedAt         time.Time `json:"created_at"`
}

// Staff represents teachers and administrative personnel (Circular 20/2023)
type Staff struct {
	ID           uuid.UUID `json:"id"`
	SchoolID     uuid.UUID `json:"school_id"`
	NationalID   string    `json:"national_id"`   // 12-digit CCCD/VNeID
	EmployeeCode string    `json:"employee_code"` // GV001...
	FullName     string    `json:"full_name"`
	Gender       string    `json:"gender"` // 'MALE', 'FEMALE'
	BirthDate    time.Time `json:"birth_date"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`

	// Job Title and Rank
	JobPosition string `json:"job_position"` // 'PRINCIPAL', 'PRIMARY_TEACHER', 'ACCOUNTANT', 'NURSE'...
	TeacherRank string `json:"teacher_rank"` // 'RANK_I', 'RANK_II', 'RANK_III'
	IsHomeroom  bool   `json:"is_homeroom"`

	// Teaching Quota according to Circular 05/2025/TT-BGDDT
	StandardAnnualHours int      `json:"standard_annual_hours"` // 805 hours/year for primary
	WeeklyTeachingHours int      `json:"weekly_teaching_hours"` // 23 hours/week average
	WeeklyReducedHours  int      `json:"weekly_reduced_hours"`  // Reduced hours (GVCN: 4, Head Subject: 3)
	ConcurrentRoles     []string `json:"concurrent_roles"`      // Kiêm nhiệm (Tổ trưởng, Tư vấn tâm lý...)

	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Student represents a primary school student (Decree 13/2023 & Law on Identification 2023)
type Student struct {
	ID                   uuid.UUID `json:"id"`
	SchoolID             uuid.UUID `json:"school_id"`
	NationalID           string    `json:"national_id"`   // Mã định danh cá nhân 12 số
	StudentCode          string    `json:"student_code"`  // Mã Bộ GD&ĐT
	FullName             string    `json:"full_name"`
	BirthDate            time.Time `json:"birth_date"`
	Gender               string    `json:"gender"`
	Ethnicity            string    `json:"ethnicity"`
	CurrentClassID       uuid.UUID `json:"current_class_id"`
	Status               string    `json:"status"` // 'STUDYING', 'TRANSFERRED', 'DROPPED', 'GRADUATED'
	IsInclusiveEducation bool      `json:"is_inclusive_education"`

	// Guardian details (VNeID Level 2)
	GuardianName         string `json:"guardian_name"`
	GuardianNationalID   string `json:"guardian_national_id"`
	GuardianPhone        string `json:"guardian_phone"`
	GuardianRelationship string `json:"guardian_relationship"` // 'FATHER', 'MOTHER', 'LEGAL_GUARDIAN'

	// Digital Consent (Decree 13/2023 Article 20 on Children Data Protection)
	HasDigitalConsent bool       `json:"has_digital_consent"`
	ConsentSignedAt   *time.Time `json:"consent_signed_at"`
	ConsentMethod     string     `json:"consent_method"` // 'VNEID_SSO', 'SMS_OTP', 'PORTAL_SIGN'

	CreatedAt time.Time `json:"created_at"`
}

// Subject represents a subject or educational activity
type Subject struct {
	ID               string `json:"id"` // 'TIENG_VIET', 'TOAN', 'TIENG_ANH', 'STEM'...
	Name             string `json:"name"`
	IsScored         bool   `json:"is_scored"` // Has exam score (1-10) or qualitative only
	ApplicableGrades []int  `json:"applicable_grades"`
}

// StudentSubjectEvaluation represents evaluations for subjects (Circular 27/2020 & 89/2026)
type StudentSubjectEvaluation struct {
	ID        uuid.UUID `json:"id"`
	StudentID uuid.UUID `json:"student_id"`
	TermID    string    `json:"term_id"`
	SubjectID string    `json:"subject_id"`
	TeacherID uuid.UUID `json:"teacher_id"`

	// Achievement level: HTT (Hoàn thành tốt), HT (Hoàn thành), CHT (Chưa hoàn thành)
	AchievementLevel string `json:"achievement_level"`

	TeacherComment     string `json:"teacher_comment"`
	AISuggestedComment string `json:"ai_suggested_comment,omitempty"`
	IsAIAssisted       bool   `json:"is_ai_assisted"`

	// Exam score (1.0 to 10.0 for scored subjects)
	ExamScore           *float64 `json:"exam_score,omitempty"`
	IsComputerBasedTest bool     `json:"is_computer_based_test"` // TT 89/2026
	IsCentralizedExam   bool     `json:"is_centralized_exam"`   // Đề chung Sở/Phòng Lớp 5 (TT 89/2026)
	UpdatedAt           time.Time `json:"updated_at"`
}

// StudentBehaviorEvaluation represents 5 qualities and 10 competencies
type StudentBehaviorEvaluation struct {
	ID        uuid.UUID `json:"id"`
	StudentID uuid.UUID `json:"student_id"`
	TermID    string    `json:"term_id"`
	TeacherID uuid.UUID `json:"teacher_id"`

	// 5 Qualities: 'T' (Tốt), 'D' (Đạt), 'C' (Cần cố gắng)
	Patriotism     string `json:"patriotism"`
	Compassion     string `json:"compassion"`
	Diligence      string `json:"diligence"`
	Honesty        string `json:"honesty"`
	Responsibility string `json:"responsibility"`

	// 3 General Competencies
	AutonomySelfLearning     string `json:"autonomy_self_learning"`
	CommunicationCooperation string `json:"communication_cooperation"`
	ProblemSolvingCreativity string `json:"problem_solving_creativity"`

	// 7 Specific Competencies
	LanguageCompetence      string `json:"language_competence"`
	MathematicalCompetence  string `json:"mathematical_competence"`
	ScientificCompetence    string `json:"scientific_competence"`
	TechnologicalCompetence string `json:"technological_competence"`
	ComputerCompetence      string `json:"computer_competence"`
	AestheticCompetence     string `json:"aesthetic_competence"`
	PhysicalCompetence      string `json:"physical_competence"`

	GeneralComment string    `json:"general_comment"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BoardingMenu represents a daily menu with nutritional facts
type BoardingMenu struct {
	ID                     uuid.UUID `json:"id"`
	SchoolID               uuid.UUID `json:"school_id"`
	MenuDate               time.Time `json:"menu_date"`
	MealType               string    `json:"meal_type"` // 'LUNCH', 'SNACK', 'BREAKFAST'
	DishesDescription      string    `json:"dishes_description"`
	TotalKcal              float64   `json:"total_kcal"` // 600 - 750 Kcal for primary lunch
	ProteinPercent         float64   `json:"protein_percent"`
	LipidPercent           float64   `json:"lipid_percent"`
	GlucidPercent          float64   `json:"glucid_percent"`
	SampleRetentionVerified bool      `json:"sample_retention_verified"` // Lưu mẫu 24h
	PhotoURL               string    `json:"photo_url"`
}

// BoardingDailyAttendance represents daily meal attendance
type BoardingDailyAttendance struct {
	ID                  uuid.UUID  `json:"id"`
	StudentID           uuid.UUID  `json:"student_id"`
	AttendanceDate      time.Time  `json:"attendance_date"`
	IsEating            bool       `json:"is_eating"`
	IsCancelledByParent bool       `json:"is_cancelled_by_parent"`
	CancellationTime    *time.Time `json:"cancellation_time"`
	IsRefundEligible    bool       `json:"is_refund_eligible"` // True if cancelled before 8:00 AM
	RefundAmount        float64    `json:"refund_amount"`
}
