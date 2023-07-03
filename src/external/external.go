package external

import (
	"encoding/json"
	"ordinal/src/entity"
	"ordinal/src/httputil"
	"strings"
)

func GetTxIdInfo(txIds []string) []string {

	// 构造请求参数
	requestParam := strings.Join(txIds, ";")
	endPoint := "https://api.blockcypher.com/v1/btc/main/txs/"
	url := endPoint + requestParam

	// 请求api
	res := httputil.GetRequest(url)

	// 解析结果
	var txInfos []entity.TxInfo
	err := json.Unmarshal(res, &txInfos)
	if err != nil {
		// 解析错误
		return nil
	}

	for _, txInfo := range txInfos {
		if len(txInfo.Inputs) != 0 {
			for _, input := range txInfo.Inputs {
				witness := input.Witness
				content := witness[1]
				// 解析content
				//
			}
		}
	}
}
