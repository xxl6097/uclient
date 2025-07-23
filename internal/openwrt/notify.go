package openwrt

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/webhook"
	"strings"
	"time"
)

func (this *openWRT) NotifySignCardEvent(working int, macAddress string) error {
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

	t1 := glog.Now()
	month := fmt.Sprintf("%d-%02d", t1.Year(), int(t1.Month()))
	day := t1.Format(time.DateOnly)
	var monthOverTimes string
	var wk *WorkTime
	works, err := getWorkTime(macAddress, wts)
	if err == nil && works != nil {
		for _, work := range works {
			if work.Month == month {
				monthOverTimes = work.OverTime
				if work.WorkTime != nil {
					for _, t := range work.WorkTime {
						if t.Date == day {
							//todayOverTimes = t.OverWorkTimes
							wk = &t
							break
						}
					}
					break
				}
			}
		}
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
	}
	if working == 0 {
		msg.Title = fmt.Sprintf("【%s】上班了", name)
		if wk != nil && wk.WorkTime1 != "" {
			return fmt.Errorf("上班已经打卡了 %v", wk.WorkTime1)
		}
	} else if working == 2 {
		msg.Title = fmt.Sprintf("【%s】下班了", name)
	} else if working == 3 {
		msg.Title = fmt.Sprintf("【%s】考勤统计", name)
	} else {
		return fmt.Errorf("当前是异常的工作时间 %v", working)
	}
	//msg.TodayOverTime = todayOverTimes
	//msg.MonthOverTime = monthOverTimes
	return webhook.Notify(msg, func(builder *strings.Builder) {
		if builder == nil {
			return
		}
		if wk != nil {
			if wk.WorkTime1 != "" {
				builder.WriteString(fmt.Sprintf("- 上班时间：%s\n ", wk.WorkTime1))
			}
			if wk.WorkTime2 != "" {
				builder.WriteString(fmt.Sprintf("- 下班时间：%s\n ", wk.WorkTime2))
			}
			if wk.OverWorkTimes != "" {
				builder.WriteString(fmt.Sprintf("- 今日加班时长：%s\n ", wk.OverWorkTimes))
			}
		}
		if monthOverTimes != "" {
			builder.WriteString(fmt.Sprintf("- 本月加班时长：%s\n ", monthOverTimes))
		}
		//if entry != nil {
		//	if entry.OnWorkTime > 0 {
		//		builder.WriteString(fmt.Sprintf("- 上班时间：%s\n ", u.TimestampToTime(entry.OnWorkTime)))
		//	}
		//	if entry.OffWorkTime > 0 {
		//		builder.WriteString(fmt.Sprintf("- 下班时间：%s\n ", u.TimestampToTime(entry.OffWorkTime)))
		//	}
		//}
	})
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

	//t := u.UTC8ToTime(client.StartTime)
	msg := webhook.WebHookMessage{
		Url:        this.webhookUrl,
		IpAddress:  client.IP,
		MacAddress: client.MAC,
		//WorkTime:   &t,
	}
	if client.Nick != nil && client.Nick.Name != "" {
		msg.DeviceName = client.Nick.Name
	} else {
		msg.DeviceName = client.Hostname
	}
	if client.Online {
		msg.Title = fmt.Sprintf("【%s】上线啦", msg.DeviceName)
	} else {
		msg.Title = fmt.Sprintf("【%s】离线了", msg.DeviceName)
	}
	return webhook.Notify(msg, nil)
}
