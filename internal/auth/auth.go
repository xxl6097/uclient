package auth

import (
	"github.com/xxl6097/go-service/pkg/ukey"
	"os"
)

var (
	authFilePath = "/etc/config/uclient/auth"
)

func GetAuthData() ([]string, error) {
	data, err := os.ReadFile(authFilePath)
	if err != nil {
		return nil, err
	}
	dataArray := make([]string, 0)
	err = ukey.GobToStruct(data, &dataArray)
	if err != nil {
		return nil, err
	}
	return dataArray, nil
}

func SetAuthData(cfg []string) error {
	if cfg == nil || len(cfg) == 0 {
		return nil
	}
	content, err := ukey.StructToGob(cfg)
	if err != nil {
		return nil
	}
	file, err := os.Create(authFilePath) // 文件不存在则创建，存在则截断
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入内容
	_, err = file.Write(content)
	return err
}

func AddAuthData(code string) error {
	if code == "" {
		return nil
	}
	codes, _ := GetAuthData()
	if codes == nil {
		codes = make([]string, 0)
	}
	codes = append(codes, code)
	return SetAuthData(codes)
}
