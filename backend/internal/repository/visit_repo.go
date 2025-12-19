package repository

import (
	"time"

	"gorm.io/gorm"

	"tidalcore-backend/internal/model"
)

type VisitRepository struct {
	db *gorm.DB
}

func NewVisitRepository(db *gorm.DB) *VisitRepository {
	return &VisitRepository{db: db}
}

// Create 创建访问记录
func (r *VisitRepository) Create(visit *model.Visit) error {
	return r.db.Create(visit).Error
}

// HasVisitedToday 检查访客今日是否已记录
func (r *VisitRepository) HasVisitedToday(visitorID string, loc *time.Location) (bool, error) {
	today := time.Now().In(loc).Format("2006-01-02")

	var count int64
	err := r.db.Model(&model.Visit{}).
		Where("visitor_id = ? AND visited_date = ?", visitorID, today).
		Count(&count).Error

	return count > 0, err
}

// GetTodayVisitCount 获取今日访问人数
func (r *VisitRepository) GetTodayVisitCount(loc *time.Location) (int64, error) {
	today := time.Now().In(loc).Format("2006-01-02")

	var count int64
	err := r.db.Model(&model.Visit{}).
		Where("visited_date = ?", today).
		Count(&count).Error

	return count, err
}

// GetHeatmap 获取访问热力图数据
func (r *VisitRepository) GetHeatmap(days int, loc *time.Location) (map[string]int, error) {
	now := time.Now().In(loc)
	startDate := now.AddDate(0, 0, -days).Format("2006-01-02")

	type Result struct {
		VisitedDate string
		Count       int
	}

	var results []Result
	err := r.db.Model(&model.Visit{}).
		Select("visited_date, COUNT(*) as count").
		Where("visited_date >= ?", startDate).
		Group("visited_date").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	heatmap := make(map[string]int)
	for _, r := range results {
		heatmap[r.VisitedDate] = r.Count
	}

	return heatmap, nil
}

// GetTotalVisits 获取总访问次数
func (r *VisitRepository) GetTotalVisits() (int64, error) {
	var count int64
	err := r.db.Model(&model.Visit{}).Count(&count).Error
	return count, err
}
