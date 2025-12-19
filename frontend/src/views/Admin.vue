<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useUserStore } from '@/store/user'
import { getUsers, deleteUser, setUserAdmin, type UsersResponse } from '@/api/admin'
import type { UserInfo } from '@/api/auth'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Key, ArrowLeft, Search, Refresh } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const users = ref<UserInfo[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')

onMounted(async () => {
  await userStore.fetchProfile()
  if (!userStore.user?.is_admin) {
    ElMessage.error('无权访问管理后台')
    router.push('/dashboard')
    return
  }
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

// 过滤用户
const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value
  const query = searchQuery.value.toLowerCase()
  return users.value.filter(user =>
    user.username.toLowerCase().includes(query) ||
    user.display_name.toLowerCase().includes(query)
  )
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
</script>

<template>
  <MainLayout>
    <div class="admin-page">
      <!-- Header -->
      <div class="admin-header">
        <div class="header-left">
          <RouterLink to="/dashboard" class="back-link">
            <el-icon><ArrowLeft /></el-icon>
            返回用户中心
          </RouterLink>
          <h1>管理后台</h1>
          <p>管理用户和系统设置</p>
        </div>
        <div class="header-stats">
          <div class="stat-item">
            <span class="stat-value">{{ total }}</span>
            <span class="stat-label">总用户数</span>
          </div>
        </div>
      </div>

      <!-- Toolbar -->
      <div class="admin-toolbar">
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
        <el-button @click="loadUsers" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>

      <!-- Users Table -->
      <el-card class="users-card" shadow="never">
        <el-table
          :data="filteredUsers"
          v-loading="loading"
          style="width: 100%"
          class="users-table"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column label="用户" min-width="200">
            <template #default="{ row }">
              <div class="user-cell">
                <div class="user-avatar">
                  {{ row.username?.[0]?.toUpperCase() || '?' }}
                </div>
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
          <el-table-column label="打卡统计" width="180">
            <template #default="{ row }">
              <div class="stats-cell">
                <span>连续 {{ row.streak }} 天</span>
                <span class="stats-divider">|</span>
                <span>累计 {{ row.total_checkin }} 次</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="注册时间" width="120">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <div class="action-cell">
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

        <!-- Pagination -->
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
    </div>
  </MainLayout>
</template>

<style scoped>
.admin-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 24px 80px;
  min-height: calc(100vh - 72px);
}

/* ===== Header ===== */
.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
  flex-wrap: wrap;
  gap: 24px;
}

.header-left {
  flex: 1;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  font-size: 14px;
  margin-bottom: 12px;
  transition: color 0.2s;
}

.back-link:hover {
  color: rgb(var(--ocean-surface));
}

.admin-header h1 {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
}

.admin-header p {
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}

.header-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  text-align: center;
  padding: 16px 24px;
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-lg);
}

.stat-value {
  display: block;
  font-size: 32px;
  font-weight: 700;
  color: rgb(var(--ocean-surface));
}

.stat-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

/* ===== Toolbar ===== */
.admin-toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.search-input {
  max-width: 300px;
}

.search-input :deep(.el-input__wrapper) {
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.15);
  box-shadow: none;
}

.search-input :deep(.el-input__wrapper:hover) {
  border-color: rgba(56, 189, 248, 0.3);
}

.search-input :deep(.el-input__wrapper.is-focus) {
  border-color: rgb(var(--ocean-surface));
}

.search-input :deep(.el-input__inner) {
  color: #fff;
}

/* ===== Users Card ===== */
.users-card {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.1) !important;
  border-radius: var(--radius-xl) !important;
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

/* ===== User Cell ===== */
.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  flex-shrink: 0;
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

/* ===== Stats Cell ===== */
.stats-cell {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
}

.stats-divider {
  margin: 0 8px;
  color: rgba(56, 189, 248, 0.3);
}

/* ===== Action Cell ===== */
.action-cell {
  display: flex;
  gap: 8px;
}

/* ===== Pagination ===== */
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
</style>
