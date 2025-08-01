package main

import (
	"github.com/xxl6097/uclient/internal/ntfy"
	"github.com/xxl6097/uclient/internal/u"
)

func main() {
	//resp, err := http.Get("http://uuxia.cn:90/uclient/json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()
	//scanner := bufio.NewScanner(resp.Body)
	//for scanner.Scan() {
	//	println(scanner.Text())
	//}

	//req, e := http.NewRequest("GET", "http://uuxia.cn:90/uclient/json", nil)
	//if e != nil {
	//	panic(e)
	//}
	//req.SetBasicAuth("admin", "het002402")
	////req.Header.Set("Authorization", "Basic tk_trk4agho2")
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(resp.StatusCode)
	//scanner := bufio.NewScanner(resp.Body)
	//for scanner.Scan() {
	//	fmt.Println(scanner.Text())
	//}
	//defer req.Body.Close()

	//go ntfy.Subscribe("http://uuxia.cn:90", "uclient", "admin", "het002402", func(s string) {
	//	fmt.Println(s)
	//})
	//ntfy.Subscribe("http://uuxia.cn:90", "uclient", "admin", "het002402", func(s string) {
	//	fmt.Println("--->", s)
	//})

	ntfy.GetInstance().Start(&u.NtfyInfo{
		Address:  "http://uuxia.cn:90",
		Topic:    "uclient",
		Username: "admin",
		Password: "het002402",
	})

}
