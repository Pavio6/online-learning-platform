package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"online-learning-platform/internal/config"
)

var centralDB *gorm.DB

// InitCentralDB 初始化中央服务器数据库连接
func InitCentralDB(cfg config.DBSettings) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to central database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != "" {
		if duration, err := time.ParseDuration(cfg.ConnMaxLifetime); err == nil {
			sqlDB.SetConnMaxLifetime(duration)
		}
	}

	centralDB = db
	return nil
}

// GetCentralDB 获取中央服务器数据库连接
func GetCentralDB() *gorm.DB {
	return centralDB
}

// CloseCentralDB 关闭中央服务器数据库连接
func CloseCentralDB() error {
	if centralDB != nil {
		sqlDB, err := centralDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

