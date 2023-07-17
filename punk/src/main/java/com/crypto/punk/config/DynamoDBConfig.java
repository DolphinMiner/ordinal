package com.crypto.punk.config;

import com.crypto.punk.common.Constants;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import software.amazon.awssdk.auth.credentials.DefaultCredentialsProvider;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;

@Configuration
public class DynamoDBConfig {

    @Bean(name = "dynamoDbClient")
    public DynamoDbClient buildDynamoDbClient() {
        return DynamoDbClient.builder()
                .region(Constants.REGION)
                .credentialsProvider(DefaultCredentialsProvider.create())
                .build();
    }

}
