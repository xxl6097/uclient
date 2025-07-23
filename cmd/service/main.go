package main

import (
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/gs"
	"github.com/xxl6097/uclient/cmd/service/service"
)

func main() {
	//defer glog.GlobalRecover()
	s := service.Service{}
	err := gs.Run(&s)
	glog.Debug("程序结束", err)
}
