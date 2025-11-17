package teacher

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/service"
)

// AuthHandler 教师认证处理器
type AuthHandler struct {
	userService *service.UserService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService: service.NewUserService(),
	}
}

// Login 教师登录
// @Summary 教师登录
// @Description 教师登录获取token
// @Tags 教师认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/teacher/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": err.Error(),
		})
		return
	}

	resp, err := h.userService.Login(&req)
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

	// 验证是否为教师角色
	if resp.Role != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    errors.ErrCodeForbidden,
			"message": "Only teachers can login here",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProfile 获取个人信息
// @Summary 获取个人信息
// @Description 获取当前登录教师的个人信息
// @Tags 教师认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.UserInfo
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/teacher/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errors.ErrCodeUnauthorized,
			"message": "User not authenticated",
		})
		return
	}

	userInfo, err := h.userService.GetUserInfo(userID.(uint))
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

	c.JSON(http.StatusOK, userInfo)
}

