package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"online-learning-platform/internal/config"
	"online-learning-platform/internal/database"
	apperrors "online-learning-platform/internal/errors"
	"online-learning-platform/internal/models"
	"online-learning-platform/pkg/utils"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BranchID  uint   `json:"branch_id" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	BranchID uint   `json:"branch_id"`
	Token    string `json:"token"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	BranchID uint   `json:"branch_id"`
	Token    string `json:"token"`
}

// UserInfo 用户信息
type UserInfo struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	BranchID  uint   `json:"branch_id"`
}

// Register 学生注册
func (s *UserService) Register(req *RegisterRequest) (*RegisterResponse, error) {
	// 验证分支是否存在
	branchDB, err := database.GetBranchDBByBranchID(req.BranchID)
	if err != nil {
		return nil, apperrors.ErrBranchNotFound
	}

	// 检查用户名是否已存在
	var existingUser models.Users
	if err := branchDB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// 检查邮箱是否已存在（需要跨分片查询）
	branchDBs := database.GetAllBranchDBs()
	for _, db := range branchDBs {
		var user models.Users
		if err := db.Where("email = ?", req.Email).First(&user).Error; err == nil {
			return nil, apperrors.ErrUserAlreadyExists
		}
	}

	// 加密密码
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := models.Users{
		BranchID:     req.BranchID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "student",
		Status:       "active",
	}

	if err := branchDB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 生成JWT token
	cfg := config.GetConfig()
	expiration, _ := time.ParseDuration(cfg.JWT.Expiration)
	token, err := utils.GenerateToken(user.UserID, user.Username, user.Role, user.BranchID, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &RegisterResponse{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		BranchID: user.BranchID,
		Token:    token,
	}, nil
}

// Login 用户登录（学生和教师）
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 在所有分支节点中查找用户
	branchDBs := database.GetAllBranchDBs()
	var user models.Users
	var foundDB *gorm.DB

	for _, db := range branchDBs {
		var u models.Users
		if err := db.Where("email = ?", req.Email).First(&u).Error; err == nil {
			user = u
			foundDB = db
			break
		}
	}

	if foundDB == nil {
		return nil, apperrors.ErrUserNotFound
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, apperrors.ErrInvalidPassword
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("user account is not active")
	}

	// 生成JWT token
	cfg := config.GetConfig()
	expiration, _ := time.ParseDuration(cfg.JWT.Expiration)
	token, err := utils.GenerateToken(user.UserID, user.Username, user.Role, user.BranchID, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		BranchID: user.BranchID,
		Token:    token,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*UserInfo, error) {
	// 根据user_id找到对应的分支节点
	branchDB, err := database.GetBranchDBByUserID(userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	var user models.Users
	if err := branchDB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &UserInfo{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
		BranchID:  user.BranchID,
	}, nil
}

// GetBranches 获取所有分支列表（用于注册时选择）
func (s *UserService) GetBranches() ([]models.Branches, error) {
	// 从所有分支节点获取branches表数据并合并
	// 每个分支节点都有自己的branches表，但应该包含所有分支的信息
	branchDBs := database.GetAllBranchDBs()
	if len(branchDBs) == 0 {
		return nil, errors.New("no branch databases available")
	}

	// 使用map去重
	branchMap := make(map[uint]models.Branches)

	// 从所有分支节点获取并合并
	for _, db := range branchDBs {
		var branches []models.Branches
		if err := db.Find(&branches).Error; err != nil {
			continue // 如果某个节点查询失败，继续查询其他节点
		}
		for _, branch := range branches {
			branchMap[branch.BranchID] = branch
		}
	}

	// 转换为切片
	result := make([]models.Branches, 0, len(branchMap))
	for _, branch := range branchMap {
		result = append(result, branch)
	}

	return result, nil
}
