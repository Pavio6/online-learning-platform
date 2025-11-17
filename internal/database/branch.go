package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"online-learning-platform/internal/config"
)

var (
	branchDBs map[uint]*gorm.DB
	branchMu  sync.RWMutex
)

// InitBranchDBs 初始化所有分支节点数据库连接
func InitBranchDBs(branches []config.BranchConfig) error {
	branchDBs = make(map[uint]*gorm.DB)

	for _, branch := range branches {
		// 获取数据库配置（使用squash展开的字段）
		dbConfig := branch.DB
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			dbConfig.Host, dbConfig.User, dbConfig.Password,
			dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return fmt.Errorf("failed to connect to branch database (branch_id=%d): %w", branch.BranchID, err)
		}

		// 配置连接池
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB for branch %d: %w", branch.BranchID, err)
		}

		if dbConfig.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		}
		if dbConfig.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		}
		if dbConfig.ConnMaxLifetime != "" {
			if duration, err := time.ParseDuration(dbConfig.ConnMaxLifetime); err == nil {
				sqlDB.SetConnMaxLifetime(duration)
			}
		}

		branchMu.Lock()
		branchDBs[branch.BranchID] = db
		branchMu.Unlock()
	}

	return nil
}

// GetBranchDB 获取指定分支节点的数据库连接
func GetBranchDB(branchID uint) (*gorm.DB, error) {
	branchMu.RLock()
	defer branchMu.RUnlock()

	db, ok := branchDBs[branchID]
	if !ok {
		return nil, fmt.Errorf("branch database not found for branch_id=%d", branchID)
	}

	return db, nil
}

// GetAllBranchDBs 获取所有分支节点数据库连接
func GetAllBranchDBs() map[uint]*gorm.DB {
	branchMu.RLock()
	defer branchMu.RUnlock()

	result := make(map[uint]*gorm.DB)
	for k, v := range branchDBs {
		result[k] = v
	}
	return result
}

// CloseBranchDBs 关闭所有分支节点数据库连接
func CloseBranchDBs() error {
	branchMu.Lock()
	defer branchMu.Unlock()

	var lastErr error
	for branchID, db := range branchDBs {
		sqlDB, err := db.DB()
		if err != nil {
			lastErr = err
			continue
		}
		if err := sqlDB.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close branch %d: %w", branchID, err)
		}
	}

	branchDBs = nil
	return lastErr
}

