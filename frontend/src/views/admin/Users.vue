<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getUsers, deleteUser, setUserAdmin, updateUserStats, type UsersResponse } from '@/api/admin'
import type { UserInfo } from '@/api/auth'
import { useUserStore } from '@/store/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Key, Search, Refresh, Edit, Download } from '@element-plus/icons-vue'
import UserAvatar from '@/components/UserAvatar.vue'

const userStore = useUserStore()

const loading = ref(false)
const users = ref<UserInfo[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 搜索和筛选
const searchQuery = ref('')
const filterAdmin = ref<string>('')
const filterTitle = ref<string>('')

// 编辑用户数据相关
const showEditDialog = ref(false)
const editLoading = ref(false)
const editingUser = ref<UserInfo | null>(null)
const editForm = ref({
  streak: 0,
  max_streak: 0,
  total_checkin: 0,
  title: ''
})

// 称号配置
const titleConfig = [
  { min: 1000, name: '海神降临', color: '#f472b6', icon: '🔱' },
  { min: 730, name: '深渊霸主', color: '#a78bfa', icon: '🦑' },
  { min: 365, name: '深海传奇', color: '#fbbf24', icon: '🌊' },
  { min: 180, name: '海洋大师', color: '#38bdf8', icon: '🐋' },
  { min: 90, name: '浪潮专家', color: '#22d3ee', icon: '🐬' },
  { min: 30, name: '潮汐进阶', color: '#34d399', icon: '🐠' },
  { min: 7, name: '入海新手', color: '#0ea5e9', icon: '🐟' },
  { min: 0, name: '初探海域', color: 'rgba(255, 255, 255, 0.6)', icon: '🐚' }
]

onMounted(async () => {
  await loadUsers()
})

async function loadUsers() {
  loading.value = true
  try {
    const res: UsersResponse = await getUsers(currentPage.value, pageSize.value)
    users.value = res.users
    total.value = res.total
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.msg || '获取用户列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  loadUsers()
}

// 根据打卡次数获取称号
function getTitleByCheckin(total: number) {
  for (const title of titleConfig) {
    if (total >= title.min) {
      return title
    }
  }
  return titleConfig[titleConfig.length - 1]!
}

// 获取用户显示的称号
function getUserTitle(user: UserInfo) {
  if (user.title) {
    const found = titleConfig.find(t => t.name === user.title)
    if (found) return found
  }
  return getTitleByCheckin(user.total_checkin || 0)
}

// 过滤用户
const filteredUsers = computed(() => {
  let result = users.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(user =>
      user.username.toLowerCase().includes(query) ||
      user.display_name.toLowerCase().includes(query)
    )
  }

  // 管理员筛选
  if (filterAdmin.value !== '') {
    const isAdmin = filterAdmin.value === 'true'
    result = result.filter(user => user.is_admin === isAdmin)
  }

  // 称号筛选
  if (filterTitle.value) {
    result = result.filter(user => {
      const userTitle = getUserTitle(user)
      return userTitle.name === filterTitle.value
    })
  }

  return result
})

// 删除用户
async function handleDeleteUser(user: UserInfo) {
  if (user.id === userStore.user?.id) {
    ElMessage.warning('不能删除自己的账号')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${user.display_name || user.username}" 吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    loading.value = true
    await deleteUser(user.id)
    ElMessage.success('用户已删除')
    await loadUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.msg || '删除失败')
    }
  } finally {
    loading.value = false
  }
}

// 切换管理员状态
async function handleToggleAdmin(user: UserInfo) {
  if (user.id === userStore.user?.id) {
    ElMessage.warning('不能修改自己的管理员状态')
    return
  }

  const action = user.is_admin ? '取消' : '设置'
  try {
    await ElMessageBox.confirm(
      `确定要${action}用户 "${user.display_name || user.username}" 的管理员权限吗？`,
      '权限修改',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    loading.value = true
    await setUserAdmin(user.id, !user.is_admin)
    ElMessage.success(`已${action}管理员权限`)
    await loadUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.msg || '操作失败')
    }
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

// 打开编辑对话框
function openEditDialog(user: UserInfo) {
  editingUser.value = user
  editForm.value = {
    streak: user.streak || 0,
    max_streak: user.max_streak || 0,
    total_checkin: user.total_checkin || 0,
    title: user.title || ''
  }
  showEditDialog.value = true
}

// 保存用户数据
async function handleSaveUserStats() {
  if (!editingUser.value) return

  editLoading.value = true
  try {
    await updateUserStats(editingUser.value.id, {
      streak: editForm.value.streak,
      max_streak: editForm.value.max_streak,
      total_checkin: editForm.value.total_checkin,
      title: editForm.value.title
    })
    ElMessage.success('用户数据更新成功')
    showEditDialog.value = false
    await loadUsers()
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.msg || '更新失败')
  } finally {
    editLoading.value = false
  }
}

// 计算编辑表单中的称号预览
const previewTitle = computed(() => {
  if (editForm.value.title) {
    const selected = titleConfig.find(t => t.name === editForm.value.title)
    if (selected) return selected
  }
  return getTitleByCheckin(editForm.value.total_checkin)
})

// 导出用户数据为 CSV
function exportToCSV() {
  const headers = ['ID', '用户名', '显示名称', '称号', '连续打卡', '最高连续', '累计打卡', '管理员', '注册时间']
  const rows = filteredUsers.value.map(user => [
    user.id,
    user.username,
    user.display_name || '',
    getUserTitle(user).name,
    user.streak || 0,
    user.max_streak || 0,
    user.total_checkin || 0,
    user.is_admin ? '是' : '否',
    formatDate(user.created_at)
  ])

  const csvContent = [headers, ...rows]
    .map(row => row.map(cell => `"${cell}"`).join(','))
    .join('\n')

  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `users_${new Date().toISOString().split('T')[0]}.csv`
  link.click()
  URL.revokeObjectURL(url)

  ElMessage.success('导出成功')
}

// 清除筛选
function clearFilters() {
  searchQuery.value = ''
  filterAdmin.value = ''
  filterTitle.value = ''
}
</script>

<template>
  <div class="users-page">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <el-input
          v-model="searchQuery"
          placeholder="搜索用户名或显示名称"
          class="search-input"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <el-select
          v-model="filterAdmin"
          placeholder="管理员状态"
          clearable
          class="filter-select"
        >
          <el-option label="全部" value="" />
          <el-option label="管理员" value="true" />
          <el-option label="普通用户" value="false" />
        </el-select>

        <el-select
          v-model="filterTitle"
          placeholder="用户称号"
          clearable
          class="filter-select"
        >
          <el-option label="全部称号" value="" />
          <el-option
            v-for="t in titleConfig"
            :key="t.name"
            :value="t.name"
            :label="t.name"
          >
            <span :style="{ color: t.color }">{{ t.icon }} {{ t.name }}</span>
          </el-option>
        </el-select>

        <el-button v-if="searchQuery || filterAdmin || filterTitle" text @click="clearFilters">
          清除筛选
        </el-button>
      </div>

      <div class="toolbar-right">
        <el-button @click="exportToCSV">
          <el-icon><Download /></el-icon>
          导出 CSV
        </el-button>
        <el-button @click="loadUsers" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="stats-bar">
      <span>共 <strong>{{ total }}</strong> 个用户</span>
      <span v-if="filteredUsers.length !== users.length">
        ，筛选显示 <strong>{{ filteredUsers.length }}</strong> 个
      </span>
    </div>

    <!-- 用户表格 -->
    <el-card class="users-card" shadow="never">
      <el-table
        :data="filteredUsers"
        v-loading="loading"
        style="width: 100%"
        class="users-table"
      >
        <el-table-column prop="id" label="ID" width="70" />

        <el-table-column label="用户" min-width="200">
          <template #default="{ row }">
            <div class="user-cell">
              <UserAvatar
                :user-id="row.id"
                :username="row.username"
                :size="40"
                :border-radius="10"
                class="user-avatar"
              />
              <div class="user-info">
                <div class="user-name">
                  {{ row.display_name || row.username }}
                  <el-tag v-if="row.is_admin" type="warning" size="small">管理员</el-tag>
                </div>
                <div class="user-username">@{{ row.username }}</div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="称号" width="140">
          <template #default="{ row }">
            <span
              class="user-title"
              :style="{ color: getUserTitle(row).color }"
            >
              {{ getUserTitle(row).icon }} {{ getUserTitle(row).name }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="打卡统计" width="180">
          <template #default="{ row }">
            <div class="stats-cell">
              <span>连续 {{ row.streak || 0 }} 天</span>
              <span class="stats-divider">|</span>
              <span>累计 {{ row.total_checkin || 0 }} 次</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="注册时间" width="120">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <div class="action-cell">
              <el-button
                size="small"
                type="info"
                @click="openEditDialog(row)"
              >
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button
                size="small"
                :type="row.is_admin ? 'warning' : 'primary'"
                @click="handleToggleAdmin(row)"
                :disabled="row.id === userStore.user?.id"
              >
                <el-icon><Key /></el-icon>
                {{ row.is_admin ? '取消管理员' : '设为管理员' }}
              </el-button>
              <el-button
                size="small"
                type="danger"
                @click="handleDeleteUser(row)"
                :disabled="row.id === userStore.user?.id"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 编辑用户数据对话框 -->
    <el-dialog
      v-model="showEditDialog"
      title=""
      width="500px"
      class="edit-dialog"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="edit-dialog-header">
          <div class="edit-user-info" v-if="editingUser">
            <UserAvatar
              :user-id="editingUser.id"
              :username="editingUser.username"
              :size="48"
              :border-radius="12"
              class="edit-user-avatar"
            />
            <div class="edit-user-text">
              <h3>编辑用户数据</h3>
              <p>{{ editingUser.display_name || editingUser.username }}</p>
            </div>
          </div>
        </div>
      </template>

      <div class="edit-form">
        <!-- 称号选择 -->
        <div class="form-item">
          <label>用户称号</label>
          <el-select
            v-model="editForm.title"
            placeholder="自动（根据打卡次数）"
            size="large"
            clearable
            class="title-select"
          >
            <el-option
              value=""
              label="自动（根据打卡次数）"
            >
              <span class="title-option">
                <span class="title-option-icon">🔄</span>
                <span>自动（根据打卡次数）</span>
              </span>
            </el-option>
            <el-option
              v-for="t in titleConfig"
              :key="t.name"
              :value="t.name"
              :label="t.name"
            >
              <span class="title-option" :style="{ color: t.color }">
                <span class="title-option-icon">{{ t.icon }}</span>
                <span>{{ t.name }}</span>
                <span class="title-option-min">({{ t.min }}+)</span>
              </span>
            </el-option>
          </el-select>
          <div class="form-hint">
            选择"自动"则根据累计打卡次数自动计算称号
          </div>
        </div>

        <!-- 称号预览 -->
        <div class="title-preview">
          <span class="preview-label">显示效果</span>
          <span class="preview-title" :style="{ color: previewTitle.color }">
            {{ previewTitle.icon }} {{ previewTitle.name }}
          </span>
        </div>

        <div class="form-item">
          <label>连续打卡天数</label>
          <el-input-number
            v-model="editForm.streak"
            :min="0"
            :max="9999"
            size="large"
            controls-position="right"
          />
        </div>

        <div class="form-item">
          <label>最高连续天数</label>
          <el-input-number
            v-model="editForm.max_streak"
            :min="0"
            :max="9999"
            size="large"
            controls-position="right"
          />
        </div>

        <div class="form-item">
          <label>累计打卡次数</label>
          <el-input-number
            v-model="editForm.total_checkin"
            :min="0"
            :max="99999"
            size="large"
            controls-position="right"
          />
          <div class="form-hint" v-if="!editForm.title">
            当前自动称号：{{ getTitleByCheckin(editForm.total_checkin).icon }} {{ getTitleByCheckin(editForm.total_checkin).name }}
          </div>
        </div>
      </div>

      <template #footer>
        <div class="edit-dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" :loading="editLoading" @click="handleSaveUserStats">
            保存修改
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.users-page {
  max-width: 1400px;
  margin: 0 auto;
}

/* ===== 工具栏 ===== */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.toolbar-right {
  display: flex;
  gap: 12px;
}

.search-input {
  width: 240px;
}

.filter-select {
  width: 140px;
}

.search-input :deep(.el-input__wrapper),
.filter-select :deep(.el-select__wrapper) {
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.15);
  box-shadow: none;
}

.search-input :deep(.el-input__wrapper:hover),
.filter-select :deep(.el-select__wrapper:hover) {
  border-color: rgba(56, 189, 248, 0.3);
}

.search-input :deep(.el-input__wrapper.is-focus),
.filter-select :deep(.el-select__wrapper.is-focused) {
  border-color: rgb(var(--ocean-surface));
}

.search-input :deep(.el-input__inner) {
  color: #fff;
}

/* ===== 统计栏 ===== */
.stats-bar {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 16px;
}

.stats-bar strong {
  color: rgb(var(--ocean-surface));
}

/* ===== 用户卡片 ===== */
.users-card {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.1) !important;
  border-radius: 16px !important;
}

.users-card :deep(.el-card__body) {
  padding: 0;
}

.users-table {
  background: transparent !important;
}

.users-table :deep(.el-table__header-wrapper th) {
  background: rgba(56, 189, 248, 0.05) !important;
  color: rgba(255, 255, 255, 0.8);
  border-bottom: 1px solid rgba(56, 189, 248, 0.1);
}

.users-table :deep(.el-table__body-wrapper tr) {
  background: transparent !important;
}

.users-table :deep(.el-table__body-wrapper td) {
  border-bottom: 1px solid rgba(56, 189, 248, 0.08);
  color: rgba(255, 255, 255, 0.8);
}

.users-table :deep(.el-table__body-wrapper tr:hover > td) {
  background: rgba(56, 189, 248, 0.05) !important;
}

/* ===== 用户单元格 ===== */
.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  min-width: 0;
}

.user-name {
  font-weight: 600;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-username {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

/* ===== 称号 ===== */
.user-title {
  font-size: 13px;
  font-weight: 500;
}

/* ===== 统计单元格 ===== */
.stats-cell {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
}

.stats-divider {
  margin: 0 8px;
  color: rgba(56, 189, 248, 0.3);
}

/* ===== 操作单元格 ===== */
.action-cell {
  display: flex;
  gap: 8px;
}

/* ===== 分页 ===== */
.pagination-wrapper {
  padding: 20px;
  display: flex;
  justify-content: center;
  border-top: 1px solid rgba(56, 189, 248, 0.08);
}

.pagination-wrapper :deep(.el-pagination button),
.pagination-wrapper :deep(.el-pager li) {
  background: transparent;
  color: rgba(255, 255, 255, 0.6);
}

.pagination-wrapper :deep(.el-pager li.is-active) {
  color: rgb(var(--ocean-surface));
}

/* ===== 编辑对话框 ===== */
.edit-dialog :deep(.el-dialog) {
  background: linear-gradient(180deg, rgba(12, 20, 38, 0.98), rgba(8, 15, 30, 0.98)) !important;
  backdrop-filter: blur(30px);
  border: 1px solid rgba(56, 189, 248, 0.12);
  border-radius: 20px;
}

.edit-dialog :deep(.el-dialog__header) {
  padding: 0;
  margin: 0;
}

.edit-dialog :deep(.el-dialog__body) {
  padding: 0 24px 24px;
}

.edit-dialog :deep(.el-dialog__footer) {
  padding: 0 24px 24px;
}

.edit-dialog-header {
  padding: 24px 24px 20px;
  border-bottom: 1px solid rgba(56, 189, 248, 0.1);
}

.edit-user-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.edit-user-text h3 {
  margin: 0 0 4px;
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.edit-user-text p {
  margin: 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-top: 20px;
}

.title-preview {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: 12px;
}

.preview-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.preview-title {
  font-size: 16px;
  font-weight: 600;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-item label {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
}

.form-item :deep(.el-input-number) {
  width: 100%;
}

.form-item :deep(.el-input-number .el-input__wrapper) {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(56, 189, 248, 0.15);
  box-shadow: none;
}

.form-item :deep(.el-input-number .el-input__wrapper:hover) {
  border-color: rgba(56, 189, 248, 0.3);
}

.form-item :deep(.el-input-number .el-input__wrapper.is-focus) {
  border-color: rgb(var(--ocean-surface));
}

.form-item :deep(.el-input-number .el-input__inner) {
  color: #fff;
}

.form-hint {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

/* ===== 称号选择 ===== */
.title-select {
  width: 100%;
}

.title-select :deep(.el-select__wrapper) {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(56, 189, 248, 0.15);
  box-shadow: none;
}

.title-select :deep(.el-select__wrapper:hover) {
  border-color: rgba(56, 189, 248, 0.3);
}

.title-select :deep(.el-select__wrapper.is-focused) {
  border-color: rgb(var(--ocean-surface));
}

.title-select :deep(.el-select__selected-item) {
  color: #fff;
}

.title-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-option-icon {
  font-size: 16px;
}

.title-option-min {
  margin-left: auto;
  font-size: 12px;
  opacity: 0.6;
}

.edit-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.edit-dialog-footer .el-button--primary {
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .toolbar-left,
  .toolbar-right {
    flex-wrap: wrap;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }

  .action-cell {
    flex-wrap: wrap;
  }
}
</style>
