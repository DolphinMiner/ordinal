package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// 执行爬虫任务
func Process() {
	url := "https://ordinalslite.com/inscriptions" // 要爬取的网页 URL

	// 发送 HTTP GET 请求
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败：", err)
		return
	}
	defer response.Body.Close()

	// 读取响应的内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败：", err)
		return
	}

	// 正则解析

	// 打印网页内容
	fmt.Println(string(body))

	// 处理结果返回交易信息
}
