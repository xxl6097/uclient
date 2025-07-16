package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Data struct {
	Timestamp int64 `json:"timestamp"`
	Connected bool  `json:"connected"`
}
type TimeData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []Data `json:"data"`
}

var testJson = "{\n    \"code\": 0,\n    \"msg\": \"获取成功\",\n    \"data\": [\n        {\n            \"timestamp\": 1752581874016,\n            \"connected\": true\n        },\n        {\n            \"timestamp\": 1752564867862,\n            \"connected\": false\n        },\n        {\n            \"timestamp\": 1752563258374,\n            \"connected\": true\n        },\n        {\n            \"timestamp\": 1752539186800,\n            \"connected\": true\n        },\n        {\n            \"timestamp\": 1752539185268,\n            \"connected\": false\n        },\n        {\n            \"timestamp\": 1752537627469,\n            \"connected\": true\n        }\n    ]\n}"

func GroupTimestampsByDay(timestamps []Data) map[time.Time][]int64 {
	// 初始化分组Map
	groups := make(map[time.Time][]int64)

	for _, ts := range timestamps {
		// 将时间戳转为time.Time类型
		t := time.UnixMilli(ts.Timestamp)
		// 归一化到当天的0点（去掉时分秒）
		day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		// 将时间戳添加到对应天的分组中
		groups[day] = append(groups[day], ts.Timestamp)
	}
	return groups
}
func main() {
	var timedata TimeData
	json.Unmarshal([]byte(testJson), &timedata)
	//fmt.Println("clife", timedata.Data)
	// 按天分组
	grouped := GroupTimestampsByDay(timedata.Data)

	// 打印结果
	for day, stamps := range grouped {
		fmt.Printf("======> %s 周%d\n", day.Format(time.DateOnly), day.Weekday())
		for _, stamp := range stamps {
			fmt.Printf("%v\n", time.UnixMilli(stamp).Format(time.DateTime))
		}
		a := stamps[0]
		b := stamps[len(stamps)-1]
		t1 := time.UnixMilli(b)
		t2 := time.UnixMilli(a)
		onWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), 9, 00, 0, 0, t1.Location())
		offWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), 18, 30, 0, 0, t1.Location())
		onWorkOverTime := onWorkTime.Sub(t1)
		offWorkOverTime := t2.Sub(offWorkTime)
		fmt.Printf("上班： %v\n", t1.Format(time.TimeOnly))
		fmt.Printf("下班： %v\n\n", t2.Format(time.TimeOnly))
		fmt.Printf("上班加班时长： %v\n", onWorkOverTime)
		fmt.Printf("下班加班时长： %v\n", offWorkOverTime)
		duration := onWorkOverTime
		if duration < 0 {
			duration = 0
		}
		if offWorkOverTime > 0 {
			duration += offWorkOverTime
		}
		fmt.Printf("加班时长： %v\n\n", duration)
		abs := duration.Seconds() / time.Hour.Seconds()
		fmt.Printf("加班时长： %.1f\n", abs)
	}
}
