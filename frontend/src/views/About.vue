<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { Timer, Check, Warning, Star, ArrowRight } from '@element-plus/icons-vue'

const mounted = ref(false)

onMounted(() => {
  setTimeout(() => {
    mounted.value = true
  }, 100)
})

const benefits = [
  {
    icon: '💪',
    title: '增强盆底肌力量',
    description: '规律训练可显著提升盆底肌群的收缩力和耐力'
  },
  {
    icon: '🩺',
    title: '改善尿失禁',
    description: '有效预防和改善压力性尿失禁等问题'
  },
  {
    icon: '❤️',
    title: '提升生活质量',
    description: '改善产后恢复，提升整体健康水平'
  },
  {
    icon: '🧘',
    title: '核心稳定性',
    description: '盆底肌是核心肌群的重要组成部分'
  },
  {
    icon: '👴',
    title: '预防老年问题',
    description: '降低老年后盆底功能障碍的风险'
  },
  {
    icon: '🔄',
    title: '促进血液循环',
    description: '改善盆腔区域的血液循环和新陈代谢'
  }
]

const steps = [
  {
    phase: '收缩',
    duration: '5秒',
    description: '想象正在憋尿，收紧盆底肌肉，感受向上提拉的力量',
    tips: '保持呼吸平稳，不要屏气'
  },
  {
    phase: '保持',
    duration: '5秒',
    description: '维持收缩状态，保持肌肉紧张',
    tips: '专注于盆底肌，避免腹部、臀部用力'
  },
  {
    phase: '放松',
    duration: '5秒',
    description: '缓慢放松肌肉，让盆底完全松弛',
    tips: '放松同样重要，不要急于开始下一次'
  }
]

const tips = [
  '训练前先排空膀胱',
  '保持正常呼吸，不要屏气',
  '只收缩盆底肌，避免腹部、臀部、大腿用力',
  '循序渐进，从每天1-2组开始',
  '坚持规律训练，效果通常在4-6周后显现',
  '可以在任何姿势下进行：坐、站、躺'
]

const warnings = [
  '如有泌尿系统感染，请先治疗后再训练',
  '孕期训练请先咨询医生',
  '训练时如感到疼痛，请立即停止并就医',
  '不要在排尿过程中练习凯格尔运动'
]
</script>

<template>
  <MainLayout>
    <div class="about-page">
      <!-- Hero Section -->
      <section class="hero-section">
        <div class="hero-content" :class="{ mounted }">
          <h1 class="hero-title">关于凯格尔运动</h1>
          <p class="hero-subtitle">
            凯格尔运动（Kegel Exercise）是一种简单有效的盆底肌训练方法，
            <br />
            由美国妇科医生阿诺德·凯格尔于1948年发明。
          </p>
        </div>
      </section>

      <!-- What is Kegel Section -->
      <section class="content-section">
        <div class="section-container">
          <div class="section-header">
            <div class="section-eyebrow">
              <span class="eyebrow-line"></span>
              <span class="eyebrow-text">基础知识</span>
              <span class="eyebrow-line"></span>
            </div>
            <h2 class="section-title">什么是凯格尔运动？</h2>
          </div>

          <div class="intro-card">
            <div class="intro-content">
              <p>
                凯格尔运动是一种针对<strong>盆底肌群</strong>的锻炼方法。盆底肌是一组位于骨盆底部的肌肉，
                像吊床一样支撑着膀胱、子宫（女性）、前列腺（男性）和直肠等器官。
              </p>
              <p>
                这项运动通过反复收缩和放松盆底肌肉来增强其力量和耐力，
                <strong>无需任何器械</strong>，可以在任何时间、任何地点进行，
                是一种安全、有效且隐蔽的锻炼方式。
              </p>
              <p class="highlight-text">
                无论男女老少，都可以从凯格尔运动中获益。
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- Benefits Section -->
      <section class="content-section benefits-section">
        <div class="section-container">
          <div class="section-header">
            <div class="section-eyebrow">
              <span class="eyebrow-line"></span>
              <span class="eyebrow-text">健康益处</span>
              <span class="eyebrow-line"></span>
            </div>
            <h2 class="section-title">凯格尔运动的好处</h2>
            <p class="section-desc">规律的盆底肌训练带来多方面的健康益处</p>
          </div>

          <div class="benefits-grid">
            <div v-for="(benefit, index) in benefits" :key="index" class="benefit-card">
              <div class="benefit-icon">{{ benefit.icon }}</div>
              <h3 class="benefit-title">{{ benefit.title }}</h3>
              <p class="benefit-desc">{{ benefit.description }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- How to Section -->
      <section class="content-section rhythm-section">
        <div class="section-container">
          <div class="section-header">
            <div class="section-eyebrow">
              <span class="eyebrow-line"></span>
              <span class="eyebrow-text">训练节奏</span>
              <span class="eyebrow-line"></span>
            </div>
            <h2 class="section-title">正确的训练方法</h2>
            <p class="section-desc">遵循「收缩-保持-放松」的潮汐节奏</p>
          </div>

          <div class="rhythm-flow">
            <div v-for="(step, index) in steps" :key="index" class="rhythm-step">
              <div class="step-header">
                <div class="step-number">{{ index + 1 }}</div>
                <div class="step-phase">{{ step.phase }}</div>
                <div class="step-duration">
                  <el-icon><Timer /></el-icon>
                  {{ step.duration }}
                </div>
              </div>
              <p class="step-description">{{ step.description }}</p>
              <div class="step-tip">
                <el-icon><Star /></el-icon>
                <span>{{ step.tips }}</span>
              </div>
              <div v-if="index < steps.length - 1" class="step-connector">
                <el-icon><ArrowRight /></el-icon>
              </div>
            </div>
          </div>

          <div class="rhythm-summary">
            <div class="summary-item">
              <span class="summary-label">每组</span>
              <span class="summary-value">10-15 次</span>
            </div>
            <div class="summary-divider"></div>
            <div class="summary-item">
              <span class="summary-label">每天</span>
              <span class="summary-value">3 组</span>
            </div>
            <div class="summary-divider"></div>
            <div class="summary-item">
              <span class="summary-label">见效周期</span>
              <span class="summary-value">4-6 周</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Tips Section -->
      <section class="content-section tips-section">
        <div class="section-container">
          <div class="tips-grid">
            <!-- Tips Card -->
            <div class="tips-card success-card">
              <div class="tips-header">
                <el-icon class="tips-icon" :size="24"><Check /></el-icon>
                <h3 class="tips-title">训练小贴士</h3>
              </div>
              <ul class="tips-list">
                <li v-for="(tip, index) in tips" :key="index">
                  <span class="tip-bullet">✓</span>
                  <span>{{ tip }}</span>
                </li>
              </ul>
            </div>

            <!-- Warnings Card -->
            <div class="tips-card warning-card">
              <div class="tips-header">
                <el-icon class="tips-icon" :size="24"><Warning /></el-icon>
                <h3 class="tips-title">注意事项</h3>
              </div>
              <ul class="tips-list">
                <li v-for="(warning, index) in warnings" :key="index">
                  <span class="tip-bullet warning">!</span>
                  <span>{{ warning }}</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      <!-- CTA Section -->
      <section class="cta-section">
        <div class="cta-container">
          <h2 class="cta-title">准备好开始训练了吗？</h2>
          <p class="cta-desc">使用 TidalCore 的潮汐计时器，跟随节奏轻松完成每日训练</p>
          <div class="cta-buttons">
            <RouterLink to="/train" class="cta-btn-primary">
              <el-icon><Timer /></el-icon>
              <span>开始训练</span>
            </RouterLink>
            <RouterLink to="/" class="cta-btn-secondary">返回首页</RouterLink>
          </div>
        </div>
      </section>
    </div>
  </MainLayout>
</template>

<style scoped>
.about-page {
  display: flex;
  flex-direction: column;
}

/* ===== Hero Section ===== */
.hero-section {
  min-height: 50vh;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 30px 24px 40px;
}

.hero-content {
  max-width: 800px;
  opacity: 0;
  transform: translateY(30px);
  transition: all 0.8s var(--ease-smooth);
}

.hero-content.mounted {
  opacity: 1;
  transform: translateY(0);
}

.hero-title {
  font-size: clamp(36px, 8vw, 56px);
  font-weight: 700;
  margin-bottom: 20px;
  background: linear-gradient(135deg, #38bdf8 0%, #22d3ee 50%, #0ea5e9 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.8;
}

/* ===== Section Common ===== */
.content-section {
  padding: 66px 0;
}

.section-container {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0 24px;
}

.section-header {
  text-align: center;
  margin-bottom: 48px;
}

.section-eyebrow {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 16px;
}

.section-eyebrow .eyebrow-line {
  width: 50px;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgb(var(--ocean-surface)));
}

.section-eyebrow .eyebrow-line:last-child {
  background: linear-gradient(90deg, rgb(var(--ocean-surface)), transparent);
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

.section-title {
  font-size: clamp(24px, 5vw, 36px);
  font-weight: 700;
  color: white;
  margin: 0 0 12px;
}

.section-desc {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.6);
}

/* ===== Intro Card ===== */
.intro-card {
  background: rgba(8, 15, 30, 0.5);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-2xl);
  padding: 40px;
}

.intro-content p {
  font-size: 16px;
  line-height: 1.9;
  color: rgba(255, 255, 255, 0.75);
  margin: 0 0 20px;
}

.intro-content p:last-child {
  margin-bottom: 0;
}

.intro-content strong {
  color: rgb(var(--ocean-surface));
}

.intro-content .highlight-text {
  font-size: 18px;
  font-weight: 500;
  color: rgb(var(--aqua-glow));
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid rgba(56, 189, 248, 0.1);
  margin-top: 24px;
}

/* ===== Benefits Grid ===== */
.benefits-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.benefit-card {
  padding: 28px;
  background: rgba(8, 15, 30, 0.4);
  border: 1px solid rgba(56, 189, 248, 0.08);
  border-radius: var(--radius-xl);
  transition: all 0.3s var(--ease-smooth);
}

.benefit-card:hover {
  background: rgba(8, 15, 30, 0.6);
  border-color: rgba(56, 189, 248, 0.15);
  transform: translateY(-4px);
}

.benefit-icon {
  font-size: 32px;
  margin-bottom: 16px;
}

.benefit-title {
  font-size: 18px;
  font-weight: 600;
  color: white;
  margin: 0 0 10px;
}

.benefit-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  line-height: 1.6;
  margin: 0;
}

/* ===== Rhythm Section ===== */
.rhythm-flow {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 48px;
}

@media (min-width: 900px) {
  .rhythm-flow {
    flex-direction: row;
    align-items: stretch;
    gap: 0;
  }
}

.rhythm-step {
  flex: 1;
  padding: 28px;
  background: rgba(8, 15, 30, 0.5);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-xl);
  position: relative;
  display: flex;
  flex-direction: column;
}

@media (min-width: 900px) {
  .rhythm-step {
    margin: 0 -1px;
  }

  .rhythm-step:first-child {
    border-radius: var(--radius-xl) 0 0 var(--radius-xl);
  }

  .rhythm-step:last-child {
    border-radius: 0 var(--radius-xl) var(--radius-xl) 0;
  }

  .rhythm-step:not(:first-child):not(:last-child) {
    border-radius: 0;
  }
}

.step-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.step-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  color: white;
}

.step-phase {
  font-size: 20px;
  font-weight: 600;
  color: white;
}

.step-duration {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
  padding: 6px 12px;
  background: rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-full);
  font-size: 13px;
  font-weight: 500;
  color: rgb(var(--ocean-surface));
}

.step-description {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.7;
  margin: 0 0 16px;
  flex: 1;
}

.step-tip {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 12px;
  background: rgba(251, 191, 36, 0.08);
  border: 1px solid rgba(251, 191, 36, 0.15);
  border-radius: var(--radius-md);
  font-size: 13px;
  color: rgba(251, 191, 36, 0.9);
  margin-top: auto;
}

.step-tip .el-icon {
  flex-shrink: 0;
  margin-top: 2px;
}

.step-connector {
  display: none;
}

@media (min-width: 900px) {
  .step-connector {
    display: flex;
    position: absolute;
    right: -12px;
    top: 50%;
    transform: translateY(-50%);
    z-index: 1;
    width: 24px;
    height: 24px;
    align-items: center;
    justify-content: center;
    background: rgba(8, 15, 30, 0.9);
    border: 1px solid rgba(56, 189, 248, 0.2);
    border-radius: 50%;
    color: rgb(var(--ocean-surface));
  }
}

.rhythm-summary {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 32px;
  padding: 24px;
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.1) 0%, rgba(34, 211, 238, 0.05) 100%);
  border: 1px solid rgba(56, 189, 248, 0.2);
  border-radius: var(--radius-xl);
  flex-wrap: wrap;
}

.summary-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.summary-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.summary-value {
  font-size: 20px;
  font-weight: 700;
  color: rgb(var(--ocean-surface));
}

.summary-divider {
  width: 1px;
  height: 40px;
  background: rgba(56, 189, 248, 0.2);
}

/* ===== Tips Section ===== */
.tips-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 24px;
}

.tips-card {
  padding: 28px;
  border-radius: var(--radius-xl);
}

.success-card {
  background: rgba(52, 211, 153, 0.05);
  border: 1px solid rgba(52, 211, 153, 0.15);
}

.warning-card {
  background: rgba(251, 191, 36, 0.05);
  border: 1px solid rgba(251, 191, 36, 0.15);
}

.tips-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.success-card .tips-icon {
  color: rgb(var(--seaweed-green));
}

.warning-card .tips-icon {
  color: #fbbf24;
}

.tips-title {
  font-size: 18px;
  font-weight: 600;
  color: white;
  margin: 0;
}

.tips-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.tips-list li {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.75);
  line-height: 1.6;
}

.tip-bullet {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: rgba(52, 211, 153, 0.15);
  color: rgb(var(--seaweed-green));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 600;
}

.tip-bullet.warning {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

/* ===== CTA Section ===== */
.cta-section {
  padding: 80px 24px;
}

.cta-container {
  max-width: 600px;
  margin: 0 auto;
  text-align: center;
  padding: 48px 40px;
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

.cta-title {
  font-size: clamp(22px, 4vw, 28px);
  font-weight: 700;
  color: white;
  margin: 0 0 12px;
}

.cta-desc {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0 0 32px;
}

.cta-buttons {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
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
  gap: 8px;
  padding: 14px 28px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  color: white;
  font-size: 15px;
  font-weight: 600;
  text-decoration: none;
  box-shadow: 0 4px 20px rgba(14, 165, 233, 0.4);
  transition: all 0.3s var(--ease-smooth);
}

.cta-btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(14, 165, 233, 0.5);
}

.cta-btn-secondary {
  padding: 14px 24px;
  border-radius: var(--radius-lg);
  background: rgba(56, 189, 248, 0.08);
  border: 1px solid rgba(56, 189, 248, 0.15);
  color: white;
  font-size: 15px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s var(--ease-smooth);
}

.cta-btn-secondary:hover {
  background: rgba(56, 189, 248, 0.12);
  border-color: rgba(56, 189, 248, 0.25);
}
</style>
