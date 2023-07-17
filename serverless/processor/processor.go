package processor

import (
	"encoding/json"
	ddb "ordinal/serverless/db"
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
							saveSuccess, err := saveInfo(imageId, txInfo)
							if err != nil {
								return false
							}
							return saveSuccess
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
	// 将所有图片信息转换为一个map,<imageHex,imageId> 初始化图片信息
	images := make(map[string]int, 10000)
	value, exists := images[imageContent]
	return value, exists
}

func saveInfo(imageId int, txInfo entity.TxInfo) (bool, error) {
	// 将数据存入DDB
	ordinal := entity.Ordinal{
		TokenID:     imageId,
		SequenceNo:  0,
		GenesisTxID: txInfo.Hash,
		// todo 爬虫爬inscription id
		InscriptionID: 0,
		CreateTime:    txInfo.Confirmed.Format("2006-01-02 15:04:05"),
	}

	err := ddb.PutOrdinal(&ordinal)
	if err != nil {
		return false, err
	}
	return true, err
}
