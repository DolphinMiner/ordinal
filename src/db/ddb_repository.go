package ddb

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tbName = "cur_index"

func GetCurIndex() *CurIndex {
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	svc := dynamodb.New(sess)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(tbName),
			},
		},
		TableName: aws.String("Music"),
	}

	response, err := svc.GetItem(input)
	if err != nil {
		log.Fatalln("fail to get cur inscription index!")
	}
	curIndex := CurIndex{}
	err = dynamodbattribute.UnmarshalMap(response.Item, &curIndex)
	if err != nil {
		log.Fatalln("fail to unmarshaMap!")
	}
	return &curIndex
}
