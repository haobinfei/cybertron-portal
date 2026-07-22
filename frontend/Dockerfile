# frontend/Dockerfile
# 多阶段构建：先编译 Vue/Vite 前端，再用 nginx 托管

# 阶段 1：编译
FROM registry.cn-hangzhou.aliyuncs.com/ss-vip/node:20-alpine AS builder

WORKDIR /app

# 配置 npm 国内镜像源，解决 CI 构建时访问 npm registry 超时
ENV NPM_CONFIG_REGISTRY=https://registry.npmmirror.com

# 复制依赖文件并安装
# --include=dev：强制安装 devDependencies（vue-tsc 等构建工具），
# 防止基础镜像携带 NODE_ENV=production 导致 npm ci 跳过 dev 依赖
COPY package*.json ./
RUN npm ci --include=dev

# 复制源码并构建
COPY . .
RUN npm run build

# 阶段 2：运行（nginx 托管静态资源）
FROM registry.cn-hangzhou.aliyuncs.com/ss-vip/nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露 80 端口
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
