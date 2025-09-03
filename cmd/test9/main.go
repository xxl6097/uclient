package main

import (
	"fmt"
	"time"
)

// CountWeekendsInMonth 统计指定年份和月份的周末天数（周六和周日）
func CountWeekendsInMonth(year int, month time.Month) int {
	// 获取指定月份的第一天（00:00:00）
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	// 获取下个月的第一天，然后减去一秒，得到当前月份的最后一天
	nextMonth := firstDay.AddDate(0, 1, 0)
	lastDay := nextMonth.Add(-time.Second)

	weekendCount := 0
	// 从月份的第一天开始遍历，直到处理完该月的所有天
	for current := firstDay; !current.After(lastDay); current = current.AddDate(0, 0, 1) {
		// 获取当前日期是星期几
		weekday := current.Weekday()
		// 如果当前是周六或周日，则增加周末计数
		//if weekday == time.Saturday || weekday == time.Sunday {
		if weekday == time.Saturday {
			weekendCount++
		}
	}

	return weekendCount
}
func main() {
	js := "{\"SiteConfig\":{\"SiteName\":\"MoonTV\",\"Announcement\":\"\xe6\x9c\xac\xe7\xbd\x91\xe7\xab\x99\xe4\xbb\x85\xe6\x8f\x90\xe4\xbe\x9b\xe5\xbd\xb1\xe8\xa7\x86\xe4\xbf\xa1\xe6\x81\xaf\xe6\x90\x9c\xe7\xb4\xa2\xe6\x9c\x8d\xe5\x8a\xa1\xef\xbc\x8c\xe6\x89\x80\xe6\x9c\x89\xe5\x86\x85\xe5\xae\xb9\xe5\x9d\x87\xe6\x9d\xa5\xe8\x87\xaa\xe7\xac\xac\xe4\xb8\x89\xe6\x96\xb9\xe7\xbd\x91\xe7\xab\x99\xe3\x80\x82\xe6\x9c\xac\xe7\xab\x99\xe4\xb8\x8d\xe5\xad\x98\xe5\x82\xa8\xe4\xbb\xbb\xe4\xbd\x95\xe8\xa7\x86\xe9\xa2\x91\xe8\xb5\x84\xe6\xba\x90\xef\xbc\x8c\xe4\xb8\x8d\xe5\xaf\xb9\xe4\xbb\xbb\xe4\xbd\x95\xe5\x86\x85\xe5\xae\xb9\xe7\x9a\x84\xe5\x87\x86\xe7\xa1\xae\xe6\x80\xa7\xe3\x80\x81\xe5\x90\x88\xe6\xb3\x95\xe6\x80\xa7\xe3\x80\x81\xe5\xae\x8c\xe6\x95\xb4\xe6\x80\xa7\xe8\xb4\x9f\xe8\xb4\xa3\xe3\x80\x82\",\"SearchDownstreamMaxPage\":5,\"SiteInterfaceCacheTime\":7200,\"SearchResultDefaultAggregate\":true},\"UserConfig\":{\"AllowRegister\":true,\"Users\":[{\"username\":\"admin\",\"role\":\"owner\"},{\"username\":\"uuxia\",\"role\":\"user\"},{\"username\":\"cp\",\"role\":\"user\"},{\"username\":\"ch\",\"role\":\"user\"}]},\"SourceConfig\":[{\"key\":\"dyttzy\",\"name\":\"\xe7\x94\xb5\xe5\xbd\xb1\xe5\xa4\xa9\xe5\xa0\x82\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"http://caiji.dyttzyapi.com/api.php/provide/vod\",\"detail\":\"http://caiji.dyttzyapi.com\",\"from\":\"config\",\"disabled\":false},{\"key\":\"heimuer\",\"name\":\"\xe9\xbb\x91\xe6\x9c\xa8\xe8\x80\xb3\",\"api\":\"https://json.heimuer.xyz/api.php/provide/vod\",\"detail\":\"https://heimuer.tv\",\"from\":\"config\",\"disabled\":false},{\"key\":\"ruyi\",\"name\":\"\xe5\xa6\x82\xe6\x84\x8f\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://cj.rycjapi.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"bfzy\",\"name\":\"\xe6\x9a\xb4\xe9\xa3\x8e\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://bfzyapi.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"tyyszy\",\"name\":\"\xe5\xa4\xa9\xe6\xb6\xaf\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://tyyszy.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"ffzy\",\"name\":\"\xe9\x9d\x9e\xe5\x87\xa1\xe5\xbd\xb1\xe8\xa7\x86\",\"api\":\"http://ffzy5.tv/api.php/provide/vod\",\"detail\":\"http://ffzy5.tv\",\"from\":\"config\",\"disabled\":false},{\"key\":\"zy360\",\"name\":\"360\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://360zy.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"iqiyi\",\"name\":\"iqiyi\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://www.iqiyizyapi.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"wolong\",\"name\":\"\xe5\x8d\xa7\xe9\xbe\x99\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://wolongzyw.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"jisu\",\"name\":\"\xe6\x9e\x81\xe9\x80\x9f\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://jszyapi.com/api.php/provide/vod\",\"detail\":\"https://jszyapi.com\",\"from\":\"config\",\"disabled\":false},{\"key\":\"dbzy\",\"name\":\"\xe8\xb1\x86\xe7\x93\xa3\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://dbzy.tv/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"mozhua\",\"name\":\"\xe9\xad\x94\xe7\x88\xaa\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://mozhuazy.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"mdzy\",\"name\":\"\xe9\xad\x94\xe9\x83\xbd\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://www.mdzyapi.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"zuid\",\"name\":\"\xe6\x9c\x80\xe5\xa4\xa7\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://api.zuidapi.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"yinghua\",\"name\":\"\xe6\xa8\xb1\xe8\x8a\xb1\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://m3u8.apiyhzy.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"wujin\",\"name\":\"\xe6\x97\xa0\xe5\xb0\xbd\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://api.wujinapi.me/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"wwzy\",\"name\":\"\xe6\x97\xba\xe6\x97\xba\xe7\x9f\xad\xe5\x89\xa7\",\"api\":\"https://wwzy.tv/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false},{\"key\":\"ikun\",\"name\":\"iKun\xe8\xb5\x84\xe6\xba\x90\",\"api\":\"https://ikunzyapi.com/api.php/provide/vod\",\"from\":\"config\",\"disabled\":false}]}"
	fmt.Println(js)

	year := 2024
	month := time.Month(9)
	fmt.Println(year, month)
	count := CountWeekendsInMonth(year, month)
	fmt.Printf("%d年%d月有%d个周末天\n", year, month, count) // 输出：2024年1月有8个周末天
}
