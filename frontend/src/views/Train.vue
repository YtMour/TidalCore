<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Timer from '@/components/Timer.vue'
import { useTrainingStore } from '@/store/training'
import { useUserStore } from '@/store/user'
import { checkin } from '@/api/checkin'
import { Setting, Clock, Aim, TrendCharts, WarningFilled, Lightning, Refresh } from '@element-plus/icons-vue'
import { DIFFICULTY_PRESETS } from '@/store/training'
import confetti from 'canvas-confetti'

const trainingStore = useTrainingStore()
const userStore = useUserStore()

const mounted = ref(false)
const showSettings = ref(false)
const showDifficulty = ref(false)
const showResultModal = ref(false)

// 难度选择
function selectDifficulty(level: string) {
  trainingStore.setDifficulty(level as any)
}

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

    // Enhanced confetti - ocean colors
    const colors = ['#38bdf8', '#22d3ee', '#0ea5e9', '#34d399', '#fbbf24']
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
          <span class="eyebrow-text">潮汐训练</span>
          <span class="eyebrow-line"></span>
        </div>

        <h1 class="train-title">
          <span class="gradient-text">潮汐呼吸</span>
        </h1>
        <p class="train-subtitle">跟随海浪节奏，完成今日训练，如潮汐般坚持</p>

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
        <Timer v-model:showSettings="showSettings" v-model:showDifficulty="showDifficulty" />
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
                <div class="stepper-control">
                  <button
                    class="stepper-btn"
                    @click="updateSetting('contractTime', Math.max(1, trainingStore.settings.contractTime - 1))"
                  >
                    <span class="stepper-icon">−</span>
                  </button>
                  <div class="stepper-value">
                    <span class="value-number">{{ trainingStore.settings.contractTime }}</span>
                    <span class="value-unit">秒</span>
                  </div>
                  <button
                    class="stepper-btn"
                    @click="updateSetting('contractTime', Math.min(30, trainingStore.settings.contractTime + 1))"
                  >
                    <span class="stepper-icon">+</span>
                  </button>
                </div>
              </div>
              <div class="setting-item">
                <label class="setting-label">保持时间</label>
                <div class="stepper-control">
                  <button
                    class="stepper-btn"
                    @click="updateSetting('holdTime', Math.max(1, trainingStore.settings.holdTime - 1))"
                  >
                    <span class="stepper-icon">−</span>
                  </button>
                  <div class="stepper-value">
                    <span class="value-number">{{ trainingStore.settings.holdTime }}</span>
                    <span class="value-unit">秒</span>
                  </div>
                  <button
                    class="stepper-btn"
                    @click="updateSetting('holdTime', Math.min(30, trainingStore.settings.holdTime + 1))"
                  >
                    <span class="stepper-icon">+</span>
                  </button>
                </div>
              </div>
              <div class="setting-item">
                <label class="setting-label">放松时间</label>
                <div class="stepper-control">
                  <button
                    class="stepper-btn"
                    @click="updateSetting('relaxTime', Math.max(1, trainingStore.settings.relaxTime - 1))"
                  >
                    <span class="stepper-icon">−</span>
                  </button>
                  <div class="stepper-value">
                    <span class="value-number">{{ trainingStore.settings.relaxTime }}</span>
                    <span class="value-unit">秒</span>
                  </div>
                  <button
                    class="stepper-btn"
                    @click="updateSetting('relaxTime', Math.min(30, trainingStore.settings.relaxTime + 1))"
                  >
                    <span class="stepper-icon">+</span>
                  </button>
                </div>
              </div>
              <div class="setting-item">
                <label class="setting-label">循环次数</label>
                <div class="stepper-control">
                  <button
                    class="stepper-btn"
                    @click="updateSetting('cycles', Math.max(1, trainingStore.settings.cycles - 1))"
                  >
                    <span class="stepper-icon">−</span>
                  </button>
                  <div class="stepper-value">
                    <span class="value-number">{{ trainingStore.settings.cycles }}</span>
                    <span class="value-unit">次</span>
                  </div>
                  <button
                    class="stepper-btn"
                    @click="updateSetting('cycles', Math.min(50, trainingStore.settings.cycles + 1))"
                  >
                    <span class="stepper-icon">+</span>
                  </button>
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

      <!-- Difficulty Panel -->
      <el-collapse-transition>
        <el-card v-if="showDifficulty" class="settings-card difficulty-card" shadow="never">
          <div class="settings-decoration"></div>

          <div class="settings-content">
            <h3 class="settings-title">
              <div class="title-icon difficulty-title-icon">
                <span>🎯</span>
              </div>
              选择训练难度
            </h3>

            <div class="difficulty-grid">
              <!-- 预设难度 -->
              <button
                v-for="preset in DIFFICULTY_PRESETS"
                :key="preset.id"
                class="difficulty-option"
                :class="{ active: trainingStore.difficulty === preset.id }"
                @click="selectDifficulty(preset.id)"
              >
                <span class="difficulty-option-icon">{{ preset.icon }}</span>
                <div class="difficulty-option-content">
                  <span class="difficulty-option-name">{{ preset.name }}</span>
                  <span class="difficulty-option-desc">{{ preset.description }}</span>
                </div>
                <div class="difficulty-option-params">
                  <span>{{ preset.settings.contractTime }}-{{ preset.settings.holdTime }}-{{ preset.settings.relaxTime }}s</span>
                  <span>{{ preset.settings.cycles }}次</span>
                </div>
              </button>

              <!-- 随机难度 -->
              <button
                class="difficulty-option random-option"
                :class="{ active: trainingStore.difficulty === 'random' }"
                @click="selectDifficulty('random')"
              >
                <span class="difficulty-option-icon">🎲</span>
                <div class="difficulty-option-content">
                  <span class="difficulty-option-name">随机挑战</span>
                  <span class="difficulty-option-desc">每次随机生成参数</span>
                </div>
                <div class="difficulty-option-params">
                  <el-icon :size="16"><Refresh /></el-icon>
                </div>
              </button>
            </div>

            <!-- 当前参数预览 -->
            <div class="difficulty-preview">
              <div class="preview-label">当前训练参数</div>
              <div class="preview-params">
                <div class="preview-param">
                  <span class="param-label">收缩</span>
                  <span class="param-value">{{ trainingStore.settings.contractTime }}s</span>
                </div>
                <div class="preview-divider">→</div>
                <div class="preview-param">
                  <span class="param-label">保持</span>
                  <span class="param-value">{{ trainingStore.settings.holdTime }}s</span>
                </div>
                <div class="preview-divider">→</div>
                <div class="preview-param">
                  <span class="param-label">放松</span>
                  <span class="param-value">{{ trainingStore.settings.relaxTime }}s</span>
                </div>
                <div class="preview-divider">×</div>
                <div class="preview-param">
                  <span class="param-label">循环</span>
                  <span class="param-value">{{ trainingStore.settings.cycles }}</span>
                </div>
              </div>
            </div>

            <p class="difficulty-tip">
              💡 选择难度后会自动应用对应参数，也可在「训练设置」中自定义调整
            </p>
          </div>
        </el-card>
      </el-collapse-transition>

      <!-- Tips Section - 潮汐阶段 -->
      <div class="tips-section" :class="{ mounted }">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="8">
            <div class="tip-card tip-contract">
              <div class="tip-icon">
                <el-icon :size="20" color="#f43f5e"><Lightning /></el-icon>
              </div>
              <div class="tip-content">
                <h4 class="tip-title">涨潮收缩</h4>
                <p class="tip-desc">如浪涌般用力收紧</p>
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="tip-card tip-hold">
              <div class="tip-icon">
                <el-icon :size="20" color="#f59e0b"><Aim /></el-icon>
              </div>
              <div class="tip-content">
                <h4 class="tip-title">平潮保持</h4>
                <p class="tip-desc">如海面般稳定维持</p>
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="tip-card tip-relax">
              <div class="tip-icon">
                <el-icon :size="20" color="#10b981"><TrendCharts /></el-icon>
              </div>
              <div class="tip-content">
                <h4 class="tip-title">退潮放松</h4>
                <p class="tip-desc">如浪退般完全释放</p>
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
              <svg viewBox="0 0 32 32" fill="none" class="result-wave-svg">
                <path d="M6 16C6 16 9 12 13 16C17 20 21 12 25 16" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
              </svg>
            </div>
            <h3 class="result-title">打卡成功!</h3>
            <p class="result-subtitle">如潮汐般坚持，继续加油！</p>

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
  padding: 40px 24px 80px;
  position: relative;
  min-height: calc(100vh - 72px);
  z-index: 1;
}

/* Header - 海洋主题 */
.train-header {
  text-align: center;
  margin-bottom: 48px;
  position: relative;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth);
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
  width: 400px;
  height: 150px;
  background: linear-gradient(to bottom, rgba(56, 189, 248, 0.2), transparent);
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
  background: linear-gradient(90deg, transparent, rgb(var(--ocean-surface)));
}

.train-eyebrow .eyebrow-line:last-child {
  background: linear-gradient(90deg, rgb(var(--ocean-surface)), transparent);
}

.train-eyebrow .eyebrow-text {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 3px;
  background: linear-gradient(135deg, rgb(var(--ocean-surface)), rgb(var(--aqua-glow)));
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
  background: linear-gradient(
    90deg,
    #38bdf8 0%,
    #5dd3f5 12%,
    #22d3ee 25%,
    #45c8e8 37%,
    #38bdf8 50%,
    #5dd3f5 62%,
    #22d3ee 75%,
    #45c8e8 87%,
    #38bdf8 100%
  );
  background-size: 300% 100%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: ocean-shimmer 10s ease-in-out infinite;
}

@keyframes ocean-shimmer {
  0%, 100% { background-position: 0% center; }
  50% { background-position: 200% center; }
}

.train-subtitle {
  color: rgba(255, 255, 255, 0.65);
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
  color: rgba(255, 255, 255, 0.6);
  padding: 8px 14px;
  border-radius: var(--radius-md);
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.1);
}

/* Timer Wrapper - 海洋光效 */
.timer-wrapper {
  position: relative;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth) 0.1s;
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
  width: 340px;
  height: 340px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.25) 0%, rgba(14, 165, 233, 0.1) 50%, transparent 70%);
  border-radius: 50%;
  filter: blur(40px);
  opacity: 0;
  transition: opacity 0.5s var(--ease-smooth);
  z-index: -1;
}

.timer-glow.active {
  opacity: 1;
  animation: tidal-pulse 4s var(--ease-tidal) infinite;
}

@keyframes tidal-pulse {
  0%, 100% { transform: translate(-50%, -50%) scale(1); opacity: 0.8; }
  50% { transform: translate(-50%, -50%) scale(1.15); opacity: 1; }
}

/* Settings Card - 海洋风格 */
.settings-card {
  margin-top: 24px;
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.1) !important;
  border-radius: var(--radius-xl) !important;
  overflow: hidden;
  position: relative;
}

.settings-card :deep(.el-card__body) {
  padding: 24px;
  overflow: hidden;
}

.settings-decoration {
  position: absolute;
  top: 0;
  right: 0;
  width: 150px;
  height: 150px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.1) 0%, transparent 70%);
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
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.2), rgba(14, 165, 233, 0.1));
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgb(var(--ocean-surface));
}

/* Settings Grid Layout */
.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px 20px;
  margin-bottom: 24px;
}

@media (min-width: 640px) {
  .settings-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 16px;
  }
}

.setting-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.setting-label {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
  white-space: nowrap;
}

/* Custom Stepper Control - 海洋风格 */
.stepper-control {
  display: flex;
  align-items: center;
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.15);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.stepper-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 44px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.6);
  transition: all 0.2s var(--ease-smooth);
  flex-shrink: 0;
}

.stepper-btn:hover {
  background: rgba(56, 189, 248, 0.2);
  color: rgb(var(--ocean-surface));
}

.stepper-btn:active {
  background: rgba(56, 189, 248, 0.3);
}

.stepper-icon {
  font-size: 18px;
  font-weight: 300;
  line-height: 1;
}

.stepper-value {
  flex: 1;
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 2px;
  padding: 0 4px;
  min-width: 0;
}

.value-number {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  line-height: 1;
}

.value-unit {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  flex-shrink: 0;
}

.training-summary {
  margin-top: 20px;
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.1), rgba(14, 165, 233, 0.05)) !important;
  border: 1px solid rgba(56, 189, 248, 0.2) !important;
  border-radius: var(--radius-lg) !important;
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
  color: rgb(var(--ocean-surface));
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
  color: rgba(255, 255, 255, 0.6);
}

/* Tips Section - 潮汐阶段 */
.tips-section {
  margin-top: 48px;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth) 0.3s;
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
  padding: 16px 18px;
  border-radius: var(--radius-lg);
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  transition: all 0.3s var(--ease-smooth);
  margin-bottom: 16px;
}

.tip-card:hover {
  border-color: rgba(56, 189, 248, 0.2);
  transform: translateY(-2px);
}

.tip-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.tip-contract .tip-icon {
  background: linear-gradient(135deg, rgba(244, 63, 94, 0.2), rgba(244, 63, 94, 0.05));
  border: 1px solid rgba(244, 63, 94, 0.2);
}

.tip-hold .tip-icon {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(245, 158, 11, 0.05));
  border: 1px solid rgba(245, 158, 11, 0.2);
}

.tip-relax .tip-icon {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(16, 185, 129, 0.05));
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.tip-content {
  flex: 1;
  min-width: 0;
}

.tip-title {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 4px;
}

.tip-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.55);
  margin: 0;
}

/* Result Dialog - 海洋风格 */
.result-dialog :deep(.el-dialog) {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.15);
  border-radius: var(--radius-2xl) !important;
}

.result-dialog :deep(.el-dialog__body) {
  padding: 48px;
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
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  box-shadow: 0 10px 30px rgba(14, 165, 233, 0.3);
}

.result-wave-svg {
  width: 40px;
  height: 40px;
  color: white;
  animation: wave-flow 2s ease-in-out infinite;
}

@keyframes wave-flow {
  0%, 100% { transform: translateX(-3px); }
  50% { transform: translateX(3px); }
}

.result-icon.error {
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.result-title {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
}

.result-subtitle {
  color: rgba(255, 255, 255, 0.65);
  margin: 0 0 24px;
}

.result-stats {
  margin-bottom: 32px;
}

.result-stat {
  padding: 18px 16px;
  border-radius: var(--radius-lg);
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.1);
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  margin-bottom: 6px;
}

.stat-value.orange {
  color: rgb(var(--sunset-amber));
}

.stat-value.purple {
  color: rgb(var(--ocean-surface));
}

.stat-value.green {
  color: rgb(var(--seaweed-green));
}

.stat-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

.login-link {
  color: rgb(var(--ocean-surface));
  text-decoration: none;
  transition: color 0.2s var(--ease-smooth);
}

.login-link:hover {
  color: rgb(var(--aqua-glow));
}

.result-btn {
  padding: 14px 36px !important;
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
  box-shadow: 0 4px 20px rgba(14, 165, 233, 0.3);
}

.result-btn:hover {
  box-shadow: 0 8px 30px rgba(14, 165, 233, 0.4);
}

/* Difficulty Card */
.difficulty-card .title-icon {
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.2), rgba(245, 158, 11, 0.1));
}

.difficulty-title-icon span {
  font-size: 16px;
}

/* Difficulty Grid */
.difficulty-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.difficulty-option {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  border-radius: var(--radius-lg);
  background: rgba(56, 189, 248, 0.03);
  border: 1px solid rgba(56, 189, 248, 0.1);
  cursor: pointer;
  transition: all 0.3s var(--ease-smooth);
  text-align: left;
}

.difficulty-option:hover {
  background: rgba(56, 189, 248, 0.08);
  border-color: rgba(56, 189, 248, 0.2);
  transform: translateX(4px);
}

.difficulty-option.active {
  background: rgba(56, 189, 248, 0.12);
  border-color: rgba(56, 189, 248, 0.35);
  box-shadow: 0 0 20px rgba(56, 189, 248, 0.15);
}

.difficulty-option-icon {
  font-size: 24px;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  background: rgba(255, 255, 255, 0.05);
  flex-shrink: 0;
}

.difficulty-option.active .difficulty-option-icon {
  background: rgba(56, 189, 248, 0.15);
}

.difficulty-option-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.difficulty-option-name {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
}

.difficulty-option-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.difficulty-option-params {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.45);
  flex-shrink: 0;
}

.difficulty-option.active .difficulty-option-params {
  color: rgb(var(--ocean-surface));
}

/* Random Option */
.random-option {
  border-style: dashed;
}

.random-option .difficulty-option-icon {
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.15), rgba(245, 158, 11, 0.05));
}

.random-option.active {
  border-color: rgba(251, 191, 36, 0.4);
  background: rgba(251, 191, 36, 0.1);
}

.random-option .difficulty-option-params {
  color: rgba(251, 191, 36, 0.7);
}

/* Difficulty Preview */
.difficulty-preview {
  padding: 16px;
  border-radius: var(--radius-lg);
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.1);
  margin-bottom: 16px;
}

.preview-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 12px;
  text-align: center;
}

.preview-params {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  flex-wrap: wrap;
}

.preview-param {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  background: rgba(255, 255, 255, 0.03);
}

.param-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.45);
}

.param-value {
  font-size: 18px;
  font-weight: 600;
  color: rgb(var(--ocean-surface));
}

.preview-divider {
  color: rgba(255, 255, 255, 0.3);
  font-size: 14px;
}

/* Difficulty Tip */
.difficulty-tip {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  text-align: center;
  margin: 0;
  line-height: 1.6;
}

@media (max-width: 480px) {
  .difficulty-option {
    padding: 12px 14px;
  }

  .difficulty-option-icon {
    width: 38px;
    height: 38px;
    font-size: 20px;
  }

  .preview-params {
    gap: 4px;
  }

  .preview-param {
    padding: 6px 8px;
  }

  .param-value {
    font-size: 16px;
  }
}
</style>
