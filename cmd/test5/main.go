package main

import (
	"fmt"
	"regexp"
)

// 定义日志结构体
type KernelLog struct {
	Timestamp    string // 时间戳
	LogLevel     string // 日志级别
	KernelTime   string // 内核时间
	Location     string // 位置标识
	FunctionName string // 函数名
	LineNumber   string // 代码行号
	Operation    string // 操作类型
	MACAddress   string // MAC地址
}

func main() {
	logStr := "Mon Jul 28 18:12:45 2025 kern.warn kernel: [34592.938265] 7981@C13L2,MacTableDeleteEntry() 1921: Del Sta:ee:af:48:c9:e6:c1"
	logStr = "Mon Jul 28 18:15:11 2025 kern.warn kernel: [34739.452859] 7981@C13L2,MacTableInsertEntry() 1559: New Sta:ee:af:48:c9:e6:c1"
	//MacTableDeleteEntry() 1921: Del Sta:ee:af:48:c9:e6:c1

	// 核心正则表达式（匹配7个关键组）
	re := regexp.MustCompile(
		`^(\w{3} \w{3} \d{2} \d{2}:\d{2}:\d{2} \d{4}) ` + // 时间戳
			`(\w+\.\w+) .+? ` + // 日志级别
			`\[([\d\.]+)\] ` + // 内核时间
			`(\d+@\w+\d+),` + // 位置标识
			`(\w+)\(\) ` + // 函数名
			`(\d+): ` + // 行号
			`(\w+ \w+):` + // 操作描述
			`([\w:]+)$`, // MAC地址
	)

	matches := re.FindStringSubmatch(logStr)
	if len(matches) < 8 {
		fmt.Println("日志格式不匹配")
		return
	}

	// 填充结构体
	logEntry := KernelLog{
		Timestamp:    matches[1],
		LogLevel:     matches[2],
		KernelTime:   matches[3],
		Location:     matches[4],
		FunctionName: matches[5],
		LineNumber:   matches[6],
		Operation:    matches[7],
		MACAddress:   matches[8],
	}

	// 打印解析结果
	fmt.Printf("时间: %s\n函数: %s (行号: %s)\n操作: %s\nMAC地址: %s\n",
		logEntry.Timestamp,
		logEntry.FunctionName,
		logEntry.LineNumber,
		logEntry.Operation,
		logEntry.MACAddress,
	)
}
