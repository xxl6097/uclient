package main

import (
	"fmt"
	"time"
)

func main() {
	//2025-07-12 18:56:57
	//timestamp := 1752346617
	//fmt.Println(timestamp)
	//fmt.Println(u.TimestampFormat(int64(timestamp)))

	// 秒级时间戳 → 格式化时间（北京时间）
	timestamp := int64(1750341355)
	utcTime := time.Unix(timestamp, 0)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	beijingTime := utcTime.In(loc)

	formatted := beijingTime.Format("2006-01-02 15:04:05")
	fmt.Println(formatted) // "2025-07-12 21:10:15"
}
