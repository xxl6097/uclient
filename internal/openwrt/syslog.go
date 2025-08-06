package openwrt

func subscribeSysLogs(fn func(string)) error {
	return command(func(s string) {
		if fn != nil {
			fn(s)
		}
	}, "logread", "-f")
}
