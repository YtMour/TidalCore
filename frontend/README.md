# TidalCore Frontend

TidalCore 前端应用 - 专业的盆底肌训练平台。

## 技术栈

- **框架**: Vue 3.5 + TypeScript
- **构建工具**: Vite 7
- **状态管理**: Pinia 3
- **路由**: Vue Router 4
- **样式**: Tailwind CSS 4
- **图标**: Lucide Vue Next
- **HTTP 客户端**: Axios
- **动画**: Canvas Confetti, VueUse Motion

## 项目结构

```
src/
├── api/                    # API 接口层
│   ├── auth.ts            # 认证相关 API
│   ├── checkin.ts         # 打卡相关 API
│   └── request.ts         # Axios 实例配置
├── assets/                 # 静态资源
├── components/             # 组件
│   ├── ui/                # 通用 UI 组件
│   │   ├── Button.vue     # 按钮组件
│   │   ├── Card.vue       # 卡片组件
│   │   ├── Input.vue      # 输入框组件
│   │   ├── Modal.vue      # 模态框组件
│   │   ├── Skeleton.vue   # 骨架屏组件
│   │   ├── StatCard.vue   # 统计卡片组件
│   │   └── index.ts       # 组件导出
│   ├── Heatmap.vue        # 热力图组件
│   ├── LeaderboardTable.vue # 排行榜表格组件
│   └── Timer.vue          # 训练计时器组件
├── layouts/                # 布局组件
│   └── MainLayout.vue     # 主布局
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
- 深色主题，护眼设计
- 流畅的过渡动画
- 玻璃拟态 (Glassmorphism) 设计风格

## 设计系统

### 颜色变量

```css
/* 主色调 */
--color-primary: 99 102 241;      /* Indigo */
--color-primary-light: 129 140 248;
--color-primary-dark: 79 70 229;

/* 强调色 */
--aurora-purple: 139 92 246;
--aurora-pink: 236 72 153;
--aurora-cyan: 34 211 238;
--aurora-emerald: 52 211 153;
--aurora-amber: 251 191 36;
```

### 圆角规范

```css
--radius-sm: 6px;
--radius-md: 8px;
--radius-lg: 12px;
--radius-xl: 16px;
--radius-2xl: 20px;
```

### 动画时长

```css
--duration-fast: 150ms;
--duration-normal: 250ms;
--duration-slow: 400ms;
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
2. 在 `src/router/index.ts` 添加路由配置
3. 如需要，在 `src/layouts/MainLayout.vue` 添加导航链接

### 添加新组件

1. 通用 UI 组件放在 `src/components/ui/`
2. 业务组件放在 `src/components/`
3. 在 `src/components/ui/index.ts` 导出通用组件

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
