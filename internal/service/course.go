package service

import (
	"bytes"
	"context"
	"fmt"

	"gorm.io/gorm"

	"online-learning-platform/internal/database"
	apperrors "online-learning-platform/internal/errors"
	"online-learning-platform/internal/models"
	"online-learning-platform/internal/oss"
)

// CourseService 课程管理服务
type CourseService struct{}

// NewCourseService 创建课程服务
func NewCourseService() *CourseService {
	return &CourseService{}
}

// CreateCourseRequest 创建课程请求
type CreateCourseRequest struct {
	CourseTitle string  `json:"course_title" binding:"required"`
	Description string  `json:"description"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Status      string  `json:"status"`
}

// CreateChapterRequest 创建章节请求
type CreateChapterRequest struct {
	ChapterTitle string `json:"chapter_title" binding:"required"`
	Description  string `json:"description"`
	ChapterOrder int    `json:"chapter_order"`
}

// CreateLessonRequest 创建课程请求
type CreateLessonRequest struct {
	LessonTitle string `json:"lesson_title" binding:"required"`
	LessonType  string `json:"lesson_type"` // video, text, quiz
	LessonOrder int    `json:"lesson_order"`
	ContentURL  string `json:"content_url"` // OSS URL，如果上传文件则自动生成
}

// CourseInfo 课程信息
type CourseInfo struct {
	CourseID    uint            `json:"course_id"`
	CourseTitle string          `json:"course_title"`
	Description string          `json:"description"`
	InstructorID uint           `json:"instructor_id"`
	StartDate   *string         `json:"start_date"`
	EndDate     *string         `json:"end_date"`
	Status      string          `json:"status"`
	Chapters    []ChapterInfo   `json:"chapters,omitempty"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

// ChapterInfo 章节信息
type ChapterInfo struct {
	ChapterID    uint          `json:"chapter_id"`
	CourseID     uint          `json:"course_id"`
	ChapterTitle string        `json:"chapter_title"`
	ChapterOrder int           `json:"chapter_order"`
	Description  string        `json:"description"`
	Lessons      []LessonInfo  `json:"lessons,omitempty"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
}

// LessonInfo 课程信息
type LessonInfo struct {
	LessonID    uint   `json:"lesson_id"`
	CourseID    uint   `json:"course_id"`
	ChapterID   uint   `json:"chapter_id"`
	LessonTitle string `json:"lesson_title"`
	ContentURL  string `json:"content_url"`
	LessonType  string `json:"lesson_type"`
	LessonOrder int    `json:"lesson_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateCourse 教师创建课程
func (s *CourseService) CreateCourse(instructorID uint, req *CreateCourseRequest) (*models.Courses, error) {
	db := database.GetCentralDB()

	course := models.Courses{
		CourseTitle:  req.CourseTitle,
		Description:  req.Description,
		InstructorID: instructorID,
		Status:       "active",
	}

	if req.Status != "" {
		course.Status = req.Status
	}

	if err := db.Create(&course).Error; err != nil {
		return nil, fmt.Errorf("failed to create course: %w", err)
	}

	return &course, nil
}

// CreateChapter 创建章节
func (s *CourseService) CreateChapter(courseID, instructorID uint, req *CreateChapterRequest) (*models.Chapters, error) {
	db := database.GetCentralDB()

	// 验证课程是否存在且属于该教师
	var course models.Courses
	if err := db.Where("course_id = ? AND instructor_id = ?", courseID, instructorID).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotCourseInstructor
		}
		return nil, fmt.Errorf("failed to verify course: %w", err)
	}

	// 如果没有指定顺序，自动获取下一个顺序
	chapterOrder := req.ChapterOrder
	if chapterOrder == 0 {
		var maxOrder int
		db.Model(&models.Chapters{}).
			Where("course_id = ?", courseID).
			Select("COALESCE(MAX(chapter_order), 0)").
			Scan(&maxOrder)
		chapterOrder = maxOrder + 1
	}

	chapter := models.Chapters{
		CourseID:     courseID,
		ChapterTitle: req.ChapterTitle,
		ChapterOrder: chapterOrder,
		Description:  req.Description,
	}

	if err := db.Create(&chapter).Error; err != nil {
		return nil, fmt.Errorf("failed to create chapter: %w", err)
	}

	return &chapter, nil
}

// CreateLesson 创建课程（上传视频到OSS）
func (s *CourseService) CreateLesson(courseID, chapterID, instructorID uint, req *CreateLessonRequest, videoFile []byte, fileName string) (*models.Lessons, error) {
	db := database.GetCentralDB()

	// 验证课程是否存在且属于该教师
	var course models.Courses
	if err := db.Where("course_id = ? AND instructor_id = ?", courseID, instructorID).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotCourseInstructor
		}
		return nil, fmt.Errorf("failed to verify course: %w", err)
	}

	// 验证章节是否存在
	var chapter models.Chapters
	if err := db.Where("chapter_id = ? AND course_id = ?", chapterID, courseID).First(&chapter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrChapterNotFound
		}
		return nil, fmt.Errorf("failed to verify chapter: %w", err)
	}

	// 上传视频到OSS
	contentURL := req.ContentURL
	if len(videoFile) > 0 && fileName != "" {
		objectKey := fmt.Sprintf("courses/%d/chapters/%d/lessons/%s", courseID, chapterID, fileName)
		ossURL, err := oss.UploadReader(context.Background(), objectKey, bytes.NewReader(videoFile))
		if err != nil {
			return nil, fmt.Errorf("failed to upload video to OSS: %w", err)
		}
		contentURL = ossURL
	}

	// 如果没有指定顺序，自动获取下一个顺序
	lessonOrder := req.LessonOrder
	if lessonOrder == 0 {
		var maxOrder int
		db.Model(&models.Lessons{}).
			Where("chapter_id = ?", chapterID).
			Select("COALESCE(MAX(lesson_order), 0)").
			Scan(&maxOrder)
		lessonOrder = maxOrder + 1
	}

	lessonType := req.LessonType
	if lessonType == "" {
		lessonType = "video"
	}

	lesson := models.Lessons{
		CourseID:    courseID,
		ChapterID:   chapterID,
		LessonTitle: req.LessonTitle,
		ContentURL:  contentURL,
		LessonType:  lessonType,
		LessonOrder: lessonOrder,
	}

	if err := db.Create(&lesson).Error; err != nil {
		return nil, fmt.Errorf("failed to create lesson: %w", err)
	}

	return &lesson, nil
}

// GetCourse 获取课程详情
func (s *CourseService) GetCourse(courseID uint, includeDetails bool) (*CourseInfo, error) {
	db := database.GetCentralDB()

	var course models.Courses
	if err := db.Where("course_id = ?", courseID).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrCourseNotFound
		}
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	courseInfo := &CourseInfo{
		CourseID:     course.CourseID,
		CourseTitle:  course.CourseTitle,
		Description:  course.Description,
		InstructorID: course.InstructorID,
		Status:       course.Status,
	}

	if course.StartDate != nil {
		startDate := course.StartDate.Format("2006-01-02 15:04:05")
		courseInfo.StartDate = &startDate
	}
	if course.EndDate != nil {
		endDate := course.EndDate.Format("2006-01-02 15:04:05")
		courseInfo.EndDate = &endDate
	}
	courseInfo.CreatedAt = course.CreatedAt.Format("2006-01-02 15:04:05")
	courseInfo.UpdatedAt = course.UpdatedAt.Format("2006-01-02 15:04:05")

	if includeDetails {
		// 获取章节和课程
		var chapters []models.Chapters
		db.Where("course_id = ?", courseID).Order("chapter_order ASC").Find(&chapters)

		chapterInfos := make([]ChapterInfo, 0, len(chapters))
		for _, ch := range chapters {
			chapterInfo := ChapterInfo{
				ChapterID:    ch.ChapterID,
				CourseID:     ch.CourseID,
				ChapterTitle: ch.ChapterTitle,
				ChapterOrder: ch.ChapterOrder,
				Description:  ch.Description,
				CreatedAt:    ch.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt:    ch.UpdatedAt.Format("2006-01-02 15:04:05"),
			}

			// 获取课程
			var lessons []models.Lessons
			db.Where("chapter_id = ?", ch.ChapterID).Order("lesson_order ASC").Find(&lessons)

			lessonInfos := make([]LessonInfo, 0, len(lessons))
			for _, le := range lessons {
				lessonInfos = append(lessonInfos, LessonInfo{
					LessonID:    le.LessonID,
					CourseID:    le.CourseID,
					ChapterID:   le.ChapterID,
					LessonTitle: le.LessonTitle,
					ContentURL:  le.ContentURL,
					LessonType:  le.LessonType,
					LessonOrder: le.LessonOrder,
					CreatedAt:   le.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdatedAt:   le.UpdatedAt.Format("2006-01-02 15:04:05"),
				})
			}
			chapterInfo.Lessons = lessonInfos
			chapterInfos = append(chapterInfos, chapterInfo)
		}
		courseInfo.Chapters = chapterInfos
	}

	return courseInfo, nil
}

// ListCourses 获取课程列表
func (s *CourseService) ListCourses(instructorID *uint, page, pageSize int) ([]CourseInfo, int64, error) {
	db := database.GetCentralDB()

	query := db.Model(&models.Courses{})
	if instructorID != nil {
		query = query.Where("instructor_id = ?", *instructorID)
	}

	var total int64
	query.Count(&total)

	var courses []models.Courses
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list courses: %w", err)
	}

	courseInfos := make([]CourseInfo, 0, len(courses))
	for _, course := range courses {
		courseInfo := CourseInfo{
			CourseID:     course.CourseID,
			CourseTitle:  course.CourseTitle,
			Description:  course.Description,
			InstructorID: course.InstructorID,
			Status:       course.Status,
			CreatedAt:    course.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    course.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if course.StartDate != nil {
			startDate := course.StartDate.Format("2006-01-02 15:04:05")
			courseInfo.StartDate = &startDate
		}
		if course.EndDate != nil {
			endDate := course.EndDate.Format("2006-01-02 15:04:05")
			courseInfo.EndDate = &endDate
		}
		courseInfos = append(courseInfos, courseInfo)
	}

	return courseInfos, total, nil
}

