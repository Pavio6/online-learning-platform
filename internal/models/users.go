package models

import (
	"time"

	"gorm.io/gorm"
)

// Users 用户表
type Users struct {
	UserID      uint           `gorm:"primaryKey;column:user_id" json:"user_id"`
	BranchID    uint           `gorm:"column:branch_id;not null;index" json:"branch_id"`
	Username    string         `gorm:"column:username;not null;uniqueIndex" json:"username"`
	Email       string         `gorm:"column:email;not null;uniqueIndex" json:"email"`
	PasswordHash string        `gorm:"column:password_hash;not null" json:"-"`
	FirstName   string         `gorm:"column:first_name" json:"first_name"`
	LastName    string         `gorm:"column:last_name" json:"last_name"`
	Role        string         `gorm:"column:role;not null;default:'student'" json:"role"` // student, teacher
	Status      string         `gorm:"column:status;default:'active'" json:"status"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	// 关联关系
	Branch *Branches `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

// TableName 指定表名
func (Users) TableName() string {
	return "users"
}

