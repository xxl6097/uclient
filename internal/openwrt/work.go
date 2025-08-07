package openwrt

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/uclient/internal/u"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type WorkEntry struct {
	OnWorkTime    int64 `json:"onWorkTime"`
	OffWorkTime   int64 `json:"offWorkTime"`
	OnWorkSignal  int   `json:"onWorkSignal"`
	OffWorkSignal int   `json:"offWorkSignal"`
	Weekday       int   `json:"weekday"`
	DayType       int   `json:"dayType"` //0工作日，1节假日，2补班日，3加班日
}

// WorkTypeSetting time.Sunday || t1.Weekday() == time.Saturday
type WorkTypeSetting struct {
	OnWorkTime     string `json:"onWorkTime"`
	OffWorkTime    string `json:"offWorkTime"`
	WebhookUrl     string `json:"webhookUrl"`
	IsSaturdayWork bool   `json:"isSaturdayWork"` //默认false，意思是周六休息
}
type WorkTime struct {
	Date                  string        `json:"date"`
	Weekday               int           `json:"weekday"`
	WorkTime1             string        `json:"workTime1"`
	WorkTime2             string        `json:"workTime2"`
	OverWorkTimes         string        `json:"overWorkTimes"`
	DayType               int           `json:"dayType"` //0工作日，1节假日，2补班日，3加班日
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

func caculeteAMWorkDay(amSignTime, workTime1 time.Time) time.Duration {
	amOverTimes := u.CompareTime(amSignTime, workTime1)
	var duration time.Duration
	if amOverTimes > 0 {
		duration = time.Duration(amOverTimes) * time.Second
	}
	return duration
}
func caculetePMWorkDay(pmSignTime, workTime2 time.Time) time.Duration {
	pmOverTimes := u.CompareTime(workTime2, pmSignTime)
	var duration time.Duration
	if pmOverTimes > 0 {
		duration += time.Duration(pmOverTimes) * time.Second
	}
	return duration
}

//func caculeteWorkDay(amSignTime, pmSignTime, workTime1, workTime2 time.Time) time.Duration {
//	amOverTimes := u.CompareTime(amSignTime, workTime1)
//	pmOverTimes := u.CompareTime(workTime2, pmSignTime)
//	var duration time.Duration
//	if amOverTimes > 0 {
//		duration = time.Duration(amOverTimes) * time.Second
//	}
//	if pmOverTimes > 0 {
//		duration += time.Duration(pmOverTimes) * time.Second
//	}
//	return duration
//}

func GetWorkTime(mac, tempFilePath string, workType *WorkTypeSetting) ([]*Work, error) {
	if mac == "" {
		return nil, fmt.Errorf("mac is empty")
	}
	if tempFilePath == "" {
		return nil, fmt.Errorf("tempFilePath is empty")
	}
	if workType == nil {
		return nil, fmt.Errorf("workType is empty")
	}
	amSignTime, err := u.TimeParse(workType.OnWorkTime)
	if amSignTime == nil || err != nil {
		return nil, fmt.Errorf("on work time is nill %+v", workType)
	}
	pmSignTime, e2 := u.TimeParse(workType.OffWorkTime)
	if pmSignTime == nil || e2 != nil {
		return nil, fmt.Errorf("off work time is nill")
	}
	works := ReadWorkTimeByMac(tempFilePath)
	if works == nil {
		return nil, fmt.Errorf("works is empty")
	}

	result := make([]*Work, 0)
	months := make(map[string]*Work)
	for day, w := range works {
		workTime1 := u.UTC8ToTime(w.OnWorkTime)  //time1.Format(time.TimeOnly)
		workTime2 := u.UTC8ToTime(w.OffWorkTime) //time2.Format(time.TimeOnly)
		month := fmt.Sprintf("%d-%02d", workTime1.Year(), int(workTime1.Month()))
		d, e := u.AutoParse(day)
		if e == nil && d != nil {
			month = d.Format("2006-01")
		}
		var duration time.Duration
		//0工作日，1节假日，2补班日
		//如果是周六，且标记周六加班，那么加班时间不按照打开时间计算
		if w.Weekday == int(time.Saturday) && workType.IsSaturdayWork && w.DayType != 2 {
			if w.OnWorkTime > 0 && w.OffWorkTime > 0 {
				duration = time.Duration(u.CompareTime(workTime2, workTime1)) * time.Second
			}
		} else if w.DayType == 0 || w.DayType == 2 {
			if w.OnWorkTime > 0 {
				duration += caculeteAMWorkDay(*amSignTime, workTime1)
			}
			if w.OffWorkTime > 0 {
				duration += caculetePMWorkDay(*pmSignTime, workTime2)
			}
		}
		wrokTimeTemp := WorkTime{
			Date:                  day,
			DayType:               w.DayType,
			Weekday:               w.Weekday,
			OverWorkTimes:         duration.String(),
			OverWorkTimesDuration: duration,
		}

		if w.OnWorkTime > 0 {
			wrokTimeTemp.WorkTime1 = workTime1.Format(time.TimeOnly)
		}
		if w.OffWorkTime > 0 {
			wrokTimeTemp.WorkTime2 = workTime2.Format(time.TimeOnly)
		}
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

	//sort.Slice(result, func(i, j int) bool {
	//	a, b := result[i], result[j]
	//	sort.Slice(a.WorkTime, func(i, j int) bool {
	//		aa, ab := a.WorkTime[i], a.WorkTime[j]
	//		return aa.Date < ab.Date
	//	})
	//	sort.Slice(b.WorkTime, func(i, j int) bool {
	//		ba, bb := b.WorkTime[i], b.WorkTime[j]
	//		return ba.Date < bb.Date
	//	})
	//	return a.Month < b.Month
	//})

	for _, w := range result {
		sort.Slice(w.WorkTime, func(i, j int) bool {
			aa, ab := w.WorkTime[i], w.WorkTime[j]
			return aa.Date > ab.Date
		})
		for _, workTime := range w.WorkTime {
			w.OverTimeDuration += workTime.OverWorkTimesDuration
		}
		w.OverTime = w.OverTimeDuration.String()
	}
	sort.Slice(result, func(i, j int) bool {
		a, b := result[i], result[j]
		return a.Month > b.Month
	})

	//temp := result[0]
	//for _, w := range temp.WorkTime {
	//	temp.OverTimeDuration += w.OverWorkTimesDuration
	//}
	//fmt.Println(temp.Month, temp.OverTimeDuration.String())
	return result, nil
}

func getWorkTime(mac string, workType *WorkTypeSetting) ([]*Work, error) {
	tempFilePath := filepath.Join(workDir, mac)
	//glog.Debug("GetWorkTime", mac)
	return GetWorkTime(mac, tempFilePath, workType)
}

func TestSetWorkTime(isDel bool, mac, workDir, day string, fn func(*WorkEntry)) error {
	_, err := setWorkTime(isDel, mac, workDir, day, fn)
	return err
}

func setWorkTime(isDel bool, mac, workDir, day string, fn func(*WorkEntry)) (*WorkEntry, error) {
	if mac == "" {
		return nil, fmt.Errorf("mac is empty")
	}
	if workDir == "" {
		return nil, fmt.Errorf("workDir is empty")
	}
	if day == "" {
		return nil, fmt.Errorf("day is empty")
	}
	//if fn == nil {
	//	return fmt.Errorf("fn is nil")
	//}
	tempFilePath := filepath.Join(workDir, mac)
	//glog.Debug("updatetWorkTime", mac)
	works := ReadWorkTimeByMac(tempFilePath)
	if works == nil {
		works = make(map[string]*WorkEntry)
	}
	tempEntry := works[day]
	if tempEntry == nil {
		tempEntry = &WorkEntry{}
	}
	if fn != nil {
		fn(tempEntry)
	}
	works[day] = tempEntry
	if isDel {
		delete(works, day)
	}
	glog.Debugf("更新打卡 %v %+v", mac, tempEntry)
	//for k, status := range works {
	//	glog.Printf("%v %+v", k, status)
	//}
	content, err := ukey.StructToGob(works)
	if err != nil {
		return tempEntry, err
	}
	err = u.CheckDirector(workDir)
	file, err := os.Create(tempFilePath) // 文件不存在则创建，存在则截断
	if err != nil {
		return tempEntry, err
	}
	defer file.Close()
	// 写入内容
	_, err = file.Write(content)
	return tempEntry, err
}

func UpdatetWorkTime(mac, day string, data map[string]interface{}) error {
	if data == nil {
		return fmt.Errorf("data map is empty")
	}
	_, err := setWorkTime(false, mac, workDir, day, func(tempEntry *WorkEntry) {
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
	return err
}
func AddWorkTime(mac string, timestamp int64, isOnWork bool) error {
	if timestamp <= 0 {
		return fmt.Errorf("timestamp is zero")
	}
	if !u.IsMillisecondTimestamp(timestamp) {
		timestamp *= 1000
	}
	t1 := u.UTC8ToTime(timestamp)
	day := t1.Format(time.DateOnly)
	_, err := setWorkTime(false, mac, workDir, day, func(t *WorkEntry) {
		t.Weekday = int(t1.Weekday())
		if isOnWork {
			t.OnWorkTime = timestamp
		} else {
			t.OffWorkTime = timestamp
		}
	})
	return err
}

func DelWorkTime(mac string, day string) error {
	_, err := setWorkTime(true, mac, workDir, day, nil)
	return err
}

func GetTodaySign(mac string) *WorkEntry {
	day := glog.Now().Format(time.DateOnly)
	tempFilePath := filepath.Join(workDir, mac)
	works := ReadWorkTimeByMac(tempFilePath)
	if works == nil {
		return &WorkEntry{}
	}
	tempEntry := works[day]
	if tempEntry == nil {
		return &WorkEntry{}
	}
	return tempEntry
}

func sysLogUpdateWorkTime(tempData *DHCPLease) (*WorkEntry, error) {
	if tempData == nil || tempData.Nick == nil && tempData.Nick.WorkType == nil && tempData.Nick.WorkType.OnWorkTime == "" {
		return nil, fmt.Errorf("参数不全 %+v", tempData)
	}
	workType := tempData.Nick.WorkType
	mac := tempData.MAC
	timestamp := tempData.Nick.StartTime
	if workType == nil {
		return nil, fmt.Errorf("考勤打卡未设置")
	}
	if mac == "" {
		return nil, fmt.Errorf("设备Mac空～")
	}
	if workType.OnWorkTime == "" || workType.OffWorkTime == "" {
		return nil, fmt.Errorf("考勤打卡时间未设置 %+v", workType)
	}
	if timestamp <= 0 {
		return nil, fmt.Errorf("timestamp is zero")
	}
	if !u.IsMillisecondTimestamp(timestamp) {
		timestamp *= 1000
	}
	workingTime, err := u.IsWorkingTime(workType.OnWorkTime, workType.OffWorkTime)
	if err != nil {
		return nil, err
	}
	t1 := u.UTC8ToTime(timestamp)
	day := t1.Format(time.DateOnly)
	//glog.Debug("系统监听更新", mac, workingTime, u.UTC8ToString(timestamp, time.DateTime))
	return setWorkTime(false, mac, workDir, day, func(t *WorkEntry) {
		t.Weekday = int(t1.Weekday())
		if workType.IsSaturdayWork && t1.Weekday() == time.Saturday {
			t.DayType = 3
		}
		if t1.Weekday() == time.Saturday || t1.Weekday() == time.Sunday {
			if t.OnWorkTime == 0 {
				t.OnWorkTime = timestamp
				t.OnWorkSignal = tempData.Signal
			} else {
				t.OffWorkTime = timestamp
				t.OffWorkSignal = tempData.Signal
			}
		} else {
			if workingTime == 0 {
				//上班打卡
				if t.OnWorkTime <= 0 {
					//说明上午未打卡
					t.OnWorkTime = timestamp
					t.OnWorkSignal = tempData.Signal
				}
			} else if workingTime == 2 {
				t.OffWorkTime = timestamp
				t.OffWorkSignal = tempData.Signal
			}
		}
	})
}
func sysLogUpdateWorkTime1(mac string, timestamp int64, workType *WorkTypeSetting) (*WorkEntry, error) {
	if workType == nil {
		return nil, fmt.Errorf("考勤打卡未设置")
	}
	if workType.OnWorkTime == "" || workType.OffWorkTime == "" {
		return nil, fmt.Errorf("考勤打卡时间未设置 %+v", workType)
	}
	if timestamp <= 0 {
		return nil, fmt.Errorf("timestamp is zero")
	}
	if !u.IsMillisecondTimestamp(timestamp) {
		timestamp *= 1000
	}
	workingTime, err := u.IsWorkingTime(workType.OnWorkTime, workType.OffWorkTime)
	if err != nil {
		return nil, err
	}
	t1 := u.UTC8ToTime(timestamp)
	day := t1.Format(time.DateOnly)
	//glog.Debug("系统监听更新", mac, workingTime, u.UTC8ToString(timestamp, time.DateTime))
	return setWorkTime(false, mac, workDir, day, func(t *WorkEntry) {
		t.Weekday = int(t1.Weekday())
		if workType.IsSaturdayWork && t1.Weekday() == time.Saturday {
			t.DayType = 3
		}
		if t1.Weekday() == time.Saturday || t1.Weekday() == time.Sunday {
			if t.OnWorkTime == 0 {
				t.OnWorkTime = timestamp
			} else {
				t.OffWorkTime = timestamp
			}
		} else {
			if workingTime == 0 {
				//上班打卡
				if t.OnWorkTime <= 0 {
					//说明上午未打卡
					t.OnWorkTime = timestamp
				}
			} else if workingTime == 2 {
				t.OffWorkTime = timestamp
			}
		}
	})
}
