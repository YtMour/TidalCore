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

// UpdateProfile 更新用户显示名称
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "无效的用户")
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效")
		return
	}

	user, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "用户不存在")
			return
		}
		response.ServerError(c, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", user)
}

// UpdateUsername 更新用户名
func (h *UserHandler) UpdateUsername(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "无效的用户")
		return
	}

	var req service.UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效")
		return
	}

	user, err := h.userService.UpdateUsername(userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "用户不存在")
			return
		}
		if errors.Is(err, service.ErrUserExists) {
			response.BadRequest(c, "用户名已存在")
			return
		}
		if errors.Is(err, service.ErrInvalidUsername) {
			response.BadRequest(c, "用户名只能包含字母、数字和下划线")
			return
		}
		response.ServerError(c, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "用户名更新成功", user)
}

// UpdatePassword 更新密码
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "无效的用户")
		return
	}

	var req service.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效")
		return
	}

	err := h.userService.UpdatePassword(userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "用户不存在")
			return
		}
		if errors.Is(err, service.ErrOldPasswordWrong) {
			response.BadRequest(c, "原密码错误")
			return
		}
		response.ServerError(c, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "密码更新成功", nil)
}

// ========== 管理员接口 ==========

// GetAllUsers 获取所有用户列表（管理员）
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	users, total, err := h.userService.GetAllUsers(page, pageSize)
	if err != nil {
		response.ServerError(c, "获取用户列表失败")
		return
	}

	response.Success(c, gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteUser 删除用户（管理员）
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	// 不能删除自己
	currentUserID := c.GetUint("user_id")
	if uint(userID) == currentUserID {
		response.BadRequest(c, "不能删除自己的账号")
		return
	}

	if err := h.userService.DeleteUser(uint(userID)); err != nil {
		response.ServerError(c, "删除用户失败")
		return
	}

	response.SuccessWithMsg(c, "用户已删除", nil)
}

// SetUserAdmin 设置用户管理员状态（管理员）
func (h *UserHandler) SetUserAdmin(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		IsAdmin bool `json:"is_admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效")
		return
	}

	// 不能修改自己的管理员状态
	currentUserID := c.GetUint("user_id")
	if uint(userID) == currentUserID {
		response.BadRequest(c, "不能修改自己的管理员状态")
		return
	}

	if err := h.userService.SetUserAdmin(uint(userID), req.IsAdmin); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "用户不存在")
			return
		}
		response.ServerError(c, "设置失败")
		return
	}

	response.SuccessWithMsg(c, "设置成功", nil)
}
