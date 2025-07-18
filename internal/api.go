package internal

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/github"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/uclient/internal/iface"
	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/sse"
	"github.com/xxl6097/uclient/internal/u"
	"net/http"
	"path/filepath"
	"sync"
)

type Api struct {
	igs    igs.Service
	sseApi iface.ISSE
	pool   *sync.Pool // use sync.Pool caching buf to reduce gc ratio
}

func NewApi(igs igs.Service) *Api {
	github.Api().SetName("xxl6097", "openwrt-client-manager")
	sseApi := sse.NewServer()
	sseApi.Start()
	a := &Api{
		igs:    igs,
		sseApi: sseApi,
		pool: &sync.Pool{
			New: func() interface{} { return make([]byte, 32*1024) },
		},
	}
	openwrt.GetInstance().Listen(a.listen)
	openwrt.GetInstance().ListenOne(a.notifySSEEvent)
	return a
}

func (this *Api) listen(list []*openwrt.DHCPLease) {
	if len(list) >= 0 && this.sseApi != nil {
		eve := iface.SSEEvent{
			Event:   "update",
			Payload: list,
		}
		this.sseApi.Broadcast(eve)
	}
}

func (this *Api) notifySSEEvent(cls *openwrt.DHCPLease) {
	if cls != nil && this.sseApi != nil {
		eve := iface.SSEEvent{
			Event:   "update-one",
			Payload: cls,
		}
		this.sseApi.Broadcast(eve)
	}
}

func (this *Api) GetClients(w http.ResponseWriter, r *http.Request) {
	//req := utils.GetReqMapData(w, r)
	//glog.Warn(req)
	//glog.Warn("getClients---->", r)
	//cls, err := getClients()
	//
	//if err != nil {
	//	glog.Error("getClients err:", err)
	//	u.Respond(w, u.Error(-1, err.Error()))
	//} else {
	//
	//}
	u.Respond(w, u.SucessWithObject(openwrt.GetInstance().GetClients()))
}

func (this *Api) ResetClients(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	openwrt.GetInstance().ResetClients()
	res.Ok("重置成功~")
}

func (this *Api) GetStatus(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	mac := r.URL.Query().Get("mac")
	if mac == "" {
		res.Error("mac地址空～")
		return
	}

	tempFilePath := filepath.Join(openwrt.StatusDir, mac)
	list := openwrt.GetInstance().GetDeviceTimeLineDatas(tempFilePath)
	if list == nil {
		res.Ok("暂无数据")
	} else {
		res.Object("获取成功", list)
	}
}

func (this *Api) Clear(w http.ResponseWriter, r *http.Request) {
	//req := utils.GetReqMapData(w, r)
	//glog.Warn(req)
	glog.Warn("Clear---->", r.URL)
	err := u.ClearTemp()
	if err != nil {
		glog.Error("Clear err:", err)
		u.Respond(w, u.Error(-1, err.Error()))
	} else {
		u.OKK(w)
	}
}

func (this *Api) SetStaticIp(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[openwrt.DHCPLease](r)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	}
	err = openwrt.SetStaticIpAddress(body.MAC, body.IP, body.Hostname)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	} else {
		res.Ok("设置成功")
	}
}

func (this *Api) SetWebhook(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[struct {
		WebHookUrl string `json:"webhookUrl"`
	}](r)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	}
	err = openwrt.GetInstance().SetWebHook(body.WebHookUrl)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	} else {
		res.Ok("设置成功")
	}
}
func (this *Api) DeleteStaticIp(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	mac := r.URL.Query().Get("mac")
	if mac == "" {
		res.Error("mac地址空～")
		return
	}
	err := openwrt.GetInstance().DeleteStaticIp(mac)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	} else {
		res.Ok("删除成功")
	}
}

func (this *Api) GetStaticIps(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	ips, err := openwrt.GetInstance().GetStaticIpMap()
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	}
	if ips == nil || len(ips) <= 0 {
		glog.Error("列表空")
		res.Err(fmt.Errorf("列表空"))
		return
	}
	res.Object("请求列表成功", ips)
}

func (this *Api) ResetNetwork(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	err := openwrt.RestartNetwork()
	if err != nil {
		res.Err(err)
		return
	}
	res.Ok("重置成功")
}

func (this *Api) SetNick(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[openwrt.NickEntry](r)
	if err != nil {
		res.Err(err)
		return
	}
	if body.WorkType != nil {
		if body.WorkType.OnWorkTime != "" {
			e := u.TestTimeParse(body.WorkType.OnWorkTime)
			if e != nil {
				res.Err(err)
				return
			}
		}
		if body.WorkType.OffWorkTime != "" {
			e := u.TestTimeParse(body.WorkType.OffWorkTime)
			if e != nil {
				res.Err(err)
				return
			}
		}
	}
	err = openwrt.GetInstance().UpdateNickName(body)
	if err != nil {
		res.Err(err)
		return
	}
	res.Ok("设置成功")
}

func (this *Api) AddWorkTime(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[struct {
		Mac       string `json:"mac"`
		Timestamp int64  `json:"timestamp"`
		IsOnWork  bool   `json:"isOnWork"`
	}](r)
	if err != nil {
		res.Err(err)
		return
	}
	if body.Mac == "" {
		res.Err(fmt.Errorf("mac is empty"))
		return
	}
	if body.Timestamp <= 0 {
		res.Err(fmt.Errorf("timestamp is zero %v", body.Timestamp))
		return
	}
	err = openwrt.AddWorkTime(body.Mac, body.Timestamp, body.IsOnWork)
	if err != nil {
		res.Err(err)
		return
	}
	res.Ok("添加成功")
}

func (this *Api) DelWorkTime(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[struct {
		Mac string `json:"mac"`
		Day string `json:"day"`
	}](r)
	if err != nil {
		res.Err(err)
		return
	}
	if body.Mac == "" {
		res.Err(fmt.Errorf("mac is empty"))
		return
	}
	if body.Day == "" {
		res.Err(fmt.Errorf("Day is empty"))
		return
	}
	err = openwrt.DelWorkTime(body.Mac, body.Day)
	if err != nil {
		res.Err(err)
		return
	}
	res.Ok("删除成功")
}

func (this *Api) GetWorkTime(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[struct {
		Mac string `json:"mac"`
	}](r)
	if err != nil {
		res.Err(err)
		return
	}
	if body.Mac == "" {
		res.Err(fmt.Errorf("mac is empty"))
		return
	}
	data, err := openwrt.GetInstance().GetWorkTime(body.Mac)
	if err != nil {
		res.Err(fmt.Errorf("GetWorkTime err %v", err))
		return
	}
	//for _, work := range data {
	//	glog.Printf("%+v\n", work)
	//}
	res.Object("获取成功", data)
}

func (this *Api) TiggerSignCardEvent(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[struct {
		Mac string `json:"mac"`
	}](r)
	if err != nil {
		res.Err(err)
		return
	}
	if body.Mac == "" {
		res.Err(fmt.Errorf("mac is empty"))
		return
	}
	err = openwrt.GetInstance().NotifySignCardEvent(3, body.Mac, glog.Now())
	if err != nil {
		res.Err(err)
		return
	}
	res.Ok("ok")
}

func (this *Api) UpdatetWorkTime(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[struct {
		Mac  string                 `json:"mac"`
		Day  string                 `json:"day"`
		Data map[string]interface{} `json:"data"`
	}](r)
	if err != nil {
		res.Err(err)
		return
	}
	if body.Mac == "" {
		res.Err(fmt.Errorf("mac is empty"))
		return
	}
	if body.Day == "" {
		res.Err(fmt.Errorf("Day is empty"))
		return
	}
	glog.Printf("修改 %+v\n", body)
	err = openwrt.UpdatetWorkTime(body.Mac, body.Day, body.Data)
	if err != nil {
		res.Err(fmt.Errorf("UpdatetWorkTime err %v", err))
		return
	}
	res.Ok("更新成功")
}

func (this *Api) GetSSE() iface.ISSE {
	return this.sseApi
}
