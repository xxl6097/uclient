package openwrt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"io"
	"net/http"
)

func WebHookMessage(webhookUrl string, payload any) error {
	//jsonData, err := json.Marshal(payload)
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	glog.Debug("webhook", string(jsonData))
	resp, err := http.Post(
		webhookUrl,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		glog.Errorf("Error: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("无法读取响应内容:", err)
		return err
	}
	glog.Println("响应内容:", resp.StatusCode, string(respBody))
	return err
}
