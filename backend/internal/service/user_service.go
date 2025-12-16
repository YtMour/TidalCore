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
	ErrUserExists      = errors.New("username already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUsername = errors.New("invalid username format")
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
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
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
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(user.ID, user.Username)
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

	token, err := auth.GenerateToken(user.ID, user.Username)
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
