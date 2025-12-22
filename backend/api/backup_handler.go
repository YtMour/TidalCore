package api

import (
	"io"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"tidalcore-backend/internal/service"
	"tidalcore-backend/pkg/response"
)

type BackupHandler struct {
	backupService *service.BackupService
}

func NewBackupHandler() *BackupHandler {
	return &BackupHandler{
		backupService: service.NewBackupService(),
	}
}

// CreateBackup 创建备份
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	backup, err := h.backupService.CreateBackup()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, backup)
}

// ListBackups 获取备份列表
func (h *BackupHandler) ListBackups(c *gin.Context) {
	backups, err := h.backupService.ListBackups()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, backups)
}

// DownloadBackup 下载备份文件
func (h *BackupHandler) DownloadBackup(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.BadRequest(c, "文件名不能为空")
		return
	}

	path, err := h.backupService.GetBackupPath(filename)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(path)
}

// RestoreBackup 从备份恢复
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.BadRequest(c, "文件名不能为空")
		return
	}

	if err := h.backupService.RestoreBackup(filename); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "恢复成功"})
}

// UploadAndRestore 上传并恢复
func (h *BackupHandler) UploadAndRestore(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请上传文件")
		return
	}
	defer file.Close()

	// 验证文件扩展名
	ext := filepath.Ext(header.Filename)
	if ext != ".sql" {
		response.BadRequest(c, "只支持 .sql 文件")
		return
	}

	// 限制文件大小 (50MB)
	if header.Size > 50*1024*1024 {
		response.BadRequest(c, "文件大小不能超过 50MB")
		return
	}

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		response.ServerError(c, "读取文件失败")
		return
	}

	// 执行恢复
	if err := h.backupService.RestoreFromSQL(string(content)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "恢复成功"})
}

// DeleteBackup 删除备份
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.BadRequest(c, "文件名不能为空")
		return
	}

	if err := h.backupService.DeleteBackup(filename); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
