package api

import (
	"log"

	"qlth/internal/admin"
	"qlth/internal/api/handlers"
	"qlth/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up Gin engine and registers all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Mount GoAdmin Web Interface at /admin
	if _, err := admin.InitAdmin(r); err != nil {
		log.Printf("Warning: failed to initialize GoAdmin: %v", err)
	}

	h := handlers.NewAPIHandler()

	// Public Health check
	r.GET("/health", h.HealthCheck)

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Teaching Quota (Circular 05/2025)
		v1.POST("/teaching/calculate-quota",
			middleware.RequireRole("BGH", "KE_TOAN", "GVCN", "GVBM"),
			h.CalculateTeachingQuota,
		)

		// KPI Evaluation (Decree 90/2020 & 48/2023)
		v1.POST("/kpi/evaluate-individual",
			middleware.RequireRole("BGH", "TO_TRUONG"),
			h.EvaluateKPI,
		)
		v1.POST("/kpi/evaluate-school-batch",
			middleware.RequireRole("BGH"),
			h.BatchEvaluateKPI,
		)

		// AI Pedagogical Copilot (Circular 27/2020 & 89/2026)
		v1.POST("/evaluations/ai-copilot",
			middleware.RequireRole("BGH", "GVCN", "GVBM"),
			h.AIPedagogicalCopilot,
		)

		// Boarding meal cancellation with automatic refund check
		v1.POST("/boarding/cancel-meal",
			middleware.RequireRole("BGH", "KE_TOAN", "GVCN", "PHHS"),
			h.CancelBoardingMeal,
		)

		// Student Onboarding & VNeID
		v1.POST("/students/onboard",
			middleware.RequireRole("BGH", "VAN_THU"),
			h.StudentOnboarding,
		)

		// Student Subject Evaluations (Circular 27/2020 & 89/2026)
		v1.POST("/evaluations/subject-scores",
			middleware.RequireRole("BGH", "GVCN", "GVBM"),
			h.RecordSubjectScore,
		)

		// Student Behavior Evaluations (Circular 27/2020 & 89/2026)
		v1.POST("/evaluations/behaviors",
			middleware.RequireRole("BGH", "GVCN"),
			h.RecordBehaviorEvaluation,
		)

		// Academic Transcripts & Snapshot Generation (Circular 27/2020 & 89/2026)
		v1.POST("/transcripts/generate",
			middleware.RequireRole("BGH", "GVCN"),
			h.GenerateTranscript,
		)

		// Lesson Plans (CV 2345/STEM TT 86)
		v1.POST("/lesson-plans/submit",
			middleware.RequireRole("BGH", "TO_TRUONG", "GVCN", "GVBM", "GV"),
			h.SubmitLessonPlan,
		)
		v1.POST("/lesson-plans/:id/approve",
			middleware.RequireRole("BGH", "TO_TRUONG"),
			h.ApproveLessonPlan,
		)
		v1.GET("/lesson-plans/:id",
			middleware.RequireRole("BGH", "TO_TRUONG", "GVCN", "GVBM", "GV"),
			h.GetLessonPlan,
		)
		v1.GET("/lesson-plans",
			middleware.RequireRole("BGH", "TO_TRUONG", "GVCN", "GVBM", "GV"),
			h.ListLessonPlans,
		)

		// Timetable Slot Configuration & Management (Circular 05/2025)
		v1.POST("/timetable/slots",
			middleware.RequireRole("BGH", "TO_TRUONG", "GVCN", "GVBM"),
			h.ConfigureTimetableSlot,
		)
		v1.GET("/timetable/slots",
			middleware.RequireRole("BGH", "TO_TRUONG", "GVCN", "GVBM", "PHHS"),
			h.GetTimetableSlots,
		)

		// Boarding Menus & 24h Food Sample Retention (Viện Dinh Dưỡng & QĐ 1246/QĐ-BYT)
		v1.GET("/boarding/menus",
			middleware.RequireRole("BGH", "KE_TOAN", "GVCN", "GVBM", "PHHS", "Y_TE", "NHAN_VIEN_Y_TE"),
			h.ListBoardingMenus,
		)
		v1.POST("/boarding/menus",
			middleware.RequireRole("BGH", "KE_TOAN", "Y_TE", "NHAN_VIEN_Y_TE", "BEP_TRUONG"),
			h.CreateBoardingMenu,
		)
		v1.POST("/boarding/menus/:id/sample-retention",
			middleware.RequireRole("BGH", "Y_TE", "NHAN_VIEN_Y_TE", "BEP_TRUONG"),
			h.VerifyFoodSampleRetention,
		)

		// Student Health Profile & Infirmary Diary (BMI, Eyesight, Medical Diary, Prescriptions)
		v1.GET("/students/:id/health",
			middleware.RequireRole("BGH", "GVCN", "Y_TE", "NHAN_VIEN_Y_TE", "PHHS"),
			h.GetStudentHealthProfile,
		)
		v1.POST("/students/:id/health",
			middleware.RequireRole("BGH", "Y_TE", "NHAN_VIEN_Y_TE"),
			h.RecordStudentHealth,
		)
		v1.POST("/students/:id/health/prescriptions",
			middleware.RequireRole("BGH", "GVCN", "PHHS", "Y_TE", "NHAN_VIEN_Y_TE"),
			h.AddPrescriptionDelivery,
		)
		v1.POST("/students/:id/health/prescriptions/:rx_id/administer",
			middleware.RequireRole("BGH", "Y_TE", "NHAN_VIEN_Y_TE"),
			h.AdministerPrescription,
		)
		v1.POST("/students/:id/health/medical-events",
			middleware.RequireRole("BGH", "Y_TE", "NHAN_VIEN_Y_TE"),
			h.AddMedicalEvent,
		)

		// External Integration Sandbox Stubs (Phân hệ 7, 9, 10)
		v1.POST("/integrations/smartca/sign",
			middleware.RequireRole("BGH", "GVCN", "GVBM", "TO_TRUONG", "GV"),
			h.SignSmartCA,
		)
		v1.POST("/integrations/vietqr/generate",
			middleware.RequireRole("BGH", "KE_TOAN", "GVCN", "PHHS"),
			h.GenerateVietQR,
		)
		v1.POST("/integrations/moet/sync",
			middleware.RequireRole("BGH", "VAN_THU"),
			h.SyncMOET,
		)
		v1.GET("/integrations/moet/queue",
			middleware.RequireRole("BGH", "VAN_THU"),
			h.GetMOETQueue,
		)
		v1.POST("/integrations/vneid/verify",
			middleware.RequireRole("BGH", "VAN_THU", "GVCN", "PHHS"),
			h.VerifyVNeID,
		)
		v1.POST("/integrations/einvoice/issue",
			middleware.RequireRole("BGH", "KE_TOAN"),
			h.IssueEInvoice,
		)
		v1.GET("/integrations/sandbox/status",
			middleware.RequireRole("BGH", "KE_TOAN", "VAN_THU", "GVCN", "GVBM", "PHHS", "Y_TE", "TO_TRUONG"),
			h.GetSandboxStatus,
		)
	}

	return r
}
