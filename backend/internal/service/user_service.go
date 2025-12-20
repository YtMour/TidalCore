package service

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"tidalcore-backend/internal/auth"
	"tidalcore-backend/internal/model"
	"tidalcore-backend/internal/repository"
)

var (
	ErrUserExists       = errors.New("username already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidUsername  = errors.New("invalid username format")
	ErrOldPasswordWrong = errors.New("old password is incorrect")
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=50"`
	Password    string `json:"password" binding:"required,min=6,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (s *UserService) Register(req *RegisterRequest) (*AuthResponse, error) {
	username := strings.TrimSpace(req.Username)
	if !usernameRegex.MatchString(username) {
		return nil, ErrInvalidUsername
	}

	displayName := strings.TrimSpace(req.DisplayName)
	if displayName == "" {
		displayName = username // 如果未提供显示名称，默认使用用户名
	}

	exists, err := s.userRepo.ExistsByUsername(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateTokenWithAdmin(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) Login(req *LoginRequest) (*AuthResponse, error) {
	username := strings.TrimSpace(req.Username)

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	token, err := auth.GenerateTokenWithAdmin(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *UserService) GetLeaderboard(limit int) ([]model.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.userRepo.GetLeaderboard(limit)
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name" binding:"omitempty,min=1,max=50"`
}

// UpdateUsernameRequest 更新用户名请求
type UpdateUsernameRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// UpdateProfile 更新用户显示名称
func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	displayName := strings.TrimSpace(req.DisplayName)
	if displayName != "" {
		user.DisplayName = displayName
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUsername 更新用户名
func (s *UserService) UpdateUsername(userID uint, req *UpdateUsernameRequest) (*model.User, error) {
	username := strings.TrimSpace(req.Username)
	if !usernameRegex.MatchString(username) {
		return nil, ErrInvalidUsername
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 检查新用户名是否已存在（排除自己）
	if username != user.Username {
		exists, err := s.userRepo.ExistsByUsername(username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrUserExists
		}
	}

	user.Username = username
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdatePassword 更新密码
func (s *UserService) UpdatePassword(userID uint, req *UpdatePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 验证旧密码
	if !auth.CheckPassword(req.OldPassword, user.PasswordHash) {
		return ErrOldPasswordWrong
	}

	// 生成新密码哈希
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

// InitAdmin 初始化管理员账号
// 如果管理员账号不存在则创建，如果存在则更新密码
func (s *UserService) InitAdmin(username, password string) error {
	if username == "" || password == "" {
		return nil // 未配置管理员，跳过
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	// 检查管理员账号是否存在
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新管理员账号
			admin := &model.User{
				Username:     username,
				DisplayName:  "管理员",
				PasswordHash: hashedPassword,
				IsAdmin:      true,
			}
			return s.userRepo.Create(admin)
		}
		return err
	}

	// 更新现有管理员账号
	user.PasswordHash = hashedPassword
	user.IsAdmin = true
	return s.userRepo.Update(user)
}

// GetAllUsers 获取所有用户列表（管理员功能）
func (s *UserService) GetAllUsers(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.GetAllUsers(page, pageSize)
}

// DeleteUser 删除用户（管理员功能）
func (s *UserService) DeleteUser(userID uint) error {
	return s.userRepo.Delete(userID)
}

// SetUserAdmin 设置用户管理员状态（管理员功能）
func (s *UserService) SetUserAdmin(userID uint, isAdmin bool) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ErrUserNotFound
	}
	user.IsAdmin = isAdmin
	return s.userRepo.Update(user)
}

// UpdateUserStatsRequest 更新用户统计数据请求
type UpdateUserStatsRequest struct {
	Streak       *int    `json:"streak"`
	MaxStreak    *int    `json:"max_streak"`
	TotalCheckin *int    `json:"total_checkin"`
	Title        *string `json:"title"`
}

// UpdateUserStats 更新用户统计数据（管理员功能）
func (s *UserService) UpdateUserStats(userID uint, req *UpdateUserStatsRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.Streak != nil {
		user.Streak = *req.Streak
	}
	if req.MaxStreak != nil {
		user.MaxStreak = *req.MaxStreak
	}
	if req.TotalCheckin != nil {
		user.TotalCheckin = *req.TotalCheckin
	}
	if req.Title != nil {
		user.Title = *req.Title
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
