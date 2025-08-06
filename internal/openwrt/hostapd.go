package openwrt

import (
	"context"
	"encoding/json"
	"github.com/xxl6097/glog/glog"
	"strings"
	"time"
)

// HostapdDevice ### 📡 **字段详细含义**
// | **字段名**   | **值示例**          | **含义与作用**                                                                 |
// |--------------|---------------------|------------------------------------------------------------------------------|
// | `address`    | `5a:a7:22:62:3d:26` | **客户端设备的MAC地址**，唯一标识接入网络的终端设备（如手机、笔记本）。 |
// | `target`     | `88:c3:97:07:55:21` | **接入点（AP）的MAC地址**，表示客户端连接的目标路由器或AP的物理地址。             |
// | `signal`     | `-44`               | **信号强度**（单位：dBm）：负值（越接近0表示信号越强）。`-44`为**优秀信号**（通常`-50`以上为良好）。 |
// | `freq`       | `5180`              | **无线频段频率**（单位：MHz）：`5180`属于**5GHz频段**（常见频段：2.4GHz范围为`2400~2483`，5GHz为`5150~5850`）。 |
type HostapdDevice struct {
	Address   string    `json:"address"`
	Target    string    `json:"target"`
	Signal    int       `json:"signal"`
	Freq      int       `json:"freq"`
	DataType  int       `json:"datatype"` //0:上线，1：离线，2：状态
	Timestamp time.Time `json:"timestamp"`
}

type Assoc struct {
	Assoc *HostapdDevice `json:"assoc"`
}
type Probe struct {
	Probe *HostapdDevice `json:"probe"`
}
type DisAssoc struct {
	Disassoc *HostapdDevice `json:"disassoc"`
}

func deviceOnline(s string) *HostapdDevice {
	var tempData Assoc
	err := json.Unmarshal([]byte(s), &tempData)
	if err == nil {
		if tempData.Assoc != nil {
			tempData.Assoc.Timestamp = glog.Now()
			tempData.Assoc.DataType = 0
		}
		return tempData.Assoc
	}
	return nil
}

func deviceOffline(s string) *HostapdDevice {
	var tempData DisAssoc
	err := json.Unmarshal([]byte(s), &tempData)
	if err == nil {
		if tempData.Disassoc != nil {
			tempData.Disassoc.Timestamp = glog.Now()
			tempData.Disassoc.DataType = 1
		}
		return tempData.Disassoc
	}
	return nil
}
func deviceStatus(s string) *HostapdDevice {
	var tempData Probe
	err := json.Unmarshal([]byte(s), &tempData)
	if err == nil {
		if tempData.Probe != nil {
			tempData.Probe.Timestamp = glog.Now()
			tempData.Probe.DataType = 2
		}
		return tempData.Probe
	}
	return nil
}
func decode(s string, fn func(*HostapdDevice)) {
	var tempData map[string]interface{}
	err := json.Unmarshal([]byte(s), &tempData)
	if err == nil {
		for key, _ := range tempData {
			var tempDevice *HostapdDevice
			switch key {
			case "assoc":
				tempDevice = deviceOnline(s)
				//glog.Debugf("上线 %+v", tempData)
				break
			case "disassoc":
				tempDevice = deviceOffline(s)
				//glog.Debugf("离线 %+v", tempData)
				break
			case "probe":
				tempDevice = deviceStatus(s)
				//glog.Debugf("状态 %+v", tempData)
				break
			}
			if fn != nil && tempDevice != nil {
				fn(tempDevice)
			}
		}
	}
}

// SubscribeHostapd ubus subscribe hostapd.phy1-ap0 hostapd.phy0-ap0
func SubscribeHostapd(ctx context.Context, fn func(*HostapdDevice)) error {
	args := []string{"subscribe", "hostapd.*"}
	args = []string{"subscribe", "hostapd.phy1-ap0", "hostapd.phy0-ap0"}
	return Command(ctx, func(s string) {
		if s == "" {
			return
		}
		decode(s, fn)
	}, "ubus", args...)
}

func Monitor() error {
	return command(func(s string) {
		if s == "" {
			return
		}
		if strings.Contains(s, "5a:a7:22:62:3d:26") {
			glog.Info(s)
		}
	}, "ubus", "monitor")
}
