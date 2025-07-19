package openwrt

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"time"
)

func (this *openWRT) ddingNotify(tempData *DHCPLease) {
	_ = this.notifyWebhookMessage(tempData)
}

func (this *openWRT) ddingWorkSign(tempData *DHCPLease) {
	if tempData.Nick != nil {
		err := sysLogUpdateWorkTime(tempData.MAC, tempData.StartTime, tempData.Nick.WorkType, func(working int, macAddress string, t1 time.Time) {
			_ = this.NotifySignCardEvent(working, macAddress, t1)
		})
		if err != nil {
			glog.Error(fmt.Errorf("updatetWorkTime Error:%v", err))
		}
	}
}
