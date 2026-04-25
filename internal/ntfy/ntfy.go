package ntfy

import (
	"fmt"
	"sync"
	"time"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/uclient/internal/u"
)

var (
	instance *Ntfy
	once     sync.Once
)

type Ntfy struct {
	client *Client
	fn     func(string) string
}

// GetInstance 返回单例实例
func GetInstance() *Ntfy {
	once.Do(func() {
		instance = &Ntfy{client: nil, fn: nil}
	})
	return instance
}

func (this *Ntfy) Start(cfg *u.NtfyInfo) {
	if cfg == nil {
		return
	}
	this.client = NewClientWithAuth(cfg.Address, cfg.Username, cfg.Password, time.Second*10)
	this.subscribe(cfg.Topic)
}

func (this *Ntfy) Stop() {
	if this.client != nil {
		this.client.Stop()
	}
}
func (this *Ntfy) SetFunc(fn func(string) string) {
	this.fn = fn
}
func (this *Ntfy) subscribe(topic string) {
	this.client.ListenWithRetry(topic, func(msg *Message) string {
		//fmt.Printf("\n[服务端] 收到请求：ID=%s | 内容=%s\n", msg.ID, msg.Message)
		z.L().Debug(fmt.Sprintf("[服务端] 收到请求：ID=%s | 内容=%s\n", msg.ID, msg.Message))
		//this.client.DispatchResponse(msg)
		// 返回结果
		var res string
		if this.fn != nil {
			res = this.fn(msg.Message)
		}
		return res
	})
}

func (this *Ntfy) Publish(data *Message) error {
	_, err := this.client.Send(data)
	if err != nil {
		return err
	}
	return nil
}

func (this *Ntfy) PublishSync(title string, data *SyncMessage) (*Message, error) {
	return this.client.SendSync(title, data)
}
