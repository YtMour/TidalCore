package repository

import (
	"time"

	"gorm.io/gorm"

	"tidalcore-backend/internal/model"
	"tidalcore-backend/pkg/database"
)

type CheckinRepository struct {
	db *gorm.DB
}

func NewCheckinRepository() *CheckinRepository {
	return &CheckinRepository{db: database.Get()}
}

func (r *CheckinRepository) Create(checkin *model.Checkin) error {
	return r.db.Create(checkin).Error
}

func (r *CheckinRepository) GetByUserID(userID uint, limit int) ([]model.Checkin, error) {
	var checkins []model.Checkin
	err := r.db.Where("user_id = ?", userID).
		Order("checked_at DESC").
		Limit(limit).
		Find(&checkins).Error
	return checkins, err
}

func (r *CheckinRepository) GetByUserIDAndDateRange(userID uint, start, end time.Time) ([]model.Checkin, error) {
	var checkins []model.Checkin
	err := r.db.Where("user_id = ? AND checked_at BETWEEN ? AND ?", userID, start, end).
		Order("checked_at ASC").
		Find(&checkins).Error
	return checkins, err
}

func (r *CheckinRepository) HasCheckedToday(userID uint, loc *time.Location) (bool, error) {
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var count int64
	err := r.db.Model(&model.Checkin{}).
		Where("user_id = ? AND checked_at >= ? AND checked_at < ?", userID, startOfDay, endOfDay).
		Count(&count).Error
	return count > 0, err
}

func (r *CheckinRepository) GetGlobalHeatmap(days int, loc *time.Location) (map[string]int, error) {
	startDate := time.Now().In(loc).AddDate(0, 0, -days)

	type Result struct {
		Date  string
		Count int
	}
	var results []Result

	err := r.db.Model(&model.Checkin{}).
		Select("DATE(checked_at) as date, COUNT(DISTINCT user_id) as count").
		Where("checked_at >= ?", startDate).
		Group("DATE(checked_at)").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	heatmap := make(map[string]int)
	for _, r := range results {
		heatmap[r.Date] = r.Count
	}
	return heatmap, nil
}
