package main

import (
	"flag"
	"fmt"
	"godsend/pkg/crawler"
	"godsend/pkg/crawler/spider"
	"log"
	"strings"
)

type search struct {
	scanner crawler.Interface
	sites   []string
	depth   int
	workers int
}

func main() {
	search := new()
	docs := search.scan()
	result(docs)
}

func new() *search {
	s := search{}
	s.scanner = spider.New()
	s.sites = []string{"https://go.dev", "https://golang.org/"}
	s.depth = 2
	s.workers = 10
	return &s
}

func (s *search) scan() []crawler.Document {
	log.Println("🚀 Scanning go sites... ✨")
	docs := []crawler.Document{}

	for _, site := range s.sites {
		d, err := s.scanner.Scan(site, 2)

		if err != nil {
			log.Println("Scan error", err)
			break
		}

		for _, i := range d {
			docs = append(docs, i)
		}
	}

	log.Println("Scan has been completed ✅")
	return docs
}

func result(docs []crawler.Document) {
	sFlag := flag.String("s", "", "search by word")
	flag.Parse()

	if *sFlag != "" {
		byWord(sFlag, docs)
		return
	}

	for i, item := range docs {
		fmt.Printf("Page %d. %s\nURL: %s\n\n", i+1, item.Title, item.URL)
	}
}

func byWord(w *string, docs []crawler.Document) {
	any := false
	n := 0

	for _, d := range docs {
		if strings.Contains(d.Title, *w) {
			if !any {
				any = true
				fmt.Printf("%s was found at:\n\n", *w)
			}
			n++
			fmt.Printf("%d. %s\n%s\n\n", n, d.Title, d.URL)
		}
	}

	if !any {
		fmt.Printf("There were no page containing '%s' found", *w)
	}
}
