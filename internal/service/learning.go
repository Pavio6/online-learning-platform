package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"online-learning-platform/internal/database"
	apperrors "online-learning-platform/internal/errors"
	"online-learning-platform/internal/models"
)

// LearningService 学习进度服务
type LearningService struct{}

// NewLearningService 创建实例
func NewLearningService() *LearningService {
	return &LearningService{}
}

// EnrollCourse 学生报名课程
func (s *LearningService) EnrollCourse(userID, branchID, courseID uint) (*models.Learning, error) {
	// 校验课程是否存在
	if err := ensureCourseExists(courseID); err != nil {
		return nil, err
	}

	branchDB, err := database.GetBranchDBByBranchID(branchID)
	if err != nil {
		return nil, err
	}

	var learning models.Learning
	if err := branchDB.Where("user_id = ? AND course_id = ?", userID, courseID).First(&learning).Error; err == nil {
		return nil, apperrors.ErrAlreadyEnrolled
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to query learning record: %w", err)
	}

	learning = models.Learning{
		UserID:             userID,
		CourseID:           courseID,
		Status:             "enrolled",
		ProgressPercentage: 0,
	}

	if err := branchDB.Create(&learning).Error; err != nil {
		return nil, fmt.Errorf("failed to enroll course: %w", err)
	}

	return &learning, nil
}

// UpdateProgressRequest 更新进度请求
type UpdateProgressRequest struct {
	Status             string `json:"status"`
	ProgressPercentage *int   `json:"progress_percentage"`
}

// UpdateProgress 更新学习进度
func (s *LearningService) UpdateProgress(userID, branchID, courseID uint, req *UpdateProgressRequest) (*models.Learning, error) {
	branchDB, err := database.GetBranchDBByBranchID(branchID)
	if err != nil {
		return nil, err
	}

	var learning models.Learning
	if err := branchDB.Where("user_id = ? AND course_id = ?", userID, courseID).First(&learning).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotEnrolled
		}
		return nil, fmt.Errorf("failed to get learning progress: %w", err)
	}

	if req.ProgressPercentage != nil {
		if *req.ProgressPercentage < 0 || *req.ProgressPercentage > 100 {
			return nil, apperrors.ErrInvalidParam
		}
		learning.ProgressPercentage = *req.ProgressPercentage
	}

	if req.Status != "" {
		learning.Status = req.Status
		if req.Status == "completed" && learning.CompletedAt == nil {
			now := time.Now()
			learning.CompletedAt = &now
		}
	}

	if err := branchDB.Save(&learning).Error; err != nil {
		return nil, fmt.Errorf("failed to update learning progress: %w", err)
	}

	return &learning, nil
}

// GetStudentProgress 获取学生进度
func (s *LearningService) GetStudentProgress(userID, branchID, courseID uint) (*models.Learning, error) {
	branchDB, err := database.GetBranchDBByBranchID(branchID)
	if err != nil {
		return nil, err
	}

	var learning models.Learning
	if err := branchDB.Where("user_id = ? AND course_id = ?", userID, courseID).First(&learning).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotEnrolled
		}
		return nil, fmt.Errorf("failed to query learning progress: %w", err)
	}

	return &learning, nil
}

// LearningProgressView 教师视角学生学习进度
type LearningProgressView struct {
	BranchID           uint       `json:"branch_id"`
	UserID             uint       `json:"user_id"`
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	CourseID           uint       `json:"course_id"`
	Status             string     `json:"status"`
	ProgressPercentage int        `json:"progress_percentage"`
	CompletedAt        *time.Time `json:"completed_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// ListCourseProgressForTeacher 教师查看课程学生进度
func (s *LearningService) ListCourseProgressForTeacher(instructorUserID, branchID, courseID uint) ([]LearningProgressView, error) {
	if err := validateCourseOwner(courseID, instructorUserID, branchID); err != nil {
		return nil, err
	}

	branchDBs := database.GetAllBranchDBs()
	results := make([]LearningProgressView, 0)

	for bID, db := range branchDBs {
		type row struct {
			UserID             uint
			Username           string
			Email              string
			CourseID           uint
			Status             string
			ProgressPercentage int
			CompletedAt        *time.Time
			UpdatedAt          time.Time
		}
		var rows []row
		if err := db.Table("learning").
			Select("learning.user_id, learning.course_id, learning.status, learning.progress_percentage, learning.completed_at, learning.updated_at, users.username, users.email").
			Joins("JOIN users ON users.user_id = learning.user_id").
			Where("learning.course_id = ?", courseID).
			Scan(&rows).Error; err != nil {
			continue
		}

		for _, r := range rows {
			results = append(results, LearningProgressView{
				BranchID:           bID,
				UserID:             r.UserID,
				Username:           r.Username,
				Email:              r.Email,
				CourseID:           r.CourseID,
				Status:             r.Status,
				ProgressPercentage: r.ProgressPercentage,
				CompletedAt:        r.CompletedAt,
				UpdatedAt:          r.UpdatedAt,
			})
		}
	}

	return results, nil
}

func ensureCourseExists(courseID uint) error {
	db := database.GetCentralDB()
	var course models.Courses
	if err := db.Where("course_id = ?", courseID).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrCourseNotFound
		}
		return fmt.Errorf("failed to query course: %w", err)
	}
	return nil
}
