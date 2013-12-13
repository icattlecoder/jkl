package jkl

import (
	"encoding/json"
	"github.com/huichen/sego"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

type FC struct {
	dict string
	save string
	ok   bool
}

var ignore = []string{"我", "人", "有", "的", "和", "主", "产", "不", "为", "这", "工", "要", "在", "地", "一", "上",
	"是", "中", "国", "的", "经", "以", "发", "了", "民", "同"}

var IGNORE = make(map[string]bool)

func (f *FC) PlugeName() string {
	return "fc"
}

func (f *FC) _init() {
	for _, v := range ignore {
		IGNORE[v] = true
	}
}

func (f *FC) SetArgs(args []string) {
	if f.ok = len(args) > 1; f.ok {
		f.dict = args[0]
		f.save = args[1]
	}
}

func (f *FC) OnPageParsed(page *Page) {

}

func (f *FC) OnSiteGenarated(site *Site) {
	if !f.ok {
		return
	}
	f.fc(site)
}

type fcr struct {
	Count int
	Links map[string]bool
}

type fcRes map[string]interface{}

func (f *FC) fc(site *Site) {

	// 载入词典
	var segmenter sego.Segmenter
	segmenter.LoadDictionary(f.dict)
	// segmenter.LoadData()

	index := make(map[string]fcr)

	root := "_site"

	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	//去除连续的换行符
	delbr, _ := regexp.Compile("\\s{2,}")

	walker := func(fn string, bs []byte) {
		rel, _ := filepath.Rel(root, fn)
		segments := segmenter.Segment(bs)
		for _, s := range segments {
			text := s.Token().Text()
			if _, ok := IGNORE[text]; ok {
				continue
			}
			// log.Println(len([]rune(text)))
			if len([]rune(text)) > 1 {
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

	for _, v := range site.PP {
		for _, p := range v {
			text := p.GetContent()
			text = re.ReplaceAllString(text, "")
			text = delbr.ReplaceAllString(text, "\n")
			walker(p.GetUrl(), []byte(p.GetContent()))
		}
	}

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
	type context struct {
		Context string `json:"c"`
		LinkID  int    `json:"l"`
	}

	type tmp struct {
		K string    `json:"value"`
		L []context `json:"data"`
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
		return
	}
	content := "var fc =" + string(fcs)
	ioutil.WriteFile(f.save, []byte(content), 0777)
	site.files = append(site.files, f.save)
	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	// fmt.Println(sego.SegmentsToString(segments, false))

}
