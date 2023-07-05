package crawler

import (
	"fmt"
	"log"
	"net/http"
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
	// curIndex := getLatestIndex()
	// TODO 需要循环爬，把所有的结果丢到一个切片里
	// doCrawler(curIndex, )
}

// 上一次获取的inscription的编号
func getLatestIndex() int {
	// 当前页 5433601 ~ 5433502
	// prev 5433501
	// next 5433701

	// 跳到前一页应该是 5433501 ~ 5433402
	// prev 5433401
	// next 5433601
	return 5433601
}

func doCrawler(beginIndex int, isLatest *bool) (int, []string) {
	client := &http.Client{}

	var targetUrl = url + "/" + strconv.Itoa(beginIndex)
	fmt.Printf("targetUrl is %s \n", targetUrl)

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

	// 当前查询到的最新的一次铭刻的index
	var currentIndex int
	reIndex := regexp.MustCompile(reInscriptionIndex)
	// 获取当前的页面的下一页的索引，prev索引一定存在
	doc.Find(".center .prev").Each(func(index int, ele *goquery.Selection) {
		prevLink, _ := ele.Attr("href")
		if prevLink == "" || len(prevLink) == 0 {
			log.Fatalln("fail to extract prev index from html!")
		}
		preIndex, err := strconv.Atoi(reIndex.FindStringSubmatch(prevLink)[1])
		if err != nil {
			log.Fatalln("fail to get prev index!")
		}
		currentIndex = preIndex + 100
	})

	// 获取当前的页面的下一页的索引，next索引不存在的时候说明当前已经是最新的铭刻交易
	doc.Find(".center .next").Each(func(index int, ele *goquery.Selection) {
		*isLatest = true
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
	return currentIndex, originTxIds
}
