package student

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/service"
)

// LearningHandler 学习进度处理器
type LearningHandler struct {
	learningService *service.LearningService
}

// NewLearningHandler 创建学习进度处理器
func NewLearningHandler() *LearningHandler {
	return &LearningHandler{
		learningService: service.NewLearningService(),
	}
}

// Enroll 报名课程
// @Summary 学生报名课程
// @Tags 学生学习
// @Security BearerAuth
// @Param id path int true "课程ID"
// @Success 200 {object} models.Learning
// @Router /api/v1/student/courses/{id}/enroll [post]
func (h *LearningHandler) Enroll(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid course id",
		})
		return
	}

	userID, _ := c.Get("user_id")
	branchID, _ := c.Get("branch_id")

	learning, err := h.learningService.EnrollCourse(userID.(uint), branchID.(uint), uint(courseID))
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

	c.JSON(http.StatusOK, learning)
}

// UpdateProgress 更新学习进度
// @Summary 更新学习进度
// @Tags 学生学习
// @Security BearerAuth
// @Accept json
// @Param id path int true "课程ID"
// @Param request body service.UpdateProgressRequest true "学习进度"
// @Success 200 {object} models.Learning
// @Router /api/v1/student/courses/{id}/progress [put]
func (h *LearningHandler) UpdateProgress(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid course id",
		})
		return
	}

	var req service.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	branchID, _ := c.Get("branch_id")

	learning, err := h.learningService.UpdateProgress(userID.(uint), branchID.(uint), uint(courseID), &req)
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

	c.JSON(http.StatusOK, learning)
}

// GetProgress 学生查询自己进度
// @Summary 获取学习进度
// @Description 获取学生学习进度，如果未报名则返回 enrolled: false
// @Tags 学生学习
// @Security BearerAuth
// @Param id path int true "课程ID"
// @Success 200 {object} models.Learning "已报名时返回学习进度"
// @Success 200 {object} map[string]bool "未报名时返回 {\"enrolled\": false}"
// @Router /api/v1/student/courses/{id}/progress [get]
func (h *LearningHandler) GetProgress(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid course id",
		})
		return
	}

	userID, _ := c.Get("user_id")
	branchID, _ := c.Get("branch_id")

	learning, err := h.learningService.GetStudentProgress(userID.(uint), branchID.(uint), uint(courseID))
	if err != nil {
		// 如果是"未报名"错误，返回正常响应表示未报名
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.ErrCodeNotEnrolled {
			c.JSON(http.StatusOK, gin.H{
				"enrolled": false,
			})
			return
		}
		// 其他错误才返回错误响应
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

	// 已报名，返回学习进度
	c.JSON(http.StatusOK, learning)
}
