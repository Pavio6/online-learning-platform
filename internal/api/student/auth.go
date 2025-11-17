package student

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/service"
)

// AuthHandler 学生认证处理器
type AuthHandler struct {
	userService *service.UserService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService: service.NewUserService(),
	}
}

// Register 学生注册
// @Summary 学生注册
// @Description 学生注册账号，需要选择校区
// @Tags 学生认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 200 {object} service.RegisterResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/student/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidParam,
			"message": err.Error(),
		})
		return
	}

	resp, err := h.userService.Register(&req)
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

	c.JSON(http.StatusOK, resp)
}

// Login 学生登录
// @Summary 学生登录
// @Description 学生登录获取token
// @Tags 学生认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/student/auth/login [post]
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

	c.JSON(http.StatusOK, resp)
}

// GetProfile 获取个人信息
// @Summary 获取个人信息
// @Description 获取当前登录学生的个人信息
// @Tags 学生认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.UserInfo
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/student/profile [get]
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

// GetBranches 获取校区列表
// @Summary 获取校区列表
// @Description 获取所有可选的校区列表，用于注册时选择
// @Tags 学生认证
// @Accept json
// @Produce json
// @Success 200 {array} models.Branches
// @Router /api/v1/student/branches [get]
func (h *AuthHandler) GetBranches(c *gin.Context) {
	branches, err := h.userService.GetBranches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeInternal,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, branches)
}
