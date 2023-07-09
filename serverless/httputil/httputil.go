package httputil

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRequest(url string) []byte {
	// 发送 HTTP GET 请求
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败：", err)
		return nil
	}
	defer response.Body.Close()

	// 读取响应的内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败：", err)
		return nil
	}
	return body
}
