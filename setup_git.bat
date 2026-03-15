@echo off
chcp 65001 >nul
setlocal

echo ======================================
echo XLog Git 仓库初始化
echo ======================================
echo.

REM 检查是否已经在Git仓库中
if exist .git (
    echo ⚠️  当前目录已经是Git仓库
    echo.
    set /p confirm="是否重新初始化？(y/N): "
    if /i not "%confirm%"=="y" (
        echo 取消操作
        exit /b 0
    )
    rmdir /s /q .git
)

REM 初始化Git仓库
echo 📦 初始化Git仓库...
git init

REM 添加所有文件
echo 📝 添加文件到暂存区...
git add .

REM 显示即将提交的文件
echo.
echo 即将提交的文件：
git status --short
echo.

REM 首次提交
echo 💾 创建首次提交...
git commit -m "feat: Initial release of XLog - Go async logging library" -m "Core features:" -m "- Async logging with channel and goroutine" -m "- Auto cleanup based on retention days" -m "- Auto compression to zip format" -m "- Log rotation by date" -m "- Concurrent-safe logging" -m "- Cross-platform support (Windows/Linux/macOS)" -m "" -m "Performance:" -m "- 876K ops/sec single-thread" -m "- 4.1M ops/sec concurrent" -m "- 1.3µs average latency" -m "" -m "Ported from C++ XSimpleLogEx"

echo.
echo ✅ Git仓库初始化完成！
echo.
echo 后续步骤：
echo.
echo 1. 在GitHub/GitLab/Gitee上创建新仓库
echo    例如: https://github.com/new
echo.
echo 2. 添加远程仓库并推送：
echo    git remote add origin https://github.com/你的用户名/xlog.git
echo    git branch -M main
echo    git push -u origin main
echo.
echo 3. 可选：创建版本标签
echo    git tag -a v1.0.0 -m "Release v1.0.0"
echo    git push origin v1.0.0
echo.
pause
