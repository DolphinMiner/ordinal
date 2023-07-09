package com.crypto.punk.ddb;


import com.crypto.punk.common.Constants;
import com.crypto.punk.entity.Ordinal;
import org.springframework.stereotype.Repository;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.model.AttributeValue;
import software.amazon.awssdk.services.dynamodb.model.QueryRequest;
import software.amazon.awssdk.services.dynamodb.model.QueryResponse;

import javax.annotation.Resource;
import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Repository
public class OrdinalMapper {
    @Resource
    private DynamoDbClient dynamoDbClient;
    public List<Ordinal> queryByTokenID(Integer tokenID) {

        QueryRequest queryRequest = QueryRequest.builder()
                .tableName(Constants.ORDINAL_TABLE_NAME)
                .keyConditionExpression("#tokenID = :tokenID")
                .expressionAttributeValues(
                        Map.of(":tokenID", AttributeValue.builder().n(tokenID.toString()).build()))
                .build();
        QueryResponse queryResponse = dynamoDbClient.query(queryRequest);

        return queryResponse.items().stream()
                .map(item -> Ordinal.builder()
                        .tokenID(Integer.parseInt(item.get("tokenID").n()))
                        .index(Integer.parseInt(item.get("index").n()))
                        .createTime(item.get("createTime").s())
                        .genesisTxID(item.get("genesisTxID").s())
                        .inscriptionID(Integer.parseInt(item.get("inscriptionID").n()))
                        .build())
                .sorted(Comparator.comparing(Ordinal::getIndex))
                .collect(Collectors.toList());
    }
}
