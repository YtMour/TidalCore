<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Heatmap from '@/components/Heatmap.vue'
import { useUserStore } from '@/store/user'
import { getHeatmap, getHistory, type CheckinRecord } from '@/api/checkin'
import { updateProfile, updateUsername, updatePassword } from '@/api/auth'
import { Timer, Calendar, Clock, CircleCheck, Pointer, Trophy, Aim, Edit, Lock, User, Setting } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()

const heatmapData = ref<Record<string, number>>({})
const recentHistory = ref<CheckinRecord[]>([])
const loading = ref(true)
const mounted = ref(false)

// 设置相关
const showSettingsDialog = ref(false)
const activeSettingTab = ref('profile')
const settingsLoading = ref(false)

// 表单数据
const profileForm = ref({
  display_name: ''
})
const usernameForm = ref({
  username: ''
})
const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

// 打开设置对话框
function openSettings() {
  profileForm.value.display_name = userStore.user?.display_name || ''
  usernameForm.value.username = userStore.user?.username || ''
  passwordForm.value = { old_password: '', new_password: '', confirm_password: '' }
  showSettingsDialog.value = true
}

// 更新显示名称
async function handleUpdateProfile() {
  if (!profileForm.value.display_name.trim()) {
    ElMessage.warning('请输入显示名称')
    return
  }
  settingsLoading.value = true
  try {
    await updateProfile({ display_name: profileForm.value.display_name })
    await userStore.refreshUser()
    ElMessage.success('显示名称更新成功')
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.msg || '更新失败')
  } finally {
    settingsLoading.value = false
  }
}

// 更新用户名
async function handleUpdateUsername() {
  if (!usernameForm.value.username.trim()) {
    ElMessage.warning('请输入用户名')
    return
  }
  if (!/^[a-zA-Z0-9_]+$/.test(usernameForm.value.username)) {
    ElMessage.warning('用户名只能包含字母、数字和下划线')
    return
  }
  settingsLoading.value = true
  try {
    await updateUsername({ username: usernameForm.value.username })
    await userStore.refreshUser()
    ElMessage.success('用户名更新成功')
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.msg || '更新失败')
  } finally {
    settingsLoading.value = false
  }
}

// 更新密码
async function handleUpdatePassword() {
  if (!passwordForm.value.old_password) {
    ElMessage.warning('请输入原密码')
    return
  }
  if (!passwordForm.value.new_password) {
    ElMessage.warning('请输入新密码')
    return
  }
  if (passwordForm.value.new_password.length < 6) {
    ElMessage.warning('新密码长度至少6位')
    return
  }
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }
  settingsLoading.value = true
  try {
    await updatePassword({
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password
    })
    ElMessage.success('密码更新成功')
    passwordForm.value = { old_password: '', new_password: '', confirm_password: '' }
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.msg || '更新失败')
  } finally {
    settingsLoading.value = false
  }
}

// 称号配置 - 海洋主题（含特效）
const titleConfig = [
  { min: 1000, name: '海神降临', color: '#f472b6', bg: 'linear-gradient(135deg, rgba(244, 114, 182, 0.25), rgba(236, 72, 153, 0.1))', icon: '🔱', effect: 'effect-divine' },
  { min: 730, name: '深渊霸主', color: '#a78bfa', bg: 'linear-gradient(135deg, rgba(167, 139, 250, 0.25), rgba(139, 92, 246, 0.1))', icon: '🦑', effect: 'effect-abyss' },
  { min: 365, name: '深海传奇', color: '#fbbf24', bg: 'linear-gradient(135deg, rgba(251, 191, 36, 0.25), rgba(245, 158, 11, 0.1))', icon: '🌊', effect: 'effect-legend' },
  { min: 180, name: '海洋大师', color: '#38bdf8', bg: 'linear-gradient(135deg, rgba(56, 189, 248, 0.25), rgba(14, 165, 233, 0.1))', icon: '🐋', effect: 'effect-master' },
  { min: 90, name: '浪潮专家', color: '#22d3ee', bg: 'linear-gradient(135deg, rgba(34, 211, 238, 0.25), rgba(6, 182, 212, 0.1))', icon: '🐬', effect: 'effect-expert' },
  { min: 30, name: '潮汐进阶', color: '#34d399', bg: 'linear-gradient(135deg, rgba(52, 211, 153, 0.25), rgba(16, 185, 129, 0.1))', icon: '🐠', effect: 'effect-advanced' },
  { min: 7, name: '入海新手', color: '#0ea5e9', bg: 'linear-gradient(135deg, rgba(14, 165, 233, 0.25), rgba(2, 132, 199, 0.1))', icon: '🐟', effect: 'effect-beginner' },
  { min: 0, name: '初探海域', color: 'rgba(255, 255, 255, 0.6)', bg: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.05))', icon: '🐚', effect: '' }
]

// 计算用户当前称号
const userLevel = computed(() => {
  // 优先使用数据库中存储的称号
  const userTitle = userStore.user?.title
  if (userTitle) {
    const customTitle = titleConfig.find(t => t.name === userTitle)
    if (customTitle) return customTitle
  }
  // 否则根据打卡次数自动计算
  const total = userStore.user?.total_checkin || 0
  for (const title of titleConfig) {
    if (total >= title.min) {
      return title
    }
  }
  return titleConfig[titleConfig.length - 1]!
})

// 判断是否使用自定义称号
const isCustomTitle = computed(() => {
  return !!userStore.user?.title
})

// 计算下一个称号及进度
const nextTitle = computed(() => {
  // 如果是自定义称号，不显示下一个称号进度
  if (isCustomTitle.value) {
    return null
  }
  const total = userStore.user?.total_checkin || 0
  const currentIndex = titleConfig.findIndex(t => total >= t.min)
  if (currentIndex <= 0) {
    return null // 已达最高称号
  }
  const next = titleConfig[currentIndex - 1]!
  const current = titleConfig[currentIndex]!
  const progress = ((total - current.min) / (next.min - current.min)) * 100
  return {
    ...next,
    required: next.min,
    progress: Math.min(progress, 100),
    remaining: next.min - total
  }
})

onMounted(async () => {
  setTimeout(() => {
    mounted.value = true
  }, 100)

  await userStore.fetchProfile()
  await loadData()
})

async function loadData() {
  loading.value = true
  try {
    const [checkins, history] = await Promise.all([
      getHeatmap(365),
      getHistory(10)
    ])

    const heatmap: Record<string, number> = {}
    checkins.forEach(c => {
      const datePart = c.checked_at.split('T')[0]
      if (datePart) {
        heatmap[datePart] = (heatmap[datePart] || 0) + 1
      }
    })
    heatmapData.value = heatmap
    recentHistory.value = history
  } catch {
    // 静默处理
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return mins > 0 ? `${mins}分${secs}秒` : `${secs}秒`
}

function getRelativeTime(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  return formatDate(dateStr)
}
</script>

<template>
  <MainLayout>
    <div class="dashboard-page">
      <!-- Profile Header -->
      <section class="profile-section" :class="{ mounted }">
        <el-card class="profile-card" shadow="never">
          <!-- Decorative background - 海洋主题 -->
          <div class="profile-decoration profile-decoration-1"></div>
          <div class="profile-decoration profile-decoration-2"></div>
          <div class="profile-wave"></div>

          <div class="profile-content">
            <!-- Avatar with level badge - 海洋风格 -->
            <div class="avatar-wrapper">
              <div class="avatar-rings">
                <div class="avatar-ring avatar-ring-1"></div>
                <div class="avatar-ring avatar-ring-2"></div>
              </div>
              <div class="avatar">
                {{ userStore.user?.username?.[0]?.toUpperCase() || '?' }}
              </div>
            </div>

            <!-- User Info -->
            <div class="user-info">
              <div class="user-name">
                <h1>{{ userStore.user?.display_name || userStore.user?.username }}</h1>
                <div class="wave-icon">
                  <svg viewBox="0 0 24 24" fill="none" width="20" height="20">
                    <path d="M2 12C2 12 5 8 8 12C11 16 14 8 17 12C20 16 22 12 22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  </svg>
                </div>
                <!-- 称号徽章 -->
                <el-tooltip
                  :disabled="!nextTitle"
                  placement="bottom"
                  effect="dark"
                >
                  <template #content>
                    <div class="title-tooltip" v-if="nextTitle">
                      <div class="tooltip-header">
                        <span>下一称号</span>
                        <span class="next-title" :style="{ color: nextTitle.color }">
                          {{ nextTitle.icon }} {{ nextTitle.name }}
                        </span>
                      </div>
                      <div class="tooltip-progress">
                        <div class="progress-bar">
                          <div class="progress-fill" :style="{ width: nextTitle.progress + '%', background: nextTitle.color }"></div>
                        </div>
                        <span class="progress-text">还需 {{ nextTitle.remaining }} 次打卡</span>
                      </div>
                    </div>
                  </template>
                  <div
                    class="title-badge"
                    :class="userLevel.effect"
                    :style="{ background: userLevel.bg, color: userLevel.color, borderColor: userLevel.color }"
                  >
                    <span class="title-icon">{{ userLevel.icon }}</span>
                    <span class="title-name">{{ userLevel.name }}</span>
                  </div>
                </el-tooltip>
              </div>
              <p class="user-motto">@{{ userStore.user?.username }}</p>

              <!-- Quick stats inline - 海洋色彩 -->
              <div class="inline-stats">
                <span class="inline-stat ocean">
                  <el-icon><Pointer /></el-icon>
                  {{ userStore.user?.streak || 0 }} 天连续
                </span>
                <span class="inline-stat aqua">
                  <el-icon><Trophy /></el-icon>
                  {{ userStore.user?.max_streak || 0 }} 最高
                </span>
                <!-- 操作按钮 -->
                <el-button class="inline-action-btn settings" size="small" @click="openSettings">
                  <el-icon :size="14"><Setting /></el-icon>
                  <span>账号设置</span>
                </el-button>
                <RouterLink v-if="userStore.user?.is_admin" to="/admin">
                  <el-button class="inline-action-btn admin" size="small">
                    <el-icon :size="14"><Setting /></el-icon>
                    <span>管理后台</span>
                  </el-button>
                </RouterLink>
              </div>
            </div>
          </div>
        </el-card>
      </section>

      <!-- Stats Cards - 海洋主题 -->
      <section class="stats-section" :class="{ mounted }">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="8">
            <div class="stat-card stat-ocean">
              <div class="stat-glow"></div>
              <div class="stat-icon">
                <el-icon :size="28" color="#38bdf8"><Pointer /></el-icon>
              </div>
              <div class="stat-value">{{ userStore.user?.streak || 0 }}</div>
              <div class="stat-label">连续打卡</div>
              <div class="stat-unit">天</div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="stat-card stat-aqua">
              <div class="stat-glow"></div>
              <div class="stat-icon">
                <el-icon :size="28" color="#22d3ee"><Trophy /></el-icon>
              </div>
              <div class="stat-value">{{ userStore.user?.max_streak || 0 }}</div>
              <div class="stat-label">最高记录</div>
              <div class="stat-unit">天</div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="stat-card stat-green">
              <div class="stat-glow"></div>
              <div class="stat-icon">
                <el-icon :size="28" color="#34d399"><Aim /></el-icon>
              </div>
              <div class="stat-value">{{ userStore.user?.total_checkin || 0 }}</div>
              <div class="stat-label">累计打卡</div>
              <div class="stat-unit">次</div>
            </div>
          </el-col>
        </el-row>
      </section>

      <!-- Heatmap Section - 海洋主题 -->
      <section class="heatmap-section" :class="{ mounted }">
        <el-card class="section-card" shadow="never">
          <div class="card-decoration"></div>

          <div class="section-header">
            <div class="header-icon">
              <el-icon :size="20" color="#38bdf8"><Calendar /></el-icon>
            </div>
            <div class="header-text">
              <h2>我的打卡热力图</h2>
              <p>记录每一次潮汐的痕迹</p>
            </div>
          </div>

          <div v-if="loading" class="loading-content">
            <el-skeleton :rows="3" animated />
          </div>
          <Heatmap v-else :data="heatmapData" :days="365" />
        </el-card>
      </section>

      <!-- Recent History - 海洋主题 -->
      <section class="history-section" :class="{ mounted }">
        <el-card class="section-card" shadow="never">
          <div class="card-decoration"></div>

          <div class="section-header">
            <div class="header-icon">
              <el-icon :size="20" color="#22d3ee"><Clock /></el-icon>
            </div>
            <div class="header-text">
              <h2>最近打卡记录</h2>
              <p>浪花般的训练足迹</p>
            </div>
          </div>

          <div v-if="loading" class="loading-content">
            <div v-for="i in 5" :key="i" class="history-skeleton">
              <el-skeleton-item variant="rect" style="width: 48px; height: 48px; border-radius: 12px" />
              <div class="skeleton-text">
                <el-skeleton-item variant="text" style="width: 60%" />
                <el-skeleton-item variant="text" style="width: 40%; height: 12px" />
              </div>
            </div>
          </div>

          <div v-else-if="recentHistory.length === 0" class="empty-state">
            <div class="empty-icon">
              <svg viewBox="0 0 48 48" fill="none" width="48" height="48">
                <circle cx="24" cy="24" r="20" stroke="rgba(56, 189, 248, 0.3)" stroke-width="2" fill="none"/>
                <path d="M12 24C12 24 16 18 20 24C24 30 28 18 32 24C36 30 40 24 40 24" stroke="rgba(56, 189, 248, 0.5)" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <p class="empty-title">暂无打卡记录</p>
            <p class="empty-desc">完成训练后，浪花将汇聚于此</p>
            <RouterLink to="/train">
              <el-button type="primary" class="empty-btn">
                <el-icon><Timer /></el-icon>
                开始第一次训练
              </el-button>
            </RouterLink>
          </div>

          <div v-else class="history-list">
            <TransitionGroup name="list">
              <div
                v-for="(record, index) in recentHistory"
                :key="record.id"
                class="history-item"
                :style="{ transitionDelay: `${index * 50}ms` }"
              >
                <div class="history-icon">
                  <el-icon :size="24" color="#34d399"><CircleCheck /></el-icon>
                </div>
                <div class="history-info">
                  <div class="history-time">{{ getRelativeTime(record.checked_at) }}</div>
                  <div class="history-detail">
                    <span class="detail-cycles">{{ record.cycles }} 个循环</span>
                    <span class="detail-dot">·</span>
                    <span class="detail-duration">{{ formatDuration(record.duration) }}</span>
                  </div>
                </div>
                <div class="history-date">
                  {{ formatDate(record.checked_at) }}
                </div>
              </div>
            </TransitionGroup>
          </div>
        </el-card>
      </section>
    </div>

    <!-- 设置对话框 - 重新设计 -->
    <el-dialog
      v-model="showSettingsDialog"
      title=""
      width="480px"
      class="settings-dialog-new"
      :close-on-click-modal="false"
      :show-close="true"
    >
      <template #header>
        <div class="settings-header">
          <div class="settings-header-icon">
            <el-icon :size="24"><Setting /></el-icon>
          </div>
          <div class="settings-header-text">
            <h3>账号设置</h3>
            <p>管理您的个人信息和安全设置</p>
          </div>
        </div>
      </template>

      <div class="settings-content">
        <el-tabs v-model="activeSettingTab" class="settings-tabs-new">
          <el-tab-pane name="profile">
            <template #label>
              <div class="tab-label">
                <el-icon><User /></el-icon>
                <span>显示名称</span>
              </div>
            </template>
            <div class="setting-form-new">
              <div class="form-header">
                <h4>修改显示名称</h4>
                <p>显示名称将在排行榜和个人主页中展示</p>
              </div>
              <div class="form-field">
                <label>显示名称</label>
                <el-input
                  v-model="profileForm.display_name"
                  placeholder="请输入显示名称"
                  maxlength="50"
                  show-word-limit
                  size="large"
                >
                  <template #prefix>
                    <el-icon><User /></el-icon>
                  </template>
                </el-input>
              </div>
              <el-button
                type="primary"
                :loading="settingsLoading"
                @click="handleUpdateProfile"
                class="submit-btn-new"
                size="large"
              >
                保存修改
              </el-button>
            </div>
          </el-tab-pane>

          <el-tab-pane name="username">
            <template #label>
              <div class="tab-label">
                <el-icon><Edit /></el-icon>
                <span>用户名</span>
              </div>
            </template>
            <div class="setting-form-new">
              <div class="form-header">
                <h4>修改用户名</h4>
                <p>用户名用于登录，只能包含字母、数字和下划线</p>
              </div>
              <div class="form-field">
                <label>新用户名</label>
                <el-input
                  v-model="usernameForm.username"
                  placeholder="请输入新用户名"
                  maxlength="50"
                  size="large"
                >
                  <template #prefix>
                    <el-icon><Edit /></el-icon>
                  </template>
                </el-input>
              </div>
              <el-button
                type="primary"
                :loading="settingsLoading"
                @click="handleUpdateUsername"
                class="submit-btn-new"
                size="large"
              >
                保存修改
              </el-button>
            </div>
          </el-tab-pane>

          <el-tab-pane name="password">
            <template #label>
              <div class="tab-label">
                <el-icon><Lock /></el-icon>
                <span>修改密码</span>
              </div>
            </template>
            <div class="setting-form-new">
              <div class="form-header">
                <h4>修改密码</h4>
                <p>为了账号安全，请定期更换密码</p>
              </div>
              <div class="form-field">
                <label>原密码</label>
                <el-input
                  v-model="passwordForm.old_password"
                  type="password"
                  placeholder="请输入原密码"
                  show-password
                  size="large"
                >
                  <template #prefix>
                    <el-icon><Lock /></el-icon>
                  </template>
                </el-input>
              </div>
              <div class="form-field">
                <label>新密码</label>
                <el-input
                  v-model="passwordForm.new_password"
                  type="password"
                  placeholder="请输入新密码（至少6位）"
                  show-password
                  size="large"
                >
                  <template #prefix>
                    <el-icon><Lock /></el-icon>
                  </template>
                </el-input>
              </div>
              <div class="form-field">
                <label>确认新密码</label>
                <el-input
                  v-model="passwordForm.confirm_password"
                  type="password"
                  placeholder="请再次输入新密码"
                  show-password
                  size="large"
                >
                  <template #prefix>
                    <el-icon><Lock /></el-icon>
                  </template>
                </el-input>
              </div>
              <el-button
                type="primary"
                :loading="settingsLoading"
                @click="handleUpdatePassword"
                class="submit-btn-new"
                size="large"
              >
                修改密码
              </el-button>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-dialog>
  </MainLayout>
</template>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 32px;
  max-width: 900px;
  margin: 0 auto;
  padding: 40px 24px 80px;
  min-height: calc(100vh - 72px);
}

/* ===== Profile Section - 海洋主题 ===== */
.profile-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth);
}

.profile-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.profile-card {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.1) !important;
  border-radius: var(--radius-2xl) !important;
  overflow: hidden;
  position: relative;
}

.profile-card :deep(.el-card__body) {
  padding: 32px;
}

.profile-decoration {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
}

.profile-decoration-1 {
  top: -50px;
  right: -50px;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.2), transparent 70%);
}

.profile-decoration-2 {
  bottom: -50px;
  left: -50px;
  width: 250px;
  height: 250px;
  background: radial-gradient(circle, rgba(34, 211, 238, 0.15), transparent 70%);
}

.profile-wave {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 80px;
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 100'%3E%3Cpath fill='rgba(56, 189, 248, 0.05)' d='M0,50 C360,100 720,0 1080,50 C1260,75 1380,25 1440,50 L1440,100 L0,100 Z'/%3E%3C/svg%3E");
  background-size: 1440px 100px;
  background-repeat: repeat-x;
  animation: wave-move 15s linear infinite;
  pointer-events: none;
}

@keyframes wave-move {
  0% { background-position-x: 0; }
  100% { background-position-x: 1440px; }
}

.profile-content {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

@media (min-width: 640px) {
  .profile-content {
    flex-direction: row;
    align-items: flex-start;
  }
}

.avatar-wrapper {
  position: relative;
  flex-shrink: 0;
}

.avatar-rings {
  position: absolute;
  inset: -12px;
  pointer-events: none;
}

.avatar-ring {
  position: absolute;
  inset: 0;
  border-radius: 30px;
  border: 2px solid;
  animation: ring-pulse 3s ease-in-out infinite;
}

.avatar-ring-1 {
  border-color: rgba(56, 189, 248, 0.3);
  animation-delay: 0s;
}

.avatar-ring-2 {
  inset: -6px;
  border-radius: 34px;
  border-color: rgba(34, 211, 238, 0.2);
  animation-delay: 0.5s;
}

@keyframes ring-pulse {
  0%, 100% { transform: scale(1); opacity: 0.5; }
  50% { transform: scale(1.05); opacity: 1; }
}

.avatar {
  width: 112px;
  height: 112px;
  border-radius: 28px;
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: 700;
  color: #fff;
  box-shadow:
    0 20px 40px rgba(14, 165, 233, 0.3),
    0 0 60px rgba(56, 189, 248, 0.2);
  position: relative;
}

/* ===== Title Badge (称号徽章) ===== */
.title-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: var(--radius-full);
  font-size: 12px;
  font-weight: 600;
  border: 1px solid;
  backdrop-filter: blur(10px);
  white-space: nowrap;
}

.title-icon {
  font-size: 12px;
}

.title-name {
  font-size: 12px;
}

/* ===== Title Effects (称号特效) ===== */
/* 海神降临 - 神圣光芒 + 彩虹流光 */
.title-badge.effect-divine {
  position: relative;
  animation: divine-glow 2s ease-in-out infinite;
  box-shadow: 0 0 15px rgba(244, 114, 182, 0.5), 0 0 30px rgba(244, 114, 182, 0.3);
}

.title-badge.effect-divine::before {
  content: '';
  position: absolute;
  inset: -2px;
  border-radius: inherit;
  background: linear-gradient(90deg, #f472b6, #a78bfa, #38bdf8, #34d399, #fbbf24, #f472b6);
  background-size: 300% 100%;
  animation: rainbow-flow 3s linear infinite;
  z-index: -1;
  opacity: 0.8;
}

.title-badge.effect-divine .title-icon {
  animation: divine-icon 1.5s ease-in-out infinite;
}

@keyframes divine-glow {
  0%, 100% { box-shadow: 0 0 15px rgba(244, 114, 182, 0.5), 0 0 30px rgba(244, 114, 182, 0.3); }
  50% { box-shadow: 0 0 25px rgba(244, 114, 182, 0.7), 0 0 50px rgba(244, 114, 182, 0.4); }
}

@keyframes rainbow-flow {
  0% { background-position: 0% 50%; }
  100% { background-position: 300% 50%; }
}

@keyframes divine-icon {
  0%, 100% { transform: scale(1) rotate(0deg); }
  50% { transform: scale(1.2) rotate(10deg); }
}

/* 深渊霸主 - 暗紫脉动 + 触手波动 */
.title-badge.effect-abyss {
  position: relative;
  animation: abyss-pulse 2.5s ease-in-out infinite;
  box-shadow: 0 0 12px rgba(167, 139, 250, 0.4), inset 0 0 8px rgba(139, 92, 246, 0.3);
}

.title-badge.effect-abyss::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: radial-gradient(circle at 30% 50%, rgba(167, 139, 250, 0.3) 0%, transparent 50%);
  animation: abyss-tentacle 3s ease-in-out infinite;
}

.title-badge.effect-abyss .title-icon {
  animation: abyss-icon 2s ease-in-out infinite;
}

@keyframes abyss-pulse {
  0%, 100% { box-shadow: 0 0 12px rgba(167, 139, 250, 0.4), inset 0 0 8px rgba(139, 92, 246, 0.3); }
  50% { box-shadow: 0 0 20px rgba(167, 139, 250, 0.6), inset 0 0 15px rgba(139, 92, 246, 0.4); }
}

@keyframes abyss-tentacle {
  0%, 100% { background-position: 30% 50%; opacity: 0.5; }
  50% { background-position: 70% 50%; opacity: 0.8; }
}

@keyframes abyss-icon {
  0%, 100% { transform: scale(1); }
  25% { transform: scale(1.1) rotate(-5deg); }
  75% { transform: scale(1.1) rotate(5deg); }
}

/* 深海传奇 - 金色闪耀 + 波浪效果 */
.title-badge.effect-legend {
  position: relative;
  animation: legend-shine 2s ease-in-out infinite;
  box-shadow: 0 0 10px rgba(251, 191, 36, 0.4);
}

.title-badge.effect-legend::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.4), transparent);
  animation: legend-sweep 2.5s ease-in-out infinite;
}

.title-badge.effect-legend .title-icon {
  animation: legend-wave 1.5s ease-in-out infinite;
}

@keyframes legend-shine {
  0%, 100% { box-shadow: 0 0 10px rgba(251, 191, 36, 0.4); }
  50% { box-shadow: 0 0 20px rgba(251, 191, 36, 0.6), 0 0 30px rgba(251, 191, 36, 0.3); }
}

@keyframes legend-sweep {
  0% { left: -100%; }
  50%, 100% { left: 200%; }
}

@keyframes legend-wave {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-2px); }
}

/* 海洋大师 - 水波纹扩散 */
.title-badge.effect-master {
  position: relative;
  animation: master-glow 2s ease-in-out infinite;
}

.title-badge.effect-master::after {
  content: '';
  position: absolute;
  inset: -4px;
  border-radius: inherit;
  border: 1px solid rgba(56, 189, 248, 0.3);
  animation: master-ripple 2s ease-out infinite;
}

.title-badge.effect-master .title-icon {
  animation: master-swim 3s ease-in-out infinite;
}

@keyframes master-glow {
  0%, 100% { box-shadow: 0 0 8px rgba(56, 189, 248, 0.3); }
  50% { box-shadow: 0 0 15px rgba(56, 189, 248, 0.5); }
}

@keyframes master-ripple {
  0% { transform: scale(1); opacity: 0.6; }
  100% { transform: scale(1.3); opacity: 0; }
}

@keyframes master-swim {
  0%, 100% { transform: translateX(0); }
  50% { transform: translateX(2px); }
}

/* 浪潮专家 - 轻微波动 */
.title-badge.effect-expert {
  animation: expert-float 2.5s ease-in-out infinite;
}

.title-badge.effect-expert .title-icon {
  animation: expert-jump 2s ease-in-out infinite;
}

@keyframes expert-float {
  0%, 100% { transform: translateY(0); box-shadow: 0 0 6px rgba(34, 211, 238, 0.3); }
  50% { transform: translateY(-1px); box-shadow: 0 0 12px rgba(34, 211, 238, 0.4); }
}

@keyframes expert-jump {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-1px); }
}

/* 潮汐进阶 - 柔和呼吸 */
.title-badge.effect-advanced {
  animation: advanced-breathe 3s ease-in-out infinite;
}

.title-badge.effect-advanced .title-icon {
  animation: advanced-swim 2.5s ease-in-out infinite;
}

@keyframes advanced-breathe {
  0%, 100% { opacity: 0.9; box-shadow: 0 0 4px rgba(52, 211, 153, 0.2); }
  50% { opacity: 1; box-shadow: 0 0 8px rgba(52, 211, 153, 0.3); }
}

@keyframes advanced-swim {
  0%, 100% { transform: translateX(0); }
  50% { transform: translateX(1px); }
}

/* 入海新手 - 轻微闪烁 */
.title-badge.effect-beginner {
  animation: beginner-twinkle 3s ease-in-out infinite;
}

@keyframes beginner-twinkle {
  0%, 100% { opacity: 0.85; }
  50% { opacity: 1; }
}

/* ===== Title Tooltip ===== */
.title-tooltip {
  padding: 4px 0;
  min-width: 180px;
}

.tooltip-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-size: 13px;
}

.tooltip-header > span:first-child {
  color: rgba(255, 255, 255, 0.6);
}

.next-title {
  font-weight: 600;
}

.tooltip-progress {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.progress-bar {
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  text-align: right;
}

.user-info {
  flex: 1;
  text-align: center;
}

@media (min-width: 640px) {
  .user-info {
    text-align: left;
  }
}

.user-name {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 4px;
}

@media (min-width: 640px) {
  .user-name {
    justify-content: flex-start;
  }
}

.user-name h1 {
  font-size: 26px;
  font-weight: 700;
  color: #fff;
  margin: 0;
}

.wave-icon {
  color: rgb(var(--ocean-surface));
  animation: wave-flow 2s ease-in-out infinite;
}

@keyframes wave-flow {
  0%, 100% { transform: translateX(-2px); }
  50% { transform: translateX(2px); }
}

.user-motto {
  color: rgba(255, 255, 255, 0.65);
  margin: 0 0 12px;
  font-size: 14px;
}

.inline-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
  font-size: 14px;
}

@media (min-width: 640px) {
  .inline-stats {
    justify-content: flex-start;
  }
}

.inline-stat {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: var(--radius-md);
  background: rgba(56, 189, 248, 0.08);
  border: 1px solid rgba(56, 189, 248, 0.15);
}

.inline-stat.ocean {
  color: rgb(var(--ocean-surface));
}

.inline-stat.aqua {
  color: rgb(var(--aqua-glow));
}

/* ===== Inline Action Buttons ===== */
.inline-action-btn {
  display: inline-flex !important;
  align-items: center !important;
  gap: 6px !important;
  padding: 6px 12px !important;
  border-radius: var(--radius-md) !important;
  font-size: 14px !important;
  font-weight: 500 !important;
  height: auto !important;
  transition: all 0.3s var(--ease-smooth) !important;
}

.inline-action-btn.settings {
  background: rgba(56, 189, 248, 0.08) !important;
  border: 1px solid rgba(56, 189, 248, 0.15) !important;
  color: rgb(var(--ocean-surface)) !important;
}

.inline-action-btn.settings:hover {
  background: rgba(56, 189, 248, 0.15) !important;
  border-color: rgba(56, 189, 248, 0.3) !important;
  color: #fff !important;
}

.inline-action-btn.admin {
  background: rgba(251, 191, 36, 0.08) !important;
  border: 1px solid rgba(251, 191, 36, 0.15) !important;
  color: #fbbf24 !important;
}

.inline-action-btn.admin:hover {
  background: rgba(251, 191, 36, 0.15) !important;
  border-color: rgba(251, 191, 36, 0.3) !important;
}

/* ===== Stats Section - 海洋主题 ===== */
.stats-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth) 0.1s;
}

.stats-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.stat-card {
  text-align: center;
  padding: 28px 24px;
  border-radius: var(--radius-xl);
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  margin-bottom: 16px;
  transition: all 0.4s var(--ease-smooth);
  position: relative;
  overflow: hidden;
}

.stat-glow {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100px;
  height: 100px;
  border-radius: 50%;
  filter: blur(40px);
  opacity: 0;
  transition: opacity 0.4s var(--ease-smooth);
  pointer-events: none;
}

.stat-ocean .stat-glow {
  background: rgba(56, 189, 248, 0.3);
}

.stat-aqua .stat-glow {
  background: rgba(34, 211, 238, 0.3);
}

.stat-green .stat-glow {
  background: rgba(52, 211, 153, 0.3);
}

.stat-card:hover {
  border-color: rgba(56, 189, 248, 0.2);
  transform: translateY(-4px);
}

.stat-card:hover .stat-glow {
  opacity: 1;
}

.stat-card:hover .stat-icon {
  transform: scale(1.1);
}

.stat-icon {
  width: 60px;
  height: 60px;
  margin: 0 auto 16px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s var(--ease-smooth);
  position: relative;
}

.stat-ocean .stat-icon {
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.2), rgba(14, 165, 233, 0.1));
  border: 1px solid rgba(56, 189, 248, 0.2);
}

.stat-aqua .stat-icon {
  background: linear-gradient(135deg, rgba(34, 211, 238, 0.2), rgba(6, 182, 212, 0.1));
  border: 1px solid rgba(34, 211, 238, 0.2);
}

.stat-green .stat-icon {
  background: linear-gradient(135deg, rgba(52, 211, 153, 0.2), rgba(16, 185, 129, 0.1));
  border: 1px solid rgba(52, 211, 153, 0.2);
}

.stat-value {
  font-size: 36px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 4px;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 2px;
}

.stat-unit {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

/* ===== Heatmap & History Section - 海洋主题 ===== */
.heatmap-section,
.history-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth);
}

.heatmap-section {
  transition-delay: 0.2s;
}

.history-section {
  transition-delay: 0.3s;
}

.heatmap-section.mounted,
.history-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.section-card {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.1) !important;
  border-radius: var(--radius-xl) !important;
  position: relative;
  overflow: hidden;
}

.section-card :deep(.el-card__body) {
  padding: 28px;
}

.card-decoration {
  position: absolute;
  top: 0;
  right: 0;
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.08) 0%, transparent 70%);
  filter: blur(40px);
  pointer-events: none;
}

.section-header {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  margin-bottom: 24px;
}

.header-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.2), rgba(14, 165, 233, 0.1));
  border: 1px solid rgba(56, 189, 248, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.header-text h2 {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 4px;
}

.header-text p {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.55);
  margin: 0;
}

.loading-content {
  padding: 32px 0;
}

/* ===== History List ===== */
.history-skeleton {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border-radius: var(--radius-lg);
  background: rgba(56, 189, 248, 0.05);
  margin-bottom: 12px;
}

.skeleton-text {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.empty-state {
  text-align: center;
  padding: 48px 0;
}

.empty-icon {
  margin: 0 auto 20px;
  animation: float 4s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.empty-title {
  color: rgba(255, 255, 255, 0.6);
  font-size: 16px;
  margin: 0 0 8px;
}

.empty-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.55);
  margin: 0 0 24px;
}

.empty-btn {
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
  padding: 12px 24px !important;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 18px;
  border-radius: var(--radius-lg);
  background: rgba(56, 189, 248, 0.04);
  border: 1px solid rgba(56, 189, 248, 0.08);
  transition: all 0.3s var(--ease-smooth);
}

.history-item:hover {
  background: rgba(56, 189, 248, 0.08);
  border-color: rgba(56, 189, 248, 0.15);
  transform: translateX(4px);
}

.history-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgba(52, 211, 153, 0.2), rgba(16, 185, 129, 0.1));
  border: 1px solid rgba(52, 211, 153, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.history-info {
  flex: 1;
  min-width: 0;
}

.history-time {
  font-weight: 600;
  color: #fff;
  margin-bottom: 4px;
}

.history-detail {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-dot {
  color: rgba(56, 189, 248, 0.4);
}

.history-date {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  flex-shrink: 0;
}

/* List transition */
.list-enter-active,
.list-leave-active {
  transition: all 0.4s var(--ease-smooth);
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* ===== Settings Dialog - 新设计 ===== */
.settings-dialog-new :deep(.el-dialog) {
  background: linear-gradient(180deg, rgba(12, 20, 38, 0.98), rgba(8, 15, 30, 0.98)) !important;
  backdrop-filter: blur(30px);
  border: 1px solid rgba(56, 189, 248, 0.12);
  border-radius: 24px;
  overflow: hidden;
  box-shadow:
    0 30px 60px rgba(0, 0, 0, 0.6),
    0 0 120px rgba(56, 189, 248, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

.settings-dialog-new :deep(.el-dialog__header) {
  padding: 0;
  margin: 0;
  border: none;
}

.settings-dialog-new :deep(.el-dialog__headerbtn) {
  top: 24px;
  right: 24px;
  width: 36px;
  height: 36px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  transition: all 0.3s ease;
}

.settings-dialog-new :deep(.el-dialog__headerbtn:hover) {
  background: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.3);
}

.settings-dialog-new :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: rgba(255, 255, 255, 0.5);
  font-size: 16px;
}

.settings-dialog-new :deep(.el-dialog__headerbtn:hover .el-dialog__close) {
  color: #ef4444;
}

.settings-dialog-new :deep(.el-dialog__body) {
  padding: 0;
}

.settings-header {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 32px 32px 24px;
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.06), rgba(34, 211, 238, 0.03));
  border-bottom: 1px solid rgba(56, 189, 248, 0.08);
  position: relative;
  overflow: hidden;
}

.settings-header::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.1), transparent 70%);
  filter: blur(40px);
  pointer-events: none;
}

.settings-header-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.2), rgba(34, 211, 238, 0.1));
  border: 1px solid rgba(56, 189, 248, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgb(var(--ocean-surface));
  flex-shrink: 0;
  box-shadow: 0 8px 20px rgba(56, 189, 248, 0.15);
}

.settings-header-text h3 {
  margin: 0 0 6px;
  font-size: 22px;
  font-weight: 600;
  color: #fff;
  letter-spacing: -0.02em;
}

.settings-header-text p {
  margin: 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.45);
}

.settings-content {
  padding: 0;
}

.settings-tabs-new :deep(.el-tabs__header) {
  margin: 0;
  padding: 0 24px;
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid rgba(56, 189, 248, 0.08);
}

.settings-tabs-new :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.settings-tabs-new :deep(.el-tabs__item) {
  padding: 18px 24px;
  height: auto;
  color: rgba(255, 255, 255, 0.45);
  transition: all 0.3s ease;
  font-size: 14px;
}

.settings-tabs-new :deep(.el-tabs__item:hover) {
  color: rgba(255, 255, 255, 0.75);
}

.settings-tabs-new :deep(.el-tabs__item.is-active) {
  color: rgb(var(--ocean-surface));
  font-weight: 500;
}

.settings-tabs-new :deep(.el-tabs__active-bar) {
  background: linear-gradient(90deg, rgb(var(--ocean-surface)), rgb(var(--aqua-glow)));
  height: 3px;
  border-radius: 3px 3px 0 0;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
}

.setting-form-new {
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-header {
  margin-bottom: 8px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(56, 189, 248, 0.06);
}

.form-header h4 {
  margin: 0 0 8px;
  font-size: 17px;
  font-weight: 600;
  color: #fff;
  letter-spacing: -0.01em;
}

.form-header p {
  margin: 0;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.45);
  line-height: 1.5;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.form-field label {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.65);
}

.setting-form-new :deep(.el-input__wrapper) {
  background: rgba(0, 0, 0, 0.25);
  border: 1px solid rgba(56, 189, 248, 0.12);
  border-radius: 12px;
  box-shadow: none;
  padding: 6px 14px;
  transition: all 0.3s ease;
}

.setting-form-new :deep(.el-input__wrapper:hover) {
  border-color: rgba(56, 189, 248, 0.25);
  background: rgba(0, 0, 0, 0.3);
}

.setting-form-new :deep(.el-input__wrapper.is-focus) {
  border-color: rgb(var(--ocean-surface));
  background: rgba(56, 189, 248, 0.08);
  box-shadow: 0 0 0 4px rgba(56, 189, 248, 0.08);
}

.setting-form-new :deep(.el-input__inner) {
  color: #fff;
}

.setting-form-new :deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.35);
}

.setting-form-new :deep(.el-input__prefix) {
  color: rgba(56, 189, 248, 0.5);
}

.setting-form-new :deep(.el-input__count-inner) {
  background: transparent;
  color: rgba(255, 255, 255, 0.35);
}

.submit-btn-new {
  margin-top: 12px;
  width: 100%;
  height: 48px;
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
  border-radius: 12px !important;
  font-weight: 600;
  font-size: 15px;
  letter-spacing: 0.02em;
  transition: all 0.3s ease !important;
  position: relative;
  overflow: hidden;
}

.submit-btn-new::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
  transition: left 0.5s ease;
}

.submit-btn-new:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 28px rgba(14, 165, 233, 0.35);
}

.submit-btn-new:hover::before {
  left: 100%;
}

.submit-btn-new:active {
  transform: translateY(0);
}
</style>
