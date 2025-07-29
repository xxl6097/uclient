package main

import (
	"fmt"
	"github.com/xxl6097/uclient/internal/openwrt"
)

func main() {
	ips, err := openwrt.GetInstance().GetStaticIpMap()
	if err != nil {
		return
	}
	if ips == nil || len(ips) <= 0 {
		return
	}

	for _, ip := range ips {
		fmt.Printf("%+v\n", ip)
	}
}
