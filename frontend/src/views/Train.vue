<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Timer from '@/components/Timer.vue'
import { useTrainingStore } from '@/store/training'
import { useUserStore } from '@/store/user'
import { checkin } from '@/api/checkin'
import { Setting, Clock, Aim, TrendCharts, WarningFilled, StarFilled, ArrowUp, ArrowDown } from '@element-plus/icons-vue'
import confetti from 'canvas-confetti'

const trainingStore = useTrainingStore()
const userStore = useUserStore()

const mounted = ref(false)
const showSettings = ref(false)
const showResultModal = ref(false)

onMounted(() => {
  setTimeout(() => {
    mounted.value = true
  }, 100)
})

// 计算预计训练时长
const estimatedDuration = computed(() => {
  const { contractTime, holdTime, relaxTime, cycles } = trainingStore.settings
  const totalSeconds = (contractTime + holdTime + relaxTime) * cycles
  const mins = Math.floor(totalSeconds / 60)
  const secs = totalSeconds % 60
  return mins > 0 ? `${mins}分${secs}秒` : `${secs}秒`
})
const checkinResult = ref<{
  streak: number
  maxStreak: number
  total: number
} | null>(null)
const checkinError = ref('')

watch(() => trainingStore.isRunning, async (isRunning, wasRunning) => {
  if (wasRunning && !isRunning && trainingStore.currentCycle > 0) {
    await handleCheckin()
  }
})

async function handleCheckin() {
  if (!userStore.isLoggedIn) {
    checkinError.value = '请先登录后再打卡'
    showResultModal.value = true
    return
  }

  try {
    const res = await checkin({
      duration: trainingStore.totalDuration,
      cycles: trainingStore.currentCycle
    })

    checkinResult.value = {
      streak: res.current_streak,
      maxStreak: res.max_streak,
      total: res.total_checkin
    }

    userStore.updateStreak(res.current_streak, res.max_streak, res.total_checkin)
    showResultModal.value = true

    // Enhanced confetti
    const colors = ['#a855f7', '#ec4899', '#f59e0b', '#10b981']
    confetti({
      particleCount: 150,
      spread: 100,
      origin: { y: 0.6 },
      colors
    })
  } catch (e) {
    checkinError.value = e instanceof Error ? e.message : '打卡失败'
    showResultModal.value = true
  }
}

function closeResult() {
  showResultModal.value = false
  checkinResult.value = null
  checkinError.value = ''
  trainingStore.reset()
}

function updateSetting(key: string, value: number) {
  trainingStore.updateSettings({ [key]: value })
}
</script>

<template>
  <MainLayout>
    <div class="train-page">
      <!-- Header with decorative elements -->
      <div class="train-header" :class="{ mounted }">
        <!-- Decorative background -->
        <div class="header-bg">
          <div class="header-glow"></div>
        </div>

        <div class="train-eyebrow">
          <span class="eyebrow-line"></span>
          <span class="eyebrow-text">专业训练</span>
          <span class="eyebrow-line"></span>
        </div>

        <h1 class="train-title">
          <span class="gradient-text">潮汐训练器</span>
        </h1>
        <p class="train-subtitle">跟随节奏，完成今日训练，坚持就是胜利</p>

        <!-- Quick stats -->
        <div class="quick-stats">
          <div class="quick-stat">
            <el-icon><Clock /></el-icon>
            <span>预计 {{ estimatedDuration }}</span>
          </div>
          <div class="quick-stat">
            <el-icon><Aim /></el-icon>
            <span>{{ trainingStore.settings.cycles }} 个循环</span>
          </div>
        </div>
      </div>

      <!-- Timer Component with enhanced container -->
      <div class="timer-wrapper" :class="{ mounted }">
        <!-- Ambient glow behind timer -->
        <div class="timer-glow" :class="{ active: trainingStore.isRunning }"></div>
        <Timer />
      </div>

      <!-- Settings Toggle -->
      <div class="settings-toggle" :class="{ mounted }">
        <button
          @click="showSettings = !showSettings"
          class="settings-toggle-btn"
          :class="{ 'is-open': showSettings }"
        >
          <span class="toggle-icon-wrapper">
            <el-icon :size="18"><Setting /></el-icon>
          </span>
          <span class="toggle-text">训练设置</span>
          <el-icon class="toggle-arrow" :size="14">
            <ArrowUp v-if="showSettings" />
            <ArrowDown v-else />
          </el-icon>
        </button>
      </div>

      <!-- Settings Panel -->
      <el-collapse-transition>
        <el-card v-if="showSettings" class="settings-card" shadow="never">
          <div class="settings-decoration"></div>

          <div class="settings-content">
            <h3 class="settings-title">
              <div class="title-icon">
                <el-icon><Setting /></el-icon>
              </div>
              自定义训练参数
            </h3>

            <div class="settings-grid">
              <div class="setting-item">
                <label class="setting-label">收缩时间</label>
                <div class="setting-input-wrapper">
                  <el-input-number
                    :model-value="trainingStore.settings.contractTime"
                    @update:model-value="updateSetting('contractTime', $event as number)"
                    :min="1"
                    :max="30"
                    controls-position="right"
                    size="large"
                    class="setting-input"
                  />
                  <span class="setting-unit">秒</span>
                </div>
              </div>
              <div class="setting-item">
                <label class="setting-label">保持时间</label>
                <div class="setting-input-wrapper">
                  <el-input-number
                    :model-value="trainingStore.settings.holdTime"
                    @update:model-value="updateSetting('holdTime', $event as number)"
                    :min="1"
                    :max="30"
                    controls-position="right"
                    size="large"
                    class="setting-input"
                  />
                  <span class="setting-unit">秒</span>
                </div>
              </div>
              <div class="setting-item">
                <label class="setting-label">放松时间</label>
                <div class="setting-input-wrapper">
                  <el-input-number
                    :model-value="trainingStore.settings.relaxTime"
                    @update:model-value="updateSetting('relaxTime', $event as number)"
                    :min="1"
                    :max="30"
                    controls-position="right"
                    size="large"
                    class="setting-input"
                  />
                  <span class="setting-unit">秒</span>
                </div>
              </div>
              <div class="setting-item">
                <label class="setting-label">循环次数</label>
                <div class="setting-input-wrapper">
                  <el-input-number
                    :model-value="trainingStore.settings.cycles"
                    @update:model-value="updateSetting('cycles', $event as number)"
                    :min="1"
                    :max="50"
                    controls-position="right"
                    size="large"
                    class="setting-input"
                  />
                  <span class="setting-unit">次</span>
                </div>
              </div>
            </div>

            <!-- Training summary -->
            <el-alert
              type="info"
              :closable="false"
              class="training-summary"
            >
              <template #title>
                <div class="summary-content">
                  <p class="summary-text">
                    <span class="summary-label">训练概览：</span>
                    每个循环包含收缩、保持、放松三个阶段
                  </p>
                  <div class="summary-stats">
                    <span class="summary-stat">
                      <el-icon><Clock /></el-icon>
                      总时长: {{ estimatedDuration }}
                    </span>
                    <span class="summary-stat">
                      <el-icon><Aim /></el-icon>
                      循环数: {{ trainingStore.settings.cycles }}
                    </span>
                  </div>
                </div>
              </template>
            </el-alert>
          </div>
        </el-card>
      </el-collapse-transition>

      <!-- Tips Section -->
      <div class="tips-section" :class="{ mounted }">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="8">
            <div class="tip-card tip-contract">
              <div class="tip-icon">
                <el-icon :size="20" color="#fb7185"><Lightning /></el-icon>
              </div>
              <div class="tip-content">
                <h4 class="tip-title">收缩阶段</h4>
                <p class="tip-desc">用力收紧盆底肌肉</p>
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="tip-card tip-hold">
              <div class="tip-icon">
                <el-icon :size="20" color="#fbbf24"><Aim /></el-icon>
              </div>
              <div class="tip-content">
                <h4 class="tip-title">保持阶段</h4>
                <p class="tip-desc">维持收缩状态</p>
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="tip-card tip-relax">
              <div class="tip-icon">
                <el-icon :size="20" color="#34d399"><TrendCharts /></el-icon>
              </div>
              <div class="tip-content">
                <h4 class="tip-title">放松阶段</h4>
                <p class="tip-desc">完全放松肌肉</p>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- Result Modal -->
      <el-dialog
        v-model="showResultModal"
        width="400px"
        :show-close="false"
        center
        class="result-dialog"
      >
        <div class="result-content">
          <template v-if="checkinResult">
            <div class="result-icon success">
              <el-icon :size="40" color="#fff"><StarFilled /></el-icon>
            </div>
            <h3 class="result-title">打卡成功!</h3>
            <p class="result-subtitle">太棒了，继续保持！</p>

            <el-row :gutter="16" class="result-stats">
              <el-col :span="8">
                <div class="result-stat">
                  <div class="stat-value orange">{{ checkinResult.streak }}</div>
                  <div class="stat-label">连续天数</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="result-stat">
                  <div class="stat-value purple">{{ checkinResult.maxStreak }}</div>
                  <div class="stat-label">最高记录</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="result-stat">
                  <div class="stat-value green">{{ checkinResult.total }}</div>
                  <div class="stat-label">累计打卡</div>
                </div>
              </el-col>
            </el-row>
          </template>

          <template v-else>
            <div class="result-icon error">
              <el-icon :size="40" color="#f87171"><WarningFilled /></el-icon>
            </div>
            <h3 class="result-title">{{ checkinError }}</h3>
            <p class="result-subtitle">
              <RouterLink to="/login" class="login-link">点击登录</RouterLink>
              后即可记录打卡
            </p>
          </template>

          <el-button
            @click="closeResult"
            type="primary"
            size="large"
            round
            class="result-btn"
          >
            确定
          </el-button>
        </div>
      </el-dialog>
    </div>
  </MainLayout>
</template>

<style scoped>
.train-page {
  max-width: 768px;
  margin: 0 auto;
}

/* Header */
.train-header {
  text-align: center;
  margin-bottom: 48px;
  position: relative;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease;
}

.train-header.mounted {
  opacity: 1;
  transform: translateY(0);
}

.header-bg {
  position: absolute;
  inset: 0;
  z-index: -1;
  overflow: hidden;
}

.header-glow {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 384px;
  height: 128px;
  background: linear-gradient(to bottom, rgba(139, 92, 246, 0.2), transparent);
  border-radius: 50%;
  filter: blur(60px);
}

.train-eyebrow {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 20px;
}

.train-eyebrow .eyebrow-line {
  width: 40px;
  height: 1px;
  background: linear-gradient(90deg, transparent, #8b5cf6);
}

.train-eyebrow .eyebrow-line:last-child {
  background: linear-gradient(90deg, #8b5cf6, transparent);
}

.train-eyebrow .eyebrow-text {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 3px;
  background: linear-gradient(135deg, #a78bfa, #c4b5fd);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.train-title {
  font-size: 36px;
  font-weight: 700;
  margin: 0 0 12px;
}

.gradient-text {
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.train-subtitle {
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 24px;
}

.quick-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
}

.quick-stat {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.4);
}

/* Timer Wrapper */
.timer-wrapper {
  position: relative;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease 0.1s;
}

.timer-wrapper.mounted {
  opacity: 1;
  transform: translateY(0);
}

.timer-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 320px;
  height: 320px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(236, 72, 153, 0.2));
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0;
  transition: opacity 0.5s ease;
  z-index: -1;
}

.timer-glow.active {
  opacity: 1;
  animation: pulse-soft 2s ease-in-out infinite;
}

@keyframes pulse-soft {
  0%, 100% { transform: translate(-50%, -50%) scale(1); }
  50% { transform: translate(-50%, -50%) scale(1.1); }
}

/* Settings Toggle */
.settings-toggle {
  display: flex;
  justify-content: center;
  margin-top: 48px;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease 0.2s;
}

.settings-toggle.mounted {
  opacity: 1;
  transform: translateY(0);
}

.settings-toggle-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.03);
  cursor: pointer;
  transition: all 0.25s ease;
}

.settings-toggle-btn:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.12);
}

.settings-toggle-btn.is-open {
  background: rgba(139, 92, 246, 0.1);
  border-color: rgba(139, 92, 246, 0.25);
}

.toggle-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(168, 85, 247, 0.1));
  color: #a78bfa;
  transition: transform 0.3s ease;
}

.settings-toggle-btn:hover .toggle-icon-wrapper {
  transform: rotate(45deg);
}

.toggle-text {
  font-size: 15px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
}

.settings-toggle-btn:hover .toggle-text {
  color: #fff;
}

.toggle-arrow {
  color: rgba(255, 255, 255, 0.4);
  transition: transform 0.2s ease;
}

.settings-toggle-btn.is-open .toggle-arrow {
  color: #a78bfa;
}

/* Settings Card */
.settings-card {
  margin-top: 24px;
  background: rgba(30, 30, 46, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 10px !important;
  overflow: hidden;
  position: relative;
}

.settings-card :deep(.el-card__body) {
  padding: 24px;
}

.settings-decoration {
  position: absolute;
  top: 0;
  right: 0;
  width: 128px;
  height: 128px;
  background: linear-gradient(225deg, rgba(139, 92, 246, 0.1), transparent);
  border-radius: 50%;
  filter: blur(40px);
  pointer-events: none;
}

.settings-content {
  position: relative;
}

.settings-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 24px;
}

.title-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(168, 85, 247, 0.1));
  display: flex;
  align-items: center;
  justify-content: center;
  color: #a78bfa;
}

/* Settings Grid Layout */
.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

@media (min-width: 640px) {
  .settings-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 24px;
  }
}

.setting-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.setting-label {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
}

.setting-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.setting-unit {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  flex-shrink: 0;
}

.setting-input {
  flex: 1;
  min-width: 0;
}

.setting-input :deep(.el-input-number__decrease),
.setting-input :deep(.el-input-number__increase) {
  background: rgba(255, 255, 255, 0.08) !important;
  border-color: rgba(255, 255, 255, 0.12) !important;
  color: rgba(255, 255, 255, 0.7) !important;
}

.setting-input :deep(.el-input-number__decrease:hover),
.setting-input :deep(.el-input-number__increase:hover) {
  background: rgba(255, 255, 255, 0.12) !important;
  color: #fff !important;
}

.setting-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 8px !important;
  box-shadow: none !important;
  padding: 4px 12px !important;
}

.setting-input :deep(.el-input__inner) {
  color: #fff !important;
  text-align: center;
  font-weight: 600;
  font-size: 18px;
}

.training-summary {
  margin-top: 20px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.1), rgba(168, 85, 247, 0.1)) !important;
  border: 1px solid rgba(139, 92, 246, 0.2) !important;
  border-radius: 8px !important;
}

.training-summary :deep(.el-alert__title) {
  color: rgba(255, 255, 255, 0.8);
}

.summary-content {
  width: 100%;
}

.summary-text {
  margin: 0 0 8px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
}

.summary-label {
  color: #a78bfa;
  font-weight: 500;
}

.summary-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.summary-stat {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

/* Tips Section */
.tips-section {
  margin-top: 48px;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease 0.3s;
}

.tips-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.tip-card {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.2s ease;
  margin-bottom: 16px;
}

.tip-card:hover {
  border-color: rgba(255, 255, 255, 0.1);
}

.tip-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.tip-contract .tip-icon {
  background: linear-gradient(135deg, rgba(244, 63, 94, 0.2), rgba(236, 72, 153, 0.1));
}

.tip-hold .tip-icon {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(249, 115, 22, 0.1));
}

.tip-relax .tip-icon {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(20, 184, 166, 0.1));
}

.tip-content {
  flex: 1;
  min-width: 0;
}

.tip-title {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  margin: 0 0 2px;
}

.tip-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  margin: 0;
}

/* Result Dialog */
.result-dialog :deep(.el-dialog) {
  background: rgba(30, 30, 46, 0.95) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px !important;
}

.result-dialog :deep(.el-dialog__body) {
  padding: 40px;
}

.result-content {
  text-align: center;
}

.result-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.result-icon.success {
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
}

.result-icon.error {
  background: rgba(239, 68, 68, 0.2);
}

.result-title {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
}

.result-subtitle {
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 24px;
}

.result-stats {
  margin-bottom: 32px;
}

.result-stat {
  padding: 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-value.orange {
  color: #fb923c;
}

.stat-value.purple {
  color: #a78bfa;
}

.stat-value.green {
  color: #34d399;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.login-link {
  color: #a78bfa;
  text-decoration: none;
}

.login-link:hover {
  color: #c4b5fd;
}

.result-btn {
  padding: 12px 32px !important;
  background: linear-gradient(135deg, #8b5cf6, #ec4899) !important;
  border: none !important;
}
</style>
