#!/bin/bash

# 配置区（用户可修改）
TARGET="www.baidu.com"       # 检测目标
LOG_FILE="network_monitor.log" # 日志路径
INTERVAL=1                   # 检测间隔(秒)

# 颜色定义
RED="🔴"
GREEN="🟢"
YELLOW=""
RESET=""

# 网络检测函数
check_network() {
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")

  # Ping测试
  ping -c1 -W1 $TARGET &>/dev/null
  if [ $? -eq 0 ]; then
    local ping_status="${GREEN}UP${RESET}"
    local ping_time=$(ping -c1 $TARGET | awk -F'/' '/min/ {print $5}')
  else
    ping_status="${RED}DOWN${RESET}"
    ping_time="N/A"
  fi

  # HTTP检测
  http_code=$(curl -s -o /dev/null -w "%{http_code}" $TARGET)
  [ "$http_code" = "200" ] && http_status="${GREEN}$http_code${RESET}" || http_status="${RED}$http_code${RESET}"

  # 输出结果
  printf "%-19s | Ping: %-7s (%-5sms) | HTTP: %-5s\n" \
    "$timestamp" "$ping_status" "${YELLOW}$ping_time${RESET}" "$http_status"

  # 记录日志
  echo "$timestamp,Ping:$ping_status,Time:${ping_time}ms,HTTP:$http_code" >> $LOG_FILE
}

# 主循环
echo -e "${YELLOW}开始网络监控 (目标: $TARGET)${RESET}"
echo "按 Ctrl+C 停止..."
while true; do
  check_network
  sleep $INTERVAL
done