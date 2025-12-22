package openwrt

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/uclient/internal/u"
)

type subData struct {
	IsSubscribed bool
	ctx          context.Context
	cancel       context.CancelFunc
}

func (this *openWRT) SetSettings(settings *u.Settings) error {
	setting := &u.Settings{
		IsSysLogListen:  true,
		IsArpListen:     true,
		IsDnsmasqListen: true,
		IsHostApdListen: true,
	}
	return utils.SaveToFile[*u.Settings](setting, settingPath)
}

func (this *openWRT) GetSettings() (*u.Settings, error) {
	if u.IsFileExist(settingPath) {
		return utils.LoadFromFile[*u.Settings](settingPath)
	}
	return nil, nil
}

func (this *openWRT) loadSettings() *u.Settings {
	info, err := this.GetSettings()
	if err == nil && info != nil {
		glog.Errorf("settingPath Error:%v", err)
		return info
	}
	setting := &u.Settings{
		IsSysLogListen:  true,
		IsArpListen:     true,
		IsDnsmasqListen: true,
		IsHostApdListen: true,
	}
	_ = this.SetSettings(setting)
	return setting
}

func (this *openWRT) subscribe() {
	settings := this.loadSettings()
	if settings.IsSysLogListen {
		glog.Warn("启动 SysLog")
		go this.subscribeSysLog()
	}
	if settings.IsArpListen {
		glog.Warn("启动 ArpEvent")
		go this.subscribeArpEvent()
	}
	go this.subscribeArpPing()
	go this.subscribeFsnotify()
	//读取状态，显示网速等信息
	go this.subscribeStatus()
	if settings.IsHostApdListen {
		glog.Warn("启动 hostapd")
		if strings.Contains(this.ulistString, "hostapd") {
			go this.subscribeHostapd()
		}
	}
	if settings.IsDnsmasqListen {
		glog.Warn("启动 dnsmasq")
		if strings.Contains(this.ulistString, "dnsmasq") {
			go this.subscribeDnsmasq()
		}
	}

	//if strings.Contains(this.ulistString, "ahsapd.sta") {
	//	go this.subscribeAhsapdsta()
	//}
}

func (this *openWRT) subscribeSysLog() {
	tryCount := 0
	for {
		select {
		case <-this.ctx.Done():
			glog.Debug("logread 监听退出...")
			return
		default:
			err := subscribeSysLogs(this.ctx, func(s string) {
				subscribeHostapdLog(s, func(event *SysLogEvent) {
					if event != nil && event.Mac != "" {
						glog.Infof("Hostapd事件:%+v", event)
						eve := &DHCPLease{
							MAC:       event.Mac,
							Online:    event.Online,
							StartTime: event.Timestamp.UnixMilli(),
							Phy:       event.Phy,
						}
						if v, ok := this.leases[eve.MAC]; ok {
							if v.Hostname != "" {
								eve.Hostname = v.Hostname
							}
						}
						go this.updateDeviceStatus("Hostapd事件", eve)
					}
				})
				subscribeHetSysLog(s, func(event *KernelLog) {
					glog.LogToFile("HetSysLog", s)
					if event != nil && event.MACAddress != "" {
						glog.Infof("HetSysLog事件:%+v", event)
						eve := &DHCPLease{
							MAC:       event.MACAddress,
							Online:    event.Online,
							StartTime: event.Timestamp,
						}
						if v, ok := this.leases[eve.MAC]; ok {
							if v.Hostname != "" {
								eve.Hostname = v.Hostname
							}
						}
						go this.updateDeviceStatus("HetSysLog事件", eve)
					}
				})
				subscribeLedLog(s)
			})
			if err != nil {
				glog.Error(fmt.Errorf("logread 监听失败 %v", err))
				time.Sleep(time.Second * 10)
				glog.Error("重新监听 logread")
				tryCount++
				if tryCount > RE_REY_MAX_COUNT {
					glog.Error("监听 logread 失败，超过最大重试次数")
					break
				}
			}
		}
	}
}

func (this *openWRT) subscribeHostapd() {
	tryCount := 0
	for {
		select {
		case <-this.ctx.Done():
			glog.Debug("Hostapd 监听退出...")
			return
		default:
			err := SubscribeHostapd(this.ctx, func(device *HostapdDevice) {
				if device != nil && device.Address != "" {
					if device.DataType == 2 {
						cls := this.getClient(device.Address)
						if cls != nil {
							cls.Signal = device.Signal
							cls.Freq = device.Freq
							cls.StartTime = device.Timestamp.UnixMilli()
							//这里只是更新信号，不在web上notify通知
							this.webUpdateOne(cls)
						}
					} else {
						glog.Infof("Hostapd事件:%+v", device)
						dhcp := &DHCPLease{
							MAC:       device.Address,
							Signal:    device.Signal,
							Freq:      device.Freq,
							StartTime: device.Timestamp.UnixMilli(),
							Online:    device.DataType == 0,
						}
						if v, ok := this.leases[dhcp.MAC]; ok {
							if v.Hostname != "" {
								dhcp.Hostname = v.Hostname
							}
						}
						go this.updateDeviceStatus("Hostapd事件", dhcp)
					}
				}
			})
			if err != nil {
				glog.Errorf("订阅失败 Hostapd %v", err)
				time.Sleep(time.Second * 5)
				glog.Error("重新订阅 Hostapd")
				tryCount++
				if tryCount > RE_REY_MAX_COUNT {
					glog.Error("订阅 Hostapd 失败，超过最大重试次数")
					break
				}
			}
		}
	}
}

func (this *openWRT) subscribeArpEvent() {
	SubscribeArpCache(this.ctx, time.Second*10, func(entrys map[string]*ARPEntry) {
		if entrys != nil {
			for mac, entry := range entrys {
				dhcp := &DHCPLease{
					MAC:       entry.MAC.String(),
					IP:        entry.IP.String(),
					Phy:       entry.Interface,
					StartTime: glog.Now().UnixMilli(),
					Flags:     entry.Flags,
				}
				if v, ok := this.leases[mac]; ok {
					if v.Hostname != "" {
						dhcp.Hostname = v.Hostname
					}
				}
				if v, ok := this.clients[mac]; ok {
					isOnline := false
					if entry.Flags == 2 {
						isOnline = true
					}
					if v.Online != isOnline && !isOnline {
						glog.Infof("Arp事件:%+v", entry)
						go this.updateDeviceStatus("Arp事件", dhcp)
					}
				} else {
					glog.Infof("Arp事件(新增):%+v", entry)
					this.clients[mac] = dhcp
				}
			}
		}
	})
}

func (this *openWRT) subscribeArpPing() {
	SubscribeArpCache(this.ctx, time.Second*10, func(entrys map[string]*ARPEntry) {
		if entrys != nil {
			for mac, entry := range entrys {
				if v, ok := this.clients[mac]; ok {
					if v.Online {
						if entry.Flags == 2 {
							u.Ping(entry.IP.String())
						}
					}
				}

			}
		}
	})
}

func (this *openWRT) subscribeDnsmasq() {
	tryCount := 0
	for {
		select {
		case <-this.ctx.Done():
			glog.Debug("Dnsmasq 监听退出...")
			return
		default:
			err := SubscribeDnsmasq(this.ctx, func(device *DnsmasqDevice) {
				if device != nil && device.Mac != "" {
					glog.Infof("Dnsmasq事件:%+v", device)
					dhcp := &DHCPLease{
						MAC:       device.Mac,
						IP:        device.Ip,
						Hostname:  device.Name,
						StartTime: device.Timestamp.UnixMilli(),
						Phy:       device.Interface,
						Online:    true,
					}
					go this.updateDeviceStatus("Dnsmasq事件", dhcp)
				}
			})
			if err != nil {
				glog.Errorf("Dnsmasq订阅失败 %v", err)
				time.Sleep(time.Second * 5)
				glog.Error("重新订阅 Dnsmasq")
				tryCount++
				if tryCount > RE_REY_MAX_COUNT {
					glog.Error("监听 Dnsmasq 失败，超过最大重试次数")
					break
				}
			}
		}
	}

}

func (this *openWRT) subscribeAhsapdsta() {
	tryCount := 0
	for {
		select {
		case <-this.ctx.Done():
			glog.Debug("Ahsapdsta 监听退出...")
			return
		default:
			err := SubscribeSta(this.ctx, func(device *StaUpDown) {
				if device != nil && device.MacAddress != "" {
					glog.Infof("ahsapd.sta事件:%+v", device)
					num, _ := strconv.Atoi(device.Rssi)
					dhcp := &DHCPLease{
						MAC:       u.MacFormat(device.MacAddress),
						Hostname:  device.HostName,
						StartTime: device.Timestamp.UnixMilli(),
						Phy:       device.Ssid,
						Online:    device.Online == 1,
						Signal:    num,
					}
					if dhcp.StartTime <= 0 {
						dhcp.StartTime = glog.Now().UnixMilli()
					}
					glog.Infof("ahsapd.sta dhcp:%+v", dhcp)
					go this.updateDeviceStatus("ahsapd事件", dhcp)
				}
			})
			if err != nil {
				glog.Errorf("ahsapd 订阅失败 %v", err)
				time.Sleep(time.Second * 5)
				glog.Error("重新订阅 ahsapd")
				tryCount++
				if tryCount > RE_REY_MAX_COUNT {
					glog.Error("监听 ahsapd 失败，超过最大重试次数")
					break
				}
			}
		}
	}
}

func (this *openWRT) subscribeFsnotify() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Error(fmt.Errorf("创建监控器失败 %v", err))
	}

	this.CheckFile(dhcpLeasesFilePath)
	//go this.listenFsnotify(watcher)
	err = watcher.Add(dhcpLeasesFilePath)
	//err = watcher.Add(hetsysinfoFilePath)
	if err != nil {
		glog.Error(fmt.Errorf("watcher add err %v", err))
	}
	//go this.SubscribeSysJsonFile(watcher)
	this.listenFsnotify(watcher)
}
