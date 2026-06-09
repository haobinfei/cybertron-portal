package service

import (
	"context"
	"testing"

	"cybertron-portal/internal/model"
	"cybertron-portal/pkg/jwt"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserRepo struct {
	users map[uint]*model.User
	byName map[string]*model.User
	nextID uint
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:  make(map[uint]*model.User),
		byName: make(map[string]*model.User),
		nextID: 1,
	}
}

func (m *mockUserRepo) Create(user *model.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.byName[user.Username] = user
	return nil
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *mockUserRepo) FindByUsername(username string) (*model.User, error) {
	user, ok := m.byName[username]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *mockUserRepo) FindAll(page, pageSize int) ([]model.User, int64, error) {
	result := make([]model.User, 0, len(m.users))
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

func (m *mockUserRepo) Update(user *model.User) error {
	m.users[user.ID] = user
	m.byName[user.Username] = user
	return nil
}

func (m *mockUserRepo) Delete(id uint) error {
	user, ok := m.users[id]
	if ok {
		delete(m.byName, user.Username)
		delete(m.users, id)
	}
	return nil
}

func (m *mockUserRepo) UpdateLastLogin(id uint) error {
	return nil
}

func TestAuthService_Login_Success(t *testing.T) {
	repo := newMockUserRepo()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	authSvc := NewAuthService(repo, rdb, "test-secret", 24)

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	repo.Create(&model.User{
		Username:     "admin",
		PasswordHash: string(hash),
		Role:         "admin",
		Status:       1,
	})

	result, err := authSvc.Login(context.Background(), "admin", "pass123")
	require.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, "admin", result.UserInfo.Username)

	claims, err := jwt.ParseToken(result.Token, "test-secret")
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "admin", claims.Username)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	authSvc := NewAuthService(repo, rdb, "test-secret", 24)

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	repo.Create(&model.User{
		Username:     "admin",
		PasswordHash: string(hash),
		Status:       1,
	})

	_, err = authSvc.Login(context.Background(), "admin", "wrongpass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户名或密码错误")
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	authSvc := NewAuthService(repo, rdb, "test-secret", 24)

	_, err = authSvc.Login(context.Background(), "nouser", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户名或密码错误")
}

func TestAuthService_Login_DisabledUser(t *testing.T) {
	repo := newMockUserRepo()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	authSvc := NewAuthService(repo, rdb, "test-secret", 24)

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	repo.Create(&model.User{
		Username:     "disabled",
		PasswordHash: string(hash),
		Status:       0,
	})

	_, err = authSvc.Login(context.Background(), "disabled", "pass123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已被禁用")
}

func TestAuthService_Logout_Success(t *testing.T) {
	repo := newMockUserRepo()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	authSvc := NewAuthService(repo, rdb, "test-secret", 24)

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo.Create(&model.User{Username: "u", PasswordHash: string(hash), Status: 1})

	result, _ := authSvc.Login(context.Background(), "u", "pass")

	err = authSvc.Logout(context.Background(), result.Token)
	require.NoError(t, err)

	_, err = authSvc.ValidateToken(context.Background(), result.Token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已失效")
}

func TestAuthService_ValidateToken_Valid(t *testing.T) {
	repo := newMockUserRepo()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	authSvc := NewAuthService(repo, rdb, "test-secret", 24)

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo.Create(&model.User{Username: "u", PasswordHash: string(hash), Status: 1})

	result, _ := authSvc.Login(context.Background(), "u", "pass")

	claims, err := authSvc.ValidateToken(context.Background(), result.Token)
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
}
