package models

import (
	"time"

	"gorm.io/gorm"
)

// Learning 学习进度表（分支节点）
type Learning struct {
	LearningID        uint           `gorm:"primaryKey;column:learning_id" json:"learning_id"`
	UserID            uint           `gorm:"column:user_id;not null;index" json:"user_id"`
	CourseID          uint           `gorm:"column:course_id;not null;index" json:"course_id"`
	Status            string         `gorm:"column:status;default:'enrolled'" json:"status"` // enrolled, in_progress, completed
	ProgressPercentage int           `gorm:"column:progress_percentage;default:0" json:"progress_percentage"`
	CompletedAt       *time.Time     `gorm:"column:completed_at" json:"completed_at"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	// 关联关系
	User   *Users `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// Course关联需要在查询时从中央服务器获取
}

// TableName 指定表名
func (Learning) TableName() string {
	return "learning"
}

