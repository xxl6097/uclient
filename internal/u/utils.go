package u

import (
	"encoding/json"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/pkg"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Error(code int, msg string) map[string]interface{} {
	return map[string]interface{}{"code": code, "msg": msg, "success": false}
}
func OK(code int, msg string) map[string]interface{} {
	return map[string]interface{}{"code": code, "msg": msg}
}

func OKK(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "msg": "成功"}) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func Sucess(code int, data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{"code": code, "data": data}
}
func SucessWithData(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{"code": 0, "data": data}
}
func SucessWithObject(data interface{}) map[string]interface{} {
	return map[string]interface{}{"code": 0, "data": data}
}
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(data) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	//w.WriteHeader(400)
}

func RespondObject(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(data) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetDataByJson[T any](r *http.Request) (*T, error) {
	var t T
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ClearTemp() error {
	tempDir := glog.TempDir()
	glog.Debug(tempDir)
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}
	for _, entry := range entries {
		fullPath := filepath.Join(tempDir, entry.Name())
		err = os.RemoveAll(fullPath)
		if err != nil {
			glog.Errorf("删除失败 %s  %v", fullPath, err)
		} else {
			glog.Debugf("删除成功 %s", fullPath)
		}
	}
	return err
}

func IsMacOs() bool {
	if strings.Compare(runtime.GOOS, "darwin") == 0 {
		return true
	}
	return false
}

func GetSelfSize() uint64 {
	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("获取可执行文件路径时出错: %v\n", err)
		return 0
	}
	// 获取文件信息
	fileInfo, err := os.Stat(exePath)
	if err != nil {
		fmt.Printf("获取文件信息时出错: %v\n", err)
		return 0
	}

	// 获取文件大小
	fileSize := fileInfo.Size()
	fmt.Printf("本程序自身大小为: %v\n", ByteCountIEC(uint64(fileSize)))
	return uint64(fileSize)
}

func ByteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func GetVersion() map[string]interface{} {
	hostName, _ := os.Hostname()
	return map[string]interface{}{
		"hostName":    hostName,
		"appName":     pkg.AppName,
		"appVersion":  pkg.AppVersion,
		"buildTime":   pkg.BuildTime,
		"gitRevision": pkg.GitRevision,
		"gitBranch":   pkg.GitBranch,
		"goVersion":   pkg.GoVersion,
		"displayName": pkg.DisplayName,
		"description": pkg.Description,
		"osType":      pkg.OsType,
		"arch":        pkg.Arch,
	}
}

func CheckDirector(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// 存在，删除
		return nil
	}
	return os.MkdirAll(path, 0755)
}

func TimestampFormat(timestamp int64) string {
	return time.UnixMilli(timestamp).Format(time.DateTime) // 0表示纳秒部分
}
