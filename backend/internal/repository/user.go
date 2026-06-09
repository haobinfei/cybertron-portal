package repository

import (
	"cybertron-portal/internal/model"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAll(page, pageSize int) ([]model.User, int64, error)
	Update(user *model.User) error
	Delete(id uint) error
	UpdateLastLogin(id uint) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn("last_login_at", gorm.Expr("NOW()")).Error
}
