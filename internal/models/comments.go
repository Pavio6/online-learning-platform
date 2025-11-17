package models

import (
	"time"

	"gorm.io/gorm"
)

// Comments 评论表（分支节点）
type Comments struct {
	CommentID      uint           `gorm:"primaryKey;column:comment_id" json:"comment_id"`
	CourseID       uint           `gorm:"column:course_id;not null;index" json:"course_id"`
	UserID         uint           `gorm:"column:user_id;not null;index" json:"user_id"`
	BranchID       uint           `gorm:"column:branch_id;not null;index" json:"branch_id"`
	CommentContent string         `gorm:"column:comment_content;type:text;not null" json:"comment_content"`
	ParentCommentID *uint         `gorm:"column:parent_comment_id;index" json:"parent_comment_id"` // 回复的评论ID
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	// 关联关系
	Branch        *Branches  `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	User          *Users     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentComment *Comments `gorm:"foreignKey:ParentCommentID" json:"parent_comment,omitempty"`
	Replies       []Comments `gorm:"foreignKey:ParentCommentID" json:"replies,omitempty"`
}

// TableName 指定表名
func (Comments) TableName() string {
	return "comments"
}

