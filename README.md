# SearchEngineByGolang总览

> 《SWJTU搜索引擎课设项目一》
>
> 文本预处理。
>
> 难点：
>
> - 网页高效爬取
> - 整体架构
> - 并发控制
> - 字符编码
> - 正则表达式字符过滤
> - Porter Stemming词干提取算法
> - 性能测试
>
> 创新点：
>
> - 自设计单机多线程(goroutine)架构
>
> - 停用词map集合设计
> - 编码自适应不同爬虫目标
>
> 成果：
>
> - 一键启动，爬取网页到文本预处理
>
> 可以优化提升的方向
>
> - 参考MapReduce架构，做分布式处理。进一步提升性能
> - 对于任务的分发还是写的比较死，这块也应该结合channel
> - 多goroutine同步还需要解决
> - 代码coverage不高，写代码的姿势还有待提高
> - 命名存在一定的问题
> - 持久化可以采用数据库存储mysql，postgreSQL等
> - 错误处理还是不妥，仍有许多unhandled err



## 目录

- 项目需求
- 项目整体思路
  - 编程语言的选取
  - 如何爬取原始文档
  - 如何对原始文档进行类Trim操作
  - 中文分词技术及开源库
  - 如何描述一个页面任务
  - 如何持久化
  - 如何提高性能
  - 网站目标文档选取
  - 中英文停用词表
- 整体结构
  - 项目架构
  - 一个爬虫任务流程
  - 项目目录结构
- 具体实现及成果
  - 下载引擎 -- Spider引擎的实现
  - 单词字符化，删除特殊字符，进行大小写转换
  - 中文分词技术和工具实现中文分词
  - 删除英文停用词
  - 删除中文停用词
  - Porter Stemming 词干提取算法实现
  - 中文文档字符化 生成搜索引擎模式字符单元
  - 英文预处理结果持久化
  - 中文文档预处理结果持久化
- 性能测试
  - 本机环境
  - 单元测试
  - 速度测试
  - 基准测试
  - 覆盖率测试
  - CPU占用测试
- 心得体会
- 参考资料



## FinalDesign需求

- 通过下载引擎(Web Crawler/Spider)自动下载至少500个英文文档/网页，以及500个中文文档/网页，越多越好，并保留原始的文档/网页备份(如: News_1_0rg. txt)

- 编程对所下载文档进行自动预处理:将各个单词进行字符化，完成删除特殊字符、大小写转换等操作

  - 调研并选择合适的中文分词技术和工具实现中文分词
  - 删除英文停用词(Stop Word)

  - 删除中文停用词调用或者编程实现英文Porter Stemming 功能
    将中文文档进行字符化，即可被搜索引擎索引的字符单元

  - 对于英文文档，经过以上处理之后，将经过处理之后所形成简化文档保存(如:News_ 1._E.txt)，以备以后的索引处理_

  - 对于中文文档，经过以上处理之后，将经过处理之后所形成简化文档保存(如:News_ 1_C.txt)，以备以后的索引处理



## 1. **编程语言的选取 -- Golang**

就搜索引擎来说Google是第一家，本项目选用Google自主开发的Golang语言实现。

***\*Golang语言相比其它编程语言有如下优点：\****

 

Ø 可直接编译成机器码，不依赖其他库，glibc的版本有一定要求，部署就是扔一个文件上去就完成了。

 

Ø 静态类型语言，但是有动态语言的感觉，静态类型的语言就是可以在编译的时候检查出来隐藏的大多数问题，动态语言的感觉就是有很多的包可以使用，写起来的效率很高。 

 

Ø 语言层面支持并发，这个就是Go最大的特色，天生的支持并发。Go就是基因里面支持的并发，可以充分的利用多核，很容易的使用并发。 

 

Ø 内置runtime，支持垃圾回收，这属于动态语言的特性之一吧，虽然目前来说GC(内存垃圾回收机制)不算完美，但是足以应付我们所能遇到的大多数情况，特别是Go1.1之后的GC。

 

Ø 简单易学，Go语言的作者都有C的基因，那么Go自然而然就有了C的基因，那么Go关键字是25个，但是表达能力很强大，几乎支持大多数你在其他语言见过的特性：继承、重载、对象等。 

 

Ø 丰富的标准库，Go目前已经内置了大量的库，特别是网络库非常强大。 

 

Ø 内置强大的工具，Go语言里面内置了很多工具链，最好的应该是gofmt工具，自动化格式化代码，能够让团队review变得如此的简单，代码格式一模一样，想不一样都很困难。 

 

Ø 跨平台编译，如果你写的Go代码不包含cgo，那么就可以做到window系统编译linux的应用，如何做到的呢？Go引用了plan9的代码，这就是不依赖系统的信息。 

 

Ø 内嵌C支持，Go里面也可以直接包含C代码，利用现有的丰富的C库。

## 2. **如何爬取原始文档 -- 自实现爬虫引擎 Spider**

***\*爬虫引擎原理：\****

（一）发送HTTP请求给服务器

（二）服务器验证通过后返回资源数据

（三）客户端拿到原始资源数据后，通过正则表达式进行Filter操作，拿到自己的目标数据

（四）进一步对目标数据做持久化或数据处理

## 3. **如何对原始文档进行类Trim操作 -- 使用go std lib**

对字符串做简单处理 可以通过使用对应编程语言提供的标志库进行处理。

```go
// 去除空格
parseResult += strings.Replace(selection.Text(), " ", "", -1)
// 去除换行符
parseResult += strings.Replace(selection.Text(), "\n", "", -1)
```





 

## 4. **如何选择中文分词器 -- gojieba开源项目整合**

***\*gojieba简介：\****

Ø 支持多种分词方式，包括: 最大概率模式, HMM新词发现模式, 搜索引擎模式, 全模式

Ø 核心算法底层由C++实现，性能高效。

Ø 字典路径可配置，NewJieba(...string), NewExtractor(...string) 可变形参，当参数为空时使用默认词典(推荐方式)

## 5. **Porter Stemming的实现 -- CGO调用C语言实现**

采用https://tartarus.org/martin/PorterStemmer/官方提供的C语言实现，ThreadSafe版本。

## 6. **如何描述一个页面任务 -- 建立Page实体**

```go
// page页面
type Page struct {
       // 目标页面Url
       Url           string

       // 目标内容ID
       ElementId      string

       // true 中文页面 false 英文页面
       PageType       bool
}
```



对于一个页面来说，采用URL以及该页面下的标签ID来获取对应内容。并且采用PageType标记是否是中文页面

## 7. **如何持久化 -- 中间文件存储**

采用操作系统提供的文件系统对数据进行进一步的处理。

采用.txt方式保存

## 8. **如何增加项目的可维护性及可扩展性 -- config包做配置文件**

***\*在config包下的config.go文件中：\****

```go
package config

import "sync"

// 常量
const (
       // 网页相关
       USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.115"
       // style ID
       CLASS_ID = ".content14"

       PREFIX_URL = "https://news.swjtu.edu.cn/shownews-"
       SUBFIX_URL = ".shtml"
       // 基础页数
       Base_URL = 20000

       // 最大爬虫数 单个最大工作数
       MAX_WORK   = 10
       MAX_WORKER = 10

       // 文档数量
       DOC_NUM = 500
       DOC_NUM_EN = 500

       // 英文文档相关
       PREFIX_URL_EN = "http://www.enread.com/novel/"
       SUBFIX_URL_EN = ".html"
       CLASS_ID_EN = "#dede_content > div"
       Base_URL_EN = 112944
       NEW_FILE_EN = "./files/New_Files_EN/"

       // files保存位置

       // CN
       NEWS_OLD_CN = "./files/Old_Files_CN/"
       NEWS_NEW_CN = "./files/New_Files_CN/"

       NEWS_OLD_EN = "./files/Old_Files_EN/"

       // 分词token路径
       SPLIT_TOKEN_CN_EN_PATH = "./files/stop_token_files_EN_CN/baidu_stopwords_en_cn.txt"
       SPLIT_TOKEN_CN_PATH = "./files/stop_token_files_EN_CN/stopwords_cn.txt"
       SPLIT_TOKEN_EN_PATH = "./files/stop_token_files_EN_CN/stopwords_en.txt"
       // 分词存储路径
       SPILIT_WORDS_CN_PATH = "./files/split_words_files_CN/"
       SPILIT_WORDS_EN_PATH = "./files/split_words_files_EN/"

)
var (
       COUNT  uint
       // 用于控制COUNT++
       COUNT_LOCK sync.Mutex
)

func init() {
       COUNT = 1
}
```



 

## 9. **如何提高性能 -- 单机多线程（goroutine）**

Go语言最大的特色就是从语言层面支持并发（Goroutine），Goroutine是Go中最基本的执行单元。事实上每一个Go程序至少有一个Goroutine：主Goroutine。当程序启动时，它会自动创建。

 

为了更好理解Goroutine，现讲一下线程和协程的概念

 

***\*线程（Thread）：\****有时被称为轻量级进程(Lightweight Process，LWP），是程序执行流的最小单元。一个标准的线程由线程ID，当前指令指针(PC），寄存器集合和堆栈组成。另外，线程是进程中的一个实体，是被系统独立调度和分派的基本单位，线程自己不拥有系统资源，只拥有一点儿在运行中必不可少的资源，但它可与同属一个进程的其它线程共享进程所拥有的全部资源。

 

线程拥有自己独立的栈和共享的堆，共享堆，不共享栈，线程的切换一般也由操作系统调度。

 

***\*协程（coroutine）：\****又称微线程与子例程（或者称为函数）一样，协程（coroutine）也是一种程序组件。相对子例程而言，协程更为一般和灵活，但在实践中使用没有子例程那样广泛。

 

和线程类似，共享堆，不共享栈，协程的切换一般由程序员在代码中显式控制。它避免了上下文切换的额外耗费，兼顾了多线程的优点，简化了高并发程序的复杂。

 

Goroutine和其他语言的协程（coroutine）在使用方式上类似，但从字面意义上来看不同（一个是Goroutine，一个是coroutine），再就是协程是一种协作任务控制机制，在最简单的意义上，协程不是并发的，而Goroutine支持并发的。因此Goroutine可以理解为一种Go语言的协程。同时它可以运行在一个或多个线程上。

## 10. **网站文档选取**

### **10.1中文文档 -- 西南交通大学新闻网**

选择了URL相比比较好拼的西南交通大学新闻网，并且学校的网站没有反爬虫机制，比较容易。

https://news.swjtu.edu.cn/shownews-20000.shtml

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps1.jpg) 

### **10.2 英文文档 -- 英文阅读网**

爬取的是英文阅读网，小说部分的内容。因为此部分的URL不好拼。就需要先将小说列表url爬取下来，再去爬取对应的URL。

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps6.jpg) 

http://www.enread.com/novel/list_1.html

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps7.jpg) 

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps4.jpg) 

 

## 11. **中英文停用词 -- 词汇表**

github资源项目：

https://github.com/DengSchoo/stopwords

将中文和英语停用词放在同一个文件中，方便后续处理。

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps8.jpg)

# 项目架构设计

## 1. **项目架构**

整个项目从爬取页面到预处理文本数据***\*一键启动\****。如果需要改动目标页面及目标数据则需要在config包下修改对应的配置数据。

 ![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps16.jpg)

## 2. **一个爬虫任务流程**

每个Spider爬虫(goroutine)在执行每一个页面任务及后续数据处理的逻辑如下：

不同的数据持久化到不同的目录下。

 ![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps15.jpg)

## 3. **项目目录结构**

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps14.jpg)

 

| ***\*Directory\****      | ***\*Dir Description\****         |
| ------------------------ | --------------------------------- |
| config                   | 项目中所有配置设置                |
| files                    | 保存所以持久化项目                |
| main                     | 项目启动入口                      |
| Page                     | 页面实体                          |
| porter_stemming_with_CGO | porter_stemmingC语言实现，CGO调用 |
| Spider                   | 自实现搜索下载引擎                |
| split_words_by_gojieba   | 整合gojieba中文分词器             |
| Stop_Token               | 用于提供分词utils                 |
| test                     | 性能测试，及基础测试              |

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps13.jpg)

| ***\*files Directory\**** | ***\*files Dir Description\**** |
| ------------------------- | ------------------------------- |
| New_Files_CN              | 中文文档被预处理的结果          |
| New_Files_EN              | 英文文档被预处理的结果          |
| Old_Files_CN              | 中文原文档                      |
| Old_Files_EN              | 英文原文档                      |
| split_words_files_CN      | 中文文档分词结果持久化          |
| split_words_files_EN      | 英文文档分词结果持久化          |
| stop_token_files_EN_CN    | 中文、英文停用词词汇表          |

# 具体实现

## 1. **搜索引擎 -- Spider引擎的实现**

### **1.1设计思路**

每一个Spider代表一个goroutine，在config中可以配置自己想要的goroutine数量。在本项目中使用的是10个Spider。

这十个Spider会将500个目标页面任务领取完毕，即每个Spider会有50个爬取任务。

首先要做的是请求页面，得到页面初始数据，并编码转换成为统一的字符编码。得到页面基本数据后，需要对页面数据进行基本处理，进行格式整理，如对于中文页面需要去除空格以及换行等。对于英文页面需要将其中隐含的中文字符去除等。

得到经过过滤的文件后，下一步就需要分词处理，分词完毕之后，再对分词结果统一做分词过滤处理。对于中文分词结果可以直接持久化最终文档，对于英文文档来说可以在对分词结果做进一步的处理，即porter stemming算法进行一个词干提取。

### 1.2 运行效果

显示爬取的任务页面URL，以及完成该任务的Spider编号：

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps28.jpg)

在执行完毕Spider自动关闭，并打印显示已完成任务数：

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps29.jpg)

 

### **1.3 源代码**

spider.go文件下：

 ```go
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
       ID       int // 标记该爬虫的ID
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
       if (resp == nil) {
              fmt.Println("RESP is NIL",)
              return  "nil"
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
                            str = selection.Text()[len(selection.Text()) - 5: len(selection.Text())]
                     }
                     for _, v:= range(str) {
                            // 其实就只用检查一次
                            if unicode.Is(unicode.Han, v) {
                                   check = true
                                   break;
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
              reg := regexp.MustCompile(`[\u4e00-\u9fa5]`)
              parseResult = reg.ReplaceAllString(parseResult, "")
       }

       return parseResult
}

// 保存到本地
func (this *Spider) Save(content string, filePath string) (error){
       //if (len(content) == 0) {
       //     return errors.New("内容为空")
       //}
       file, err := os.OpenFile(filePath, os.O_CREATE | os.O_WRONLY, 0666)
       if err != nil {
              fmt.Println("Open file err = ", err)
              return  err;
       }
       file.Write([]byte(content)) // 将字符串转为数组存放
       defer file.Close()
       return nil;
}

func RemoveStopToken(words []string) []string{

       newWords := make([]string, 0)
       for _,word :=  range words {
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

              err = this.Save(content, config.NEWS_OLD_CN + "SWJTU_NEWS_OLD_FILES_" + strconv.Itoa(int(config.COUNT)) + ".txt")
              // 4. 将数据进行分词处理 并去除停用词
              myjieba.SplitWords(content, config.SPILIT_WORDS_CN_PATH + "SPLITED_WORDS_" + strconv.Itoa(int(config.COUNT)) + ".txt")
              // 5. 将数据进行去除停用词
              //words = RemoveStopToken(words)
              //this.Save(strings.Join(words, "\n"), config.SPILIT_WORDS_EN_PATH + "SPILIED_WORDS_" + strconv.Itoa(int(config.COUNT)) + ".txt")
              // 6.将数据转换成为可以被搜索引擎模式
              myjieba.SplitWords(content, config.NEWS_NEW_CN + "SearchMode_New_Files_CN_" + strconv.Itoa(int(config.COUNT)) + ".txt")

       } else {
              err = this.Save(content, config.NEWS_OLD_EN + "EN_READING_OLD_FILES_" + strconv.Itoa(int(config.COUNT)) + ".txt")
              words := strings.Fields(content)
              log.Printf("before remove %d,",len(words))
              words = RemoveStopToken(words)
              log.Printf(" after remove %d,",len(words))
              this.Save(strings.Join(words, "\n"), config.SPILIT_WORDS_EN_PATH + "SPILIED_WORDS_" + strconv.Itoa(int(config.COUNT)) + ".txt")

              // PorterStemming提取词干
              stemWords := make([]string, len(words))
              for idx, word := range words {

                     stem := porterstemmer.Stem([]rune(word))

                     stemWords[idx] = string(stem)

                     fmt.Printf("The word [%s] has the stem [%s].\n", string(word), string(stem))
              }
              this.Save(strings.Join(stemWords, "\n"), config.NEW_FILE_EN + "SPILIED_WORDS_" + strconv.Itoa(int(config.COUNT)) + ".txt")

       }
       if err != nil {
              log.Println("在处理【" +page.Url + "】时 写入文件失败")
       }
       config.COUNT++

       config.COUNT_LOCK.Unlock()

       //log.Println(words)
       // 6.保存到
}

// 启动该爬虫机器人
func (this *Spider) Run() {
       cnt := 0
       for url:= range this.WorkList {
              //url := <-this.WorkList
              this.DoWork(url)
              // 需要增加Timer用于回收Worker
              log.Printf("Spider[%d] has done the [%s] works now[%s]\n", this.ID, url, time.Now())

       }
       log.Printf("Spider[%d] has done the [%d] works now[%s]\n", this.ID, cnt, time.Now())
}

func (this *Spider) ShutDown() {
       close(this.WorkList)
       log.Printf("Spider[%d] is shutdown, now\n", this.ID)
}
 ```



## 2. **将单词字符化，删除特殊字符，进行大小写转换**

### **2.1设计思路**

简单使用正则表达式结合标准库提供的字符串替换函数就能实现。

### **2.2 运行效果**

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps30.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps31.jpg)

 

 

### **2.3 源代码**

```go
if pageType {
       parseResult = strings.Replace(parseResult, " ", "", -1)
       parseResult = strings.Replace(parseResult, "\n", "", -1)
} else {
       parseResult = strings.Replace(parseResult, ".", " ", -1)
       parseResult = strings.Replace(parseResult, "?", " ", -1)
       parseResult = strings.Replace(parseResult, "!", " ", -1)
       parseResult = strings.Replace(parseResult, ",", " ", -1)
       reg := regexp.MustCompile(`[\u4e00-\u9fa5]`)
       parseResult = reg.ReplaceAllString(parseResult, "")
}
```



 

## 3. **中文分词技术和工具实现中文分词**

### **3.1设计思路**

使用中文分词开源框架gojieba实现中文分词。

### **3.2 运行效果**

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps32.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps33.jpg)

### **3.3 源代码**

```go
package split_words_by_gojieba

import (
       ST "SearchEngineByGolang/Stop_Token"
       "SearchEngineByGolang/config"
       "fmt"
       "github.com/yanyiwu/gojieba"
       "log"
       "os"
       "strings"
)

var st ST.StopTokens
func init() {
       st.Init(config.SPLIT_TOKEN_CN_EN_PATH)
}
func RemoveStopToken(words []string) []string{

       newWords := make([]string, 0)
       for _,word :=  range words {
              if !st.IsStopToken(word) {
                     newWords = append(newWords, word)
                     //fmt.Println(newWords)
              }
       }

       return newWords
}

func saveSpilitResult(content , path string) error{
       file, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0666)
       if err != nil {
              fmt.Println("Open file err = ", err)
              return  err;
       }
       file.Write([]byte(content)) // 将字符串转为数组存放
       defer file.Close()
       return nil;
}
var jieba = gojieba.NewJieba()
func SplitWords(content , path string) {
       var words []string

       words = jieba.Cut(content, true)
       //log.Printf("before remove %d,",len(words))
       before := len(words)
       words = RemoveStopToken(words)
       after := len(words)
       //log.Printf(" after remove %d,",len(words))
       log.Printf("target : [%s] was removed %d stop words!...\n",path, before - after)
       saveSpilitResult(strings.Join(words, "\n"), path)
       //return words;
}
func SplitWordsSearchMode(content , path string) {
       var words []string
       jieba := gojieba.NewJieba()
       words = jieba.CutAll(content)

       words = RemoveStopToken(words)
       saveSpilitResult(strings.Join(words, "\n"), path)
       words = nil
       //return words;
}
```





 

## 4. **删除英文停用词**

### **4.1设计思路**

给定停用词表，读取每一行停用词，使用go 的map数据结构构建一个set集合。具体代码为：map[string]bool。

传入分词结果，将分词中的每一个词语都与之匹配，如果存在就舍弃加入到新分词表。如果不存在则加入到其中。最后返回新的词表。最后log一下移除了多少个词语。

### **4.2 运行效果**

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps34.jpg)

### **4.3 源代码**

```go
var st ST.StopTokens
func init() {
       st.Init(config.SPLIT_TOKEN_CN_EN_PATH)
}
func RemoveStopToken(words []string) []string{

       newWords := make([]string, 0)
       for _,word :=  range words {
              if !st.IsStopToken(word) {
                     newWords = append(newWords, word)
                     //fmt.Println(newWords)
              }
       }

       return newWords
}
package Stop_Token

import (
       "bufio"
       "log"
       "os"
)


type StopTokens struct {
       stopTokens map[string]bool
}

// 从stopTokenFile中读入停用词，一个词一行
// 文档索引建立时会跳过这些停用词
func (st *StopTokens) Init(stopTokenFile string) {
       st.stopTokens = make(map[string]bool)
       if stopTokenFile == "" {
              return
       }

       file, err := os.Open(stopTokenFile)
       if err != nil {
              log.Fatal(err)
       }
       defer file.Close()

       scanner := bufio.NewScanner(file)
       for scanner.Scan() {
              text := scanner.Text()
              //fmt.Println(text)
              if text != "" {
                     st.stopTokens[text] = true
                     //log.Println(text)
              }
       }

}

func (st *StopTokens) IsStopToken(token string) bool {
       //_, found := st.stopTokens[token]
       if st.stopTokens[token] {
              return true
       }
       return false
}
```



 

## 5. **删除中文停用词**

### **5.1设计思路**

给定停用词表，读取每一行停用词，使用go 的map数据结构构建一个set集合。具体代码为：map[string]bool。

传入分词结果，将分词中的每一个词语都与之匹配，如果存在就舍弃加入到新分词表。如果不存在则加入到其中。最后返回新的词表。最后log一下移除了多少个词语。

 

### **5.2 运行效果**

显示保存到的目标文件的数据中有多少个停用词与被删除：

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps35.jpg)

### **5.3 源代码**

```go
var st ST.StopTokens
func init() {
       st.Init(config.SPLIT_TOKEN_CN_EN_PATH)
}
func RemoveStopToken(words []string) []string{

       newWords := make([]string, 0)
       for _,word :=  range words {
              if !st.IsStopToken(word) {
                     newWords = append(newWords, word)
                     //fmt.Println(newWords)
              }
       }

       return newWords
}

package Stop_Token

import (
       "bufio"
       "log"
       "os"
)


type StopTokens struct {
       stopTokens map[string]bool
}

// 从stopTokenFile中读入停用词，一个词一行
// 文档索引建立时会跳过这些停用词
func (st *StopTokens) Init(stopTokenFile string) {
       st.stopTokens = make(map[string]bool)
       if stopTokenFile == "" {
              return
       }

       file, err := os.Open(stopTokenFile)
       if err != nil {
              log.Fatal(err)
       }
       defer file.Close()

       scanner := bufio.NewScanner(file)
       for scanner.Scan() {
              text := scanner.Text()
              //fmt.Println(text)
              if text != "" {
                     st.stopTokens[text] = true
                     //log.Println(text)
              }
       }

}

func (st *StopTokens) IsStopToken(token string) bool {
       //_, found := st.stopTokens[token]
       if st.stopTokens[token] {
              return true
       }
       return false
}
```



 

## 6.编程实现Porter Stemming功能

### **6.1设计思路**

使用CGO调用官方提供的C语言实现的Thread Safe版本的就可以了。

### **6.2 运行效果**

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps36.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps37.jpg)

 

### **6.3 源代码**



 

 

## 7. **将中文文档字符化 生成可被搜索引擎索引的字符单元**

### **7.1设计思路**

小明硕士毕业于中国科学院计算所，后在日本京都大学深造

搜索引擎模式: 小明/硕士/毕业/于/中国/科学/学院/科学院/中国科学院/计算/计算所/，/后/在/日本/京都/大学/日本京都大学/深造

### **7.2 运行效果**

 

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps38.jpg)

### **7.3 源代码**



 

## 8. **将英文文档预处理结果持久化**

### **8.1设计思路**

调用标准库中的文件操作，将文本内容，转为byte数组一次性写入到文件中即可。

### **8.2 运行效果**

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps39.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps40.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps41.jpg)



### **8.3 源代码**

```go
err = this.Save(content, config.NEWS_OLD_EN + "EN_READING_OLD_FILES_" + strconv.Itoa(int(config.COUNT)) + ".txt")
words := strings.Fields(content)
before := len(words)
words = RemoveStopToken(words)
after := len(words)
log.Printf("target : [%s] was removed %d stop words!...\n",page.Url, before - after)
this.Save(strings.Join(words, "\n"), config.SPILIT_WORDS_EN_PATH + "SPILIED_WORDS_" + strconv.Itoa(int(config.COUNT)) + ".txt")

// PorterStemming提取词干
stemWords := make([]string, len(words))
for idx, word := range words {

       stem := porterstemmer.Stem([]rune(word))

       stemWords[idx] = string(stem)

       fmt.Printf("The word [%s] has the stem [%s].\n", string(word), string(stem))
}
this.Save(strings.Join(stemWords, "\n"), config.NEW_FILE_EN + "SearchMode_New_Files_EN_" + strconv.Itoa(int(config.COUNT)) + ".txt")

// 保存到本地
func (this *Spider) Save(content string, filePath string) (error){
       //if (len(content) == 0) {
       //     return errors.New("内容为空")
       //}
       file, err := os.OpenFile(filePath, os.O_CREATE | os.O_WRONLY, 0666)
       if err != nil {
              fmt.Println("Open file err = ", err)
              return  err;
       }
       file.Write([]byte(content)) // 将字符串转为数组存放
       defer file.Close()
       return nil;
}
```



 

## 9. **将中文文档预处理结果持久化**

### **9.1设计思路**

调用标准库中的文件操作，将文本内容，转为byte数组一次性写入到文件中即可。

 

### **9.2 运行效果**

此处乱码是因为直接按照字节流的格式写入的。打开时要转码。

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps42.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps43.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps44.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps45.jpg)

### **9.3 源代码**

```go
err = this.Save(content, config.NEWS_OLD_CN + "SWJTU_NEWS_OLD_FILES_" + strconv.Itoa(int(config.COUNT)) + ".txt")
// 4. 将数据进行分词处理 并去除停用词
myjieba.SplitWords(content, config.SPILIT_WORDS_CN_PATH + "SPLITED_WORDS_" + strconv.Itoa(int(config.COUNT)) + ".txt")
// 5. 将数据进行去除停用词
//words = RemoveStopToken(words)
//this.Save(strings.Join(words, "\n"), config.SPILIT_WORDS_EN_PATH + "SPILIED_WORDS_" + strconv.Itoa(int(config.COUNT)) + ".txt")
// 6.将数据转换成为可以被搜索引擎模式
myjieba.SplitWords(content, config.NEWS_NEW_CN + "SearchMode_New_Files_CN_" + strconv.Itoa(int(config.COUNT)) + ".txt")

// 保存到本地
func (this *Spider) Save(content string, filePath string) (error){
       //if (len(content) == 0) {
       //     return errors.New("内容为空")
       //}
       file, err := os.OpenFile(filePath, os.O_CREATE | os.O_WRONLY, 0666)
       if err != nil {
              fmt.Println("Open file err = ", err)
              return  err;
       }
       file.Write([]byte(content)) // 将字符串转为数组存放
       defer file.Close()
       return nil;
}
```

# 性能测试

## 1.本机环境

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps47.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps48.jpg)

## 2.单元测试

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps49.jpg)

## 3.速度测试

请求交大新闻网500个页面从爬取数据到，预处理及持久化总共耗时32.986355s

![img](C:\Users\联想\AppData\Local\Temp\ksohtml18420\wps50.jpg) 

请求英语阅读网500个网页得到的结果：耗时4m14.9390929s。比较慢的原因：网站的服务器的性能不够好。除此之外还做了反爬机制。会封IP，让你TimeOut或者拒绝搁置你的请求等。增加耗时。故对于请求500个url的页面增加sleep(1s) 和每个任务开始sleep(0.5s)

实际时间应该是 = 4 * 60 + 15s  - 500 * 0.5 = 265 - 250 =15s，即如果在该英文阅读网服务器允许的情况下从爬取到文本预处理，总共耗时为15s就能完成。

![img](C:\Users\联想\AppData\Local\Temp\ksohtml18420\wps51.jpg)

## 4.性能测试

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps52.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps53.jpg)

## 5.覆盖率测试

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps54.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps55.jpg)

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps56.jpg)

## 6.CPU占用率测试

![img](https://gitee.com/DengSchoo374/img/raw/master/images/wps57.jpg)

# 参考资料及开源地址

## **参考资料：**

porter stemming 算法: https://tartarus.org/martin/PorterStemmer/

开源停用词：https://github.com/goto456/stopwords

中文分词：https://github.com/yanyiwu/gojieba

go-porter-stemmer，go语言版本的实现:https://github.com/reiver/go-porterstemmer

 

## **开源地址：**

我的github：https://github.com/DengSchoo/SearchEngineByGolang

