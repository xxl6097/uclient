package openwrt

import (
	"encoding/json"
	"github.com/xxl6097/glog/glog"
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

func SubscribeSta(fn func(*StaUpDown)) error {
	return command(func(s string) {
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
