package service

import (
	"fmt"

	"gorm.io/gorm"

	"online-learning-platform/internal/database"
	apperrors "online-learning-platform/internal/errors"
	"online-learning-platform/internal/models"
)

// ensureInstructorRecord 确保教师在中央服务器的instructors表中存在，并返回记录
func ensureInstructorRecord(instructorUserID, branchID uint) (*models.Instructors, error) {
	db := database.GetCentralDB()

	var instructor models.Instructors
	if err := db.Where("branch_id = ? AND branch_user_id = ?", branchID, instructorUserID).
		First(&instructor).Error; err == nil {
		return &instructor, nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to query instructor record: %w", err)
	}

	// 不存在，则从分支节点读取用户信息并创建
	branchDB, err := database.GetBranchDBByBranchID(branchID)
	if err != nil {
		return nil, err
	}

	var user models.Users
	if err := branchDB.Where("user_id = ?", instructorUserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch teacher info from branch: %w", err)
	}

	if user.Role != "teacher" {
		return nil, apperrors.ErrNotCourseInstructor
	}

	instructor = models.Instructors{
		BranchID:     branchID,
		BranchUserID: instructorUserID,
		Username:     user.Username,
		Email:        user.Email,
	}

	if err := db.Create(&instructor).Error; err != nil {
		return nil, fmt.Errorf("failed to create instructor record: %w", err)
	}

	return &instructor, nil
}

// getInstructorByID 根据中央服务器的instructor_id获取记录
func getInstructorByID(instructorID uint) (*models.Instructors, error) {
	db := database.GetCentralDB()
	var instructor models.Instructors
	if err := db.Where("instructor_id = ?", instructorID).First(&instructor).Error; err != nil {
		return nil, err
	}
	return &instructor, nil
}

func validateCourseOwner(courseID, instructorUserID, branchID uint) error {
	instructorRecord, err := ensureInstructorRecord(instructorUserID, branchID)
	if err != nil {
		return err
	}

	db := database.GetCentralDB()
	var course models.Courses
	if err := db.Where("course_id = ? AND instructor_id = ?", courseID, instructorRecord.InstructorID).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrNotCourseInstructor
		}
		return fmt.Errorf("failed to verify course owner: %w", err)
	}
	return nil
}
