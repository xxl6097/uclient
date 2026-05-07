package main

import (
	"fmt"
	"net"

	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/u"
	"github.com/xxl6097/uclient/pkg"
)

func main() {
	a := &openwrt.DHCPLease{
		MAC: "aaaaa",
		IP:  net.ParseIP("127.0.0.1").String(),
	}
	fmt.Printf("--->%p %+v\n", a, a)
	b := u.DeepCopyGob[openwrt.DHCPLease](a)
	fmt.Printf("--->%p %+v\n", b, b)

	u.AppandText("./a.txt", "aaaaaa")
	u.AppandText("./a.txt", "bbbbbbb")

	for {
		pkg.DedupDo("event:foo", func() {
			// 1 秒内重复触发只会执行一次
			fmt.Printf("--->%p %+v\n", a, a)
		})
	}
}
