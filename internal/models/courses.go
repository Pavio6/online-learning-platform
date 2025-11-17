package models

import (
	"time"

	"gorm.io/gorm"
)

// Courses 课程表（中央服务器）
type Courses struct {
	CourseID    uint           `gorm:"primaryKey;column:course_id" json:"course_id"`
	CourseTitle string         `gorm:"column:course_title;not null" json:"course_title"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	InstructorID uint           `gorm:"column:instructor_id;not null;index" json:"instructor_id"`
	StartDate   *time.Time     `gorm:"column:start_date" json:"start_date"`
	EndDate     *time.Time     `gorm:"column:end_date" json:"end_date"`
	Status      string         `gorm:"column:status;default:'active'" json:"status"` // active, archived
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	// 关联关系（这些关联在分支节点查询时需要从中央服务器获取）
	Chapters []Chapters `gorm:"foreignKey:CourseID" json:"chapters,omitempty"`
	Lessons  []Lessons  `gorm:"foreignKey:CourseID" json:"lessons,omitempty"`
}

// TableName 指定表名
func (Courses) TableName() string {
	return "courses"
}

