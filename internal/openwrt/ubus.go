package openwrt

var (
	//强制断开客户端
	//addr：客户端MAC地址
	//reason：断开原因代码（如5表示协议错误）
	//deauth：是否发送解认证帧
	offline = "ubus call hostapd.* del_client '{\"addr\":\"5a:a7:22:62:3d:26\", \"reason\":5, \"deauth\":true}'"
	//list    = "ubus call hostapd.* get_clients"
)

//ubus -v list hostapd.*
