import request from './request'

export interface BackupInfo {
  filename: string
  size: number
  created_at: string
}

// 创建备份
export function createBackup(): Promise<BackupInfo> {
  return request.post('/admin/backup')
}

// 获取备份列表
export function getBackups(): Promise<BackupInfo[]> {
  return request.get('/admin/backups')
}

// 下载备份文件
export function downloadBackup(filename: string): void {
  const token = localStorage.getItem('token')
  const url = `/api/v1/admin/backup/${encodeURIComponent(filename)}`

  // 创建一个隐藏的 a 标签来下载文件
  const link = document.createElement('a')
  link.href = url
  link.download = filename

  // 需要带上 token，使用 fetch 下载
  fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`
    }
  })
    .then(response => response.blob())
    .then(blob => {
      const blobUrl = window.URL.createObjectURL(blob)
      link.href = blobUrl
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(blobUrl)
    })
}

// 恢复备份
export function restoreBackup(filename: string): Promise<{ message: string }> {
  return request.post(`/admin/backup/restore/${encodeURIComponent(filename)}`)
}

// 上传并恢复
export function uploadAndRestore(file: File): Promise<{ message: string }> {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/admin/backup/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 删除备份
export function deleteBackup(filename: string): Promise<{ message: string }> {
  return request.delete(`/admin/backup/${encodeURIComponent(filename)}`)
}
