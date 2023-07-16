package ddb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"ordinal/serverless/entity"
	"strconv"
	"time"
)

const OrdinalTableName = "ordinal"

var ctx = context.Background()

func GetLatestIndex() *entity.Ordinal {
	result, err := ordinalClient().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(OrdinalTableName),
		Key: map[string]types.AttributeValue{
			"ordinalID": &types.AttributeValueMemberN{Value: strconv.Itoa(-1)},
			"index":     &types.AttributeValueMemberN{Value: "0"},
		},
	})
	if err != nil {
		log.Fatalln("fail to get cur inscription index!")
	}

	ordinal := entity.Ordinal{}
	err = attributevalue.UnmarshalMap(result.Item, &ordinal)
	if err != nil {
		log.Fatalln("fail to unmarshalMap")
	}
	return &ordinal
}

func UpdateInscriptionID(newInscriptionID int, genesisTxID string) error {

	update := &dynamodb.UpdateItemInput{
		TableName: aws.String(OrdinalTableName),
		Key: map[string]types.AttributeValue{
			"tokenID":    &types.AttributeValueMemberN{Value: strconv.Itoa(-1)},
			"sequenceNo": &types.AttributeValueMemberN{Value: strconv.Itoa(-1)},
		},
		UpdateExpression: aws.String("SET #inscriptionID = :inscriptionID and #genesisTxID = :genesisTxID "),
		ExpressionAttributeNames: map[string]string{
			"#inscriptionID": "inscriptionID",
			"#genesisTxID":   "genesisTxID",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newInscriptionID": &types.AttributeValueMemberN{Value: strconv.Itoa(newInscriptionID)},
			":genesisTxID":      &types.AttributeValueMemberS{Value: genesisTxID},
		},
	}

	_, err := ordinalClient().UpdateItem(ctx, update)
	return err
}

func PutOrdinal(ordinalRequest *entity.Ordinal) error {
	// 查询ordinal数量作为index
	ordinalsSize, err := queryOrdinalSize(ordinalRequest.TokenID)
	if err != nil {
		log.Fatalln("fail to UnmarshalMap ordinal record")
		return err
	}
	ordinal := &entity.Ordinal{
		TokenID:       ordinalRequest.TokenID,
		SequenceNo:    ordinalsSize,
		GenesisTxID:   ordinalRequest.GenesisTxID,
		InscriptionID: ordinalRequest.InscriptionID,
		CreateTime:    time.Now().Format("2006-01-02 15:04:05"),
	}
	item, err := attributevalue.MarshalMap(ordinal)
	if err != nil {
		return err
	}

	_, err = ordinalClient().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(OrdinalTableName),
		Item:      item,
	})
	return err
}

func queryOrdinalSize(tokenID int) (int, error) {

	result, err := ordinalClient().Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(OrdinalTableName),
		KeyConditionExpression: aws.String("tokenID = :tokenID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			"tokenID": &types.AttributeValueMemberN{Value: strconv.Itoa(tokenID)},
		},
	})

	if err != nil {
		return 0, err
	}
	return len(result.Items), nil
}

func ordinalClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return dynamodb.NewFromConfig(cfg)
}
