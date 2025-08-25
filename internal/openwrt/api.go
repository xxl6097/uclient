package openwrt

import (
	"errors"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/uclient/internal/ntfy"
	"github.com/xxl6097/uclient/internal/u"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (this *openWRT) initData() error {
	staInfo := GetStaInfo(strings.Contains(this.ulistString, "ahsapd.sta"))
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
		var sysLogMap map[string][]*Status
		if strings.Contains(this.ulistString, "hostapd") {
			sysLogMap, _ = getStatusFromSysLog()
		}
		nickMap, e4 := getNickData()
		if e4 == nil {
			this.nicks = nickMap
			glog.Debug("\n✅ nickMap：")
			for _, temp := range nickMap {
				if temp != nil {
					glog.Debugf("%+v %+v", temp, temp.WorkType)
				}
			}
		}

		stcMap, e5 := getStaticIpMap()
		if e4 != nil {
			nickMap = map[string]*NickEntry{}
		}
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
					//if lease.StartTime <= 0 {
					//	lease.StartTime = glog.Now().UnixMilli()
					//}
					//item.StartTime = lease.StartTime
					item.Hostname = lease.Hostname
				}
			}
			if staInfo != nil {
				//glog.Debugf("---->%+v", staInfo)
				sta := staInfo[mac]
				if sta != nil {
					item.Vendor = sta.StaVendor
					if item.Hostname == "" || item.Hostname == "*" {
						item.Hostname = sta.HostName
					}
					num, _ := strconv.Atoi(sta.Rssi)
					item.Signal = num
					item.StaType = sta.StaType
					item.Ssid = sta.Ssid
					if item.StartTime == 0 {
						item.StartTime = sta.Timestamp
					}
				}
			}

			if sysLogMap != nil && len(sysLogMap) > 0 {
				list := sysLogMap[mac]
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

			//if item.IP != "" {
			//	item.Online = u.Ping(item.IP)
			//}
			//dataMap[mac] = item
			if v, ok := this.clients[mac]; ok {
				v.IP = item.IP
				v.MAC = item.MAC
				v.Phy = item.Phy
				v.Hostname = item.Hostname
				v.Online = item.Online
				v.Signal = item.Signal
				v.Freq = item.Freq
				v.Nick = item.Nick
				v.Static = item.Static
			} else {
				this.clients[mac] = item
			}
		}
		if e4 != nil {
			err := updateNicksData(nickMap)
			if err != nil {
				glog.Errorf("NickData Save Error:%v", err)
			}
		}

		glog.Debug("\n✅ dataMap：")
		for _, temp := range this.clients {
			glog.Debugf("%+v", temp)
		}
		return nil
	}
	return e1
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

func (this *openWRT) GetStaticIpMap() ([]*DHCPHost, error) {
	arr, err := GetUCIOutput()
	if err != nil {
		glog.Printf("GetStaticIpMap Error: %v\n", err)
		return nil, err
	}
	//for _, entry := range arr {
	//	this.clients[entry.MAC].Static = entry
	//}
	glog.Println("GetStaticIpMap：", len(arr))
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

func (this *openWRT) SetNtfy(info *u.NtfyInfo) error {
	if info == nil {
		return fmt.Errorf("NtfyInfo is nil")
	}
	go ntfy.GetInstance().Start(info)
	return utils.SaveWithGob[u.NtfyInfo](*info, ntfyFilePath)
}

func (this *openWRT) GetWebHook() string {
	data, err := os.ReadFile(webhookFilePath)
	if err != nil {
		return ""
	}
	return string(data)
}

func (this *openWRT) GetWorkTimeAndCaculate(mac string) ([]*Work, error) {
	if mac == "" {
		return nil, fmt.Errorf("mac is empty")
	}
	if v, ok := this.clients[mac]; ok {
		if v.Nick != nil {
			return getWorkTimeAndCaculate(mac, v.Nick.WorkType)
		}
	}
	return nil, fmt.Errorf("client not found mac")
}

func (this *openWRT) CheckFile(file string) {
	if !u.IsFileExist(file) {
		tryCount := 0
		for {
			time.Sleep(time.Second * 10)
			if u.IsFileExist(file) {
				return
			}
			tryCount++
			if tryCount > RE_REY_MAX_COUNT {
				glog.Error("监听  失败，超过最大重试次数")
				return
			}
		}
	}
}
