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
	s.find()
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
			log.Println(i)
			docs = append(docs, i)
		}
	}
	return docs
}

func (s *search) find() {
	sFlag := flag.String("s", "", "search by word")
	flag.Parse()

	if *sFlag != "" {
		ids := s.index.Search(*sFlag)
		docs := s.index.Docs()

		fmt.Printf("\n'%s' was found in:\n\n", *sFlag)

		for _, id := range ids {
			fmt.Printf("%v\n", docs[id])
		}
	}
}
