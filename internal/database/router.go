package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"online-learning-platform/internal/models"
)

var (
	userBranchCache map[uint]uint // user_id -> branch_id 缓存
	cacheMu         sync.RWMutex
	cacheExpiry     = 5 * time.Minute
	lastCacheUpdate = make(map[uint]time.Time)
)

// GetBranchIDByUserID 根据user_id获取branch_id（带缓存）
func GetBranchIDByUserID(userID uint) (uint, error) {
	// 先检查缓存
	cacheMu.RLock()
	if branchID, ok := userBranchCache[userID]; ok {
		if lastUpdate, exists := lastCacheUpdate[userID]; exists {
			if time.Since(lastUpdate) < cacheExpiry {
				cacheMu.RUnlock()
				return branchID, nil
			}
		}
	}
	cacheMu.RUnlock()

	// 缓存未命中或过期，查询所有分支节点
	branchDBs := GetAllBranchDBs()
	for branchID, db := range branchDBs {
		var user models.Users
		if err := db.Where("user_id = ?", userID).First(&user).Error; err == nil {
			// 找到用户，更新缓存
			cacheMu.Lock()
			if userBranchCache == nil {
				userBranchCache = make(map[uint]uint)
			}
			userBranchCache[userID] = branchID
			lastCacheUpdate[userID] = time.Now()
			cacheMu.Unlock()
			return branchID, nil
		}
	}

	return 0, fmt.Errorf("user not found: user_id=%d", userID)
}

// GetBranchDBByUserID 根据user_id获取对应的分支节点数据库连接
func GetBranchDBByUserID(userID uint) (*gorm.DB, error) {
	branchID, err := GetBranchIDByUserID(userID)
	if err != nil {
		return nil, err
	}

	return GetBranchDB(branchID)
}

// GetBranchDBByBranchID 根据branch_id直接获取分支节点数据库连接
func GetBranchDBByBranchID(branchID uint) (*gorm.DB, error) {
	return GetBranchDB(branchID)
}

// ClearUserCache 清除用户缓存（当用户信息更新时调用）
func ClearUserCache(userID uint) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	delete(userBranchCache, userID)
	delete(lastCacheUpdate, userID)
}

// ClearAllCache 清除所有缓存
func ClearAllCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	userBranchCache = make(map[uint]uint)
	lastCacheUpdate = make(map[uint]time.Time)
}

