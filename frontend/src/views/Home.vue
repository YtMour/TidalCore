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
    description: '如潮汐般规律的收缩-保持-放松循环',
    color: '#38bdf8'
  },
  {
    icon: Pointer,
    title: '浪迹打卡',
    description: '记录每一次训练的浪花',
    color: '#22d3ee'
  },
  {
    icon: Trophy,
    title: '深海排行',
    description: '与全站用户一起潜入深海',
    color: '#fbbf24'
  },
  {
    icon: Lock,
    title: '隐私守护',
    description: '如深海般守护你的隐私',
    color: '#34d399'
  },
  {
    icon: Lightning,
    title: '轻盈高效',
    description: '每天几分钟，如海风般自在',
    color: '#0ea5e9'
  },
  {
    icon: StarFilled,
    title: '健康习惯',
    description: '建立如潮汐般规律的习惯',
    color: '#f472b6'
  }
]

const stats = ref([
  { value: '-', label: '活跃用户', icon: User, suffix: '' },
  { value: '-', label: '累计打卡', icon: Pointer, suffix: '' },
  { value: '-', label: '最长连续', icon: Trophy, suffix: '天' }
])

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
    const [heatmapResult, leaderboardResult] = await Promise.all([
      getGlobalHeatmap(365),
      getLeaderboard(100)
    ])

    heatmapData.value = heatmapResult

    const totalCheckins = Object.values(heatmapResult || {}).reduce((sum, count) => sum + count, 0)
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
      <!-- Hero Section - 全屏沉浸式 -->
      <section class="hero-section">
        <div class="hero-content" :class="{ mounted }">
          <!-- TidalCore Logo -->
          <div class="hero-logo" :class="{ mounted }">
            <div class="logo-container">
              <!-- 涟漪效果 -->
              <div class="logo-ripples">
                <div class="ripple ripple-1"></div>
                <div class="ripple ripple-2"></div>
                <div class="ripple ripple-3"></div>
              </div>
              <!-- 核心圆 -->
              <div class="logo-core">
                <!-- 海浪涌动SVG -->
                <svg viewBox="0 0 60 60" fill="none" class="core-wave-svg">
                  <defs>
                    <linearGradient id="coreWaveGrad" x1="0%" y1="100%" x2="0%" y2="0%">
                      <stop offset="0%" stop-color="#0ea5e9" />
                      <stop offset="100%" stop-color="#38bdf8" />
                    </linearGradient>
                    <clipPath id="coreClip">
                      <circle cx="30" cy="30" r="26" />
                    </clipPath>
                  </defs>
                  <!-- 海浪填充效果 - 从下往上涌动 -->
                  <g clip-path="url(#coreClip)">
                    <path class="wave-fill wave-fill-1" d="M0 35 Q15 30, 30 35 T60 35 L60 60 L0 60 Z" fill="url(#coreWaveGrad)" opacity="0.9"/>
                    <path class="wave-fill wave-fill-2" d="M0 40 Q15 35, 30 40 T60 40 L60 60 L0 60 Z" fill="#22d3ee" opacity="0.6"/>
                    <path class="wave-fill wave-fill-3" d="M0 45 Q15 40, 30 45 T60 45 L60 60 L0 60 Z" fill="#06b6d4" opacity="0.4"/>
                  </g>
                  <!-- 波浪线条 -->
                  <path class="wave-line" d="M4 30 Q17 22, 30 30 T56 30" stroke="white" stroke-width="2" stroke-linecap="round" fill="none" opacity="0.8"/>
                </svg>
                <!-- 光晕效果 -->
                <div class="core-glow"></div>
              </div>
            </div>
          </div>

          <!-- 标题 -->
          <h1 class="hero-title" :class="{ mounted }">
            <span class="title-tidal">Tidal</span><span class="title-core">Core</span>
          </h1>

          <p class="hero-tagline" :class="{ mounted }">潮汐核心 · 盆底肌训练平台</p>

          <!-- 特性标签 -->
          <div class="hero-badge" :class="{ mounted }">
            <span class="badge-wave"></span>
            <span>开源 · 免费 · 隐私优先</span>
          </div>

          <!-- 描述 -->
          <p class="hero-subtitle" :class="{ mounted }">
            如潮汐般规律，如海浪般自然。
            <br />
            <span class="highlight">每天几分钟，随潮汐一起呼吸。</span>
          </p>

          <!-- CTA 按钮 -->
          <div class="hero-actions" :class="{ mounted }">
            <RouterLink to="/train" class="cta-primary">
              <el-icon><Timer /></el-icon>
              <span>开始训练</span>
              <svg viewBox="0 0 24 24" fill="none" class="arrow-icon">
                <path d="M5 12H19M19 12L12 5M19 12L12 19" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </RouterLink>
            <RouterLink to="/leaderboard" class="cta-secondary">
              <el-icon :style="{ color: '#fbbf24' }"><Trophy /></el-icon>
              <span>查看排行</span>
            </RouterLink>
          </div>

          <!-- 统计数据 -->
          <div class="stats-row" :class="{ mounted }">
            <div v-for="stat in stats" :key="stat.label" class="stat-item">
              <el-icon :size="18" class="stat-icon"><component :is="stat.icon" /></el-icon>
              <span class="stat-value">{{ stat.value }}<span class="stat-suffix">{{ stat.suffix }}</span></span>
              <span class="stat-label">{{ stat.label }}</span>
            </div>
          </div>
        </div>

        <!-- 向下滚动指示 -->
        <div class="scroll-indicator" :class="{ mounted }">
          <div class="scroll-mouse">
            <div class="scroll-wheel"></div>
          </div>
          <span>向下滚动探索更多</span>
        </div>
      </section>

      <!-- Features Section - 全屏展示 -->
      <section class="features-section">
        <div class="section-container">
          <div class="section-header">
            <span class="section-eyebrow">核心功能</span>
            <h2 class="section-title">为什么选择 TidalCore</h2>
            <p class="section-desc">如潮汐般自然的训练体验</p>
          </div>

          <div class="features-grid">
            <div v-for="(feature, index) in features" :key="feature.title" class="feature-item" :style="{ animationDelay: index * 0.1 + 's' }">
              <div class="feature-icon" :style="{ '--feature-color': feature.color }">
                <el-icon :size="24"><component :is="feature.icon" /></el-icon>
              </div>
              <h3 class="feature-title">{{ feature.title }}</h3>
              <p class="feature-desc">{{ feature.description }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Heatmap Section -->
      <section class="heatmap-section">
        <div class="section-container">
          <div class="section-header">
            <span class="section-eyebrow">社区动态</span>
            <h2 class="section-title">全站打卡热力图</h2>
            <p class="section-desc">社区成员的坚持轨迹</p>
          </div>

          <div class="heatmap-wrapper">
            <div class="heatmap-stats" v-if="heatmapData">
              <div class="hm-stat">
                <span class="hm-value">{{ Object.values(heatmapData || {}).reduce((sum, count) => sum + count, 0) }}</span>
                <span class="hm-label">次打卡</span>
              </div>
              <div class="hm-divider"></div>
              <div class="hm-stat">
                <span class="hm-value">{{ Object.values(heatmapData || {}).filter(count => count > 0).length }}</span>
                <span class="hm-label">天活跃</span>
              </div>
              <div class="live-badge">
                <span class="live-dot"></span>
                实时
              </div>
            </div>

            <div v-if="loading" class="heatmap-loading">
              <el-skeleton :rows="3" animated />
            </div>
            <Heatmap v-else :data="heatmapData || {}" :days="365" />
          </div>
        </div>
      </section>

      <!-- Steps Section -->
      <section class="steps-section">
        <div class="section-container">
          <div class="section-header">
            <span class="section-eyebrow">快速上手</span>
            <h2 class="section-title">三步开启潮汐之旅</h2>
            <p class="section-desc">无需复杂设置，跟随潮汐节奏</p>
          </div>

          <div class="steps-grid">
            <div class="step-item">
              <div class="step-num">1</div>
              <h3 class="step-title">注册账号</h3>
              <p class="step-desc">仅需用户名和密码，30秒即可入海</p>
            </div>
            <div class="step-connector">
              <svg viewBox="0 0 100 20" preserveAspectRatio="none">
                <path d="M0,10 Q25,5 50,10 T100,10" fill="none" stroke="rgba(56, 189, 248, 0.3)" stroke-width="2"/>
              </svg>
            </div>
            <div class="step-item">
              <div class="step-num">2</div>
              <h3 class="step-title">跟随潮汐</h3>
              <p class="step-desc">跟随计时器节奏，完成每日训练</p>
            </div>
            <div class="step-connector">
              <svg viewBox="0 0 100 20" preserveAspectRatio="none">
                <path d="M0,10 Q25,15 50,10 T100,10" fill="none" stroke="rgba(56, 189, 248, 0.3)" stroke-width="2"/>
              </svg>
            </div>
            <div class="step-item">
              <div class="step-num">3</div>
              <h3 class="step-title">浪迹打卡</h3>
              <p class="step-desc">记录进度，冲击深海排行榜</p>
            </div>
          </div>
        </div>
      </section>

      <!-- CTA Section -->
      <section class="cta-section">
        <div class="cta-container">
          <div class="cta-icon">
            <svg viewBox="0 0 60 60" fill="none" class="cta-wave-svg">
              <defs>
                <linearGradient id="ctaWaveGrad" x1="0%" y1="100%" x2="0%" y2="0%">
                  <stop offset="0%" stop-color="#0284c7" />
                  <stop offset="100%" stop-color="#38bdf8" />
                </linearGradient>
                <clipPath id="ctaClip">
                  <rect x="4" y="4" width="52" height="52" rx="12" />
                </clipPath>
              </defs>
              <!-- 海浪涌动填充 -->
              <g clip-path="url(#ctaClip)">
                <rect x="4" y="4" width="52" height="52" rx="12" fill="rgba(255,255,255,0.1)"/>
                <path class="cta-wave-fill cta-wave-1" d="M0 32 Q15 26, 30 32 T60 32 L60 60 L0 60 Z" fill="url(#ctaWaveGrad)" opacity="0.8"/>
                <path class="cta-wave-fill cta-wave-2" d="M0 38 Q15 32, 30 38 T60 38 L60 60 L0 60 Z" fill="#22d3ee" opacity="0.5"/>
                <path class="cta-wave-fill cta-wave-3" d="M0 44 Q15 38, 30 44 T60 44 L60 60 L0 60 Z" fill="#0ea5e9" opacity="0.3"/>
              </g>
              <!-- 波浪线条 -->
              <path class="cta-wave-line" d="M8 28 Q22 20, 36 28 T52 28" stroke="white" stroke-width="2.5" stroke-linecap="round" fill="none" opacity="0.9"/>
            </svg>
          </div>
          <h2 class="cta-title">准备好入海了吗？</h2>
          <p class="cta-desc">加入数千名用户，开启你的潮汐训练之旅。完全免费，永久开源。</p>
          <div class="cta-buttons">
            <RouterLink to="/register" class="cta-btn-primary">
              立即注册，免费使用
              <svg viewBox="0 0 24 24" fill="none" class="arrow-icon">
                <path d="M5 12H19M19 12L12 5M19 12L12 19" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </RouterLink>
            <RouterLink to="/train" class="cta-btn-secondary">先体验一下</RouterLink>
          </div>
        </div>
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
}

/* ===== Hero Section - 全屏沉浸式 ===== */
.hero-section {
  min-height: calc(100vh - 72px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px 24px;
  position: relative;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
}

/* Logo 动画 */
.hero-logo {
  margin-bottom: 32px;
  opacity: 0;
  transform: translateY(40px) scale(0.8);
  transition: all 1s var(--ease-smooth);
}

.hero-logo.mounted {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.logo-container {
  position: relative;
  width: 140px;
  height: 140px;
  margin: 0 auto;
}

/* 涟漪效果 - 从中心向外扩散 - 速度减慢，解决加载重叠 */
.logo-ripples {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.logo-ripples .ripple {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2px solid rgba(56, 189, 248, 0.35);
  opacity: 0;
  animation: ripple-out 6s ease-out infinite;
}

.logo-ripples .ripple-1 {
  animation-delay: 0.5s;
}

.logo-ripples .ripple-2 {
  animation-delay: 2.5s;
}

.logo-ripples .ripple-3 {
  animation-delay: 4.5s;
}

@keyframes ripple-out {
  0% {
    transform: scale(0.7);
    opacity: 0;
  }
  10% {
    opacity: 0.6;
  }
  100% {
    transform: scale(1.6);
    opacity: 0;
  }
}

/* Logo核心 */
.logo-core {
  position: absolute;
  inset: 20%;
  border-radius: 50%;
  background: linear-gradient(180deg, #0c4a6e 0%, #0369a1 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  box-shadow:
    0 0 40px rgba(56, 189, 248, 0.4),
    0 0 80px rgba(56, 189, 248, 0.2),
    inset 0 -20px 40px rgba(56, 189, 248, 0.3);
}

.core-wave-svg {
  width: 100%;
  height: 100%;
  position: absolute;
  inset: 0;
}

/* 海浪涌动动画 - 上下起伏 */
.core-wave-svg .wave-fill-1 {
  animation: wave-surge-1 3s ease-in-out infinite;
}

.core-wave-svg .wave-fill-2 {
  animation: wave-surge-2 3.5s ease-in-out infinite;
}

.core-wave-svg .wave-fill-3 {
  animation: wave-surge-3 4s ease-in-out infinite;
}

.core-wave-svg .wave-line {
  animation: wave-line-surge 2.5s ease-in-out infinite;
}

@keyframes wave-surge-1 {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-8px);
  }
}

@keyframes wave-surge-2 {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-6px);
  }
}

@keyframes wave-surge-3 {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-4px);
  }
}

@keyframes wave-line-surge {
  0%, 100% {
    transform: translateY(0);
    d: path("M4 30 Q17 22, 30 30 T56 30");
  }
  50% {
    transform: translateY(-3px);
    d: path("M4 30 Q17 38, 30 30 T56 30");
  }
}

/* 核心光晕 */
.core-glow {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  background: radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.3), transparent 60%);
  animation: glow-pulse 3s ease-in-out infinite;
}

@keyframes glow-pulse {
  0%, 100% { opacity: 0.8; }
  50% { opacity: 1; }
}

/* 标题 */
.hero-title {
  font-size: clamp(48px, 10vw, 80px);
  font-weight: 800;
  margin-bottom: 12px;
  letter-spacing: -0.03em;
  opacity: 0;
  transform: translateY(30px);
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
  animation: shimmer 4s ease-in-out infinite;
}

.title-core {
  background: linear-gradient(135deg, #0ea5e9, #0284c7, #0369a1);
  background-size: 200% auto;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: shimmer 4s ease-in-out infinite 0.5s;
}

@keyframes shimmer {
  0%, 100% { background-position: 0% center; }
  50% { background-position: 200% center; }
}

.hero-tagline {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 24px;
  letter-spacing: 0.15em;
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
  gap: 12px;
  padding: 12px 20px;
  border-radius: var(--radius-full);
  background: rgba(56, 189, 248, 0.1);
  border: 1px solid rgba(56, 189, 248, 0.2);
  margin-bottom: 28px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  opacity: 0;
  transform: translateY(30px);
  transition: all 0.8s var(--ease-smooth) 0.25s;
}

.hero-content.mounted .hero-badge {
  opacity: 1;
  transform: translateY(0);
}

.badge-wave {
  width: 24px;
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
  max-width: 500px;
  margin: 0 auto 48px;
  line-height: 1.8;
  opacity: 0;
  transform: translateY(30px);
  transition: all 1s var(--ease-smooth) 0.3s;
}

.hero-content.mounted .hero-subtitle {
  opacity: 1;
  transform: translateY(0);
}

.hero-subtitle .highlight {
  color: rgb(var(--aqua-glow));
}

.hero-actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  margin-bottom: 64px;
  opacity: 0;
  transform: translateY(30px);
  transition: all 1s var(--ease-smooth) 0.4s;
}

.hero-content.mounted .hero-actions {
  opacity: 1;
  transform: translateY(0);
}

@media (min-width: 480px) {
  .hero-actions {
    flex-direction: row;
    justify-content: center;
  }
}

.cta-primary {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 32px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  color: white;
  font-size: 16px;
  font-weight: 600;
  text-decoration: none;
  box-shadow: 0 4px 20px rgba(14, 165, 233, 0.4);
  transition: all 0.3s var(--ease-smooth);
}

.cta-primary:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 30px rgba(14, 165, 233, 0.5);
}

.cta-primary .arrow-icon {
  width: 18px;
  height: 18px;
  transition: transform 0.2s;
}

.cta-primary:hover .arrow-icon {
  transform: translateX(4px);
}

.cta-secondary {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 28px;
  border-radius: var(--radius-lg);
  background: rgba(56, 189, 248, 0.08);
  border: 1px solid rgba(56, 189, 248, 0.15);
  color: white;
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s var(--ease-smooth);
}

.cta-secondary:hover {
  background: rgba(56, 189, 248, 0.12);
  border-color: rgba(56, 189, 248, 0.25);
}

/* 统计 */
.stats-row {
  display: flex;
  justify-content: center;
  gap: 32px;
  flex-wrap: wrap;
  opacity: 0;
  transform: translateY(30px);
  transition: all 1s var(--ease-smooth) 0.5s;
}

.hero-content.mounted .stats-row {
  opacity: 1;
  transform: translateY(0);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 16px 24px;
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-lg);
  min-width: 100px;
}

.stat-icon {
  color: rgb(var(--aqua-glow));
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: white;
}

.stat-suffix {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin-left: 2px;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

/* 滚动指示器 */
.scroll-indicator {
  position: absolute;
  bottom: 40px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: rgba(255, 255, 255, 0.4);
  font-size: 12px;
  opacity: 0;
  transition: opacity 1s var(--ease-smooth) 1s;
}

.scroll-indicator.mounted {
  opacity: 1;
}

.scroll-mouse {
  width: 24px;
  height: 38px;
  border: 2px solid rgba(56, 189, 248, 0.3);
  border-radius: 12px;
  position: relative;
}

.scroll-wheel {
  position: absolute;
  top: 8px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 8px;
  background: rgba(56, 189, 248, 0.5);
  border-radius: 2px;
  animation: scroll-wheel 2s ease-in-out infinite;
}

@keyframes scroll-wheel {
  0%, 100% { transform: translateX(-50%) translateY(0); opacity: 1; }
  50% { transform: translateX(-50%) translateY(10px); opacity: 0.3; }
}

/* ===== Section 通用样式 ===== */
.section-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.section-header {
  text-align: center;
  margin-bottom: 64px;
}

.section-eyebrow {
  display: inline-block;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 3px;
  color: rgb(var(--aqua-glow));
  margin-bottom: 16px;
}

.section-title {
  font-size: clamp(28px, 5vw, 40px);
  font-weight: 700;
  color: white;
  margin: 0 0 12px;
}

.section-desc {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.5);
}

/* ===== Features Section ===== */
.features-section {
  padding: 120px 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
}

.feature-item {
  padding: 32px;
  background: rgba(8, 15, 30, 0.4);
  border: 1px solid rgba(56, 189, 248, 0.08);
  border-radius: var(--radius-xl);
  text-align: center;
  transition: all 0.4s var(--ease-smooth);
}

.feature-item:hover {
  background: rgba(8, 15, 30, 0.6);
  border-color: rgba(56, 189, 248, 0.15);
  transform: translateY(-6px);
}

.feature-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.15), rgba(56, 189, 248, 0.05));
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--feature-color);
  transition: transform 0.3s var(--ease-smooth);
}

.feature-item:hover .feature-icon {
  transform: scale(1.1);
}

.feature-title {
  font-size: 18px;
  font-weight: 600;
  color: white;
  margin: 0 0 10px;
}

.feature-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  line-height: 1.6;
  margin: 0;
}

/* ===== Heatmap Section ===== */
.heatmap-section {
  padding: 120px 0;
}

.heatmap-wrapper {
  background: rgba(8, 15, 30, 0.4);
  border: 1px solid rgba(56, 189, 248, 0.08);
  border-radius: var(--radius-2xl);
  padding: 32px;
}

.heatmap-stats {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.hm-stat {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.hm-value {
  font-size: 24px;
  font-weight: 700;
  color: rgb(var(--ocean-surface));
}

.hm-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
}

.hm-divider {
  width: 1px;
  height: 24px;
  background: rgba(56, 189, 248, 0.2);
}

.live-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: rgba(52, 211, 153, 0.1);
  border: 1px solid rgba(52, 211, 153, 0.15);
  border-radius: var(--radius-full);
  font-size: 12px;
  color: rgb(var(--seaweed-green));
  margin-left: auto;
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
  50% { opacity: 0.5; transform: scale(1.2); }
}

.heatmap-loading {
  padding: 40px 0;
}

/* ===== Steps Section ===== */
.steps-section {
  padding: 120px 0;
}

.steps-grid {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  flex-wrap: wrap;
}

.step-item {
  text-align: center;
  padding: 0 24px;
  flex: 0 0 auto;
}

.step-num {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: white;
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
  color: white;
  margin: 0 0 8px;
}

.step-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  max-width: 180px;
  margin: 0 auto;
}

.step-connector {
  width: 80px;
  height: 20px;
  flex-shrink: 0;
}

.step-connector svg {
  width: 100%;
  height: 100%;
}

@media (max-width: 768px) {
  .step-connector {
    display: none;
  }

  .steps-grid {
    flex-direction: column;
    gap: 40px;
  }
}

/* ===== CTA Section ===== */
.cta-section {
  padding: 120px 24px;
}

.cta-container {
  max-width: 700px;
  margin: 0 auto;
  text-align: center;
  padding: 64px 40px;
  background: rgba(8, 15, 30, 0.5);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-2xl);
  position: relative;
  overflow: hidden;
}

.cta-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgb(var(--ocean-surface)), transparent);
}

.cta-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 28px;
  border-radius: var(--radius-xl);
  background: linear-gradient(180deg, #0c4a6e, #0369a1);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 40px rgba(14, 165, 233, 0.3);
  overflow: hidden;
  position: relative;
}

.cta-wave-svg {
  width: 100%;
  height: 100%;
  position: absolute;
  inset: 0;
}

/* CTA海浪涌动动画 */
.cta-wave-svg .cta-wave-1 {
  animation: cta-surge-1 3s ease-in-out infinite;
}

.cta-wave-svg .cta-wave-2 {
  animation: cta-surge-2 3.5s ease-in-out infinite;
}

.cta-wave-svg .cta-wave-3 {
  animation: cta-surge-3 4s ease-in-out infinite;
}

.cta-wave-svg .cta-wave-line {
  animation: cta-line-surge 2.5s ease-in-out infinite;
}

@keyframes cta-surge-1 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}

@keyframes cta-surge-2 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

@keyframes cta-surge-3 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-2px); }
}

@keyframes cta-line-surge {
  0%, 100% {
    transform: translateY(0);
    d: path("M8 28 Q22 20, 36 28 T52 28");
  }
  50% {
    transform: translateY(-2px);
    d: path("M8 28 Q22 36, 36 28 T52 28");
  }
}

.cta-title {
  font-size: clamp(24px, 4vw, 32px);
  font-weight: 700;
  color: white;
  margin: 0 0 16px;
}

.cta-desc {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 36px;
}

.cta-buttons {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

@media (min-width: 480px) {
  .cta-buttons {
    flex-direction: row;
    justify-content: center;
  }
}

.cta-btn-primary {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 32px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  color: white;
  font-size: 16px;
  font-weight: 600;
  text-decoration: none;
  box-shadow: 0 4px 20px rgba(14, 165, 233, 0.4);
  transition: all 0.3s var(--ease-smooth);
}

.cta-btn-primary:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 30px rgba(14, 165, 233, 0.5);
}

.cta-btn-primary .arrow-icon {
  width: 18px;
  height: 18px;
  transition: transform 0.2s;
}

.cta-btn-primary:hover .arrow-icon {
  transform: translateX(4px);
}

.cta-btn-secondary {
  padding: 16px 28px;
  border-radius: var(--radius-lg);
  background: rgba(56, 189, 248, 0.08);
  border: 1px solid rgba(56, 189, 248, 0.15);
  color: white;
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s var(--ease-smooth);
}

.cta-btn-secondary:hover {
  background: rgba(56, 189, 248, 0.12);
  border-color: rgba(56, 189, 248, 0.25);
}
</style>
