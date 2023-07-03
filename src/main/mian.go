package main

import (
	"ordinal/src/crawler"
	"ordinal/src/entity"
)

func main() {
	// 执行爬虫任务
	crawler.Process()
	// lambda.Start(Handler)
}

func Handler(event entity.InputEvent) {
	// 执行爬虫任务
	crawler.Process()
}
