package openwrt

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/uclient/internal/u"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	arpFilePath           = "/proc/net/arp"
	dhcpLeasesFilePath    = "/tmp/dhcp.leases"
	hetsysinfoFilePath    = "/tmp/hetsysinfo.json"
	dhcpCfgFilePath       = "/etc/config/dhcp"
	brLanString           = "br-lan"
	apStaDisConnectString = "AP-STA-DISCONNECTED"
	apStaConnectString    = "AP-STA-CONNECTED"
	MAX_SIZE              = 12000
	MAX_WORK_SIZE         = 3600
	RE_REY_MAX_COUNT      = 5
)
var (
	StatusDir       = "/etc/config/uclient/status"
	workDir         = "/etc/config/uclient/work"
	nickFilePath    = "/etc/config/uclient/nick"
	webhookFilePath = "/etc/config/uclient/webhook"
	ntfyFilePath    = "/etc/config/uclient/ntfy"
)

type SysLogEvent struct {
	Online    bool      `json:"online"`
	Timestamp time.Time `json:"timestamp"`
	Mac       string    `json:"mac"`
	Phy       string    `json:"phy"`
	RawData   string    `json:"-"`
}
type DeviceTimeLine struct {
	DateTime  string `json:"dateTime"`
	Timestamp int64  `json:"timestamp"`
	Connected bool   `json:"connected"`
	Ago       string `json:"ago"`
}

type Status struct {
	Timestamp int64 `json:"timestamp"`
	Connected bool  `json:"connected"`
}
type NickEntry struct {
	Name      string           `json:"name"`
	IsPush    bool             `json:"isPush"`
	MAC       string           `json:"mac"`
	IP        string           `json:"ip"`
	StartTime int64            `json:"starTime"`
	Hostname  string           `json:"hostname"`
	WorkType  *WorkTypeSetting `json:"workType"`
}
type DHCPLease struct {
	IP        string     `json:"ip"`  //DHCP 服务器分配给客户端的 IP
	MAC       string     `json:"mac"` //设备的物理地址，格式为 xx:xx:xx:xx:xx:xx
	Phy       string     `json:"phy"`
	Hostname  string     `json:"hostname"` //客户端上报的主机名（可能为空或 *）
	StartTime int64      `json:"starTime"` //租约失效的精确时间（秒级精度）
	Online    bool       `json:"online"`
	Signal    int        `json:"signal"`
	Freq      int        `json:"freq"`
	StaType   string     `json:"staType"`
	Ssid      string     `json:"ssid"`
	UpRate    string     `json:"upRate,omitempty"`
	DownRate  string     `json:"downRate,omitempty"`
	Device    *u.Device  `json:"device,omitempty"`
	Nick      *NickEntry `json:"nick"` //
	Static    *DHCPHost  `json:"static"`
}
type ARPEntry struct {
	IP        net.IP           //设备的 IPv4 地址
	HWType    uint8            // 硬件类型（通常为 0x1，表示以太网）。
	Flags     uint8            // ARP 表项状态标志：0x0：无效（离线）;0x2：有效（在线），表示设备可达。
	MAC       net.HardwareAddr //设备的 MAC 地址
	Mask      string           // 子网掩码（通常为 *，表示未使用）
	Interface string           // 关联的网络接口（如 br-lan、eth0）
	Timestamp time.Time
}

func getDataFromSysLog(pattern string, args ...string) (map[string][]DHCPLease, error) {
	dataMap := make(map[string][]DHCPLease)
	// 2. 编译正则表达式（匹配连接/断开事件）
	//pattern := `AP-STA-(CONNECTED|DISCONNECTED)`
	re := regexp.MustCompile(pattern)
	return dataMap, command(func(data string) {
		if re.MatchString(data) {
			//fmt.Println("[事件] ", data) // 输出匹配行
			macAddr := ParseMacAddr(data)
			mac, err := net.ParseMAC(macAddr)
			if err == nil {
				t, e := parseTimer(data)
				if e == nil {
					var status bool
					if strings.Contains(data, apStaConnectString) {
						status = true
					} else if strings.Contains(data, apStaDisConnectString) {
						status = false
					}
					element := DHCPLease{
						MAC:       mac.String(),
						StartTime: t.Unix(),
						Online:    status,
					}
					v, ok := dataMap[element.MAC]
					if ok {
						v = append(v, element)
					} else {
						v = []DHCPLease{element}
					}
					dataMap[element.MAC] = v
				}
			}
		}

	}, "logread", args...)
}

func subscribeSysLog(fn func(*SysLogEvent)) error {
	//args := []string{"-f", "|", "grep", "hostapd.*"}
	pattern := `hostapd.*`
	re := regexp.MustCompile(pattern)
	return command(func(s string) {
		if re.MatchString(s) {
			tempData := parseSysLog(s)
			if fn != nil && tempData != nil {
				fn(tempData)
			}
			//if status == 0 {
			//	if fn != nil {
			//		fn(timestamp, macAddr, phy, false)
			//	}
			//} else if status == 1 {
			//	if fn != nil {
			//		fn(timestamp, macAddr, phy, true)
			//	}
			//} else {
			//	//fmt.Printf("未知类型 %s\n", s)
			//}
		}

	}, "logread", "-f")
}

func getStatusFromSysLog() (map[string][]*Status, error) {
	pattern := `AP-STA-(CONNECTED|DISCONNECTED)`
	data, err := getDataFromSysLog(pattern)
	//fmt.Printf("GetDisconnectFromSysLog %+v\n", data)
	if err == nil && len(data) > 0 {
		times := make(map[string][]*Status)
		for k, v := range data {
			value := make([]*Status, 0)
			for _, entry := range v {
				value = append(value, &Status{
					//MAC:       entry.MAC,
					Timestamp: entry.StartTime,
					Connected: entry.Online,
				})
			}
			times[k] = value
		}
		return times, err
	}
	return nil, err
}

func getLeaseTime() time.Duration {
	data, err := os.ReadFile(dhcpCfgFilePath)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return time.Duration(3600 * 12)
	}

	// 正则匹配leasetime选项
	re := regexp.MustCompile(`option leasetime ['"]([\dhms]+)['"]`)
	match := re.FindStringSubmatch(string(data))
	if len(match) < 2 {
		fmt.Println("未找到leasetime配置")
		return time.Duration(3600 * 12)
	}

	// 解析时间字符串（如"12h"）
	leaseStr := match[1]
	leaseDuration, err := time.ParseDuration(leaseStr)
	if err != nil {
		fmt.Println("解析时间失败:", err)
		return time.Duration(3600 * 12)
	}
	//fmt.Printf("DHCP租约时间: %v（%d秒）\n", leaseStr, int(leaseDuration.Seconds()))
	return leaseDuration
}

func parseTime(logLine string) int64 {
	//re := regexp.MustCompile(`^(\w+\s+\w+\s+\d+\s+\d+:\d+:\d+\s+\d+)`)
	//matches := re.FindStringSubmatch(logLine)
	//if len(matches) > 1 {
	//	timeStr := matches[1]
	//	t, err := autoParse(timeStr)
	//	if err == nil {
	//		return t.Format(time.DateTime)
	//	}
	//}
	//return ""
	t, err := parseTimer(logLine)
	if err == nil {
		return t.Unix()
	}
	return 0
}

func parseTime1(logLine string) string {
	t, err := parseTimer(logLine)
	if err == nil {
		return t.Format(time.DateTime)
	}
	return ""
}

func ParseMacAddr(logLine string) string {
	// 1. 定义MAC地址正则表达式（兼容冒号/短横线分隔）
	pattern := `(?:[0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2}`
	re := regexp.MustCompile(pattern)
	// 2. 提取所有匹配的MAC地址
	macAddresses := re.FindAllString(logLine, -1)
	// 3. 输出结果
	if len(macAddresses) > 0 {
		return macAddresses[0]
	}
	return ""
}
func parseSysLog(data string) *SysLogEvent {
	phy := parsePhy(data)
	//timestamp := parseTime(data)
	macAddr := ParseMacAddr(data)
	//timestamp := glog.Now().UnixMilli() //time.Now().UnixMilli()
	timestamp := glog.Now()
	// 1. 检查字符串是否包含目标字段
	eve := SysLogEvent{
		Timestamp: timestamp,
		Mac:       macAddr,
		Phy:       phy,
	}
	if strings.Contains(data, apStaDisConnectString) { //AP-STA-DISCONNECTED
		eve.Online = false
		return &eve
	} else if strings.Contains(data, apStaConnectString) { //AP-STA-CONNECTED
		eve.Online = true
		return &eve
	}
	return nil
}

func parseLeaseLine(line string, leasetime time.Duration) (DHCPLease, error) {
	// 示例行: 1693837890 00:11:22:33:44:55 192.168.1.100 hostname-1
	fields := strings.Fields(line)
	if len(fields) < 4 { // 至少包含时间戳、MAC、IP、主机名
		return DHCPLease{}, fmt.Errorf("字段不足")
	}
	//fmt.Println("fields", fields)
	// 解析时间戳（Unix时间）
	now := glog.Now()
	var startTime = now
	startSec, e := strconv.ParseInt(fields[0], 10, 64)
	if e != nil {
		//beijingLoc, err := time.LoadLocation("Asia/Shanghai")
		//if err != nil {
		//	// 备选方案：手动创建东八区时区
		//	//fmt.Println(err, "备选方案：手动创建东八区时区")
		//	//beijingLoc = time.FixedZone("CST", 8*60*60) // UTC+8
		//	beijingLoc = time.FixedZone("UTC+8", 8*60*60)
		//}
		//utcTime := time.Unix(startSec, 0) //// 解析为 UTC 时间
		//startTime = utcTime.In(beijingLoc).Add(-leasetime)
		startTime = u.UTC8ToTime(startSec)
	}
	//startTime := time.Unix(startSec, 0) //// 解析为 UTC 时间
	if startTime.UnixMilli() > now.UnixMilli() {
		startTime = now
	}

	// 解析MAC地址
	mac, err := net.ParseMAC(fields[1])
	if err != nil {
		return DHCPLease{}, fmt.Errorf("MAC格式错误: %v", err)
	}

	// 解析IP地址
	ip := net.ParseIP(fields[2])
	if ip == nil {
		return DHCPLease{}, fmt.Errorf("IP格式错误")
	}

	// 主机名（可能包含空格，合并剩余字段）
	//hostname := strings.Join(fields[3:], " ")
	hostname := fields[3]

	return DHCPLease{
		IP:        ip.String(),
		MAC:       mac.String(),
		Hostname:  hostname,
		StartTime: startTime.UnixMilli(),
		//IsActive:  time.Now().Before(startTime.Add(time.Second * time.Duration(leaseDuration))),
	}, nil
}

func getNickData() (map[string]*NickEntry, error) {
	data, err := os.ReadFile(nickFilePath)
	if err != nil {
		return nil, err
	}
	dataMap := map[string]*NickEntry{}
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		return nil, err
	}
	return dataMap, nil
}

func setNickData(dataMap map[string]*NickEntry) error {
	if dataMap == nil || len(dataMap) == 0 {
		return nil
	}
	content, err := json.MarshalIndent(dataMap, "", "  ")
	//content, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}
	file, err := os.Create(nickFilePath) // 文件不存在则创建，存在则截断
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入内容
	_, err = file.Write(content)
	return err
}

func updateNicksData(dataMap map[string]*NickEntry) error {
	if dataMap == nil || len(dataMap) == 0 {
		return fmt.Errorf("dataMap is empty")
	}
	tempData, err := getNickData()
	if err != nil || tempData == nil {
		return setNickData(dataMap)
	}
	for k, v := range dataMap {
		tempData[k] = v
	}
	return setNickData(tempData)
}

func updateNickData(mac string, data *NickEntry) error {
	if data == nil {
		return fmt.Errorf("data is nil")
	}

	return updateNicksData(map[string]*NickEntry{mac: data})
}

func getArp(deviceInterfaceName string) ([]string, error) {
	data, err := os.ReadFile(arpFilePath)
	if err != nil {
		return nil, err
	}
	arpText := make([]string, 0)
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue // 跳过标题行和空行
		}
		if !strings.HasSuffix(line, deviceInterfaceName) {
			continue // 根据Device过滤
		}
		arpText = append(arpText, strings.TrimSpace(line))
	}
	return arpText, nil
}

func parsePhy(logLine string) string {
	re := regexp.MustCompile(`hostapd:\s+(phy[\w-]+?):`)
	// 提取匹配结果
	matches := re.FindStringSubmatch(logLine)
	if len(matches) < 2 {
		return ""
	}
	phyField := matches[1] // 捕获组索引为1
	return phyField
}

func parseArpLines(lines []string) (map[string]*ARPEntry, error) {
	if lines == nil || len(lines) == 0 {
		return nil, nil
	}
	entries := make(map[string]*ARPEntry)
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue // 跳过标题行和空行
		}
		entry, err := parseARPLine(line)
		if err != nil {
			return nil, err
		}
		entries[entry.MAC.String()] = entry
	}
	return entries, nil
}
func parseDHCPLeases(filePath string) (map[string]*DHCPLease, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	leaseTime := getLeaseTime()
	entries := make(map[string]*DHCPLease)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // 跳过注释和空行
		}
		lease, err := parseLeaseLine(line, leaseTime)
		if err != nil {
			log.Printf("解析失败: %v | 行: %s", err, line)
			continue
		}
		entries[lease.MAC] = &lease
	}
	return entries, nil
}

func getClientsByDhcp() (map[string]*DHCPLease, error) {
	return parseDHCPLeases(dhcpLeasesFilePath)
}

func parseTimer(logLine string) (*time.Time, error) {
	re := regexp.MustCompile(`^(\w+\s+\w+\s+\d+\s+\d+:\d+:\d+\s+\d+)`)
	matches := re.FindStringSubmatch(logLine)
	if len(matches) > 1 {
		timeStr := matches[1]
		t, err := u.AutoParse(timeStr)
		return t, err
	}
	return nil, nil
}

func Command(fu func(string), name string, arg ...string) error {
	return command(fu, name, arg...)
}

func command(fu func(string), name string, arg ...string) error {
	glog.Println(name, arg)
	// 创建ubus命令对象
	cmd := exec.Command(name, arg...)

	// 创建标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("创建管道失败: %v\n", err)
		return err
	}

	// 启动命令
	if e := cmd.Start(); e != nil {
		fmt.Printf("启动命令失败: %v\n", e)
		return err
	}
	defer cmd.Process.Kill() // 确保退出时终止进程

	// 实时读取输出流
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		rawEvent := scanner.Text()
		//fmt.Printf("原始事件: %s\n", rawEvent)
		fu(rawEvent)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取错误: %v\n", err)
		return err
	}
	return cmd.Wait() // 等待命令退出
}

func RunCMD(name string, args ...string) ([]byte, error) {
	//glog.Println(name, args)
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput() // 合并stdout和stderr
	if err != nil {
		return nil, fmt.Errorf("执行失败: %v, 输出: %s", err, string(output))
	}
	return output, nil
}

func getStatusByMac(mac string) []*Status {
	if mac == "" {
		return nil
	}
	_ = u.CheckDirector(StatusDir)
	tempFilePath := filepath.Join(StatusDir, mac)
	//byteArray, err := os.ReadFile(tempFilePath)
	//if err != nil {
	//	return nil
	//}
	//var cfg []*Status
	//err = ukey.GobToStruct(byteArray, &cfg)
	//if err != nil {
	//	return nil
	//}
	return readStatusByMac(tempFilePath)
}

func readTimeLineByMac(tempFilePath string) []*DeviceTimeLine {
	if tempFilePath == "" {
		return nil
	}
	byteArray, err := os.ReadFile(tempFilePath)
	if err != nil {
		return nil
	}
	var cfg []*DeviceTimeLine
	err = ukey.GobToStruct(byteArray, &cfg)
	if err != nil {
		return nil
	}
	return cfg
}

func readStatusByMac(tempFilePath string) []*Status {
	if tempFilePath == "" {
		return nil
	}
	byteArray, err := os.ReadFile(tempFilePath)
	if err != nil {
		return nil
	}
	var cfg []*Status
	err = ukey.GobToStruct(byteArray, &cfg)
	if err != nil {
		return nil
	}
	return cfg
}

func setStatusByMac(mac string, statusList []*Status) error {
	if mac == "" {
		return nil
	}
	if statusList == nil || len(statusList) == 0 {
		return nil
	}

	content, err := ukey.StructToGob(statusList)
	if err != nil {
		return err
	}

	_ = u.CheckDirector(StatusDir)
	tempFilePath := filepath.Join(StatusDir, mac)
	file, err := os.Create(tempFilePath) // 文件不存在则创建，存在则截断
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入内容
	_, err = file.Write(content)
	return err
}
