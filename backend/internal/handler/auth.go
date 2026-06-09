package handler

import (
	"strings"

	"cybertron-portal/internal/service"
	"cybertron-portal/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.SuccessWithMessage(c, "登录成功", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		response.Error(c, 400, "缺少令牌")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		response.Error(c, 500, "退出登录失败")
		return
	}

	response.SuccessWithMessage(c, "已退出登录", nil)
}
