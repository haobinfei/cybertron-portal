package handler

import (
	"strconv"

	"cybertron-portal/internal/middleware"
	"cybertron-portal/internal/service"
	"cybertron-portal/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)

	user, err := h.userService.GetByID(userID)
	if err != nil {
		response.Error(c, 404, "用户不存在")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.userService.FindAll(page, pageSize)
	if err != nil {
		response.Error(c, 500, "查询用户列表失败")
		return
	}

	response.Page(c, users, total, page, pageSize)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Create(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建用户成功", user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的用户ID")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Update(uint(id), &req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新用户成功", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的用户ID")
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		response.Error(c, 500, "删除用户失败")
		return
	}

	response.SuccessWithMessage(c, "删除用户成功", nil)
}
