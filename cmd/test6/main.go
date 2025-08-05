package main

import (
	"fmt"
	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/u"
)

type Roaming struct {
	AhsapdRoaming struct {
		Radio   string `json:"radio"`
		Rssi    string `json:"rssi"`
		Channel string `json:"channel"`
	} `json:"ahsapd.roaming"`
}

var js = "{\"MEM\":93,\"UPTIME\":442905,\"CPU\":3,\"WANMAC\":\"405EE15456AC\",\"OPMODE\":\"router\",\"WANIP\":\"10.6.50.76\",\"NETMASK\":\"255.255.254.0\",\"GATEWAY\":\"10.6.50.254\",\"DNS\":\"114.114.114.114\",\"WAN_UPTIME\":442885,\"PROTO\":\"dhcp\",\"IPV6ENABLE\":1,\"NETSTATUS\":1,\"ifaRate\":{\"AverRxRate\":\"0.0000\",\"AverTxRate\":\"0.0000\",\"MaxRxRate\":\"0.0000\",\"MaxTxRate\":\"0.0000\"},\"if6Rate\":{\"AverRxRate\":\"0.0000\",\"AverTxRate\":\"0.0000\",\"MaxRxRate\":\"0.0000\",\"MaxTxRate\":\"0.0000\"},\"PORTINFO\":{\"PORT0\":{\"LINKSTATUS\":2,\"MacAddress\":\"1C20DB8D9676\",\"UpTime\":442905,\"TxBytes_rt\":\"1.1193\",\"RxBytes_rt\":\"48.2314\",\"TxBytes\":11382171966,\"RxBytes\":238134789493},\"PORT1\":{\"LINKSTATUS\":3,\"MacAddress\":\"8CEC4B588109\",\"IPADDR\":\"192.168.1.3\",\"HOSTNAME\":\"fnos\",\"UpTime\":442831,\"TxBytes_rt\":\"0.3098\",\"RxBytes_rt\":\"0.2008\",\"TxBytes\":23137672506,\"RxBytes\":8149320224,\"TxPkts\":30622496,\"RxPkts\":31123565,\"ifaRate\":{\"AverRxRate\":\"0.0000\",\"AverTxRate\":\"0.0000\",\"MaxRxRate\":\"0.0000\",\"MaxTxRate\":\"0.0000\"},\"if6Rate\":{\"AverRxRate\":\"0.0000\",\"AverTxRate\":\"0.0000\",\"MaxRxRate\":\"0.0000\",\"MaxTxRate\":\"0.0000\"}},\"PORT2\":{\"LINKSTATUS\":0}},\"2G\":{\"ra0\":{\"SSID\":\"abcc\",\"NUMBER\":0,\"stainfo\":[]},\"ra1\":{\"SSID\":\"CMCC-GUIDE-LINK\",\"NUMBER\":0,\"stainfo\":[]},\"NUMBER\":0},\"5G\":{\"rax0\":{\"SSID\":\"abcc-5G\",\"NUMBER\":2,\"stainfo\":[{\"UpTime\":20422,\"RSSI\":-37,\"TxDataRate\":2401,\"RxDataRate\":2401,\"TxBytes_rt\":\"8.2274\",\"RxBytes_rt\":\"0.1898\",\"TxBytes\":1721712800,\"RxBytes\":722741188,\"TxPkts\":2243840,\"RxPkts\":1167747,\"MacAddress\":\"92D75B753D9A\",\"IPADDR\":\"192.168.1.2\",\"HOSTNAME\":\"M4\",\"ifaRate\":{\"AverRxRate\":\"0.0000\",\"AverTxRate\":\"0.0000\",\"MaxRxRate\":\"0.0000\",\"MaxTxRate\":\"0.0000\"},\"if6Rate\":{\"AverRxRate\":\"0.0000\",\"AverTxRate\":\"0.0000\",\"MaxRxRate\":\"0.0000\",\"MaxTxRate\":\"0.0000\"}},{\"UpTime\":1409,\"RSSI\":-29,\"TxDataRate\":2401,\"RxDataRate\":1921,\"TxBytes_rt\":\"15.5783\",\"RxBytes_rt\":\"0.5379\",\"TxBytes\":3990957029,\"RxBytes\":78857129,\"TxPkts\":3232236,\"RxPkts\":832568,\"MacAddress\":\"16006F8335E1\",\"IPADDR\":\"192.168.1.4\",\"HOSTNAME\":\"xiaomi15\",\"ifaRate\":{\"AverRxRate\":\"0.0000\",\"AverTxRate\":\"0.0000\",\"MaxRxRate\":\"0.0000\",\"MaxTxRate\":\"0.0000\"},\"if6Rate\":{\"AverRxRate\":\"0.0000\",\"AverTxRate\":\"0.0000\",\"MaxRxRate\":\"0.0000\",\"MaxTxRate\":\"0.0000\"}}]},\"NUMBER\":2}}"

// RxBytes 上传
// TxBytes 下载
func main() {

	a := 22357088240
	b := 22357470636
	speed1 := float64(b-a) / float64(5*1000)
	fmt.Println(speed1) // 输出：你好

	c := 7867738418
	d := 7868179949
	speed2 := float64(d-c) / float64(5)
	fmt.Println(u.ByteCountSpeed(uint64(speed2)))
	fmt.Println(speed2) // 输出：你好

	da := openwrt.GetInstance().DecodeDevice([]byte(js))
	for _, d := range da {
		fmt.Println(d)
	}
}
