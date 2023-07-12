package com.crypto.punk.ddb;


import com.crypto.punk.common.Constants;
import com.crypto.punk.entity.Ordinal;
import org.springframework.stereotype.Repository;
import org.springframework.util.StringUtils;
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials;
import software.amazon.awssdk.auth.credentials.StaticCredentialsProvider;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.model.AttributeValue;
import software.amazon.awssdk.services.dynamodb.model.QueryRequest;
import software.amazon.awssdk.services.dynamodb.model.QueryResponse;

import javax.annotation.Resource;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Repository
public class OrdinalMapper {

    @Resource
    private DynamoDbClient dynamoDbClient;

    public List<Ordinal> queryByTokenID(Integer tokenID) {

        AwsBasicCredentials credentials = AwsBasicCredentials.create("s", "x");
        StaticCredentialsProvider credentialsProvider = StaticCredentialsProvider.create(credentials);
        DynamoDbClient client = DynamoDbClient.builder()
                .region(Region.US_EAST_1)
                .credentialsProvider(credentialsProvider)
                .build();

        QueryRequest queryRequest = QueryRequest.builder()
                .tableName(Constants.ORDINAL_TABLE_NAME)
                .keyConditionExpression("#tokenID = :tokenID")
                .expressionAttributeValues(
                        Map.of(":tokenID", AttributeValue.builder().n(tokenID.toString()).build()))
                .build();
        QueryResponse queryResponse = client.query(queryRequest);

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

    public List<Ordinal> queryRecentMint() {

        AwsBasicCredentials credentials = AwsBasicCredentials.create("x", "x");
        StaticCredentialsProvider credentialsProvider = StaticCredentialsProvider.create(credentials);
        DynamoDbClient client = DynamoDbClient.builder()
                .region(Region.US_EAST_1)
                .credentialsProvider(credentialsProvider)
                .build();

        // 定义查询请求
        String tableName = Constants.ORDINAL_TABLE_NAME;
        String partitionKey = "tokenID";
        String sortKey = "index";

        Map<String, AttributeValue> expressionAttributeValues = new HashMap<>();

        expressionAttributeValues.put(":indexValue", AttributeValue.builder().n("0").build());
        expressionAttributeValues.put(":tokenID", AttributeValue.builder().n("0").build());

        Map<String, String> expressionAttributeNames = new HashMap<>();
        expressionAttributeNames.put("#tokenID", partitionKey);
        expressionAttributeNames.put("#index", sortKey);
        expressionAttributeNames.put("#genesisTxID", "genesisTxID");

        QueryRequest queryRequest = QueryRequest.builder()
                .tableName(tableName)
                .keyConditionExpression("#index = :indexValue and #tokenID = :tokenID")
                .projectionExpression("#tokenID, #index, #genesisTxID")
                .expressionAttributeValues(expressionAttributeValues)
                .expressionAttributeNames(expressionAttributeNames)
                .build();

        // 执行查询请求
        QueryResponse queryResponse = client.query(queryRequest);


        // 处理结果，返回mint的nft按时间先后
        String pattern = "yyyy-MM-dd HH:mm:ss SSSSSZ";
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern(pattern);

        return queryResponse.items().stream()
                .map(item -> Ordinal.builder()
                        .tokenID(Integer.parseInt(item.get("tokenID").n()))
                        .index(Integer.parseInt(item.get("index").n()))
                        .createTime(item.get("createTime").s())
                        .genesisTxID(item.get("genesisTxID").s())
                        .inscriptionID(Integer.parseInt(item.get("inscriptionID").n()))
                        .build())
                .filter(ordinal -> !StringUtils.isEmpty(ordinal.getCreateTime()))
                .sorted(Comparator.comparing(ordinal -> LocalDateTime.parse(ordinal.getCreateTime(), formatter)))
                .collect(Collectors.toList());
    }

}
