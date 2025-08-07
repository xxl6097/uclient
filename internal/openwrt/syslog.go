package openwrt

import (
	"context"
	"github.com/xxl6097/glog/glog"
	"os"
)

func subscribeSysLogs(ctx context.Context, exitFun func(process *os.Process), fn func(string)) error {
	glog.Debug("subscribeSysLogs...")
	return Command(ctx, exitFun, func(s string) {
		if fn != nil {
			fn(s)
		}
	}, "logread", "-f")
}
