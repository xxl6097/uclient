package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/ntfy"
	"github.com/xxl6097/uclient/internal/u"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type WebHookMessage struct {
	Url        string `json:"url"`
	Title      string `json:"title"`
	DeviceName string `json:"deviceName"`
	EventName  string `json:"eventName"`
	IpAddress  string `json:"ipAddress"`
	MacAddress string `json:"macAddress"`
	Signal     int    `json:"signal"`
	Vendor     string `json:"vendor"`
	//WorkTime      *time.Time `json:"dutyTime"`
	//TodayOverTime string `json:"todayOverTime"`
	//MonthOverTime string `json:"monthOverTime"`
}

func Notify(msg WebHookMessage, fn func(*strings.Builder)) error {
	if msg.Url == "" {
		return fmt.Errorf("webhook url is empty")
	}
	if msg.Title == "" {
		return fmt.Errorf("title is empty")
	}
	text := strings.Builder{}
	text.WriteString(fmt.Sprintf("#### %s \n ", msg.Title))
	now := glog.Now()
	text.WriteString(fmt.Sprintf("- 今天是 %s %s\n ", now.Format(time.DateOnly), u.GetWeekName(now.Weekday())))
	hostName, err := os.Hostname()
	if err == nil && hostName != "" {
		text.WriteString(fmt.Sprintf("- 设备：%s\n ", hostName))
	}
	if msg.EventName != "" {
		text.WriteString(fmt.Sprintf("- 类型：%s\n ", msg.EventName))
	}
	if msg.Signal != 0 {
		text.WriteString(fmt.Sprintf("- 信号：%d\n ", msg.Signal))
	}
	if msg.Vendor != "" {
		text.WriteString(fmt.Sprintf("- 品牌：%v\n ", msg.Vendor))
	}
	if msg.IpAddress != "" {
		text.WriteString(fmt.Sprintf("- IP地址：%s\n ", msg.IpAddress))
	}
	if msg.MacAddress != "" {
		text.WriteString(fmt.Sprintf("- Mac地址：%s\n ", msg.MacAddress))
	}
	//if msg.TodayOverTime != "" {
	//	text.WriteString(fmt.Sprintf("- 今日加班时长：%s\n ", msg.TodayOverTime))
	//}
	//if msg.MonthOverTime != "" {
	//	text.WriteString(fmt.Sprintf("- 本月累计加班时长：%s\n ", msg.MonthOverTime))
	//}
	if fn != nil {
		fn(&text)
	}
	//if msg.WorkTime != nil {
	//	text.WriteString(fmt.Sprintf("- 打卡时间：%s\n ", msg.WorkTime.Format(time.DateTime)))
	//}
	text.WriteString(fmt.Sprintf("- 消息时间：%s\n ", u.TimestampToMilliTime(glog.Now().UnixMilli())))
	markdown := make(map[string]interface{})
	markdown["title"] = msg.Title
	markdown["text"] = text.String()
	payload := map[string]interface{}{"msgtype": "markdown", "markdown": markdown}
	glog.Debug("webhook", msg.Title)

	e := ntfy.GetInstance().Publish(&u.NtfyEventData{
		Topic:    "work",
		Title:    msg.Title,
		Message:  text.String(),
		Markdown: true,
	})
	if e != nil {
		glog.Error("ntfy", e)
	}

	return WebHook(msg.Url, payload)
}
func WebHook(webhookUrl string, payload any) error {
	//jsonData, err := json.Marshal(payload)
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	//glog.Debug("webhook", string(jsonData))
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
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%s %s", resp.Status, string(respBody)))
	}
	//glog.Println("响应内容:", resp.StatusCode, string(respBody))
	return err
}
