package u

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/pkg"
	"net"
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

func IsMillisecondTimestamp(ts int64) bool {
	// 毫秒级：13位（≥1e12），秒级：10位（<1e12）
	return ts >= 1_000_000_000_000
}

func UTC8ToString(timestamp int64, layout string) string {
	loc, err := time.LoadLocation("Asia/Shanghai") // 等价于 UTC+8
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) // 东八区
	}
	if IsMillisecondTimestamp(timestamp) {
		if loc != nil {
			return time.UnixMilli(timestamp).In(loc).Format(layout)
		}
		return time.UnixMilli(timestamp).Format(layout)
	}
	if loc != nil {
		return time.Unix(timestamp, 0).In(loc).Format(layout)
	}
	return time.Unix(timestamp, 0).Format(layout)
}
func UTC8ToTime(timestamp int64) time.Time {
	loc, err := time.LoadLocation("Asia/Shanghai") // 等价于 UTC+8
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) // 东八区
	}
	if IsMillisecondTimestamp(timestamp) {
		if loc != nil {
			return time.UnixMilli(timestamp).In(loc)
		}
		return time.UnixMilli(timestamp)
	}
	if loc != nil {
		return time.Unix(timestamp, 0).In(loc)
	}
	return time.Unix(timestamp, 0)
}

func TimestampToTime(timestamp int64) string {
	loc, err := time.LoadLocation("Asia/Shanghai") // 等价于 UTC+8
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) // 东八区
	}
	if IsMillisecondTimestamp(timestamp) {
		if loc != nil {
			return time.UnixMilli(timestamp).In(loc).Format(fmt.Sprintf("%s.000", time.TimeOnly)) // 0表示纳秒部分
		}
		return time.UnixMilli(timestamp).Format(fmt.Sprintf("%s.000", time.TimeOnly)) // 0表示纳秒部分
	}
	if loc != nil {
		return time.Unix(timestamp, 0).In(loc).Format(time.TimeOnly) // 0表示纳秒部分
	}
	return time.Unix(timestamp, 0).Format(time.TimeOnly) // 0表示纳秒部分
}

func TimestampToDateTime(timestamp int64) string {
	loc, err := time.LoadLocation("Asia/Shanghai") // 等价于 UTC+8
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) // 东八区
	}
	if IsMillisecondTimestamp(timestamp) {
		if loc != nil {
			return time.UnixMilli(timestamp).In(loc).Format(fmt.Sprintf("%s.000", time.DateTime)) // 0表示纳秒部分
		}
		return time.UnixMilli(timestamp).Format(fmt.Sprintf("%s.000", time.DateTime)) // 0表示纳秒部分
	}
	if loc != nil {
		return time.Unix(timestamp, 0).In(loc).Format(time.DateTime) // 0表示纳秒部分
	}
	return time.Unix(timestamp, 0).Format(time.DateTime) // 0表示纳秒部分
}
func TimestampFormatToMonth(timestamp int64) string {
	monthFormat := "2006-01"
	loc, err := time.LoadLocation("Asia/Shanghai") // 等价于 UTC+8
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) // 东八区
	}
	if IsMillisecondTimestamp(timestamp) {
		if loc != nil {
			return time.UnixMilli(timestamp).In(loc).Format(monthFormat) // 0表示纳秒部分
		}
		return time.UnixMilli(timestamp).Format(monthFormat) // 0表示纳秒部分
	}
	if loc != nil {
		return time.Unix(timestamp, 0).In(loc).Format(monthFormat) // 0表示纳秒部分
	}
	return time.Unix(timestamp, 0).Format(monthFormat) // 0表示纳秒部分
}

func GetLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai") // 等价于 UTC+8
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) // 东八区
	}
	return loc
}

func GetTime(timeString string, loc *time.Location) *time.Time {
	t, e := time.ParseInLocation(time.TimeOnly, timeString, loc) // 按北京时间解析
	if e == nil {
		return &t
	}
	return nil
}

func GetDay(timestamp int64) string {
	if !IsMillisecondTimestamp(timestamp) {
		timestamp *= 1000
	}
	return UTC8ToString(timestamp, time.DateOnly)
}

func AutoParse(timeStr string) (*time.Time, error) {
	var layouts = []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}
	for _, layout := range layouts {
		loc := GetLocation()
		t, e := time.ParseInLocation(layout, timeStr, loc) // 按北京时间解析
		if e == nil {
			return &t, nil // 解析成功
		}
	}
	return nil, fmt.Errorf("无法识别的格式")
}

// CompareTime 小于目标时间
func CompareTime(now, target time.Time) int {
	nowSeconds := now.Hour()*3600 + now.Minute()*60 + now.Second()
	tarSeconds := target.Hour()*3600 + target.Minute()*60 + target.Second()
	return nowSeconds - tarSeconds
}

// IsWorkingTime
// 0：上班打卡
// 1：工作时间，不能打卡
// 2：下班打卡
func IsWorkingTime(time1, time2 string) (int, error) {
	t1, e1 := time.Parse(time.TimeOnly, time1)
	if e1 != nil {
		return -1, e1
	}
	t2, e2 := time.Parse(time.TimeOnly, time2)
	if e2 != nil {
		return -1, e2
	}
	now := glog.Now()
	if CompareTime(now, t1) <= 0 {
		//小于等于上班时间
		return 0, nil
	} else if CompareTime(now, t2) < 0 {
		//工作时间内
		return 1, nil
	} else {
		//大于等于下班时间
		return 2, nil
	}
}

func IsOnWorked(time1 string) bool {
	t1, e1 := time.Parse(time.TimeOnly, time1)
	if e1 != nil {
		return false
	}
	now := glog.Now()
	if CompareTime(now, t1) > 0 {
		return true
	} else {
		return false
	}
}

func DateParse(timestr string) (*time.Time, error) {
	t, e1 := time.Parse(time.DateOnly, timestr)
	if e1 == nil {
		return &t, e1
	}
	return nil, nil
}

func TimeParse(timestr string) (*time.Time, error) {
	t, e1 := time.Parse(time.TimeOnly, timestr)
	if e1 == nil {
		return &t, e1
	}
	return nil, nil
}
func TestTimeParse(timestr string) error {
	_, e1 := time.Parse(time.TimeOnly, timestr)
	if e1 != nil {
		return e1
	}
	return nil
}

func GetWeekName(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "星期一"
	case time.Tuesday:
		return "星期二"
	case time.Wednesday:
		return "星期三"
	case time.Thursday:
		return "星期四"
	case time.Friday:
		return "星期五"
	case time.Saturday:
		return "星期六"
	case time.Sunday:
		return "星期日"
	}
	return ""
}
func GetLocalMac() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("获取网络接口失败：", err)
		return ""
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.HardwareAddr != nil {
			devMac := strings.ReplaceAll(iface.HardwareAddr.String(), ":", "")
			fmt.Println(iface.Name, ":", devMac)
			return devMac
		}
	}
	return ""
}

func IsFileExist(file string) bool {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) { // 明确检查“不存在”错误
		fmt.Println("文件不存在")
	} else if err != nil {
		fmt.Println("其他错误:", err)
	} else {
		fmt.Println("文件存在")
		return true
	}
	return false
}
