package main

import (
	"github.com/xxl6097/uclient/internal/openwrt"
)

func tee() {
	//var hostapdCmd = "ubus subscribe hostapd.*"
	_ = openwrt.SubscribeHostapd()
}

func main() {
	tee()
}
