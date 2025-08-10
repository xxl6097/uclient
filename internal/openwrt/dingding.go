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
	//1. 判断具备打卡条件
	hasSignCondition, working, now := this.isSignTime(tempData)
	if !hasSignCondition {
		//不具备打卡条件或者不在打开时间范围内（工作时间不打卡），退出
		return
	}
	//2. 打卡时间戳要正确
	timestamp := tempData.StartTime
	if timestamp == 0 {
		return
	}

	//满足条件后：
	//1. 读取今天打卡信息；
	//2. 如果是周六或周日，那么每次上线离线，都属于打卡
	//3. 如果是工作日
	signTime := u.UTC8ToTime(timestamp)
	mac := tempData.MAC
	signDatas := GetSignData(mac)
	if signDatas == nil {
		signDatas = make(map[string]*WorkEntry)
	}
	todaySignData := signDatas[now.Format(time.DateOnly)]
	if todaySignData == nil {
		todaySignData = &WorkEntry{}
		if signTime.Weekday() == time.Sunday {
			//默认情况周日是节假日
			todaySignData.DayType = 1
		}
	}

	//设置周几
	todaySignData.Weekday = int(signTime.Weekday())
	//如果设置周六位加班日，则设置DayType=3
	if tempData.Nick.WorkType.IsSaturdayWork && signTime.Weekday() == time.Saturday {
		todaySignData.DayType = 3
	}
	needUpdateSign := false
	if signTime.Weekday() == time.Saturday || signTime.Weekday() == time.Sunday {
		if todaySignData.OnWorkTime == 0 {
			todaySignData.OnWorkTime = timestamp
			todaySignData.OnWorkSignal = tempData.Signal
		} else {
			todaySignData.OffWorkTime = timestamp
			todaySignData.OffWorkSignal = tempData.Signal
		}
		needUpdateSign = true
	} else {
		if working == 0 {
			//上班打卡
			if todaySignData.OnWorkTime <= 0 {
				//说明上午未打卡
				todaySignData.OnWorkTime = timestamp
				todaySignData.OnWorkSignal = tempData.Signal
				needUpdateSign = true
			}
		} else if working == 2 {
			todaySignData.OffWorkTime = timestamp
			todaySignData.OffWorkSignal = tempData.Signal
			needUpdateSign = true
		}
	}

	if needUpdateSign {
		err := SetSignData(mac, signDatas)
		if err != nil {
			glog.Errorf("打卡更新失败 %v %+v", err, tempData)
		} else {
			glog.Errorf("打卡更新成功 %v %+v", err, tempData)
			e := this.NotifyDingSign(tempData.Signal, tempData.MAC, now, todaySignData, tempData.Nick.WorkType)
			if e != nil {
				glog.Errorf("钉钉通知打卡失败 %v %+v", e, tempData)
			}
		}
	}

}

//func (this *openWRT) updateWorkTime1(tempData *DHCPLease, canNotifyDing func(int, *WorkEntry)) {
//	hasSignCondition, working := this.isSignTime(tempData)
//	if !hasSignCondition {
//		//不具备打卡条件或者不在打开时间范围内（工作时间不打卡），退出
//		return
//	}
//	timestamp := tempData.StartTime
//	if timestamp == 0 {
//		return
//	}
//	mac := tempData.MAC
//	ti := u.UTC8ToTime(timestamp)
//	todayDate := ti.Format(time.DateOnly)
//	signWork, err := UpdateWorkTime(mac, todayDate, func(todayData *WorkEntry) {
//		if ti.Weekday() == time.Sunday {
//			//默认情况周日是节假日
//			if todayData.OnWorkTime == 0 &&
//				todayData.OffWorkTime == 0 &&
//				todayData.OffWorkSignal == 0 &&
//				todayData.OnWorkSignal == 0 &&
//				todayData.DayType == 0 &&
//				todayData.Weekday == 0 {
//				todayData.DayType = 1
//			}
//		}
//		todayData.Weekday = int(ti.Weekday())
//		if tempData.Nick.WorkType.IsSaturdayWork && ti.Weekday() == time.Saturday {
//			todayData.DayType = 3
//		}
//		if ti.Weekday() == time.Saturday || ti.Weekday() == time.Sunday {
//			if todayData.OnWorkTime == 0 {
//				todayData.OnWorkTime = timestamp
//				todayData.OnWorkSignal = tempData.Signal
//			} else {
//				todayData.OffWorkTime = timestamp
//				todayData.OffWorkSignal = tempData.Signal
//			}
//		} else {
//			if working == 0 {
//				//上班打卡
//				if todayData.OnWorkTime <= 0 {
//					//说明上午未打卡
//					todayData.OnWorkTime = timestamp
//					todayData.OnWorkSignal = tempData.Signal
//				}
//			} else if working == 2 {
//				todayData.OffWorkTime = timestamp
//				todayData.OffWorkSignal = tempData.Signal
//			}
//		}
//	})
//	if signWork != nil && err == nil && canNotifyDing != nil {
//		canNotifyDing(working, signWork)
//	}
//}

// 判断设备具备打开条件，也就是是否设置了上线班时间
func (this *openWRT) hasSignCondition(tempData *DHCPLease) bool {
	if tempData != nil && tempData.MAC != "" {
		if tempData.Nick != nil && tempData.Nick.WorkType != nil && tempData.Nick.WorkType.OnWorkTime != "" && tempData.Nick.WorkType.OffWorkTime != "" {
			return true
		}
	}
	return false
}

func (this *openWRT) isSignTime(tempData *DHCPLease) (bool, int, time.Time) {
	if this.hasSignCondition(tempData) {
		now := glog.Now()
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			return true, -1, now
		} else {
			working, e1 := u.IsWorkingTime(tempData.Nick.WorkType.OnWorkTime, tempData.Nick.WorkType.OffWorkTime)
			if e1 == nil {
				switch working {
				case 0:
					return true, working, now
				case 2:
					return true, working, now
				}
			} else {
				glog.Error("判断工作时间错误❌", e1)
			}
		}
	}
	//不具备打卡条件，返回false
	return false, -1, time.Time{}
}

func (this *openWRT) ddingSignByRSSI(tempData *DHCPLease) {
	hasSignCondition, working, _ := this.isSignTime(tempData)
	if !hasSignCondition {
		//不具备打卡条件或者不在打开时间范围内（工作时间不打卡），退出
		return
	}
	if this.isWeekend() {
		wk := GetTodaySignData(tempData.MAC)
		if wk.OnWorkTime <= 0 {
			this.ddingWorkSign(tempData)
		} else if wk.OffWorkTime <= 0 {
			this.ddingWorkSign(tempData)
		} else if wk.OffWorkTime > 0 {
			this.ddingWorkSign(tempData)
		}
	} else {
		switch working {
		case 0:
			//这里要设置信号门槛，不然在楼下就连上触发了打卡不行
			//TODO 如果打卡的时候，信号确实小于-80，如何处理
			if tempData.Online && tempData.Signal >= -80 {
				wk := GetTodaySignData(tempData.MAC)
				if wk.OnWorkTime == 0 {
					this.ddingWorkSign(tempData)
				}
			}
			break
		case 2:
			wk := GetTodaySignData(tempData.MAC)
			if wk.OffWorkTime <= 0 {
				this.ddingWorkSign(tempData)
			} else if wk.OffWorkTime > 0 && tempData.Signal != wk.OffWorkSignal {
				this.ddingWorkSign(tempData)
			}
			break
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
	hasSignCondition, _, _ := this.isSignTime(tempData)
	if !hasSignCondition {
		//不具备打卡条件或者不在打开时间范围内（工作时间不打卡），退出
		return
	}
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
