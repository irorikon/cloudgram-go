#!/bin/bash

set -e

# 配置变量
IMAGE_NAME="${1:-cloudgram-go}"
TAG="${2:-latest}"
PUSH="${3:-true}"  # 第三个参数控制是否推送，默认为 true
PLATFORMS="linux/amd64,linux/arm64,linux/arm/v7"

echo "Building multi-architecture Docker image: $IMAGE_NAME:$TAG"
echo "Platforms: $PLATFORMS"
echo "Push to registry: $PUSH"

# 检查是否已登录到 Docker Hub（仅在需要推送时检查）
if [[ "$PUSH" == "true" ]]; then
    if ! docker info | grep -q "Username"; then
        echo "Please login to Docker Hub first: docker login"
        exit 1
    fi
fi

# 创建并使用新的 builder 实例（支持多平台）
if ! docker buildx inspect mybuilder >/dev/null 2>&1; then
    echo "Creating new buildx builder..."
    docker buildx create --name mybuilder --use
else
    docker buildx use mybuilder
fi

# Bootstrap the builder
docker buildx inspect --bootstrap

# 构建多架构镜像
BUILD_ARGS=(
  --platform "$PLATFORMS"
  --tag "$IMAGE_NAME:$TAG"
  --cache-from "type=registry,ref=$IMAGE_NAME:buildcache"  # 从远程缓存读取
  --cache-to "type=registry,ref=$IMAGE_NAME:buildcache,mode=max"  # 写入远程缓存
)

# 根据 PUSH 参数决定是否添加 --push 标志
if [[ "$PUSH" == "true" ]]; then
    BUILD_ARGS+=(--push)
    echo "Building and pushing multi-architecture image..."
else
    BUILD_ARGS+=(--load)  # 加载到本地 Docker 守护进程
    echo "Building multi-architecture image locally (no push)..."
fi

docker buildx build "${BUILD_ARGS[@]}" .

if [[ "$PUSH" == "true" ]]; then
    echo "Multi-architecture build completed successfully!"
    echo "Image pushed to: $IMAGE_NAME:$TAG"
else
    echo "Multi-architecture build completed successfully!"
    echo "Image loaded locally as: $IMAGE_NAME:$TAG"
fi