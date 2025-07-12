package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	cmd := exec.Command("uci", "show", "dhcp")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("命令执行失败: %v\n", err)
		return
	}
	leases := parseUciShowDHCP(out.String())
	for _, lease := range leases {
		fmt.Printf("索引号: %s\n", lease.Index)
		fmt.Printf("IP地址: %s\n", lease.IP)
		fmt.Printf("MAC地址: %s\n", lease.MAC)
		fmt.Printf("主机名: %s\n", lease.Hostname)
		fmt.Println("-------------------")
	}
}

type DHCPHost struct {
	Index    string
	IP       string
	MAC      string
	Hostname string
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
