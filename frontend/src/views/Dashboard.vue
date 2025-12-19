<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Heatmap from '@/components/Heatmap.vue'
import { useUserStore } from '@/store/user'
import { getHeatmap, getHistory, type CheckinRecord } from '@/api/checkin'
import { updateProfile, updateUsername, updatePassword } from '@/api/auth'
import { Timer, Calendar, Clock, CircleCheck, ArrowRight, Pointer, Trophy, Aim, Edit, Lock, User, Setting } from '@element-plus/icons-vue'
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

// 计算用户等级 - 海洋主题
const userLevel = computed(() => {
  const total = userStore.user?.total_checkin || 0
  if (total >= 365) return { name: '深海传奇', color: '#fbbf24', bg: 'linear-gradient(135deg, rgba(251, 191, 36, 0.25), rgba(245, 158, 11, 0.1))', icon: '🌊' }
  if (total >= 180) return { name: '海洋大师', color: '#38bdf8', bg: 'linear-gradient(135deg, rgba(56, 189, 248, 0.25), rgba(14, 165, 233, 0.1))', icon: '🐋' }
  if (total >= 90) return { name: '浪潮专家', color: '#22d3ee', bg: 'linear-gradient(135deg, rgba(34, 211, 238, 0.25), rgba(6, 182, 212, 0.1))', icon: '🐬' }
  if (total >= 30) return { name: '潮汐进阶', color: '#34d399', bg: 'linear-gradient(135deg, rgba(52, 211, 153, 0.25), rgba(16, 185, 129, 0.1))', icon: '🐠' }
  if (total >= 7) return { name: '入海新手', color: '#0ea5e9', bg: 'linear-gradient(135deg, rgba(14, 165, 233, 0.25), rgba(2, 132, 199, 0.1))', icon: '🐟' }
  return { name: '初探海域', color: 'rgba(255, 255, 255, 0.6)', bg: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.05))', icon: '🐚' }
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
              <!-- Level badge -->
              <div
                class="level-badge"
                :style="{ background: userLevel.bg, color: userLevel.color }"
              >
                <span class="level-icon">{{ userLevel.icon }}</span>
                {{ userLevel.name }}
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
                <el-button class="settings-btn" circle size="small" @click="openSettings">
                  <el-icon :size="14"><Setting /></el-icon>
                </el-button>
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
              </div>

              <div class="action-buttons">
                <RouterLink to="/train">
                  <el-button type="primary" size="large" round class="train-btn">
                    <el-icon><Timer /></el-icon>
                    开始今日训练
                    <el-icon class="arrow-icon"><ArrowRight /></el-icon>
                  </el-button>
                </RouterLink>
                <RouterLink v-if="userStore.user?.is_admin" to="/admin">
                  <el-button size="large" round class="admin-btn">
                    <el-icon><Setting /></el-icon>
                    管理后台
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

    <!-- 设置对话框 -->
    <el-dialog
      v-model="showSettingsDialog"
      title="账号设置"
      width="500px"
      class="settings-dialog"
      :close-on-click-modal="false"
    >
      <el-tabs v-model="activeSettingTab" class="settings-tabs">
        <el-tab-pane label="显示名称" name="profile">
          <div class="setting-form">
            <p class="setting-desc">显示名称将在排行榜和个人主页中展示</p>
            <el-input
              v-model="profileForm.display_name"
              placeholder="请输入显示名称"
              maxlength="50"
              show-word-limit
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
            <el-button
              type="primary"
              :loading="settingsLoading"
              @click="handleUpdateProfile"
              class="setting-submit-btn"
            >
              保存修改
            </el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane label="用户名" name="username">
          <div class="setting-form">
            <p class="setting-desc">用户名用于登录，只能包含字母、数字和下划线</p>
            <el-input
              v-model="usernameForm.username"
              placeholder="请输入新用户名"
              maxlength="50"
            >
              <template #prefix>
                <el-icon><Edit /></el-icon>
              </template>
            </el-input>
            <el-button
              type="primary"
              :loading="settingsLoading"
              @click="handleUpdateUsername"
              class="setting-submit-btn"
            >
              保存修改
            </el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="password">
          <div class="setting-form">
            <p class="setting-desc">为了账号安全，请定期更换密码</p>
            <el-input
              v-model="passwordForm.old_password"
              type="password"
              placeholder="请输入原密码"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
            <el-input
              v-model="passwordForm.new_password"
              type="password"
              placeholder="请输入新密码（至少6位）"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
            <el-input
              v-model="passwordForm.confirm_password"
              type="password"
              placeholder="请再次输入新密码"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
            <el-button
              type="primary"
              :loading="settingsLoading"
              @click="handleUpdatePassword"
              class="setting-submit-btn"
            >
              修改密码
            </el-button>
          </div>
        </el-tab-pane>
      </el-tabs>
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

.level-badge {
  position: absolute;
  bottom: -10px;
  right: -10px;
  padding: 6px 14px;
  border-radius: var(--radius-full);
  font-size: 12px;
  font-weight: 600;
  border: 1px solid rgba(56, 189, 248, 0.2);
  display: flex;
  align-items: center;
  gap: 4px;
  backdrop-filter: blur(10px);
}

.level-icon {
  font-size: 14px;
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

.train-btn {
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
  padding: 14px 32px !important;
  box-shadow: 0 10px 30px rgba(14, 165, 233, 0.3);
  transition: all 0.3s var(--ease-smooth) !important;
}

.train-btn:hover {
  box-shadow: 0 15px 40px rgba(14, 165, 233, 0.4);
  transform: translateY(-2px);
}

.train-btn .arrow-icon {
  margin-left: 8px;
  transition: transform 0.2s var(--ease-smooth);
}

.train-btn:hover .arrow-icon {
  transform: translateX(4px);
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

/* ===== Settings Button ===== */
.settings-btn {
  background: rgba(56, 189, 248, 0.1) !important;
  border: 1px solid rgba(56, 189, 248, 0.2) !important;
  color: rgba(255, 255, 255, 0.7) !important;
  transition: all 0.3s var(--ease-smooth) !important;
}

.settings-btn:hover {
  background: rgba(56, 189, 248, 0.2) !important;
  border-color: rgba(56, 189, 248, 0.3) !important;
  color: #fff !important;
}

/* ===== Action Buttons ===== */
.action-buttons {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
}

@media (min-width: 640px) {
  .action-buttons {
    justify-content: flex-start;
  }
}

.admin-btn {
  background: rgba(251, 191, 36, 0.1) !important;
  border: 1px solid rgba(251, 191, 36, 0.3) !important;
  color: #fbbf24 !important;
}

.admin-btn:hover {
  background: rgba(251, 191, 36, 0.2) !important;
  border-color: rgba(251, 191, 36, 0.4) !important;
}

/* ===== Settings Dialog ===== */
.settings-dialog :deep(.el-dialog) {
  background: rgba(15, 23, 42, 0.95) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-xl);
}

.settings-dialog :deep(.el-dialog__header) {
  border-bottom: 1px solid rgba(56, 189, 248, 0.1);
  padding: 20px 24px;
}

.settings-dialog :deep(.el-dialog__title) {
  color: #fff;
  font-weight: 600;
}

.settings-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.settings-tabs :deep(.el-tabs__header) {
  margin-bottom: 24px;
}

.settings-tabs :deep(.el-tabs__nav-wrap::after) {
  background: rgba(56, 189, 248, 0.1);
}

.settings-tabs :deep(.el-tabs__item) {
  color: rgba(255, 255, 255, 0.6);
}

.settings-tabs :deep(.el-tabs__item.is-active) {
  color: rgb(var(--ocean-surface));
}

.settings-tabs :deep(.el-tabs__active-bar) {
  background: rgb(var(--ocean-surface));
}

.setting-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.setting-desc {
  color: rgba(255, 255, 255, 0.55);
  font-size: 14px;
  margin: 0;
}

.setting-form :deep(.el-input__wrapper) {
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.15);
  box-shadow: none;
}

.setting-form :deep(.el-input__wrapper:hover) {
  border-color: rgba(56, 189, 248, 0.3);
}

.setting-form :deep(.el-input__wrapper.is-focus) {
  border-color: rgb(var(--ocean-surface));
}

.setting-form :deep(.el-input__inner) {
  color: #fff;
}

.setting-form :deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.4);
}

.setting-submit-btn {
  margin-top: 8px;
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
}
</style>
