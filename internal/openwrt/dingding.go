package openwrt

import (
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
)

func (this *openWRT) ddingNotify(tempData *DHCPLease) {
	err := this.notifyWebhookMessage(tempData)
	if err != nil {
		glog.Errorf("钉钉通知失败 %v %+v", err, tempData)
	}
}
func (this *openWRT) ddingWorkSign(tempData *DHCPLease) {
	if tempData.Nick != nil && tempData.Nick.WorkType != nil && tempData.Nick.WorkType.OnWorkTime != "" {
		_, err := sysLogUpdateWorkTime(tempData.MAC, tempData.StartTime, tempData.Nick.WorkType)
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
	if tempData != nil && !tempData.Online {
		this.ddingWorkSign(tempData)
	}
}
func (this *openWRT) ddingWorkOnSign(tempData *DHCPLease) {
	if tempData != nil && tempData.Online && !tempData.IsOnWorkSign && tempData.Signal >= -80 {
		tempData.IsOnWorkSign = true
		this.ddingWorkSign(tempData)
	}
}
