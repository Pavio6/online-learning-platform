package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"online-learning-platform/internal/config"
	"online-learning-platform/internal/database"
	"online-learning-platform/internal/models"
)

// SyncService 负责中央到分支的数据同步
type SyncService struct {
	lastSync time.Time
}

// NewSyncService 创建同步服务
func NewSyncService() *SyncService {
	return &SyncService{}
}

// RunReplication 执行中央 -> 分支同步
func (s *SyncService) RunReplication() {
	cfg := config.GetConfig()
	if !cfg.Sync.Replication.Enabled {
		return
	}

	start := time.Now()
	if s.lastSync.IsZero() {
		s.lastSync = start.Add(-24 * time.Hour)
	}

	centralDB := database.GetCentralDB()

	var data struct {
		Courses  []models.Courses
		Chapters []models.Chapters
		Lessons  []models.Lessons
		Tasks    []models.Tasks
	}

	tables := cfg.Sync.Replication.Tables
	contains := func(name string) bool {
		for _, t := range tables {
			if t == name {
				return true
			}
		}
		return false
	}

	if contains("courses") {
		if err := centralDB.Where("updated_at >= ?", s.lastSync).Find(&data.Courses).Error; err != nil {
			fmt.Printf("replication: failed to fetch courses: %v\n", err)
		}
	}
	if contains("chapters") {
		if err := centralDB.Where("updated_at >= ?", s.lastSync).Find(&data.Chapters).Error; err != nil {
			fmt.Printf("replication: failed to fetch chapters: %v\n", err)
		}
	}
	if contains("lessons") {
		if err := centralDB.Where("updated_at >= ?", s.lastSync).Find(&data.Lessons).Error; err != nil {
			fmt.Printf("replication: failed to fetch lessons: %v\n", err)
		}
	}
	if contains("tasks") {
		if err := centralDB.Where("updated_at >= ?", s.lastSync).Find(&data.Tasks).Error; err != nil {
			fmt.Printf("replication: failed to fetch tasks: %v\n", err)
		}
	}

	for branchID, branchDB := range database.GetAllBranchDBs() {
		tx := branchDB.Begin()
		if tx.Error != nil {
			fmt.Printf("replication: failed to start tx for branch %d: %v\n", branchID, tx.Error)
			continue
		}

		rollback := func() {
			if err := tx.Rollback().Error; err != nil {
				fmt.Printf("replication: rollback error for branch %d: %v\n", branchID, err)
			}
		}

		if err := upsertCourses(tx, data.Courses); err != nil {
			fmt.Printf("replication: branch %d courses error: %v\n", branchID, err)
			rollback()
			continue
		}
		if err := upsertChapters(tx, data.Chapters); err != nil {
			fmt.Printf("replication: branch %d chapters error: %v\n", branchID, err)
			rollback()
			continue
		}
		if err := upsertLessons(tx, data.Lessons); err != nil {
			fmt.Printf("replication: branch %d lessons error: %v\n", branchID, err)
			rollback()
			continue
		}
		if err := upsertTasks(tx, data.Tasks); err != nil {
			fmt.Printf("replication: branch %d tasks error: %v\n", branchID, err)
			rollback()
			continue
		}

		if err := tx.Commit().Error; err != nil {
			fmt.Printf("replication: commit error for branch %d: %v\n", branchID, err)
		}
	}

	s.lastSync = start
	fmt.Printf("replication finished at %s\n", start.Format(time.RFC3339))
}

func upsertCourses(db *gorm.DB, courses []models.Courses) error {
	for _, course := range courses {
		c := course
		if err := db.Table("courses").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "course_id"}},
			UpdateAll: true,
		}).Create(&c).Error; err != nil {
			return err
		}
	}
	return nil
}

func upsertChapters(db *gorm.DB, chapters []models.Chapters) error {
	for _, chapter := range chapters {
		c := chapter
		if err := db.Table("chapters").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "chapter_id"}},
			UpdateAll: true,
		}).Create(&c).Error; err != nil {
			return err
		}
	}
	return nil
}

func upsertLessons(db *gorm.DB, lessons []models.Lessons) error {
	for _, lesson := range lessons {
		l := lesson
		if err := db.Table("lessons").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "lesson_id"}},
			UpdateAll: true,
		}).Create(&l).Error; err != nil {
			return err
		}
	}
	return nil
}

func upsertTasks(db *gorm.DB, tasks []models.Tasks) error {
	for _, task := range tasks {
		t := task
		if err := db.Table("tasks").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "task_id"}},
			UpdateAll: true,
		}).Create(&t).Error; err != nil {
			return err
		}
	}
	return nil
}
