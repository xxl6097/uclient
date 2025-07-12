package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

// 静态IP配置结构体
type DHCPHost struct {
	Index int    `json:"index"` // 配置节点索引号
	MAC   string `json:"mac"`   // MAC地址
	IP    string `json:"ip"`    // IP地址
	Name  string `json:"name"`  // 设备名称（可选）
}

func main() {
	// 模拟从`uci show dhcp`读取输入（实际使用时可从命令管道读取）
	input := `dhcp.cfg01411c=dnsmasq
dhcp.cfg01411c.domainneeded='1'
dhcp.cfg01411c.localise_queries='1'
dhcp.cfg01411c.rebind_protection='1'
dhcp.cfg01411c.rebind_localhost='1'
dhcp.cfg01411c.local='/lan/'
dhcp.cfg01411c.domain='lan'
dhcp.cfg01411c.expandhosts='1'
dhcp.cfg01411c.cachesize='8000'
dhcp.cfg01411c.authoritative='1'
dhcp.cfg01411c.readethers='1'
dhcp.cfg01411c.leasefile='/tmp/dhcp.leases'
dhcp.cfg01411c.resolvfile='/tmp/resolv.conf.d/resolv.conf.auto'
dhcp.cfg01411c.localservice='1'
dhcp.cfg01411c.mini_ttl='3600'
dhcp.cfg01411c.dns_redirect='1'
dhcp.cfg01411c.ednspacket_max='1232'
root@Clife:/tmp# 
root@Clife:/tmp# 
root@Clife:/tmp# uci show dhcp.@dnsmasq
uci: Invalid argument
root@Clife:/tmp# uci show dhcp
dhcp.@dnsmasq[0]=dnsmasq
dhcp.@dnsmasq[0].domainneeded='1'
dhcp.@dnsmasq[0].localise_queries='1'
dhcp.@dnsmasq[0].rebind_protection='1'
dhcp.@dnsmasq[0].rebind_localhost='1'
dhcp.@dnsmasq[0].local='/lan/'
dhcp.@dnsmasq[0].domain='lan'
dhcp.@dnsmasq[0].expandhosts='1'
dhcp.@dnsmasq[0].cachesize='8000'
dhcp.@dnsmasq[0].authoritative='1'
dhcp.@dnsmasq[0].readethers='1'
dhcp.@dnsmasq[0].leasefile='/tmp/dhcp.leases'
dhcp.@dnsmasq[0].resolvfile='/tmp/resolv.conf.d/resolv.conf.auto'
dhcp.@dnsmasq[0].localservice='1'
dhcp.@dnsmasq[0].mini_ttl='3600'
dhcp.@dnsmasq[0].dns_redirect='1'
dhcp.@dnsmasq[0].ednspacket_max='1232'
dhcp.lan=dhcp
dhcp.lan.interface='lan'
dhcp.lan.start='100'
dhcp.lan.limit='150'
dhcp.lan.leasetime='12h'
dhcp.lan.dhcpv4='server'
dhcp.lan.dhcpv6='hybrid'
dhcp.lan.ra='hybrid'
dhcp.lan.ra_flags='managed-config' 'other-config'
dhcp.lan.ndp='hybrid'
dhcp.lan.ra_management='1'
dhcp.wan=dhcp
dhcp.wan.interface='wan'
dhcp.wan.ignore='1'
dhcp.odhcpd=odhcpd
dhcp.odhcpd.maindhcp='0'
dhcp.odhcpd.leasefile='/tmp/hosts/odhcpd'
dhcp.odhcpd.leasetrigger='/usr/sbin/odhcpd-update'
dhcp.odhcpd.loglevel='4'
dhcp.@host[0]=host
dhcp.@host[0].name='M4'
dhcp.@host[0].mac='EA:E6:51:97:81:F6'
dhcp.@host[0].ip='192.168.0.199'
dhcp.@host[1]=host
dhcp.@host[1].name='xiaomi15'
dhcp.@host[1].mac='EA:E6:51:97:81:F6'
dhcp.@host[1].ip='192.168.0.6'
dhcp.@host[2]=host
dhcp.@host[2].name='fnos'
dhcp.@host[2].mac='8C:EC:4B:58:81:09'
dhcp.@host[2].ip='192.168.0.3'`

	hosts := parseUciShowDHCP(input)
	fmt.Println("✅ 解析结果：")
	for _, host := range hosts {
		fmt.Printf("索引: %d | MAC: %s | IP: %s | 设备名: %s\n",
			host.Index, host.MAC, host.IP, host.Name)
	}
}

func parseUciShowDHCP(input string) []DHCPHost {
	var hosts []DHCPHost
	scanner := bufio.NewScanner(strings.NewReader(input))

	// 正则匹配索引号（如 @host[0]）
	indexRegex := regexp.MustCompile(`dhcp\.@host\[(\d+)\]`)
	// 匹配键值对（兼容带引号和不带引号）
	kvRegex := regexp.MustCompile(`dhcp\.@host\[\d+\]\.(\w+)=['"]?([^'"]+)['"]?`)

	currentIndex := -1
	currentHost := DHCPHost{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 1. 匹配索引行（如 dhcp.@host[0]=host）
		if matches := indexRegex.FindStringSubmatch(line); len(matches) > 1 {
			// 保存上一个Host配置
			if currentIndex >= 0 {
				hosts = append(hosts, currentHost)
			}
			// 初始化新Host
			fmt.Sscanf(matches[1], "%d", &currentIndex)
			currentHost = DHCPHost{Index: currentIndex}
		}

		// 2. 匹配键值对（MAC/IP/Name）
		if matches := kvRegex.FindStringSubmatch(line); len(matches) > 2 {
			key, value := matches[1], matches[2]
			switch key {
			case "mac":
				currentHost.MAC = value
			case "ip":
				currentHost.IP = value
			case "name":
				currentHost.Name = value
			}
		}
	}

	// 添加最后一个Host
	if currentIndex >= 0 {
		hosts = append(hosts, currentHost)
	}
	return hosts
}
