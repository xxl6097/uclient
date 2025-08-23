package openwrt

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/uclient/internal/ntfy"
	"github.com/xxl6097/uclient/internal/u"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	instance *openWRT
	once     sync.Once
)

type openWRT struct {
	clients      map[string]*DHCPLease
	nicks        map[string]*NickEntry
	leases       map[string]*DHCPLease
	task         map[string]*u.CountdownTask[*SignData]
	mu           sync.Mutex
	fnEvent      func(int, any)
	webhookUrl   string
	ulistString  string
	ctx          context.Context
	cancel       context.CancelFunc
	statusRuning bool
}

// GetInstance 返回单例实例
func GetInstance() *openWRT {
	once.Do(func() {
		instance = &openWRT{
			clients:      make(map[string]*DHCPLease),
			nicks:        make(map[string]*NickEntry),
			leases:       make(map[string]*DHCPLease),
			task:         make(map[string]*u.CountdownTask[*SignData]),
			statusRuning: false,
		}
		instance.init()
	})
	return instance
}

func (this *openWRT) init() {
	if u.IsMacOs() {
		return
	}
	this.ctx, this.cancel = context.WithCancel(context.Background())
	this.ulistString = UbusList()
	this.initClients()
	//go this.subscribeSysLog()
	go this.subscribeArpEvent()
	go this.subscribeArpPing()
	go this.subscribeFsnotify()
	go this.subscribeStatus()
	if strings.Contains(this.ulistString, "hostapd") {
		go this.subscribeHostapd()
	}
	//if strings.Contains(this.ulistString, "dnsmasq") {
	//	go this.subscribeDnsmasq()
	//}
	if strings.Contains(this.ulistString, "ahsapd.sta") {
		go this.subscribeAhsapdsta()
	}
	this.initNtfy()
}

func (this *openWRT) Close() {
	if this.cancel != nil {
		glog.Debug("close openWRT")
		this.cancel()
	}
	ntfy.GetInstance().Stop()
	_ = glog.Flush()
}

func (this *openWRT) ResetClients() {
	_ = this.initData()
	this.webUpdateAll(this.GetClients())
}
func (this *openWRT) initClients() {
	err := this.initData()
	if err != nil {
		glog.Errorf("initClients Error:%v", err)
		time.Sleep(5 * time.Second)
		glog.Error("5 seconds later and try...")
		this.initClients()
	}
	if this.clients == nil || len(this.clients) == 0 {
		glog.Error("dataMap is empty, 1 seconds later and try...")
		time.Sleep(16 * time.Second)
		this.initClients()
	}
}

func (this *openWRT) initNtfy() {
	if u.IsFileExist(ntfyFilePath) {
		info, err := utils.LoadWithGob[*u.NtfyInfo](ntfyFilePath)
		if err != nil {
			glog.Errorf("initNtfy Error:%v", err)
		} else {
			go ntfy.GetInstance().Start(info)
		}
	}
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
					if v.Online != isOnline {
						glog.Infof("Arp事件:%+v", entry)
						this.updateDeviceStatus("Arp事件", dhcp)
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

// 检测变化并告警
//func (this *openWRT) checkARPDiff(fn func([]string)) {
//	if this.arpList == nil || len(this.arpList) == 0 {
//		return
//	}
//	if fn == nil {
//		return
//	}
//	arpList, err := getArp(brLanString)
//	if err != nil {
//		return
//	}
//	if arpList == nil || len(arpList) == 0 {
//		return
//	}
//	arp1 := strings.Join(arpList, ",")
//	arp2 := strings.Join(this.arpList, ",")
//	if strings.Compare(arp1, arp2) != 0 {
//		this.arpList = arpList
//		fn(arpList)
//	}
//}

func (this *openWRT) listenFsnotify(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			//glog.Debug("listenFsnotify:", event)
			if event.Has(fsnotify.Write) {
				//filePath := event.Name
				//switch event.Name {
				//case dhcpLeasesFilePath:
				//	this.updateClientsByDHCP()
				//	break
				//}
				if strings.EqualFold(event.Name, dhcpLeasesFilePath) {
					go this.updateClientsByDHCP()
				}
				this.webUpdateAll(this.GetClients())
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			glog.Error("error:", err)
		case <-this.ctx.Done():
			glog.Debug("Fsnotify 监听退出...")
			return
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

func (p *openWRT) updateClientsByDHCP() {
	clientArray, err := getClientsByDhcp()
	nickMap, e2 := getNickData()
	if e2 == nil {
		p.nicks = nickMap
	}
	if err != nil {
		glog.Println(fmt.Errorf("getClientsByDhcp Error:%v", err))
	} else {
		glog.Printf("DHCP变化，客户端数量 %+v\n", len(clientArray))
		arpMap, e1 := getClientsByArp(brLanString)
		staInfo := GetStaInfo()
		for _, client := range clientArray {
			mac := client.MAC
			if e1 == nil && arpMap != nil {
				arp := arpMap[mac]
				if arp != nil {
					client.Online = arp.Flags == 2
					client.Phy = arp.Interface
				}
			}
			if e2 == nil && nickMap != nil {
				nick := nickMap[mac]
				if nick != nil {
					client.Nick = nick
				}
			}
			if staInfo != nil {
				sta := staInfo[client.MAC]
				if sta != nil {
					client.Vendor = sta.StaVendor
					if client.Hostname == "" || client.Hostname == "*" {
						client.Hostname = sta.HostName
					}
					num, _ := strconv.Atoi(sta.Rssi)
					client.Signal = num
					client.StaType = sta.StaType
					client.Ssid = sta.Ssid
				}
			}
			if v, ok := p.clients[mac]; ok && v != nil {
				if client.Phy != "" {
					v.Phy = client.Phy
				}
				if client.IP != "" {
					v.IP = client.IP
				}
				if client.Nick != nil {
					v.Nick = client.Nick
				}
				//if client.StartTime > 0 {
				//	v.StartTime = client.StartTime
				//} else {
				//	v.StartTime = glog.Now().UnixMilli()
				//}
				v.Online = client.Online

				v.Vendor = client.Vendor
				if v.Hostname == "" || v.Hostname == "*" {
					v.Hostname = client.Hostname
				}
				v.Signal = client.Signal
				v.StaType = client.StaType
				v.Ssid = client.Ssid
			} else {
				p.clients[mac] = client
			}
			p.leases[mac] = client
			p.updateDeviceStatus("dhcp", client)
		}
	}
}

func (this *openWRT) refreshClients(new *DHCPLease) (*DHCPLease, map[string]*u.StaDevice) {
	if new == nil {
		return nil, nil
	}
	staInfo := GetStaInfo()
	if staInfo != nil {
		sta := staInfo[new.MAC]
		if sta != nil {
			new.Vendor = sta.StaVendor
			if new.Hostname == "" || new.Hostname == "*" {
				new.Hostname = sta.HostName
			}
			num, _ := strconv.Atoi(sta.Rssi)
			new.Signal = num
			new.StaType = sta.StaType
			new.Ssid = sta.Ssid
			if sta.Timestamp != 0 && new.StartTime == 0 {
				new.StartTime = sta.Timestamp
			}
		} else {
			new.Signal = 0
		}
	}
	old := this.getClient(new.MAC)
	if old != nil {
		if new.IP != "" {
			old.IP = new.IP
		}
		if new.Ssid != "" {
			old.Ssid = new.Ssid
		}
		if new.Phy != "" {
			old.Phy = new.Phy
		}
		if new.Hostname != "" {
			old.Hostname = new.Hostname
		}
		if new.StartTime > 0 {
			old.StartTime = new.StartTime
		}
		//if new.Signal != 0 {
		//	old.Signal = new.Signal
		//}
		if new.Vendor != "" {
			old.Vendor = new.Vendor
		}
		if new.Freq != 0 {
			old.Freq = new.Freq
		}
		if new.StaType != "" {
			old.StaType = new.StaType
		}
		if new.UpRate != "" {
			old.UpRate = new.UpRate
		}
		if new.DownRate != "" {
			old.DownRate = new.DownRate
		}
		if new.Device != nil {
			old.Device = new.Device
		}
		if new.Nick != nil {
			old.Nick = new.Nick
		}
		if new.Static != nil {
			old.Static = new.Static
		}
		old.Flags = new.Flags
		old.Signal = new.Signal
	} else {
		glog.Debugf("新设备：%+v", new)
		if this.nicks == nil {
			nickMap, e2 := getNickData()
			if e2 == nil {
				this.nicks = nickMap
			}
		}
		if this.nicks != nil {
			new.Nick = this.nicks[new.MAC]
		}
		this.clients[new.MAC] = new
	}
	return old, staInfo
}

func (this *openWRT) updateUserTimeLineData(macAddr string, newList []*Status) {
	tempList := getStatusByMac(macAddr)
	if tempList == nil {
		tempList = newList
	} else {
		element := tempList[len(tempList)-1]
		if element != nil {
			for i, n := range newList {
				if n.Timestamp >= element.Timestamp {
					if n.Timestamp == element.Timestamp && n.Connected == element.Connected {
						continue
					}
					tempList = append(tempList, newList[i:]...)
					break
				}
			}
		}
	}
	if tempList == nil {
		return
	}
	size := len(tempList)
	if len(tempList) > MAX_SIZE {
		tempSize := size - MAX_SIZE
		tempList = tempList[tempSize:]
	}
	_ = setStatusByMac(macAddr, tempList)
}
