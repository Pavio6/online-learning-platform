package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"online-learning-platform/internal/database"
	apperrors "online-learning-platform/internal/errors"
	"online-learning-platform/internal/models"
)

// CommentService 评论服务
type CommentService struct{}

// NewCommentService 创建实例
func NewCommentService() *CommentService {
	return &CommentService{}
}

// AddCommentRequest 评论请求
type AddCommentRequest struct {
	Content         string `json:"content" binding:"required"`
	ParentCommentID *uint  `json:"parent_comment_id"`
}

// CommentView 评论返回结构
type CommentView struct {
	CommentID       uint      `json:"comment_id"`
	CourseID        uint      `json:"course_id"`
	UserID          uint      `json:"user_id"`
	BranchID        uint      `json:"branch_id"`
	Username        string    `json:"username"`
	CommentContent  string    `json:"comment_content"`
	ParentCommentID *uint     `json:"parent_comment_id"`
	CreatedAt       time.Time `json:"created_at"`
}

// AddComment 学生发表评论（需要验证是否报名）
func (s *CommentService) AddComment(userID, branchID, courseID uint, req *AddCommentRequest) (*models.Comments, error) {
	if err := ensureCourseExists(courseID); err != nil {
		return nil, err
	}

	if err := ensureStudentEnrolled(userID, branchID, courseID); err != nil {
		return nil, err
	}

	return s.createComment(userID, branchID, courseID, req)
}

// AddCommentAsTeacher 教师发表评论（需要验证是否是课程教师）
func (s *CommentService) AddCommentAsTeacher(instructorUserID, branchID, courseID uint, req *AddCommentRequest) (*models.Comments, error) {
	if err := ensureCourseExists(courseID); err != nil {
		return nil, err
	}

	// 验证教师是否是课程创建者
	if err := validateCourseOwner(courseID, instructorUserID, branchID); err != nil {
		return nil, err
	}

	return s.createComment(instructorUserID, branchID, courseID, req)
}

// createComment 创建评论的通用方法
func (s *CommentService) createComment(userID, branchID, courseID uint, req *AddCommentRequest) (*models.Comments, error) {
	branchDB, err := database.GetBranchDBByBranchID(branchID)
	if err != nil {
		return nil, err
	}

	parentCommentID := req.ParentCommentID
	if parentCommentID != nil {
		if *parentCommentID == 0 {
			parentCommentID = nil
		} else {
			// 验证父评论是否存在（需要跨分片查询，因为父评论可能在任意分支）
			var parentFound bool
			branchDBs := database.GetAllBranchDBs()
			for _, db := range branchDBs {
				var parent models.Comments
				if err := db.Where("comment_id = ? AND course_id = ?", *parentCommentID, courseID).First(&parent).Error; err == nil {
					parentFound = true
					break
				}
			}
			if !parentFound {
				return nil, apperrors.ErrNotFound
			}
		}
	}

	comment := models.Comments{
		CourseID:        courseID,
		UserID:          userID,
		BranchID:        branchID,
		CommentContent:  req.Content,
		ParentCommentID: parentCommentID,
	}

	if err := branchDB.Create(&comment).Error; err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}
	return &comment, nil
}

// ListComments 跨分片查询所有评论
func (s *CommentService) ListComments(courseID uint) ([]CommentView, error) {
	if err := ensureCourseExists(courseID); err != nil {
		return nil, err
	}

	branchDBs := database.GetAllBranchDBs()
	results := make([]CommentView, 0)

	for branchID, db := range branchDBs {
		type row struct {
			CommentID       uint
			CourseID        uint
			UserID          uint
			CommentContent  string
			ParentCommentID *uint
			CreatedAt       time.Time
			Username        string
		}
		var rows []row
		if err := db.Table("comments").
			Select("comments.comment_id, comments.course_id, comments.user_id, comments.comment_content, comments.parent_comment_id, comments.created_at, users.username").
			Joins("JOIN users ON users.user_id = comments.user_id").
			Where("comments.course_id = ?", courseID).
			Order("comments.created_at DESC").
			Scan(&rows).Error; err != nil {
			continue
		}

		for _, r := range rows {
			results = append(results, CommentView{
				CommentID:       r.CommentID,
				CourseID:        r.CourseID,
				UserID:          r.UserID,
				BranchID:        branchID,
				Username:        r.Username,
				CommentContent:  r.CommentContent,
				ParentCommentID: r.ParentCommentID,
				CreatedAt:       r.CreatedAt,
			})
		}
	}

	return results, nil
}

func ensureStudentEnrolled(userID, branchID, courseID uint) error {
	branchDB, err := database.GetBranchDBByBranchID(branchID)
	if err != nil {
		return err
	}

	var learning models.Learning
	if err := branchDB.Where("user_id = ? AND course_id = ?", userID, courseID).First(&learning).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrCannotComment
		}
		return fmt.Errorf("failed to verify enrollment: %w", err)
	}
	return nil
}
