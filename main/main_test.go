package main

import (
	"log"
	"regexp"
	"testing"
)

func TestBeginGet500CN(t *testing.T) {
	BeginGet500CN()
}
func BenchmarkBeginGet500CN(b *testing.B) {
	BeginGet500CN()
}

func BenchmarkBeginGet500EN(b *testing.B) {
	urls := getPageList()
	BeginGet500EN(urls)
}

func TestParse(t *testing.T) {

	parseResult := "邓世豪，。 123。 等12！ 12等Q22，。《》"
	reg := regexp.MustCompile(`[^a-zA-Z_ \\u4e00-\\u9fa5]`)
	parseResult = reg.ReplaceAllString(parseResult, "")
	log.Println(parseResult)
}
