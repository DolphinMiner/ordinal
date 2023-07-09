package com.crypto.punk.common;

public enum ResultCode implements IErrorCode {
    SUCCESS(200, "success"),
    FAILED(500, "failed"),
    VALIDATE_FAILED(404, "invalid param");
    private Integer code;
    private String message;

    ResultCode(Integer code, String message) {
        this.code = code;
        this.message = message;
    }

    public Integer getCode() {
        return code;
    }

    public String getMessage() {
        return message;
    }
}
