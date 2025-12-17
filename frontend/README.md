# TidalCore Frontend

TidalCore 前端应用 - 专业的盆底肌训练平台。

## 技术栈

- **框架**: Vue 3.5 + TypeScript + Composition API
- **构建工具**: Vite 7
- **状态管理**: Pinia 3
- **路由**: Vue Router 4
- **UI 组件库**: Element Plus 2.x (暗色主题)
- **样式**: Tailwind CSS 4 + 自定义 CSS
- **图标**: @element-plus/icons-vue
- **HTTP 客户端**: Axios
- **动画**: Canvas Confetti, VueUse Motion
- **工具库**: VueUse

## Element Plus 配置

项目使用 Element Plus 作为主要 UI 组件库，并配置了暗色主题：

```typescript
// main.ts
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 全局注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(ElementPlus)
```

### 常用组件

| Element Plus 组件 | 用途 |
|------------------|------|
| `el-card` | 卡片容器 |
| `el-button` | 按钮 |
| `el-input` | 输入框 |
| `el-form` / `el-form-item` | 表单 |
| `el-dialog` | 模态框 |
| `el-drawer` | 抽屉 |
| `el-menu` | 导航菜单 |
| `el-row` / `el-col` | 栅格布局 |
| `el-tag` | 标签 |
| `el-alert` | 提示 |
| `el-skeleton` | 骨架屏 |
| `el-divider` | 分割线 |
| `el-collapse-transition` | 折叠过渡 |

### 图标使用

```vue
<script setup>
import { Timer, Trophy, User } from '@element-plus/icons-vue'
</script>

<template>
  <el-icon><Timer /></el-icon>
  <el-icon :size="20" color="#a78bfa"><Trophy /></el-icon>
</template>
```

## 项目结构

```
src/
├── api/                    # API 接口层
│   ├── auth.ts            # 认证相关 API
│   ├── checkin.ts         # 打卡相关 API
│   └── request.ts         # Axios 实例配置
├── assets/                 # 静态资源
├── components/             # 组件
│   ├── ui/                # 通用 UI 组件 (保留兼容)
│   ├── Heatmap.vue        # 热力图组件
│   ├── LeaderboardTable.vue # 排行榜表格组件
│   └── Timer.vue          # 训练计时器组件
├── layouts/                # 布局组件
│   └── MainLayout.vue     # 主布局 (Element Plus)
├── router/                 # 路由配置
│   └── index.ts
├── store/                  # 状态管理
│   ├── training.ts        # 训练状态
│   └── user.ts            # 用户状态
├── views/                  # 页面视图
│   ├── Dashboard.vue      # 个人中心
│   ├── Home.vue           # 首页
│   ├── Leaderboard.vue    # 排行榜
│   ├── Login.vue          # 登录
│   ├── Register.vue       # 注册
│   └── Train.vue          # 训练页面
├── App.vue                 # 根组件
├── main.ts                 # 入口文件
└── style.css              # 全局样式
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

### 预览生产版本

```bash
npm run preview
```

## 功能特性

### 核心功能

- **训练计时器**: 支持自定义收缩、保持、放松时间和循环次数
- **打卡记录**: 记录每次训练的时长和循环数
- **热力图**: 可视化展示打卡记录
- **排行榜**: 展示用户连续打卡天数排名
- **个人中心**: 查看个人训练统计和历史记录

### UI/UX 特性

- 响应式设计，支持移动端和桌面端
- Element Plus 暗色主题，护眼设计
- 流畅的过渡动画和微交互
- 玻璃拟态 (Glassmorphism) 设计风格
- 渐变色装饰和光晕效果

## 设计系统

### 主题配色

```css
/* 主渐变色 */
background: linear-gradient(135deg, #8b5cf6, #ec4899);  /* 紫粉渐变 */
background: linear-gradient(135deg, #fbbf24, #f97316);  /* 金橙渐变 */
background: linear-gradient(135deg, #10b981, #14b8a6);  /* 绿青渐变 */

/* 强调色 */
--aurora-purple: #8b5cf6;   /* 紫色 */
--aurora-pink: #ec4899;     /* 粉色 */
--aurora-cyan: #22d3ee;     /* 青色 */
--aurora-emerald: #34d399;  /* 绿色 */
--aurora-amber: #fbbf24;    /* 琥珀色 */
--aurora-orange: #fb923c;   /* 橙色 */
```

### 暗色背景

```css
/* 卡片背景 */
background: rgba(30, 30, 46, 0.8);
border: 1px solid rgba(255, 255, 255, 0.06);

/* 装饰光晕 */
background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), transparent);
filter: blur(60px);
```

### Element Plus 样式覆盖

使用 `:deep()` 选择器覆盖 Element Plus 默认样式：

```css
.my-card :deep(.el-card__body) {
  padding: 32px;
}

.my-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: none !important;
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
2. 使用 Element Plus 组件构建 UI
3. 在 `src/router/index.ts` 添加路由配置
4. 如需要，在 `src/layouts/MainLayout.vue` 添加导航链接

### 页面模板示例

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { Timer, Star } from '@element-plus/icons-vue'

const mounted = ref(false)

onMounted(() => {
  setTimeout(() => {
    mounted.value = true
  }, 100)
})
</script>

<template>
  <MainLayout>
    <div class="page-container">
      <section class="hero-section" :class="{ mounted }">
        <el-tag type="primary" effect="plain" round>
          <el-icon><Star /></el-icon>
          标签文字
        </el-tag>
        <h1 class="page-title">
          <span class="gradient-text">页面标题</span>
        </h1>
      </section>

      <section class="content-section" :class="{ mounted }">
        <el-card shadow="never" class="content-card">
          <!-- 内容 -->
        </el-card>
      </section>
    </div>
  </MainLayout>
</template>

<style scoped>
.page-container {
  display: flex;
  flex-direction: column;
  gap: 48px;
}

.hero-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease;
}

.hero-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.gradient-text {
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.content-card {
  background: rgba(30, 30, 46, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 16px !important;
}
</style>
```

### 状态管理

使用 Pinia 进行状态管理：

```typescript
// 定义 Store
export const useExampleStore = defineStore('example', () => {
  const state = ref<Type>(initialValue)

  function action() {
    // ...
  }

  return { state, action }
})

// 使用 Store
const store = useExampleStore()
```

### API 调用

```typescript
// 在 src/api/ 中定义接口
export function getData(): Promise<DataType> {
  return request.get('/endpoint')
}

// 在组件中使用
import { getData } from '@/api/example'

const data = await getData()
```

## 浏览器支持

- Chrome >= 90
- Firefox >= 90
- Safari >= 14
- Edge >= 90

## 许可证

MIT License
