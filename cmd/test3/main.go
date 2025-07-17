package main

import (
	"github.com/xxl6097/glog/glog"
	"time"
)

// CompareTime 小于目标时间
func CompareTime(now, target time.Time) int {
	nowSeconds := now.Hour()*3600 + now.Minute()*60 + now.Second()
	tarSeconds := target.Hour()*3600 + target.Minute()*60 + target.Second()
	return nowSeconds - tarSeconds
}

// IsWorkingTime
// 0：可以上班打卡
// 1：工作时间，不能打卡
// 2：可以下班打卡
func IsWorkingTime(time1, time2 string) (int, error) {
	t1, e1 := time.Parse(time.TimeOnly, time1)
	if e1 != nil {
		return -1, e1
	}
	t2, e2 := time.Parse(time.TimeOnly, time2)
	if e2 != nil {
		return -1, e2
	}
	now := glog.Now()
	if CompareTime(now, t1) <= 0 {
		//小于等于上班时间
		return 0, nil
	} else if CompareTime(now, t2) < 0 {
		//工作时间内
		return 1, nil
	} else {
		//大于等于下班时间
		return 2, nil
	}
}

func main() {
	timestring := "09:00:00"
	t, _ := time.Parse(time.TimeOnly, timestring)
	now := glog.Now()
	CompareTime(t, now)
}
