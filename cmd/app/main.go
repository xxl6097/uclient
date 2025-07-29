package main

import (
	"fmt"
	"github.com/xxl6097/go-http/pkg/httpserver"
	assets "github.com/xxl6097/uclient/assets/openwrt"
	"github.com/xxl6097/uclient/internal"
	"github.com/xxl6097/uclient/internal/u"
	"github.com/xxl6097/uclient/pkg"
)

func init() {
	if u.IsMacOs() {
		pkg.AppVersion = "v0.0.3"
		pkg.BinName = "openwrt-client-manager_v0.0.20_darwin_arm64"
	}
}
func main() {
	fmt.Println("Hello World")
	httpserver.New().
		CORSMethodMiddleware().
		BasicAuth("admin", "admin").
		AddRoute(internal.NewRoute(internal.NewApi(nil, "admin", "admin"))).
		AddRoute(assets.NewRoute()).
		Done(7000)
}
