package ddb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"ordinal/src/entity"
	"strconv"
)

const OrdinalTableName = "ordinal"

func GetCurInscription(ctx context.Context) (int, error) {
	result, err := ordinalClient(ctx).GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(OrdinalTableName),
		Key: map[string]types.AttributeValue{
			"ordinalID":  &types.AttributeValueMemberN{Value: strconv.Itoa(-1)},
			"createTime": &types.AttributeValueMemberS{Value: "0"},
		},
	})
	if err != nil {
		return 0, err
	}

	ordinal := entity.Ordinal{}
	err = attributevalue.UnmarshalMap(result.Item, &ordinal)
	if err != nil {
		return 0, err
	}
	return ordinal.InscriptionID, nil
}

func QueryOrdinal(ctx context.Context, tokenID int) ([]entity.Ordinal, error) {

	result, err := ordinalClient(ctx).Query(ctx, &dynamodb.QueryInput{
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

	return ordinals, nil
}

func UpdateInscriptionID(ctx context.Context, newInscriptionID int) error {

	update := &dynamodb.UpdateItemInput{
		TableName: aws.String(OrdinalTableName),
		Key: map[string]types.AttributeValue{
			"tokenID":    &types.AttributeValueMemberN{Value: strconv.Itoa(-1)},
			"createTime": &types.AttributeValueMemberN{Value: strconv.Itoa(0)},
		},
		UpdateExpression: aws.String("SET #inscriptionID = :inscriptionID"),
		ExpressionAttributeNames: map[string]string{
			"#inscriptionID": "inscriptionID",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newInscriptionID": &types.AttributeValueMemberN{Value: strconv.Itoa(newInscriptionID)},
		},
	}

	_, err := ordinalClient(ctx).UpdateItem(ctx, update)
	return err
}

func PutOrdinal(ctx context.Context, ordinal *entity.Ordinal) error {

	item, err := attributevalue.MarshalMap(ordinal)
	if err != nil {
		return err
	}

	_, err = ordinalClient(ctx).PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(OrdinalTableName),
		Item:      item,
	})
	return err
}

func ordinalClient(ctx context.Context) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return dynamodb.NewFromConfig(cfg)
}
