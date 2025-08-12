package openwrt

import (
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"time"
)

func (this *openWRT) ding(eveName string, tempData *DHCPLease) {
	if this.hasNotifyCondition(tempData) {
		go func() {
			err := this.notifyWebhookMessage(eveName, tempData)
			if err != nil {
				glog.Errorf("钉钉通知失败 %v %+v", err, tempData)
			}
		}()
	}
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
		signDatas[now.Format(time.DateOnly)] = todaySignData
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
		signDatas[now.Format(time.DateOnly)] = todaySignData
		e := SetSignData(mac, signDatas)
		if e != nil {
			glog.Errorf("打卡更新失败 %v %+v", e, tempData)
		} else {
			glog.Errorf("打卡更新成功  %+v", tempData)
			e1 := this.NotifyDingSign(tempData, eveName, now, todaySignData)
			if e1 != nil {
				glog.Errorf("钉钉打卡失败 %v %+v", e1, tempData)
			}
		}
	}

}

func (this *openWRT) hasNotifyCondition(tempData *DHCPLease) bool {
	if tempData != nil && tempData.MAC != "" {
		if tempData.Nick != nil && tempData.Nick.IsPush {
			return true
		}
	}
	return false
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
				default:
					return false, working, now
				}
			} else {
				glog.Error("判断工作时间错误❌", e1)
			}
		}
	}
	//不具备打卡条件，返回false
	return false, -1, time.Time{}
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
	hasSignCondition, working, now := this.isSignTime(tempData)
	if !hasSignCondition {
		//不具备打卡条件或者不在打开时间范围内（工作时间不打卡），退出
		if working == 1 {
			staInfo := GetStaInfo()
			if staInfo != nil {
				sta := staInfo[tempData.MAC]
				if sta != nil {
					glog.Debug(tempData.Hostname, tempData.IP, tempData.MAC, tempData.Signal)
					glog.Debugf("sta %+v", sta)
				}
			}

			delete(this.tempOffline, tempData.MAC)
		}
		return
	}

	if working == 0 {
		//上班时间
		if tempData.Signal != 0 && tempData.Signal >= -81 {
			if v, ok := this.tempOffline[tempData.MAC]; ok {
				if v.OnWorkTime > 0 {
					if !u.IsTimestampToday(v.OnWorkTime) {
						v.OnWorkTime = tempData.StartTime
						v.OnWorkSignal = tempData.Signal
					}
				} else {
					v.OnWorkTime = tempData.StartTime
					v.OnWorkSignal = tempData.Signal
				}
			} else {
				this.tempOffline[tempData.MAC] = &WorkEntry{
					OnWorkTime:   tempData.StartTime,
					OnWorkSignal: tempData.Signal,
				}
			}
		}
	} else if working == 2 {
		//下班时间
		if tempData.Signal != 0 && tempData.Signal < -81 {
			if v, ok := this.tempOffline[tempData.MAC]; ok {
				if v.OffWorkTime > 0 {
					if u.IsTimestampToday(v.OffWorkTime) {
						v.OffWorkTime = tempData.StartTime
						v.OffWorkSignal = tempData.Signal
					}
				} else {
					v.OffWorkTime = tempData.StartTime
					v.OffWorkSignal = tempData.Signal
				}
			} else {
				this.tempOffline[tempData.MAC] = &WorkEntry{
					OffWorkTime:   tempData.StartTime,
					OffWorkSignal: tempData.Signal,
				}
			}
		}
	} else {
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			if tempData.Signal != 0 && tempData.Signal < -80 && tempData.Signal > -90 {
				if v, ok := this.tempOffline[tempData.MAC]; ok {
					if v.OnWorkTime > 0 {
						if !u.IsTimestampToday(v.OnWorkTime) {
							v.OnWorkTime = tempData.StartTime
							v.OnWorkSignal = tempData.Signal
						} else {
							v.OffWorkTime = tempData.StartTime
							v.OffWorkSignal = tempData.Signal
						}
					} else {
						v.OnWorkTime = tempData.StartTime
						v.OnWorkSignal = tempData.Signal
					}
				} else {
					this.tempOffline[tempData.MAC] = &WorkEntry{
						OnWorkTime:   tempData.StartTime,
						OnWorkSignal: tempData.Signal,
					}
				}
			}
		} else {
			delete(this.tempOffline, tempData.MAC)
		}
	}
}
