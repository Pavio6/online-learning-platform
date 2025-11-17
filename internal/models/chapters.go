package models

import (
	"time"

	"gorm.io/gorm"
)

// Chapters 章节表（中央服务器）
type Chapters struct {
	ChapterID    uint           `gorm:"primaryKey;column:chapter_id" json:"chapter_id"`
	CourseID     uint           `gorm:"column:course_id;not null;index" json:"course_id"`
	ChapterTitle string         `gorm:"column:chapter_title;not null" json:"chapter_title"`
	ChapterOrder int            `gorm:"column:chapter_order;not null" json:"chapter_order"`
	Description  string         `gorm:"column:description;type:text" json:"description"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	// 关联关系
	Course  *Courses  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Lessons []Lessons `gorm:"foreignKey:ChapterID" json:"lessons,omitempty"`
}

// TableName 指定表名
func (Chapters) TableName() string {
	return "chapters"
}
