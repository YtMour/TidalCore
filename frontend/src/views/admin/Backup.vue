<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Download, RefreshRight, Delete, FolderAdd } from '@element-plus/icons-vue'
import {
  createBackup,
  getBackups,
  downloadBackup,
  restoreBackup,
  uploadAndRestore,
  deleteBackup,
  type BackupInfo
} from '@/api/backup'

const loading = ref(false)
const backups = ref<BackupInfo[]>([])
const uploading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const maxBackups = 10

const canCreateBackup = computed(() => backups.value.length < maxBackups)

onMounted(() => {
  loadBackups()
})

async function loadBackups() {
  loading.value = true
  try {
    const result = await getBackups()
    backups.value = result || []
  } catch (error: any) {
    ElMessage.error(error.message || '加载备份列表失败')
    backups.value = []
  } finally {
    loading.value = false
  }
}

async function handleCreateBackup() {
  if (!canCreateBackup.value) {
    ElMessage.warning(`备份数量已达上限 (${maxBackups})，请先删除旧备份`)
    return
  }

  loading.value = true
  try {
    const backup = await createBackup()
    ElMessage.success(`备份创建成功: ${backup.filename}`)
    await loadBackups()
  } catch (error: any) {
    ElMessage.error(error.message || '创建备份失败')
  } finally {
    loading.value = false
  }
}

function handleDownload(backup: BackupInfo) {
  downloadBackup(backup.filename)
  ElMessage.success('开始下载备份文件')
}

async function handleRestore(backup: BackupInfo) {
  try {
    await ElMessageBox.confirm(
      `确定要从备份 "${backup.filename}" 恢复数据吗？此操作将清空当前数据并替换为备份数据，不可撤销！`,
      '危险操作',
      {
        confirmButtonText: '确定恢复',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    loading.value = true
    await restoreBackup(backup.filename)
    ElMessage.success('数据恢复成功')
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '恢复失败')
    }
  } finally {
    loading.value = false
  }
}

async function handleDelete(backup: BackupInfo) {
  try {
    await ElMessageBox.confirm(
      `确定要删除备份 "${backup.filename}" 吗？`,
      '删除确认',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    loading.value = true
    await deleteBackup(backup.filename)
    ElMessage.success('备份已删除')
    await loadBackups()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  } finally {
    loading.value = false
  }
}

async function handleUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  // 验证文件类型
  if (!file.name.endsWith('.sql')) {
    ElMessage.error('只支持 .sql 文件')
    input.value = ''
    return
  }

  // 验证文件大小 (50MB)
  if (file.size > 50 * 1024 * 1024) {
    ElMessage.error('文件大小不能超过 50MB')
    input.value = ''
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要从文件 "${file.name}" 恢复数据吗？此操作将清空当前数据并替换为备份数据，不可撤销！`,
      '危险操作',
      {
        confirmButtonText: '确定恢复',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    uploading.value = true
    await uploadAndRestore(file)
    ElMessage.success('数据恢复成功')
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '恢复失败')
    }
  } finally {
    uploading.value = false
    input.value = ''
  }
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

function triggerFileInput() {
  fileInput.value?.click()
}
</script>

<template>
  <div class="backup-page">
    <!-- 操作区 -->
    <div class="action-bar">
      <div class="action-left">
        <el-button
          type="primary"
          :icon="FolderAdd"
          :loading="loading"
          :disabled="!canCreateBackup"
          @click="handleCreateBackup"
        >
          创建备份
        </el-button>

        <div class="upload-wrapper">
          <input
            ref="fileInput"
            type="file"
            accept=".sql"
            class="file-input"
            :disabled="uploading"
            @change="handleUpload"
          />
          <el-button
            type="warning"
            :icon="Upload"
            :loading="uploading"
            @click="triggerFileInput"
          >
            上传恢复
          </el-button>
        </div>
      </div>

      <div class="action-right">
        <span class="backup-count">
          {{ backups.length }} / {{ maxBackups }} 个备份
        </span>
        <el-button
          :icon="RefreshRight"
          :loading="loading"
          @click="loadBackups"
        >
          刷新
        </el-button>
      </div>
    </div>

    <!-- 备份列表 -->
    <div class="backup-list" v-loading="loading">
      <div v-if="backups.length === 0" class="empty-state">
        <div class="empty-icon">📦</div>
        <div class="empty-text">暂无备份</div>
        <div class="empty-hint">点击"创建备份"按钮创建第一个备份</div>
      </div>

      <div
        v-for="backup in backups"
        :key="backup.filename"
        class="backup-item"
      >
        <div class="backup-info">
          <div class="backup-name">{{ backup.filename }}</div>
          <div class="backup-meta">
            <span class="meta-item">
              <span class="meta-label">大小:</span>
              <span class="meta-value">{{ formatSize(backup.size) }}</span>
            </span>
            <span class="meta-item">
              <span class="meta-label">创建时间:</span>
              <span class="meta-value">{{ formatDate(backup.created_at) }}</span>
            </span>
          </div>
        </div>

        <div class="backup-actions">
          <el-button
            size="small"
            :icon="Download"
            @click="handleDownload(backup)"
          >
            下载
          </el-button>
          <el-button
            size="small"
            type="warning"
            :icon="RefreshRight"
            @click="handleRestore(backup)"
          >
            恢复
          </el-button>
          <el-button
            size="small"
            type="danger"
            :icon="Delete"
            @click="handleDelete(backup)"
          >
            删除
          </el-button>
        </div>
      </div>
    </div>

    <!-- 提示信息 -->
    <div class="tips-section">
      <h3 class="tips-title">使用说明</h3>
      <ul class="tips-list">
        <li>备份包含所有用户数据、打卡记录和访问统计</li>
        <li>最多保留 {{ maxBackups }} 个备份，超出后需手动删除旧备份</li>
        <li>恢复操作会清空当前所有数据，请谨慎操作</li>
        <li>可以下载备份文件保存到本地，也可以上传本地备份文件进行恢复</li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.backup-page {
  max-width: 900px;
  margin: 0 auto;
}

/* ===== 操作栏 ===== */
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.action-left {
  display: flex;
  gap: 12px;
  align-items: center;
}

.action-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

.backup-count {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.upload-wrapper {
  display: inline-block;
}

.upload-wrapper .file-input {
  display: none;
}

/* ===== 备份列表 ===== */
.backup-list {
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: 16px;
  min-height: 200px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
}

.backup-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(56, 189, 248, 0.08);
  transition: background 0.2s;
}

.backup-item:last-child {
  border-bottom: none;
}

.backup-item:hover {
  background: rgba(56, 189, 248, 0.05);
}

.backup-info {
  flex: 1;
  min-width: 0;
}

.backup-name {
  font-size: 15px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 6px;
  word-break: break-all;
}

.backup-meta {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.meta-item {
  font-size: 13px;
}

.meta-label {
  color: rgba(255, 255, 255, 0.5);
  margin-right: 4px;
}

.meta-value {
  color: rgba(255, 255, 255, 0.7);
}

.backup-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  margin-left: 16px;
}

/* ===== 提示信息 ===== */
.tips-section {
  margin-top: 24px;
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: 16px;
  padding: 20px;
}

.tips-title {
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 12px;
}

.tips-list {
  margin: 0;
  padding-left: 20px;
}

.tips-list li {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  line-height: 1.8;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .action-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .action-left,
  .action-right {
    justify-content: center;
  }

  .backup-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .backup-actions {
    margin-left: 0;
    width: 100%;
    justify-content: flex-start;
  }

  .backup-meta {
    flex-direction: column;
    gap: 4px;
  }
}
</style>
