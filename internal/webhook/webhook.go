package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"io"
	"net/http"
	"strings"
	"time"
)

type WebHookMessage struct {
	Url           string     `json:"url"`
	Title         string     `json:"title"`
	DeviceName    string     `json:"deviceName"`
	IpAddress     string     `json:"ipAddress"`
	MacAddress    string     `json:"macAddress"`
	TimeNow       *time.Time `json:"timeNow"`
	DutyTime      *time.Time `json:"dutyTime"`
	TodayOverTime string     `json:"todayOverTime"`
	MonthOverTime string     `json:"monthOverTime"`
}

func Notify(msg WebHookMessage) error {
	if msg.Url == "" {
		return fmt.Errorf("webhook url is empty")
	}
	if msg.Title == "" {
		return fmt.Errorf("title is empty")
	}
	text := strings.Builder{}
	text.WriteString(fmt.Sprintf("#### %s \n ", msg.Title))
	if msg.TimeNow != nil {
		text.WriteString(fmt.Sprintf("- 今天是 %s %s\n ", msg.TimeNow.Format(time.DateOnly), u.GetWeekName(msg.TimeNow.Weekday())))
	}
	if msg.DutyTime != nil {
		text.WriteString(fmt.Sprintf("- 打卡时间：%s\n ", msg.DutyTime.Format(time.DateTime)))
	}
	if msg.IpAddress != "" {
		text.WriteString(fmt.Sprintf("- IP地址：%s\n ", msg.IpAddress))
	}
	if msg.MacAddress != "" {
		text.WriteString(fmt.Sprintf("- Mac地址：%s\n ", msg.MacAddress))
	}
	if msg.TodayOverTime != "" {
		text.WriteString(fmt.Sprintf("- 今日加班时长：%s\n ", msg.TodayOverTime))
	}
	if msg.MonthOverTime != "" {
		text.WriteString(fmt.Sprintf("- 本月累计加班时长：%s\n ", msg.MonthOverTime))
	}
	if msg.TimeNow != nil {
		text.WriteString(fmt.Sprintf("- 消息时间：%s\n ", msg.TimeNow.Format(time.DateTime)))
	}
	markdown := make(map[string]interface{})
	markdown["title"] = msg.Title
	markdown["text"] = text.String()
	payload := map[string]interface{}{"msgtype": "markdown", "markdown": markdown}
	return WebHook(msg.Url, payload)
}
func WebHook(webhookUrl string, payload any) error {
	//jsonData, err := json.Marshal(payload)
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	glog.Debug("webhook", string(jsonData))
	resp, err := http.Post(
		webhookUrl,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		glog.Errorf("Error: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("无法读取响应内容:", err)
		return err
	}
	glog.Println("响应内容:", resp.StatusCode, string(respBody))
	return err
}
