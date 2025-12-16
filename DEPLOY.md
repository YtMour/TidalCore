# TidalCore 宝塔面板部署指南

## 部署架构

```
用户请求 → 宝塔Nginx(SSL/域名) → 反向代理 → Docker容器
                                    ├── frontend:3000 (前端)
                                    └── backend:8080  (后端API)
                                          └── mysql:3306 (数据库)
```

## 部署流程

### 1. 上传项目到服务器

将整个项目文件夹上传到宝塔站点目录，例如：

```
/www/wwwroot/tidalcore.yourdomain.com/
├── backend/
├── frontend/
├── docker-compose.yml
└── ...
```

### 2. 修改 docker-compose.yml 配置

编辑 `docker-compose.yml`，修改以下配置：

```yaml
services:
  mysql:
    environment:
      MYSQL_ROOT_PASSWORD: 你的强密码        # 必改
      MYSQL_PASSWORD: 你的强密码             # 必改

  backend:
    environment:
      DB_PASSWORD: 你的强密码                # 与上面 MYSQL_PASSWORD 一致
      JWT_SECRET: 随机字符串                 # 必改，用于加密token
      ALLOWED_ORIGINS: https://你的域名.com  # 必改，你的实际域名
```

**生成 JWT_SECRET：**
```bash
openssl rand -base64 32
```

### 3. 启动 Docker 容器

SSH 进入服务器，执行：

```bash
cd /www/wwwroot/tidalcore.yourdomain.com

# 构建并启动所有容器
docker-compose up -d --build

# 查看运行状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

首次启动需要几分钟构建镜像，等待看到：
```
tidalcore-backend | Server starting on :8080
```

### 4. 宝塔面板创建站点

1. **网站 → 添加站点**
   - 域名：`tidalcore.yourdomain.com`
   - 根目录：随意（反向代理不需要）
   - PHP版本：纯静态
   - 数据库：不创建

2. **SSL → Let's Encrypt**
   - 申请证书
   - 开启强制 HTTPS

### 5. 配置反向代理

**网站 → 站点设置 → 反向代理 → 添加反向代理**

| 配置项 | 值 |
|--------|-----|
| 代理名称 | tidalcore |
| 目标URL | http://127.0.0.1:3000 |
| 发送域名 | $host |

点击**保存**。

### 6. 验证部署

- 访问 `https://你的域名.com` 查看前端页面
- 访问 `https://你的域名.com/health` 验证后端连接
- 尝试注册账号测试完整功能

---

## 常用命令

```bash
cd /www/wwwroot/tidalcore.yourdomain.com

# 查看容器状态
docker-compose ps

# 查看日志
docker-compose logs -f
docker-compose logs -f backend   # 只看后端
docker-compose logs -f mysql     # 只看数据库

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 重新构建并启动
docker-compose up -d --build

# 进入数据库
docker exec -it tidalcore-mysql mysql -u root -p
```

---

## 故障排查

### 502 Bad Gateway

```bash
# 检查容器是否运行
docker-compose ps

# 检查后端日志
docker-compose logs backend

# 测试后端是否响应
curl http://127.0.0.1:8080/health
curl http://127.0.0.1:3000
```

### 数据库连接失败

```bash
# 查看 MySQL 日志
docker-compose logs mysql

# 确认密码配置一致
grep PASSWORD docker-compose.yml
```

### CORS 跨域错误

检查 `docker-compose.yml` 中的 `ALLOWED_ORIGINS` 是否设置为你的域名（包含 https://）。

### 前端页面刷新 404

前端容器内的 Nginx 已配置 SPA 路由支持，如果仍有问题，检查宝塔反向代理配置。

---

## 数据备份

```bash
# 备份数据库
docker exec tidalcore-mysql mysqldump -u root -p'你的密码' tidalcore > backup_$(date +%Y%m%d).sql

# 恢复数据库
docker exec -i tidalcore-mysql mysql -u root -p'你的密码' tidalcore < backup.sql
```

---

## 更新部署

```bash
cd /www/wwwroot/tidalcore.yourdomain.com

# 拉取最新代码（如果用 git）
git pull

# 重新构建并启动
docker-compose up -d --build
```
