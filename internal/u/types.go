package u

type NtfyInfo struct {
	Address  string `json:"address"`
	Topic    string `json:"topic"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//	type NtfyEventData struct {
//		Id    string `json:"id,omitempty"`
//		Time  int    `json:"time,omitempty"`
//		Event string `json:"event,omitempty"`
//		Topic string `json:"topic,omitempty"`
//	}

const (
	OPEN      = "open"      // 0
	KEEPALIVE = "keepalive" // 1
	MESSAGE   = "message"   // 1
)

type NtfyEventData struct {
	Id       string `json:"id,omitempty"`
	Time     int64  `json:"time,omitempty"`
	Expires  int64  `json:"expires,omitempty"`
	Event    string `json:"event,omitempty"`
	Topic    string `json:"topic,omitempty"`
	Title    string `json:"title,omitempty"`
	Message  string `json:"message,omitempty"`
	Markdown bool   `json:"markdown,omitempty"`
}

type StaDevice struct {
	IpAddress            string `json:"ipAddress"`
	HostName             string `json:"hostName"`
	MacAddress           string `json:"macAddress"`
	VmacAddress          string `json:"vmacAddress"`
	UpTime               string `json:"upTime"`
	AccessTime           string `json:"accessTime"`
	Rssi                 string `json:"rssi"`
	RxRate               string `json:"rxRate"`
	TxRate               string `json:"txRate"`
	RxRateRt             string `json:"rxRate_rt"`
	TxRateRt             string `json:"txRate_rt"`
	TotalPacketsSent     int    `json:"totalPacketsSent"`
	TotalPacketsReceived int    `json:"totalPacketsReceived"`
	TotalBytesSent       int    `json:"totalBytesSent"`
	TotalBytesReceived   int    `json:"totalBytesReceived"`
	RateArray            []struct {
		MaxRxRate  string `json:"maxRxRate"`
		MaxTxRate  string `json:"maxTxRate"`
		AverRxRate string `json:"averRxRate"`
		AverTxRate string `json:"averTxRate"`
		Interface  string `json:"interface"`
	} `json:"rateArray"`
	StaType   string `json:"staType"`
	Radio     string `json:"radio"`
	Channel   string `json:"channel"`
	Ssid      string `json:"ssid"`
	StaVendor string `json:"staVendor,omitempty"`
}
type StaInfo struct {
	AhsapdSta struct {
		StaDevices []StaDevice `json:"staDevices"`
	} `json:"ahsapd.sta"`
}

type Device struct {
	Lan        string `json:"lan"`
	UpTime     int    `json:"UpTime"`
	RSSI       int    `json:"RSSI"`
	TxDataRate int    `json:"TxDataRate"`
	RxDataRate int    `json:"RxDataRate"`
	TxBytesRt  string `json:"TxBytes_rt"`
	RxBytesRt  string `json:"RxBytes_rt"`
	TxBytes    int    `json:"TxBytes"`
	RxBytes    int    `json:"RxBytes"`
	TxPkts     int    `json:"TxPkts"`
	RxPkts     int    `json:"RxPkts"`
	MacAddress string `json:"MacAddress"`
	IPADDR     string `json:"IPADDR"`
	HOSTNAME   string `json:"HOSTNAME"`
}

type PORTINFO struct {
	PORT0 *Device `json:"PORT0"`
	PORT1 *Device `json:"PORT1"`
	PORT2 *Device `json:"PORT2"`
}
type Ra struct {
	SSID    string    `json:"SSID"`
	NUMBER  int       `json:"NUMBER"`
	Stainfo []*Device `json:"stainfo"`
}

type G2 struct {
	Ra0    *Ra `json:"ra0"`
	Ra1    *Ra `json:"ra1"`
	NUMBER int `json:"NUMBER"`
}
type G5 struct {
	Rax0   *Ra `json:"rax0"`
	NUMBER int `json:"NUMBER"`
}

type DeviceStatus struct {
	MEM        int       `json:"MEM"`
	UPTIME     int       `json:"UPTIME"`
	CPU        int       `json:"CPU"`
	WANMAC     string    `json:"WANMAC"`
	OPMODE     string    `json:"OPMODE"`
	WANIP      string    `json:"WANIP"`
	NETMASK    string    `json:"NETMASK"`
	GATEWAY    string    `json:"GATEWAY"`
	DNS        string    `json:"DNS"`
	WANUPTIME  int       `json:"WAN_UPTIME"`
	PROTO      string    `json:"PROTO"`
	IPV6ENABLE int       `json:"IPV6ENABLE"`
	NETSTATUS  int       `json:"NETSTATUS"`
	PORTINFO   *PORTINFO `json:"PORTINFO"`
	G2         *G2       `json:"2G"`
	G5         *G5       `json:"5G"`
}
