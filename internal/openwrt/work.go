package openwrt

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/uclient/internal/u"
	"os"
	"path/filepath"
	"time"
)

type WorkEntry struct {
	OnWorkTime  int64 `json:"onWorkTime"`
	OffWorkTime int64 `json:"offWorkTime"`
	Weekday     int   `json:"weekday"`
	IsWeekDay   bool  `json:"isWeekDay"`
}

type WorkType struct {
	OnWorkTime  string `json:"onWorkTime"`
	OffWorkTime string `json:"offWorkTime"`
}
type WorkTime struct {
	Date                  string        `json:"date"`
	Weekday               int           `json:"weekday"`
	WorkTime1             string        `json:"workTime1"`
	WorkTime2             string        `json:"workTime2"`
	OverWorkTimes         string        `json:"overWorkTimes"`
	IsWeekDay             bool          `json:"isWeekDay"`
	OverWorkTimesDuration time.Duration `json:"-"`
}

type Work struct {
	Month            string        `json:"month"`
	OverTime         string        `json:"overtime"`
	WorkTime         []WorkTime    `json:"workTime"`
	OverTimeDuration time.Duration `json:"-"`
}

func groupTimestampsByDay(timestamps []*Status) map[time.Time][]int64 {
	// 初始化分组Map
	groups := make(map[time.Time][]int64)

	for _, ts := range timestamps {
		// 将时间戳转为time.Time类型
		if ts.Timestamp < 1_000_000_000_000 {
			ts.Timestamp *= 1000
		}
		t := time.UnixMilli(ts.Timestamp)
		fmt.Println(t.Format("2006-01-02 15:04:05"))
		// 归一化到当天的0点（去掉时分秒）
		day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		// 将时间戳添加到对应天的分组中
		groups[day] = append(groups[day], ts.Timestamp)
	}
	return groups
}

func CaculeteWork(timestamps []*Status) {
	if timestamps == nil || len(timestamps) == 0 {
		return
	}
	grouped := groupTimestampsByDay(timestamps)
	// 打印结果
	for day, stamps := range grouped {
		fmt.Printf("\n======> %s 周%d\n", day.Format(time.DateOnly), day.Weekday())
		//for _, stamp := range stamps {
		//	fmt.Printf("%v\n", time.UnixMilli(stamp).Format(time.DateTime))
		//}
		a := stamps[0]
		b := stamps[len(stamps)-1]
		t1 := time.UnixMilli(b)
		t2 := time.UnixMilli(a)
		onWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), 9, 30, 0, 0, t1.Location())
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
		abs := duration.Seconds() / time.Hour.Seconds()
		fmt.Printf("加班时长： %s【%.1f】\n", duration, abs)
	}
}

func ReadWorkTimeByMac(tempFilePath string) map[string]*WorkEntry {
	if tempFilePath == "" {
		return nil
	}
	byteArray, err := os.ReadFile(tempFilePath)
	if err != nil {
		return nil
	}
	var cfg map[string]*WorkEntry
	err = ukey.GobToStruct(byteArray, &cfg)
	if err != nil {
		return nil
	}
	return cfg
}

func GetWorkTime(mac, tempFilePath string, workType *WorkType) ([]*Work, error) {
	if mac == "" {
		return nil, fmt.Errorf("mac is empty")
	}
	if tempFilePath == "" {
		return nil, fmt.Errorf("tempFilePath is empty")
	}
	if workType == nil {
		return nil, fmt.Errorf("workType is empty")
	}
	on := u.GetTime(workType.OnWorkTime, u.GetLocation())
	if on == nil {
		return nil, fmt.Errorf("on work time is nill")
	}
	off := u.GetTime(workType.OffWorkTime, u.GetLocation())
	if off == nil {
		return nil, fmt.Errorf("off work time is nill")
	}
	works := ReadWorkTimeByMac(tempFilePath)
	if works == nil {
		return nil, fmt.Errorf("works is empty")
	}

	result := make([]*Work, 0)
	months := make(map[string]*Work)
	//offWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), off.Hour(), off.Minute(), off.Second(), off.Nanosecond(), t1.Location())
	for day, w := range works {
		time1 := time.UnixMilli(w.OnWorkTime)
		time2 := time.UnixMilli(w.OffWorkTime)
		workTime1 := time1.Format(time.TimeOnly)
		workTime2 := time2.Format(time.TimeOnly)

		onWorkTime := time.Date(time1.Year(), time1.Month(), time1.Day(), on.Hour(), on.Minute(), on.Second(), 0, time1.Location())
		offWorkTime := time.Date(time2.Year(), time2.Month(), time2.Day(), off.Hour(), off.Minute(), off.Second(), 0, time2.Location())
		onWorkOverTime := onWorkTime.Sub(time1)
		offWorkOverTime := time2.Sub(offWorkTime)
		duration := onWorkOverTime
		if duration < 0 {
			duration = 0
		}
		if offWorkOverTime > 0 {
			duration += offWorkOverTime
		}
		wrokTimeTemp := WorkTime{
			Date:                  day,
			IsWeekDay:             w.IsWeekDay,
			Weekday:               w.Weekday,
			WorkTime1:             workTime1,
			WorkTime2:             workTime2,
			OverWorkTimes:         duration.String(),
			OverWorkTimesDuration: duration,
		}

		month := fmt.Sprintf("%d-%02d", time1.Year(), int(time1.Month()))
		work, monthOk := months[month]
		if !monthOk {
			work = &Work{
				Month: month,
			}
			result = append(result, work)
		}

		if work.WorkTime == nil {
			work.WorkTime = make([]WorkTime, 0)
		}
		work.WorkTime = append(work.WorkTime, wrokTimeTemp)

		months[month] = work

	}
	return result, nil
}

func getWorkTime(mac string, workType *WorkType) ([]*Work, error) {
	tempFilePath := filepath.Join(workDir, mac)
	glog.Debug("GetWorkTime", mac)
	return GetWorkTime(mac, tempFilePath, workType)
}

func setWorkTime(mac, workDir, day string, fn func(*WorkEntry)) error {
	if mac == "" {
		return fmt.Errorf("mac is empty")
	}
	if workDir == "" {
		return fmt.Errorf("workDir is empty")
	}
	if day == "" {
		return fmt.Errorf("day is empty")
	}
	if fn == nil {
		return fmt.Errorf("fn is nil")
	}
	tempFilePath := filepath.Join(workDir, mac)
	glog.Debug("updatetWorkTime", mac)
	works := ReadWorkTimeByMac(tempFilePath)
	if works == nil {
		works = make(map[string]*WorkEntry)
	}
	tempEntry := works[day]
	if tempEntry == nil {
		tempEntry = &WorkEntry{}
	}
	fn(tempEntry)
	works[day] = tempEntry
	for k, status := range works {
		glog.Printf("%v %+v", k, status)
	}
	content, err := ukey.StructToGob(works)
	if err != nil {
		return err
	}
	err = u.CheckDirector(workDir)
	file, err := os.Create(tempFilePath) // 文件不存在则创建，存在则截断
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入内容
	_, err = file.Write(content)
	return err
}

func UpdatetWorkTime(mac, day string, data map[string]interface{}) error {
	if data == nil {
		return fmt.Errorf("data map is empty")
	}
	return setWorkTime(mac, workDir, day, func(tempEntry *WorkEntry) {
		if v, ok := data["workTime1"]; ok {
			if vv, okk := v.(string); okk {
				t, err := u.AutoParse(vv)
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
				t, err := u.AutoParse(vv)
				if err == nil && t != nil {
					timestamp := t.UnixMilli()
					if !u.IsMillisecondTimestamp(timestamp) {
						timestamp *= 1000
					}
					tempEntry.OnWorkTime = timestamp
				}
			}
		}
		if v, ok := data["weekday"]; ok {
			if vv, okk := v.(int); okk {
				tempEntry.Weekday = vv
			}
		}
		if v, ok := data["isWeekDay"]; ok {
			if vv, okk := v.(bool); okk {
				tempEntry.IsWeekDay = vv
			}
		}
	})
}

//func sysLogUpdateWorkTime(mac string, timestamp int64) error {
//	if !u.IsMillisecondTimestamp(timestamp) {
//		timestamp *= 1000
//	}
//	ts := time.UnixMilli(timestamp)
//	return setWorkTime(mac, workDir, ts.Format(time.DateOnly), func(t *WorkEntry) {
//		if t.OnWorkTime == 0 {
//			t.OnWorkTime = timestamp
//			t.IsWeekDay = ts.Weekday() == time.Sunday || ts.Weekday() == time.Saturday
//		} else {
//			t.OffWorkTime = timestamp
//		}
//		t.Weekday = int(ts.Weekday())
//	})
//}

func sysLogUpdateWorkTime(mac string, timestamp int64, workType *WorkType) error {
	if workType == nil {
		return fmt.Errorf("workType is nil")
	}
	if timestamp <= 0 {
		return fmt.Errorf("timestamp is zero")
	}
	if !u.IsMillisecondTimestamp(timestamp) {
		timestamp *= 1000
	}
	off := u.GetTime(workType.OffWorkTime, u.GetLocation())
	if off == nil {
		return fmt.Errorf("off work time is nill")
	}
	t1 := time.UnixMilli(timestamp)
	day := t1.Format(time.DateOnly)
	offWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), off.Hour(), off.Minute(), off.Second(), off.Nanosecond(), t1.Location())
	isOffWorkTime := false
	if t1.Compare(offWorkTime) >= 0 {
		//下班
		isOffWorkTime = true
	}
	return setWorkTime(mac, workDir, day, func(t *WorkEntry) {
		if t.OnWorkTime == 0 {
			t.IsWeekDay = t1.Weekday() == time.Sunday || t1.Weekday() == time.Saturday
		}
		if isOffWorkTime {
			t.OffWorkTime = timestamp
		} else {
			t.OnWorkTime = timestamp
		}
		t.Weekday = int(t1.Weekday())
	})
}
