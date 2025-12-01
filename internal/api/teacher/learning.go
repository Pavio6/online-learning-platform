package teacher

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/service"
)

// LearningHandler 教师学习进度查询
type LearningHandler struct {
	learningService *service.LearningService
}

// NewLearningHandler 创建
func NewLearningHandler() *LearningHandler {
	return &LearningHandler{
		learningService: service.NewLearningService(),
	}
}

// ListCourseLearning 教师查看学生进度
// @Summary 查看课程学习进度
// @Tags 教师学习
// @Security BearerAuth
// @Produce json
// @Param id path int true "课程ID"
// @Success 200 {array} service.LearningProgressView
// @Router /api/v1/teacher/courses/{id}/learning [get]
func (h *LearningHandler) ListCourseLearning(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid course id",
		})
		return
	}

	instructorID, _ := c.Get("user_id")
	branchID, _ := c.Get("branch_id")

	progress, err := h.learningService.ListCourseProgressForTeacher(instructorID.(uint), branchID.(uint), uint(courseID))
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

	c.JSON(http.StatusOK, progress)
}
