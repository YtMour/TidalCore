package api

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"tidalcore-backend/internal/service"
	"tidalcore-backend/pkg/response"
)

type CheckinHandler struct {
	checkinService *service.CheckinService
}

func NewCheckinHandler() *CheckinHandler {
	return &CheckinHandler{
		checkinService: service.NewCheckinService(),
	}
}

func (h *CheckinHandler) Checkin(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "无效的用户")
		return
	}

	var req service.CheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效")
		return
	}

	resp, err := h.checkinService.Checkin(userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyCheckedIn) {
			response.BadRequest(c, "今日已打卡")
			return
		}
		response.ServerError(c, "打卡失败")
		return
	}

	response.SuccessWithMsg(c, "打卡成功", resp)
}

func (h *CheckinHandler) GetHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "无效的用户")
		return
	}

	limitStr := c.DefaultQuery("limit", "30")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 30
	}
	if limit > 365 {
		limit = 365
	}

	checkins, err := h.checkinService.GetHistory(userID, limit)
	if err != nil {
		response.ServerError(c, "获取打卡记录失败")
		return
	}

	response.Success(c, checkins)
}

func (h *CheckinHandler) GetHeatmap(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "无效的用户")
		return
	}

	daysStr := c.DefaultQuery("days", "365")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 365
	}
	if days > 365 {
		days = 365
	}

	checkins, err := h.checkinService.GetHeatmap(userID, days)
	if err != nil {
		response.ServerError(c, "获取热力图数据失败")
		return
	}

	response.Success(c, checkins)
}

func (h *CheckinHandler) GetGlobalHeatmap(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "365")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 365
	}
	if days > 365 {
		days = 365
	}

	heatmap, err := h.checkinService.GetGlobalHeatmap(days)
	if err != nil {
		response.ServerError(c, "获取全站热力图失败")
		return
	}

	response.Success(c, heatmap)
}
