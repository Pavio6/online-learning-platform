package models

import (
	"time"

	"gorm.io/gorm"
)

// Answers 答案表（分支节点）
type Answers struct {
	AnswerID      uint           `gorm:"primaryKey;column:answer_id" json:"answer_id"`
	TaskID        uint           `gorm:"column:task_id;not null;index" json:"task_id"`
	BranchID      uint           `gorm:"column:branch_id;not null;index" json:"branch_id"`
	UserID        uint           `gorm:"column:user_id;not null;index" json:"user_id"`
	GradedBy      *uint          `gorm:"column:graded_by;index" json:"graded_by"` // 批改老师的ID
	AnswerContent string         `gorm:"column:answer_content;type:text" json:"answer_content"` // 文本内容或图片URL
	Type          string         `gorm:"column:type;default:'text'" json:"type"` // text, image_url
	Score         int            `gorm:"column:score;default:0" json:"score"`
	IsGraded      bool           `gorm:"column:is_graded;default:false" json:"is_graded"`
	SubmittedAt   time.Time      `gorm:"column:submitted_at" json:"submitted_at"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	// 关联关系
	Branch   *Branches `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	User     *Users    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Grader   *Users    `gorm:"foreignKey:GradedBy" json:"grader,omitempty"` // 批改老师
}

// TableName 指定表名
func (Answers) TableName() string {
	return "answers"
}

