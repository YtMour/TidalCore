<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Heatmap from '@/components/Heatmap.vue'
import { getGlobalHeatmap, getLeaderboard } from '@/api/checkin'
import { Timer, Pointer, Trophy, Lock, Lightning, StarFilled, User } from '@element-plus/icons-vue'

const heatmapData = ref<Record<string, number>>()
const loading = ref(true)
const mounted = ref(false)

const features = [
  {
    icon: Timer,
    title: '潮汐计时',
    description: '如潮汐般规律的收缩-保持-放松循环，专业节奏引导',
    color: '#38bdf8',
    bg: 'linear-gradient(135deg, rgba(56, 189, 248, 0.2), rgba(14, 165, 233, 0.1))'
  },
  {
    icon: Pointer,
    title: '浪迹打卡',
    description: '记录每一次训练的浪花，见证坚持的轨迹',
    color: '#22d3ee',
    bg: 'linear-gradient(135deg, rgba(34, 211, 238, 0.2), rgba(56, 189, 248, 0.1))'
  },
  {
    icon: Trophy,
    title: '深海排行',
    description: '匿名排行榜，与全站用户一起潜入深海',
    color: '#fbbf24',
    bg: 'linear-gradient(135deg, rgba(251, 191, 36, 0.2), rgba(245, 158, 11, 0.1))'
  },
  {
    icon: Lock,
    title: '隐私守护',
    description: '如深海般守护隐私，仅需用户名密码',
    color: '#34d399',
    bg: 'linear-gradient(135deg, rgba(52, 211, 153, 0.2), rgba(16, 185, 129, 0.1))'
  },
  {
    icon: Lightning,
    title: '轻盈高效',
    description: '每天几分钟，如海风般轻盈自在',
    color: '#0ea5e9',
    bg: 'linear-gradient(135deg, rgba(14, 165, 233, 0.2), rgba(2, 132, 199, 0.1))'
  },
  {
    icon: StarFilled,
    title: '健康习惯',
    description: '科学的训练方法，建立如潮汐般规律的习惯',
    color: '#f472b6',
    bg: 'linear-gradient(135deg, rgba(244, 114, 182, 0.2), rgba(236, 72, 153, 0.1))'
  }
]

const stats = ref([
  { value: '-', label: '活跃用户', icon: User, suffix: '' },
  { value: '-', label: '累计打卡', icon: Pointer, suffix: '' },
  { value: '-', label: '最长连续', icon: Trophy, suffix: '天' }
])

// 格式化数字显示
function formatNumber(num: number): { value: string, suffix: string } {
  if (num >= 10000) {
    return { value: (num / 10000).toFixed(1).replace(/\.0$/, ''), suffix: '万' }
  } else if (num >= 1000) {
    return { value: (num / 1000).toFixed(1).replace(/\.0$/, ''), suffix: 'k' }
  }
  return { value: num.toString(), suffix: '' }
}

onMounted(async () => {
  setTimeout(() => {
    mounted.value = true
  }, 100)

  try {
    // 并行获取热力图和排行榜数据
    const [heatmapResult, leaderboardResult] = await Promise.all([
      getGlobalHeatmap(365),
      getLeaderboard(100) // 获取更多用户以计算真实活跃用户数
    ])

    heatmapData.value = heatmapResult

    // 从热力图数据计算总打卡次数
    const totalCheckins = Object.values(heatmapResult || {}).reduce((sum, count) => sum + count, 0)

    // 从排行榜获取真实用户数和最长连续天数
    const activeUsers = leaderboardResult?.length || 0
    const maxStreak = leaderboardResult?.length > 0
      ? Math.max(...leaderboardResult.map(u => u.max_streak || 0))
      : 0

    const formatted1 = formatNumber(activeUsers)
    const formatted2 = formatNumber(totalCheckins)

    stats.value = [
      { value: formatted1.value, label: '活跃用户', icon: User, suffix: formatted1.suffix },
      { value: formatted2.value, label: '累计打卡', icon: Pointer, suffix: formatted2.suffix },
      { value: maxStreak.toString(), label: '最长连续', icon: Trophy, suffix: '天' }
    ]
  } catch {
    // 静默处理错误
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <MainLayout>
    <div class="home-page">
      <!-- Hero Section - TidalCore 品牌设计 -->
      <section class="hero-section">
        <!-- 增强的装饰背景 -->
        <div class="hero-bg">
          <!-- 潮汐波纹效果 -->
          <div class="tidal-ripple tidal-ripple-1"></div>
          <div class="tidal-ripple tidal-ripple-2"></div>
          <div class="tidal-ripple tidal-ripple-3"></div>

          <!-- 极光渐变 -->
          <div class="hero-aurora"></div>

          <!-- 装饰光斑 -->
          <div class="hero-blob hero-blob-1"></div>
          <div class="hero-blob hero-blob-2"></div>
        </div>

        <div class="hero-content" :class="{ mounted }">
          <!-- TidalCore Logo 动画 -->
          <div class="hero-logo" :class="{ mounted }">
            <div class="logo-wrapper">
              <!-- 核心图标 -->
              <div class="core-icon">
                <div class="core-ring core-ring-1"></div>
                <div class="core-ring core-ring-2"></div>
                <div class="core-ring core-ring-3"></div>
                <div class="core-center">
                  <svg viewBox="0 0 24 24" fill="none" class="core-wave">
                    <path d="M2 12C2 12 5 8 8 12C11 16 14 8 17 12C20 16 22 12 22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  </svg>
                </div>
              </div>
              <div class="logo-glow"></div>
            </div>
          </div>

          <!-- 标题 - TidalCore 潮汐核心 -->
          <h1 class="hero-title" :class="{ mounted }">
            <span class="title-tidal">Tidal</span><span class="title-core">Core</span>
          </h1>

          <!-- 副标题 -->
          <p class="hero-tagline" :class="{ mounted }">潮汐核心 · 盆底肌训练平台</p>

          <!-- 特性标签 -->
          <div class="hero-badge" :class="{ mounted }">
            <span class="badge-wave"></span>
            <span>开源 · 免费 · 隐私优先</span>
          </div>

          <!-- 描述文字 -->
          <p class="hero-subtitle" :class="{ mounted }">
            如潮汐般规律，如海浪般自然。
            <br class="hidden-mobile" />
            <span class="highlight">每天几分钟，随潮汐一起呼吸。</span>
          </p>

          <!-- CTA 按钮 -->
          <div class="hero-actions" :class="{ mounted }">
            <RouterLink to="/train">
              <el-button type="primary" size="large" class="cta-primary">
                <el-icon><Timer /></el-icon>
                <span>开始训练</span>
                <el-icon class="arrow-icon"><ArrowRight /></el-icon>
              </el-button>
            </RouterLink>
            <RouterLink to="/leaderboard">
              <el-button size="large" class="cta-secondary">
                <el-icon :style="{ color: '#fbbf24' }"><Trophy /></el-icon>
                <span>查看排行</span>
              </el-button>
            </RouterLink>
          </div>

          <!-- 统计数据 -->
          <el-row :gutter="16" class="stats-row" :class="{ mounted }">
            <el-col :xs="8" :sm="8" v-for="stat in stats" :key="stat.label">
              <div class="stat-card">
                <el-icon :size="20" class="stat-icon"><component :is="stat.icon" /></el-icon>
                <div class="stat-info">
                  <div class="stat-value">{{ stat.value }}<span class="stat-suffix">{{ stat.suffix }}</span></div>
                  <div class="stat-label">{{ stat.label }}</div>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
      </section>

      <!-- Global Heatmap Section -->
      <section class="heatmap-section" :class="{ mounted }">
        <el-card class="heatmap-card" shadow="never">
          <div class="card-decoration"></div>

          <div class="card-content">
            <div class="card-header">
              <div class="header-left">
                <div class="header-icon">
                  <el-icon :size="20" color="#34d399"><Pointer /></el-icon>
                </div>
                <div class="header-text">
                  <h2>全站打卡热力图</h2>
                  <p>社区成员的坚持轨迹</p>
                </div>
              </div>
              <div class="header-right">
                <!-- 统计信息 -->
                <div class="hm-stats" v-if="heatmapData">
                  <div class="stat-item">
                    <span class="stat-value">{{ Object.values(heatmapData || {}).reduce((sum, count) => sum + count, 0) }}</span>
                    <span class="stat-label">次打卡</span>
                  </div>
                  <div class="stat-divider"></div>
                  <div class="stat-item">
                    <span class="stat-value">{{ Object.values(heatmapData || {}).filter(count => count > 0).length }}</span>
                    <span class="stat-label">天活跃</span>
                  </div>
                </div>
                <span class="live-indicator">
                  <span class="live-dot"></span>
                  实时
                </span>
              </div>
            </div>

            <div v-if="loading" class="heatmap-loading">
              <el-skeleton :rows="3" animated />
            </div>
            <Heatmap v-else :data="heatmapData || {}" :days="365" />
          </div>
        </el-card>
      </section>

      <!-- Features Section -->
      <section class="features-section" :class="{ mounted }">
        <div class="section-header">
          <div class="section-eyebrow section-eyebrow-ocean">
            <span class="eyebrow-line"></span>
            <span class="eyebrow-text">核心功能</span>
            <span class="eyebrow-line"></span>
          </div>
          <h2 class="section-title">为什么选择 TidalCore</h2>
          <p class="section-subtitle">如潮汐般自然的训练体验，帮助你养成规律的健康习惯</p>
        </div>

        <el-row :gutter="24">
          <el-col :xs="24" :sm="12" :lg="8" v-for="feature in features" :key="feature.title">
            <el-card class="feature-card" shadow="hover">
              <div class="feature-icon" :style="{ background: feature.bg }">
                <el-icon :size="24" :style="{ color: feature.color }">
                  <component :is="feature.icon" />
                </el-icon>
              </div>
              <div class="feature-content">
                <h3 class="feature-title">{{ feature.title }}</h3>
                <p class="feature-desc">{{ feature.description }}</p>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </section>

      <!-- How It Works Section -->
      <section class="howto-section" :class="{ mounted }">
        <div class="section-header">
          <div class="section-eyebrow section-eyebrow-aqua">
            <span class="eyebrow-line"></span>
            <span class="eyebrow-text">快速上手</span>
            <span class="eyebrow-line"></span>
          </div>
          <h2 class="section-title">三步开启潮汐之旅</h2>
          <p class="section-subtitle">无需复杂设置，跟随潮汐节奏开始健康训练</p>
        </div>

        <el-row :gutter="32">
          <el-col :xs="24" :md="8" v-for="(step, index) in [
            { num: 1, title: '注册账号', desc: '仅需用户名和密码，30秒即可入海', color: ['#0ea5e9', '#22d3ee'] },
            { num: 2, title: '跟随潮汐', desc: '跟随计时器节奏，完成每日训练', color: ['#38bdf8', '#0284c7'] },
            { num: 3, title: '浪迹打卡', desc: '记录进度，冲击深海排行榜', color: ['#06b6d4', '#0891b2'] }
          ]" :key="step.num">
            <div class="step-item">
              <div class="step-num" :style="{ background: `linear-gradient(135deg, ${step.color[0]}, ${step.color[1]})` }">
                {{ step.num }}
              </div>
              <h3 class="step-title">{{ step.title }}</h3>
              <p class="step-desc">{{ step.desc }}</p>
              <div v-if="index < 2" class="step-connector"></div>
            </div>
          </el-col>
        </el-row>
      </section>

      <!-- Bottom CTA -->
      <section class="cta-section">
        <el-card class="cta-card" shadow="never">
          <div class="cta-decoration cta-decoration-1"></div>
          <div class="cta-decoration cta-decoration-2"></div>
          <div class="cta-wave-bg"></div>

          <div class="cta-content">
            <div class="cta-icon">
              <svg viewBox="0 0 32 32" fill="none" class="cta-wave-svg">
                <path d="M4 16C4 16 8 10 14 16C20 22 26 10 28 16" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
              </svg>
            </div>
            <h3 class="cta-title">准备好入海了吗？</h3>
            <p class="cta-text">加入数千名用户，开启你的潮汐训练之旅。完全免费，永久开源。</p>
            <div class="cta-buttons">
              <RouterLink to="/register">
                <el-button type="primary" size="large" class="cta-primary">
                  立即注册，免费使用
                  <el-icon class="arrow-icon"><ArrowRight /></el-icon>
                </el-button>
              </RouterLink>
              <RouterLink to="/train">
                <el-button size="large" class="cta-secondary">先体验一下</el-button>
              </RouterLink>
            </div>
          </div>
        </el-card>
      </section>
    </div>
  </MainLayout>
</template>

<script lang="ts">
import { ArrowRight } from '@element-plus/icons-vue'
export default {
  components: { ArrowRight }
}
</script>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 96px;
}

/* ===== Hero Section - 海洋潮汐主题 ===== */
.hero-section {
  position: relative;
  padding: 80px 0 120px;
  text-align: center;
  overflow: hidden;
}

.hero-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

/* 潮汐波纹效果 */
.tidal-ripple {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  border: 1px solid rgba(56, 189, 248, 0.15);
  animation: tidal-expand 8s ease-out infinite;
}

.tidal-ripple-1 {
  width: 300px;
  height: 300px;
  animation-delay: 0s;
}

.tidal-ripple-2 {
  width: 500px;
  height: 500px;
  animation-delay: 2s;
}

.tidal-ripple-3 {
  width: 700px;
  height: 700px;
  animation-delay: 4s;
}

@keyframes tidal-expand {
  0% {
    transform: translate(-50%, -50%) scale(0.5);
    opacity: 0.6;
  }
  100% {
    transform: translate(-50%, -50%) scale(2);
    opacity: 0;
  }
}

/* 海洋极光背景 */
.hero-aurora {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse 80% 60% at 50% -20%, rgba(56, 189, 248, 0.15) 0%, transparent 60%),
    radial-gradient(ellipse 60% 50% at 100% 50%, rgba(34, 211, 238, 0.1) 0%, transparent 50%),
    radial-gradient(ellipse 60% 50% at 0% 80%, rgba(14, 165, 233, 0.08) 0%, transparent 50%);
}

.hero-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
}

.hero-blob-1 {
  top: 10%;
  left: 20%;
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.25), rgba(34, 211, 238, 0.15));
  animation: float 8s ease-in-out infinite;
}

.hero-blob-2 {
  bottom: 10%;
  right: 20%;
  width: 350px;
  height: 350px;
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.2), rgba(2, 132, 199, 0.1));
  animation: float 10s ease-in-out infinite;
  animation-delay: -3s;
}

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-25px) rotate(3deg); }
}

.hero-content {
  position: relative;
  z-index: 10;
}

/* TidalCore Logo 核心动画 */
.hero-logo {
  margin-bottom: 32px;
  opacity: 0;
  transform: translateY(32px) scale(0.9);
  transition: all 1s var(--ease-smooth);
}

.hero-logo.mounted {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.logo-wrapper {
  position: relative;
  display: inline-block;
}

.core-icon {
  position: relative;
  width: 120px;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.core-ring {
  position: absolute;
  border-radius: 50%;
  border: 2px solid;
  animation: core-pulse 3s ease-in-out infinite;
}

.core-ring-1 {
  width: 100%;
  height: 100%;
  border-color: rgba(56, 189, 248, 0.4);
  animation-delay: 0s;
}

.core-ring-2 {
  width: 80%;
  height: 80%;
  border-color: rgba(34, 211, 238, 0.5);
  animation-delay: 0.3s;
}

.core-ring-3 {
  width: 60%;
  height: 60%;
  border-color: rgba(14, 165, 233, 0.6);
  animation-delay: 0.6s;
}

@keyframes core-pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.6;
  }
  50% {
    transform: scale(1.1);
    opacity: 1;
  }
}

.core-center {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #38bdf8, #0ea5e9);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    0 0 30px rgba(56, 189, 248, 0.5),
    0 0 60px rgba(56, 189, 248, 0.3),
    inset 0 0 20px rgba(255, 255, 255, 0.2);
  animation: core-glow 2s ease-in-out infinite;
}

.core-wave {
  width: 28px;
  height: 28px;
  color: white;
  animation: wave-flow 2s ease-in-out infinite;
}

@keyframes wave-flow {
  0%, 100% { transform: translateX(-2px); }
  50% { transform: translateX(2px); }
}

@keyframes core-glow {
  0%, 100% {
    box-shadow:
      0 0 30px rgba(56, 189, 248, 0.5),
      0 0 60px rgba(56, 189, 248, 0.3);
  }
  50% {
    box-shadow:
      0 0 40px rgba(56, 189, 248, 0.7),
      0 0 80px rgba(56, 189, 248, 0.4);
  }
}

.logo-glow {
  position: absolute;
  inset: -20px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.3) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
  animation: glow-pulse 4s ease-in-out infinite;
}

@keyframes glow-pulse {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.1); }
}

/* 标题样式 - 海洋渐变 */
.hero-title {
  font-size: 72px;
  font-weight: 800;
  margin-bottom: 16px;
  letter-spacing: -0.03em;
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s var(--ease-smooth) 0.15s;
}

.hero-content.mounted .hero-title {
  opacity: 1;
  transform: translateY(0);
}

.title-tidal {
  background: linear-gradient(135deg, #38bdf8, #22d3ee, #0ea5e9);
  background-size: 200% auto;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: ocean-shimmer 4s ease-in-out infinite;
}

.title-core {
  background: linear-gradient(135deg, #0ea5e9, #0284c7, #0369a1);
  background-size: 200% auto;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: ocean-shimmer 4s ease-in-out infinite;
  animation-delay: 0.5s;
}

@keyframes ocean-shimmer {
  0%, 100% { background-position: 0% center; }
  50% { background-position: 200% center; }
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 48px;
  }

  .core-icon {
    width: 100px;
    height: 100px;
  }

  .core-center {
    width: 42px;
    height: 42px;
  }
}

/* 副标题 */
.hero-tagline {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 24px;
  letter-spacing: 0.1em;
  opacity: 0;
  transform: translateY(20px);
  transition: all 0.8s var(--ease-smooth) 0.2s;
}

.hero-content.mounted .hero-tagline {
  opacity: 1;
  transform: translateY(0);
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 10px 18px;
  border-radius: var(--radius-full);
  background: rgba(56, 189, 248, 0.1);
  border: 1px solid rgba(56, 189, 248, 0.2);
  margin-bottom: 24px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  opacity: 0;
  transform: translateY(32px);
  transition: all 0.7s var(--ease-smooth) 0.25s;
}

.hero-content.mounted .hero-badge {
  opacity: 1;
  transform: translateY(0);
}

.badge-wave {
  width: 20px;
  height: 8px;
  background: linear-gradient(90deg, #38bdf8, #22d3ee, #38bdf8);
  background-size: 200% 100%;
  border-radius: 4px;
  animation: wave-badge 2s ease-in-out infinite;
}

@keyframes wave-badge {
  0%, 100% { background-position: 0% center; }
  50% { background-position: 200% center; }
}

.hero-subtitle {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.6);
  max-width: 600px;
  margin: 0 auto 48px;
  line-height: 1.8;
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s var(--ease-smooth) 0.3s;
}

.hero-content.mounted .hero-subtitle {
  opacity: 1;
  transform: translateY(0);
}

@media (max-width: 768px) {
  .hero-subtitle {
    font-size: 16px;
  }
}

.hero-subtitle .highlight {
  color: rgb(var(--aqua-glow));
}

.hidden-mobile {
  display: none;
}

@media (min-width: 640px) {
  .hidden-mobile {
    display: block;
  }
}

.hero-actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  margin-bottom: 64px;
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s var(--ease-smooth) 0.4s;
}

.hero-content.mounted .hero-actions {
  opacity: 1;
  transform: translateY(0);
}

@media (min-width: 640px) {
  .hero-actions {
    flex-direction: row;
    justify-content: center;
  }
}

.cta-primary {
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
  padding: 14px 32px !important;
  font-size: 16px !important;
  height: auto !important;
  border-radius: var(--radius-lg) !important;
  box-shadow: 0 4px 20px rgba(14, 165, 233, 0.3);
  transition: all 0.3s var(--ease-smooth) !important;
}

.cta-primary:hover {
  box-shadow: 0 8px 30px rgba(14, 165, 233, 0.4);
  transform: translateY(-2px);
}

.cta-primary .arrow-icon {
  margin-left: 8px;
  transition: transform 0.2s var(--ease-smooth);
}

.cta-primary:hover .arrow-icon {
  transform: translateX(4px);
}

.cta-secondary {
  background: rgba(56, 189, 248, 0.08) !important;
  border: 1px solid rgba(56, 189, 248, 0.15) !important;
  padding: 14px 32px !important;
  font-size: 16px !important;
  height: auto !important;
  border-radius: var(--radius-lg) !important;
  color: #fff !important;
}

.cta-secondary:hover {
  background: rgba(56, 189, 248, 0.12) !important;
  border-color: rgba(56, 189, 248, 0.25) !important;
}

.stats-row {
  max-width: 768px;
  margin: 0 auto !important;
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s var(--ease-smooth) 0.5s;
}

.hero-content.mounted .stats-row {
  opacity: 1;
  transform: translateY(0);
}

.stat-card {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
  padding: 16px 18px;
  border-radius: var(--radius-lg);
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.1);
  transition: all 0.3s var(--ease-smooth);
  margin-bottom: 16px;
}

.stat-card:hover {
  background: rgba(56, 189, 248, 0.1);
  border-color: rgba(56, 189, 248, 0.2);
  transform: translateY(-2px);
}

.stat-icon {
  color: rgb(var(--aqua-glow));
  flex-shrink: 0;
}

.stat-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  line-height: 1.2;
}

.stat-suffix {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  margin-left: 2px;
}

@media (max-width: 768px) {
  .stat-value {
    font-size: 18px;
  }
  .stat-suffix {
    font-size: 12px;
  }
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

/* ===== Heatmap Section - 海洋风格 ===== */
.heatmap-section {
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s var(--ease-smooth) 0.6s;
}

.heatmap-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.heatmap-card {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.1) !important;
  border-radius: var(--radius-xl) !important;
  overflow: hidden;
  position: relative;
}

.heatmap-card :deep(.el-card__body) {
  padding: 32px;
}

.card-decoration {
  position: absolute;
  top: 0;
  right: 0;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.1) 0%, transparent 70%);
  filter: blur(60px);
  pointer-events: none;
}

.card-content {
  position: relative;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32px;
  flex-wrap: wrap;
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.header-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgba(52, 211, 153, 0.2), rgba(16, 185, 129, 0.1));
  border: 1px solid rgba(52, 211, 153, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-text h2 {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 4px;
}

.header-text p {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.4);
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 热力图统计信息 */
.hm-stats {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.12), rgba(56, 189, 248, 0.04));
  border: 1px solid rgba(56, 189, 248, 0.2);
  border-radius: var(--radius-lg);
}

.hm-stats .stat-item {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.hm-stats .stat-value {
  font-size: 18px;
  font-weight: 700;
  color: rgb(var(--ocean-surface));
}

.hm-stats .stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.hm-stats .stat-divider {
  width: 1px;
  height: 20px;
  background: rgba(56, 189, 248, 0.2);
}

.live-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: rgb(var(--seaweed-green));
  padding: 8px 14px;
  border-radius: var(--radius-full);
  background: rgba(52, 211, 153, 0.1);
  border: 1px solid rgba(52, 211, 153, 0.15);
}

.live-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgb(var(--seaweed-green));
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(1.2); }
}

.heatmap-loading {
  padding: 32px 0;
}

/* ===== Features Section - 海洋风格 ===== */
.features-section {
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s var(--ease-smooth) 0.7s;
}

.features-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.section-header {
  text-align: center;
  margin-bottom: 64px;
}

.section-eyebrow {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 20px;
}

.section-eyebrow .eyebrow-line {
  width: 40px;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(56, 189, 248, 0.5));
}

.section-eyebrow .eyebrow-line:last-child {
  background: linear-gradient(90deg, rgba(56, 189, 248, 0.5), transparent);
}

.section-eyebrow .eyebrow-text {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 3px;
  background: linear-gradient(135deg, rgb(var(--ocean-surface)), rgb(var(--aqua-glow)));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.section-eyebrow-ocean .eyebrow-line {
  background: linear-gradient(90deg, transparent, rgb(var(--ocean-surface)));
}

.section-eyebrow-ocean .eyebrow-line:last-child {
  background: linear-gradient(90deg, rgb(var(--ocean-surface)), transparent);
}

.section-eyebrow-aqua .eyebrow-line {
  background: linear-gradient(90deg, transparent, rgb(var(--aqua-glow)));
}

.section-eyebrow-aqua .eyebrow-line:last-child {
  background: linear-gradient(90deg, rgb(var(--aqua-glow)), transparent);
}

.section-eyebrow-aqua .eyebrow-text {
  background: linear-gradient(135deg, rgb(var(--aqua-glow)), rgb(var(--ocean-surface)));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.section-title {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 16px;
}

.section-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.5);
  max-width: 512px;
  margin: 0 auto;
}

.feature-card {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(16px);
  border: 1px solid rgba(56, 189, 248, 0.1) !important;
  border-radius: var(--radius-xl) !important;
  margin-bottom: 24px;
  transition: all 0.4s var(--ease-smooth) !important;
  position: relative;
  overflow: hidden;
}

.feature-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(56, 189, 248, 0.3), transparent);
  opacity: 0;
  transition: opacity 0.3s var(--ease-smooth);
}

.feature-card:hover {
  border-color: rgba(56, 189, 248, 0.2) !important;
  transform: translateY(-4px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.feature-card:hover::before {
  opacity: 1;
}

.feature-card :deep(.el-card__body) {
  padding: 24px;
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  gap: 16px;
}

.feature-icon {
  width: 52px;
  height: 52px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform 0.3s var(--ease-smooth);
}

.feature-card:hover .feature-icon {
  transform: scale(1.1);
}

.feature-content {
  flex: 1;
  min-width: 0;
}

.feature-title {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 6px;
  transition: color 0.2s var(--ease-smooth);
}

.feature-card:hover .feature-title {
  color: rgb(var(--ocean-surface));
}

.feature-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  line-height: 1.6;
  margin: 0;
}

/* ===== How To Section - 海洋风格 ===== */
.howto-section {
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s var(--ease-smooth) 0.8s;
}

.howto-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.step-item {
  position: relative;
  text-align: center;
  margin-bottom: 32px;
}

.step-num {
  width: 60px;
  height: 60px;
  margin: 0 auto 20px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  box-shadow: 0 10px 30px rgba(14, 165, 233, 0.3);
  position: relative;
}

.step-num::after {
  content: '';
  position: absolute;
  inset: -4px;
  border-radius: inherit;
  border: 1px solid rgba(56, 189, 248, 0.2);
  animation: step-pulse 2s ease-in-out infinite;
}

@keyframes step-pulse {
  0%, 100% { transform: scale(1); opacity: 0.5; }
  50% { transform: scale(1.1); opacity: 0; }
}

.step-title {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 8px;
}

.step-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  margin: 0;
}

.step-connector {
  display: none;
}

@media (min-width: 768px) {
  .step-connector {
    display: block;
    position: absolute;
    top: 32px;
    left: 60%;
    width: 80%;
    height: 2px;
    background: linear-gradient(to right, rgba(56, 189, 248, 0.4), transparent);
  }
}

/* ===== CTA Section - 海洋风格 ===== */
.cta-section {
  padding-bottom: 64px;
}

.cta-card {
  max-width: 800px;
  margin: 0 auto;
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.15) !important;
  border-radius: var(--radius-2xl) !important;
  overflow: hidden;
  position: relative;
}

.cta-card :deep(.el-card__body) {
  padding: 56px;
}

.cta-decoration {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
}

.cta-decoration-1 {
  top: -100px;
  left: -100px;
  width: 250px;
  height: 250px;
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.25), transparent);
}

.cta-decoration-2 {
  bottom: -100px;
  right: -100px;
  width: 250px;
  height: 250px;
  background: linear-gradient(135deg, rgba(34, 211, 238, 0.2), transparent);
}

.cta-wave-bg {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 100px;
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 100'%3E%3Cpath fill='rgba(56, 189, 248, 0.05)' d='M0,50 C360,100 720,0 1080,50 C1260,75 1380,25 1440,50 L1440,100 L0,100 Z'/%3E%3C/svg%3E");
  background-size: 1440px 100px;
  background-repeat: repeat-x;
  animation: wave-move 15s linear infinite;
}

@keyframes wave-move {
  0% { background-position-x: 0; }
  100% { background-position-x: 1440px; }
}

.cta-content {
  position: relative;
  text-align: center;
}

.cta-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 28px;
  border-radius: var(--radius-xl);
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(14, 165, 233, 0.3);
}

.cta-wave-svg {
  width: 32px;
  height: 32px;
  color: white;
  animation: wave-flow 2s ease-in-out infinite;
}

.cta-title {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 16px;
}

@media (max-width: 768px) {
  .cta-title {
    font-size: 24px;
  }
}

.cta-text {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.5);
  max-width: 480px;
  margin: 0 auto 36px;
}

.cta-buttons {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

@media (min-width: 640px) {
  .cta-buttons {
    flex-direction: row;
    justify-content: center;
  }
}
</style>
