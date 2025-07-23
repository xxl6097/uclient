package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const sockPath = "/var/run/hostapd/phy1-ap0" // 根据实际接口名修改

func connectHostapd() (net.Conn, error) {
	//addr := net.UnixAddr{Name: sockPath, Net: "unix"}
	//conn, err := net.DialUnix("unix", nil, &addr)
	//if err != nil {
	//	fmt.Printf("Error connecting to hostapd: %v\n", err)
	//	return nil, err
	//}
	// 正确连接 hostapd（流式套接字）
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()
	return conn, nil
}
func getConnectedClients(conn net.Conn) ([]string, error) {
	// 发送命令并读取响应
	if _, err := conn.Write([]byte("STA-FIRST\n")); err != nil {
		fmt.Printf("getConnectedClients: %v\n", err)
		return nil, err
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	// 解析响应（示例格式：aa:bb:cc:dd:ee:ff\n）
	clients := parseMacList(string(buf[:n]))
	return clients, nil
}

// 解析MAC地址列表
func parseMacList(data string) []string {
	lines := strings.Split(data, "\n")
	var macs []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			macs = append(macs, strings.TrimSpace(line))
		}
	}
	return macs
}
func monitorEvents(conn net.Conn) {
	// 进入事件监听模式
	if _, err := conn.Write([]byte("ATTACH\n")); err != nil {
		log.Fatal(err)
	}

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		msg := string(buf[:n])
		switch {
		case strings.Contains(msg, "AP-STA-CONNECTED"):
			mac := strings.TrimPrefix(msg, "AP-STA-CONNECTED ")
			log.Printf("Device connected: %s", mac)
		case strings.Contains(msg, "AP-STA-DISCONNECTED"):
			mac := strings.TrimPrefix(msg, "AP-STA-DISCONNECTED ")
			log.Printf("Device disconnected: %s", mac)
		}
	}
}
func main() {
	// 1. 建立命令连接
	connCmd, _ := connectHostapd()
	defer connCmd.Close()

	// 2. 获取当前在线设备
	clients, _ := getConnectedClients(connCmd)
	log.Println("Online devices:", clients)

	// 3. 建立事件连接并监听
	connEvent, _ := connectHostapd()
	defer connEvent.Close()
	go monitorEvents(connEvent)

	// 保持主线程运行
	select {}
}
