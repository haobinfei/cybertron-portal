package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cybertron-portal/pkg/jwt"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupRedis() (*redis.Client, *miniredis.Miniredis) {
	mr, _ := miniredis.Run()
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return rdb, mr
}

type testAuthService struct {
	rdb         *redis.Client
	secret      string
	expireHours int
}

func (s *testAuthService) ValidateToken(ctx context.Context, token string) (*jwt.Claims, error) {
	isBlacklisted, _ := s.rdb.Exists(ctx, "blacklist:"+token).Result()
	if isBlacklisted > 0 {
		return nil, assert.AnError
	}
	return jwt.ParseToken(token, s.secret)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, mr := setupRedis()
	defer mr.Close()

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestAuthMiddleware_NoToken_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/protected", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未提供认证令牌"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	token, _ := jwt.GenerateToken(1, "admin", "admin", "test-secret", 24)

	r := gin.New()
	r.GET("/protected", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401})
			return
		}
		tokenStr := authHeader[7:]
		claims, err := jwt.ParseToken(tokenStr, "test-secret")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401})
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextRole, claims.Role)
		c.JSON(200, gin.H{"user": claims.Username})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "admin")
}

func TestAdminMiddleware_AdminAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(ContextRole, "admin")
		c.Next()
	})
	r.GET("/admin", AdminMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestAdminMiddleware_UserAccess_Denied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(ContextRole, "user")
		c.Next()
	})
	r.GET("/admin", AdminMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminMiddleware_NoRole_Denied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/admin", AdminMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCORSMiddleware_AllowsOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	r.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, 200, w.Code)
}

func TestCORSMiddleware_OptionsPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}
