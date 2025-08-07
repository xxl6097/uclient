package openwrt

import (
	"context"
	"github.com/xxl6097/glog/glog"
)

func subscribeSysLogs(ctx context.Context, fn func(string)) error {
	glog.Debug("subscribeSysLogs...")
	return Command(ctx, func(s string) {
		if fn != nil {
			fn(s)
		}
	}, "logread", "-f")
}
