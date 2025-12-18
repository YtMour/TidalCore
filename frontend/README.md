# TidalCore Frontend

TidalCore 前端应用 - 专业的盆底肌训练平台，采用全屏沉浸式海洋主题设计。

## 技术栈

- **框架**: Vue 3.5 + TypeScript + Composition API
- **构建工具**: Vite 7
- **状态管理**: Pinia 3
- **路由**: Vue Router 4
- **UI 组件库**: Element Plus 2.x (深色海洋主题)
- **样式**: Tailwind CSS 4 + 自定义海洋设计系统
- **图标**: @element-plus/icons-vue
- **HTTP 客户端**: Axios
- **动画**: Canvas Confetti, Three.js, GSAP
- **工具库**: VueUse

## 设计理念

TidalCore 采用**全屏沉浸式海洋主题**设计，整个页面如同置身深海之中：

- **全屏海洋背景**: 多层动态海浪、涟漪、气泡效果
- **海浪涌动动画**: 波浪上下起伏，模拟真实海浪涌动
- **玻璃态组件**: 透明背景，内容悬浮在海洋之上
- **海洋色系**: 深海蓝到浅海青的渐变色调

## 主要特性

### 全屏海洋背景系统

```css
/* 多层背景架构 */
.ocean-bg          /* 深海渐变基底 */
.ocean-gradient    /* 光线穿透效果 */
.animated-bg       /* 动态水体效果 */
.particles-bg      /* 海洋光斑粒子 */
.wave-layer        /* 三层海浪动画 */
.ripple-layer      /* 涟漪效果 */
.bubbles-container /* 上升气泡 */
```

### 海洋色系

```css
:root {
  /* 海洋主色 */
  --ocean-surface: 56 189 248;    /* #38bdf8 - 海面蓝 */
  --ocean-shallow: 14 165 233;    /* #0ea5e9 - 浅海蓝 */
  --ocean-mid: 2 132 199;         /* #0284c7 - 中海蓝 */
  --ocean-deep: 3 105 161;        /* #0369a1 - 深海蓝 */
  --ocean-abyss: 12 74 110;       /* #0c4a6e - 海渊蓝 */

  /* 海洋强调色 */
  --aqua-glow: 34 211 238;        /* #22d3ee 水光 */
  --seaweed-green: 52 211 153;    /* #34d399 海藻绿 */
  --sunset-amber: 251 191 36;     /* #fbbf24 落日金 */
  --coral-pink: 244 114 182;      /* #f472b6 珊瑚粉 */
}
```

### 海浪涌动动画

核心动画效果 - 波浪上下起伏而非左右移动：

```css
/* 海浪涌动 - 不同层次不同速度 */
@keyframes wave-surge-1 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

@keyframes wave-surge-2 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}

@keyframes wave-surge-3 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
```

### SVG 海浪组件

Logo、CTA图标、奖杯等均使用海浪填充效果：

```html
<!-- 海浪填充SVG示例 -->
<svg viewBox="0 0 60 60">
  <defs>
    <linearGradient id="waveGrad" x1="0%" y1="100%" x2="0%" y2="0%">
      <stop offset="0%" stop-color="#0ea5e9" />
      <stop offset="100%" stop-color="#38bdf8" />
    </linearGradient>
    <clipPath id="clip">
      <circle cx="30" cy="30" r="26" />
    </clipPath>
  </defs>
  <g clip-path="url(#clip)">
    <path class="wave-1" d="M0 35 Q15 30, 30 35 T60 35 L60 60 L0 60 Z" fill="url(#waveGrad)"/>
    <path class="wave-2" d="M0 40 Q15 35, 30 40 T60 40 L60 60 L0 60 Z" fill="#22d3ee" opacity="0.6"/>
    <path class="wave-3" d="M0 45 Q15 40, 30 45 T60 45 L60 60 L0 60 Z" fill="#06b6d4" opacity="0.4"/>
  </g>
</svg>
```

## Element Plus 深色海洋主题

项目对 Element Plus 进行了深度定制，与海洋主题融为一体：

```css
.dark {
  --el-bg-color: transparent;
  --el-bg-color-page: transparent;
  --el-fill-color: rgba(56, 189, 248, 0.05);
  --el-border-color: rgba(56, 189, 248, 0.1);
  --el-color-primary: rgb(var(--ocean-surface));
}

/* 所有组件背景透明，显示海洋背景 */
.el-card {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
}
```

## 项目结构

```
src/
├── api/                    # API 接口层
│   ├── auth.ts            # 认证相关 API
│   ├── checkin.ts         # 打卡相关 API
│   └── request.ts         # Axios 实例配置
├── components/             # 组件
│   ├── Heatmap.vue        # 热力图组件
│   ├── LeaderboardTable.vue # 排行榜表格组件
│   ├── TidalBackground.vue  # Three.js 海洋背景
│   └── Timer.vue          # 训练计时器组件
├── layouts/                # 布局组件
│   └── MainLayout.vue     # 主布局 (全屏海洋背景、透明导航)
├── router/                 # 路由配置
│   └── index.ts
├── store/                  # 状态管理
│   ├── training.ts        # 训练状态
│   └── user.ts            # 用户状态
├── views/                  # 页面视图
│   ├── Dashboard.vue      # 个人中心 (全屏)
│   ├── Home.vue           # 首页 (全屏英雄区)
│   ├── Leaderboard.vue    # 排行榜 (奖杯海浪效果)
│   ├── Login.vue          # 登录 (全屏)
│   ├── Register.vue       # 注册 (全屏)
│   └── Train.vue          # 训练页面 (Three.js 背景)
├── App.vue                 # 根组件
├── main.ts                 # 入口文件
└── style.css              # 全局样式 (海洋设计系统)
```

## 快速开始

### 环境要求

- Node.js >= 18.0.0
- npm >= 9.0.0 或 pnpm >= 8.0.0

### 安装依赖

```bash
npm install
# 或
pnpm install
```

### 开发模式

```bash
npm run dev
```

应用将在 `http://localhost:5173` 启动。

### 构建生产版本

```bash
npm run build
```

构建产物将输出到 `dist/` 目录。

## 功能特性

### 核心功能

- **训练计时器**: 支持自定义收缩、保持、放松时间和循环次数
- **打卡记录**: 记录每次训练的时长和循环数
- **热力图**: 可视化展示打卡记录
- **排行榜**: 展示用户连续打卡天数排名
- **个人中心**: 查看个人训练统计和历史记录

### UI/UX 特性

- 全屏沉浸式海洋背景
- 多层海浪动画（上下涌动效果）
- 动态气泡上升效果
- 涟漪向外扩散效果
- 玻璃态透明组件
- 透明悬浮导航栏
- 海洋色系渐变
- 响应式设计

## 页面说明

### Home.vue 首页

- 全屏英雄区 (`min-height: calc(100vh - 72px)`)
- Logo 带涟漪扩散和海浪填充效果
- 统计数据卡片
- 功能特性网格
- 全站热力图
- CTA 区域（海浪涌动图标）

### Train.vue 训练页

- Three.js 动态海洋背景
- 环形计时器
- 潮汐阶段提示（涨潮、平潮、退潮）
- 打卡成功庆祝动画

### Leaderboard.vue 排行榜

- 奖杯图标（内部海浪填充动画）
- Top 3 领奖台展示
- 统计卡片
- 排行榜表格

### Dashboard.vue 个人中心

- 用户信息卡片
- 连续打卡统计
- 个人热力图
- 最近训练记录

## 组件说明

### MainLayout.vue 主布局

负责全屏海洋背景系统：

```vue
<template>
  <div class="main-layout">
    <!-- 全屏海洋背景层 -->
    <div class="ocean-bg"></div>
    <div class="ocean-gradient"></div>
    <div class="animated-bg"></div>
    <div class="wave-layer">
      <div class="wave wave-1"></div>
      <div class="wave wave-2"></div>
      <div class="wave wave-3"></div>
    </div>
    <div class="ripple-layer"></div>
    <div class="bubbles-container">
      <!-- 动态气泡 -->
    </div>

    <!-- 透明悬浮导航 -->
    <header class="main-header" :class="{ scrolled }">
      <!-- ... -->
    </header>

    <!-- 内容区 -->
    <main class="main-content">
      <slot />
    </main>
  </div>
</template>
```

### Timer.vue 训练计时器

- 环形进度条显示当前阶段进度
- 相位颜色：收缩(红)、保持(黄)、放松(绿)
- 海洋主题按钮样式

### TidalBackground.vue

Three.js 动态海洋背景，用于训练页面：

- 波浪网格动画
- 根据训练相位调整颜色和强度
- 粒子效果

## 样式规范

### 页面容器

所有页面使用全屏高度：

```css
.page {
  min-height: calc(100vh - 72px);
  padding: 40px 24px 80px;
}
```

### 玻璃态卡片

```css
.glass-card {
  background: rgba(8, 15, 30, 0.4);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.08);
  border-radius: 18px;
}
```

### 海洋渐变文字

```css
.gradient-text {
  background: linear-gradient(135deg, #38bdf8, #22d3ee, #0ea5e9);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
```

### 海洋按钮

```css
.btn-ocean {
  background: linear-gradient(135deg, #0ea5e9, #0284c7);
  box-shadow: 0 4px 20px rgba(14, 165, 233, 0.4);
}
```

## API 配置

API 基础 URL 通过环境变量配置：

```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080/api

# .env.production
VITE_API_BASE_URL=/api
```

## Docker 部署

### 构建镜像

```bash
docker build -t tidalcore-frontend .
```

### 运行容器

```bash
docker run -d -p 80:80 tidalcore-frontend
```

## 开发指南

### 添加新页面

1. 在 `src/views/` 创建新的 Vue 组件
2. 使用全屏布局和海洋主题样式
3. 在 `src/router/index.ts` 添加路由配置

### 页面模板示例

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'

const mounted = ref(false)

onMounted(() => {
  setTimeout(() => {
    mounted.value = true
  }, 100)
})
</script>

<template>
  <MainLayout>
    <div class="page">
      <section class="header-section" :class="{ mounted }">
        <h1 class="gradient-text">页面标题</h1>
        <p class="subtitle">页面描述</p>
      </section>

      <section class="content-section" :class="{ mounted }">
        <!-- 内容 -->
      </section>
    </div>
  </MainLayout>
</template>

<style scoped>
.page {
  min-height: calc(100vh - 72px);
  padding: 40px 24px 80px;
}

.header-section {
  text-align: center;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth);
}

.header-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.gradient-text {
  background: linear-gradient(135deg, rgb(var(--ocean-surface)), rgb(var(--aqua-glow)));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  color: rgba(255, 255, 255, 0.5);
}
</style>
```

## 浏览器支持

- Chrome >= 90
- Firefox >= 90
- Safari >= 14
- Edge >= 90

## 许可证

MIT License
