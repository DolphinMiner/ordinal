package entity

type Ordinal struct {
	TokenID       int    `dynamodbav:"tokenID" json:"tokenID"`
	SequenceNo    int    `dynamodbav:"sequenceNo" json:"sequenceNo"`
	GenesisTxID   string `dynamodbav:"genesisTxID" json:"genesisTxID"`
	InscriptionID int    `dynamodbav:"inscriptionID" json:"inscriptionID"`
	CreateTime    string `dynamodbav:"createTime" json:"createTime"`
}
