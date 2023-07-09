package crawler

import (
	"fmt"
	"log"
	"net/http"
	ddb "ordinal/src/db"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	url                = "https://ordinalslite.com/inscriptions"
	reTxId             = "/inscription/(.*?)i0"
	reInscriptionIndex = "/inscriptions/(.*?)"
)

// 执行爬虫任务
func Process() []string {
	var allTxIds []string

	curIndex, curTxId := getLatestIndex()
	if curIndex == 0 || curTxId == "" {
		log.Fatal("fail to get curIndex or curTxId from ddb!")
	}
	fmt.Printf("crawler last time index is %d, txId is %s \n", curIndex, curTxId)

	crawlerIndex, isFirstTime, latestIndex, latestTxId := 0, true, 0, ""
	for {
		txIds := doCrawler(&crawlerIndex, &isFirstTime)

		// 检查上次最后爬到的txId在本次结果中的位置
		i := isContain(curTxId, txIds)
		if i != -1 {
			// 需要截断，当前页面已经包含了上次的结果
			txIds = txIds[0:i]
		}
		fmt.Printf("txids len is %d, all txids is %v \n", len(txIds), txIds)
		allTxIds = append(allTxIds, txIds...)

		// 第一次爬取到的是当前最新的交易索引
		if isFirstTime {
			latestIndex = crawlerIndex + 100
			latestTxId = allTxIds[0]
		}
		fmt.Printf("current crawler latest index is %v, and now isFirst time is %v \n", latestIndex, isFirstTime)

		// 如果当前已经爬取到的索引位置大于下一次爬取的索引，说明当前已经爬到了上一次的位置
		if curIndex >= crawlerIndex {
			break
		}
	}
	fmt.Printf("crawler end, get all txids len is %d, get all txids is %v \n", len(allTxIds), allTxIds)

	saveCurIndexAndTxId(latestIndex, latestTxId)
	return allTxIds
}

// TODO from ddb 保存当前爬到的最新的index和txId
func saveCurIndexAndTxId(curIndex int, curTxId string) {
	fmt.Printf("save latest index is %d and latest txId is %s \n", curIndex, curTxId)
	ddb.UpdateInscriptionID(curIndex, curTxId)
}

// TODO from ddb 上一次获取的inscription的编号和txid
func getLatestIndex() (int, string) {
	inscription := ddb.GetLatestIndex()
	return inscription.InscriptionID, inscription.GenesisTxID
}

func doCrawler(crawlerIndex *int, isFirstTime *bool) []string {
	// 构建爬虫的目标url
	var targetUrl = url
	if *crawlerIndex != 0 {
		// 不是第一次爬取了
		*isFirstTime = false
		// 后缀索引不是0的时候增加目标地址后缀翻页爬取
		targetUrl = targetUrl + "/" + strconv.Itoa(*crawlerIndex)

	}
	fmt.Printf("targetUrl is %s \n", targetUrl)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", targetUrl, nil)
	// 自定义header
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36")
	req.Header.Set("accept", "*/*")
	req.Header.Set("sec-fetch-mode", "cors")
	// do call
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("call website error!", err)
	}
	defer resp.Body.Close()

	// 解析HTML文档
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln("fail to parser ordinalslite html!")
	}

	reIndex := regexp.MustCompile(reInscriptionIndex)
	// 获取当前的页面的下一页的索引，prev索引一定存在
	doc.Find(".center .prev").Each(func(index int, ele *goquery.Selection) {
		prevLink, _ := ele.Attr("href")
		if prevLink == "" || len(prevLink) == 0 {
			log.Fatalln("fail to extract prev index from html!")
		}
		fmt.Printf("get prevLink regex is %s \n", reIndex.FindStringSubmatch(prevLink)[1])
		preIndex, err := strconv.Atoi(reIndex.FindStringSubmatch(prevLink)[1])
		if err != nil {
			log.Fatalln("fail to get prev index!", err)
		}
		// 下一次要爬取的后缀索引
		*crawlerIndex = preIndex
	})

	var originTxIds = []string{}
	reTx := regexp.MustCompile(reTxId)
	// 获取当前所有的txId链接
	doc.Find(".thumbnails").Find("a").Each(func(index int, ele *goquery.Selection) {
		inscriptionLink, _ := ele.Attr("href")
		originTxId := reTx.FindStringSubmatch(inscriptionLink)
		if len(originTxId) > 0 {
			originTxIds = append(originTxIds, originTxId[1])
		}
	})
	return originTxIds
}

func isContain(ele string, src []string) int {
	for i, e := range src {
		if e == ele {
			return i
		}
	}
	return -1
}
