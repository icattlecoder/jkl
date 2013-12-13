package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type menuItem struct {
	Title string     `omitempty`
	Href  string     `omitempty`
	Child []menuItem `omitempty`
	Order int
}

// By is the type of a "less" function that defines the ordering of its menuItem arguments.
type By func(p1, p2 *menuItem) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(menus []menuItem) {
	ps := &menuSorter{
		menus: menus,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// menuSorter joins a By function and a slice of menuItems to be sorted.
type menuSorter struct {
	menus []menuItem
	by    func(p1, p2 *menuItem) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *menuSorter) Len() int {
	return len(s.menus)
}

// Swap is part of sort.Interface.
func (s *menuSorter) Swap(i, j int) {
	s.menus[i], s.menus[j] = s.menus[j], s.menus[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *menuSorter) Less(i, j int) bool {
	return s.by(&s.menus[i], &s.menus[j])
}

const (
	nodeName = "fname.txt"
)

type jsonResult struct {
	Name  string
	Items []menuItem
}

func readPage(filename string) (page map[string]interface{}, err error) {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	page = map[string]interface{}{}
	err = goyaml.Unmarshal(c, &page)
	return
}

func skip2(d os.FileInfo) bool {
	skipfiles := []string{"img", "index.markdown", "index.md"}
	name := d.Name()
	for _, f := range skipfiles {
		if f == name {
			return true
		}
	}
	if d.IsDir() && name != "img" {
		return false
	}
	slash := filepath.ToSlash(name)
	if strings.Contains(slash, ".markdown") || strings.Contains(slash, ".md") {
		return false
	}
	return true
}

func skip(dir, name string) bool {
	skipfiles := []string{"img", "index.markdown", "index.md"}
	for _, f := range skipfiles {
		if f == name {
			return true
		}
	}
	slash := filepath.ToSlash(name)

	if strings.Contains(slash, ".") && (strings.Contains(slash, ".markdown") || strings.Contains(slash, ".md")) {
		return false
	} else {
		return true
	}
}

func rread(dir string) (items []menuItem) {

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		os.Exit(-1)
		fmt.Println("error=", err)
	}

	l := 0
	for _, d := range fs {
		if !skip2(d) {
			l++
		}
	}

	items = make([]menuItem, l) // make(menuItem,len(fs))

	i := 0
	for _, d := range fs {
		if skip2(d) {
			continue
		}
		n := path.Join(dir, d.Name())
		item := menuItem{}
		if d.IsDir() {
			page, err := readPage(path.Join(n, "index.markdown"))
            if err !=nil{
                page, err = readPage(path.Join(n, "index.md"))
            }
			if err != nil {
				cc, err := ioutil.ReadFile(path.Join(n, nodeName))
				if err != nil {
					item.Title = d.Name()
				} else {
					item.Title = string(cc)
				}

			} else {
				if str, ok := page["title"].(string); ok {
					item.Title = str
					item.Href = n + "/"
				}

				if order, ok := page["order"].(int); ok {
					item.Order = order
				}

			}
			item.Child = rread(n)
		} else if strings.Contains(n, ".markdown") || strings.Contains(n, ".md") {
			page, err := readPage(n)
			if err != nil {
				continue
			}
			if str, ok := page["title"].(string); ok {
				item.Title = str
			}
            if strings.Contains(n,".markdown"){
                item.Href = n[0:len(n)-8] + "html"
            }else{
                item.Href = n[0:len(n)-2] + "html"
            }
			if ord, ok := page["order"].(int); ok {
				item.Order = ord
			} else {
				item.Order = 0
			}

		} else {
			continue
		}
		items[i] = item
		i++
	}
	name := func(m1, m2 *menuItem) bool {
		return m1.Order > m2.Order
	}
	By(name).Sort(items)
	return
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println(` usage: genMenu <path>`)
		os.Exit(-1)
	}

	dir := os.Args[1]
	items := rread(dir)

	res := jsonResult{Name: dir, Items: items}

	bs, err := json.MarshalIndent(res, "", "\t")

	if err != nil {
		os.Exit(-1)
	}
	fmt.Println(string(bs))
}
