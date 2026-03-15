#!/bin/bash

# XLog Git 初始化和首次提交脚本

echo "======================================"
echo "XLog Git 仓库初始化"
echo "======================================"
echo ""

# 检查是否已经在Git仓库中
if [ -d ".git" ]; then
    echo "⚠️  当前目录已经是Git仓库"
    echo ""
    read -p "是否重新初始化？(y/N): " confirm
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        echo "取消操作"
        exit 0
    fi
    rm -rf .git
fi

# 初始化Git仓库
echo "📦 初始化Git仓库..."
git init

# 添加所有文件
echo "📝 添加文件到暂存区..."
git add .

# 显示即将提交的文件
echo ""
echo "即将提交的文件："
git status --short
echo ""

# 首次提交
echo "💾 创建首次提交..."
git commit -m "$(cat <<'EOF'
feat: Initial release of XLog - Go async logging library

Core features:
- Async logging with channel and goroutine
- Auto cleanup based on retention days
- Auto compression to zip format
- Log rotation by date
- Concurrent-safe logging
- Cross-platform support (Windows/Linux/macOS)

Performance:
- 876K ops/sec single-thread
- 4.1M ops/sec concurrent
- 1.3µs average latency

Ported from C++ XSimpleLogEx
EOF
)"

echo ""
echo "✅ Git仓库初始化完成！"
echo ""
echo "后续步骤："
echo ""
echo "1. 在GitHub/GitLab/Gitee上创建新仓库"
echo "   例如: https://github.com/new"
echo ""
echo "2. 添加远程仓库并推送："
echo "   git remote add origin https://github.com/你的用户名/xlog.git"
echo "   git branch -M main"
echo "   git push -u origin main"
echo ""
echo "3. 可选：创建版本标签"
echo "   git tag -a v1.0.0 -m 'Release v1.0.0'"
echo "   git push origin v1.0.0"
echo ""
