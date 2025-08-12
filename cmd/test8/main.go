package main

import (
	"fmt"
	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/u"
	"net"
)

func main() {
	a := &openwrt.DHCPLease{
		MAC: "aaaaa",
		IP:  net.ParseIP("127.0.0.1").String(),
	}
	fmt.Printf("--->%p %+v\n", a, a)
	b := u.DeepCopyGob[openwrt.DHCPLease](a)
	fmt.Printf("--->%p %+v\n", b, b)
}
