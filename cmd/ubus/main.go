package main

import (
	"fmt"
	"github.com/xxl6097/uclient/internal/openwrt"
	"strings"
)

func main() {
	list := openwrt.UbusList()
	fmt.Println(list)
	if strings.Contains(list, "hostapd") {
		fmt.Println("hostapd存在")
	}
	if strings.Contains(list, "ahsapd.sta") {
		fmt.Println("ahsapd.sta存在")
	}
}
