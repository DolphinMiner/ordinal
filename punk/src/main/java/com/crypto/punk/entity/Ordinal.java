package com.crypto.punk.entity;

import lombok.Builder;
import lombok.Value;

@Builder
@Value
public class Ordinal {
    Integer tokenID;
    Integer index;
    String createTime;
    String genesisTxID;
    Integer inscriptionID;
}
