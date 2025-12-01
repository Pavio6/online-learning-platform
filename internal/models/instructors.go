package models

import (
	"time"
)

// Instructors 中央服务器教师表
type Instructors struct {
	InstructorID uint      `gorm:"primaryKey;column:instructor_id" json:"instructor_id"`
	BranchID     uint      `gorm:"column:branch_id;not null" json:"branch_id"`
	BranchUserID uint      `gorm:"column:branch_user_id;not null" json:"branch_user_id"`
	Username     string    `gorm:"column:username;not null" json:"username"`
	Email        string    `gorm:"column:email;not null" json:"email"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Instructors) TableName() string {
	return "instructors"
}
