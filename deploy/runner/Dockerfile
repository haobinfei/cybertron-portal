# deploy/runner/Dockerfile
# ARC (Actions Runner Controller) self-hosted runner 镜像
# 预装 CI 所需工具，避免 workflow 每次运行时重复下载

FROM ghcr.io/actions/actions-runner:latest

USER root

# 系统依赖与常用工具
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    git \
    jq \
    unzip \
    wget \
    && rm -rf /var/lib/apt/lists/*

# 安装 Docker CLI（用于 docker/build-push-action）
RUN install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc && \
    chmod a+r /etc/apt/keyrings/docker.asc && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" \
      > /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends docker-ce-cli docker-buildx-plugin && \
    rm -rf /var/lib/apt/lists/*

# 安装 kubectl（用于 create-image-pull-secret 等步骤）
ARG KUBECTL_VERSION=v1.30.0
RUN curl -fsSL "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" \
      -o /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl

# 安装 Helm（用于 package-helm-chart 步骤）
ARG HELM_VERSION=v3.14.0
RUN curl -fsSL "https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz" | tar xz && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    rm -rf linux-amd64

# 预装 Buildx 到 runner 用户的 docker cli-plugins 目录
# 同时支持 root 和 runner 用户
ARG BUILDX_VERSION=v0.35.0
RUN mkdir -p /usr/libexec/docker/cli-plugins /home/runner/.docker/cli-plugins && \
    curl -fsSL "https://github.com/docker/buildx/releases/download/${BUILDX_VERSION}/buildx-${BUILDX_VERSION}.linux-amd64" \
      -o /usr/libexec/docker/cli-plugins/docker-buildx && \
    chmod +x /usr/libexec/docker/cli-plugins/docker-buildx && \
    cp /usr/libexec/docker/cli-plugins/docker-buildx /home/runner/.docker/cli-plugins/docker-buildx && \
    chmod +x /home/runner/.docker/cli-plugins/docker-buildx && \
    chown -R runner:runner /home/runner/.docker

# 验证安装
RUN docker --version && \
    docker buildx version && \
    kubectl version --client && \
    helm version --short

USER runner
