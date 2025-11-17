package teacher

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/service"
)

// TaskHandler 教师任务管理处理器
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		taskService: service.NewTaskService(),
	}
}

// CreateTask 创建任务
// @Summary 创建任务
// @Description 为课程创建新任务
// @Tags 教师任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "课程ID"
// @Param request body service.CreateTaskRequest true "任务信息"
// @Success 200 {object} models.Tasks
// @Router /api/v1/teacher/lessons/:id/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	lessonID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid lesson id",
		})
		return
	}

	instructorID, _ := c.Get("user_id")

	var req service.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": err.Error(),
		})
		return
	}

	task, err := h.taskService.CreateTask(uint(lessonID), instructorID.(uint), &req)
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

	c.JSON(http.StatusOK, task)
}

// GetTask 获取任务详情
// @Summary 获取任务详情
// @Description 获取任务详细信息
// @Tags 教师任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} service.TaskInfo
// @Router /api/v1/teacher/tasks/:id [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid task id",
		})
		return
	}

	task, err := h.taskService.GetTask(uint(taskID))
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

	c.JSON(http.StatusOK, task)
}

// ListTasksByCourse 获取课程的所有任务
// @Summary 获取课程的所有任务
// @Description 获取指定课程的所有任务列表
// @Tags 教师任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "课程ID"
// @Success 200 {array} service.TaskInfo
// @Router /api/v1/teacher/courses/:id/tasks [get]
func (h *TaskHandler) ListTasksByCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid course id",
		})
		return
	}

	tasks, err := h.taskService.ListTasksByCourse(uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeInternal,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

