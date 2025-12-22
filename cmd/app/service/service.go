package service

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	"github.com/xxl6097/glog/glog"
	_ "github.com/xxl6097/go-service/assets/buffer"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/uclient/internal"
	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/u"
	"github.com/xxl6097/uclient/pkg"
)

type Service struct {
	timestamp string
	gs        igs.Service
}

func (this *Service) OnFinish() {
	//openwrt.GetInstance().Close()
}
func (this *Service) OnStop() {
	openwrt.GetInstance().Close()
}

func (this *Service) OnShutdown() {
}

func load() (*u.Config, error) {
	defer glog.Flush()
	byteArray, err := ukey.Load()
	if err != nil {
		return nil, err
	}
	var cfg u.Config
	err = ukey.GobToStruct(byteArray, &cfg)
	//err = json.Unmarshal(byteArray, &cfg)
	if err != nil {
		glog.Println("ClientConfig解析错误", err)
		return nil, err
	}
	pkg.Version()
	return &cfg, nil
}

func (this *Service) OnConfig() *service.Config {
	cfg := service.Config{
		Name:        pkg.AppName,
		DisplayName: pkg.DisplayName,
		Description: pkg.Description,
	}
	return &cfg
}

func (this *Service) OnVersion() string {
	pkg.Version()
	//cfg, err := load()
	//if err == nil {
	//	glog.Debugf("cfg:%+v", cfg)
	//}
	return pkg.AppVersion
}

func (this *Service) OnRun(service igs.Service) error {
	this.gs = service
	cfg, err := load()
	if err != nil {
		return err
	}
	glog.Debug("程序运行", os.Args)
	//server := httpserver.New().
	//	CORSMethodMiddleware().
	//	AddRoute(internal.NewRoute(internal.NewApi(service, cfg.Username, cfg.Password))).
	//	AddRoute(assets.NewRoute()).
	//	BasicAuth(cfg.Username, cfg.Password, "oIin3168TLKg1X8OU2xBBWLlMEdI").
	//	Done(cfg.ServerPort)
	//defer server.Stop()
	//server.Wait()
	internal.Bootstrap(cfg, service)
	return nil
}

func (this *Service) GetAny(s2 string) []byte {
	return this.menu()
}

func (this *Service) menu() []byte {
	port := utils.InputIntDefault(fmt.Sprintf("输入服务端口(%d)：", 7000), 7000)
	username := utils.InputStringEmpty(fmt.Sprintf("输入管理用户名(%s)：", "admin"), "admin")
	password := utils.InputStringEmpty(fmt.Sprintf("输入管理密码(%s)：", "admin"), "admin")
	cfg := &u.Config{ServerPort: port, Username: username, Password: password}
	bb, e := ukey.StructToGob(cfg)
	if e != nil {
		return nil
	}
	return bb
}
