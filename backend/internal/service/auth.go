package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cybertron-portal/internal/model"
	"cybertron-portal/internal/repository"
	jwtpkg "cybertron-portal/pkg/jwt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, tokenString string) (*jwtpkg.Claims, error)
	GetUserByID(id uint) (*model.User, error)
}

type AuthService struct {
	userRepo    repository.UserRepositoryInterface
	redis       *redis.Client
	jwtSecret   string
	expireHours int
}

func NewAuthService(userRepo repository.UserRepositoryInterface, redis *redis.Client, jwtSecret string, expireHours int) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		redis:       redis,
		jwtSecret:   jwtSecret,
		expireHours: expireHours,
	}
}

type LoginResult struct {
	Token    string     `json:"token"`
	UserInfo model.User `json:"user_info"`
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Username, user.Role, s.jwtSecret, s.expireHours)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	ttl := time.Duration(s.expireHours) * time.Hour
	if err := s.redis.Set(ctx, tokenKey(token), user.ID, ttl).Err(); err != nil {
		return nil, fmt.Errorf("缓存令牌失败: %w", err)
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	return &LoginResult{
		Token:    token,
		UserInfo: *user,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	ttl := time.Duration(s.expireHours) * time.Hour
	return s.redis.Set(ctx, blacklistKey(token), "1", ttl).Err()
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*jwtpkg.Claims, error) {
	isBlacklisted, err := s.redis.Exists(ctx, blacklistKey(tokenString)).Result()
	if err != nil {
		return nil, fmt.Errorf("检查令牌黑名单失败: %w", err)
	}
	if isBlacklisted > 0 {
		return nil, errors.New("令牌已失效")
	}

	claims, err := jwtpkg.ParseToken(tokenString, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *AuthService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func tokenKey(token string) string {
	return "token:" + token
}

func blacklistKey(token string) string {
	return "blacklist:" + token
}
