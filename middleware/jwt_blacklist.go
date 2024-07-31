package middlewares

import (
	"net/http"
	"strings"
	"sync"

	"latihan/utils"

	"github.com/gin-gonic/gin"
)

var (
	blacklistedTokens = make(map[string]struct{})
	mu                sync.Mutex
)

func JWTBlacklistMiddleware() gin.HandlerFunc {
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

		mu.Lock()
		_, exists := blacklistedTokens[token]
		mu.Unlock()

		if exists {
			utils.RespondFailed(c, http.StatusUnauthorized, "Token has been revoked", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func AddToBlacklist(token string) {
	mu.Lock()
	blacklistedTokens[token] = struct{}{}
	mu.Unlock()
}
