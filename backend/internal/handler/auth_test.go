package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cybertron-portal/internal/model"
	"cybertron-portal/internal/service"
	jwtpkg "cybertron-portal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type testUserRepo struct {
	users  map[uint]*model.User
	byName map[string]*model.User
	nextID uint
}

func newTestUserRepo() *testUserRepo {
	return &testUserRepo{
		users:  make(map[uint]*model.User),
		byName: make(map[string]*model.User),
		nextID: 1,
	}
}

func (m *testUserRepo) Create(user *model.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.byName[user.Username] = user
	return nil
}

func (m *testUserRepo) FindByID(id uint) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (m *testUserRepo) FindByUsername(username string) (*model.User, error) {
	u, ok := m.byName[username]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (m *testUserRepo) FindAll(page, pageSize int) ([]model.User, int64, error) {
	result := make([]model.User, 0)
	for _, u := range m.users {
		result = append(result, *u)
	}
	total := int64(len(result))
	offset := (page - 1) * pageSize
	if offset >= len(result) {
		return []model.User{}, total, nil
	}
	end := offset + pageSize
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], total, nil
}

func (m *testUserRepo) Update(user *model.User) error {
	m.users[user.ID] = user
	m.byName[user.Username] = user
	return nil
}

func (m *testUserRepo) Delete(id uint) error {
	u, ok := m.users[id]
	if ok {
		delete(m.byName, u.Username)
		delete(m.users, id)
	}
	return nil
}

func (m *testUserRepo) UpdateLastLogin(id uint) error {
	return nil
}

type testAuthService struct {
	repo        *testUserRepo
	secret      string
	expireHours int
}

func (s *testAuthService) Login(ctx context.Context, username, password string) (*service.LoginResult, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user.Status != 1 {
		return nil, gorm.ErrRecordNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	token, _ := jwtpkg.GenerateToken(user.ID, user.Username, user.Role, s.secret, s.expireHours)
	return &service.LoginResult{Token: token, UserInfo: *user}, nil
}

func (s *testAuthService) Logout(ctx context.Context, token string) error {
	return nil
}

func (s *testAuthService) ValidateToken(ctx context.Context, token string) (*jwtpkg.Claims, error) {
	return jwtpkg.ParseToken(token, s.secret)
}

func (s *testAuthService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newTestUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	repo.Create(&model.User{Username: "admin", PasswordHash: string(hash), Role: "admin", Status: 1})

	svc := &testAuthService{repo: repo, secret: "test-secret", expireHours: 24}
	handler := NewAuthHandler(svc)

	r := gin.New()
	r.POST("/api/auth/login", handler.Login)

	body := map[string]string{"username": "admin", "password": "admin123"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "登录成功", resp.Message)
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newTestUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	repo.Create(&model.User{Username: "admin", PasswordHash: string(hash), Status: 1})

	svc := &testAuthService{repo: repo, secret: "test-secret", expireHours: 24}
	handler := NewAuthHandler(svc)

	r := gin.New()
	r.POST("/api/auth/login", handler.Login)

	body := map[string]string{"username": "admin", "password": "wrong"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 401, resp.Code)
}

func TestAuthHandler_Login_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newTestUserRepo()
	svc := &testAuthService{repo: repo, secret: "test-secret", expireHours: 24}
	handler := NewAuthHandler(svc)

	r := gin.New()
	r.POST("/api/auth/login", handler.Login)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", nil)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 400, resp.Code)
}

func TestAuthHandler_Logout_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newTestUserRepo()
	svc := &testAuthService{repo: repo, secret: "test-secret", expireHours: 24}
	handler := NewAuthHandler(svc)

	r := gin.New()
	r.POST("/api/auth/logout", handler.Logout)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 400, resp.Code)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newTestUserRepo()
	svc := &testAuthService{repo: repo, secret: "test-secret", expireHours: 24}
	handler := NewAuthHandler(svc)

	r := gin.New()
	r.POST("/api/auth/logout", handler.Logout)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	r.ServeHTTP(w, req)

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "已退出登录", resp.Message)
}
