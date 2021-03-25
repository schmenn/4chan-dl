package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var c http.Client

func main() {
	var wg sync.WaitGroup
	res, err := http.Get(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}
	doc.Find("a.fileThumb").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func() {
			a, _ := s.Attr("href")
			fmt.Println(a)
			downloadandsave(a)
			wg.Done()
		}()
	})
	
	wg.Wait()
}

func downloadandsave(s string) {
	str := "https:" + s
	fname := strings.Split(str, "/")[len(strings.Split(str, "/"))-1]
	if _, err := os.Stat(fname); !os.IsNotExist(err) {
		return
	}
	res, err := c.Get(str)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	err = ioutil.WriteFile("./"+fname, body, 0666)
	if err != nil {
		panic(err)
	}
}
