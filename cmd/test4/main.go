package main

import (
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/uclient/internal/openwrt"
	"github.com/xxl6097/uclient/internal/u"
	"math/rand"
	"path/filepath"
	"sort"
	"time"
)

func tee() {
	timeLines := make([]openwrt.DeviceTimeLine, 0)
	t := glog.Now()
	for i := 0; i < 10; i++ {
		t1 := time.Date(t.Year(), t.Month(), t.Day(), rand.Intn(11), rand.Intn(60), 0, 0, time.Local)
		timeLines = append(timeLines, openwrt.DeviceTimeLine{
			DateTime:  t1.Format(time.DateTime),
			Timestamp: t1.UnixMilli(),
			Connected: i%2 == 0,
			Ago:       t.Sub(t1).String(),
		})
		time.Sleep(time.Millisecond * 100)
	}
	for _, line := range timeLines {
		fmt.Println(line)
	}

	fmt.Println("排序.....")
	sort.Slice(timeLines, func(i, j int) bool {
		a := u.UTC8ToTime(timeLines[i].Timestamp)
		b := u.UTC8ToTime(timeLines[j].Timestamp)
		timeLines[i].Ago = t.Sub(a).String()
		timeLines[j].Ago = t.Sub(b).String()
		return timeLines[i].Timestamp < timeLines[j].Timestamp
	})

	for _, line := range timeLines {
		fmt.Println(line)
	}
}

func tee1() {
	mac := "5a:a7:22:62:3d:26"
	tempFilePath := filepath.Join("/Users/uuxia/Downloads/192.168.1.1/202507181215", mac)
	timeLines := openwrt.GetInstance().GetDeviceTimeLineDatas(tempFilePath)
	for _, line := range timeLines {
		fmt.Println(line)
	}
}

func tee2() {
	t := glog.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), rand.Intn(11), rand.Intn(60), 0, 0, time.Local)
	sub := t.Sub(t1)
	fmt.Println(sub.String())
	du := time.Duration(sub.Seconds()) * time.Second
	fmt.Println(du.String())
}
func main() {
	//tee()
	//tee1()
	tee2()
}
