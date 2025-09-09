package internal

import (
	"github.com/xxl6097/go-http/pkg/httpserver"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	assets "github.com/xxl6097/uclient/assets/openwrt"
	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/u"
	"net/http"
	"net/url"
)

func Bootstrap(cfg *u.Config, service igs.Service) {
	server := httpserver.New().
		CORSMethodMiddleware().
		AddRoute(NewRoute(NewApi(service, cfg.Username, cfg.Password))).
		AddRoute(assets.NewRoute()).
		//BasicAuth(cfg.Username, cfg.Password, "oIin3168TLKg1X8OU2xBBWLlMEdI").
		BasicAuthFunc(cfg.Username, cfg.Password, func(r *http.Request) bool {
			//glog.Info("Basic Auth:", r.URL.String())
			autoCode := r.URL.Query().Get("auth_code")
			if autoCode == "" {
				query, err := url.Parse(r.Referer())
				if err == nil && query != nil {
					if query.Query().Has("auth_code") {
						autoCode = query.Query().Get("auth_code")
					}
				}
			}
			return openwrt.GetInstance().CheckAuth(autoCode)
		}).
		Done(cfg.ServerPort)
	defer server.Stop()
	server.Wait()
}
