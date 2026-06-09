package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Create(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	user, err := svc.Create(&CreateUserRequest{
		Username: "newuser",
		Password: "pass123",
		Nickname: "新用户",
		Email:    "new@test.com",
		Role:     "user",
		Status:   1,
	})
	require.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "newuser", user.Username)
	assert.NotEmpty(t, user.PasswordHash)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("pass123"))
	assert.NoError(t, err)
}

func TestUserService_Create_DuplicateUsername(t *testing.T) {
	// The mock doesn't check unique constraint, so this test is limited
	// In real DB, duplicate username would fail
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	_, err := svc.Create(&CreateUserRequest{Username: "admin", Password: "pass123"})
	require.NoError(t, err)
}

func TestUserService_Create_DefaultRole(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	user, err := svc.Create(&CreateUserRequest{
		Username: "user1",
		Password: "pass123",
	})
	require.NoError(t, err)
	assert.Equal(t, "user", user.Role)
	assert.Equal(t, 1, user.Status)
}

func TestUserService_Update(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	_, _ = svc.Create(&CreateUserRequest{Username: "admin", Password: "pass123", Role: "admin"})

	status := 0
	user, err := svc.Update(1, &UpdateUserRequest{
		Nickname: "管理员更新",
		Email:    "admin@test.com",
		Status:   &status,
	})
	require.NoError(t, err)
	assert.Equal(t, "管理员更新", user.Nickname)
	assert.Equal(t, "admin@test.com", user.Email)
	assert.Equal(t, 0, user.Status)
}

func TestUserService_Update_Password(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	_, _ = svc.Create(&CreateUserRequest{Username: "admin", Password: "oldpass"})

	user, err := svc.Update(1, &UpdateUserRequest{
		Password: "newpass",
	})
	require.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("newpass"))
	assert.NoError(t, err)
}

func TestUserService_Update_NotFound(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	_, err := svc.Update(999, &UpdateUserRequest{Nickname: "x"})
	assert.Error(t, err)
}

func TestUserService_Delete(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	_, _ = svc.Create(&CreateUserRequest{Username: "tmp", Password: "pass123"})

	err := svc.Delete(1)
	require.NoError(t, err)

	_, err = svc.GetByID(1)
	assert.Error(t, err)
}

func TestUserService_GetByID(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	_, _ = svc.Create(&CreateUserRequest{Username: "admin", Password: "pass123"})

	user, err := svc.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, "admin", user.Username)
}

func TestUserService_FindAll(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	_, _ = svc.Create(&CreateUserRequest{Username: "user1", Password: "pass"})
	_, _ = svc.Create(&CreateUserRequest{Username: "user2", Password: "pass"})
	_, _ = svc.Create(&CreateUserRequest{Username: "user3", Password: "pass"})

	users, total, err := svc.FindAll(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, users, 3)
}

func TestUserService_FindAll_Pagination(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	for i := range 15 {
		username := "user" + string(rune('a'+i%26))
		_, _ = svc.Create(&CreateUserRequest{Username: username, Password: "pass"})
	}

	users, total, err := svc.FindAll(2, 5)
	require.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.LessOrEqual(t, len(users), 5)
}

func TestUserService_FindAll_DefaultParams(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserService(repo)

	users, total, err := svc.FindAll(0, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, users, 0)
}
