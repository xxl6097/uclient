package openwrt

import "github.com/xxl6097/glog/glog"

//func (this *openWRT) checkMessage(macAddress string, s *Status) bool {
//	defer this.mu.Unlock()
//	this.mu.Lock()
//	if macAddress == "" {
//		glog.Error("openWRT checkMessage - macAddress is empty")
//		return false
//	}
//	if v, okk := this.clientStatus[macAddress]; okk {
//		if v == nil {
//			this.clientStatus[macAddress] = s
//			glog.Error("v is nil", macAddress)
//			return false
//		} else {
//			//t1 := glog.Now()
//			//t2 := u.UTC8ToTime(v.Timestamp)
//			//du := t1.Sub(t2)
//			this.clientStatus[macAddress] = s
//			glog.Debugf("checkMessage:%v %v %v %v", s.Connected, v.Connected) //, du.String(), du.Milliseconds()
//			if s.Connected == v.Connected {                                   // && du.Milliseconds() < 1000
//				return true
//			}
//			return false
//		}
//	} else {
//		this.clientStatus[macAddress] = s
//		glog.Error("macAddress not in clientStatus", macAddress)
//		return false
//	}
//}

func (this *openWRT) updateDeviceStatus(typeEvent string, device *DHCPLease) {
	defer this.mu.Unlock()
	this.mu.Lock()
	if device == nil {
		return
	}
	macAddress := device.MAC
	if macAddress == "" {
		return
	}
	cls := this.getClient(macAddress)
	//glog.Debugf("新数据：%+v，老数据：%+v", device, cls)
	if cls == nil {
		nickMap, e2 := getNickData()
		if e2 == nil {
			this.nicks = nickMap
			device.Nick = nickMap[macAddress]
		}
		this.clients[macAddress] = device
		cls = device
	} else {
		if device.Online == cls.Online {
			//glog.Warnf("[%s]状态相同，不更新，%s[%s] 旧：%v,新：%v", typeEvent, cls.Hostname, cls.MAC, cls.Online, device.Online)
			return
		}
		cls.Online = device.Online
		if device.Signal != 0 {
			cls.Signal = device.Signal
		}
		if device.Freq != 0 {
			cls.Freq = device.Freq
		}
		if device.StartTime != 0 {
			cls.StartTime = device.StartTime
		}
		if device.IP != "" {
			cls.IP = device.IP
		}
		if device.Phy != "" {
			cls.Phy = device.Phy
		}
		if device.Hostname != "" {
			cls.Hostname = device.Hostname
		}
		//需要webnotify通知、钉钉notify、签到
		this.ddingWorkSign(cls)
	}
	s := Status{}
	s.Timestamp = device.StartTime
	s.Connected = device.Online
	this.ddingNotify(cls)
	this.webNotify(cls)
	this.updateUserTimeLineData(macAddress, []*Status{&s})
	glog.Debugf("%s 状态更新 %+v", typeEvent, device)
}

//func (this *openWRT) updateStatusByHostapd(device *HostapdDevice) {
//	if device == nil {
//		return
//	}
//	var macAddr string
//	macAddr = device.Address
//	cls := this.getClient(macAddr)
//	if device.DataType == 2 {
//		if cls != nil {
//			cls.Signal = device.Signal
//			cls.Freq = device.Freq
//			cls.StartTime = device.Timestamp.UnixMilli()
//			//这里只是更新信号，不在web上notify通知
//			this.webUpdateOne(cls)
//		} else {
//			//到这里，说明这个设备没有连上路由器，cls为nil
//		}
//	} else {
//		s := Status{}
//		s.Timestamp = device.Timestamp.UnixMilli()
//		s.Connected = device.DataType == 0
//
//		glog.Infof("【监听hostapd-1】:%v %v %+v ", s.Connected, u.TimestampToDateTime(s.Timestamp), device.Address)
//		if this.checkMessage(macAddr, &s) {
//			return
//		}
//		glog.Infof("【监听hostapd-2】:%v %v %+v", s.Connected, u.TimestampToDateTime(s.Timestamp), device.Address)
//		//在线、离线事件
//		if cls != nil {
//			cls.Signal = device.Signal
//			cls.Freq = device.Freq
//			cls.Online = s.Connected
//			cls.StartTime = s.Timestamp
//			//需要webnotify通知、钉钉notify、签到
//			this.ddingWorkSign(cls)
//		} else {
//			cls = &DHCPLease{
//				MAC:       macAddr,
//				StartTime: s.Timestamp,
//				Online:    s.Connected,
//				Signal:    device.Signal,
//				Freq:      device.Freq,
//			}
//			//需要webnotify通知、钉钉notify
//		}
//		this.ddingNotify(cls)
//		this.webNotify(cls)
//		this.updateUserTimeLineData(macAddr, []*Status{&s})
//	}
//}
//
//func (this *openWRT) updateStatusBySysLog(sysEvent *SysLogEvent) {
//	if sysEvent == nil {
//		return
//	}
//	s := Status{}
//	var macAddr string
//	s.Timestamp = sysEvent.Timestamp.UnixMilli()
//	s.Connected = sysEvent.Online
//	macAddr = sysEvent.Mac
//
//	glog.Infof("【监听日志1】:%v %v %+v", s.Connected, u.TimestampToDateTime(s.Timestamp), sysEvent)
//	if this.checkMessage(macAddr, &s) {
//		return
//	}
//	glog.Infof("【监听日志2】:%v %v %+v", s.Connected, u.TimestampToDateTime(s.Timestamp), sysEvent)
//	cls := this.getClient(macAddr)
//	if cls != nil {
//		cls.Online = s.Connected
//		cls.Phy = sysEvent.Phy
//		cls.StartTime = s.Timestamp
//		//需要web上notify通知、webhook通知、签到
//		this.ddingWorkSign(cls)
//	} else {
//		cls = &DHCPLease{
//			MAC:       macAddr,
//			StartTime: s.Timestamp,
//			Online:    s.Connected,
//		}
//		//需要web上notify通知和webhook通知
//	}
//	this.ddingNotify(cls)
//	this.webNotify(cls)
//	this.updateUserTimeLineData(macAddr, []*Status{&s})
//}
//
//func (this *openWRT) updateStatusByDnsmasq(dnsData *DnsmasqDevice) {
//	if dnsData == nil {
//		return
//	}
//	s := Status{}
//	var macAddr string
//	s.Timestamp = dnsData.Timestamp.UnixMilli()
//	s.Connected = true
//
//	if this.checkMessage(macAddr, &s) {
//		return
//	}
//	macAddr = dnsData.Mac
//	cls := this.getClient(macAddr)
//	//需要web上notify通知、webhook通知
//	glog.Infof("DNS监听:%v %+v", u.TimestampToDateTime(s.Timestamp), dnsData)
//	if cls != nil {
//		cls.Hostname = dnsData.Name
//		cls.Phy = dnsData.Interface
//		cls.IP = dnsData.Ip
//		cls.StartTime = s.Timestamp
//		cls.Online = s.Connected
//	} else {
//		cls = &DHCPLease{}
//		cls.Hostname = dnsData.Name
//		cls.Phy = dnsData.Interface
//		cls.IP = dnsData.Ip
//		cls.StartTime = s.Timestamp
//		cls.Online = s.Connected
//	}
//	this.ddingWorkSign(cls)
//	this.ddingNotify(cls)
//	this.webNotify(cls)
//	this.updateUserTimeLineData(macAddr, []*Status{&s})
//}
