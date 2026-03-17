package main

import (
	"github.com/xxl6097/glog/pkg/z"
	assets "github.com/xxl6097/uclient/assets/openwrt"
	"io/fs"
	"net/http"
)

func main() {
	subFs, _ := fs.Sub(assets.StaticFS, "static")
	http.Handle("/", http.FileServer(http.FS(subFs)))
	z.Debug("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
