package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cybertron-portal/internal/middleware"
	"cybertron-portal/internal/model"
	"cybertron-portal/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type testUserServiceForHandler struct {
	repo *testUserRepo
}

func (s *testUserServiceForHandler) GetByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *testUserServiceForHandler) FindAll(page, pageSize int) ([]model.User, int64, error) {
	return s.repo.FindAll(page, pageSize)
}

func (s *testUserServiceForHandler) Create(req *service.CreateUserRequest) (*model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		Nickname:     req.Nickname,
		Email:        req.Email,
		Role:         req.Role,
		Status:       req.Status,
	}
	if user.Role == "" {
		user.Role = "user"
	}
	if user.Status == 0 {
		user.Status = 1
	}
	s.repo.Create(user)
	return user, nil
}

func (s *testUserServiceForHandler) Update(id uint, req *service.UpdateUserRequest) (*model.User, error) {
	user, _ := s.repo.FindByID(id)
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user.PasswordHash = string(hash)
	}
	s.repo.Update(user)
	return user, nil
}

func (s *testUserServiceForHandler) Delete(id uint) error {
	return s.repo.Delete(id)
}

func TestUserHandler_GetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newTestUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo.Create(&model.User{Username: "testuser", PasswordHash: string(hash), Nickname: "测试用户", Role: "admin", Status: 1})

	handler := NewUserHandler(&testUserServiceForHandler{repo: repo})

	r := gin.New()
	r.GET("/api/user/me", func(c *gin.Context) {
		c.Set(middleware.ContextUserID, uint(1))
		c.Next()
	}, handler.GetMe)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/me", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int        `json:"code"`
		Message string     `json:"message"`
		Data    model.User `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "testuser", resp.Data.Username)
}

func TestUserHandler_GetMe_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewUserHandler(&testUserServiceForHandler{repo: newTestUserRepo()})

	r := gin.New()
	r.GET("/api/user/me", func(c *gin.Context) {
		c.Set(middleware.ContextUserID, uint(999))
		c.Next()
	}, handler.GetMe)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/me", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 404, resp.Code)
}

func TestUserHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newTestUserRepo()
	for i := range 5 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
		repo.Create(&model.User{
			Username:     "user" + string(rune('a'+i)),
			PasswordHash: string(hash),
			Role:         "user",
			Status:       1,
		})
	}

	handler := NewUserHandler(&testUserServiceForHandler{repo: repo})

	r := gin.New()
	r.GET("/api/users", handler.List)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			List  []model.User `json:"list"`
			Total int64        `json:"total"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, int64(5), resp.Data.Total)
	assert.Len(t, resp.Data.List, 5)
}

func TestUserHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewUserHandler(&testUserServiceForHandler{repo: newTestUserRepo()})

	r := gin.New()
	r.POST("/api/users", handler.Create)

	body := map[string]interface{}{
		"username": "newuser",
		"password": "pass123",
		"nickname": "新用户",
		"role":     "user",
		"status":   1,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int        `json:"code"`
		Message string     `json:"message"`
		Data    model.User `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "创建用户成功", resp.Message)
	assert.Equal(t, "newuser", resp.Data.Username)
}

func TestUserHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewUserHandler(&testUserServiceForHandler{repo: newTestUserRepo()})

	r := gin.New()
	r.POST("/api/users", handler.Create)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 400, resp.Code)
}

func TestUserHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newTestUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo.Create(&model.User{Username: "editme", PasswordHash: string(hash), Nickname: "old", Role: "user", Status: 1})

	handler := NewUserHandler(&testUserServiceForHandler{repo: repo})

	r := gin.New()
	r.PUT("/api/users/:id", handler.Update)

	body := map[string]interface{}{
		"nickname": "updated",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/users/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int        `json:"code"`
		Message string     `json:"message"`
		Data    model.User `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "updated", resp.Data.Nickname)
}

func TestUserHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newTestUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo.Create(&model.User{Username: "deleteme", PasswordHash: string(hash)})

	handler := NewUserHandler(&testUserServiceForHandler{repo: repo})

	r := gin.New()
	r.DELETE("/api/users/:id", handler.Delete)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/users/1", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "删除用户成功", resp.Message)
}

func TestUserHandler_Delete_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewUserHandler(&testUserServiceForHandler{repo: newTestUserRepo()})

	r := gin.New()
	r.DELETE("/api/users/:id", handler.Delete)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/users/abc", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 400, resp.Code)
}
