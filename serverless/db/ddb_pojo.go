package ddb

type CurIndex struct {
	Key      string `dynamodbav:"key" json:"key"`
	CurIndex int    `dynamodbav:"cur_index" json:"curIndex"`
}
