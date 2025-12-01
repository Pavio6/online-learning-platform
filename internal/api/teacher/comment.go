package teacher

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/models"
	"online-learning-platform/internal/service"
)

var _ = models.Comments{}

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler 创建
func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		commentService: service.NewCommentService(),
	}
}

// AddComment 教师发表评论
// @Summary 教师发表评论
// @Tags 教师评论
// @Security BearerAuth
// @Accept json
// @Param id path int true "课程ID"
// @Param request body service.AddCommentRequest true "评论"
// @Success 200 {object} models.Comments
// @Router /api/v1/teacher/courses/{id}/comments [post]
func (h *CommentHandler) AddComment(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid course id",
		})
		return
	}

	var req service.AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	branchID, _ := c.Get("branch_id")

	comment, err := h.commentService.AddCommentAsTeacher(userID.(uint), branchID.(uint), uint(courseID), &req)
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

	c.JSON(http.StatusOK, comment)
}

