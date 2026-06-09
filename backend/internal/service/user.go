package service

import (
	"fmt"

	"cybertron-portal/internal/model"
	"cybertron-portal/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	GetByID(id uint) (*model.User, error)
	FindAll(page, pageSize int) ([]model.User, int64, error)
	Create(req *CreateUserRequest) (*model.User, error)
	Update(id uint, req *UpdateUserRequest) (*model.User, error)
	Delete(id uint) error
}

type UserService struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   int    `json:"status"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   *int   `json:"status"`
	Password string `json:"password"`
}

func (s *UserService) Create(req *CreateUserRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	if req.Role == "" {
		req.Role = "user"
	}
	if req.Status == 0 {
		req.Status = 1
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Nickname:     req.Nickname,
		Email:        req.Email,
		Role:         req.Role,
		Status:       req.Status,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(id uint, req *UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

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
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("密码加密失败: %w", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	return user, nil
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) FindAll(page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.userRepo.FindAll(page, pageSize)
}
