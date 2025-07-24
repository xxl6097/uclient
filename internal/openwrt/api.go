package openwrt

import (
	"errors"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"os"
	"sort"
	"time"
)

//	func (this *openWRT) Listen(fn func([]*DHCPLease)) {
//		this.fnWatcher = func() {
//			if fn != nil {
//				fn(this.GetClients())
//			}
//		}
//	}
//func (this *openWRT) ListenOne(fn func(int, *DHCPLease)) {
//	this.fnNewOne = func(dataType int, cls *DHCPLease) {
//		if fn != nil {
//			fn(dataType, cls)
//		}
//	}
//}

func (this *openWRT) initData() (map[string]*DHCPLease, error) {
	this.webhookUrl = this.GetWebHook()
	arpList, e1 := getClientsByArp(brLanString)
	if e1 == nil {
		glog.Debug("\n✅ arpList：")
		for _, temp := range arpList {
			glog.Debugf("%+v", temp)
		}
	}
	if e1 == nil {
		dhcpMap, e2 := getClientsByDhcp()
		if e2 == nil {
			this.leases = dhcpMap
			glog.Debug("\n✅ dhcpMap：")
			for _, temp := range dhcpMap {
				glog.Debugf("%+v", temp)
			}
		}

		sysLogMap, e3 := getStatusFromSysLog()
		nickMap, e4 := getNickData()
		if e4 == nil {
			this.nicks = nickMap
			glog.Debug("\n✅ nickMap：")
			for _, temp := range nickMap {
				glog.Debugf("%+v", temp)
			}
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
				this.updateUserTimeLineData(mac, list)
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

		glog.Debug("\n✅ dataMap：")
		for _, temp := range dataMap {
			glog.Debugf("%+v", temp)
		}
		return dataMap, nil
	}
	return nil, e1
}

func (this *openWRT) getClient(macAddr string) *DHCPLease {
	if cls, ok := this.clients[macAddr]; ok {
		return cls
	}
	return nil
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

func (this *openWRT) SetFunc(fn func(int, any)) {
	this.fnEvent = func(i int, obj any) {
		if fn != nil {
			fn(i, obj)
		}
	}
}

func (this *openWRT) GetStatusByMac(mac string) []*Status {
	list := getStatusByMac(mac)
	if list == nil || len(list) == 0 {
		return nil
	}
	sort.Slice(list, func(i, j int) bool {
		//a := list[i]
		//a.TimeLine = u.TimestampToDateTime(a.Timestamp)
		//b := list[j]
		//b.TimeLine = u.TimestampToDateTime(b.Timestamp)
		return list[i].Timestamp > list[j].Timestamp
	})
	return list
}

func (this *openWRT) GetDeviceTimeLineDatas(tempFilePath string) []*DeviceTimeLine {
	list := readTimeLineByMac(tempFilePath)
	if list == nil || len(list) == 0 {
		return nil
	}
	t := glog.Now()
	sort.Slice(list, func(i, j int) bool {
		a := u.UTC8ToTime(list[i].Timestamp)
		b := u.UTC8ToTime(list[j].Timestamp)
		du1 := t.Sub(a)
		du2 := t.Sub(b)
		ago1 := time.Duration(du1.Seconds()) * time.Second
		ago2 := time.Duration(du2.Seconds()) * time.Second
		list[i].Ago = ago1.String()
		list[j].Ago = ago2.String()
		list[i].DateTime = a.Format(time.DateTime)
		list[j].DateTime = b.Format(time.DateTime)
		return list[i].Timestamp > list[j].Timestamp
	})
	return list
}

func (this *openWRT) GetClients() []*DHCPLease {
	data := make([]*DHCPLease, 0)
	for _, cls := range this.clients {
		data = append(data, cls)
	}
	sort.Slice(data, func(i, j int) bool {
		// 在线状态优先：在线(true) > 离线(false)
		if data[i].Online != data[j].Online {
			return data[i].Online
		}
		return data[i].Hostname < data[j].Hostname
	})
	return data
}

func (this *openWRT) UpdateNickName(obj *NickEntry) error {
	if obj == nil {
		return errors.New("DHCPLease obj is nill")
	}
	mac := obj.MAC
	var nick *NickEntry
	client, ok := this.clients[mac]
	if ok {
		nick = client.Nick
	}
	if nick == nil {
		nick = &NickEntry{}
	}
	nick.Name = obj.Name
	nick.IsPush = obj.IsPush
	nick.WorkType = obj.WorkType
	nick.Hostname = obj.Hostname
	nick.IP = obj.IP
	nick.MAC = obj.MAC
	nick.StartTime = obj.StartTime
	if v, okk := this.clients[mac]; okk {
		v.Nick = nick
	}
	if this.fnEvent != nil {
		this.fnEvent(0, this.GetClients())
	}
	this.nicks[mac] = nick
	return updateNickData(mac, nick)
}

func (this *openWRT) ResetClients() {
	this.initClients()
	this.webUpdateAll(this.GetClients())
}

func (this *openWRT) GetStaticIpMap() ([]DHCPHost, error) {
	arr, err := GetUCIOutput()
	if err != nil {
		glog.Printf("Error: %v\n", err)
		return nil, err
	}
	for _, entry := range arr {
		this.clients[entry.MAC].Static = &entry
	}
	return arr, nil
}
func (this *openWRT) DeleteStaticIp(mac string) error {
	err := deleteStaticIpAddress(mac)
	if v, ok := this.clients[mac]; ok && v != nil {
		v.Static = nil
	}
	return err
}

func (this *openWRT) SetWebHook(webhookUrl string) error {
	if webhookUrl == "" {
		return fmt.Errorf("webhook is empty")
	}
	file, err := os.Create(webhookFilePath) // 文件不存在则创建，存在则截断
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入内容
	_, err = file.Write([]byte(webhookUrl))
	if err == nil {
		this.webhookUrl = webhookUrl
	}
	return err
}

func (this *openWRT) GetWebHook() string {
	data, err := os.ReadFile(webhookFilePath)
	if err != nil {
		return ""
	}
	return string(data)
}

func (this *openWRT) GetWorkTime(mac string) ([]*Work, error) {
	if mac == "" {
		return nil, fmt.Errorf("mac is empty")
	}
	if v, ok := this.clients[mac]; ok {
		if v.Nick != nil {
			return getWorkTime(mac, v.Nick.WorkType)
		}
	}
	return nil, fmt.Errorf("client not found mac")
}
