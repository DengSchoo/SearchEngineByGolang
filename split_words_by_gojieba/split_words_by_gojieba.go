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

func saveSpilitResult(content, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Open file err = ", err)
		return err
	}
	file.Write([]byte(content)) // 将字符串转为数组存放
	defer file.Close()
	return nil
}

var jieba = gojieba.NewJieba()

func SplitWords(content, path string) {
	var words []string

	words = jieba.Cut(content, true)
	//log.Printf("before remove %d,",len(words))
	before := len(words)
	words = RemoveStopToken(words)
	after := len(words)
	//log.Printf(" after remove %d,",len(words))
	log.Printf("target : [%s] was removed %d stop words!...\n", path, before-after)
	saveSpilitResult(strings.Join(words, "\n"), path)
	//return words;
}
func SplitWordsSearchMode(content, path string) {
	var words []string
	jieba := gojieba.NewJieba()
	words = jieba.CutAll(content)

	words = RemoveStopToken(words)
	saveSpilitResult(strings.Join(words, "\n"), path)
	words = nil
	//return words;
}
