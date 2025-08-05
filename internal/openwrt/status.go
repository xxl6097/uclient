package openwrt

import (
	"context"
	"encoding/json"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"os"
	"time"
)

func (this *openWRT) StartStatus() {
	glog.Debug("启动tatus", this.statusRuning)
	if this.statusRuning {
		return
	}
	this.subscribeStatus()
}
func (this *openWRT) StopStatus() {
	glog.Debug("停止status", this.statusRuning)
	this.statusRuning = false
	if this.cancel != nil {
		this.cancel()
	}
}
func (this *openWRT) subscribeStatus() {
	this.statusRuning = true
	this.ctx, this.cancel = context.WithCancel(context.Background())
	glog.Debug("开始计算", this.statusRuning)
	for {
		select {
		case <-this.ctx.Done():
			glog.Debug("退出状态计算", this.statusRuning)
			return
		default:
			this.readStatus()
			if !this.statusRuning {
				glog.Debug("退出状态计算", this.statusRuning)
				return
			}
			time.Sleep(time.Second * 5)
		}
	}
}

func (this *openWRT) mergeStatus(list []*u.Device) {
	if list == nil || len(list) == 0 {
		return
	}
	for _, device := range list {
		if device == nil {
			continue
		}
		if device.MacAddress == "" {
			continue
		}

		mac := u.MacFormat(device.MacAddress)
		if v, ok := this.clients[mac]; ok {
			v.Signal = device.RSSI
			if device.Lan != "" {
				v.Ssid = device.Lan
			}
			if v.Device == nil {
				v.Device = device
			} else {
				upRate := float64(device.RxBytes-v.Device.RxBytes) / float64(5*2)
				downRate := float64(device.TxBytes-v.Device.TxBytes) / float64(5*2)
				v.UpRate = u.ByteCountSpeed(uint64(upRate))
				v.DownRate = u.ByteCountSpeed(uint64(downRate))
				this.webUpdateAll(this.GetClients())
				v.Device = device
			}
		} else {
			glog.Debug("mergeStatus 不存在", device)
		}
	}
}

func (this *openWRT) readStatus() {
	_, err := os.Stat(hetsysinfoFilePath)
	// 判断是否为文件不存在的错误
	if os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(hetsysinfoFilePath)
	if err != nil {
		return
	}
	var res u.DeviceStatus
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}

	this.mergeStatus(this.decodeStatus(&res))
}

func (this *openWRT) DecodeDevice(data []byte) []*u.Device {
	var res u.DeviceStatus
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil
	}
	return this.decodeStatus(&res)
}

func (this *openWRT) decodeStatus(sd *u.DeviceStatus) []*u.Device {
	if sd == nil {
		return nil
	}
	list := make([]*u.Device, 0)
	if sd.PORTINFO != nil {
		info := sd.PORTINFO
		if info.PORT0 != nil {
			if info.PORT0.IPADDR != "" {
				info.PORT0.Lan = "PORT0"
				list = append(list, info.PORT0)
			}
		}
		if info.PORT1 != nil {
			if info.PORT1.IPADDR != "" {
				info.PORT0.Lan = "PORT1"
				list = append(list, info.PORT1)
			}
		}
		if info.PORT2 != nil {
			if info.PORT2.IPADDR != "" {
				info.PORT0.Lan = "PORT2"
				list = append(list, info.PORT2)
			}
		}
	}
	if sd.G2 != nil {
		if sd.G2.Ra0 != nil && sd.G2.Ra0.Stainfo != nil {
			list = append(list, sd.G2.Ra0.Stainfo...)
		}
		if sd.G2.Ra1 != nil && sd.G2.Ra1.Stainfo != nil {
			list = append(list, sd.G2.Ra1.Stainfo...)
		}
	}
	if sd.G5 != nil {
		if sd.G5.Rax0 != nil && sd.G5.Rax0.Stainfo != nil {
			list = append(list, sd.G5.Rax0.Stainfo...)
		}
	}
	return list
}
