package Spider

import (
	myjieba "SearchEngineByGolang/split_words_by_gojieba"
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"SearchEngineByGolang/Page"
	ST "SearchEngineByGolang/Stop_Token"
	"SearchEngineByGolang/config"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type Spider struct {
	ID       int             // 标记该爬虫的ID
	WorkList chan *Page.Page // 发送给该爬虫的有缓冲的工作
	//TargetUrl string //目标Url
}

var st ST.StopTokens

func init() {

	st.Init(config.SPLIT_TOKEN_CN_EN_PATH)
}

// 创建对象接口
func NewSpider(Id int) *Spider {

	spider := &Spider{}
	spider.ID = Id
	spider.WorkList = make(chan *Page.Page, config.MAX_WORK) // 默认一个Worker最多做10个任务 不是缓冲的将会阻塞

	return spider
}

// 创建对象接口
func NewSpiderCached(Id int, MaxSize int) *Spider {

	spider := &Spider{}
	spider.ID = Id
	spider.WorkList = make(chan *Page.Page, MaxSize) // 默认一个Worker最多做10个任务 不是缓冲的将会阻塞

	return spider
}

// 添加一个任务到队列
func (this *Spider) AddWork(page *Page.Page) {

	if page != nil {
		this.WorkList <- page
	}
}

// 获取单个页面
func (this *Spider) Fetch(url string) string {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	// 解决418问题 添加一个header
	reqest.Header.Add("User-Agent", config.USER_AGENT)

	resp, err := client.Do(reqest)
	if resp == nil {
		fmt.Println("RESP is NIL")
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

	// 得到二进制数据
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// 编码转换
	var sourceCode = string(result)
	if !utf8.Valid(result) {
		data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(result)
		sourceCode = string(data)
	}
	//fmt.Printf("%s", sourceCode)
	return sourceCode
}

// 对页面进行解析 获取到对应的数据 如：中文文章
func (this *Spider) Parse(html, elmID string, pageType bool) string {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}
	var parseResult string

	dom.Find(elmID).Each(func(i int, selection *goquery.Selection) {
		//fmt.Println(selection.Text())

		if pageType {
			parseResult += selection.Text()
		}

		////str = strings.Replace(selection.Text(), "\n", "", -1)
		////parseResult += strings.Trim(selection.Text(), "\n")
		if !pageType {
			check := false
			str := selection.Text()
			if len(selection.Text()) >= 5 {
				str = selection.Text()[len(selection.Text())-5 : len(selection.Text())]
			}
			for _, v := range str {
				// 其实就只用检查一次
				if unicode.Is(unicode.Han, v) {
					check = true
					break
				}
				//break;
			}
			if !check {
				parseResult += selection.Text()
			}

		}

	})
	if pageType {
		parseResult = strings.Replace(parseResult, " ", "", -1)
		parseResult = strings.Replace(parseResult, "\n", "", -1)
	} else {
		parseResult = strings.Replace(parseResult, ".", " ", -1)
		parseResult = strings.Replace(parseResult, "?", " ", -1)
		parseResult = strings.Replace(parseResult, "!", " ", -1)
		parseResult = strings.Replace(parseResult, ",", " ", -1)
		reg := regexp.MustCompile("[\u4e00-\u9fa5]")
		parseResult = reg.ReplaceAllString(parseResult, " ")
	}

	return parseResult
}

// 保存到本地
func (this *Spider) Save(content string, filePath string) error {
	//if (len(content) == 0) {
	//	return errors.New("内容为空")
	//}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Open file err = ", err)
		return err
	}
	file.Write([]byte(content)) // 将字符串转为数组存放
	defer file.Close()
	return nil
}

func RemoveStopToken(words []string) []string {

	newWords := make([]string, 0)
	for _, word := range words {
		if !st.IsStopToken(word) {
			newWords = append(newWords, word)
			//fmt.Println(newWords)
		}
	}

	return newWords
}

// 入口
func (this *Spider) DoWork(page *Page.Page) {
	var err error
	// 1. 获取该页面的Text文档
	html := this.Fetch(page.Url)

	// 2. 使用自己的正则表达式规则匹配 得到具体数据
	content := this.Parse(html, page.ElementId, page.PageType)

	// 3. 将原始数据保存到本地
	config.COUNT_LOCK.Lock()
	if page.PageType { // 中文文档处理

		err = this.Save(content, config.NEWS_OLD_CN+"SWJTU_NEWS_OLD_FILES_"+strconv.Itoa(int(config.COUNT))+".txt")
		// 4. 将数据进行分词处理 并去除停用词
		myjieba.SplitWords(content, config.SPILIT_WORDS_CN_PATH+"SPLITED_WORDS_"+strconv.Itoa(int(config.COUNT))+".txt")
		// 5. 将数据进行去除停用词
		//words = RemoveStopToken(words)
		//this.Save(strings.Join(words, "\n"), config.SPILIT_WORDS_EN_PATH + "SPILIED_WORDS_" + strconv.Itoa(int(config.COUNT)) + ".txt")
		// 6.将数据转换成为可以被搜索引擎模式
		myjieba.SplitWords(content, config.NEWS_NEW_CN+"SearchMode_New_Files_CN_"+strconv.Itoa(int(config.COUNT))+".txt")

	} else {
		err = this.Save(content, config.NEWS_OLD_EN+"EN_READING_OLD_FILES_"+strconv.Itoa(int(config.COUNT))+".txt")
		words := strings.Fields(content)
		before := len(words)
		words = RemoveStopToken(words)
		after := len(words)
		log.Printf("target : [%s] was removed %d stop words!...\n", page.Url, before-after)
		this.Save(strings.Join(words, "\n"), config.SPILIT_WORDS_EN_PATH+"SPILIED_WORDS_"+strconv.Itoa(int(config.COUNT))+".txt")

		// PorterStemming提取词干
		stemWords := make([]string, len(words))
		for idx, word := range words {

			stem := porterstemmer.Stem([]rune(word))

			stemWords[idx] = string(stem)

			//fmt.Printf("The word [%s] has the stem [%s].\n", string(word), string(stem))
		}
		this.Save(strings.Join(stemWords, "\n"), config.NEW_FILE_EN+"SearchMode_New_Files_EN_"+strconv.Itoa(int(config.COUNT))+".txt")

	}
	if err != nil {
		log.Println("在处理【" + page.Url + "】时 写入文件失败")
	}
	config.COUNT++

	config.COUNT_LOCK.Unlock()

	//log.Println(words)
	// 6.保存到
}

// 启动该爬虫机器人
func (this *Spider) Run() {
	cnt := 0
	for url := range this.WorkList {
		//url := <-this.WorkList
		this.DoWork(url)
		// 需要增加Timer用于回收Worker
		log.Printf("Spider[%d] has done the [%s] works now[%s]\n", this.ID, url.Url, time.Now())
		cnt++
	}
	log.Printf("Spider[%d] has done the [%d] works now[%s]\n", this.ID, cnt, time.Now())
}

func (this *Spider) ShutDown() {
	close(this.WorkList)
	log.Printf("Spider[%d] is shutdown, now\n", this.ID)
}
