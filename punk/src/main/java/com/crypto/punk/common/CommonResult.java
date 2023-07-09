package com.crypto.punk.common;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class CommonResult<T> {

    private Integer code;
    private String message;
    private T data;

    public static <T> CommonResult<T> success(T data){
        return (CommonResult<T>) CommonResult.builder()
                .code(ResultCode.SUCCESS.getCode())
                .message( ResultCode.SUCCESS.getMessage())
                .data(data)
                .build();
    }

    /**
     * 失败返回结果
     * @param message 提示信息
     */
    public static <T> CommonResult<T> failed(String message) {
        return (CommonResult<T>) CommonResult.builder()
                .code(ResultCode.FAILED.getCode())
                .message(message)
                .build();
    }

    /**
     * 失败返回结果
     */
    public static <T> CommonResult<T> failed() {
        return failed(ResultCode.FAILED);
    }

    /**
     * 失败返回结果
     * @param errorCode 错误码
     */
    public static <T> CommonResult<T> failed(IErrorCode errorCode) {
        return (CommonResult<T>) CommonResult.builder()
                .code(errorCode.getCode())
                .message(errorCode.getMessage())
                .build();
    }
}
