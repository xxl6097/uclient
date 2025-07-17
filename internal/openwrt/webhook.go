package openwrt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"io"
	"net/http"
	"time"
)

func WebHookMessage(webhookUrl string, payload any) error {
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

func (this *openWRT) NotifySignCardEvent(working int, macAddress string, t1 time.Time) error {
	cls, ok := this.clients[macAddress]
	if !ok {
		return fmt.Errorf("设备【%s】不存在内存", macAddress)
	}
	if cls == nil {
		return fmt.Errorf("设备【%s】对象不存在", macAddress)
	}
	if cls.Nick == nil {
		return fmt.Errorf("设备【%s】未设置打卡", macAddress)
	}
	wts := cls.Nick.WorkType
	if wts == nil {
		return fmt.Errorf("设备【%s】未设置打卡时间", macAddress)
	}
	webhookUrl := wts.WebhookUrl
	if webhookUrl == "" {
		return fmt.Errorf("设备【%s】未设置webhook", macAddress)
	}
	var title string
	name, ip := this.getDeviceName(macAddress)
	if t1.Weekday() == time.Saturday || t1.Weekday() == time.Sunday {
		title = fmt.Sprintf("【%s】打卡了", name)
	} else {
		if working == 0 {
			title = fmt.Sprintf("【%s】上班了", name)
		} else if working == 2 {
			title = fmt.Sprintf("【%s】下班了", name)
		} else {
			title = fmt.Sprintf("【%s】加班统计", name)
		}
	}
	month := fmt.Sprintf("%d-%02d", t1.Year(), int(t1.Month()))
	day := t1.Format(time.DateOnly)
	var todayOverTimes, monthOverTimes string
	works, err := getWorkTime(macAddress, wts)
	if err == nil && works != nil {
		for _, work := range works {
			if work.Month == month {
				monthOverTimes = work.OverTime
				if work.WorkTime != nil {
					for _, t := range work.WorkTime {
						if t.Date == day {
							todayOverTimes = t.OverWorkTimes
							break
						}
					}
					break
				}
			}
		}
	}
	markdown := make(map[string]interface{})
	markdown["title"] = title
	markdown["text"] = fmt.Sprintf("#### %s \n - 今天是 %s 星期：%s\n - IP地址：%s \n- Mac地址：%s \n- 打卡时间：%s  \n- 今日加班时长：%s \n- 本月累计加班时长：%s \n",
		title,
		t1.Format(time.DateOnly),
		u.GetWeekName(t1.Weekday()),
		ip,
		macAddress,
		t1.Format(time.DateOnly),
		todayOverTimes,
		monthOverTimes,
	)
	payload := map[string]interface{}{"msgtype": "markdown", "markdown": markdown}
	_ = WebHookMessage(webhookUrl, payload)
	return nil
}
