package openwrt

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"strings"
	"sync"
	"time"
)

var (
	instance *openWRT
	once     sync.Once
)

type openWRT struct {
	clients    map[string]*DHCPLease
	nicks      map[string]*NickEntry
	mu         sync.Mutex
	fnEvent    func(int, any)
	webhookUrl string
}

// GetInstance 返回单例实例
func GetInstance() *openWRT {
	once.Do(func() {
		instance = &openWRT{
			clients: make(map[string]*DHCPLease),
			nicks:   make(map[string]*NickEntry),
		}
		instance.init()
		glog.Println("Singleton instance created")
	})
	return instance
}

func (this *openWRT) init() {
	if u.IsMacOs() {
		return
	}
	this.webhookUrl = this.GetWebHook()
	this.initClients()
	time.Sleep(time.Second)
	go this.subscribeSysLog()
	go this.subscribeHostapd()
	go this.subscribeArpEvent()
	go this.subscribeDnsmasq()
	go this.subscribeFsnotify()
	//time.AfterFunc(time.Second*10, func() {
	//	this.ResetClients()
	//})
}

func (this *openWRT) getName(macAddr string) string {
	temp := this.clients[macAddr]
	if temp != nil {
		if temp.Nick != nil && temp.Nick.Name != "" {
			return temp.Nick.Name
		} else {
			return temp.Hostname
		}
	} else {
		return macAddr
	}
}

func (this *openWRT) getDeviceName(macAddr string) (string, string) {
	temp := this.clients[macAddr]
	if temp != nil {
		if temp.Nick != nil && temp.Nick.Name != "" {
			return temp.Nick.Name, temp.IP
		} else {
			return temp.Hostname, temp.IP
		}
	} else {
		return macAddr, ""
	}
}

func (this *openWRT) subscribeSysLog() {
	for {
		err := subscribeSysLog(func(event *SysLogEvent) {
			//glog.Infof("SysLog事件:%+v", event)
			if event != nil && event.Mac != "" {
				this.updateDeviceStatus("SysLog事件", &DHCPLease{
					MAC:       event.Mac,
					Online:    event.Online,
					StartTime: event.Timestamp.UnixMilli(),
					Phy:       event.Phy,
				})
			}
		})
		if err != nil {
			glog.Error(fmt.Errorf("系统日志监听失败 %v", err))
			time.Sleep(time.Second * 10)
		}
	}

}

func (this *openWRT) getClient(macAddr string) *DHCPLease {
	if cls, ok := this.clients[macAddr]; ok {
		return cls
	}
	return nil
}

func (this *openWRT) subscribeHostapd() {
	for {
		err := SubscribeHostapd(func(device *HostapdDevice) {
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
					//glog.Infof("Hostapd事件:%+v", device)
					this.updateDeviceStatus("Hostapd事件", &DHCPLease{
						MAC:       device.Address,
						Signal:    device.Signal,
						Freq:      device.Freq,
						StartTime: device.Timestamp.UnixMilli(),
						Online:    device.DataType == 0,
					})
				}
			}
		})
		if err != nil {
			glog.Error(fmt.Errorf("订阅失败 Hostapd %v", err))
			time.Sleep(time.Second * 10)
			glog.Error(fmt.Errorf("重新订阅 Hostapd "))
		}
	}
}

//func (this *openWRT) subscribeArp() {
//	for {
//		err := SubscribeArp(time.Second*10, func(entry *ARPEntry) {
//			glog.Infof("Arp事件:%+v", entry)
//			if entry != nil && entry.MAC != nil {
//				this.updateDeviceStatus("Arp事件", &DHCPLease{
//					MAC:       entry.MAC.String(),
//					IP:        entry.IP.String(),
//					StartTime: entry.Timestamp.UnixMilli(),
//					Phy:       entry.Interface,
//				})
//			}
//
//		})
//		if err != nil {
//			glog.Error(fmt.Errorf("订阅失败 Hostapd %v", err))
//			time.Sleep(time.Second * 10)
//		}
//	}
//}

func (this *openWRT) subscribeArpEvent() {
	SubscribeArpCache(time.Second*10, func(entrys map[string]*ARPEntry) {
		if entrys != nil {
			for mac, entry := range entrys {
   dhcp:=&DHCPLease{
							MAC:       entry.MAC.String(),
							IP:        entry.IP.String(),
							StartTime: entry.Timestamp.UnixMilli(),
							Phy:       entry.Interface,
						}
				if v, ok := this.clients[mac]; ok {
					if v.Online != (entry.Flags == 2) {
						//glog.Infof("Arp事件:%+v", entry)
						this.updateDeviceStatus("Arp事件", dhcp)
					}
				}else{
this.clients[mac]=dhcp
}
			}
		}
	})
}
func (this *openWRT) subscribeDnsmasq() {
	err := SubscribeDnsmasq(func(device *DnsmasqDevice) {
		//glog.Infof("Dnsmasq事件:%+v", device)
		if device != nil && device.Mac != "" {
			this.updateDeviceStatus("Dnsmasq事件", &DHCPLease{
				MAC:       device.Mac,
				IP:        device.Ip,
				Hostname:  device.Name,
				StartTime: device.Timestamp.UnixMilli(),
				Phy:       device.Interface,
				Online:    true,
			})
		}
	})
	if err != nil {
		glog.Error(fmt.Errorf("subscribeHostapd Error:%v", err))
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
			glog.Debug("listenFsnotify:", event)
			if event.Has(fsnotify.Write) {
				//filePath := event.Name
				if strings.EqualFold(event.Name, dhcpLeasesFilePath) {
					this.updateClientsByDHCP()
				}
				this.webUpdateAll(this.GetClients())
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			glog.Error("error:", err)
		}
	}
}

func (this *openWRT) subscribeFsnotify() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Error(fmt.Errorf("创建监控器失败 %v", err))
	}
	//go this.listenFsnotify(watcher)
	err = watcher.Add(dhcpLeasesFilePath)
	if err != nil {
		glog.Error(fmt.Errorf("watcher add err %v", err))
	}
	this.listenFsnotify(watcher)
}

func (this *openWRT) initClients() {
	dataMap, err := this.initData()
	if err != nil {
		glog.Errorf("initClients Error:%v", err)
		time.Sleep(5 * time.Second)
		glog.Error("5 seconds later and try...")
		this.initClients()
	}
	this.clients = dataMap
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
		glog.Printf("DHCP更新客户端 %+v\n", len(clientArray))
		arpMap, e1 := getClientsByArp(brLanString)
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
				if client.StartTime > 0 {
					v.StartTime = client.StartTime
				}
				v.Online = client.Online
			} else {
				p.clients[mac] = client
			}
			p.updateDeviceStatus("dhcp", client)
		}
	}
}

func (this *openWRT) updateStatusList(macAddr string, newList []*Status) {
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

func (this *openWRT) initData() (map[string]*DHCPLease, error) {
	arpList, e1 := getClientsByArp(brLanString)
	glog.Println("✅ arpList：")
	for _, temp := range arpList {
		glog.Debugf("%+v", temp)
	}
	if e1 == nil {
		dhcpMap, e2 := getClientsByDhcp()
		glog.Println("✅ dhcpMap：")
		for _, temp := range dhcpMap {
			glog.Debugf("%+v", temp)
		}
		sysLogMap, e3 := getStatusFromSysLog()
		nickMap, e4 := getNickData()

		this.nicks = nickMap
		glog.Println("✅ nickMap：")
		for _, temp := range nickMap {
			glog.Debugf("%+v", temp)
		}
		stcMap, e5 := getStaticIpMap()
		if e4 != nil {
			nickMap = map[string]*NickEntry{}
		}
		dataMap := make(map[string]*DHCPLease)
		for _, entry := range arpList {
			mac := entry.MAC.String()
			item := &DHCPLease{
				IP:     entry.IP.String(),
				MAC:    mac,
				Phy:    entry.Interface,
				Online: entry.Flags == 2,
			}
			if e2 == nil {
				if lease, ok := dhcpMap[mac]; ok {
					item.StartTime = lease.StartTime
					item.Hostname = lease.Hostname
				}
			}
			if e3 == nil {
				//item.StatusList = status[mac]
				list := sysLogMap[mac]
				//_ = setStatusByMac(mac, list)
				this.updateStatusList(mac, list)
			}
			if e4 == nil {
				if nick, ok := nickMap[mac]; ok {
					//item.NickName = nick.Name
					item.Nick = nick
				}
			} else {
				nick := &NickEntry{
					StartTime: item.StartTime,
					Hostname:  item.Hostname,
					IP:        item.IP,
					MAC:       mac,
				}
				nickMap[mac] = nick
			}

			if e5 == nil {
				if ip, ok := stcMap[mac]; ok {
					item.Static = ip
				}
			}
			dataMap[mac] = item
		}
		if e4 != nil {
			err := updateNicksData(nickMap)
			if err != nil {
				glog.Errorf("NickData Save Error:%v", err)
			}
		}

		glog.Println("✅ dataMap：")
		for _, temp := range dataMap {
			glog.Debugf("%+v", temp)
		}
		return dataMap, nil
	}
	return nil, e1
}
