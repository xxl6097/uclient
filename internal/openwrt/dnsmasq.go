package openwrt

import (
	"encoding/json"
	"github.com/xxl6097/glog/glog"
	"time"
)

var dnsmasq = "ubus subscribe dnsmasq"

type DnsmasqDevice struct {
	Mac       string    `json:"mac"`
	Ip        string    `json:"ip"`
	Name      string    `json:"name"`
	Interface string    `json:"interface"`
	Timestamp time.Time `json:"timestamp"`
}
type Dnsmasq struct {
	DhcpAck *DnsmasqDevice `json:"dhcp.ack"`
}

func SubscribeDnsmasq(fn func(*DnsmasqDevice)) error {
	return command(func(s string) {
		if s == "" {
			return
		}
		var tempData Dnsmasq
		err := json.Unmarshal([]byte(s), &tempData)
		if err == nil && tempData.DhcpAck != nil {
			if tempData.DhcpAck != nil {
				tempData.DhcpAck.Timestamp = glog.Now()
			}
			if fn != nil {
				fn(tempData.DhcpAck)
			}
		}
	}, "ubus", "subscribe", "dnsmasq")
}
