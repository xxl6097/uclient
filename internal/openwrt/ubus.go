package openwrt

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"os/exec"
	"syscall"
)

var (
	//强制断开客户端
	//addr：客户端MAC地址
	//reason：断开原因代码（如5表示协议错误）
	//deauth：是否发送解认证帧
	offline = "ubus call hostapd.* del_client '{\"addr\":\"5a:a7:22:62:3d:26\", \"reason\":5, \"deauth\":true}'"
	//list    = "ubus call hostapd.* get_clients"
)

// OfflineDevice ubus -v list hostapd.*
// ubus call hostapd.* del_client '{"addr":"5a:a7:22:62:3d:26", "reason":5, "deauth":true}'
func OfflineDevice(macAddr string) error {
	// 构造 ubus 命令
	cmd := exec.Command(
		"ubus",
		"call",
		"hostapd.*",
		"del_client",
		fmt.Sprintf(`{"addr":"%s", "reason":5, "deauth":true}`, macAddr),
	)

	// 设置系统级错误处理（权限提升）
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // 避免子进程被父进程信号中断
	}

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		glog.Printf("❌ 强制下线失败: %v\n输出: %s\n", err, string(output))
		return err
	}

	glog.Printf("✅ 设备 %s 已强制下线\n", macAddr)
	return nil

}

func UbusList() string {
	// 构造 ubus 命令
	cmd := exec.Command(
		"ubus",
		"list",
	)

	// 设置系统级错误处理（权限提升）
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // 避免子进程被父进程信号中断
	}

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		glog.Printf("❌ %v\n输出: %s\n", err, string(output))
		return ""
	}
	return string(output)
}
