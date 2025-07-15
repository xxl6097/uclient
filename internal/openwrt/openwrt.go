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
	//nickMap            map[string]*NickEntry
	clients            map[string]*DHCPLease
	sysLogClientStatus map[string]bool
	mu                 sync.RWMutex
	fnWatcher          func()
	webhookUrl         string
}

// GetInstance 返回单例实例
func GetInstance() *openWRT {
	once.Do(func() {
		instance = &openWRT{
			//nickMap:            make(map[string]*NickEntry),
			sysLogClientStatus: make(map[string]bool),
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
	go this.initListenSysLog()
	this.initListenFsnotify()
}

func (this *openWRT) initListenSysLog() {
	err := listenSysLog(func(timestamp int64, macAddr string, phy string, status int, rawData string) {
		switch status {
		case 0:
			glog.Debugf("设备【%s】断开了", macAddr)
			glog.Debug(rawData)
			this.updateClientsBySysLog(timestamp, macAddr, phy, false)
			break
		case 1:
			glog.Debugf("设备【%s】连上了", macAddr)
			glog.Debug(rawData)
			this.updateClientsBySysLog(timestamp, macAddr, phy, true)
			break
		default:
			//glog.Warnf("未知数据 %v", rawData)
			break
		}

	})
	if err != nil {
		glog.Error(fmt.Errorf("listenSysLog Error:%v", err))
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
			glog.Debug("event:", event)
			if event.Has(fsnotify.Write) {
				//filePath := event.Name
				if strings.EqualFold(event.Name, dhcpLeasesFilePath) {
					this.updateClientsByDHCP()
				}
				if this.fnWatcher != nil {
					this.fnWatcher()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			glog.Error("error:", err)
		}
	}
}

func (this *openWRT) initListenFsnotify() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Error(fmt.Errorf("创建监控器失败 %v", err))
	}
	go this.listenFsnotify(watcher)
	err = watcher.Add(dhcpLeasesFilePath)
	if err != nil {
		glog.Error(fmt.Errorf("watcher add err %v", err))
	}
}

func (this *openWRT) initClients() {
	dataMap, err := this.initClientsFromDHCPAndArpAndSysLogAndNick()
	if err != nil {
		glog.Errorf("initClients Error:%v", err)
		time.Sleep(5 * time.Second)
		glog.Error("5 seconds later and try...")
		this.initClients()
	}
	this.clients = dataMap
}

func (this *openWRT) updateClientsBySysLog(timestamp int64, macAddr string, phy string, status bool) {
	if cls, ok := this.clients[macAddr]; ok {
		cls.Online = status
		cls.Phy = phy
		cls.StartTime = timestamp
		this.sysLogClientStatus[macAddr] = status
		glog.Infof("updateClientsBySysLog:%v", cls)
		this.notifyWebhookMessage(cls)
		if this.fnWatcher != nil {
			this.fnWatcher()
		}
	}
	s := Status{
		Timestamp: timestamp,
		Connected: status,
	}
	this.updateStatusList(macAddr, []*Status{&s})
}

func (p *openWRT) updateClientsByDHCP() {
	clientArray, err := getClientsByDhcp()
	if err != nil {
		glog.Println(fmt.Errorf("getClientsByDhcp Error:%v", err))
	} else {
		glog.Printf("DHCP更新客户端 %+v\n", len(clientArray))
		arpMap, e1 := getClientsByArp(brLanString)
		for _, client := range clientArray {
			mac := client.MAC
			//读取/tmp/dhcp.leases列表，这个列表没有状态，需要从syslog中获取
			if status, okk := p.sysLogClientStatus[mac]; okk {
				client.Online = status
			} else {
				if e1 == nil && arpMap != nil {
					itemData := arpMap[mac]
					if itemData != nil {
						client.Online = itemData.Flags == 2
					}
				}
			}
			var nick *NickEntry
			if v, ok := p.clients[mac]; ok && v != nil && v.Nick != nil {
				nick = v.Nick
			}

			//读取/tmp/dhcp.leases列表，获取名称和昵称
			//if p.clients != nil && p.clients.nickMap != nil {
			//	if nick, ok := p.nickMap[mac]; ok {
			//		client.NickName = nick.Name
			//		if nick.Hostname != "" && nick.Hostname != "*" {
			//			client.Hostname = nick.Hostname
			//		}
			//	}
			//}
			//缓存列表，存在就更新状态，不存在就添加
			if v, ok := p.clients[mac]; ok {
				if client.Hostname != "" && client.Hostname != "*" {
					v.Hostname = client.Hostname
				}
				v.IP = client.IP
				v.StartTime = client.StartTime
				//v.NickName = client.NickName
				v.Nick = nick
				v.Online = client.Online

			} else {
				p.clients[mac] = client
			}
		}
	}
}

//func (this *openWRT) updateClientsByARP() {
//	clientArray, err := getClientsByArp(brLanString)
//	if err != nil {
//		glog.Println(fmt.Errorf("getClientsByArp Error:%v", err))
//	} else {
//		glog.Printf("ARP更新客户端 %+v\n", len(clientArray))
//		for _, client := range clientArray {
//			mac := client.MAC.String()
//			item := &DHCPLease{
//				IP:     client.IP.String(),
//				MAC:    mac,
//				Online: client.Flags == 2,
//			}
//			if this.nickMap != nil {
//				if nick, ok := this.nickMap[mac]; ok {
//					item.NickName = nick.Name
//					if nick.Hostname != "" && nick.Hostname != "*" {
//						item.Hostname = nick.Hostname
//					}
//				}
//			}
//
//			if v, ok := this.clients[mac]; ok {
//				if item.Hostname != "" && item.Hostname != "*" {
//					v.Hostname = item.Hostname
//				}
//				v.IP = item.IP
//				v.NickName = item.NickName
//				v.Online = item.Online
//			} else {
//				this.clients[mac] = item
//			}
//		}
//	}
//}

func (this *openWRT) updateStatusList(macAddr string, newList []*Status) {
	this.mu.Lock()
	defer this.mu.Unlock()
	list := getStatusByMac(macAddr)
	if list == nil {
		list = newList
	} else {
		element := list[len(list)-1]
		if element != nil {
			for i, n := range newList {
				if n.Timestamp >= element.Timestamp {
					if n.Timestamp == element.Timestamp && n.Connected == element.Connected {
						continue
					}
					list = append(list, newList[i:]...)
					break
				}
			}
		}
	}
	if list == nil {
		return
	}
	size := len(list)
	if len(list) > MAX_SIZE {
		tempSize := size - MAX_SIZE
		list = list[tempSize:]
	}
	_ = setStatusByMac(macAddr, list)
}

func (this *openWRT) initClientsFromDHCPAndArpAndSysLogAndNick() (map[string]*DHCPLease, error) {
	entries, e1 := getClientsByArp(brLanString)
	if e1 == nil {
		leases, e2 := getClientsByDhcp()
		status, e3 := getStatusFromSysLog()
		nicks, e4 := getNickData()
		ips, e5 := getStaticIpMap()
		//this.nickMap = nicks
		glog.Errorf("getNickData Error:%v", e4)
		if e4 != nil {
			nicks = map[string]*NickEntry{}
		} else {
			for _, nick := range nicks {
				glog.Debugf("NickData:%+v", nick)
			}
		}
		dataMap := make(map[string]*DHCPLease)
		for _, entry := range entries {
			mac := entry.MAC.String()
			item := &DHCPLease{
				IP:     entry.IP.String(),
				MAC:    mac,
				Online: entry.Flags == 2,
			}
			if e2 == nil {
				if lease, ok := leases[mac]; ok {
					item.StartTime = lease.StartTime
					item.Hostname = lease.Hostname
				}
			}
			if e3 == nil {
				//item.StatusList = status[mac]
				list := status[mac]
				//_ = setStatusByMac(mac, list)
				this.updateStatusList(mac, list)
			}
			if e4 == nil {
				if nick, ok := nicks[mac]; ok {
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
				nicks[mac] = nick
			}

			if e5 == nil {
				if ip, ok := ips[mac]; ok {
					item.Static = ip
				}
			}
			dataMap[mac] = item
		}
		if e4 != nil {
			err := updateNicksData(nicks)
			if err != nil {
				glog.Errorf("NickData Save Error:%v", err)
			}
		}
		return dataMap, nil
	}
	return nil, e1
}
