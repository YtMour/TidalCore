package model

import (
	"time"
)

type Checkin struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Duration  int       `gorm:"not null" json:"duration"` // 训练时长(秒)
	Cycles    int       `gorm:"not null" json:"cycles"`   // 完成循环数
	CheckedAt time.Time `gorm:"index;not null" json:"checked_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (Checkin) TableName() string {
	return "checkins"
}
