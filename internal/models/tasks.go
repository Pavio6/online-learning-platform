package models

import (
	"time"

	"gorm.io/gorm"
)

// Tasks 任务表（中央服务器）
type Tasks struct {
	TaskID      uint           `gorm:"primaryKey;column:task_id" json:"task_id"`
	LessonID    uint           `gorm:"column:lesson_id;not null;index" json:"lesson_id"`
	TaskTitle   string         `gorm:"column:task_title;not null" json:"task_title"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	TaskType    string         `gorm:"column:task_type;default:'essay'" json:"task_type"` // essay, quiz, upload
	MaxScore    int            `gorm:"column:max_score;default:100" json:"max_score"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"column:deleted_at;index" json:"-"`

	// 关联关系
	Lesson *Lessons `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
}

// TableName 指定表名
func (Tasks) TableName() string {
	return "tasks"
}

