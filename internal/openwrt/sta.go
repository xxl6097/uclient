package openwrt

import (
	"context"
	"encoding/json"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"time"
)

//	{
//	   "staUpDown": {
//	       "online": 1,
//	       "staType": "Phone",
//	       "radio": "5G",
//	       "channel": "52",
//	       "ssid": "abcc-5G",
//	       "hostName": "xiaomi15",
//	       "macAddress": "16006F8335E1",
//	       "vmacAddress": "16006F8335E1",
//	       "staVendor": "",
//	       "accessTime": "2025-08-01 08:12:13",
//	       "rssi": "-39",
//	       "rxRate": "2161",
//	       "txRate": "2401"
//	   }
//	}
type StaUpDown struct {
	Online      int       `json:"online"` //0离线，1上线
	MacAddress  string    `json:"macAddress"`
	VmacAddress string    `json:"vmacAddress"`
	RxRate      string    `json:"rxRate"`
	TxRate      string    `json:"txRate"`
	Ssid        string    `json:"ssid"`
	Rssi        string    `json:"rssi"`
	HostName    string    `json:"hostName"`
	StaVendor   string    `json:"staVendor"` //品牌
	AccessTime  string    `json:"accessTime"`
	Timestamp   time.Time `json:"timestamp"`
	Radio       string    `json:"radio"`
	StaType     string    `json:"staType"`
}
type Sta struct {
	StaUpDown *StaUpDown `json:"staUpDown"`
}

func SubscribeSta(ctx context.Context, fn func(*StaUpDown)) error {
	return Command(ctx, func(s string) {
		if s == "" {
			return
		}
		var tempData Sta
		err := json.Unmarshal([]byte(s), &tempData)
		if err == nil && tempData.StaUpDown != nil {
			tempData.StaUpDown.Timestamp = glog.Now()
			if fn != nil {
				fn(tempData.StaUpDown)
			}
		}
	}, "ubus", "subscribe", "ahsapd.sta")
}

// GetStaInfo ubus call ahsapd.sta getStaInfo
func getStaInfo() *u.StaInfo {
	data, err := RunCMD("ubus", "call", "ahsapd.sta", "getStaInfo")
	if err != nil {
		glog.Errorf("Get sta info error: %v", err)
		return nil
	}
	if data == nil {
		glog.Errorf("Get sta info error")
		return nil
	}
	sta := u.StaInfo{}
	e := json.Unmarshal(data, &sta)
	if e != nil {
		glog.Errorf("Get sta info error: %v", e)
		return nil
	}
	return &sta
}

func GetStaInfo() map[string]*u.StaDevice {
	data := getStaInfo()
	if data == nil {
		return nil
	}
	if data.AhsapdSta.StaDevices == nil {
		return nil
	}
	devices := make(map[string]*u.StaDevice)
	for _, device := range data.AhsapdSta.StaDevices {
		mac := u.MacFormat(device.MacAddress)
		if mac != "" {
			devices[mac] = &device
		}
	}
	return devices
}
