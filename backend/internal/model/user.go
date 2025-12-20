package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	DisplayName  string         `gorm:"size:50;not null" json:"display_name"` // 显示名称，支持中文和符号
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	IsAdmin      bool           `gorm:"default:false" json:"is_admin"` // 是否为管理员
	Title        string         `gorm:"size:50;default:''" json:"title"`  // 用户称号，空字符串表示自动根据打卡次数计算
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
