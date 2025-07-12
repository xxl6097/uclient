package assets

import (
	"embed"
	"github.com/gorilla/mux"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-http/pkg/httpserver"
	"github.com/xxl6097/go-http/pkg/ihttpserver"
	"github.com/xxl6097/go-http/pkg/util"
	"io/fs"
	"net/http"
)

//go:embed static/*
var StaticFS embed.FS
var FileSystem http.FileSystem

func init() {
	subFs, err := fs.Sub(StaticFS, "static")
	if err != nil {
		glog.Fatal("静态资源加载失败", err)
	}
	FileSystem = http.FS(subFs)
}

type StaticRoute struct {
}

func (s StaticRoute) Setup(router *mux.Router) {
	httpserver.RouterUtil.AddNoAuthPrefix("/")
	httpserver.RouterUtil.AddNoAuthPrefix("static")

	router.Handle("/favicon.ico", http.FileServer(FileSystem)).Methods(http.MethodGet, http.MethodOptions)
	router.PathPrefix("/").Handler(util.MakeHTTPGzipHandler(http.StripPrefix("/", http.FileServer(FileSystem)))).Methods(http.MethodGet, http.MethodOptions)
}

func NewRoute() ihttpserver.IRoute {
	opt := &StaticRoute{}
	return opt
}
