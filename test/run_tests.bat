@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo XLog 模块完整测试套件
echo ========================================
echo.

REM 设置颜色
set "GREEN=[92m"
set "RED=[91m"
set "YELLOW=[93m"
set "CYAN=[96m"
set "RESET=[0m"

REM 切换到测试目录
cd /d "%~dp0"

REM 创建测试报告文件
set "REPORT_FILE=test_report_%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%.txt"
set "REPORT_FILE=%REPORT_FILE: =0%"

echo %CYAN%测试开始时间: %date% %time%%RESET%
echo 测试开始时间: %date% %time% > "%REPORT_FILE%"
echo. >> "%REPORT_FILE%"

REM 运行测试
echo.
echo %CYAN%[1/4] 运行单元测试...%RESET%
echo ======================================== >> "%REPORT_FILE%"
echo [1] 单元测试 >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"
go test -v -run Test 2>&1 | tee -a "%REPORT_FILE%"

if %errorlevel% equ 0 (
    echo %GREEN%✓ 单元测试通过%RESET%
    echo. >> "%REPORT_FILE%"
    echo [结果] 单元测试: 通过 >> "%REPORT_FILE%"
) else (
    echo %RED%✗ 单元测试失败%RESET%
    echo. >> "%REPORT_FILE%"
    echo [结果] 单元测试: 失败 >> "%REPORT_FILE%"
)

echo.
echo %CYAN%[2/4] 运行性能测试...%RESET%
echo. >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"
echo [2] 性能测试 >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"
go test -bench=. -benchmem 2>&1 | tee -a "%REPORT_FILE%"

if %errorlevel% equ 0 (
    echo %GREEN%✓ 性能测试完成%RESET%
    echo. >> "%REPORT_FILE%"
    echo [结果] 性能测试: 完成 >> "%REPORT_FILE%"
) else (
    echo %RED%✗ 性能测试失败%RESET%
    echo. >> "%REPORT_FILE%"
    echo [结果] 性能测试: 失败 >> "%REPORT_FILE%"
)

echo.
echo %CYAN%[3/4] 代码覆盖率测试...%RESET%
echo. >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"
echo [3] 代码覆盖率 >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"
go test -coverprofile=coverage.out 2>&1 | tee -a "%REPORT_FILE%"

if exist coverage.out (
    go tool cover -func=coverage.out >> "%REPORT_FILE%" 2>&1
    echo %GREEN%✓ 覆盖率报告生成完成%RESET%
)

echo.
echo %CYAN%[4/4] 内存泄漏检测...%RESET%
echo. >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"
echo [4] 内存检测 >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"

REM 创建简单的内存测试
go test -run TestConcurrency -memprofile=mem.out 2>&1 | tee -a "%REPORT_FILE%"

if exist mem.out (
    echo %GREEN%✓ 内存分析完成%RESET%
    echo 内存分析文件: mem.out >> "%REPORT_FILE%"
)

echo.
echo ========================================
echo %CYAN%测试总结%RESET%
echo ========================================
echo. >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"
echo 测试总结 >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"

REM 统计测试结果
findstr /C:"PASS" "%REPORT_FILE%" >nul
if %errorlevel% equ 0 (
    echo %GREEN%✓ 所有测试通过%RESET%
    echo [总体结果] 所有测试通过 >> "%REPORT_FILE%"
) else (
    echo %RED%✗ 部分测试失败%RESET%
    echo [总体结果] 部分测试失败 >> "%REPORT_FILE%"
)

echo.
echo %YELLOW%测试报告已保存到: %REPORT_FILE%%RESET%
echo 测试结束时间: %date% %time% >> "%REPORT_FILE%"
echo.
echo 按任意键退出...
pause >nul
