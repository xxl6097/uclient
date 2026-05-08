package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	dedupMu       sync.Mutex
	dedupLastSeen = make(map[string]time.Time)
)

// DedupDo 在 1 秒内相同 key 的事件只执行一次。
// 返回 true 表示本次执行，false 表示被去重忽略。
func DedupDo(key string, fn func()) bool {
	return DedupDoWithin(key, time.Second, fn)
}

// DedupDoWithin 在指定时间窗口内相同 key 的事件只执行一次。
func DedupDoWithin(key string, window time.Duration, fn func()) bool {
	now := time.Now()
	dedupMu.Lock()
	if last, ok := dedupLastSeen[key]; ok && now.Sub(last) < window {
		dedupMu.Unlock()
		return false
	}
	dedupLastSeen[key] = now
	for k, t := range dedupLastSeen {
		if now.Sub(t) > window*10 {
			delete(dedupLastSeen, k)
		}
	}
	dedupMu.Unlock()

	if fn != nil {
		fn()
	}
	return true
}

// GetDeviceID 获取 OpenWrt 设备唯一标识。
// 按优先级尝试多个来源，最终返回稳定的 sha256 摘要（16 字节十六进制）。
func GetDeviceID() string {
	if id := readDeviceIDRaw(); id != "" {
		sum := sha256.Sum256([]byte(id))
		return hex.EncodeToString(sum[:16])
	}
	return ""
}

// GetDeviceIDRaw 返回未经哈希处理的原始标识来源（调试用）。
func GetDeviceIDRaw() string {
	return readDeviceIDRaw()
}

func readDeviceIDRaw() string {
	// 1. machine-id
	if v := readTrim("/etc/machine-id"); v != "" {
		return "machine-id:" + v
	}
	if v := readTrim("/var/lib/dbus/machine-id"); v != "" {
		return "machine-id:" + v
	}

	// 2. OpenWrt 设备树序列号
	if v := readTrim("/sys/firmware/devicetree/base/serial-number"); v != "" {
		return "dts-serial:" + strings.TrimRight(v, "\x00")
	}

	// 3. /proc/cpuinfo Serial / Hardware
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		var serial, hardware string
		for _, line := range strings.Split(string(data), "\n") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			switch strings.ToLower(k) {
			case "serial":
				serial = v
			case "hardware", "machine":
				hardware = v
			}
		}
		if serial != "" && serial != "0000000000000000" {
			return "cpu-serial:" + serial
		}
		if hardware != "" {
			if mac := primaryMAC(); mac != "" {
				return "hw+mac:" + hardware + "|" + mac
			}
		}
	}

	// 4. ubus 获取 board 信息（OpenWrt 专属）
	if out, err := exec.Command("ubus", "call", "system", "board").Output(); err == nil {
		if strings.TrimSpace(string(out)) != "" {
			if mac := primaryMAC(); mac != "" {
				return "ubus+mac:" + mac
			}
		}
	}

	// 5. fallback: 首个物理网卡 MAC
	if mac := primaryMAC(); mac != "" {
		return "mac:" + mac
	}

	// 6. fallback: hostname
	if h, err := os.Hostname(); err == nil && h != "" {
		return "host:" + h
	}
	return ""
}

func readTrim(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// primaryMAC 返回稳定的物理网卡 MAC（按名称排序，跳过 lo/虚拟接口）。
func primaryMAC() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	var candidates []net.Interface
	for _, ifi := range ifaces {
		if ifi.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(ifi.HardwareAddr) == 0 {
			continue
		}
		name := strings.ToLower(ifi.Name)
		if strings.HasPrefix(name, "docker") || strings.HasPrefix(name, "veth") ||
			strings.HasPrefix(name, "br-") || strings.HasPrefix(name, "vmnet") ||
			strings.HasPrefix(name, "tun") || strings.HasPrefix(name, "tap") ||
			strings.HasPrefix(name, "wg") {
			continue
		}
		candidates = append(candidates, ifi)
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})
	if len(candidates) > 0 {
		return candidates[0].HardwareAddr.String()
	}
	return ""
}
