package main

import (
	"encoding/json"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/u"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func tee1() {
	tempFilePath := filepath.Join("/Users/uuxia/Downloads/192.168.1.1/202507151912", "5a:a7:22:62:3d:26")
	byteArray, err := os.ReadFile(tempFilePath)
	if err != nil {
		return
	}
	var cfg []*openwrt.Status
	err = ukey.GobToStruct(byteArray, &cfg)
	if err != nil {
		return
	}
	for _, status := range cfg {
		fmt.Println(u.TimestampToDateTime(status.Timestamp))
	}
	fmt.Println("header", cfg[0].Timestamp, u.TimestampToDateTime(cfg[0].Timestamp))
	fmt.Println("tailer", cfg[len(cfg)-1].Timestamp, u.TimestampToDateTime(cfg[len(cfg)-1].Timestamp))
}

func tee4() {
	//workDir := "/Users/uuxia/Downloads/192.168.1.1/202507161649"
	//mac := "5a:a7:22:62:3d:26"
	//tempFilePath := filepath.Join(workDir, mac)
	//data := openwrt.ReadWorkTimeByMac(tempFilePath)
	//for k, status := range data {
	//	fmt.Printf("%v %+v %+v", k, status[0], status[1])
	//}
	//openwrt.UpdatetWorkTime(mac, workDir, openwrt.Status{
	//	Timestamp: time.Now().UnixMilli(),
	//})
}

func getStatusByMac(mac string) []*openwrt.Status {
	statusDir := "/Users/uuxia/Downloads/192.168.1.1/202507151912"
	_ = u.CheckDirector(statusDir)
	tempFilePath := filepath.Join(statusDir, mac)
	byteArray, err := os.ReadFile(tempFilePath)
	if err != nil {
		return nil
	}
	var cfg []*openwrt.Status
	err = ukey.GobToStruct(byteArray, &cfg)
	if err != nil {
		return nil
	}
	//var month string
	//for _, status := range cfg {
	//	tempMonth := u.TimestampFormatToMonth(status.Timestamp)
	//	if month != tempMonth {
	//		month = tempMonth
	//		fmt.Println("\n==========================================>", tempMonth)
	//	}
	//	fmt.Println(u.TimestampToDateTime(status.Timestamp))
	//}
	//fmt.Println("header", cfg[0].Timestamp, u.TimestampToDateTime(cfg[0].Timestamp))
	//fmt.Println("tailer", cfg[len(cfg)-1].Timestamp, u.TimestampToDateTime(cfg[len(cfg)-1].Timestamp))
	return cfg
}

func setStatusByMac(mac string, statusList []*openwrt.Status) error {
	if mac == "" {
		return nil
	}
	if statusList == nil || len(statusList) == 0 {
		return nil
	}

	content, err := ukey.StructToGob(statusList)
	if err != nil {
		return err
	}
	statusDir := "/Users/uuxia/Downloads/192.168.1.1/202507151912"
	_ = u.CheckDirector(statusDir)
	tempFilePath := filepath.Join(statusDir, mac)
	file, err := os.Create(tempFilePath) // 文件不存在则创建，存在则截断
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入内容
	_, err = file.Write(content)
	return err
}

func setWorkByMac(mac string, data any) error {
	if mac == "" {
		return nil
	}
	if data == nil {
		return nil
	}

	//content, err := json.MarshalIndent(data, "", " ")
	content, err := ukey.StructToGob(data)
	if err != nil {
		return err
	}
	statusDir := "/Users/uuxia/Downloads/192.168.1.1"
	_ = u.CheckDirector(statusDir)
	tempFilePath := filepath.Join(statusDir, mac)
	file, err := os.Create(tempFilePath) // 文件不存在则创建，存在则截断
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入内容
	_, err = file.Write(content)
	return err
}

func tee2() {
	cfg := make([]*openwrt.Status, 0)
	now := time.Date(2024, 12, 1, 0, 0, 0, 0, u.GetLocation())
	for i := 0; i < 7; i++ {
		// 每月1号 00:00:00
		for j := 0; j < 31; j++ {
			//t := time.Date(now.Year(), now.Month()+time.Month(i+1), now.Day()+j, 0, 0, 0, 0, now.Location())
			//fmt.Printf("第%d个月首日: %s (时间戳: %d)\n", i+1, t.Format(time.DateTime), t.Unix())
			t1 := time.Date(now.Year(), now.Month()+time.Month(i+1), now.Day()+j, 8, rand.Intn(60), rand.Intn(60), 0, now.Location())
			t2 := time.Date(now.Year(), now.Month()+time.Month(i+1), now.Day()+j, 13, rand.Intn(60), rand.Intn(60), 0, now.Location())
			t3 := time.Date(now.Year(), now.Month()+time.Month(i+1), now.Day()+j, 19, rand.Intn(60), rand.Intn(60), 0, now.Location())
			fmt.Println(t3.Format(time.DateTime))
			fmt.Println(t2.Format(time.DateTime))
			fmt.Println(t1.Format(time.DateTime))
			cfg = append(cfg, &openwrt.Status{
				Timestamp: t3.UnixMilli(),
				Connected: true,
			}, &openwrt.Status{
				Timestamp: t2.UnixMilli(),
				Connected: true,
			}, &openwrt.Status{
				Timestamp: t1.UnixMilli(),
				Connected: false,
			})
		}

	}

	fmt.Println("size", len(cfg))
	setStatusByMac("cfg", cfg)
	//getStatusByMac("cfg")

}

func tee5() {
	now := time.Date(2024, 12, 1, 0, 0, 0, 0, u.GetLocation())
	days := make(map[string]*openwrt.WorkEntry)
	for i := 0; i < 7; i++ {
		for j := 0; j < 31; j++ {
			//t := time.Date(now.Year(), now.Month()+time.Month(i+1), now.Day()+j, 0, 0, 0, 0, now.Location())
			//fmt.Printf("第%d个月首日: %s (时间戳: %d)\n", i+1, t.Format(time.DateTime), t.Unix())
			t1 := time.Date(now.Year(), now.Month()+time.Month(i+1), now.Day()+j, 8, rand.Intn(60), rand.Intn(60), 0, now.Location())
			t3 := time.Date(now.Year(), now.Month()+time.Month(i+1), now.Day()+j, 19, rand.Intn(60), rand.Intn(60), 0, now.Location())
			fmt.Println(t3.Format(time.DateTime))
			fmt.Println(t1.Format(time.DateTime))
			item := openwrt.WorkEntry{
				OnWorkTime:  t1.UnixMilli(),
				OffWorkTime: t3.UnixMilli(),
				Weekday:     int(t3.Weekday()),
				//IsWeekDay:   t1.Weekday() == time.Saturday || t1.Weekday() == time.Sunday,
			}
			day := t1.Format(time.DateOnly)
			days[day] = &item
		}
	}

	setWorkByMac("work", days)
	//getStatusByMac("cfg")

}

func groupTimestampsByDay(timestamps []*openwrt.Status, workType openwrt.WorkTypeSetting) []*openwrt.Work {
	on := u.GetTime(workType.OnWorkTime, u.GetLocation())
	off := u.GetTime(workType.OffWorkTime, u.GetLocation())
	if on == nil || off == nil {
		return nil
	}
	// 初始化分组Map
	works := make([]*openwrt.Work, 0)
	//groups := make(map[int64]map[int64][]int64)
	months := make(map[string]*openwrt.Work)
	days := make(map[string][]int64)
	var day string
	for idx, ts := range timestamps {
		if ts.Timestamp < 1_000_000_000_000 {
			ts.Timestamp *= 1000
		}
		t := time.UnixMilli(ts.Timestamp)
		month := fmt.Sprintf("%d-%02d", t.Year(), int(t.Month())) //int(t.Month()) + t.Year() //time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, t.Location())
		work, monthOk := months[month]
		if !monthOk {
			//monthMap = make(map[int64][]int64)
			work = &openwrt.Work{Month: month}
			works = append(works, work)
		}
		//fmt.Println(t.Format(time.DateTime))
		//OffWorkTime: t.Format(time.TimeOnly)
		day = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Format(time.DateOnly)
		dayList, dayOk := days[day]
		if !dayOk {
			if idx > 0 {
				worksSize := len(works)
				if monthOk {
					worksSize -= 1
				} else {
					worksSize -= 2
				}
				tempWork := works[worksSize]
				if tempWork != nil {
					tempItem := timestamps[idx-1]
					if tempItem != nil {
						//tmpTime := time.UnixMilli(tempItem.Timestamp).Format(time.TimeOnly)
						workTimeSize := len(tempWork.WorkTime)
						//fmt.Println(workTimeSize)
						//if workTimeSize == 0 {
						//	fmt.Println(workTimeSize)
						//}
						tempWorkTime := tempWork.WorkTime[workTimeSize-1]
						if tempWorkTime.Date != "" {
							tempDay := days[tempWorkTime.Date]
							t1 := time.UnixMilli(tempDay[len(tempDay)-1])
							t2 := time.UnixMilli(tempDay[0])
							tempWorkTime.WorkTime1 = t1.Format(time.TimeOnly)
							tempWorkTime.WorkTime2 = t2.Format(time.TimeOnly)
							tempWorkTime.Weekday = int(t2.Weekday())
							//weekIndex := int(t1.Weekday())
							//isWeekDay := weekIndex == 0 || weekIndex == 6
							//tempWorkTime.IsWeekDay = isWeekDay
							//if isWeekDay && !ts.IsWeekDay {
							//	tempWorkTime.IsWeekDay = !ts.IsWeekDay
							//}

							onWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), on.Hour(), on.Minute(), on.Second(), 0, t1.Location())
							offWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), off.Hour(), off.Minute(), off.Second(), 0, t1.Location())
							onWorkOverTime := onWorkTime.Sub(t1)
							offWorkOverTime := t2.Sub(offWorkTime)

							//fmt.Printf("上班： %v\n", t1.Format(time.TimeOnly))
							//fmt.Printf("下班： %v\n\n", t2.Format(time.TimeOnly))
							//fmt.Printf("上班加班时长： %v\n", onWorkOverTime)
							//fmt.Printf("下班加班时长： %v\n", offWorkOverTime)

							duration := onWorkOverTime
							if duration < 0 {
								duration = 0
							}
							if offWorkOverTime > 0 {
								duration += offWorkOverTime
							}
							tempWorkTime.OverWorkTimes = duration.String()
							tempWorkTime.OverWorkTimesDuration = duration
							tempWork.OverTimeDuration += duration
							tempWork.OverTime = tempWork.OverTimeDuration.String()
						}
						tempWork.WorkTime[workTimeSize-1] = tempWorkTime
					}
				}
				if tempWork == nil {
					//fmt.Println(workTimeSize)
				}
			}

			if work.WorkTime == nil {
				work.WorkTime = make([]openwrt.WorkTime, 0)
			}
			work.WorkTime = append(work.WorkTime, openwrt.WorkTime{
				Date:      day,
				Weekday:   int(t.Weekday()),
				WorkTime2: t.Format(time.TimeOnly),
			})

		}
		dayList = append(dayList, ts.Timestamp)
		//dayList = append(dayList, ts.Timestamp)
		//monthMap[day.UnixMilli()] = dayList
		days[day] = dayList
		months[month] = work
	}

	tempDay, dayOk := days[day]
	if dayOk {
		t1 := time.UnixMilli(tempDay[len(tempDay)-1])
		t2 := time.UnixMilli(tempDay[0])
		onWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), on.Hour(), on.Minute(), on.Second(), 0, t1.Location())
		offWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), off.Hour(), off.Minute(), off.Second(), 0, t1.Location())

		onWorkOverTime := onWorkTime.Sub(t1)
		offWorkOverTime := t2.Sub(offWorkTime)

		tempWork := works[len(works)-1]
		if tempWork != nil {
			tempWorkTime := tempWork.WorkTime[len(tempWork.WorkTime)-1]
			tempWorkTime.WorkTime1 = t1.Format(time.TimeOnly)
			tempWorkTime.WorkTime2 = t2.Format(time.TimeOnly)
			tempWorkTime.Weekday = int(t2.Weekday())
			duration := onWorkOverTime
			if duration < 0 {
				duration = 0
			}
			if offWorkOverTime > 0 {
				duration += offWorkOverTime
			}
			tempWorkTime.OverWorkTimes = duration.String()
			tempWork.OverTimeDuration += duration
			tempWork.OverTime = tempWork.OverTimeDuration.String()
			tempWork.WorkTime[len(tempWork.WorkTime)-1] = tempWorkTime
		}
	}
	return works
}

func tee3() {
	list := getStatusByMac("cfg")
	//var month string
	//for _, status := range list {
	//	tempMonth := u.TimestampFormatToMonth(status.Timestamp)
	//	if month != tempMonth {
	//		month = tempMonth
	//		fmt.Println("\n==========================================>", tempMonth)
	//	}
	//	fmt.Println(u.TimestampToDateTime(status.Timestamp))
	//}
	data := groupTimestampsByDay(list, openwrt.WorkTypeSetting{
		OnWorkTime:  "09:00:00",
		OffWorkTime: "18:30:00",
	})
	for _, w := range data {
		fmt.Printf("%+v\n", w)
	}
	jsonBytes, _ := json.Marshal(data)
	fmt.Println(string(jsonBytes))
	//for month, status := range data {
	//	fmt.Println("\n\n@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@", time.UnixMilli(month).Format("2006-01"))
	//	for day, ts := range status {
	//		fmt.Println(">>>>>>>>>", time.UnixMilli(day).Format(time.DateOnly), "<<<<<<<<")
	//		for _, t := range ts {
	//			fmt.Println("||", u.TimestampToDateTime(t), "||")
	//		}
	//	}
	//}
	//openwrt.CaculeteWork(list)
}

func tee6() {
	statusDir := "/Users/uuxia/Downloads/192.168.1.1"
	mac := "work"
	tempFilePath := filepath.Join(statusDir, mac)
	d, err := openwrt.GetWorkTimeAndCaculate(mac, tempFilePath, &openwrt.WorkTypeSetting{
		OnWorkTime:  "09:00:00",
		OffWorkTime: "18:30:00",
	})
	if err != nil {
		fmt.Println(err)
	} else {
		//jsonBytes, _ := json.MarshalIndent(d, "", " ")
		//fmt.Println(string(jsonBytes))
		for _, work := range d {
			fmt.Printf("%+v\n", work)
		}
	}
}

func tee7() {
	mac := "16:00:6f:83:35:e1"
	tempFilePath := filepath.Join("/Users/uuxia/Downloads/192.168.1.1/202508072016", mac)
	d, err := openwrt.GetWorkTimeAndCaculate(mac, tempFilePath, &openwrt.WorkTypeSetting{
		OnWorkTime:  "09:00:00",
		OffWorkTime: "18:30:00",
	})
	if err != nil {
		fmt.Println(err)
	} else {
		//jsonBytes, _ := json.MarshalIndent(d, "", " ")
		//fmt.Println(string(jsonBytes))
		for _, work := range d {
			fmt.Printf("%+v %+v\n", len(work.WorkTime), work)
		}
	}
	fmt.Println("size", len(d))
}

//	{
//		"mac": "5a:a7:22:62:3d:26",
//		"day": "2025-07-17",
//			"data": {
//			"date": "2025-07-17",
//			"weekday": 4,
//			"workTime1": "08:17:00",
//			"workTime2": "08:00:00",
//			"overWorkTimes": "1h0m0s",
//			"dayType": 0
//		}
//	}
func tee8() {
	jsonStr := "{\n    \"mac\": \"5a:a7:22:62:3d:26\",\n    \"day\": \"2025-07-17\",\n    \"data\": {\n        \"date\": \"2025-07-17\",\n        \"weekday\": 4,\n        \"workTime1\": \"08:17:00\",\n        \"workTime2\": \"19:00:00\",\n        \"overWorkTimes\": \"1h0m0s\",\n        \"dayType\": 0\n    }\n}"
	fmt.Println(jsonStr)

	jsonStr = "{\n    \"mac\": \"5a:a7:22:62:3d:26\",\n    \"day\": \"2025-07-12\",\n    \"data\": {\n        \"date\": \"2025-07-12\",\n        \"weekday\": 6,\n        \"workTime1\": \"06:28:00\",\n        \"workTime2\": \"14:42:00\",\n        \"overWorkTimes\": \"8h14m0s\",\n        \"dayType\": 2,\n        \"showSelect\": false\n    }\n}"

	workDir := "/Users/uuxia/Downloads/192.168.1.1/202507181533"
	//202507171747
	var body struct {
		Mac  string                 `json:"mac"`
		Day  string                 `json:"day"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal([]byte(jsonStr), &body)
	fmt.Println(body)

	data := body.Data
	day := body.Day
	openwrt.TestSetWorkTime(false, body.Mac, workDir, body.Day, func(tempEntry *openwrt.WorkEntry) {
		if v, ok := data["workTime1"]; ok {
			if vv, okk := v.(string); okk {
				t, err := u.AutoParse(fmt.Sprintf("%s %s", day, vv))
				if err == nil && t != nil {
					timestamp := t.UnixMilli()
					if !u.IsMillisecondTimestamp(timestamp) {
						timestamp *= 1000
					}
					tempEntry.OnWorkTime = timestamp
				}
			}
		}
		if v, ok := data["workTime2"]; ok {
			if vv, okk := v.(string); okk {
				t, err := u.AutoParse(fmt.Sprintf("%s %s", day, vv))
				if err == nil && t != nil {
					timestamp := t.UnixMilli()
					if !u.IsMillisecondTimestamp(timestamp) {
						timestamp *= 1000
					}
					tempEntry.OffWorkTime = timestamp
				}
			}
		}
		if floatVal, ok := data["weekday"].(float64); ok {
			intVal := int(floatVal) // 显式转换为 int
			tempEntry.Weekday = intVal
		} else {
			glog.Println("值非数字类型 weekday", data["weekday"])
		}
		if floatVal, ok := data["dayType"].(float64); ok {
			intVal := int(floatVal) // 显式转换为 int
			tempEntry.DayType = intVal
		} else {
			glog.Println("值非数字类型 dayType", data["dayType"])
		}
	})
}
func main() {
	//tee5()
	//tee6()
	tee7()
	//tee8()
	//day := "2025-07-17"
	//d, _ := u.AutoParse(day)
	//fmt.Println(d.Format("2006-01"))
}
