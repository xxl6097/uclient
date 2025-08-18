package leak

import (
	"fmt"
	"github.com/xxl6097/go-http/pkg/httpserver"
	assets "github.com/xxl6097/uclient/assets/openwrt"
	"github.com/xxl6097/uclient/internal"
	"github.com/xxl6097/uclient/pkg"
	"go.uber.org/goleak"
	"net/http"
	"os"
	"testing"
)

// go test -v -run TestLeakyEndpoint
func TestHealthyEndpoint(t *testing.T) {
	defer goleak.VerifyNone(t) // 泄漏检测
	pkg.AppVersion = "v0.0.3"
	pkg.BinName = "openwrt-client-manager_v0.0.20_darwin_arm64"
	fmt.Println("Hello World", os.Getpid())
	router := httpserver.New().
		CORSMethodMiddleware().
		BasicAuth("admin", "admin").
		AddRoute(internal.NewRoute(internal.NewApi(nil, "admin", "admin"))).
		AddRoute(assets.NewRoute())
	//router.Handle("/metrics", promhttp.Handler())
	server := router.Done(8081)
	defer server.Stop()

	// 发送健康检查请求
	//resp, err := http.Get("http://localhost" + server.server.Addr + "/healthy")
	resp, err := http.Get("http://localhost:8081")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m) // 全局泄漏检测
}
