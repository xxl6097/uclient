package openwrt

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"regexp"
	"strings"
	"time"
)

// KernelLog logStr := "Mon Jul 28 18:12:45 2025 kern.warn kernel: [34592.938265] 7981@C13L2,MacTableDeleteEntry() 1921: Del Sta:ee:af:48:c9:e6:c1"
// logStr = "Mon Jul 28 18:15:11 2025 kern.warn kernel: [34739.452859] 7981@C13L2,MacTableInsertEntry() 1559: New Sta:ee:af:48:c9:e6:c1"
// KernelLog 定义日志结构体
type KernelLog struct {
	Timestamp  int64  // 时间戳
	MACAddress string // MAC地址
	Online     bool   // MAC地址
}

func parseHetSysLog(logStr string) *KernelLog {
	online := false
	if strings.Contains(logStr, "MacTableDeleteEntry") {
		online = false
	} else if strings.Contains(logStr, "MacTableInsertEntry") {
		online = true
	} else {
		return nil
	}
	macAddress := ParseMacAddr(logStr)
	if macAddress == "" {
		glog.Errorf("parseHetSysLog: invalid macAddress")
		return nil
	}
	// 填充结构体
	logEntry := KernelLog{
		Timestamp:  glog.Now().UnixMilli(),
		MACAddress: macAddress,
		Online:     online,
	}
	return &logEntry
}

func subscribeHetSysLog(fn func(*KernelLog)) error {
	//args := []string{"-f", "|", "grep", "hostapd.*"}
	pattern := `MacTable.*`
	re := regexp.MustCompile(pattern)
	return command(func(s string) {
		if re.MatchString(s) {
			if strings.Contains(s, "MacTableDeleteEntry") || strings.Contains(s, "MacTableInsertEntry") {
				tempData := parseHetSysLog(s)
				glog.Debug("--->", tempData)
				if fn != nil && tempData != nil {
					fn(tempData)
				}
			}
		}

	}, "logread", "-f")
}

func (this *openWRT) subscribeHetSysLog() {
	tryCount := 0
	for {
		err := subscribeHetSysLog(func(event *KernelLog) {
			if event != nil && event.MACAddress != "" {
				glog.Infof("HetSysLog事件:%+v", event)
				eve := &DHCPLease{
					MAC:       event.MACAddress,
					Online:    event.Online,
					StartTime: glog.Now().UnixMilli(),
				}
				if v, ok := this.leases[eve.MAC]; ok {
					if v.Hostname != "" {
						eve.Hostname = v.Hostname
					}
				}
				this.updateDeviceStatus("HetSysLog事件", eve)
			}
		})
		if err != nil {
			glog.Error(fmt.Errorf("SysLog监听失败 %v", err))
			time.Sleep(time.Second * 10)
			glog.Error("重新监听 SysLog")
			tryCount++
			if tryCount > RE_REY_MAX_COUNT {
				glog.Error("监听 SysLog 失败，超过最大重试次数")
				break
			}
		}
	}
}
