package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if the authenticated user has one of the required roles
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract role from context or header for demonstration
		role := c.GetHeader("X-User-Role")
		if role == "" {
			role = "GUEST"
		}

		schoolID := c.GetHeader("X-School-ID")
		if schoolID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    401,
				"message": "Thiếu mã định danh trường học (X-School-ID)",
			})
			c.Abort()
			return
		}

		isAllowed := false
		for _, r := range allowedRoles {
			if strings.EqualFold(r, role) || role == "SUPER_ADMIN" {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    403,
				"message": "Không có quyền thực hiện thao tác này",
			})
			c.Abort()
			return
		}

		c.Set("user_role", role)
		c.Set("school_id", schoolID)
		c.Next()
	}
}
