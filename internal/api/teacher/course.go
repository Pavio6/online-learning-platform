package teacher

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/service"
)

// CourseHandler 教师课程管理处理器
type CourseHandler struct {
	courseService *service.CourseService
}

// NewCourseHandler 创建课程处理器
func NewCourseHandler() *CourseHandler {
	return &CourseHandler{
		courseService: service.NewCourseService(),
	}
}

// CreateCourse 创建课程
// @Summary 创建课程
// @Description 教师创建新课程
// @Tags 教师课程管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateCourseRequest true "课程信息"
// @Success 200 {object} models.Courses
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/teacher/courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	instructorID, _ := c.Get("user_id")

	var req service.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": err.Error(),
		})
		return
	}

	course, err := h.courseService.CreateCourse(instructorID.(uint), &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.HTTPStatus(), gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeInternal,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, course)
}

// ListCourses 获取我的课程列表
// @Summary 获取我的课程列表
// @Description 获取当前教师创建的所有课程
// @Tags 教师课程管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/teacher/courses [get]
func (h *CourseHandler) ListCourses(c *gin.Context) {
	instructorID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	instructorIDUint := instructorID.(uint)
	courses, total, err := h.courseService.ListCourses(&instructorIDUint, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeInternal,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
		"total":   total,
		"page":    page,
		"page_size": pageSize,
	})
}

// GetCourse 获取课程详情
// @Summary 获取课程详情
// @Description 获取课程详细信息，包含章节和课程
// @Tags 教师课程管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "课程ID"
// @Success 200 {object} service.CourseInfo
// @Router /api/v1/teacher/courses/:id [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid course id",
		})
		return
	}

	course, err := h.courseService.GetCourse(uint(courseID), true)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.HTTPStatus(), gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeInternal,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, course)
}

// CreateChapter 创建章节
// @Summary 创建章节
// @Description 为课程创建新章节
// @Tags 教师课程管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "课程ID"
// @Param request body service.CreateChapterRequest true "章节信息"
// @Success 200 {object} models.Chapters
// @Router /api/v1/teacher/courses/:id/chapters [post]
func (h *CourseHandler) CreateChapter(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid course id",
		})
		return
	}

	instructorID, _ := c.Get("user_id")

	var req service.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": err.Error(),
		})
		return
	}

	chapter, err := h.courseService.CreateChapter(uint(courseID), instructorID.(uint), &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.HTTPStatus(), gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeInternal,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, chapter)
}

// CreateLesson 创建课程（上传视频）
// @Summary 创建课程
// @Description 为章节创建新课程，可以上传视频文件
// @Tags 教师课程管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "课程ID"
// @Param chapter_id path int true "章节ID"
// @Param lesson_title formData string true "课程标题"
// @Param lesson_type formData string false "课程类型" default(video)
// @Param lesson_order formData int false "课程顺序"
// @Param video_file formData file false "视频文件"
// @Success 200 {object} models.Lessons
// @Router /api/v1/teacher/courses/:id/chapters/:chapter_id/lessons [post]
func (h *CourseHandler) CreateLesson(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid course id",
		})
		return
	}

	chapterID, err := strconv.ParseUint(c.Param("chapter_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid chapter id",
		})
		return
	}

	instructorID, _ := c.Get("user_id")

	// 从form-data获取数据
	req := service.CreateLessonRequest{
		LessonTitle: c.PostForm("lesson_title"),
		LessonType:  c.PostForm("lesson_type"),
	}

	if orderStr := c.PostForm("lesson_order"); orderStr != "" {
		if order, err := strconv.Atoi(orderStr); err == nil {
			req.LessonOrder = order
		}
	}

	if req.LessonTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "lesson_title is required",
		})
		return
	}

	// 处理文件上传
	var videoFile []byte
	var fileName string
	file, err := c.FormFile("video_file")
	if err == nil && file != nil {
		f, err := file.Open()
		if err == nil {
			defer f.Close()
			videoFile = make([]byte, file.Size)
			f.Read(videoFile)
			fileName = file.Filename
		}
	}

	lesson, err := h.courseService.CreateLesson(uint(courseID), uint(chapterID), instructorID.(uint), &req, videoFile, fileName)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.HTTPStatus(), gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeInternal,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

