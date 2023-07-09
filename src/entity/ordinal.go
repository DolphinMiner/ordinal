package entity

type Ordinal struct {
	TokenID       int    `dynamodbav:"tokenID" json:"tokenID"`
	CreateTime    string `dynamodbav:"createTime" json:"createTime"`
	GenesisTxID   int    `dynamodbav:"genesisTxID" json:"genesisTxID"`
	InscriptionID int    `dynamodbav:"inscriptionID" json:"inscriptionID"`
}
