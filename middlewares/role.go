package middlewares

import (
	"net/http"
	"recyco/utils"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			utils.RespondFailed(c, http.StatusForbidden, "User role not found in context", nil)
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		utils.RespondFailed(c, http.StatusForbidden, "You do not have permission to access this resource", nil)
		c.Abort()
	}
}
