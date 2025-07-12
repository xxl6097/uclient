package openwrt

import (
	"errors"
	"github.com/xxl6097/uclient/internal/u"
	"sort"
)

func (this *openWRT) Listen(fn func([]*DHCPLease)) {
	this.fnWatcher = func() {
		if fn != nil {
			fn(this.GetClients())
		}
	}
}

func (this *openWRT) GetStatusByMac(mac string) []*Status {
	list := getStatusByMac(mac)
	if list == nil || len(list) == 0 {
		return nil
	}
	sort.Slice(list, func(i, j int) bool {
		a := list[i]
		a.TimeLine = u.TimestampFormat(a.Timestamp)
		b := list[j]
		b.TimeLine = u.TimestampFormat(b.Timestamp)
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

func (this *openWRT) UpdateNickName(obj *DHCPLease) error {
	if obj == nil {
		return errors.New("DHCPLease obj is nill")
	}
	mac := obj.MAC
	nick, ok := this.nickMap[mac]
	if ok {
		nick.Name = obj.NickName
	} else {
		nick = &NickEntry{
			Hostname: obj.Hostname,
			IP:       obj.IP,
		}
	}
	if v, ok := this.clients[mac]; ok {
		v.NickName = obj.NickName
	}
	if this.fnWatcher != nil {
		this.fnWatcher()
	}
	return updateNickData(mac, nick)
}

func (this *openWRT) ResetClients() {
	this.initClients()
	if this.fnWatcher != nil {
		this.fnWatcher()
	}
}
