package main

import (
	"fmt"
	"github.com/xxl6097/uclient/internal/u"
	"time"
)

// Sat Jul 12 21:24:08 2025 daemon.notice hostapd: phy1-ap0: AP-STA-CONNECTED 28:f0:76:39:85:88 auth_alg=open
func main() {
	//2025-07-12 18:56:57
	//timestamp := 1752346617
	//fmt.Println(timestamp)
	//fmt.Println(u.TimestampFormat(int64(timestamp)))

	// 秒级时间戳 → 格式化时间（北京时间）
	//timestamp := int64(1750341355)
	//utcTime := time.Unix(timestamp, 0)
	//loc, _ := time.LoadLocation("Asia/Shanghai")
	//beijingTime := utcTime.In(loc)
	//
	//formatted := beijingTime.Format("2006-01-02 15:04:05")
	//fmt.Println(formatted) // "2025-07-12 21:10:15"

	//timeStr := "Sat Jul 12 21:24:08 2025"
	//t, err := time.Parse(time.ANSIC, timeStr)
	//if err == nil {
	//	fmt.Println(t.Format(time.DateTime))
	//	timestamp := t.Unix()
	//	// 秒级时间戳 → 格式化时间（北京时间）
	//	utcTime := time.Unix(timestamp, 0)
	//	loc, _ := time.LoadLocation("Asia/Shanghai")
	//	beijingTime := utcTime.In(loc)
	//
	//	formatted := beijingTime.Format("2006-01-02 15:04:05")
	//	fmt.Println(formatted) // "2025-07-12 21:10:15"
	//}

	//timeStr := "Sat Jul 12 21:24:08 2025"
	////t, err := time.Parse(time.ANSIC, timeStr) // 解析为 UTC 时间
	//loc, _ := time.LoadLocation("Asia/Shanghai")
	//t, err := time.ParseInLocation(time.ANSIC, timeStr, loc) // 按北京时间解析
	//if err == nil {
	//	fmt.Println(t.Format(time.DateTime)) // 输出：2025-07-12 21:24:08（UTC）
	//	timestamp := t.Unix()
	//	utcTime := time.Unix(timestamp, 0) // 仍是 UTC 时间
	//	loc, _ := time.LoadLocation("Asia/Shanghai")
	//	beijingTime := utcTime.In(loc)                         // UTC→UTC+8
	//	fmt.Println(beijingTime.Format("2006-01-02 15:04:05")) // 输出：2025-07-13 05:24:08
	//}

	fmt.Println(u.TimestampFormat(1752553006670))

	t1 := time.Now().UnixMilli()
	fmt.Println("t1", u.TimestampFormat(t1))
}
