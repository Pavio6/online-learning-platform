package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

const (
	// 通用错误码
	ErrCodeInternal     ErrorCode = 1000 // 内部服务器错误
	ErrCodeInvalidParam ErrorCode = 1001 // 参数错误
	ErrCodeNotFound     ErrorCode = 1002 // 资源不存在
	ErrCodeUnauthorized ErrorCode = 1003 // 未授权
	ErrCodeForbidden    ErrorCode = 1004 // 禁止访问

	// 用户相关错误码
	ErrCodeUserNotFound      ErrorCode = 2001 // 用户不存在
	ErrCodeUserAlreadyExists ErrorCode = 2002 // 用户已存在
	ErrCodeInvalidPassword   ErrorCode = 2003 // 密码错误
	ErrCodeInvalidRole       ErrorCode = 2004 // 角色无效

	// 课程相关错误码
	ErrCodeCourseNotFound     ErrorCode = 3001 // 课程不存在
	ErrCodeChapterNotFound    ErrorCode = 3002 // 章节不存在
	ErrCodeLessonNotFound     ErrorCode = 3003 // 课程不存在
	ErrCodeTaskNotFound       ErrorCode = 3004 // 任务不存在
	ErrCodeNotCourseInstructor ErrorCode = 3005 // 不是课程教师

	// 学习相关错误码
	ErrCodeNotEnrolled        ErrorCode = 4001 // 未报名课程
	ErrCodeAlreadyEnrolled    ErrorCode = 4002 // 已报名课程
	ErrCodeCannotComment      ErrorCode = 4003 // 不能评论（未报名）

	// 作业相关错误码
	ErrCodeAnswerNotFound     ErrorCode = 5001 // 作业不存在
	ErrCodeAnswerAlreadyGraded ErrorCode = 5002 // 作业已评分
	ErrCodeInvalidScore       ErrorCode = 5003 // 分数无效

	// 数据库相关错误码
	ErrCodeDatabaseError      ErrorCode = 6001 // 数据库错误
	ErrCodeBranchNotFound     ErrorCode = 6002 // 分支不存在
)

// AppError 应用错误
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatus 返回HTTP状态码
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeInvalidParam:
		return http.StatusBadRequest
	case ErrCodeNotFound, ErrCodeUserNotFound, ErrCodeCourseNotFound,
		ErrCodeChapterNotFound, ErrCodeLessonNotFound, ErrCodeTaskNotFound,
		ErrCodeAnswerNotFound:
		return http.StatusNotFound
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden, ErrCodeNotCourseInstructor, ErrCodeCannotComment:
		return http.StatusForbidden
	case ErrCodeUserAlreadyExists, ErrCodeAlreadyEnrolled:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// NewAppError 创建新的应用错误
func NewAppError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// WrapError 包装错误
func WrapError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误
var (
	ErrInternal      = NewAppError(ErrCodeInternal, "内部服务器错误")
	ErrInvalidParam  = NewAppError(ErrCodeInvalidParam, "参数错误")
	ErrNotFound      = NewAppError(ErrCodeNotFound, "资源不存在")
	ErrUnauthorized  = NewAppError(ErrCodeUnauthorized, "未授权")
	ErrForbidden     = NewAppError(ErrCodeForbidden, "禁止访问")

	ErrUserNotFound      = NewAppError(ErrCodeUserNotFound, "用户不存在")
	ErrUserAlreadyExists = NewAppError(ErrCodeUserAlreadyExists, "用户已存在")
	ErrInvalidPassword   = NewAppError(ErrCodeInvalidPassword, "密码错误")
	ErrInvalidRole       = NewAppError(ErrCodeInvalidRole, "角色无效")

	ErrCourseNotFound      = NewAppError(ErrCodeCourseNotFound, "课程不存在")
	ErrChapterNotFound     = NewAppError(ErrCodeChapterNotFound, "章节不存在")
	ErrLessonNotFound      = NewAppError(ErrCodeLessonNotFound, "课程不存在")
	ErrTaskNotFound        = NewAppError(ErrCodeTaskNotFound, "任务不存在")
	ErrNotCourseInstructor = NewAppError(ErrCodeNotCourseInstructor, "不是课程教师")

	ErrNotEnrolled   = NewAppError(ErrCodeNotEnrolled, "未报名课程")
	ErrAlreadyEnrolled = NewAppError(ErrCodeAlreadyEnrolled, "已报名课程")
	ErrCannotComment = NewAppError(ErrCodeCannotComment, "不能评论，请先报名课程")

	ErrAnswerNotFound     = NewAppError(ErrCodeAnswerNotFound, "作业不存在")
	ErrAnswerAlreadyGraded = NewAppError(ErrCodeAnswerAlreadyGraded, "作业已评分")
	ErrInvalidScore       = NewAppError(ErrCodeInvalidScore, "分数无效")

	ErrDatabaseError  = NewAppError(ErrCodeDatabaseError, "数据库错误")
	ErrBranchNotFound = NewAppError(ErrCodeBranchNotFound, "分支不存在")
)

