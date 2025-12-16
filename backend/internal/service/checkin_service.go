package service

import (
	"errors"
	"time"

	"tidalcore-backend/config"
	"tidalcore-backend/internal/model"
	"tidalcore-backend/internal/repository"
)

var (
	ErrAlreadyCheckedIn = errors.New("already checked in today")
)

type CheckinService struct {
	checkinRepo *repository.CheckinRepository
	userRepo    *repository.UserRepository
	location    *time.Location
}

func NewCheckinService() *CheckinService {
	loc := time.Local
	cfg := config.Get()
	if cfg != nil && cfg.Server.Timezone != "" {
		if l, err := time.LoadLocation(cfg.Server.Timezone); err == nil {
			loc = l
		}
	}

	return &CheckinService{
		checkinRepo: repository.NewCheckinRepository(),
		userRepo:    repository.NewUserRepository(),
		location:    loc,
	}
}

type CheckinRequest struct {
	Duration int `json:"duration" binding:"required,min=1,max=7200"`
	Cycles   int `json:"cycles" binding:"required,min=1,max=100"`
}

type CheckinResponse struct {
	Checkin       *model.Checkin `json:"checkin"`
	CurrentStreak int            `json:"current_streak"`
	MaxStreak     int            `json:"max_streak"`
	TotalCheckin  int            `json:"total_checkin"`
}

func (s *CheckinService) Checkin(userID uint, req *CheckinRequest) (*CheckinResponse, error) {
	hasChecked, err := s.checkinRepo.HasCheckedToday(userID, s.location)
	if err != nil {
		return nil, err
	}
	if hasChecked {
		return nil, ErrAlreadyCheckedIn
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(s.location)
	checkin := &model.Checkin{
		UserID:    userID,
		Duration:  req.Duration,
		Cycles:    req.Cycles,
		CheckedAt: now,
	}

	if err := s.checkinRepo.Create(checkin); err != nil {
		return nil, err
	}

	s.updateStreak(user, now)
	user.TotalCheckin++
	user.LastCheckin = &now

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &CheckinResponse{
		Checkin:       checkin,
		CurrentStreak: user.Streak,
		MaxStreak:     user.MaxStreak,
		TotalCheckin:  user.TotalCheckin,
	}, nil
}

func (s *CheckinService) updateStreak(user *model.User, now time.Time) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, s.location)

	if user.LastCheckin == nil {
		user.Streak = 1
	} else {
		lastCheckinLocal := user.LastCheckin.In(s.location)
		lastDate := time.Date(lastCheckinLocal.Year(), lastCheckinLocal.Month(), lastCheckinLocal.Day(), 0, 0, 0, 0, s.location)

		diff := int(today.Sub(lastDate).Hours() / 24)

		if diff == 0 {
			// 同一天，不更新 streak
		} else if diff == 1 {
			user.Streak++
		} else {
			user.Streak = 1
		}
	}

	if user.Streak > user.MaxStreak {
		user.MaxStreak = user.Streak
	}
}

func (s *CheckinService) GetHistory(userID uint, limit int) ([]model.Checkin, error) {
	if limit <= 0 || limit > 365 {
		limit = 30
	}
	return s.checkinRepo.GetByUserID(userID, limit)
}

func (s *CheckinService) GetHeatmap(userID uint, days int) ([]model.Checkin, error) {
	if days <= 0 || days > 365 {
		days = 365
	}
	start := time.Now().In(s.location).AddDate(0, 0, -days)
	end := time.Now().In(s.location)
	return s.checkinRepo.GetByUserIDAndDateRange(userID, start, end)
}

func (s *CheckinService) GetGlobalHeatmap(days int) (map[string]int, error) {
	if days <= 0 || days > 365 {
		days = 365
	}
	return s.checkinRepo.GetGlobalHeatmap(days, s.location)
}
