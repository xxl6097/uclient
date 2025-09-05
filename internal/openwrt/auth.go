package openwrt

import (
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-http/pkg/util"
	"github.com/xxl6097/uclient/internal/auth"
)

func (this *openWRT) LoadAuth() {
	codes, err := auth.GetAuthData()
	if err == nil && codes != nil && len(codes) > 0 {
		this.authcode = codes
	}
}

func (this *openWRT) CheckAuth(authcode string) bool {
	glog.Info("Basic Auth:", authcode, this.authcode)
	if this.authcode != nil {
		if util.Contains1[string](this.authcode, authcode) {
			glog.Info("Basic Auth OK:", authcode, this.authcode)
			return true
		}
	}
	return false
}
