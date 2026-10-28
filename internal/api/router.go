package api

import (
	"qlth/internal/api/handlers"
	"qlth/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up Gin engine and registers all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

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
	}

	return r
}
