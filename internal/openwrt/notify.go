package openwrt

import (
	"fmt"
	"github.com/xxl6097/uclient/internal/u"
	"github.com/xxl6097/uclient/internal/webhook"
	"time"
)

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

	name, ip := this.getDeviceName(macAddress)
	msg := webhook.WebHookMessage{
		Url:        webhookUrl,
		DeviceName: name,
		IpAddress:  ip,
		MacAddress: macAddress,
		TimeNow:    &t1,
		DutyTime:   &t1,
	}
	if t1.Weekday() == time.Saturday || t1.Weekday() == time.Sunday {
		msg.Title = fmt.Sprintf("【%s】打卡了", name)
	} else {
		if working == 0 {
			msg.Title = fmt.Sprintf("【%s】上班了", name)
		} else if working == 2 {
			msg.Title = fmt.Sprintf("【%s】下班了", name)
		} else {
			msg.Title = fmt.Sprintf("【%s】考勤统计", name)
			msg.DutyTime = nil
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
	msg.TodayOverTime = todayOverTimes
	msg.MonthOverTime = monthOverTimes
	return webhook.Notify(msg)
}

func (this *openWRT) notifyWebhookMessage(client *DHCPLease) error {
	if this.webhookUrl == "" {
		return fmt.Errorf("webhookUrl is empty")
	}
	if client == nil {
		return fmt.Errorf("client is nil")
	}
	if client.Nick == nil {
		return fmt.Errorf("client.Nick is empty")
	}
	if !client.Nick.IsPush {
		return fmt.Errorf("client.Nick is not push")
	}

	t := u.UTC8ToTime(client.StartTime)
	msg := webhook.WebHookMessage{
		Url:        this.webhookUrl,
		IpAddress:  client.IP,
		MacAddress: client.MAC,
		TimeNow:    &t,
	}
	if client.Nick != nil && client.Nick.Name != "" {
		msg.DeviceName = client.Nick.Name
	} else {
		msg.DeviceName = client.Hostname
	}
	if client.Online {
		msg.Title = fmt.Sprintf("%s上线啦", msg.DeviceName)
	} else {
		msg.Title = fmt.Sprintf("%s离线了", msg.DeviceName)
	}
	return webhook.Notify(msg)
}
