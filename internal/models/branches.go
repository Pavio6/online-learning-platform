package models

import (
	"time"

	"gorm.io/gorm"
)

// Branches 分支表
type Branches struct {
	BranchID   uint           `gorm:"primaryKey;column:branch_id" json:"branch_id"`
	BranchName string         `gorm:"column:branch_name;not null" json:"branch_name"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName 指定表名
func (Branches) TableName() string {
	return "branches"
}

