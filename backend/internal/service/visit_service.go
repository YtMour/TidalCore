package service

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"tidalcore-backend/internal/model"
	"tidalcore-backend/internal/repository"
	"tidalcore-backend/pkg/database"
)

type VisitService struct {
	repo *repository.VisitRepository
	loc  *time.Location
}

func NewVisitService() *VisitService {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return &VisitService{
		repo: repository.NewVisitRepository(database.Get()),
		loc:  loc,
	}
}

// generateVisitorID 根据 IP 和 User-Agent 生成访客标识
func (s *VisitService) generateVisitorID(ip, userAgent string) string {
	data := ip + "|" + userAgent
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16]) // 取前 16 字节
}

// RecordVisit 记录访问
func (s *VisitService) RecordVisit(ip, userAgent string) error {
	visitorID := s.generateVisitorID(ip, userAgent)

	// 检查今日是否已记录
	visited, err := s.repo.HasVisitedToday(visitorID, s.loc)
	if err != nil {
		return err
	}

	if visited {
		// 今日已记录，不重复记录
		return nil
	}

	// 创建访问记录
	visit := &model.Visit{
		VisitorID:   visitorID,
		UserAgent:   userAgent,
		VisitedDate: time.Now().In(s.loc).Format("2006-01-02"),
	}

	return s.repo.Create(visit)
}

// GetTodayCount 获取今日访问人数
func (s *VisitService) GetTodayCount() (int64, error) {
	return s.repo.GetTodayVisitCount(s.loc)
}

// GetHeatmap 获取访问热力图
func (s *VisitService) GetHeatmap(days int) (map[string]int, error) {
	if days <= 0 {
		days = 365
	}
	if days > 365 {
		days = 365
	}
	return s.repo.GetHeatmap(days, s.loc)
}

// GetStats 获取访问统计
func (s *VisitService) GetStats() (todayCount int64, totalCount int64, err error) {
	todayCount, err = s.repo.GetTodayVisitCount(s.loc)
	if err != nil {
		return 0, 0, err
	}

	totalCount, err = s.repo.GetTotalVisits()
	if err != nil {
		return 0, 0, err
	}

	return todayCount, totalCount, nil
}
