package openwrt

import "github.com/xxl6097/glog/glog"

func subscribeSysLogs(fn func(string)) error {
	glog.Debug("subscribeSysLogs...")
	return command(func(s string) {
		if fn != nil {
			fn(s)
		}
	}, "logread", "-f")
}
