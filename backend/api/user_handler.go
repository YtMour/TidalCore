package api

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"tidalcore-backend/internal/service"
	"tidalcore-backend/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效")
		return
	}

	resp, err := h.userService.Register(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			response.BadRequest(c, "用户名已存在")
			return
		}
		if errors.Is(err, service.ErrInvalidUsername) {
			response.BadRequest(c, "用户名只能包含字母、数字和下划线")
			return
		}
		response.ServerError(c, "注册失败")
		return
	}

	response.SuccessWithMsg(c, "注册成功", resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效")
		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrInvalidPassword) {
			response.Unauthorized(c, "用户名或密码错误")
			return
		}
		response.ServerError(c, "登录失败")
		return
	}

	response.SuccessWithMsg(c, "登录成功", resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "无效的用户")
		return
	}

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.ServerError(c, "获取用户信息失败")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	users, err := h.userService.GetLeaderboard(limit)
	if err != nil {
		response.ServerError(c, "获取排行榜失败")
		return
	}

	response.Success(c, users)
}
