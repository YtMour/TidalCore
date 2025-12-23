# 🌊 TidalCore (潮汐核心)

**TidalCore** 是一款开源、免费、纯粹的 Web 盆底肌（提肛/凯格尔运动）训练平台。
其命名灵感源于盆底肌收缩与放松时的律动，犹如潮汐一般有节奏且蕴含力量。

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Vue](https://img.shields.io/badge/Frontend-Vue%203-brightgreen.svg)](https://vuejs.org/)
[![Gin](https://img.shields.io/badge/Backend-Gin%20(Go)-00ADD8.svg)](https://gin-gonic.com/)

---

**在线演示**：[https://tidal.ytmour.baby/]

## ✨ 项目愿景
我们致力于消除"健康焦虑"与"病耻感"，提供一个无需复杂注册、保护隐私、且具有良好交互反馈的训练环境。通过开源的力量，让盆底肌健康管理变得像呼吸一样自然。

## 🌟 核心功能

- ⏱️ **潮汐训练器**：精心设计的交互式计时器，支持自定义"收缩-保持-放松"循环。
- 🎉 **胜利庆祝**：完成挑战时触发炫酷的 Canvas 五彩纸屑特效，赋予运动仪式感。
- 📊 **匿名互助圈**：
  - **打卡热力图**：直观展示全站及个人的坚持轨迹。
  - **毅力排行榜**：基于连续打卡天数的匿名排名，激发动力。
- 🏅 **称号系统**：根据累计打卡次数自动获得海洋主题称号，每个称号拥有独特的视觉特效。
- 🔒 **极致隐私**：仅需用户名与密码即可注册，不收集邮箱、手机号等任何个人敏感信息。
- 📱 **纯粹 Web 体验**：采用 Tailwind CSS 响应式设计，完美适配手机与 PC 浏览器，无需下载 App。
- ⚙️ **管理后台**：管理员可管理用户数据、设置称号、调整打卡统计等。

## 🏅 称号系统

根据累计打卡次数，用户将自动获得对应的海洋主题称号：

| 打卡次数 | 称号 | 特效 |
|---------|------|------|
| 1000+ | 🔱 海神降临 | 神圣光芒 + 彩虹流光 |
| 730+ | 🦑 深渊霸主 | 暗紫脉动 + 触手波动 |
| 365+ | 🌊 深海传奇 | 金色闪耀 + 光线扫过 |
| 180+ | 🐋 海洋大师 | 水波纹扩散 |
| 90+ | 🐬 浪潮专家 | 轻微浮动 |
| 30+ | 🐠 潮汐进阶 | 柔和呼吸 |
| 7+ | 🐟 入海新手 | 轻微闪烁 |
| 0+ | 🐚 初探海域 | - |


## 🛠️ 技术栈

### 前端 (Frontend)
- **框架**: Vue 3 (Composition API)
- **构建**: Vite
- **状态管理**: Pinia
- **UI 组件**: Element Plus
- **样式**: Tailwind CSS
- **动画/特效**: `canvas-confetti` & `Motion One`

### 后端 (Backend)
- **语言**: Golang 1.2x
- **Web 框架**: Gin
- **数据库**: MySQL 8.0
- **持久层**: GORM
- **认证**: JWT (JSON Web Token)

## 📂 目录结构

```text
TidalCore/
├── backend/            # Go 后端服务 (Gin 架构)
│   ├── api/            # API 路由和处理器
│   ├── internal/       # 内部模块
│   │   ├── model/      # 数据模型
│   │   ├── repository/ # 数据访问层
│   │   └── service/    # 业务逻辑层
│   └── middleware/     # 中间件
├── frontend/           # Vue 3 前端应用 (Vite)
│   ├── src/
│   │   ├── api/        # API 接口
│   │   ├── components/ # 组件
│   │   ├── views/      # 页面视图
│   │   └── store/      # 状态管理
├── scripts/            # 数据库初始化 SQL 及工具脚本
└── docker-compose.yml  # 一键部署配置
```

## 🗄️ 数据库结构

### 用户表 (users)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| username | string | 用户名（唯一） |
| display_name | string | 显示名称 |
| password_hash | string | 密码哈希 |
| is_admin | bool | 是否为管理员 |
| title | string | 用户称号（空则自动计算） |
| streak | int | 当前连续打卡天数 |
| max_streak | int | 最高连续打卡天数 |
| total_checkin | int | 累计打卡次数 |
| last_checkin | time | 最后打卡时间 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

## 🔌 API 接口

### 管理员接口

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/users` | 获取用户列表 |
| DELETE | `/api/v1/admin/users/:id` | 删除用户 |
| PUT | `/api/v1/admin/users/:id/admin` | 设置管理员权限 |
| PUT | `/api/v1/admin/users/:id/stats` | 更新用户统计数据和称号 |

### 用户统计数据更新参数

```json
{
  "streak": 10,
  "max_streak": 15,
  "total_checkin": 100,
  "title": "海洋大师"
}
```

- `title` 为空字符串时，称号将根据 `total_checkin` 自动计算
- `title` 设置为具体称号名称时，将覆盖自动计算的结果