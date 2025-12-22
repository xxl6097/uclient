package internal

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-http/pkg/util"
	"github.com/xxl6097/go-service/pkg/github"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-sse/pkg/sse"
	"github.com/xxl6097/go-sse/pkg/sse/isse"
	"github.com/xxl6097/uclient/internal/auth"
	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/u"
)

type Api struct {
	igs    igs.Service
	sseApi isse.ISseServer
	pool   *sync.Pool // use sync.Pool caching buf to reduce gc ratio
}

func NewApi(igs igs.Service, username, password string) *Api {
	github.Api().SetName("xxl6097", "uclient")
	//initSSEClient(username, password)
	a := &Api{
		igs:    igs,
		sseApi: initSSE(),
		pool: &sync.Pool{
			New: func() interface{} { return make([]byte, 32*1024) },
		},
	}
	openwrt.GetInstance().SetFunc(func(dataType int, obj any) {
		if obj != nil && a.sseApi != nil {
			eve := isse.Event{
				Payload: obj,
			}
			//eve := isse.SSEEvent{
			//	Payload: obj,
			//}
			switch dataType {
			case 0:
				eve.Event = "updateAll"
				break
			case 1:
				eve.Event = "updateOne"
				break
			case 2:
				eve.Event = "showNotify"
				break

			}
			a.sseApi.Broadcast(eve)
		}
	})
	return a
}

func initSSE() isse.ISseServer {
	//serv := sse.New().
	//	InvalidateFun(func(request *http.Request) (string, error) {
	//		return time.Now().Format("20060102150405.999999999"), nil
	//	}).
	//	Register(func(server iface.ISseServer, client *iface.Client) {
	//		//server.Stream("内置丰富的开发模板，包括前后端开发所需的所有工具，如pycharm、idea、navicat、vscode以及XTerminal远程桌面管理工具等模板，用户可以轻松部署和管理各种应用程序和工具", time.Millisecond*500)
	//	}).
	//	UnRegister(nil).
	//	Done()
	return sse.
		New().
		Register(func(server isse.ISseServer, client *isse.Client) {
			//glog.Debug("sse新链接", client)
			//openwrt.GetInstance().StartStatus()
		}).
		UnRegister(func(server isse.ISseServer, client *isse.Client) {
			cls := server.GetClients()
			if cls != nil {
				//glog.Debug("sse链接断开", len(cls), client)
				if len(cls) == 0 {
					//表示没有客户端了
					//openwrt.GetInstance().StopStatus()
				}
			}
		}).
		Done()
}

func initSSEClient(username string, password string) *sse.Client {
	//serv := sse.New().
	//	InvalidateFun(func(request *http.Request) (string, error) {
	//		return time.Now().Format("20060102150405.999999999"), nil
	//	}).
	//	Register(func(server iface.ISseServer, client *iface.Client) {
	//		//server.Stream("内置丰富的开发模板，包括前后端开发所需的所有工具，如pycharm、idea、navicat、vscode以及XTerminal远程桌面管理工具等模板，用户可以轻松部署和管理各种应用程序和工具", time.Millisecond*500)
	//	}).
	//	UnRegister(nil).
	//	Done()
	url := "http://uuxia.cn:7001/api/sse"
	return sse.NewClient(url).
		BasicAuth(username, password).
		ListenFunc(func(s string) {
			glog.Debugf("SSE: %s", s)
		}).Header(func(header *http.Header) {
		header.Add("Sse-Event-IP-Address", util.GetHostIp())
		header.Add("Sse-Event-MAC-Address", u.GetLocalMac())
	}).Done()
}

//func (this *Api) listen(list []*openwrt.DHCPLease) {
//	if len(list) >= 0 && this.sseApi != nil {
//		eve := iface.SSEEvent{
//			Event:   "update",
//			Payload: list,
//		}
//		this.sseApi.Broadcast(eve)
//	}
//}
//
//func (this *Api) eventFunc(dataType int, obj any) {
//	if obj != nil && this.sseApi != nil {
//		eve := iface.SSEEvent{
//			Event:   "update-one",
//			Payload: obj,
//		}
//		if dataType == 2 {
//			eve.Event = "update-status"
//		}
//
//		this.sseApi.Broadcast(eve)
//	}
//}

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

func (this *Api) Reboot(w http.ResponseWriter, r *http.Request) {
	//req := utils.GetReqMapData(w, r)
	//glog.Warn(req)
	glog.Warn("Reboot---->", r.URL)
	err := this.igs.Restart()
	if err != nil {
		glog.Error("Reboot err:", err)
		u.Respond(w, u.Error(-1, err.Error()))
	} else {
		u.OKK(w)
	}
}

func (this *Api) AddStaticIp(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[openwrt.DHCPLease](r)
	if err != nil {
		res.Err(err)
		return
	}
	err = openwrt.SetStaticIpAddress(body.MAC, body.IP, body.Hostname)
	if err != nil {
		res.Err(err)
		return
	} else {
		res.Ok("设置成功")
	}
}

func (this *Api) SetNtfy(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[u.NtfyInfo](r)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	}
	err = openwrt.GetInstance().SetNtfy(body)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	} else {
		res.Ok("设置成功")
	}
}

func (this *Api) AddAuthCode(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	body, err := u.GetDataByJson[struct {
		Authcode string `json:"authcode"`
	}](r)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	}
	if body.Authcode == "" {
		glog.Error("Authcode is empty")
		res.Error("Authcode is empty")
		return
	}

	err = auth.AddAuthData(body.Authcode)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	} else {
		res.Ok("设置成功")
		openwrt.GetInstance().LoadAuth()
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

func (this *Api) SetSettings(w http.ResponseWriter, r *http.Request) {
	res, ff := Response(r)
	defer ff(w)
	body, err := u.GetDataByJson[u.Settings](r)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	}
	err = openwrt.GetInstance().SetSettings(body)
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	} else {
		res.Ok("设置成功")
	}
}

func (this *Api) GetSettings(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	settings, err := openwrt.GetInstance().GetSettings()
	if err != nil {
		glog.Error(err)
		res.Err(err)
		return
	}
	res.Sucess("获取成功", settings)
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
		res.Err(err)
		return
	}
	if ips == nil || len(ips) <= 0 {
		res.Err(fmt.Errorf("列表空"))
		return
	}
	res.Object("请求列表成功", ips)
	//res.Any(ips)
}

func (this *Api) GetLedLog(w http.ResponseWriter, r *http.Request) {
	data, err := u.ReadFile(openwrt.LedEventLog)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(400)
		return
	}
	_, _ = w.Write(data)
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
	data, err := openwrt.GetInstance().GetWorkTimeAndCaculate(body.Mac)
	if err != nil {
		res.Err(fmt.Errorf("GetWorkTime err %v", err))
		return
	}
	//for _, work := range data {
	//	glog.Printf("%+v\n", work)
	//}
	res.Object("获取成功", data)
}

func (this *Api) OfflineDevice(w http.ResponseWriter, r *http.Request) {
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
	err = openwrt.OfflineDevice(body.Mac)
	if err != nil {
		res.Err(err)
		return
	}
	res.Ok("ok")
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
	err = openwrt.GetInstance().TiggerSignCardEvent(body.Mac)
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
	err = openwrt.ApiUpdateWorkTime(body.Mac, body.Day, body.Data)
	if err != nil {
		res.Err(fmt.Errorf("UpdatetWorkTime err %v", err))
		return
	}
	res.Ok("更新成功")
}

func (this *Api) GetSSE() isse.ISseServer {
	return this.sseApi
}
