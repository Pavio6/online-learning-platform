package models

import (
	"time"

	"gorm.io/gorm"
)

// Lessons 课程表（中央服务器）
type Lessons struct {
	LessonID    uint           `gorm:"primaryKey;column:lesson_id" json:"lesson_id"`
	CourseID    uint           `gorm:"column:course_id;not null;index" json:"course_id"`
	ChapterID   uint           `gorm:"column:chapter_id;not null;index" json:"chapter_id"`
	LessonTitle string         `gorm:"column:lesson_title;not null" json:"lesson_title"`
	ContentURL  string         `gorm:"column:content_url" json:"content_url"` // OSS视频链接
	LessonType  string         `gorm:"column:lesson_type;default:'video'" json:"lesson_type"` // video, text, quiz
	LessonOrder int            `gorm:"column:lesson_order;not null" json:"lesson_order"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	// 关联关系
	Course  *Courses `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Chapter *Chapters `gorm:"foreignKey:ChapterID" json:"chapter,omitempty"`
	Tasks   []Tasks  `gorm:"foreignKey:LessonID" json:"tasks,omitempty"`
}

// TableName 指定表名
func (Lessons) TableName() string {
	return "lessons"
}

