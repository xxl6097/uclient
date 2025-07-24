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
	leases     map[string]*DHCPLease
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
			leases:  make(map[string]*DHCPLease),
		}
		instance.init()
	})
	return instance
}

func (this *openWRT) init() {
	if u.IsMacOs() {
		return
	}
	this.initClients()
	go this.subscribeSysLog()
	go this.subscribeHostapd()
	go this.subscribeArpEvent()
	go this.subscribeDnsmasq()
	go this.subscribeFsnotify()
}

func (this *openWRT) initClients() {
	dataMap, err := this.initData()
	if err != nil {
		glog.Errorf("initClients Error:%v", err)
		time.Sleep(5 * time.Second)
		glog.Error("5 seconds later and try...")
		this.initClients()
	}
	if dataMap == nil || len(dataMap) == 0 {
		glog.Error("dataMap is empty, 1 seconds later and try...")
		time.Sleep(16 * time.Second)
		this.initClients()
	}
	this.clients = dataMap
}

func (this *openWRT) subscribeSysLog() {
	for {
		err := subscribeSysLog(func(event *SysLogEvent) {
			if event != nil && event.Mac != "" {
				glog.Infof("SysLog事件:%+v", event)
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
				this.updateDeviceStatus("SysLog事件", eve)
			}
		})
		if err != nil {
			glog.Error(fmt.Errorf("SysLog监听失败 %v", err))
			time.Sleep(time.Second * 5)
			glog.Error("重新g监听 SysLog")
		}
	}
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
					this.updateDeviceStatus("Hostapd事件", dhcp)
				}
			}
		})
		if err != nil {
			glog.Errorf("订阅失败 Hostapd %v", err)
			time.Sleep(time.Second * 5)
			glog.Error("重新订阅 Hostapd")
		}
	}
}

func (this *openWRT) subscribeArpEvent() {
	SubscribeArpCache(time.Second*10, func(entrys map[string]*ARPEntry) {
		if entrys != nil {
			for mac, entry := range entrys {
				dhcp := &DHCPLease{
					MAC:       entry.MAC.String(),
					IP:        entry.IP.String(),
					StartTime: entry.Timestamp.UnixMilli(),
					Phy:       entry.Interface,
				}
				if v, ok := this.leases[mac]; ok {
					if v.Hostname != "" {
						dhcp.Hostname = v.Hostname
					}
				}
				if v, ok := this.clients[mac]; ok {
					if v.Online != (entry.Flags == 2) {
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
func (this *openWRT) subscribeDnsmasq() {
	for {
		err := SubscribeDnsmasq(func(device *DnsmasqDevice) {
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
				this.updateDeviceStatus("Dnsmasq事件", dhcp)
			}
		})
		if err != nil {
			glog.Errorf("Dnsmasq订阅失败 %v", err)
			time.Sleep(time.Second * 5)
			glog.Error("重新订阅 Dnsmasq")
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
			p.leases[mac] = client
			p.updateDeviceStatus("dhcp", client)
		}
	}
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
