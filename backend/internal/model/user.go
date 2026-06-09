package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash string         `gorm:"size:256;not null" json:"-"`
	Nickname     string         `gorm:"size:64" json:"nickname"`
	Email        string         `gorm:"size:128" json:"email"`
	Avatar       string         `gorm:"size:256" json:"avatar"`
	Role         string         `gorm:"size:32;default:user" json:"role"`
	Status       int            `gorm:"default:1" json:"status"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
