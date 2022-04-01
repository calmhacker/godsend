package main

import (
	"flag"
	"fmt"
	"godsend/pkg/crawler"
	"godsend/pkg/crawler/spider"
	"godsend/pkg/index"
	"godsend/pkg/index/hash"
	"log"
)

type search struct {
	scanner crawler.Interface
	sites   []string
	index   index.Interface
}

func main() {
	search := new()
	search.run()
}

func new() *search {
	s := search{}
	s.scanner = spider.New()
	s.sites = []string{"https://go.dev", "https://golang.org/"}
	s.index = hash.New()
	return &s
}

func (s *search) run() {
	docs := s.scan()
	s.index.Add(docs)

	sFlag := flag.String("s", "", "search by word")
	flag.Parse()

	// s.find(*sFlag)
	s.bsearchFind(*sFlag)
}

func (s *search) scan() []crawler.Document {
	log.Println("🚀 Scanning go sites... ✨")
	docs := []crawler.Document{}

	for _, site := range s.sites {
		d, err := s.scanner.Scan(site, 2)

		if err != nil {
			log.Println("Scan error", err)
			continue
		}

		for _, i := range d {
			docs = append(docs, i)
		}
	}
	return docs
}

func (s *search) find(str string) {
	if str != "" {
		ids := s.index.Search(str)
		docs := s.index.Docs()

		fmt.Printf("\n'%s' was found in:\n\n", str)

		for _, id := range ids {
			fmt.Printf("%v\n", docs[id])
		}
	}
}

func (s *search) bsearchFind(str string) {
	if str != "" {
		ids := s.index.Search(str)
		docs := s.index.Docs()

		fmt.Printf("\n'%s' was found in:\n\n", str)

		for _, id := range ids {
			if bsearch(id, docs) {
				fmt.Printf("%v\n", docs[id])
			}
		}
	}
}

func bsearch(s int, arr []crawler.Document) bool {
	low := 0
	high := len(arr) - 1

	for low <= high {
		median := (low + high) / 2

		if arr[median].ID < s {
			low = median + 1
			continue
		}

		high = median - 1
	}

	if low == len(arr) || arr[low].ID != s {
		return false
	}

	return true
}
