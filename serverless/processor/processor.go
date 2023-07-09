package processor

import (
	"encoding/json"
	"ordinal/serverless/entity"
	"ordinal/serverless/httputil"
	"strings"
	"time"
)

func GetTxIdInfo(txIds []string) bool {

	endPoint := "https://api.blockcypher.com/v1/ltc/main/txs/"

	// 分组处理(每组数量不超过3个)
	groups := make([][]string, 0)
	currentGroup := make([]string, 0)

	for _, str := range txIds {
		currentGroup = append(currentGroup, str)

		if len(currentGroup) == 3 {
			groups = append(groups, currentGroup)
			currentGroup = make([]string, 0)
		}
	}
	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	// 处理每一组,每组处理完后延迟1s
	for _, txGroup := range groups {

		requestParam := strings.Join(txGroup, ";")
		url := endPoint + requestParam
		// 请求api
		res := httputil.GetRequest(url)
		// 解析结果
		var txInfos []entity.TxInfo
		err := json.Unmarshal(res, &txInfos)
		if err != nil {
			// 解析错误,返回处理失败
			return false
		}
		// 解析交易内容
		for _, txInfo := range txInfos {
			if len(txInfo.Inputs) != 0 {
				for _, input := range txInfo.Inputs {
					witness := input.Witness
					content := witness[1]
					// 检索png字节数组开头的这一段
					// index := strings.Index(content, "89504E470D0A1A0A")
					index := strings.Index(content, "8950ac")
					if index != -1 {
						imageContent := content[index : len(content)-2]
						// 判断图片16进制数据是否是项目集合中的16进制数据
						if imageId, exits := judgeImage(imageContent); exits {
							// 将txId信息存入DDB
							saveInfo(imageContent, imageId, txInfo.Hash)
						}
					}
				}
			}
		}
		// 延时1s
		time.Sleep(1 * time.Second)
	}
	return true
}

func judgeImage(imageContent string) (int, bool) {
	// 将所有图片信息转换为一个map,<imageHex,imageId>
	images := make(map[string]int, 10000)
	value, exists := images[imageContent]
	return value, exists
}

func saveInfo(imageContent string, imageId int, txHash string) {
	// 将数据存入DDB
}
