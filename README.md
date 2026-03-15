# XLog - Go 异步日志库

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

XLog 是一个高性能的 Go 语言异步日志库，从 C++ 的 [XSimpleLogEx](https://github.com/chen1994/xlog) 移植而来。

> **注意**: 本项目从 C++ 版本的 XSimpleLogEx 移植，保留了核心功能并针对 Go 语言特性进行了优化。

## 特性

- **异步日志写入**: 使用 channel 和后台 goroutine 实现异步写入，避免阻塞主线程
- **自动清理**: 支持按保留天数自动清理过期日志文件
- **自动压缩**: 支持定期将旧日志文件压缩为 .zip 格式，节省磁盘空间
- **日志轮转**: 支持按日期自动创建新的日志文件
- **灵活配置**: 提供丰富的配置选项，满足不同场景需求

## 安装

```bash
go get github.com/chen1994/xlog
```

## 快速开始

### 基础用法

```go
package main

import (
    "github.com/chen1994/xlog"
)

func main() {
    // 使用默认配置
    config := xlog.DefaultConfig()

    // 创建日志记录器
    logger, err := xlog.NewXLogger(config)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    // 记录日志
    logger.Log("这是一条日志消息")
    logger.Log("带参数的日志: %d, %s", 123, "test")
}
```

### 自定义配置

```go
config := xlog.DefaultConfig()

// 基础配置
config.LogDir = "./logs"           // 日志目录
config.LogBaseName = "myapp"       // 日志文件名前缀
config.LogSuffix = "log"           // 日志文件扩展名
config.HasDate = true              // 文件名包含日期
config.HasProcessID = false        // 文件名不包含进程ID
config.UTF8Format = true           // 使用 UTF-8 编码

// 异步写入配置
config.WriteInterval = 5 * time.Second  // 写入间隔
config.BufferSize = 10000               // 缓冲区大小

// 自动清理配置
config.AutoCleanup = true               // 启用自动清理
config.CleanupCycle = 6 * time.Minute   // 清理检查周期
config.RetainDays = 3                   // 保留最近3天的日志

// 自动压缩配置
config.AutoCompress = true              // 启用自动压缩
config.CompressCycle = 10 * time.Minute // 压缩检查周期
config.CompressExclude = true           // 排除当前日志文件

logger, err := xlog.NewXLogger(config)
```

### 日志级别

```go
const (
    LevelError = 0x01  // 错误
    LevelWarn  = 0x02  // 警告
    LevelInfo  = 0x04  // 信息
    LevelDebug = 0x08  // 调试
)

// 使用不同级别记录日志
logger.LogEx(xlog.LevelInfo, "信息日志")
logger.LogEx(xlog.LevelWarn, "警告日志")
logger.LogEx(xlog.LevelError, "错误日志")
logger.LogEx(xlog.LevelDebug, "调试日志")

// Log() 方法默认使用 LevelInfo
logger.Log("普通日志")
```

## 配置说明

### Config 结构体

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| LogDir | string | 可执行文件所在目录 | 日志文件目录 |
| LogBaseName | string | 可执行文件名 | 日志文件基本名称 |
| LogSuffix | string | "log" | 日志文件扩展名 |
| HasDate | bool | true | 文件名是否包含日期（格式：YYYYMMDD） |
| HasProcessID | bool | false | 文件名是否包含进程ID |
| UTF8Format | bool | true | 是否使用 UTF-8 编码（false 则使用 UTF-16） |
| WriteInterval | time.Duration | 5秒 | 日志写入间隔 |
| BufferSize | int | 10000 | 日志缓冲区大小 |
| AutoCleanup | bool | false | 是否启用自动清理 |
| CleanupCycle | time.Duration | 6分钟 | 清理检查周期 |
| RetainDays | int | 3 | 日志保留天数 |
| AutoCompress | bool | false | 是否启用自动压缩 |
| CompressCycle | time.Duration | 10分钟 | 压缩检查周期 |
| CompressExclude | bool | true | 是否排除当前日志文件 |

## 文件命名规则

日志文件名格式：`{LogBaseName}.{YYYYMMDD}.{ProcessID}.{LogSuffix}`

示例：
- `myapp.20260315.log` （HasDate=true, HasProcessID=false）
- `myapp.20260315.12345.log` （HasDate=true, HasProcessID=true）
- `myapp.log` （HasDate=false, HasProcessID=false）

## 工作原理

### 异步写入机制

1. 日志调用将消息发送到 channel（非阻塞）
2. 后台 goroutine 定期从 channel 读取日志
3. 批量写入到文件（减少磁盘 I/O）

### 自动清理机制

1. 每隔 `CleanupCycle` 检查一次日志文件
2. 删除修改时间超过 `RetainDays` 的日志文件
3. 支持按磁盘空间清理（`CleanupBySize` 方法）

### 自动压缩机制

1. 每隔 `CompressCycle` 检查一次日志文件
2. 将符合条件的老日志压缩为 .zip 格式
3. 压缩成功后删除原文件
4. 支持排除当前正在使用的日志文件

## 性能特点

- **非阻塞**: 日志写入使用 channel，主线程不会阻塞
- **批量写入**: 定期批量刷新，减少磁盘 I/O 次数
- **并发安全**: 使用 mutex 保护文件句柄操作
- **资源友好**: 支持自动清理和压缩，避免磁盘空间耗尽

## 注意事项

1. 调用 `logger.Close()` 确保所有日志都已写入文件
2. 如果 channel 满了，新的日志会被丢弃（避免阻塞）
3. 压缩功能需要足够的磁盘空间临时存储 .zip 文件
4. 清理和压缩操作在后台 goroutine 中执行

## 从 C++ XSimpleLogEx 迁移

主要功能对照：

| C++ 功能 | Go 功能 | 状态 |
|----------|---------|------|
| 异步日志写入 | ✓ | 已实现 |
| 自动清理 | ✓ | 已实现 |
| 自动压缩 | ✓ | 已实现 |
| UTF-16 编码 | ✓ | 已实现（默认 UTF-8） |
| 日志级别 | ✓ | 已实现 |
| 日期轮转 | ✓ | 已实现 |

## 许可证

MIT License

## 作者

从 C++ XSimpleLogEx 移植到 Go
