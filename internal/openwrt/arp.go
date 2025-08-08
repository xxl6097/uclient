package openwrt

import (
	"context"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var tempMap = make(map[string]*ARPEntry)

func parseARPLine(line string) (*ARPEntry, error) {
	fields := strings.Fields(line)
	if len(fields) < 6 {
		return nil, fmt.Errorf("invalid ARP line: expected 6 fields, got %d", len(fields))
	}

	// 解析 IP 地址
	ip := net.ParseIP(fields[0])
	if ip == nil {
		return nil, fmt.Errorf("invalid IP: %s", fields[0])
	}

	// 解析十六进制数值（HWType 和 Flags）
	hwType, _ := strconv.ParseUint(strings.TrimPrefix(fields[1], "0x"), 16, 8)
	flags, _ := strconv.ParseUint(strings.TrimPrefix(fields[2], "0x"), 16, 8)

	// 解析 MAC 地址
	mac, err := net.ParseMAC(fields[3])
	if err != nil {
		return nil, fmt.Errorf("invalid MAC: %v", err)
	}
	if mac.String() == "00:00:00:00:00:00" {
		return nil, fmt.Errorf("error MAC")
	}

	return &ARPEntry{
		IP:        ip,
		HWType:    uint8(hwType),
		Flags:     uint8(flags),
		MAC:       mac,
		Mask:      fields[4],
		Interface: fields[5],
	}, nil
}

func getClientsByArp(deviceInterfaceName string) (map[string]*ARPEntry, error) {
	data, err := os.ReadFile(arpFilePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	//glog.Debugf("\n%s", content)
	lines := strings.Split(content, "\n")
	entries := make(map[string]*ARPEntry)
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue // 跳过标题行和空行
		}
		if !strings.HasSuffix(line, deviceInterfaceName) {
			continue // 根据Device过滤
		}
		entry, e := parseARPLine(line)
		if e != nil {
			//return nil, err
			//glog.Error("parseARPLine error", e, line)
			continue
		}
		mac := entry.MAC.String()
		if _, ok := entries[mac]; ok {
			temp := entries[mac]
			if temp.Flags != entry.Flags && entry.Flags == 2 {
				entries[mac] = entry
			}
		} else {
			entries[mac] = entry
		}
	}
	return entries, nil
}

// 比较两个ARP表，返回新增、删除和修改的条目
func compareARPTables(new map[string]*ARPEntry) (added, changed []*ARPEntry) {
	// 检查新增和修改的条目
	for mac, newEntry := range new {
		oldEntry, exists := tempMap[mac]
		if !exists {
			added = append(added, newEntry)
		} else if !entriesEqual(oldEntry, newEntry) {
			changed = append(changed, newEntry)
		}
		tempMap[mac] = newEntry
	}
	return
}

// 检查两个ARP条目是否相等
func entriesEqual(a, b *ARPEntry) bool {
	return a.IP.String() == b.IP.String() &&
		a.HWType == b.HWType &&
		a.Flags == b.Flags &&
		a.MAC.String() == b.MAC.String() &&
		a.Mask == b.Mask &&
		a.Interface == b.Interface
}

func SubscribeArp(interval time.Duration, fn func(entry *ARPEntry)) error {
	glog.Println("开始监控/proc/net/arp变化...")
	// 获取初始ARP表
	prevEntries, err := getClientsByArp(brLanString)
	if err != nil {
		err = fmt.Errorf("读取ARP表失败: %v\n", err)
		return err
	}
	tempMap = prevEntries
	// 打印初始ARP表
	glog.Println("初始ARP表:")
	for _, entry := range prevEntries {
		glog.Printf("%-15v %-6v %-6v %-17v %-6v %s\n",
			entry.IP.String(), entry.HWType, entry.Flags, entry.MAC.String(), entry.Mask, entry.Interface)
	}
	// 持续监控
	for {
		time.Sleep(interval)
		// 获取当前ARP表
		currentEntries, e1 := getClientsByArp(brLanString)
		if e1 != nil {
			glog.Printf("读取ARP表失败: %v\n", err)
			continue
		}

		// 比较差异
		added, changed := compareARPTables(currentEntries)

		// 输出变化
		if len(added) > 0 || len(changed) > 0 {
			glog.Printf("\n[%s] 检测到ARP表变化:\n", time.Now().Format(time.DateTime))
		}

		for _, entry := range added {
			glog.Printf("± 新增   %-15s MAC: %-17s 设备: %s Flags: %v", entry.IP, entry.MAC, entry.Interface, entry.Flags)
			if fn != nil {
				fn(entry)
			}
		}

		for _, entry := range changed {
			glog.Printf("± 修改   %-15s MAC: %-17s 设备: %s Flags: %v", entry.IP, entry.MAC, entry.Interface, entry.Flags)
			if fn != nil {
				fn(entry)
			}
		}
		// 更新上一次的ARP表
		prevEntries = currentEntries
	}
}

func SubscribeArpCache(ctx context.Context, interval time.Duration, fn func(entry map[string]*ARPEntry)) {
	for {
		select {
		case <-ctx.Done():
			glog.Debug("Arp监听退出...")
			return
		default:
			time.Sleep(interval)
			// 获取当前ARP表
			currentEntries, e1 := getClientsByArp(brLanString)
			if e1 != nil {
				glog.Errorf("读取ARP表失败: %v\n", e1)
				continue
			}
			if fn != nil {
				fn(currentEntries)
			}
		}
	}
}
