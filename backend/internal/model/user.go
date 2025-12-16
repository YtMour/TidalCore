package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Streak       int            `gorm:"default:0" json:"streak"`
	MaxStreak    int            `gorm:"default:0" json:"max_streak"`
	TotalCheckin int            `gorm:"default:0" json:"total_checkin"`
	LastCheckin  *time.Time     `json:"last_checkin"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
