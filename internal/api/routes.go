package api

import (
	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/api/middleware"
	"online-learning-platform/internal/api/student"
	"online-learning-platform/internal/api/teacher"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine) {
	// 初始化处理器
	studentAuthHandler := student.NewAuthHandler()
	studentCourseHandler := student.NewCourseHandler()
	studentTaskHandler := student.NewTaskHandler()
	teacherAuthHandler := teacher.NewAuthHandler()
	teacherCourseHandler := teacher.NewCourseHandler()
	teacherTaskHandler := teacher.NewTaskHandler()

	// 学生端API
	studentAPI := r.Group("/api/v1/student")
	{
		// 认证相关（不需要认证）
		auth := studentAPI.Group("/auth")
		{
			auth.POST("/register", studentAuthHandler.Register)
			auth.POST("/login", studentAuthHandler.Login)
		}

		// 获取校区列表（不需要认证）
		studentAPI.GET("/branches", studentAuthHandler.GetBranches)

		// 课程相关（不需要认证）
		studentAPI.GET("/courses", studentCourseHandler.ListCourses)
		studentAPI.GET("/courses/:id", studentCourseHandler.GetCourse)
		studentAPI.GET("/courses/:id/tasks", studentTaskHandler.ListTasksByCourse)
		studentAPI.GET("/tasks/:id", studentTaskHandler.GetTask)

		// 需要认证的接口
		studentAPI.Use(middleware.AuthMiddleware())
		{
			studentAPI.GET("/profile", studentAuthHandler.GetProfile)
		}
	}

	// 教师端API
	teacherAPI := r.Group("/api/v1/teacher")
	{
		// 认证相关（不需要认证）
		auth := teacherAPI.Group("/auth")
		{
			auth.POST("/login", teacherAuthHandler.Login)
		}

		// 需要认证的接口
		teacherAPI.Use(middleware.AuthMiddleware())
		teacherAPI.Use(middleware.RequireRole("teacher"))
		{
			// 个人信息
			teacherAPI.GET("/profile", teacherAuthHandler.GetProfile)

			// 课程管理
			teacherAPI.POST("/courses", teacherCourseHandler.CreateCourse)
			teacherAPI.GET("/courses", teacherCourseHandler.ListCourses)
			teacherAPI.GET("/courses/:id", teacherCourseHandler.GetCourse)
			teacherAPI.POST("/courses/:id/chapters", teacherCourseHandler.CreateChapter)
			teacherAPI.POST("/courses/:id/chapters/:chapter_id/lessons", teacherCourseHandler.CreateLesson)

			// 任务管理
			teacherAPI.POST("/lessons/:id/tasks", teacherTaskHandler.CreateTask)
			teacherAPI.GET("/tasks/:id", teacherTaskHandler.GetTask)
			teacherAPI.GET("/courses/:id/tasks", teacherTaskHandler.ListTasksByCourse)
		}
	}
}

