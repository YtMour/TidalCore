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
    title: '科学计时',
    description: '自定义收缩-保持-放松循环，专业训练节奏引导',
    color: '#a78bfa',
    bg: 'linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(168, 85, 247, 0.1))'
  },
  {
    icon: Pointer,
    title: '连续打卡',
    description: '记录你的坚持轨迹，见证每一天的进步与成长',
    color: '#fb923c',
    bg: 'linear-gradient(135deg, rgba(249, 115, 22, 0.2), rgba(251, 146, 60, 0.1))'
  },
  {
    icon: Trophy,
    title: '毅力排行',
    description: '匿名排行榜，与全站用户一起坚持，互相激励',
    color: '#fbbf24',
    bg: 'linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(251, 191, 36, 0.1))'
  },
  {
    icon: Lock,
    title: '隐私保护',
    description: '仅需用户名密码，不收集任何敏感个人信息',
    color: '#34d399',
    bg: 'linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(52, 211, 153, 0.1))'
  },
  {
    icon: Lightning,
    title: '轻量高效',
    description: '每天只需几分钟，随时随地开始训练',
    color: '#22d3ee',
    bg: 'linear-gradient(135deg, rgba(6, 182, 212, 0.2), rgba(34, 211, 238, 0.1))'
  },
  {
    icon: StarFilled,
    title: '健康习惯',
    description: '科学的训练方法，帮助建立长期健康习惯',
    color: '#f472b6',
    bg: 'linear-gradient(135deg, rgba(236, 72, 153, 0.2), rgba(244, 114, 182, 0.1))'
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
      <!-- Hero Section -->
      <section class="hero-section">
        <!-- Enhanced Decorative elements -->
        <div class="hero-bg">
          <div class="hero-blob hero-blob-1"></div>
          <div class="hero-blob hero-blob-2"></div>
          <div class="hero-blob hero-blob-3"></div>
          <div class="hero-ring hero-ring-1"></div>
          <div class="hero-ring hero-ring-2"></div>
        </div>

        <div class="hero-content" :class="{ mounted }">
          <!-- Logo with enhanced animation -->
          <div class="hero-logo" :class="{ mounted }">
            <div class="logo-wrapper">
              <span class="logo-emoji">🌊</span>
              <div class="logo-glow"></div>
            </div>
          </div>

          <!-- Title with enhanced gradient -->
          <h1 class="hero-title" :class="{ mounted }">
            <span class="gradient-text">TidalCore</span>
          </h1>

          <!-- Tagline badge -->
          <div class="hero-badge" :class="{ mounted }">
            <span class="badge-dot"></span>
            <span>开源 · 免费 · 隐私优先</span>
          </div>

          <!-- Subtitle -->
          <p class="hero-subtitle" :class="{ mounted }">
            专业的盆底肌训练平台，帮助你建立健康的训练习惯。
            <br class="hidden-mobile" />
            <span class="highlight">每天几分钟，坚持就是胜利。</span>
          </p>

          <!-- CTA Buttons -->
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

          <!-- Stats Row -->
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
          <div class="section-eyebrow">
            <span class="eyebrow-line"></span>
            <span class="eyebrow-text">核心功能</span>
            <span class="eyebrow-line"></span>
          </div>
          <h2 class="section-title">为什么选择 TidalCore</h2>
          <p class="section-subtitle">简单、专注、有效的训练体验，帮助你养成健康习惯</p>
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
          <div class="section-eyebrow section-eyebrow-pink">
            <span class="eyebrow-line"></span>
            <span class="eyebrow-text">快速上手</span>
            <span class="eyebrow-line"></span>
          </div>
          <h2 class="section-title">三步开始训练</h2>
          <p class="section-subtitle">无需复杂设置，立即开始你的健康之旅</p>
        </div>

        <el-row :gutter="32">
          <el-col :xs="24" :md="8" v-for="(step, index) in [
            { num: 1, title: '注册账号', desc: '仅需用户名和密码，30秒完成注册', color: ['#8b5cf6', '#a855f7'] },
            { num: 2, title: '开始训练', desc: '跟随计时器节奏，完成每日训练', color: ['#ec4899', '#f43f5e'] },
            { num: 3, title: '坚持打卡', desc: '记录进度，冲击排行榜', color: ['#f59e0b', '#f97316'] }
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

          <div class="cta-content">
            <div class="cta-icon">
              <el-icon :size="32" color="#fff"><StarFilled /></el-icon>
            </div>
            <h3 class="cta-title">准备好开始了吗？</h3>
            <p class="cta-text">加入数千名用户，开始你的健康训练之旅。完全免费，永久开源。</p>
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

/* Hero Section */
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

.hero-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
}

.hero-blob-1 {
  top: 0;
  left: 25%;
  width: 384px;
  height: 384px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.3), rgba(168, 85, 247, 0.2));
  animation: float 6s ease-in-out infinite;
}

.hero-blob-2 {
  bottom: 0;
  right: 25%;
  width: 320px;
  height: 320px;
  background: linear-gradient(135deg, rgba(236, 72, 153, 0.25), rgba(244, 63, 94, 0.15));
  animation: float 8s ease-in-out infinite;
  animation-delay: -2s;
}

.hero-blob-3 {
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 600px;
  height: 600px;
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.1), rgba(59, 130, 246, 0.05));
}

.hero-ring {
  position: absolute;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.hero-ring-1 {
  top: 80px;
  right: 80px;
  width: 128px;
  height: 128px;
  animation: pulse 4s ease-in-out infinite;
}

.hero-ring-2 {
  bottom: 128px;
  left: 64px;
  width: 96px;
  height: 96px;
  border-color: rgba(139, 92, 246, 0.1);
  animation: float 5s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.05); }
}

.hero-content {
  position: relative;
  z-index: 10;
}

.hero-logo {
  margin-bottom: 32px;
  opacity: 0;
  transform: translateY(32px) scale(0.9);
  transition: all 1s ease;
}

.hero-logo.mounted {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.logo-wrapper {
  position: relative;
  display: inline-block;
}

.logo-emoji {
  font-size: 128px;
  display: inline-block;
  animation: float 3s ease-in-out infinite;
  filter: drop-shadow(0 20px 40px rgba(0, 0, 0, 0.3));
}

@media (max-width: 768px) {
  .logo-emoji {
    font-size: 96px;
  }
}

.logo-glow {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(34, 211, 238, 0.3), rgba(59, 130, 246, 0.2));
  border-radius: 50%;
  filter: blur(40px);
  transform: scale(1.5);
  z-index: -1;
}

.hero-title {
  font-size: 64px;
  font-weight: 800;
  margin-bottom: 24px;
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s ease 0.15s;
}

.hero-content.mounted .hero-title {
  opacity: 1;
  transform: translateY(0);
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 48px;
  }
}

.gradient-text {
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  margin-bottom: 24px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  opacity: 0;
  transform: translateY(32px);
  transition: all 0.7s ease 0.2s;
}

.hero-content.mounted .hero-badge {
  opacity: 1;
  transform: translateY(0);
}

.badge-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #34d399;
  animation: pulse 2s ease-in-out infinite;
}

.hero-subtitle {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.6);
  max-width: 600px;
  margin: 0 auto 48px;
  line-height: 1.8;
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s ease 0.3s;
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
  color: rgba(255, 255, 255, 0.8);
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
  transition: all 1s ease 0.4s;
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
  background: linear-gradient(135deg, #8b5cf6, #ec4899) !important;
  border: none !important;
  padding: 12px 32px !important;
  font-size: 16px !important;
  height: auto !important;
  border-radius: 8px !important;
}

.cta-primary .arrow-icon {
  margin-left: 8px;
  transition: transform 0.2s ease;
}

.cta-primary:hover .arrow-icon {
  transform: translateX(4px);
}

.cta-secondary {
  background: rgba(255, 255, 255, 0.04) !important;
  border: 1px solid rgba(255, 255, 255, 0.08) !important;
  padding: 12px 32px !important;
  font-size: 16px !important;
  height: auto !important;
  border-radius: 8px !important;
  color: #fff !important;
}

.cta-secondary:hover {
  background: rgba(255, 255, 255, 0.08) !important;
  border-color: rgba(255, 255, 255, 0.15) !important;
}

.stats-row {
  max-width: 768px;
  margin: 0 auto !important;
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s ease 0.5s;
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
  padding: 14px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: all 0.2s ease;
  margin-bottom: 16px;
}

.stat-card:hover {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.1);
}

.stat-icon {
  color: #a78bfa;
  flex-shrink: 0;
}

.stat-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.stat-value {
  font-size: 22px;
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

/* Heatmap Section */
.heatmap-section {
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s ease 0.6s;
}

.heatmap-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.heatmap-card {
  background: rgba(30, 30, 46, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 10px !important;
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
  width: 256px;
  height: 256px;
  background: linear-gradient(225deg, rgba(139, 92, 246, 0.1), transparent);
  border-radius: 50%;
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
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(52, 211, 153, 0.1));
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
  padding: 8px 14px;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.12), rgba(16, 185, 129, 0.04));
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: 8px;
}

.hm-stats .stat-item {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.hm-stats .stat-value {
  font-size: 18px;
  font-weight: 700;
  color: #10b981;
}

.hm-stats .stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.hm-stats .stat-divider {
  width: 1px;
  height: 20px;
  background: rgba(255, 255, 255, 0.1);
}

.live-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #34d399;
  padding: 6px 12px;
  border-radius: 20px;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.15);
}

.live-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #34d399;
  animation: pulse 2s ease-in-out infinite;
}

.heatmap-loading {
  padding: 32px 0;
}

/* Features Section */
.features-section {
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s ease 0.7s;
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
  background: linear-gradient(90deg, transparent, #8b5cf6);
}

.section-eyebrow .eyebrow-line:last-child {
  background: linear-gradient(90deg, #8b5cf6, transparent);
}

.section-eyebrow .eyebrow-text {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 3px;
  background: linear-gradient(135deg, #a78bfa, #c4b5fd);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.section-eyebrow-pink .eyebrow-line {
  background: linear-gradient(90deg, transparent, #ec4899);
}

.section-eyebrow-pink .eyebrow-line:last-child {
  background: linear-gradient(90deg, #ec4899, transparent);
}

.section-eyebrow-pink .eyebrow-text {
  background: linear-gradient(135deg, #f472b6, #fda4af);
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
  background: rgba(30, 30, 46, 0.6) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 10px !important;
  margin-bottom: 24px;
  transition: all 0.3s ease !important;
}

.feature-card:hover {
  border-color: rgba(255, 255, 255, 0.1) !important;
  transform: translateY(-4px);
}

.feature-card :deep(.el-card__body) {
  padding: 20px;
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  gap: 16px;
}

.feature-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform 0.3s ease;
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
  margin: 0 0 4px;
  transition: color 0.2s ease;
}

.feature-card:hover .feature-title {
  color: #ddd6fe;
}

.feature-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  line-height: 1.5;
  margin: 0;
}

/* How To Section */
.howto-section {
  opacity: 0;
  transform: translateY(32px);
  transition: all 1s ease 0.8s;
}

.howto-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.howto-section .section-tag {
  background: rgba(236, 72, 153, 0.1) !important;
  border-color: rgba(236, 72, 153, 0.2) !important;
}

.step-item {
  position: relative;
  text-align: center;
  margin-bottom: 32px;
}

.step-num {
  width: 56px;
  height: 56px;
  margin: 0 auto 20px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
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
    height: 1px;
    background: linear-gradient(to right, rgba(139, 92, 246, 0.5), transparent);
  }
}

/* CTA Section */
.cta-section {
  padding-bottom: 64px;
}

.cta-card {
  max-width: 768px;
  margin: 0 auto;
  background: rgba(30, 30, 46, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 10px !important;
  overflow: hidden;
  position: relative;
}

.cta-card :deep(.el-card__body) {
  padding: 48px;
}

.cta-decoration {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  pointer-events: none;
}

.cta-decoration-1 {
  top: -96px;
  left: -96px;
  width: 192px;
  height: 192px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), transparent);
}

.cta-decoration-2 {
  bottom: -96px;
  right: -96px;
  width: 192px;
  height: 192px;
  background: linear-gradient(135deg, rgba(236, 72, 153, 0.2), transparent);
}

.cta-content {
  position: relative;
  text-align: center;
}

.cta-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 24px;
  border-radius: 10px;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(139, 92, 246, 0.3);
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
  max-width: 448px;
  margin: 0 auto 32px;
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
