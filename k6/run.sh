#!/bin/bash

# K6 负载测试运行脚本

BASE_URL=${BASE_URL:-"http://localhost:8080"}
TEST_TYPE=${TEST_TYPE:-"load"}

echo "=========================================="
echo "K6 负载测试"
echo "=========================================="
echo "目标服务器: $BASE_URL"
echo "测试类型: $TEST_TYPE"
echo "=========================================="
echo ""

# 检查 k6 是否安装
if ! command -v k6 &> /dev/null; then
    echo "错误: k6 未安装"
    echo "请访问 https://k6.io/docs/getting-started/installation/ 安装 k6"
    exit 1
fi

# 根据测试类型选择脚本
case $TEST_TYPE in
    smoke)
        echo "运行冒烟测试..."
        BASE_URL=$BASE_URL k6 run smoke_test.js
        ;;
    stress)
        echo "运行压力测试..."
        BASE_URL=$BASE_URL k6 run stress_test.js
        ;;
    load|*)
        echo "运行负载测试..."
        BASE_URL=$BASE_URL k6 run load_test.js
        ;;
esac

echo ""
echo "=========================================="
echo "测试完成"
echo "=========================================="

