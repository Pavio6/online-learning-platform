package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// InitLogger 初始化日志系统
func InitLogger(logLevel string) {
	logger = logrus.New()

	// 设置日志格式为JSON（生产环境）或文本（开发环境）
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	// 设置日志级别
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 设置输出到标准输出
	logger.SetOutput(os.Stdout)
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	if logger == nil {
		InitLogger("info")
	}
	return logger
}

// WithField 添加字段到日志
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithFields 添加多个字段到日志
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

