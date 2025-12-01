package student

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/models"
	"online-learning-platform/internal/service"
)

// 引入models以供Swagger引用
var _ models.Answers

// TaskHandler 学生任务处理器
type TaskHandler struct {
	taskService   *service.TaskService
	answerService *service.AnswerService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		taskService:   service.NewTaskService(),
		answerService: service.NewAnswerService(),
	}
}

// GetTask 获取任务详情
// @Summary 获取任务详情
// @Description 获取任务详细信息
// @Tags 学生任务
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} service.TaskInfo
// @Router /api/v1/student/tasks/{id} [get]
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
// @Tags 学生任务
// @Accept json
// @Produce json
// @Param id path int true "课程ID"
// @Success 200 {array} service.TaskInfo
// @Router /api/v1/student/courses/{id}/tasks [get]
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

// SubmitAnswer 学生提交作业
// @Summary 学生提交作业
// @Description 学生上传作业文本或图片
// @Tags 学生任务
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param answer_content formData string false "作业文本内容"
// @Param type formData string false "作业类型(text/image_url)，默认为text"
// @Param file formData file false "作业图片文件"
// @Success 200 {object} models.Answers
// @Router /api/v1/student/tasks/{id}/answers [post]
func (h *TaskHandler) SubmitAnswer(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid task id",
		})
		return
	}

	userID, _ := c.Get("user_id")
	branchID, _ := c.Get("branch_id")

	req := &service.SubmitAnswerRequest{
		AnswerContent: c.PostForm("answer_content"),
		Type:          c.PostForm("type"),
	}

	var fileBytes []byte
	var fileName string
	if fileHeader, err := c.FormFile("file"); err == nil && fileHeader != nil {
		f, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    errors.ErrCodeInvalidParam,
				"message": "failed to open uploaded file",
			})
			return
		}
		defer f.Close()

		fileBytes, err = io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    errors.ErrCodeInternal,
				"message": "failed to read uploaded file",
			})
			return
		}
		fileName = fileHeader.Filename
	}

	answer, err := h.answerService.SubmitAnswer(userID.(uint), branchID.(uint), uint(taskID), req, fileBytes, fileName)
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

	c.JSON(http.StatusOK, answer)
}

// GetMyAnswer 学生查询自己的作业
// @Summary 学生查询作业
// @Description 获取当前学生在指定任务下的作业，如果未提交则返回 null
// @Tags 学生任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} models.Answers "已提交时返回作业信息"
// @Success 200 {object} map[string]interface{} "未提交时返回 {\"submitted\": false}"
// @Router /api/v1/student/tasks/{id}/answers [get]
func (h *TaskHandler) GetMyAnswer(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": "invalid task id",
		})
		return
	}

	userID, _ := c.Get("user_id")
	branchID, _ := c.Get("branch_id")

	answer, err := h.answerService.GetStudentAnswer(userID.(uint), branchID.(uint), uint(taskID))
	if err != nil {
		// 如果错误是"作业不存在"，说明学生未提交作业，这是正常状态，返回友好响应
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.ErrCodeAnswerNotFound {
			c.JSON(http.StatusOK, gin.H{
				"submitted": false,
			})
			return
		}
		// 其他错误正常返回
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

	c.JSON(http.StatusOK, answer)
}
