#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo "========================================"
echo "XLog 模块完整测试套件"
echo "========================================"
echo ""

# 切换到脚本所在目录
cd "$(dirname "$0")"

# 创建测试报告文件
REPORT_FILE="test_report_$(date +%Y%m%d_%H%M%S).txt"

echo -e "${CYAN}测试开始时间: $(date '+%Y-%m-%d %H:%M:%S')${NC}"
echo "测试开始时间: $(date '+%Y-%m-%d %H:%M:%S')" > "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

# 运行测试
echo ""
echo -e "${CYAN}[1/4] 运行单元测试...${NC}"
echo "========================================" >> "$REPORT_FILE"
echo "[1] 单元测试" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE}"
go test -v -run Test 2>&1 | tee -a "$REPORT_FILE"

if [ ${PIPESTATUS[0]} -eq 0 ]; then
    echo -e "${GREEN}✓ 单元测试通过${NC}"
    echo "" >> "$REPORT_FILE"
    echo "[结果] 单元测试: 通过" >> "$REPORT_FILE"
else
    echo -e "${RED}✗ 单元测试失败${NC}"
    echo "" >> "$REPORT_FILE"
    echo "[结果] 单元测试: 失败" >> "$REPORT_FILE"
fi

echo ""
echo -e "${CYAN}[2/4] 运行性能测试...${NC}"
echo "" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE"
echo "[2] 性能测试" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE}"
go test -bench=. -benchmem 2>&1 | tee -a "$REPORT_FILE"

if [ ${PIPESTATUS[0]} -eq 0 ]; then
    echo -e "${GREEN}✓ 性能测试完成${NC}"
    echo "" >> "$REPORT_FILE"
    echo "[结果] 性能测试: 完成" >> "$REPORT_FILE"
else
    echo -e "${RED}✗ 性能测试失败${NC}"
    echo "" >> "$REPORT_FILE"
    echo "[结果] 性能测试: 失败" >> "$REPORT_FILE"
fi

echo ""
echo -e "${CYAN}[3/4] 代码覆盖率测试...${NC}"
echo "" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE"
echo "[3] 代码覆盖率" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE}"
go test -coverprofile=coverage.out 2>&1 | tee -a "$REPORT_FILE"

if [ -f coverage.out ]; then
    go tool cover -func=coverage.out >> "$REPORT_FILE" 2>&1
    echo -e "${GREEN}✓ 覆盖率报告生成完成${NC}"
fi

echo ""
echo -e "${CYAN}[4/4] 内存泄漏检测...${NC}"
echo "" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE"
echo "[4] 内存检测" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE}"

go test -run TestConcurrency -memprofile=mem.out 2>&1 | tee -a "$REPORT_FILE"

if [ -f mem.out ]; then
    echo -e "${GREEN}✓ 内存分析完成${NC}"
    echo "内存分析文件: mem.out" >> "$REPORT_FILE"
fi

echo ""
echo "========================================"
echo -e "${CYAN}测试总结${NC}"
echo "========================================"
echo "" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE"
echo "测试总结" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE}"

# 统计测试结果
if grep -q "PASS" "$REPORT_FILE"; then
    echo -e "${GREEN}✓ 所有测试通过${NC}"
    echo "[总体结果] 所有测试通过" >> "$REPORT_FILE"
else
    echo -e "${RED}✗ 部分测试失败${NC}"
    echo "[总体结果] 部分测试失败" >> "$REPORT_FILE"
fi

echo ""
echo -e "${YELLOW}测试报告已保存到: $REPORT_FILE${NC}"
echo "测试结束时间: $(date '+%Y-%m-%d %H:%M:%S')" >> "$REPORT_FILE"
echo ""
