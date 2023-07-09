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
	"sort"
	"strconv"
	"strings"
)

const OrdinalTableName = "ordinal"

var ctx = context.Background()

func GetLatestIndex() *entity.Ordinal {
	result, err := ordinalClient().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(OrdinalTableName),
		Key: map[string]types.AttributeValue{
			"ordinalID":  &types.AttributeValueMemberN{Value: strconv.Itoa(-1)},
			"createTime": &types.AttributeValueMemberS{Value: "0"},
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

func QueryOrdinal(tokenID int) ([]entity.Ordinal, error) {

	result, err := ordinalClient().Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(OrdinalTableName),
		KeyConditionExpression: aws.String("tokenID = :tokenID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			"tokenID": &types.AttributeValueMemberN{Value: strconv.Itoa(tokenID)},
		},
	})

	if err != nil {
		return nil, err
	}

	var ordinals []entity.Ordinal

	for _, item := range result.Items {
		var ordinal entity.Ordinal
		err := attributevalue.UnmarshalMap(item, &ordinal)
		if err != nil {
			log.Fatalln("fail to UnmarshalMap ordinal record")
		}
		ordinals = append(ordinals, ordinal)
	}

	// 根据 createTime 进行排序
	sort.Slice(ordinals, func(i, j int) bool {
		return strings.Compare(ordinals[i].CreateTime, ordinals[j].CreateTime) < 0
	})

	return ordinals, nil
}

//func QueryFullOrdinal() ([]entity.Ordinal, error) {
//
//}

func UpdateInscriptionID(newInscriptionID int, genesisTxID string) error {

	update := &dynamodb.UpdateItemInput{
		TableName: aws.String(OrdinalTableName),
		Key: map[string]types.AttributeValue{
			"tokenID":    &types.AttributeValueMemberN{Value: strconv.Itoa(-1)},
			"createTime": &types.AttributeValueMemberN{Value: strconv.Itoa(0)},
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

func PutOrdinal(ordinal *entity.Ordinal) error {

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

func ordinalClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return dynamodb.NewFromConfig(cfg)
}
