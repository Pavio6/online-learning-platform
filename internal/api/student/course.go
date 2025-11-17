package student

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/service"
)

// CourseHandler 学生课程处理器
type CourseHandler struct {
	courseService *service.CourseService
}

// NewCourseHandler 创建课程处理器
func NewCourseHandler() *CourseHandler {
	return &CourseHandler{
		courseService: service.NewCourseService(),
	}
}

// ListCourses 获取课程列表
// @Summary 获取课程列表
// @Description 获取所有可用课程列表
// @Tags 学生课程
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/student/courses [get]
func (h *CourseHandler) ListCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	courses, total, err := h.courseService.ListCourses(nil, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeInternal,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses":   courses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetCourse 获取课程详情
// @Summary 获取课程详情
// @Description 获取课程详细信息，包含章节和课程
// @Tags 学生课程
// @Accept json
// @Produce json
// @Param id path int true "课程ID"
// @Success 200 {object} service.CourseInfo
// @Router /api/v1/student/courses/:id [get]
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

