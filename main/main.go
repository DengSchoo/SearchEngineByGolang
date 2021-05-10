package main

import (
	"SearchEngineByGolang/Page"
	"SearchEngineByGolang/Spider"
	"SearchEngineByGolang/config"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Close(spiders []*Spider.Spider) {
	for _, spider := range spiders {
		spider.ShutDown()
	}
}

func BeginGet500CN() bool {
	spiders := make([]*Spider.Spider, config.MAX_WORKER)
	for i := 0; i < config.MAX_WORKER; i++ {
		spiders[i] = Spider.NewSpiderCached(i+1, config.MAX_WORK)
		//spiders = append(spiders, Spider.NewSpiderCached(i + 1, config.MAX_WORK))
		go spiders[i].Run()
	}
	//go spider.Run();
	for j := 0; j < config.DOC_NUM/config.MAX_WORKER; j++ {
		for i := 0; i < config.MAX_WORKER; i++ {
			page := Page.Page{
				config.PREFIX_URL + strconv.Itoa(config.Base_URL+j*config.MAX_WORKER+i+1) + config.SUBFIX_URL,
				config.CLASS_ID,
				true,
			}
			spiders[i].AddWork(&page)
		}
	}
	time.Sleep(20 * time.Second)
	Close(spiders)
	return true
}

func testConnection(url string) string {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	// 解决418问题 添加一个header
	reqest.Header.Add("User-Agent", config.USER_AGENT)

	resp, err := client.Do(reqest)
	if resp == nil {
		return "nil"
	}
	if err != nil {
		fmt.Println("Error Status is :", resp.StatusCode)
		return "nil"
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error Status is :", resp.StatusCode)
		return "nil"
	}
	time.Sleep(1000)
	return "ok"
}

func getPageList() []string {
	spider := new(Spider.Spider)
	pages := make([]string, 500)
	baseUrl := "http://www.enread.com/essays/list_"
	num := 25
	end := ".html"
	slector := "body > div.wrap > div.main > table > tbody > tr > td.left > div > div > div.list > div > div.title > h2 > a"
	for i := 1; i <= num; i++ {
		content := spider.Fetch(baseUrl + strconv.Itoa(i) + end)
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			log.Fatalln(err)
		}
		var parseResult string
		cnt := 0
		dom.Find(slector).Each(func(idx int, selection *goquery.Selection) {
			//fmt.Println(selection.Text())
			// 去除空格
			//parseResult += strings.Replace(selection.Text(), " ", "", -1)
			//// 去除换行符
			//parseResult += strings.Replace(selection.Text(), "\n", "", -1)
			////str = strings.Replace(selection.Text(), "\n", "", -1)
			////parseResult += strings.Trim(selection.Text(), "\n")
			temp, _ := selection.Attr("href")
			pages[(i-1)*20+cnt] = "http://www.enread.com" + temp
			cnt = cnt + 1
			print("http://www.enread.com" + temp + "\n")
		})
		print(parseResult)
	}
	return pages
}

func getEnPages() []string {
	//spider := new(Spider.Spider)

	pages := make([]string, 500)
	cnt := 0
	base := config.Base_URL_EN
	//dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	for {
		//content := spider.Fetch(config.NEWS_LIST_PREFIX + strconv.Itoa(base) +  config.NEWS_LIST_SUBFIX)

		if testConnection(config.PREFIX_URL_EN+strconv.Itoa(base)+config.SUBFIX_URL_EN) == "ok" {
			pages[cnt] = config.PREFIX_URL_EN + strconv.Itoa(base) + config.SUBFIX_URL_EN
			cnt++
			//dom, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
			//dom.Find(config.NEWS_LSIT_ID).Each(func(i int, selection *goquery.Selection) {
			//	a, _ := selection.Attr("href")
			//	pages[cnt] = a
			//	cnt++
			log.Println(pages[cnt-1])
			if cnt == 500 {
				return pages
			}
		}

		if base == 0 {
			return pages
		}
		base--
		time.Sleep(1 * time.Second)
		//time.Sleep(1000)

	}
	return pages
}
func BeginGet500EN(pages []string) bool {
	spiders := make([]*Spider.Spider, config.MAX_WORKER)
	for i := 0; i < config.MAX_WORKER; i++ {
		spiders[i] = Spider.NewSpiderCached(i+1, config.MAX_WORK)
		//spiders = append(spiders, Spider.NewSpiderCached(i + 1, config.MAX_WORK))
		go spiders[i].Run()
	}

	//go spider.Run();
	for j := 0; j < config.DOC_NUM_EN/config.MAX_WORKER; j++ {
		for i := 0; i < config.MAX_WORKER; i++ {
			page := Page.Page{
				pages[j*config.MAX_WORKER+i],
				config.CLASS_ID_EN,
				false,
			}
			time.Sleep(500 * time.Millisecond)
			spiders[i].AddWork(&page)
		}
	}

	Close(spiders)
	return true
}

func main() {
	// 1. 开始爬取 并处理中文页面

	//start := time.Now()
	//BeginGet500CN()
	//log.Println("500 中文文档耗时 ", time.Since(start))
	////简单同步 time.Sleep(10 * time.Second)
	//time.Sleep(10 * time.Second)

	//// 2. 开始爬取 并处理英文页面
	start := time.Now()
	config.COUNT = 1
	pages_url := getPageList()
	BeginGet500EN(pages_url)
	log.Println("全过程耗时 ", time.Since(start))

	//parseResult := "邓世豪 123 等12 12等22"
	//reg := regexp.MustCompile("[\u4e00-\u9fa5]")
	//parseResult = reg.ReplaceAllString(parseResult, "")
	//log.Println(parseResult)
	time.Sleep(200 * time.Second)
	//test_reg()
}
