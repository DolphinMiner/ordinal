package com.crypto.punk.api;

import com.crypto.punk.common.CommonResult;
import com.crypto.punk.ddb.OrdinalMapper;
import com.crypto.punk.entity.Ordinal;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.Resource;
import java.util.List;

@RequestMapping("v1/crypto/punks")
@RestController
public class OrdinalApi {

    @Resource
    private OrdinalMapper ordinalMapper;

    @GetMapping("/queryByID")
    public CommonResult<List<Ordinal>> queryByID(@RequestParam Integer tokenID) {
        List<Ordinal> ordinals = ordinalMapper.queryByTokenID(tokenID);
        return CommonResult.success(ordinals);
    }


    @GetMapping("/queryRecentMint")
    public CommonResult<List<Ordinal>> queryRecentMint() {
        List<Ordinal> ordinals = ordinalMapper.queryRecentMint();
        return CommonResult.success(ordinals);
    }

}
