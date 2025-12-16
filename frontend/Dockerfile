# 构建阶段
FROM node:20-alpine AS builder

WORKDIR /app

# 复制 package.json 和 lock 文件
COPY package.json package-lock.json* ./

# 安装依赖
RUN npm ci || npm install

# 复制源代码（排除 node_modules，见 .dockerignore）
COPY . .

# 构建生产版本
RUN npm run build

# 运行阶段 - 使用 Nginx 提供静态文件
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
