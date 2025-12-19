package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"tidalcore-backend/internal/service"
	"tidalcore-backend/pkg/response"
)

type VisitHandler struct {
	service *service.VisitService
}

func NewVisitHandler() *VisitHandler {
	return &VisitHandler{
		service: service.NewVisitService(),
	}
}

// RecordVisit 记录访问
func (h *VisitHandler) RecordVisit(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	if err := h.service.RecordVisit(ip, userAgent); err != nil {
		response.ServerError(c, "记录访问失败")
		return
	}

	// 获取今日访问数
	todayCount, _ := h.service.GetTodayCount()

	response.Success(c, gin.H{
		"today_visits": todayCount,
	})
}

// GetStats 获取访问统计
func (h *VisitHandler) GetStats(c *gin.Context) {
	todayCount, totalCount, err := h.service.GetStats()
	if err != nil {
		response.ServerError(c, "获取统计失败")
		return
	}

	response.Success(c, gin.H{
		"today_visits": todayCount,
		"total_visits": totalCount,
	})
}

// GetHeatmap 获取访问热力图
func (h *VisitHandler) GetHeatmap(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "365")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 365
	}

	heatmap, err := h.service.GetHeatmap(days)
	if err != nil {
		response.ServerError(c, "获取热力图失败")
		return
	}

	response.Success(c, heatmap)
}
