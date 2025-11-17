package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/api"
	"online-learning-platform/internal/api/middleware"
	"online-learning-platform/internal/config"
	"online-learning-platform/internal/database"
	"online-learning-platform/internal/logger"
	ossclient "online-learning-platform/internal/oss"
	"online-learning-platform/pkg/utils"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.InitLogger(cfg.App.LogLevel)
	logger.Info("Logger initialized")

	// 初始化中央服务器数据库连接
	if err := database.InitCentralDB(cfg.Database.Central); err != nil {
		logger.Fatalf("Failed to initialize central database: %v", err)
	}
	logger.Info("Central database connected")

	// 初始化分支节点数据库连接
	if err := database.InitBranchDBs(cfg.Branches); err != nil {
		logger.Fatalf("Failed to initialize branch databases: %v", err)
	}
	logger.Infof("Branch databases connected: %d branches", len(cfg.Branches))

	// 初始化OSS客户端
	if err := ossclient.InitOSSClient(cfg.OSS); err != nil {
		logger.Fatalf("Failed to initialize OSS client: %v", err)
	}
	logger.Info("OSS client initialized")

	// 初始化JWT
	utils.InitJWT(cfg.JWT.Secret)
	logger.Info("JWT initialized")

	// 设置Gin模式
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin路由
	r := gin.New()

	// 注册中间件
	r.Use(middleware.RequestLogger()) // 请求日志
	r.Use(middleware.ErrorHandler())  // 错误处理
	r.Use(gin.Recovery())             // 恢复panic

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// 注册路由
	api.SetupRoutes(r)

	// 启动HTTP服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	// 优雅关闭
	go func() {
		if err := r.Run(":" + port); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Infof("Server started on port %s", port)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 清理资源
	if err := database.CloseCentralDB(); err != nil {
		logger.Errorf("Failed to close central database: %v", err)
	}
	if err := database.CloseBranchDBs(); err != nil {
		logger.Errorf("Failed to close branch databases: %v", err)
	}

	logger.Info("Server stopped")
}
