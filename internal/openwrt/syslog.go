package openwrt

import (
	"context"
	"github.com/xxl6097/glog/pkg/z"
)

func subscribeSysLogs(ctx context.Context, fn func(string)) error {
	z.Debug("subscribeSysLogs...")
	return Command(ctx, func(s string) {
		if fn != nil {
			fn(s)
		}
	}, "logread", "-f")
}
