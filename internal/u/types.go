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
