package com.crypto.punk.api;

import com.crypto.punk.common.CommonResult;
import com.crypto.punk.ddb.OrdinalMapper;
import com.crypto.punk.entity.Ordinal;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.Resource;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

@RequestMapping("api/crypto/punks")
@RestController
public class OrdinalApi {

    @Resource
    private OrdinalMapper ordinalMapper;

    @GetMapping("/queryByID")
    public CommonResult<List<Ordinal>> queryByID(@RequestParam Integer tokenID) {
        List<Ordinal> ordinals = ordinalMapper.queryByTokenID(tokenID);
        if(ordinals == null) {
            CommonResult.failed("Query failed!");
        }
        return CommonResult.success(ordinals);
    }

    @GetMapping("/querySortedByRandom")
    public CommonResult<List<Ordinal>> querySortedByRandom() {
        List<Ordinal> ordinals = ordinalMapper.queryAllOrdinals();
        if(ordinals == null) {
            CommonResult.failed("Query failed!");
        }
        return CommonResult.success(ordinals);
    }

    @GetMapping("/querySortedByrRecentMint")
    public CommonResult<List<Ordinal>> querySortedByrRecentMint() {
        List<Ordinal> ordinals = ordinalMapper.queryAllOrdinals().stream()
                .sorted(Comparator.comparing(Ordinal::getCreateTime).reversed())
                .collect(Collectors.toList());
        if(ordinals == null) {
            CommonResult.failed("Query failed!");
        }
        return CommonResult.success(ordinals);
    }

    @GetMapping("/querySortedByPunkID")
    public CommonResult<List<Ordinal>> querySortedByPunkID() {
        List<Ordinal> ordinals = ordinalMapper.queryAllOrdinals().stream()
                .sorted(Comparator.comparing(Ordinal::getTokenID))
                .collect(Collectors.toList());
        if(ordinals == null) {
            CommonResult.failed("Query failed!");
        }
        return CommonResult.success(ordinals);
    }
}
