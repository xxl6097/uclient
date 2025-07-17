package openwrt

import (
	"errors"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"os"
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
		//a := list[i]
		//a.TimeLine = u.TimestampToDateTime(a.Timestamp)
		//b := list[j]
		//b.TimeLine = u.TimestampToDateTime(b.Timestamp)
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

func (this *openWRT) notifyWebhookMessage(client *DHCPLease) {
	if this.webhookUrl == "" {
		return
	}
	if client == nil {
		return
	}
	if client.Nick == nil {
		return
	}
	if !client.Nick.IsPush {
		return
	}
	markdown := make(map[string]interface{})
	var name, title string
	if client.Nick != nil && client.Nick.Name != "" {
		name = client.Nick.Name
	} else {
		name = client.Hostname
	}
	if client.Online {
		title = fmt.Sprintf("%s上线啦", name)
	} else {
		title = fmt.Sprintf("%s离线了", name)
	}
	format := "#### %s \n - 名称：%s\n - IP地址：%s \n- Mac地址：%s \n- 时间：%s \n"
	markdown["title"] = title
	markdown["text"] = fmt.Sprintf(format, title, name, client.IP, client.MAC, u.TimestampToDateTime(client.StartTime))
	payload := map[string]interface{}{"msgtype": "markdown", "markdown": markdown}
	_ = WebHookMessage(this.webhookUrl, payload)
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
