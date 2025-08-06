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
	if tempData.Nick != nil && tempData.Nick.WorkType != nil && tempData.Nick.WorkType.OnWorkTime != "" {
		_, err := sysLogUpdateWorkTime(tempData)
		if err != nil {
			glog.Errorf("更新时间失败 %v %+v", err, tempData)
		} else {
			working, e1 := u.IsWorkingTime(tempData.Nick.WorkType.OnWorkTime, tempData.Nick.WorkType.OffWorkTime)
			if e1 != nil {
				glog.Error(e1)
			}
			e := this.NotifySignCardEvent(working, tempData.Signal, tempData.MAC)
			if e != nil {
				glog.Errorf("钉钉通知打卡失败 %v %+v", e, tempData)
			}
		}
	} else {
		glog.Errorf("未设置打卡 %+v", tempData)
	}
}

func (this *openWRT) ddingWorkOffSign(tempData *DHCPLease) {
	if tempData != nil && !tempData.Online && this.isSignTime(tempData) {
		glog.Errorf("ddingWorkOffSign %+v", tempData)
		this.ddingWorkSign(tempData)
	}
}

func (this *openWRT) isSignTime(tempData *DHCPLease) bool {
	if tempData != nil {
		if tempData.Nick != nil && tempData.Nick.WorkType != nil && tempData.Nick.WorkType.OnWorkTime != "" {
			working, e1 := u.IsWorkingTime(tempData.Nick.WorkType.OnWorkTime, tempData.Nick.WorkType.OffWorkTime)
			if e1 == nil {
				switch working {
				case 0:
					return true
				case 2:
					return true
				default:
					t1 := glog.Now()
					if t1.Weekday() == time.Saturday || t1.Weekday() == time.Sunday {
						return true
					}
				}
			}
		}
	}
	return false
}

func (this *openWRT) ddingSign(tempData *DHCPLease) {
	if tempData != nil {
		if tempData.Nick != nil && tempData.Nick.WorkType != nil && tempData.Nick.WorkType.OnWorkTime != "" {
			working, e1 := u.IsWorkingTime(tempData.Nick.WorkType.OnWorkTime, tempData.Nick.WorkType.OffWorkTime)
			if e1 == nil {
				switch working {
				case 0:
					if tempData.Online && tempData.Signal >= -80 {
						wk := GetTodaySign(tempData.MAC)
						if wk.OnWorkTime == 0 {
							this.ddingWorkSign(tempData)
						}
					}
					break
				case 2:
					if tempData.Signal < -80 && tempData.Signal > -90 {
						wk := GetTodaySign(tempData.MAC)
						if wk.OffWorkTime <= 0 {
							this.ddingWorkSign(tempData)
						} else if wk.OffWorkTime > 0 && tempData.Signal != wk.OffWorkSignal {
							this.ddingWorkSign(tempData)
						}
					}
					break
				default:
					t1 := glog.Now()
					if t1.Weekday() == time.Saturday || t1.Weekday() == time.Sunday {
						if (tempData.Signal < -80 && tempData.Signal > -90) || (tempData.Online && tempData.Signal >= -80) {
							wk := GetTodaySign(tempData.MAC)
							if wk.OnWorkTime <= 0 {
								this.ddingWorkSign(tempData)
							} else if wk.OffWorkTime <= 0 {
								this.ddingWorkSign(tempData)
							} else if wk.OffWorkTime > 0 && tempData.Signal != wk.OffWorkSignal {
								this.ddingWorkSign(tempData)
							}

						}
					}
				}
			}
		}
	}
}
