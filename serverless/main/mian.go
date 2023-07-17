package main

import (
	"ordinal/serverless/crawler"
	"ordinal/serverless/entity"
	"ordinal/serverless/processor"
)

func main() {
	// lambda.Start(Handler)

	Handler(entity.InputEvent{})
}

func Handler(event entity.InputEvent) entity.OutputEvent {
	// 执行爬虫任务
	txIds := crawler.Process()

	// mock crawler
	// txIds := []string{
	// 	"9772cfaeadd15d0c37323ca027ecc88095eadee476b6ceaf58a9b5dd23ca3968",
	// 	"9050f4c0d40ea2fa39f3431cf034864317db36c9d3950e8d31c9246bfc68ca99",
	// 	"cbab6155ab74f57efb3632cc530e833c55f9d83e4b71ce7a34826aac602abccf",
	// 	"305735922eb0bb14c52f56a90d394c56ce31fae8a74569e2bc74456bc3f0f2b7",
	// 	"b7805d8a4acd03e3dbd1f06f13c689248d3b5d39c55e4d992451ad7db10566a1",
	// }

	// 处理爬虫结果
	result := processor.GetTxIdInfo(txIds)

	// 返回执行结果
	return buildRes(result)
}

func buildRes(boolValue bool) entity.OutputEvent {
	res := entity.OutputEvent{}

	if boolValue {
		res.Code = 0
		res.Message = "success"
	} else {
		res.Code = 1
		res.Message = "failed"
	}
	return res
}
