package openwrt

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// DHCPHost 表示dhcp配置中的一个主机条目
type DHCPHost struct {
	Index    string `json:"index"`
	Hostname string `json:"hostname"`
	MAC      string `json:"mac"`
	IP       string `json:"ip"`
}

func setStaticLease(mac, ip, name string) error {
	// 添加或更新DHCP配置
	cmdAdd := exec.Command("uci", "add", "dhcp", "host")
	if err := cmdAdd.Run(); err != nil {
		glog.Printf("新增host条目失败（可能已存在）: %v", err)
		return err
	}

	// 设置参数
	exec.Command("uci", "set", "dhcp.@host[-1].name="+name).Run()
	exec.Command("uci", "set", "dhcp.@host[-1].mac="+mac).Run()
	exec.Command("uci", "set", "dhcp.@host[-1].ip="+ip).Run()

	// 提交配置变更
	if err := exec.Command("uci", "commit", "dhcp").Run(); err != nil {
		return err
	}
	return nil
}

func deleteStaticLease(index string) error {
	// 添加或更新DHCP配置
	cmdAdd := exec.Command("uci", "delete", fmt.Sprintf("dhcp.%s", index))
	if err := cmdAdd.Run(); err != nil {
		glog.Printf("新增host条目失败（可能已存在）: %v", err)
		return err
	}

	// 提交配置变更
	if err := exec.Command("uci", "commit", "dhcp").Run(); err != nil {
		return err
	}
	return nil
}

func restartDNSMasq() error {
	cmd := exec.Command("/etc/init.d/dnsmasq", "restart")
	return cmd.Run()
}
func RestartNetwork() error {
	cmd := exec.Command("/etc/init.d/network", "restart")
	return cmd.Run()
}

func isStaticIpUsed(ipAddress, macAddress, name string) error {
	entries, err := GetUCIOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil
	}
	for _, entry := range entries {
		//判断IP地址相同 && Mac地址不同，说明这个IP被另外一个设备(mac)占用
		if entry.IP == ipAddress {
			if entry.MAC != macAddress {
				return fmt.Errorf("ip address %s is used by %s[%s]", ipAddress, entry.Hostname, entry.MAC)
			}
		} else if strings.EqualFold(entry.Hostname, name) {
			//判断名称相同的 && Mac不同，说明这个名称被占用
			if entry.MAC != macAddress {
				return fmt.Errorf("host %s is used by %s[%s]", name, entry.Hostname, entry.MAC)
			}
		}
	}
	return nil
}

func SetStaticIpAddress(mac, ip, name string) error {
	if err := isStaticIpUsed(ip, mac, name); err != nil {
		return err
	}
	err := setStaticLease(mac, ip, name)
	if err != nil {
		return err
	}
	return restartDNSMasq()
}

func DeleteStaticIpAddress(mac string) error {
	entries, err := GetUCIOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil
	}
	for _, entry := range entries {
		if strings.EqualFold(entry.MAC, mac) {
			err = deleteStaticLease(entry.Index)
			if err != nil {
				return err
			}
			return restartDNSMasq()
		}
	}

	return fmt.Errorf("not found mac address %s", mac)
}

func TestStaticLease() {
	mac := "AA:BB:CC:DD:EE:FF" // 替换为实际MAC
	ip := "192.168.1.50"       // 需分配的静态IP
	name := "My-Phone"         // 设备标识

	if err := setStaticLease(mac, ip, name); err != nil {
		log.Fatalf("配置失败: %v", err)
	}
	if err := restartDNSMasq(); err != nil {
		log.Fatalf("服务重启失败: %v", err)
	}
	log.Println("静态IP设置成功！")
}

func parseUciShowDHCP(output string) []DHCPHost {
	var leases []DHCPHost
	lines := strings.Split(output, "\n")
	var currentIndex string
	var currentLease DHCPHost

	reConfig := regexp.MustCompile(`^dhcp\.([^=]+)=host$`)
	reIP := regexp.MustCompile(`^dhcp\.([^.]+)\.ip='([^']+)'$`)
	reMAC := regexp.MustCompile(`^dhcp\.([^.]+)\.mac='([^']+)'$`)
	reHostname := regexp.MustCompile(`^dhcp\.([^.]+)\.name='([^']+)'$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if match := reConfig.FindStringSubmatch(line); match != nil {
			if currentIndex != "" {
				leases = append(leases, currentLease)
			}
			currentIndex = match[1]
			currentLease = DHCPHost{Index: currentIndex}
		} else if currentIndex != "" {
			if match := reIP.FindStringSubmatch(line); match != nil && match[1] == currentIndex {
				currentLease.IP = match[2]
			} else if match := reMAC.FindStringSubmatch(line); match != nil && match[1] == currentIndex {
				currentLease.MAC = match[2]
			} else if match := reHostname.FindStringSubmatch(line); match != nil && match[1] == currentIndex {
				currentLease.Hostname = match[2]
			}
		}
	}

	if currentIndex != "" {
		leases = append(leases, currentLease)
	}

	return leases
}

func GetUCIOutput() ([]DHCPHost, error) {
	cmd := exec.Command("uci", "show", "dhcp")
	output, err := cmd.Output()
	if err != nil {
		glog.Printf("Error: %v\n", err)
		return nil, err
	}
	hosts := parseUciShowDHCP(string(output))
	glog.Println("✅ 解析结果：")
	for _, host := range hosts {
		glog.Printf("索引: %v | MAC: %s | IP: %s | 设备名: %s\n",
			host.Index, host.MAC, host.IP, host.Hostname)
	}
	return hosts, nil
}
