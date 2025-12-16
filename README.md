# 🌊 TidalCore (潮汐核心)

**TidalCore** 是一款开源、免费、纯粹的 Web 盆底肌（提肛/凯格尔运动）训练平台。  
其命名灵感源于盆底肌收缩与放松时的律动，犹如潮汐一般有节奏且蕴含力量。

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Vue](https://img.shields.io/badge/Frontend-Vue%203-brightgreen.svg)](https://vuejs.org/)
[![Gin](https://img.shields.io/badge/Backend-Gin%20(Go)-00ADD8.svg)](https://gin-gonic.com/)

---

## ✨ 项目愿景
我们致力于消除“健康焦虑”与“病耻感”，提供一个无需复杂注册、保护隐私、且具有良好交互反馈的训练环境。通过开源的力量，让盆底肌健康管理变得像呼吸一样自然。

## 🌟 核心功能

- ⏱️ **潮汐训练器**：精心设计的交互式计时器，支持自定义“收缩-保持-放松”循环。
- 🎉 **胜利庆祝**：完成挑战时触发炫酷的 Canvas 五彩纸屑特效，赋予运动仪式感。
- 📊 **匿名互助圈**：
  - **打卡热力图**：直观展示全站及个人的坚持轨迹。
  - **毅力排行榜**：基于连续打卡天数的匿名排名，激发动力。
- 🔒 **极致隐私**：仅需用户名与密码即可注册，不收集邮箱、手机号等任何个人敏感信息。
- 📱 **纯粹 Web 体验**：采用 Tailwind CSS 响应式设计，完美适配手机与 PC 浏览器，无需下载 App。

## 🛠️ 技术栈

### 前端 (Frontend)
- **框架**: Vue 3 (Composition API)
- **构建**: Vite
- **状态管理**: Pinia
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
├── frontend/           # Vue 3 前端应用 (Vite)
├── scripts/            # 数据库初始化 SQL 及工具脚本
└── docker-compose.yml  # 一键部署配置