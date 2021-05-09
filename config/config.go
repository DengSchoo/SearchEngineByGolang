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
	DOC_NUM    = 500
	DOC_NUM_EN = 500

	// 英文文档相关
	PREFIX_URL_EN = "http://www.enread.com/essays/"
	SUBFIX_URL_EN = ".html"
	CLASS_ID_EN   = "#dede_content > div"
	Base_URL_EN   = 112944
	NEW_FILE_EN   = "./files/New_Files_EN/"

	// files保存位置

	// CN
	NEWS_OLD_CN = "./files/Old_Files_CN/"
	NEWS_NEW_CN = "./files/New_Files_CN/"

	NEWS_OLD_EN = "./files/Old_Files_EN/"

	// 分词token路径
	SPLIT_TOKEN_CN_EN_PATH = "./files/stop_token_files_EN_CN/baidu_stopwords_en_cn.txt"
	SPLIT_TOKEN_CN_PATH    = "./files/stop_token_files_EN_CN/stopwords_cn.txt"
	SPLIT_TOKEN_EN_PATH    = "./files/stop_token_files_EN_CN/stopwords_en.txt"
	// 分词存储路径
	SPILIT_WORDS_CN_PATH = "./files/split_words_files_CN/"
	SPILIT_WORDS_EN_PATH = "./files/split_words_files_EN/"
)

var (
	COUNT uint
	// 用于控制COUNT++
	COUNT_LOCK sync.Mutex
)

func init() {
	COUNT = 1
}
