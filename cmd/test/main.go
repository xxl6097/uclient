package main

import (
	"fmt"
	"regexp"
	"time"
)

func autoParse(timeStr string) (time.Time, error) {
	var layouts = []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			return t, nil // 解析成功
		}
	}
	return time.Time{}, fmt.Errorf("无法识别的格式")
}
func parseTimer(logLine string) (*time.Time, error) {
	re := regexp.MustCompile(`^(\w+\s+\w+\s+\d+\s+\d+:\d+:\d+\s+\d+)`)
	matches := re.FindStringSubmatch(logLine)
	if len(matches) > 1 {
		timeStr := matches[1]
		t, err := autoParse(timeStr)
		return &t, err
	}
	return nil, nil
}

func parseTime1(logLine string) string {
	t, err := parseTimer(logLine)
	if err == nil {
		return t.Format(time.DateTime)
	}
	return ""
}

func main() {
	logLine := "Wed Jul  9 14:57:55 2025 daemon.notice hostapd: phy1-ap0: AP-STA-DISCONNECTED 7a:34:62:d5:a4:18"
	logLine = "Fri Jul 11 21:26:42 2025 daemon.notice hostapd: phy1-ap0: AP-STA-DISCONNECTED 5a:a7:22:62:3d:26"
	t, _ := parseTimer(logLine)
	fmt.Println(t.UnixMicro())
	fmt.Println(t.UnixMilli())
	fmt.Println(t.UnixNano())
	fmt.Println(t.Unix())
	fmt.Println(time.Unix(t.Unix(), 0).Format("2006-01-02 15:04:05"))
	fmt.Println("----", parseTime1(logLine))

	// 编译正则表达式 [3,5](@ref)
	re := regexp.MustCompile(`hostapd:\s+(phy[\w-]+?):`)

	// 提取匹配结果
	matches := re.FindStringSubmatch(logLine)
	if len(matches) < 2 {
		fmt.Println("未找到匹配字段")
		return
	}
	phyField := matches[1]             // 捕获组索引为1
	fmt.Println("解析结果:", phyField) // 输出: phy1-ap0

	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	a = append(a, 1, 2, 3)
	size := len(a)
	tempSize := size - 3
	fmt.Println(a[tempSize:])
}
