#!/bin/bash

# Sivi Go SDK Git发布脚本
# 用法: ./release.sh [版本号]
# 示例: ./release.sh v1.2.3

set -e

echo "开始发布流程..."

# 检查是否有未提交的更改
if [[ -n $(git status --porcelain) ]]; then
    echo "发现未提交的更改，正在提交..."
    git add .
    git commit -m "Auto commit: $(date '+%Y-%m-%d %H:%M:%S')"
else
    echo "没有未提交的更改"
fi

# 获取当前commit ID
COMMIT_ID=$(git rev-parse --short HEAD)

# 获取GMT+8时间戳
TIMESTAMP=$(TZ='Asia/Shanghai' date '+%Y%m%d%H%M%S')

# 版本号设置
if [[ -n "$1" ]]; then
    # 使用指定版本号 + 时间戳 + commit ID
    VERSION="$1.${TIMESTAMP}.${COMMIT_ID}"
    echo "使用指定版本号: ${VERSION}"
else
    # 使用 0.0.0 + 时间戳 + commit ID
    VERSION="0.0.0.${TIMESTAMP}.${COMMIT_ID}"
    echo "使用默认版本号: ${VERSION}"
fi

# 检查标签是否已存在
if git rev-parse --verify "refs/tags/${VERSION}" >/dev/null 2>&1; then
    echo "⚠️  警告: 标签 ${VERSION} 已存在"
    echo "是否要覆盖? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "❌ 取消发布"
        exit 1
    fi
    echo "🔄 将覆盖已存在的标签"
fi

# 检查代码格式
echo "检查代码格式..."
go fmt ./...

# 整理依赖
echo "整理依赖..."
go mod tidy

# 推送到远程仓库
echo "推送到远程仓库..."
git push origin master

# 创建并推送标签
echo "创建版本标签: ${VERSION}"
git tag -f ${VERSION}  # -f 强制覆盖已存在的标签
git push origin ${VERSION} --force  # 强制推送标签

# 获取远程仓库地址
REPO_URL=$(git config --get remote.origin.url)
if [[ $REPO_URL == git@* ]]; then
    # SSH格式转换为HTTPS
    REPO_URL=$(echo $REPO_URL | sed 's/git@github.com:/https:\/\/github.com\//')
    REPO_URL=$(echo $REPO_URL | sed 's/\.git$//')
fi

echo ""
echo "✅ 发布成功！"
echo "📦 版本: ${VERSION}"
echo "🔗 仓库地址: ${REPO_URL}"
echo "📋 使用方法:"
echo "   go get ${REPO_URL#https://}@${VERSION}"
echo ""
echo "或在go.mod中添加:"
echo "   require ${REPO_URL#https://} ${VERSION}"
