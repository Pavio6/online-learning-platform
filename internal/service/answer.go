package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"online-learning-platform/internal/database"
	apperrors "online-learning-platform/internal/errors"
	"online-learning-platform/internal/models"
	"online-learning-platform/internal/oss"
)

// AnswerService 作业提交/评分服务
type AnswerService struct{}

// NewAnswerService 创建实例
func NewAnswerService() *AnswerService {
	return &AnswerService{}
}

// SubmitAnswerRequest 学生提交作业请求
type SubmitAnswerRequest struct {
	AnswerContent string `json:"answer_content"`
	Type          string `json:"type"` // text, image_url
}

// GradeAnswerRequest 教师评分请求
type GradeAnswerRequest struct {
	BranchID uint `json:"branch_id" binding:"required"` // 答案所在的分支ID
	Score    int  `json:"score" binding:"required"`
}

// SubmitAnswer 学生提交作业
func (s *AnswerService) SubmitAnswer(userID, branchID, taskID uint, req *SubmitAnswerRequest, fileBytes []byte, fileName string) (*models.Answers, error) {
	// 校验任务是否存在
	centralDB := database.GetCentralDB()
	var task models.Tasks
	if err := centralDB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to verify task: %w", err)
	}

	branchDB, err := database.GetBranchDBByBranchID(branchID)
	if err != nil {
		return nil, err
	}

	if req.AnswerContent == "" && len(fileBytes) == 0 {
		return nil, apperrors.ErrInvalidParam
	}

	answerContent := req.AnswerContent
	answerType := req.Type
	objectKey := ""

	// 如果上传了图片则优先使用OSS
	if len(fileBytes) > 0 && fileName != "" {
		objectKey = fmt.Sprintf("answers/%d/%d/%d/%d_%s", branchID, userID, taskID, time.Now().Unix(), fileName)
		ossURL, err := oss.UploadReader(context.Background(), objectKey, bytes.NewReader(fileBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to upload answer file to OSS: %w", err)
		}
		answerContent = ossURL
		answerType = "image_url"
	}

	if answerType == "" {
		answerType = "text"
	}

	now := time.Now()

	var answer models.Answers
	tx := branchDB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&answer)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		if objectKey != "" {
			_ = oss.DeleteObject(objectKey)
		}
		return nil, fmt.Errorf("failed to check existing answer: %w", tx.Error)
	}

	if tx.Error == nil {
		// 更新已有作业
		answer.AnswerContent = answerContent
		answer.Type = answerType
		answer.SubmittedAt = now
		answer.Score = 0
		answer.IsGraded = false
		answer.GradedBy = nil
		answer.BranchID = branchID

		if err := branchDB.Save(&answer).Error; err != nil {
			if objectKey != "" {
				_ = oss.DeleteObject(objectKey)
			}
			return nil, fmt.Errorf("failed to update answer: %w", err)
		}
		return &answer, nil
	}

	// 创建新作业
	answer = models.Answers{
		TaskID:        taskID,
		BranchID:      branchID,
		UserID:        userID,
		AnswerContent: answerContent,
		Type:          answerType,
		SubmittedAt:   now,
	}

	if err := branchDB.Create(&answer).Error; err != nil {
		if objectKey != "" {
			_ = oss.DeleteObject(objectKey)
		}
		return nil, fmt.Errorf("failed to create answer: %w", err)
	}

	return &answer, nil
}

// GetStudentAnswer 学生查询自己的作业
func (s *AnswerService) GetStudentAnswer(userID, branchID, taskID uint) (*models.Answers, error) {
	branchDB, err := database.GetBranchDBByBranchID(branchID)
	if err != nil {
		return nil, err
	}

	var answer models.Answers
	if err := branchDB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&answer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrAnswerNotFound
		}
		return nil, fmt.Errorf("failed to query answer: %w", err)
	}
	return &answer, nil
}

// AnswerWithStudentInfo 包含学生信息的作业
type AnswerWithStudentInfo struct {
	models.Answers
	StudentFirstName string `json:"student_first_name"`
	StudentLastName  string `json:"student_last_name"`
}

// ListAnswersForTask 教师查看任务的所有作业（跨所有分支，因为课程是共享的）
func (s *AnswerService) ListAnswersForTask(instructorUserID, branchID, taskID uint) ([]AnswerWithStudentInfo, error) {
	if _, err := ensureInstructorRecord(instructorUserID, branchID); err != nil {
		return nil, err
	}

	// 校验任务属于该教师
	if err := validateTaskOwner(taskID, instructorUserID, branchID); err != nil {
		return nil, err
	}

	// 查询所有分支的作业（因为课程是共享的，教师应该能看到所有报名学生的作业）
	branchDBs := database.GetAllBranchDBs()
	allAnswers := make([]AnswerWithStudentInfo, 0)

	// 使用JOIN查询获取学生姓名
	type result struct {
		models.Answers
		FirstName string
		LastName  string
	}

	for _, branchDB := range branchDBs {
		var results []result
		if err := branchDB.Table("answers").
			Select("answers.*, users.first_name, users.last_name").
			Joins("JOIN users ON users.user_id = answers.user_id").
			Where("answers.task_id = ?", taskID).
			Order("answers.submitted_at DESC").
			Scan(&results).Error; err != nil {
			// 如果某个分支查询失败，继续查询其他分支
			continue
		}

		// 转换为AnswerWithStudentInfo
		for _, r := range results {
			allAnswers = append(allAnswers, AnswerWithStudentInfo{
				Answers:          r.Answers,
				StudentFirstName: r.FirstName,
				StudentLastName:  r.LastName,
			})
		}
	}

	return allAnswers, nil
}

// GradeAnswer 教师评分
// answerBranchID: 答案所在的分支ID（从请求中获取，确保找到正确的答案）
func (s *AnswerService) GradeAnswer(instructorUserID, instructorBranchID, answerID, answerBranchID uint, score int) (*models.Answers, error) {
	if _, err := ensureInstructorRecord(instructorUserID, instructorBranchID); err != nil {
		return nil, err
	}

	// 根据答案所在的分支ID，直接查询对应的分支数据库
	// 这样可以避免因为不同分支中 answer_id 重复而找到错误的答案
	branchDB, err := database.GetBranchDBByBranchID(answerBranchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch database: %w", err)
	}

	var answer models.Answers
	// 同时匹配 answer_id 和 branch_id，确保找到正确的答案
	if err := branchDB.Where("answer_id = ? AND branch_id = ?", answerID, answerBranchID).First(&answer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrAnswerNotFound
		}
		return nil, fmt.Errorf("failed to get answer: %w", err)
	}

	// 更新答案
	answer.Score = score
	answer.IsGraded = true
	
	// 注意：graded_by 字段有外键约束，如果教师在其他分支，该分支的users表中可能没有该教师
	// 因此，我们需要检查该分支是否有该教师，如果没有则不设置 graded_by
	// 或者，我们可以尝试在该分支中查找或创建教师记录
	// 但最简单的方法是：如果教师在该分支存在，则设置 graded_by；否则设置为 nil
	var teacherUser models.Users
	if err := branchDB.Where("user_id = ? AND branch_id = ?", instructorUserID, answerBranchID).First(&teacherUser).Error; err == nil {
		// 教师在该分支存在，可以设置 graded_by
		answer.GradedBy = &instructorUserID
	} else {
		// 教师不在该分支，不设置 graded_by（避免外键约束错误）
		// 但我们可以通过其他方式记录评分教师（比如在注释中，或者不记录）
		answer.GradedBy = nil
	}

	if err := branchDB.Save(&answer).Error; err != nil {
		return nil, fmt.Errorf("failed to grade answer: %w", err)
	}

	return &answer, nil
}

// validateTaskOwner 确认任务属于当前教师
func validateTaskOwner(taskID, instructorUserID, branchID uint) error {
	instructorRecord, err := ensureInstructorRecord(instructorUserID, branchID)
	if err != nil {
		return err
	}

	db := database.GetCentralDB()
	type result struct {
		TaskID       uint
		InstructorID uint
	}
	var res result

	query := `SELECT tasks.task_id, courses.instructor_id
			  FROM tasks
			  JOIN lessons ON lessons.lesson_id = tasks.lesson_id
			  JOIN courses ON courses.course_id = lessons.course_id
			  WHERE tasks.task_id = ?`

	if err := db.Raw(query, taskID).Scan(&res).Error; err != nil {
		if err == gorm.ErrRecordNotFound || res.TaskID == 0 {
			return apperrors.ErrTaskNotFound
		}
		return fmt.Errorf("failed to verify task owner: %w", err)
	}

	if res.InstructorID != instructorRecord.InstructorID {
		return apperrors.ErrNotCourseInstructor
	}
	return nil
}
