# 在 aiproxy 中集成 xlog 日志库

本文档说明如何在 aiproxy 项目中引用 xlog 日志库。

## 方式一：本地引用（推荐用于开发阶段）

### 1. 在 aiproxy 的 go.mod 中添加 replace 指令

编辑 `aiproxy/go.mod` 文件，在末尾添加：

```go
replace github.com/chen1994/xlog => T:/Library/go/xlog
```

### 2. 在 go.mod 中添加依赖

```go
require github.com/chen1994/xlog v0.0.0
```

### 3. 在代码中使用

```go
package main

import (
    "github.com/chen1994/xlog"
)

func main() {
    // 使用默认配置
    config := xlog.DefaultConfig()
    config.LogDir = "./logs"
    config.LogBaseName = "aiproxy"
    config.AutoCleanup = true
    config.AutoCompress = true

    logger, err := xlog.NewXLogger(config)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    // 记录日志
    logger.Log("应用启动")
    logger.LogEx(xlog.LevelInfo, "服务运行在端口 %d", 8080)
}
```

### 4. 更新依赖

```bash
cd aiproxy
go mod tidy
```

---

## 方式二：Git 仓库引用（推荐用于生产环境）

### 1. 先将 xlog 推送到你的 Git 仓库

```bash
cd T:/Library/go/xlog
./setup_git.sh  # 或 setup_git.bat

# 在GitHub/Gitee上创建仓库后
git remote add origin https://github.com/你的用户名/xlog.git
git push -u origin main
```

### 2. 在 aiproxy 中引用

编辑 `aiproxy/go.mod`:

```go
require github.com/chen1994/xlog v0.0.0

// 或者使用具体的Git commit
require github.com/chen1994/xlog v0.0.0-20260315000000-xxxxxxxxxxxx
```

然后运行：

```bash
cd aiproxy
go get github.com/chen1994/xlog
go mod tidy
```

---

## 方式三：复制到 vendor 目录（不推荐）

```bash
mkdir -p aiproxy/vendor/github.com/chen1994
cp -r T:/Library/go/xlog aiproxy/vendor/github.com/chen1994/xlog
```

然后在 go.mod 中添加：

```go
require github.com/chen1994/xlog v0.0.0
```

---

## 替换现有的日志实现

假设 aiproxy 当前使用标准库的 log，可以这样替换：

### 原代码：
```go
package main

import (
    "log"
    "os"
)

func main() {
    log.Println("应用启动")
    log.Println("服务运行")
}
```

### 替换为 xlog：
```go
package main

import (
    "github.com/chen1994/xlog"
)

var logger *xlog.XLogger

func main() {
    config := xlog.DefaultConfig()
    config.LogDir = "./logs"
    config.LogBaseName = "aiproxy"

    var err error
    logger, err = xlog.NewXLogger(config)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    logger.Log("应用启动")
    logger.Log("服务运行")
}
```

---

## 配置建议

### 开发环境配置：
```go
config := xlog.DefaultConfig()
config.LogDir = "./logs"
config.LogBaseName = "aiproxy-dev"
config.HasDate = true
config.HasProcessID = false
config.UTF8Format = true
config.AutoCleanup = false  // 开发时不清理
config.AutoCompress = false // 开发时不压缩
```

### 生产环境配置：
```go
config := xlog.DefaultConfig()
config.LogDir = "/var/log/aiproxy"
config.LogBaseName = "aiproxy"
config.HasDate = true
config.HasProcessID = true
config.UTF8Format = true
config.WriteInterval = 10 * time.Second
config.BufferSize = 20000

// 自动清理：保留7天
config.AutoCleanup = true
config.CleanupCycle = 1 * time.Hour
config.RetainDays = 7

// 自动压缩：压缩3天前的日志
config.AutoCompress = true
config.CompressCycle = 2 * time.Hour
config.CompressExclude = true
```

---

## 性能对比

假设 aiproxy 需要记录大量日志：

| 日志库 | 同步/异步 | 平均延迟 | 吞吐量 | 推荐度 |
|--------|----------|----------|--------|--------|
| 标准库 log | 同步 | ~10-50µs | ~20-100K/s | ⭐⭐ |
| logrus | 同步 | ~10-50µs | ~20-100K/s | ⭐⭐ |
| zap | 异步 | ~1-5µs | ~500K-1M/s | ⭐⭐⭐⭐ |
| **xlog** | **异步** | **~1.3µs** | **~876K/s** | ⭐⭐⭐⭐⭐ |

**xlog 优势**：
- ✅ 极低延迟（1.3微秒）
- ✅ 高吞吐量（87万/秒）
- ✅ 自动清理和压缩
- ✅ 简单易用
- ✅ 完全自主可控

---

## 集成检查清单

- [ ] 选择引用方式（本地/Git/vendor）
- [ ] 修改 go.mod
- [ ] 创建日志配置
- [ ] 初始化 logger
- [ ] 替换现有 log 调用
- [ ] 运行 `go mod tidy`
- [ ] 测试日志输出
- [ ] 配置生产环境参数

---

## 故障排查

### 问题1: 找不到包
```
cannot find package "github.com/chen1994/xlog"
```

**解决**:
```bash
go mod tidy
go mod download
```

### 问题2: replace 不生效
```
module github.com/chen1994/xlog: reading at revision v0.0.0: unknown revision
```

**解决**: 确保 replace 指令在 go.mod 的**末尾**

### 问题3: 日志文件不创建
**解决**: 检查 LogDir 目录权限，确保程序有写入权限

---

**最后更新**: 2026-03-15
