package main

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-http/pkg/httpserver"
	"github.com/xxl6097/go-service/pkg/gs"
	assets "github.com/xxl6097/uclient/assets/openwrt"
	"github.com/xxl6097/uclient/cmd/app/service"
	"github.com/xxl6097/uclient/internal"
	"github.com/xxl6097/uclient/internal/u"
	"github.com/xxl6097/uclient/pkg"
	"os"
)

//func prome() http.Handler {
//	prometheus.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
//		Name: "go_memory_usage_bytes",
//		Help: "Current memory usage",
//	}, func() float64 {
//		var m runtime.MemStats
//		runtime.ReadMemStats(&m)
//		return float64(m.HeapAlloc)
//	}))
//	return promhttp.Handler()
//}

func init() {
	//go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()
	if u.IsMacOs() {
		pkg.AppVersion = "v0.0.3"
		pkg.BinName = "openwrt-client-manager_v0.0.20_darwin_arm64"
		fmt.Println("Hello World", os.Getpid())
		router := httpserver.New().
			CORSMethodMiddleware().
			BasicAuth("admin", "admin").
			AddRoute(internal.NewRoute(internal.NewApi(nil, "admin", "admin"))).
			AddRoute(assets.NewRoute())
		//router.Handle("/metrics", promhttp.Handler())
		server := router.Done(7000)
		defer server.Stop()
		server.Wait()
	}
}

func main() {
	//defer glog.GlobalRecover()
	s := service.Service{}
	err := gs.Run(&s)
	glog.Debug("程序结束", err)
}
