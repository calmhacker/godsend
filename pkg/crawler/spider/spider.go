package spider

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"godsend/pkg/crawler"
)

type Service struct{}

func New() *Service {
	s := Service{}
	return &s
}

func (s *Service) BatchScan(urls []string, depth int, workers int) (<-chan crawler.Document, <-chan error) {
	chURLs := make(chan string)
	chOut := make(chan crawler.Document)
	chErr := make(chan error)
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for url := range chURLs {
				data, err := s.Scan(url, depth)
				if err != nil {
					log.Println("error:", err)
					chErr <- err
					return
				}
				for _, doc := range data {
					log.Println("scanned doc:", doc)
					chOut <- doc
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		log.Println("closing err and output channels")
		close(chErr)
		close(chOut)
	}()

	go func() {
		for _, url := range urls {
			chURLs <- url
		}
		log.Println("closing link channels")
		close(chURLs)
	}()

	return chOut, chErr
}

func (s *Service) Scan(url string, depth int) (data []crawler.Document, err error) {
	pages := make(map[string]string)

	parse(url, url, depth, pages)

	for url, title := range pages {
		item := crawler.Document{
			URL:   url,
			Title: title,
		}
		data = append(data, item)
	}

	return data, nil
}

func parse(url, baseurl string, depth int, data map[string]string) error {
	if depth == 0 {
		return nil
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	page, err := html.Parse(response.Body)
	if err != nil {
		return err
	}

	data[url] = pageTitle(page)

	if depth == 1 {
		return nil
	}
	links := pageLinks(nil, page)
	for _, link := range links {
		link = strings.TrimSuffix(link, "/")
		if strings.HasPrefix(link, "/") && len(link) > 1 {
			link = baseurl + link
		}
		if data[link] != "" {
			continue
		}
		if strings.HasPrefix(link, baseurl) {
			parse(link, baseurl, depth-1, data)
		}
	}

	return nil
}

func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

func pageLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if !sliceContains(links, a.Val) {
					links = append(links, a.Val)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = pageLinks(links, c)
	}
	return links
}

func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
