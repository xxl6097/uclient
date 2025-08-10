package openwrt

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"github.com/xxl6097/uclient/internal/webhook"
	"strings"
	"time"
)

func (this *openWRT) TiggerSignCardEvent(macAddress string) error {
	if v, ok := this.clients[macAddress]; ok {
		if v.Nick != nil && v.Nick.WorkType != nil {
			return this.NotifyDingSign(0, macAddress, glog.Now(), nil, v.Nick.WorkType)
		}
	}
	return nil
}

func (this *openWRT) NotifyDingSign(signal int, macAddress string, now time.Time, wrk *WorkEntry, wts *WorkTypeSetting) error {
	//if wrk == nil {
	//	return fmt.Errorf("设备【%s】没今天打卡数据 ", macAddress)
	//}
	if wts == nil {
		return fmt.Errorf("设备【%s】未设置打卡时间", macAddress)
	}

	//t1 := glog.Now()
	month := fmt.Sprintf("%d-%02d", now.Year(), int(now.Month()))
	day := now.Format(time.DateOnly)
	var monthOverTimes string
	var wk *WorkTime
	works, err := getWorkTimeAndCaculate(macAddress, wts)
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
		Signal:     signal,
	}
	if wrk != nil {
		if wrk.OnWorkTime > 0 && wrk.OffWorkTime == 0 {
			msg.Title = fmt.Sprintf("【%s】上班了", name)
		} else {
			msg.Title = fmt.Sprintf("【%s】下班了", name)
		}
	} else {
		msg.Title = fmt.Sprintf("【%s】考勤统计", name)
	}

	//if working == 0 {
	//	msg.Title = fmt.Sprintf("【%s】上班了", name)
	//	if wrk.OnWorkTime > 0 {
	//		return fmt.Errorf("上班已经打卡了 %v", u.TimestampToTime(wrk.OnWorkTime))
	//	}
	//} else if working == 2 {
	//	msg.Title = fmt.Sprintf("【%s】下班了", name)
	//} else if working == 3 {
	//	msg.Title = fmt.Sprintf("【%s】考勤统计", name)
	//} else {
	//	return fmt.Errorf("当前是异常的工作时间 %v", working)
	//}
	//msg.TodayOverTime = todayOverTimes
	//msg.MonthOverTime = monthOverTimes
	return webhook.Notify(msg, func(builder *strings.Builder) {
		if builder == nil {
			return
		}
		if wrk != nil {
			if wrk.OnWorkTime > 0 {
				builder.WriteString(fmt.Sprintf("- 上班时间：%s\n ", u.TimestampToTime(wrk.OnWorkTime)))
			}
			if wrk.OffWorkTime > 0 {
				builder.WriteString(fmt.Sprintf("- 下班时间：%s\n ", u.TimestampToTime(wrk.OffWorkTime)))
			}
		}
		if wk != nil && wk.OverWorkTimes != "" {
			builder.WriteString(fmt.Sprintf("- 今日加班时长：%s\n ", wk.OverWorkTimes))
		}
		if monthOverTimes != "" {
			builder.WriteString(fmt.Sprintf("- 本月加班时长：%s\n ", monthOverTimes))
		}
	})
}

func (this *openWRT) notifyWebhookMessage(eveName string, client *DHCPLease) error {
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
		EventName:  eveName,
		//WorkTime:   &t,
	}
	if client.Nick != nil && client.Nick.Name != "" {
		msg.DeviceName = client.Nick.Name
	} else {
		msg.DeviceName = client.Hostname
	}
	msg.Signal = client.Signal
	if client.Online {
		msg.Title = fmt.Sprintf("【%s】上线啦", msg.DeviceName)
	} else {
		msg.Title = fmt.Sprintf("【%s】离线了", msg.DeviceName)
	}
	glog.Debugf("ding通知 %+v", client)
	return webhook.Notify(msg, nil)
}
