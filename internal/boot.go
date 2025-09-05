package internal

import (
	"github.com/xxl6097/go-http/pkg/httpserver"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	assets "github.com/xxl6097/uclient/assets/openwrt"
	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/u"
	"net/http"
)

func Bootstrap(cfg *u.Config, service igs.Service) {
	server := httpserver.New().
		CORSMethodMiddleware().
		AddRoute(NewRoute(NewApi(service, cfg.Username, cfg.Password))).
		AddRoute(assets.NewRoute()).
		//BasicAuth(cfg.Username, cfg.Password, "oIin3168TLKg1X8OU2xBBWLlMEdI").
		BasicAuthFunc(cfg.Username, cfg.Password, func(r *http.Request) bool {
			return openwrt.GetInstance().CheckAuth(r.URL.Query().Get("auth_code"))
		}).
		Done(cfg.ServerPort)
	defer server.Stop()
	server.Wait()
}
