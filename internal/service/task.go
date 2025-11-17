package service

import (
	"fmt"

	"gorm.io/gorm"

	"online-learning-platform/internal/database"
	apperrors "online-learning-platform/internal/errors"
	"online-learning-platform/internal/models"
)

// TaskService 任务管理服务
type TaskService struct{}

// NewTaskService 创建任务服务
func NewTaskService() *TaskService {
	return &TaskService{}
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	TaskTitle   string `json:"task_title" binding:"required"`
	Description string `json:"description"`
	TaskType    string `json:"task_type"` // essay, quiz, upload
	MaxScore    int    `json:"max_score"`
}

// TaskInfo 任务信息
type TaskInfo struct {
	TaskID      uint   `json:"task_id"`
	LessonID    uint   `json:"lesson_id"`
	TaskTitle   string `json:"task_title"`
	Description string `json:"description"`
	TaskType    string `json:"task_type"`
	MaxScore    int    `json:"max_score"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateTask 教师创建任务
func (s *TaskService) CreateTask(lessonID, instructorID uint, req *CreateTaskRequest) (*models.Tasks, error) {
	db := database.GetCentralDB()

	// 验证课程是否存在且属于该教师
	var lesson models.Lessons
	if err := db.Where("lesson_id = ?", lessonID).First(&lesson).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrLessonNotFound
		}
		return nil, fmt.Errorf("failed to verify lesson: %w", err)
	}

	// 验证课程是否属于该教师
	var course models.Courses
	if err := db.Where("course_id = ? AND instructor_id = ?", lesson.CourseID, instructorID).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotCourseInstructor
		}
		return nil, fmt.Errorf("failed to verify course: %w", err)
	}

	taskType := req.TaskType
	if taskType == "" {
		taskType = "essay"
	}

	maxScore := req.MaxScore
	if maxScore == 0 {
		maxScore = 100
	}

	task := models.Tasks{
		LessonID:    lessonID,
		TaskTitle:   req.TaskTitle,
		Description: req.Description,
		TaskType:    taskType,
		MaxScore:    maxScore,
	}

	if err := db.Create(&task).Error; err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &task, nil
}

// GetTask 获取任务详情
func (s *TaskService) GetTask(taskID uint) (*TaskInfo, error) {
	db := database.GetCentralDB()

	var task models.Tasks
	if err := db.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &TaskInfo{
		TaskID:      task.TaskID,
		LessonID:    task.LessonID,
		TaskTitle:   task.TaskTitle,
		Description: task.Description,
		TaskType:    task.TaskType,
		MaxScore:    task.MaxScore,
		CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ListTasksByCourse 获取课程的所有任务
func (s *TaskService) ListTasksByCourse(courseID uint) ([]TaskInfo, error) {
	db := database.GetCentralDB()

	// 通过lessons表关联查询
	var tasks []models.Tasks
	if err := db.Joins("JOIN lessons ON tasks.lesson_id = lessons.lesson_id").
		Where("lessons.course_id = ?", courseID).
		Order("tasks.created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	taskInfos := make([]TaskInfo, 0, len(tasks))
	for _, task := range tasks {
		taskInfos = append(taskInfos, TaskInfo{
			TaskID:      task.TaskID,
			LessonID:    task.LessonID,
			TaskTitle:   task.TaskTitle,
			Description: task.Description,
			TaskType:    task.TaskType,
			MaxScore:    task.MaxScore,
			CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   task.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return taskInfos, nil
}

// ListTasksByLesson 获取课程的所有任务
func (s *TaskService) ListTasksByLesson(lessonID uint) ([]TaskInfo, error) {
	db := database.GetCentralDB()

	var tasks []models.Tasks
	if err := db.Where("lesson_id = ?", lessonID).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	taskInfos := make([]TaskInfo, 0, len(tasks))
	for _, task := range tasks {
		taskInfos = append(taskInfos, TaskInfo{
			TaskID:      task.TaskID,
			LessonID:    task.LessonID,
			TaskTitle:   task.TaskTitle,
			Description: task.Description,
			TaskType:    task.TaskType,
			MaxScore:    task.MaxScore,
			CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   task.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return taskInfos, nil
}

