package com.crypto.punk.ddb;


import com.crypto.punk.common.Constants;
import com.crypto.punk.entity.Ordinal;
import org.apache.log4j.Logger;
import org.springframework.stereotype.Component;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.model.*;

import javax.annotation.Resource;
import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Component
public class OrdinalMapper {

    private static final Logger log = Logger.getLogger(OrdinalMapper.class);
    @Resource
    private DynamoDbClient dynamoDbClient;
    public List<Ordinal> queryByTokenID(Integer tokenID) {

        QueryRequest queryRequest = QueryRequest.builder()
                .tableName(Constants.ORDINAL_TABLE_NAME)
                .keyConditionExpression("#tokenID = :tokenID")
                .expressionAttributeValues(
                        Map.of(":tokenID", AttributeValue.builder().n(tokenID.toString()).build()))
                .build();
        try{
            QueryResponse queryResponse = dynamoDbClient.query(queryRequest);

            return queryResponse.items().stream()
                    .map(OrdinalMapper::convertToOrdinals)
                    .sorted(Comparator.comparing(Ordinal::getIndex))
                    .collect(Collectors.toList());
        } catch (Exception e) {
            log.error("【QueryError】" + e);
        }
        return null;
    }

    public List<Ordinal> queryAllOrdinals() {
        ScanRequest scanRequest = ScanRequest.builder()
                .tableName(Constants.ORDINAL_TABLE_NAME)
                .filterExpression("index = :index")
                .expressionAttributeValues(
                        Map.of(":index", AttributeValue.builder().n("0").build()))
                .build();
        try {
            ScanResponse scanResponse = dynamoDbClient.scan(scanRequest);
            return scanResponse.items().stream()
                    .map(OrdinalMapper::convertToOrdinals)
                    .collect(Collectors.toList());
        } catch (Exception e) {
            log.error("【ScanError】" + e);
        }
        return null;
    }

    private static Ordinal convertToOrdinals(Map<String, AttributeValue> attrMap) {
        return Ordinal.builder()
                        .tokenID(Integer.parseInt(attrMap.get("tokenID").n()))
                        .index(Integer.parseInt(attrMap.get("index").n()))
                        .createTime(attrMap.get("createTime").s())
                        .genesisTxID(attrMap.get("genesisTxID").s())
                        .inscriptionID(Integer.parseInt(attrMap.get("inscriptionID").n()))
                        .build();
    }
}
