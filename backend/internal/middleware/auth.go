package middleware

import (
	"net/http"
	"strings"

	"cybertron-portal/internal/service"
	"cybertron-portal/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "userID"
	ContextUsername = "username"
	ContextRole     = "role"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    401,
				Message: "未提供认证令牌",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    401,
				Message: "认证令牌格式错误",
			})
			return
		}

		claims, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    401,
				Message: "认证令牌无效或已过期",
			})
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextRole, claims.Role)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRole)
		if !exists || role.(string) != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Response{
				Code:    403,
				Message: "无权限访问",
			})
			return
		}
		c.Next()
	}
}
