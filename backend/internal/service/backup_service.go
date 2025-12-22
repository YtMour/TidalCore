package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"tidalcore-backend/internal/model"
	"tidalcore-backend/pkg/database"
)

const (
	BackupDir     = "backups"
	MaxBackups    = 10
	BackupPrefix  = "backup_"
	BackupSuffix  = ".sql"
	TimeFormat    = "20060102_150405"
)

type BackupInfo struct {
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

type BackupService struct{}

func NewBackupService() *BackupService {
	return &BackupService{}
}

// ensureBackupDir 确保备份目录存在
func (s *BackupService) ensureBackupDir() error {
	return os.MkdirAll(BackupDir, 0755)
}

// CreateBackup 创建数据库备份
func (s *BackupService) CreateBackup() (*BackupInfo, error) {
	if err := s.ensureBackupDir(); err != nil {
		return nil, fmt.Errorf("创建备份目录失败: %w", err)
	}

	// 检查备份数量限制
	backups, err := s.ListBackups()
	if err != nil {
		return nil, err
	}
	if len(backups) >= MaxBackups {
		return nil, fmt.Errorf("备份数量已达上限 (%d)，请先删除旧备份", MaxBackups)
	}

	// 生成备份文件名
	filename := fmt.Sprintf("%s%s%s", BackupPrefix, time.Now().Format(TimeFormat), BackupSuffix)
	filepath := filepath.Join(BackupDir, filename)

	// 生成 SQL 内容
	sqlContent, err := s.generateSQL()
	if err != nil {
		return nil, fmt.Errorf("生成 SQL 失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filepath, []byte(sqlContent), 0644); err != nil {
		return nil, fmt.Errorf("写入备份文件失败: %w", err)
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	return &BackupInfo{
		Filename:  filename,
		Size:      fileInfo.Size(),
		CreatedAt: fileInfo.ModTime(),
	}, nil
}

// generateSQL 生成所有表的 SQL 语句
func (s *BackupService) generateSQL() (string, error) {
	var sb strings.Builder

	// 写入头部注释
	sb.WriteString("-- TidalCore Database Backup\n")
	sb.WriteString(fmt.Sprintf("-- Generated at: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString("-- This file can be used to restore the database\n\n")

	sb.WriteString("SET FOREIGN_KEY_CHECKS = 0;\n\n")

	// 导出 users 表
	usersSQL, err := s.exportUsersTable()
	if err != nil {
		return "", err
	}
	sb.WriteString(usersSQL)

	// 导出 checkins 表
	checkinsSQL, err := s.exportCheckinsTable()
	if err != nil {
		return "", err
	}
	sb.WriteString(checkinsSQL)

	// 导出 visits 表
	visitsSQL, err := s.exportVisitsTable()
	if err != nil {
		return "", err
	}
	sb.WriteString(visitsSQL)

	sb.WriteString("SET FOREIGN_KEY_CHECKS = 1;\n")

	return sb.String(), nil
}

// exportUsersTable 导出用户表
func (s *BackupService) exportUsersTable() (string, error) {
	var sb strings.Builder
	var users []model.User

	// 使用 Unscoped 获取包括软删除的所有记录
	if err := database.Get().Unscoped().Find(&users).Error; err != nil {
		return "", fmt.Errorf("查询用户表失败: %w", err)
	}

	sb.WriteString("-- Table: users\n")
	sb.WriteString("DELETE FROM `users`;\n")

	if len(users) > 0 {
		sb.WriteString("INSERT INTO `users` (`id`, `username`, `display_name`, `password_hash`, `is_admin`, `title`, `streak`, `max_streak`, `total_checkin`, `last_checkin`, `created_at`, `updated_at`, `deleted_at`) VALUES\n")

		for i, user := range users {
			lastCheckin := "NULL"
			if user.LastCheckin != nil {
				lastCheckin = fmt.Sprintf("'%s'", user.LastCheckin.Format("2006-01-02 15:04:05"))
			}

			deletedAt := "NULL"
			if user.DeletedAt.Valid {
				deletedAt = fmt.Sprintf("'%s'", user.DeletedAt.Time.Format("2006-01-02 15:04:05"))
			}

			isAdmin := 0
			if user.IsAdmin {
				isAdmin = 1
			}

			sb.WriteString(fmt.Sprintf("(%d, '%s', '%s', '%s', %d, '%s', %d, %d, %d, %s, '%s', '%s', %s)",
				user.ID,
				escapeString(user.Username),
				escapeString(user.DisplayName),
				escapeString(user.PasswordHash),
				isAdmin,
				escapeString(user.Title),
				user.Streak,
				user.MaxStreak,
				user.TotalCheckin,
				lastCheckin,
				user.CreatedAt.Format("2006-01-02 15:04:05"),
				user.UpdatedAt.Format("2006-01-02 15:04:05"),
				deletedAt,
			))

			if i < len(users)-1 {
				sb.WriteString(",\n")
			} else {
				sb.WriteString(";\n")
			}
		}
	}

	sb.WriteString("\n")
	return sb.String(), nil
}

// exportCheckinsTable 导出打卡表
func (s *BackupService) exportCheckinsTable() (string, error) {
	var sb strings.Builder
	var checkins []model.Checkin

	if err := database.Get().Find(&checkins).Error; err != nil {
		return "", fmt.Errorf("查询打卡表失败: %w", err)
	}

	sb.WriteString("-- Table: checkins\n")
	sb.WriteString("DELETE FROM `checkins`;\n")

	if len(checkins) > 0 {
		sb.WriteString("INSERT INTO `checkins` (`id`, `user_id`, `duration`, `cycles`, `checked_at`, `created_at`) VALUES\n")

		for i, checkin := range checkins {
			sb.WriteString(fmt.Sprintf("(%d, %d, %d, %d, '%s', '%s')",
				checkin.ID,
				checkin.UserID,
				checkin.Duration,
				checkin.Cycles,
				checkin.CheckedAt.Format("2006-01-02 15:04:05"),
				checkin.CreatedAt.Format("2006-01-02 15:04:05"),
			))

			if i < len(checkins)-1 {
				sb.WriteString(",\n")
			} else {
				sb.WriteString(";\n")
			}
		}
	}

	sb.WriteString("\n")
	return sb.String(), nil
}

// exportVisitsTable 导出访问表
func (s *BackupService) exportVisitsTable() (string, error) {
	var sb strings.Builder
	var visits []model.Visit

	if err := database.Get().Find(&visits).Error; err != nil {
		return "", fmt.Errorf("查询访问表失败: %w", err)
	}

	sb.WriteString("-- Table: visits\n")
	sb.WriteString("DELETE FROM `visits`;\n")

	if len(visits) > 0 {
		sb.WriteString("INSERT INTO `visits` (`id`, `visitor_id`, `user_agent`, `visited_date`, `created_at`) VALUES\n")

		for i, visit := range visits {
			sb.WriteString(fmt.Sprintf("(%d, '%s', '%s', '%s', '%s')",
				visit.ID,
				escapeString(visit.VisitorID),
				escapeString(visit.UserAgent),
				escapeString(visit.VisitedDate),
				visit.CreatedAt.Format("2006-01-02 15:04:05"),
			))

			if i < len(visits)-1 {
				sb.WriteString(",\n")
			} else {
				sb.WriteString(";\n")
			}
		}
	}

	sb.WriteString("\n")
	return sb.String(), nil
}

// ListBackups 获取备份列表
func (s *BackupService) ListBackups() ([]BackupInfo, error) {
	if err := s.ensureBackupDir(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(BackupDir)
	if err != nil {
		return nil, fmt.Errorf("读取备份目录失败: %w", err)
	}

	var backups []BackupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasPrefix(entry.Name(), BackupPrefix) || !strings.HasSuffix(entry.Name(), BackupSuffix) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, BackupInfo{
			Filename:  entry.Name(),
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	// 按时间倒序排列
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return backups, nil
}

// GetBackupPath 获取备份文件路径
func (s *BackupService) GetBackupPath(filename string) (string, error) {
	// 验证文件名格式
	if !strings.HasPrefix(filename, BackupPrefix) || !strings.HasSuffix(filename, BackupSuffix) {
		return "", fmt.Errorf("无效的备份文件名")
	}

	// 防止路径遍历攻击
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return "", fmt.Errorf("无效的备份文件名")
	}

	path := filepath.Join(BackupDir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("备份文件不存在")
	}

	return path, nil
}

// RestoreBackup 从备份文件恢复
func (s *BackupService) RestoreBackup(filename string) error {
	path, err := s.GetBackupPath(filename)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取备份文件失败: %w", err)
	}

	return s.executeSQL(string(content))
}

// RestoreFromSQL 从 SQL 内容恢复
func (s *BackupService) RestoreFromSQL(sqlContent string) error {
	// 基本验证
	if !strings.Contains(sqlContent, "TidalCore Database Backup") {
		return fmt.Errorf("无效的备份文件：缺少 TidalCore 标识")
	}

	return s.executeSQL(sqlContent)
}

// executeSQL 执行 SQL 语句
func (s *BackupService) executeSQL(sqlContent string) error {
	db := database.Get()

	// 先禁用外键检查
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("禁用外键检查失败: %w", err)
	}

	// 确保最后恢复外键检查
	defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// 按顺序删除表数据（先删除有外键依赖的表）
	if err := db.Exec("DELETE FROM `checkins`").Error; err != nil {
		return fmt.Errorf("清空 checkins 表失败: %w", err)
	}
	if err := db.Exec("DELETE FROM `visits`").Error; err != nil {
		return fmt.Errorf("清空 visits 表失败: %w", err)
	}
	if err := db.Exec("DELETE FROM `users`").Error; err != nil {
		return fmt.Errorf("清空 users 表失败: %w", err)
	}

	// 分割 SQL 语句
	statements := splitSQL(sqlContent)

	// 执行 INSERT 语句
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		// 跳过 SET 和 DELETE 语句，只执行 INSERT
		upperStmt := strings.ToUpper(stmt)
		if strings.HasPrefix(upperStmt, "SET") || strings.HasPrefix(upperStmt, "DELETE") {
			continue
		}

		if strings.HasPrefix(upperStmt, "INSERT") {
			if err := db.Exec(stmt).Error; err != nil {
				return fmt.Errorf("执行 SQL 失败: %w", err)
			}
		}
	}

	return nil
}

// DeleteBackup 删除备份文件
func (s *BackupService) DeleteBackup(filename string) error {
	path, err := s.GetBackupPath(filename)
	if err != nil {
		return err
	}

	return os.Remove(path)
}

// escapeString 转义 SQL 字符串
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\x00", "")
	return s
}

// splitSQL 分割 SQL 语句
func splitSQL(content string) []string {
	var statements []string
	var current strings.Builder
	inString := false
	stringChar := byte(0)

	for i := 0; i < len(content); i++ {
		c := content[i]

		if !inString {
			if c == '\'' || c == '"' {
				inString = true
				stringChar = c
			} else if c == ';' {
				statements = append(statements, current.String())
				current.Reset()
				continue
			}
		} else {
			if c == stringChar && (i == 0 || content[i-1] != '\\') {
				inString = false
			}
		}

		current.WriteByte(c)
	}

	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}
