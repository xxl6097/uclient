package openwrt

import (
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"time"
)

func (this *openWRT) ddingNotify(eveName string, tempData *DHCPLease) {
	glog.Errorf("ddingNotify %+v", tempData)
	err := this.notifyWebhookMessage(eveName, tempData)
	if err != nil {
		glog.Errorf("钉钉通知失败 %v %+v", err, tempData)
	}
}
func (this *openWRT) ddingWorkSign(tempData *DHCPLease) {
	this.updateWorkTime(tempData, func(working int, wrk *WorkEntry) {
		e := this.NotifySignCardEvent(working, tempData.Signal, tempData.MAC, wrk)
		if e != nil {
			glog.Errorf("钉钉通知打卡失败 %v %+v", e, tempData)
		}
	})
}

func (this *openWRT) updateWorkTime(tempData *DHCPLease, canNotifyDing func(int, *WorkEntry)) {
	if !this.isSignTime(tempData) {
		return
	}
	timestamp := tempData.StartTime
	if timestamp == 0 {
		return
	}
	mac := tempData.MAC
	ti := u.UTC8ToTime(timestamp)
	todayDate := ti.Format(time.DateOnly)
	workingTime, err := u.IsWorkingTime(tempData.Nick.WorkType.OnWorkTime, tempData.Nick.WorkType.OffWorkTime)
	if err != nil {
		glog.Error("判断工作时间错误❌", err)
	}
	signWork, err := UpdateWorkTime(mac, todayDate, func(todayData *WorkEntry) {
		if ti.Weekday() == time.Sunday {
			//默认情况周日是节假日
			if todayData.OnWorkTime == 0 &&
				todayData.OffWorkTime == 0 &&
				todayData.OffWorkSignal == 0 &&
				todayData.OnWorkSignal == 0 &&
				todayData.DayType == 0 &&
				todayData.Weekday == 0 {
				todayData.DayType = 1
			}
		}
		todayData.Weekday = int(ti.Weekday())
		if tempData.Nick.WorkType.IsSaturdayWork && ti.Weekday() == time.Saturday {
			todayData.DayType = 3
		}
		if ti.Weekday() == time.Saturday || ti.Weekday() == time.Sunday {
			if todayData.OnWorkTime == 0 {
				todayData.OnWorkTime = timestamp
				todayData.OnWorkSignal = tempData.Signal
			} else {
				todayData.OffWorkTime = timestamp
				todayData.OffWorkSignal = tempData.Signal
			}
		} else {
			if workingTime == 0 {
				//上班打卡
				if todayData.OnWorkTime <= 0 {
					//说明上午未打卡
					todayData.OnWorkTime = timestamp
					todayData.OnWorkSignal = tempData.Signal
				}
			} else if workingTime == 2 {
				todayData.OffWorkTime = timestamp
				todayData.OffWorkSignal = tempData.Signal
			}
		}
	})
	if signWork != nil && err == nil && canNotifyDing != nil {
		canNotifyDing(workingTime, signWork)
	}
}

// 判断设备具备打开条件，也就是是否设置了上线班时间
func (this *openWRT) hasSignCondition(tempData *DHCPLease) bool {
	if tempData != nil && tempData.MAC != "" {
		if tempData.Nick != nil && tempData.Nick.WorkType != nil && tempData.Nick.WorkType.OnWorkTime != "" && tempData.Nick.WorkType.OffWorkTime != "" {
			return true
		}
	}
	return false
}

func (this *openWRT) isSignTime(tempData *DHCPLease) bool {
	if this.hasSignCondition(tempData) {
		t1 := glog.Now()
		if t1.Weekday() == time.Saturday || t1.Weekday() == time.Sunday {
			return true
		} else {
			working, e1 := u.IsWorkingTime(tempData.Nick.WorkType.OnWorkTime, tempData.Nick.WorkType.OffWorkTime)
			if e1 == nil {
				switch working {
				case 0:
					return true
				case 2:
					return true
				}
			}
		}
	}
	//不具备打卡条件，返回false
	return false
}

func (this *openWRT) ddingSignByRSSI(tempData *DHCPLease) {
	if !this.isSignTime(tempData) {
		//不具备打卡条件或者不在打开时间范围内（工作时间不打卡），退出
		return
	}
	if this.isWeekend() {
		wk := GetTodaySign(tempData.MAC)
		if wk.OnWorkTime <= 0 {
			this.ddingWorkSign(tempData)
		} else if wk.OffWorkTime <= 0 {
			this.ddingWorkSign(tempData)
		} else if wk.OffWorkTime > 0 {
			this.ddingWorkSign(tempData)
		}
	} else {
		working, e1 := u.IsWorkingTime(tempData.Nick.WorkType.OnWorkTime, tempData.Nick.WorkType.OffWorkTime)
		if e1 == nil {
			switch working {
			case 0:
				//这里要设置信号门槛，不然在楼下就连上触发了打卡不行
				//TODO 如果打卡的时候，信号确实小于-80，如何处理
				if tempData.Online && tempData.Signal >= -80 {
					wk := GetTodaySign(tempData.MAC)
					if wk.OnWorkTime == 0 {
						this.ddingWorkSign(tempData)
					}
				}
				break
			case 2:
				wk := GetTodaySign(tempData.MAC)
				if wk.OffWorkTime <= 0 {
					this.ddingWorkSign(tempData)
				} else if wk.OffWorkTime > 0 && tempData.Signal != wk.OffWorkSignal {
					this.ddingWorkSign(tempData)
				}
				break
			}
		} else {
			glog.Debug(e1)
		}
	}
}

// 判断是否为周末
func (this *openWRT) isWeekend() bool {
	t1 := glog.Now()
	if t1.Weekday() == time.Saturday || t1.Weekday() == time.Sunday {
		return true
	}
	return false
}

// 具备打卡条件，而且信号变弱，故判断可能离线了
func (this *openWRT) signalWeak(tempData *DHCPLease) {
	if this.isSignTime(tempData) {
		if _, ok := this.tempOffline[tempData.MAC]; ok {
			return
		}
		if tempData.Signal != 0 && tempData.Signal < -80 {
			if !u.Ping(tempData.IP) {
				//ping不通，估计离线了
				glog.Warnf("已经离线了 %+v", tempData)
				this.tempOffline[tempData.MAC] = &DHCPLease{
					MAC:       tempData.MAC,
					IP:        tempData.IP,
					Signal:    tempData.Signal,
					Ssid:      tempData.Ssid,
					Hostname:  tempData.Hostname,
					StartTime: glog.Now().UnixMilli(),
					Online:    false,
				}
			}
		}
	}
}
