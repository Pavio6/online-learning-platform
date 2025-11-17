package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码
type ErrorCode int

const (
	// 通用错误码
	ErrCodeInternalServer ErrorCode = 1000
	ErrCodeBadRequest     ErrorCode = 1001
	ErrCodeUnauthorized   ErrorCode = 1002
	ErrCodeForbidden      ErrorCode = 1003
	ErrCodeNotFound       ErrorCode = 1004

	// 用户相关错误码
	ErrCodeUserNotFound      ErrorCode = 2001
	ErrCodeUserAlreadyExists ErrorCode = 2002
	ErrCodeInvalidCredentials ErrorCode = 2003

	// 课程相关错误码
	ErrCodeCourseNotFound ErrorCode = 3001
	ErrCodeCourseAccessDenied ErrorCode = 3002

	// 任务相关错误码
	ErrCodeTaskNotFound ErrorCode = 4001

	// 作业相关错误码
	ErrCodeAnswerNotFound ErrorCode = 5001
	ErrCodeAnswerAlreadySubmitted ErrorCode = 5002

	// 学习相关错误码
	ErrCodeNotEnrolled ErrorCode = 6001
	ErrCodeAlreadyEnrolled ErrorCode = 6002

	// 评论相关错误码
	ErrCodeCommentNotFound ErrorCode = 7001
	ErrCodeCommentNotAllowed ErrorCode = 7002
)

// AppError 应用错误
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Status  int       `json:"-"`
	Err     error     `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError 创建新的应用错误
func NewAppError(code ErrorCode, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// WrapError 包装错误
func WrapError(err error, code ErrorCode, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// 预定义错误
var (
	ErrInternalServer = NewAppError(ErrCodeInternalServer, "内部服务器错误", http.StatusInternalServerError)
	ErrBadRequest     = NewAppError(ErrCodeBadRequest, "请求参数错误", http.StatusBadRequest)
	ErrUnauthorized   = NewAppError(ErrCodeUnauthorized, "未授权", http.StatusUnauthorized)
	ErrForbidden      = NewAppError(ErrCodeForbidden, "禁止访问", http.StatusForbidden)
	ErrNotFound       = NewAppError(ErrCodeNotFound, "资源不存在", http.StatusNotFound)

	ErrUserNotFound      = NewAppError(ErrCodeUserNotFound, "用户不存在", http.StatusNotFound)
	ErrUserAlreadyExists = NewAppError(ErrCodeUserAlreadyExists, "用户已存在", http.StatusConflict)
	ErrInvalidCredentials = NewAppError(ErrCodeInvalidCredentials, "用户名或密码错误", http.StatusUnauthorized)

	ErrCourseNotFound   = NewAppError(ErrCodeCourseNotFound, "课程不存在", http.StatusNotFound)
	ErrCourseAccessDenied = NewAppError(ErrCodeCourseAccessDenied, "无权访问该课程", http.StatusForbidden)

	ErrTaskNotFound = NewAppError(ErrCodeTaskNotFound, "任务不存在", http.StatusNotFound)

	ErrAnswerNotFound = NewAppError(ErrCodeAnswerNotFound, "作业不存在", http.StatusNotFound)
	ErrAnswerAlreadySubmitted = NewAppError(ErrCodeAnswerAlreadySubmitted, "作业已提交", http.StatusConflict)

	ErrNotEnrolled   = NewAppError(ErrCodeNotEnrolled, "未报名该课程", http.StatusForbidden)
	ErrAlreadyEnrolled = NewAppError(ErrCodeAlreadyEnrolled, "已报名该课程", http.StatusConflict)

	ErrCommentNotFound = NewAppError(ErrCodeCommentNotFound, "评论不存在", http.StatusNotFound)
	ErrCommentNotAllowed = NewAppError(ErrCodeCommentNotAllowed, "未参加课程，无法评论", http.StatusForbidden)
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	RequestID string  `json:"request_id,omitempty"`
}

