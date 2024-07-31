package middlewares

import (
	"net/http"
	"recyco/config"
	"recyco/models"
	"recyco/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondFailed(c, http.StatusUnauthorized, "Authorization header required", nil)
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondFailed(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
			c.Abort()
			return
		}
		token := parts[1]
		claims, err := utils.ValidateToken(token)
		if err != nil {
			utils.RespondFailed(c, http.StatusUnauthorized, "Invalid token", nil)
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)

		var user models.User
		if err := config.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			utils.RespondFailed(c, http.StatusUnauthorized, "User not found", nil)
			c.Abort()
			return
		}

		c.Set("userRole", user.Role)
		c.Next()
	}
}
