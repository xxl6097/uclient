package openwrt

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/glog/pkg/zutil"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/uclient/internal/u"
)

type subData struct {
	IsSubscribed bool
	ctx          context.Context
	cancel       context.CancelFunc
}

func (this *openWRT) SetSettings(settings *u.Settings) error {
	z.Debug("setting file path:", settingPath)
	z.Debug("settings:", settings)
	return utils.SaveToFile[*u.Settings](settings, settingPath)
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
		z.Debugf("读取默认系统配置：%v", info)
		return info
	}

	z.Errorf("读取默认系统配置 Error:%v  info:%v", err, info)
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
	if settings == nil {
		z.Warnf("系统设置孔：%+v", settings)
		return
	}
	z.Warnf("系统设置：%+v", settings)
	if settings.IsSysLogListen {
		z.Warn("启动 SysLog")
		go this.subscribeSysLog()
	}
	if settings.IsArpListen {
		z.Warn("启动 ArpEvent")
		go this.subscribeArpEvent()
	}
	go this.subscribeArpPing()
	go this.subscribeFsnotify()
	//读取状态，显示网速等信息
	go this.subscribeStatus()
	if settings.IsHostApdListen {
		z.Warn("启动 hostapd")
		if strings.Contains(this.ulistString, "hostapd") {
			go this.subscribeHostapd()
		}
	}
	if settings.IsDnsmasqListen {
		z.Warn("启动 dnsmasq")
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
			z.Debug("logread 监听退出...")
			return
		default:
			err := subscribeSysLogs(this.ctx, func(s string) {
				subscribeHostapdLog(s, func(event *SysLogEvent) {
					if event != nil && event.Mac != "" {
						z.Infof("Hostapd事件:%+v", event)
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
					z.Debug("HetSysLog", s)
					if event != nil && event.MACAddress != "" {
						z.Infof("HetSysLog事件:%+v", event)
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
				z.Error(fmt.Errorf("logread 监听失败 %v", err))
				time.Sleep(time.Second * 10)
				z.Error("重新监听 logread")
				tryCount++
				if tryCount > RE_REY_MAX_COUNT {
					z.Error("监听 logread 失败，超过最大重试次数")
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
			z.Debug("Hostapd 监听退出...")
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
						z.Infof("Hostapd事件:%+v", device)
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
				z.Errorf("订阅失败 Hostapd %v", err)
				time.Sleep(time.Second * 5)
				z.Error("重新订阅 Hostapd")
				tryCount++
				if tryCount > RE_REY_MAX_COUNT {
					z.Error("订阅 Hostapd 失败，超过最大重试次数")
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
					StartTime: zutil.Now().UnixMilli(),
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
						z.Infof("Arp事件:%+v", entry)
						go this.updateDeviceStatus("Arp事件", dhcp)
					}
				} else {
					z.Infof("Arp事件(新增):%+v", entry)
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
			z.Debug("Dnsmasq 监听退出...")
			return
		default:
			err := SubscribeDnsmasq(this.ctx, func(device *DnsmasqDevice) {
				if device != nil && device.Mac != "" {
					z.Infof("Dnsmasq事件:%+v", device)
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
				z.Errorf("Dnsmasq订阅失败 %v", err)
				time.Sleep(time.Second * 5)
				z.Error("重新订阅 Dnsmasq")
				tryCount++
				if tryCount > RE_REY_MAX_COUNT {
					z.Error("监听 Dnsmasq 失败，超过最大重试次数")
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
			z.Debug("Ahsapdsta 监听退出...")
			return
		default:
			err := SubscribeSta(this.ctx, func(device *StaUpDown) {
				if device != nil && device.MacAddress != "" {
					z.Infof("ahsapd.sta事件:%+v", device)
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
						dhcp.StartTime = zutil.Now().UnixMilli()
					}
					z.Infof("ahsapd.sta dhcp:%+v", dhcp)
					go this.updateDeviceStatus("ahsapd事件", dhcp)
				}
			})
			if err != nil {
				z.Errorf("ahsapd 订阅失败 %v", err)
				time.Sleep(time.Second * 5)
				z.Error("重新订阅 ahsapd")
				tryCount++
				if tryCount > RE_REY_MAX_COUNT {
					z.Error("监听 ahsapd 失败，超过最大重试次数")
					break
				}
			}
		}
	}
}

func (this *openWRT) subscribeFsnotify() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		z.Error(fmt.Errorf("创建监控器失败 %v", err))
	}

	this.CheckFile(dhcpLeasesFilePath)
	//go this.listenFsnotify(watcher)
	err = watcher.Add(dhcpLeasesFilePath)
	//err = watcher.Add(hetsysinfoFilePath)
	if err != nil {
		z.Error(fmt.Errorf("watcher add err %v", err))
	}
	//go this.SubscribeSysJsonFile(watcher)
	this.listenFsnotify(watcher)
}
