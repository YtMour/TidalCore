package model

import (
	"time"
)

// Visit 记录站点访问
type Visit struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VisitorID   string    `gorm:"index;size:64;not null" json:"visitor_id"` // 访客标识 (可以是 IP hash 或 fingerprint)
	UserAgent   string    `gorm:"size:512" json:"user_agent"`
	VisitedDate string    `gorm:"index;size:10;not null" json:"visited_date"` // 格式: 2006-01-02
	CreatedAt   time.Time `json:"created_at"`
}

func (Visit) TableName() string {
	return "visits"
}
