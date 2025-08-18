package internal

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/github"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/go-service/pkg/utils/util"
	"github.com/xxl6097/uclient/internal/u"
	"github.com/xxl6097/uclient/pkg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func (this *Api) ApiHeap(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	res.Ok(fmt.Sprintf("当前堆内存: %v MB", memStats.HeapAlloc/1024/1024))
	//glog.Println("操作系统:", runtime.GOOS)     // 如 "linux", "windows"
	//glog.Println("CPU 架构:", runtime.GOARCH) // 如 "amd64", "arm64"
	//glog.Println("CPU 核心数:", runtime.NumCPU())
}

func (this *Api) ApiVersion(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	res.Sucess("获取成功", u.GetVersion())
	//glog.Println("操作系统:", runtime.GOOS)     // 如 "linux", "windows"
	//glog.Println("CPU 架构:", runtime.GOARCH) // 如 "amd64", "arm64"
	//glog.Println("CPU 核心数:", runtime.NumCPU())
}

func (this *Api) ApiCheckVersion(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	data, err := github.Api().CheckUpgrade(pkg.BinName)
	if err != nil {
		res.Err(err)
	} else {
		glog.Debug("version:", data)
		res.Any(data)
	}
}

func (this *Api) ApiUpdate(w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	ctx := r.Context()
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	updir := glog.AppHome()
	_, _, free, _ := util.GetDiskUsage(updir)
	if free < u.GetSelfSize()*2 {
		if err := utils.ClearTemp(); err != nil {
			fmt.Println("/tmp清空失败:", err)
		} else {
			fmt.Println("/tmp清空完成")
		}
	}

	var newFilePath string
	switch r.Method {
	case "PUT", "put":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			res.Response(400, fmt.Sprintf("read request body error: %v", err))
			glog.Warnf("%s", res.Msg)
			return
		}
		if len(body) == 0 {
			res.Response(400, "升级URL空的哦～")
			glog.Warnf("%s", res.Msg)
			return
		}
		binUrl := string(body)
		glog.Debugf("upgrade by url: %s", binUrl)
		newUrl := utils.DownloadFileWithCancelByUrls(github.Api().GetProxyUrls(binUrl))
		newFilePath = newUrl
		break
	case "POST", "post":
		// 获取上传的文件
		file, handler, err := r.FormFile("file")
		if err != nil {
			res.Error("body no file")
			return
		}
		defer file.Close()
		dstFilePath := filepath.Join(glog.AppHome("temp", "upgrade"), handler.Filename)
		//dstFilePath 名称为上传文件的原始名称
		dst, err := os.Create(dstFilePath)
		if err != nil {
			res.Error(fmt.Sprintf("create file %s error: %v", handler.Filename, err))
			return
		}
		buf := this.pool.Get().([]byte)
		defer this.pool.Put(buf)
		_, err = io.CopyBuffer(dst, file, buf)
		dst.Close()
		if err != nil {
			res.Error(err.Error())
			return
		}
		newFilePath = dstFilePath
		break
	default:
		res.Error("位置请求方法")
	}
	if newFilePath != "" {
		glog.Debugf("开始升级 %s", newFilePath)
		err := this.igs.Upgrade(ctx, newFilePath)
		glog.Debug("---->升级", err)
		if err == nil {
			res.Ok("升级成功～")
		} else {
			res.Error(fmt.Sprintf("更新失败～%v", err))
		}

	}
}
