package ntfy

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/u"
	"net/http"
	"sync"
	"time"
)

var (
	instance *Ntfy
	once     sync.Once
)

type Ntfy struct {
	mu       sync.Mutex
	resp     *http.Response
	isClosed bool
	ctx      context.Context
	cancel   context.CancelFunc
	fnArray  []func(string)
	info     *u.NtfyInfo
}

// GetInstance 返回单例实例
func GetInstance() *Ntfy {
	once.Do(func() {
		instance = &Ntfy{
			resp:     nil,
			isClosed: false,
			fnArray:  make([]func(string), 0),
			info:     nil,
		}
	})
	return instance
}

func (this *Ntfy) Start(info *u.NtfyInfo) {
	if info == nil {
		return
	}
	this.info = info
	this.isClosed = false
	this.subscribe(info)
}

func (this *Ntfy) AddFunc(fn func(string)) {
	this.fnArray = append(this.fnArray, fn)
}

func (this *Ntfy) Stop() {
	if this.resp != nil {
		this.isClosed = true
		glog.Debugf("Ntfy stop url: %s", this.resp.Request.URL.String())
		_ = this.resp.Body.Close()
		if this.cancel != nil {
			this.cancel()
		}
	}
}

func (this *Ntfy) Subscribe(address, topic string, username, password string) error {
	//req, e := http.NewRequest("GET", fmt.Sprintf("%s/%s/json?poll=1", address, topic), nil)
	req, e := http.NewRequest("GET", fmt.Sprintf("%s/%s/json", address, topic), nil)
	if e != nil {
		return e
	}
	req.SetBasicAuth(username, password)
	//req.Header.Set("Authorization", "Basic tk_trk4agho2")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	this.resp = resp
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	glog.Info("NTFY订阅成功", address, resp.StatusCode)
	scanner := bufio.NewScanner(resp.Body)
	if e1 := scanner.Err(); e1 != nil {
		return e1
	}
	for scanner.Scan() {
		text := scanner.Text()
		//glog.Info(text)
		if this.fnArray != nil && len(this.fnArray) > 0 {
			for _, fn := range this.fnArray {
				fn(text)
			}
		}
	}
	return scanner.Err()
}

func (this *Ntfy) subscribe(info *u.NtfyInfo) {
	this.ctx, this.cancel = context.WithCancel(context.Background())
	for {
		select {
		case <-this.ctx.Done():
			glog.Debug("ntfy 监听退出...")
			return
		default:
			err := this.Subscribe(info.Address, info.Topic, info.Username, info.Password)
			if err != nil {
				glog.Error(err)
				if this.isClosed {
					return
				}
				time.Sleep(time.Second * 10)
			}
		}
	}
}

func (this *Ntfy) Publish(data *u.NtfyEventData) error {
	if data == nil {
		return errors.New("NtfyEventData is nil")
	}
	if this.info == nil {
		return errors.New("NtfyInfo is nil")
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, e := http.NewRequest("POST", fmt.Sprintf("%s", this.info.Address), bytes.NewReader(body))
	if e != nil {
		return e
	}
	req.SetBasicAuth(this.info.Username, this.info.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	this.resp = resp
	glog.Debug("Ntfy publish: ", resp.Status)
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}
