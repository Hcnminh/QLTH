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
	}

	return r
}
