package main

import (
	"fmt"
	"github.com/xxl6097/uclient/internal/u"
	"time"
)

func main() {
	//markdown := make(map[string]interface{})
	//markdown["title"] = "张三上线了"
	//format := "#### %s \n - 名称：%s\n - IP地址：%s \n- Mac地址：%s \n- 时间：%s \n"
	//markdown["text"] = fmt.Sprintf(format, "张三上线了", "张三", "192.168.1.2", "AC:CC:BB:11:22:33", u.TimestampFormat(time.Now().UnixMilli()))
	//payload := map[string]interface{}{"msgtype": "markdown", "markdown": markdown}
	//webHookUrl := "https://oapi.dingtalk.com/robot/send?access_token=122512eee3d8e359643b4b38961c4a729319f3f518e4faa0168fc803abde66bf"
	//_ = openwrt.WebHookMessage(webHookUrl, payload)

	t1 := time.Now()
	on := u.GetTime("18:30:00", u.GetLocation())
	onWorkTime := time.Date(t1.Year(), t1.Month(), t1.Day(), on.Hour(), on.Minute(), on.Second(), 0, t1.Location())
	a := time.Date(t1.Year(), t1.Month(), t1.Day(), 18, 31, on.Second(), 0, t1.Location())
	fmt.Println(onWorkTime.Format(time.DateTime))
	fmt.Println(a.Format(time.DateTime))
	fmt.Println(a.Compare(onWorkTime))
}
