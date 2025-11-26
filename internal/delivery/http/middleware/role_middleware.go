package middleware

import (
	"event-campus-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if user has one of the required roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(403, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   "User role not found",
			})
			c.Abort()
			return
		}

		role := userRole.(string)

		// Check if user has one of the required roles
		hasRole := false
		for _, requiredRole := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(403, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOrganisasi checks if user is organisasi or admin
func RequireOrganisasi() gin.HandlerFunc {
	return RequireRole(domain.RoleOrganisasi, domain.RoleAdmin)
}

// RequireAdmin checks if user is admin
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(domain.RoleAdmin)
}
