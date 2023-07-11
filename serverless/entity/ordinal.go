package entity

type Ordinal struct {
	TokenID       int    `dynamodbav:"tokenID" json:"tokenID"`
	Index         int    `dynamodbav:"index" json:"index"`
	GenesisTxID   string `dynamodbav:"genesisTxID" json:"genesisTxID"`
	InscriptionID int    `dynamodbav:"inscriptionID" json:"inscriptionID"`
	CreateTime    string `dynamodbav:"createTime" json:"createTime"`
}
