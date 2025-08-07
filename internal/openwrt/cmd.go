package openwrt

import (
	"bufio"
	"context"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/utils/util"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

var procs = make([]*os.Process, 0)

func init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		for _, proc := range procs {
			if proc != nil && proc.Pid > 0 {
				_ = proc.Kill()
			}
		}
		os.Exit(1)
	}()
}

func writePid(name string, pid int) {
	pidData := []byte(strconv.Itoa(pid))
	if err := os.WriteFile(filepath.Join(glog.AppHome("pid"), name), pidData, 0644); err != nil {
		glog.Error("写入PID文件失败:", err)
	}
}

func readPid(name string) int {
	data, err := os.ReadFile(filepath.Join(glog.AppHome("pid"), name))
	if err != nil {
		glog.Error("读取PID失败:", err)
		return 0
	}
	savedPid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0
	}
	return savedPid
}

func killPid(pid int) {
	if pid == 0 {
		return
	}
	// 4. 根据PID终止进程
	if runtime.GOOS == "windows" {
		killCmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
		if err := killCmd.Run(); err != nil {
			log.Fatal("终止失败:", err)
		}
	} else {
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
			log.Fatal("终止失败:", err)
		}
	}
}

func Command(ctx context.Context, fu func(string), name string, arg ...string) error {
	glog.Println(name, arg)
	//ctx, cancel := context.WithCancel(context.Background())
	// 创建ubus命令对象
	ccc := exec.CommandContext(ctx, name, arg...)
	//ccc.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // 创建进程组
	util.SetPlatformSpecificAttrs(ccc)
	// 创建标准输出管道
	stdout, err := ccc.StdoutPipe()
	if err != nil {
		fmt.Printf("创建管道失败: %v\n", err)
		return err
	}

	// 启动命令
	if e := ccc.Start(); e != nil {
		fmt.Printf("启动命令失败: %v\n", e)
		return err
	}
	procs = append(procs, ccc.Process)
	procName := fmt.Sprintf("%s%s", name, strings.Join(arg, "."))
	pid := readPid(procName)
	killPid(pid)
	writePid(procName, ccc.Process.Pid)
	defer ccc.Process.Kill() // 确保退出时终止进程

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
	return ccc.Wait() // 等待命令退出
}

func RunCMD(name string, args ...string) ([]byte, error) {
	//glog.Println(name, args)
	cmd := exec.Command(name, args...)
	//cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // 创建进程组
	util.SetPlatformSpecificAttrs(cmd)
	output, err := cmd.CombinedOutput() // 合并stdout和stderr
	if err != nil {
		return nil, fmt.Errorf("执行失败: %v, 输出: %s", err, string(output))
	}
	return output, nil
}
