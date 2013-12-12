package main

import (
	"encoding/json"
	"fmt"
	"github.com/huichen/sego"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var SKIP = []string{".", "_", "static", "assert"}
var EXAM = []string{".html"}

type fcr struct {
	Count int
	Links map[string]bool
}

type fcRes map[string]interface{}

func main() {

	if len(os.Args) < 2 {
		log.Fatalln("usage: fc <dict> <dir>")
	}
	// 载入词典
	var segmenter sego.Segmenter
	segmenter.LoadDictionary(os.Args[1])
	// segmenter.LoadData()

	// 分词
	// text := []byte(`词典用前缀树实现， 分词器算法为基于词频的最短路径加动态规划。 支持普通和搜索引擎两种分词模式，支持用户词典、词性标注，可运行JSON RPC服务。 分词速度单线程2.7MB/s，goroutines并发13MB/s, 处理器Core i7-3615QM 2.30GHz 8核`) // segments := segmenter.Segment(text)

	index := make(map[string]fcr)

	root := os.Args[2]

	skip := func(fi os.FileInfo) bool {
		if fi.IsDir() {
			return true
		}
		name := fi.Name()
		for _, s := range SKIP {
			if strings.Index(name, s) == 0 {
				return true
			}
		}
		for _, e := range EXAM {
			if strings.Index(name, e) > 0 {
				return false
			}
		}
		return true
	}

	walker := func(fn string, fi os.FileInfo, err error) error {
		rel, _ := filepath.Rel(root, fn)
		if !skip(fi) {
			bs, err := ioutil.ReadFile(fn)
			if err != nil {
				return err
			}
			segments := segmenter.Segment(bs)
			for _, s := range segments {
				text := s.Token().Text()
				if len(text) > 1 {
					r, ok := index[text]
					if ok {
						r.Count = r.Count + 1
						r.Links[rel] = true
						index[text] = r
					} else {
						r := fcr{Count: 1, Links: make(map[string]bool)}
						r.Links[rel] = true
						index[text] = r
					}
				}
			}
		}
		return nil
	}

	filepath.Walk(root, walker)
	murls := make(map[string]int)
	for _, v := range index {
		for w, _ := range v.Links {
			murls[w] = 0
		}
	}
	aurls := make([]string, len(murls))
	i := 0
	for k, _ := range murls {
		murls[k] = i
		aurls[i] = k
		i++
	}

	type tmp struct {
		K string `json:"value"`
		L []int  `json:"data"`
	}
	tt := make([]tmp, len(index))
	j := 0
	for k, v := range index {
		ws := make([]int, len(v.Links))
		i := 0
		for w, _ := range v.Links {
			ws[i] = murls[w]
			i++
		}
		t := tmp{K: k, L: ws}
		tt[j] = t
		j++
	}

	jsn := make(map[string]interface{})
	jsn["links"] = aurls
	jsn["fc"] = tt
	fcs, err := json.Marshal(jsn)
	if err != nil {
		log.Fatalln("err:", err)
		os.Exit(-1)
	}
	fmt.Println("var fc = " + string(fcs))

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	// fmt.Println(sego.SegmentsToString(segments, false))

}
