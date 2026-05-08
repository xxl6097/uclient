package u

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func Ping(target string) bool {
	conn, _ := net.Dial("ip4:icmp", target)
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Body: &icmp.Echo{ID: 0, Seq: 1, Data: []byte("")},
	}
	data, _ := msg.Marshal(nil)

	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	_, err := conn.Write(data)
	if err != nil {
		return false
	}

	reply := make([]byte, 1500)
	_, err = conn.Read(reply)
	return err == nil
}

func ByteCountSpeed(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B/s", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB/s", float64(b)/float64(div), "KMGTPE"[exp])
}

// GetMacAddress 获取不带冒号的 MAC 地址
func GetMacAddress(iface string) (string, error) {
	macPath := filepath.Join("/sys/class/net", iface, "address")

	file, err := os.Open(macPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		mac := strings.TrimSpace(scanner.Text())
		// 去掉所有冒号
		mac = strings.ReplaceAll(mac, ":", "")
		return mac, nil
	}

	return "", fmt.Errorf("读取MAC失败")
}

func GetEth0Mac1() string {
	iface := "eth0"
	mac, err := GetMacAddress(iface)
	if err != nil {
		fmt.Printf("获取失败: %v\n", err)
		return ""
	}
	return mac
}

// ParseDeviceStr 解析 0a5e697602d5-getList 格式字符串
// 返回: mac地址, 指令, 错误
func ParseDeviceStr(s string) (string, string, error) {
	// 按 "-" 分割成两部分
	parts := strings.Split(s, "-")

	// 校验格式必须是 [mac]-[cmd]
	if len(parts) != 2 {
		return "", "", fmt.Errorf("格式错误，必须是 无冒号MAC-指令 格式")
	}

	mac := parts[0]
	cmd := parts[1]

	// 可选：校验MAC长度（12位）
	if len(mac) != 12 {
		return "", "", fmt.Errorf("MAC长度错误，必须是12位无冒号格式")
	}

	return mac, cmd, nil
}

// ParseDeviceSnStr 解析 0a5e697602d5-getList 格式字符串
// 返回: deviceSn地址, 指令, 错误
func ParseDeviceSnStr(s string) (string, string, error) {
	// 按 "-" 分割成两部分
	parts := strings.Split(s, "-")

	// 校验格式必须是 [mac]-[cmd]
	if len(parts) != 2 {
		return "", "", fmt.Errorf("格式错误，必须是 无冒号MAC-指令 格式")
	}

	mac := parts[0]
	cmd := parts[1]

	return mac, cmd, nil
}
